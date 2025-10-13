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
