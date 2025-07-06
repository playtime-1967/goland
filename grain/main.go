package main

import (
	"grain/handlers"
	"grain/repositories"
	"grain/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func main() {
	RunSamples()
	return
	//Connect to Cassandra
	cluster := gocql.NewCluster("127.0.0.1") // Replace with your Cassandra IP
	cluster.Keyspace = "grainapp"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra!:", err)
	}

	defer session.Close()

	r := gin.Default()

	repo := repositories.NewUserRepository(session)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)
	r.POST("/users", handler.Create)
	r.GET("/users/:id", handler.Get)
	r.Run(":8080")
}
