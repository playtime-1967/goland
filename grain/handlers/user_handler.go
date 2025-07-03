package handlers

import (
	"grain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(c *gin.Context) {

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Create(input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// func CreateUser(session *gocql.Session) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user entities.User
// 		if err := c.ShouldBindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		user.ID = gocql.TimeUUID()

// 		if err := session.Query(`INSERT INTO users (id, name, email) VALUES (?, ?, ?)`,
// 			user.ID, user.Name, user.Email).Exec(); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusCreated, user)
// 	}
// }

// func GetUser(session *gocql.Session) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		idStr := c.Param("id")
// 		id, err := gocql.ParseUUID(idStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
// 			return
// 		}

// 		var user entities.User
// 		if err := session.Query(`SELECT id, name, email FROM users WHERE id = ?`, id).
// 			Scan(&user.ID, &user.Name, &user.Email); err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, user)
// 	}
// }
