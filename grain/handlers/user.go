package handlers

import (
	"grain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func CreateUser(session *gocql.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user entities.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.ID = gocql.TimeUUID()

		if err := session.Query(`INSERT INTO users (id, name, email) VALUES (?, ?, ?)`,
			user.ID, user.Name, user.Email).Exec(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func GetUser(session *gocql.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		var user entities.User
		if err := session.Query(`SELECT id, name, email FROM users WHERE id = ?`, id).
			Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
