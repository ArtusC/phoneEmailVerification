package main

import (
	"context"
	"fmt"

	// collector "github.com/ArtusC/phoneEmailVerification/api"
	api "github.com/ArtusC/phoneEmailVerification/api"
	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	phoneNumberUseCase "github.com/ArtusC/phoneEmailVerification/usecases/phoneNumber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoSession mongo.Session
)

func init() {
	mongoUrl := "mongodb://root:root@localhost:27018"
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	fmt.Println("MongoDB connection established!")
	if mongoSession, err = client.StartSession(); err != nil {
		panic(err.Error())
	}
}

func main() {
	fmt.Println("Starting aplication")

	// Starts mongo repository.
	mongoRepo := repository.NewMongoRepository(mongoSession)
	defer mongoSession.EndSession(context.TODO())

	phoneUseCases := phoneNumberUseCase.NewPhoneUseCases(mongoRepo)

	api := api.NewApi(phoneUseCases)

	if err := api.StartServer(); err != nil {
		panic(fmt.Sprintf("error to start server due to %s", err.Error()))
	}

	// mongoDB, err := mongodb.CreateMongoClient(mongoURL)
	// if err != nil {
	// 	panic(fmt.Sprintf("error to connect with mongoDB due to %s", err.Error()))
	// }

	// LIDAR COM ESSA PARTE DAS ROTAS
	// e := echo.New()
	// e.POST("/api/phone", createPhone)
	// http.HandleFunc("/api/phone", mongoRepo.CreateRecord())
	// http.ListenAndServe(":8080", nil)

}
