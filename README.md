# rest2go

<p align="center">
    <img src="./rest2go.png">
</p>

**rest2go** is a lightweight starter for building REST API microservices in Go. Library provides complete foundation 
(from application configuration to database connection) so you can focus on writing business logic right away.

Library was created to address internal needs od **Ternary Software Solutions**. After developing several microservices 
written in almost pure Go, growing amount of duplicated boilerplate code was noticed - configuration loading, HTTP setup, 
database integration and errors handling. Each new project required additional time just to prepare repository and basic 
code infrastructure. **rest2go** was built to solve that problem by providing reusable, consistent starting point for 
all future services.

**Included out of the box**:

- Microservice configuration via **YAML**
- Basic ready to use **HTTP server**
- Set of core **middlewares**
- Authorization via **Api-Key header**
- Consistent **API error handling**
- **Database connection** and **migrations**
- Built in **filtering** and pagination **support**

**Requirements**:

- GoLang >= 1.25.0

**Table of contents**:

1. [Settings](#Settings)
2. [Logs](#Logs)
3. [HTTP server](#HTTP-server)
4. [Middlewares](#Middlewares)
5. [Errors handling](#Errors-handling)
6. [Database connection](#Database-connection)
7. [Database migrations](#Database-migrations)

## Settings

`rest2go` provides settings functionality for application using library. Idea behind this is to have single YAML file 
with application configuration. Library components are also configurable, so it provides basic settings structure. This 
structure can be easily extended in the application with `embed`. First of all, configuration is read from two localizations:

1. `app.yml` - first try, file in application root directory
2. `configs/app.yml` - second try, directory in application root directory

If configuration file does not exist, defaults will be loaded.

Settings provided by `rest2go` with their defaults:

```yaml
# Logs configuration (slog)
logs:
  # Logging level (info/debug/warn/error)
  level: "info"
  # Logging to rotable file enabled?
  file-enabled: false
  # Rotable logs max file size (MB)
  max-size: 10
  # Rotable logs file max age (days)
  max-age: 7

# HTTP server configuration
server:
  # HTTP server host
  host: "0.0.0.0"
  # HTTP server port
  port: 8080
  # Indicates if global HTTP 404 should be handled by rest2go errors handler
  not-found-handler: false

# Authorization configuration
authorization:
  # Authorization by Api-Key header
  header:
    # Api-Key header authorization enabled?
    enabled: false
    # Authorization key value (header value)
    key: ""
    # List of ant patterns for public API
    public:
      - "[ant_pattern]" 

# Database connection configuration
database:
  # Database driver (sqlite3/postgres)
  driver: "sqlite3"
  # Database host
  host: "./data/app.db"
  # Database port
  port: 0
  # Database username
  user: ""
  # Database password
  password: ""
  # Database name
  name: ""
  # Database schema
  schema: ""
```

Application using library can have its own configuration. For proper work with configurable library components it should 
be handled with `embed`:

```go
type AppSettings struct {
	settings.Settings `yaml:",inline"`
	Foo               string `yaml:"foo"`
}
```

In snippet above, `foo` represents application configuration. Then default configuration should be provided (what values 
to use, if YAML file do not exist):

```go

func (s *AppSettings) SetDefaults() {
	s.Settings.SetDefaults() // Here, defaults for library components are set
	s.Foo = "bar" // Here, defaults for application are set
}
```

Additionally, if you do not like how `slog` is configured, implement interface method `ConfigureLogs` for `AppSettings` 
to handle logs configuration:

```go
type Defaults interface {
	ConfigureLogs()
}

func (s *AppSettings) ConfigureLogs() {
  // Here you can configure logs by yourself or make it empty to omit library defaults
}
```

Finally, application settings can be loaded:

```go
settings, err := settings.Load[settings.Settings]()
```

Loading uses generics, so when `AppSettings` are used, just provide proper type. `Load` works as synchronized singleton 
(configuration can be load once to immutable struct).

## Logs 

`rest2go` provides default logging configuration with usage of `slog` and `lumberjack`. Application using library 
could log to console or console & rotable file. Configuration & how to change default behaviour is described in 
[Settings](#Settings) chapter.

## HTTP server

`rest2go` provides preconfigured HTTP server. Configuration in details is described in [Settings](#Settings) chapter. 
By default, HTTP server needs some configuration, routes and set of middlewares. Example below shows how to use 
preconfigured server.

```go
settings, err := settings.Load[settings.Settings]()

if err != nil {
  // Handle error
}

router := http.NewServeMux()
router.HandleFunc("GET /", [handle_func])
server := rest2go.NewServer(settings.Server, router, LogRequestAndResponseMiddleware, ApiKeyAuthMiddleware)

if err := server.Run(); err != nil {
  // Handle error
}
```

Passing middlewares is optional. Library also gives ability to handle HTTP 404 Not Found error globally with usage 
of [Errors handling](#Errors-handling).

## Middlewares

`rest2go` provides set of basic middlewares and ability to combine them in middlewares chain. Idea is simple - execute 
some chunks of code before HTTP request is handled by business logic. Snippet below shows how to create middlewares 
chain.

```go
router := http.NewServeMux()
server := &http.Server{
  Addr:    addr,
  Handler: Middlewares(<m1>, <m2>, ..., <mx>)(router),
}
```

### Log request & response

`LogRequestAndResponseMiddleware` is created to log details about incoming HTTP request and out coming response. 
`slog` is used to print HTTP method and path. Body is logged only for `application/json` content type. Additionally 
there is UUID that indicates that given logs was written for single REST API call.

### ApiKeyAuthMiddleware

`ApiKeyAuthMiddleware` is created to provide basic authorization mechanism for microservice. If enabled, every HTTP 
request is checked if there was `Api-Key` header with proper secret key value. Additionally, it can be configured 
with ant patterns to make some API public. Detailed configuration is described in [Settings](#Settings) chapter.

## Errors handling

`rest2go` provides standard for REST API errors handling. All errors created during request handling are processed by 
library, any other error not connected to the request is handled by GoLang language layer. All REST API errors will be 
returned in standardized structure:

```json
{
    "timestamp": "2024-06-28T16:47:23Z",
    "code": "err.invalid-request",
    "message": "Invalid Request",
    "details": [
        {
            "field": "y",
            "code": "val.division_by_zero",
            "message": "Division by zero",
            "value": "0.0",
            "expected": "y != 0"
        }
    ]
}
```

To handle error, just use `rest2go.HandleError(err, response)`. Snippet below shows, how to create request validation 
errors that can be handled:

```go
func validate(veh VehicleRequestDto) error {
  vehTypes := []string{TypeCar, TypeMoto}
  errors := []rest2go.FieldError{}

  if strings.TrimSpace(veh.Type) != "" {
    if !slices.Contains(vehTypes, veh.Type) {
      errors = append(errors, rest2go.NewDetailedFieldError("type", "val.invalid", "Invalid type", veh.Type, strings.Join(vehTypes, ",")))
    }
  } else {
    errors = append(errors, rest2go.NewFieldError("type", "val.required", "Type is required"))
  }

  if len(errors) != 0 {
    return rest2go.NewApiError(400, "vehicle validation failed", errors...)
  }

  return nil
}
```

There are two types of field errors:

1. `NewFieldError` - basic error with field indication
2. `NewDetailedFieldError` - extended error for field with given value and what is expected

**WARNING**: errors handling mechanism uses `slog` by default.

## Database connection

`rest2go` provides utilities to initialize database connection. All available settings are described in 
[Settings](#Settings) chapter. Supported drivers are described below. Idea behind this is to provide database driver 
from parent application (`blank import` in `main.go`) and use it with set of tools provided by library. First of all, 
there is database connection provider that works as singleton:

```go
settings, err := settings.Load[settings.Settings]()

if err != nil {
  // Handle error
}

provider, err := rest2go.NewDbProvider(settings.Database)

if err != nil {
  // Handle error
}

defer provider.CloseConnection()
```

After initialization, stores can be created. Every store should implement interface visible below:

```go
// Interface
type DbStore interface {
	Begin() (*DbCtx, error)
	Commit(context *DbCtx) error
	Rollback(context *DbCtx) error
}

// Implementation
type vehiclesStore struct {
  db *sql.DB
}

func NewVehiclesStore(db *sql.DB) *vehiclesStore {
  return &vehiclesStore{
    db: db,
  }
}

func (s *vehiclesStore) Begin() (*rest2go.DbCtx, error) {
  tx, err := s.db.Begin()

  if err != nil {
    return nil, err
  }

  return rest2go.NewDbContext(tx), nil
}

func (s *vehiclesStore) Commit(context *rest2go.DbCtx) error {
  if err := context.Tx.Commit(); err != nil {
    return err
  }

  return nil
}

func (s *vehiclesStore) Rollback(context *rest2go.DbCtx) error {
  if err := context.Tx.Rollback(); err != nil {
    return err
  }

  return nil
}
```

`db` for store can be obtained from `provider` with call `provider.Db()`. By default, library is configured to handle 
SQLite database that exists in `./data/app.db`.

### SQLite

By default, library is configured to handle SQLite database. Connection can be configured as follows:

```yml
database:
  driver: "sqlite3"
  host: "./data/app.db" # Path to database file (given example is library default)
```

Additionally, connection is automatically prepared to respect foreign keys - `PRAGMA FOREIGN_KEYS=ON;`.

### Postgres

Library supports Postgres databases. Connection can be configured as follows:

```yml
database:
  driver: "postgres"
  host: "" # Database URL
  port: 5432 # Database port
  user: "" # Database username
  password: "" # Database password
  name: "" # Database name
  schema: "public" # Database schema
```

## Database migrations

`rest2go` is ready to handle database migrations with usage of [Goose](https://github.com/pressly/goose). All you need 
to do, is to place SQL migrations in `./migrations/${driver}` directory. Driver configured according to [Settings](#Settings) 
chapter also indicates what `Goose` should use. Snippet below shows how to migrate database within application.

```go

settings, err := settings.Load[settings.Settings]()

if err != nil {
  // Handle error
}

provider, err := rest2go.NewDbProvider(settings.Database)

if err != nil {
  // Handle error
}

defer provider.CloseConnection()

if err := provider.MigrateDatabase(); err != nil {
  // Handle error
}
```
