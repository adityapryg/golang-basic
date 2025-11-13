package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=1,max=150"`
}

var (
	users   = make(map[int]User)
	nextID  = 1
	usersMu sync.RWMutex
)

func init() {
	users[1] = User{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30}
	users[2] = User{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25}
	users[3] = User{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Age: 35}
	nextID = 4
}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Printf("[%s] %s %s\n", t.Format("15:04:05"), c.Request.Method, c.Request.URL.Path)
		c.Next()
		latency := time.Since(t)
		fmt.Printf("  Status: %d - Latency: %v\n", c.Writer.Status(), latency)
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or missing token",
			})
			return
		}
		c.Next()
	}
}

func main() {
	port := flag.Int("port", 8080, "Server port")
	host := flag.String("host", "localhost", "Server host")
	debug := flag.Bool("debug", true, "Debug mode")
	flag.Parse()

	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(CustomLogger())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Framework API",
			"version": "1.0.0",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUsers)
		v1.GET("/users/:id", getUser)

		protected := v1.Group("")
		protected.Use(AuthRequired())
		{
			protected.POST("/users", createUser)
			protected.PUT("/users/:id", updateUser)
			protected.DELETE("/users/:id", deleteUser)
		}
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("\n🚀 Server running on http://%s\n", addr)
	fmt.Println("📚 Docs: GET /api/v1/users")
	fmt.Println("🔐 Auth: Authorization: Bearer secret-token")

	r.Run(addr)
}

func getUsers(c *gin.Context) {
	usersMu.RLock()
	defer usersMu.RUnlock()

	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(userList),
		"data":    userList,
	})
}

func getUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	usersMu.RLock()
	user, exists := users[id]
	usersMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersMu.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = user
	usersMu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created",
		"data":    user,
	})
}

func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	if _, exists := users[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	users[id] = user

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated",
		"data":    user,
	})
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	if _, exists := users[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	delete(users, id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted",
	})
}
