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

**Table of contents::

1. [Settings](#Settings)
2. [Logs](#Logs)

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
