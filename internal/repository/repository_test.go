package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tp "github.com/ArtusC/phoneEmailVerification/types"
)

type fixture struct {
	logger          zerolog.Logger
	mongoSession    mongo.Session
	mongoRepository MongoRepository
}

var (
	mongoSession mongo.Session
)

const (
	mongoUrl       = "mongodb://root:root@localhost:27018"
	dbName         = "phoneDb"
	collectionName = "phone-collection"
)

func setUp() *fixture {

	log := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	//TODO: change the docker-compose that will up the database (to not sobrescribe), dont put login/password
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Panic().Msgf("[MongoDBTest] Error to start the client: %s", err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic().Msgf("[MongoDBTest] Errot to ping the client: %s", err.Error())
	}

	log.Info().Msg("[MongoDBTest] Connection established!")

	if mongoSession, err = client.StartSession(); err != nil {
		log.Panic().Msgf("[MongoDBTest] Error to start session: %s", err.Error())
	}

	repository := NewMongoRepository(log, mongoSession)
	log.Info().Msg("[MongoDBTest] Instantiate the Mongo repository")

	return &fixture{
		logger:          log,
		mongoSession:    mongoSession,
		mongoRepository: repository,
	}
}

func (f *fixture) tearDown() {
	f.logger.Info().Msg("[MongoDBTest] Connection finished!")
	d := f.mongoSession.Client().Database(dbName)
	d.Drop(context.Background())
	defer f.mongoSession.EndSession(context.Background())
}

func TestPrepareUpdateOperation(t *testing.T) {
	f := setUp()
	defer f.tearDown()

	v1 := tp.Country{
		Name:    "Brazil",
		M49Code: 840,
		IsoAdminLanguages: tp.IsoAdminLanguages{
			{NativeName: "v11", IsoName: "v12", IsoAlpha3: "v13", IsoAlpha2: "v14"},
		},
	}

	v2 := tp.Country{
		Name:    "Brazil",
		M49Code: 841,
		IsoAdminLanguages: tp.IsoAdminLanguages{
			{NativeName: "v21", IsoName: "v22", IsoAlpha2: "v14"},
		},
	}

	// Call the function you want to test
	result := PrepareUpdateOperation(v1, v2)

	assert.NotNil(t, result)
}
