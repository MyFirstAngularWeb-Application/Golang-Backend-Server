// config/config.go
package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://sachinayeshmantha:MR3eWG8SVIJuLCOP@gocluster1.ctslqzm.mongodb.net/?retryWrites=true&w=majority&appName=goCluster1")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("devops_community") // Replace with your database name
	fmt.Println("Connected to MongoDB!")
}
