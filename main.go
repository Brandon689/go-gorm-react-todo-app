package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Todo model
type Todo struct {
	ID     string `gorm:"type:uuid;primary_key" json:"id"`
	Text   string    `json:"text"`
	Status string    `json:"status"`
	gorm.Model
}

var db *gorm.DB

func init() {
	// Connect to SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// AutoMigrate the Todo model
	db.AutoMigrate(&Todo{})
}

// Handlers
func getAllTodos(c echo.Context) error {
	var todos []Todo
	db.Find(&todos)
	return c.JSON(http.StatusOK, todos)
}



func createTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}

	// Generate a new UUID
	uuid, err := uuid.NewRandom()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate UUID"})
	}

	todo.ID = uuid.String()

	db.Create(&todo)
	return c.JSON(http.StatusCreated, todo)
}

func getTodoByID(c echo.Context) error {
	id := c.Param("id")
	var todo Todo
	if err := db.First(&todo, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, todo)
}

func updateTodoByID(c echo.Context) error {
	id := c.Param("id")
	var todo Todo
	if err := db.First(&todo, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	newTodo := new(Todo)
	if err := c.Bind(newTodo); err != nil {
		return err
	}

	todo.Text = newTodo.Text
	todo.Status = newTodo.Status

	db.Save(&todo)
	return c.JSON(http.StatusOK, todo)
}

func deleteTodoByID(c echo.Context) error {
	id := c.Param("id")
	var todo Todo
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db.Delete(&todo)
	return c.NoContent(http.StatusNoContent)
}


func main() {
	e := echo.New()

	// Enable CORS
	//e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	
	e.Static("/", "my-react-vite-app/dist")

	// Routes
	e.GET("/todos", getAllTodos)
	e.POST("/todos", createTodo)
	e.GET("/todos/:id", getTodoByID)
	e.PUT("/todos/:id", updateTodoByID)
	e.DELETE("/todos/:id", deleteTodoByID)

	port := 8000
	fmt.Printf("Server is running on :%d...\n", port)
	e.Start(fmt.Sprintf(":%d", port))
}
