package main

import (
	"context"
	"os"
	"time"

	api "github.com/ArtusC/phoneEmailVerification/api"
	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	phoneNumberUseCase "github.com/ArtusC/phoneEmailVerification/usecases/phoneNumber"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rs/zerolog"
)

var (
	mongoSession mongo.Session
	logger       zerolog.Logger
	api_bdc_key  string
)

func init() {

	api_bdc_key = os.Getenv("API_BDC_KEY")
	if api_bdc_key == "" {
		panic("Please, create and export on the shell the necessary KEY to get data from BDC API (more details on README)")
	}

	logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	mongoUrl := "mongodb://root:root@localhost:27018"
	logger.Info().Msgf("[MongoDb] Starting connection at %s!\n", mongoUrl)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		logger.Panic().Msgf("[MongoDb] Error to start the client: %s", err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Panic().Msgf("[MongoDB] Errot to ping the client: %s", err.Error())
	}

	logger.Info().Msg("[MongoDB] Connection established!")
	if mongoSession, err = client.StartSession(); err != nil {
		logger.Panic().Msgf("[MongoDB] Error to start session: %s", err.Error())
	}
}

// Generate a unique correlation ID
func traceID() string {
	return uuid.New().String()
}

func main() {

	logger.Info().Msg("Starting aplication!")

	logger.Info().Msg("Starting mongo repository.")

	logger = logger.With().Str("traceId", traceID()).Caller().Logger()

	mongoRepo := repository.NewMongoRepository(logger, mongoSession)
	defer mongoSession.EndSession(context.TODO())

	logger.Info().Msg("Instatiating phone number use case.")
	phoneUseCases := phoneNumberUseCase.NewPhoneUseCases(logger, mongoRepo, api_bdc_key)

	logger.Info().Msg("Starting API.")
	api := api.NewApi(logger, phoneUseCases)

	if err := api.StartServer(); err != nil {
		logger.Fatal().Msgf("error to start server due to %s", err.Error())
	}

}
