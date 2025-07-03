package main

import (
	"grain/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func main() {
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

	r.POST("/users", handlers.CreateUser(session))
	r.GET("/users/:id", handlers.GetUser(session))

	r.Run(":8080")
}
