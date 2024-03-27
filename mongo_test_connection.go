package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoUrl := "mongodb://root:root@localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		fmt.Println(err.Error())
	}

	collection := client.Database("phoneDb").Collection("phone-collection")

	fmt.Println(collection.Name())

	data := map[string]interface{}{
		"ola": "mundo",
	}

	req, err := collection.InsertOne(ctx, data)

	if err != nil {
		fmt.Println(err.Error())
	}

	insertedId := req.InsertedID

	res := map[string]interface{}{
		"data": map[string]interface{}{
			"insertedId": insertedId,
		},
	}

	fmt.Printf("RES: ", res)
}
