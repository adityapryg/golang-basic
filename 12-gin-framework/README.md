# Gin Framework Documentation

## Overview
Gin is a web framework written in Go that simplifies the process of building web applications. It is known for its speed and productivity.

## Installation
To install Gin, use the following command:
```
go get github.com/gin-gonic/gin
```

## Basic Routing
Gin allows you to define routes easily. Here’s a simple example:
```go
r := gin.Default()
r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Hello, World!"})
})
```

## URL Parameters
You can bind URL parameters to variables:
```go
r.GET("/user/:name", func(c *gin.Context) {
    name := c.Param("name")
    c.JSON(200, gin.H{"user": name})
})
```

## Request Binding with Validation Tags
You can easily bind and validate request data:
```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

var user User
if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

## Response Handling with JSON/XML/YAML
You can handle different response formats:
```go
// JSON
c.JSON(200, gin.H{"message": "Success!"})

// XML
c.XML(200, gin.H{"message": "Success!"})

// YAML
c.YAML(200, gin.H{"message": "Success!"})
```

## Middleware
### Built-in Middleware
Gin comes with several built-in middleware, like logging and recovery:
```go
r.Use(gin.Logger())
r.Use(gin.Recovery())
```

### Custom Middleware
You can also create custom middleware:
```go
r.Use(func(c *gin.Context) {
    // Before request
    c.Next() // Proceed to the next middleware/handler
    // After request
})
```

### CORS Middleware Example
To handle CORS, you can use the following:
```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins: true,
}))
```

## Route Grouping
You can group routes to organize your code:
```go
v1 := r.Group("/v1") {
    v1.GET("/users", getUsers)
}
```

## Comparison between net/http and Gin
Using Gin results in cleaner and more readable code:
```go
// net/http
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
})

// Gin
r.GET("/", func(c *gin.Context) {
    c.String(200, "Hello, World!")
})
```

## Best Practices
### DO's
- Keep handlers small and focused.
- Use middleware to handle repetitive tasks.

### DON'Ts
- Don't ignore error handling.
- Avoid complex logic in handlers.

## Examples of GET/POST/PUT/DELETE Handlers
```go
// GET
r.GET("/users", getUsers)
// POST
r.POST("/users", createUser)
// PUT
r.PUT("/users/:id", updateUser)
// DELETE
r.DELETE("/users/:id", deleteUser)
```

## Authentication Middleware Example
```go
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check authentication
        c.Next()
    }
}
```

## Testing with curl Commands
You can test your API with curl:
```sh
# GET
curl -X GET http://localhost:8080/users
# POST
curl -X POST http://localhost:8080/users -d '{"name":"John"}' -H 'Content-Type: application/json'
```

## Conclusion
Gin is a powerful framework for building web applications in Go, providing developers with tools to create complex systems with minimal overhead.
