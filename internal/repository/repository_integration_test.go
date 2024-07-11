//go:build unit || integration
// +build unit integration

package repository_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	tp "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// arrumar aqui
type fixture struct {
	logger          zerolog.Logger
	mongoSession    mongo.Session
	mongoRepository repository.MongoRepository
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

	repository := repository.NewMongoRepository(log, mongoSession)
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

// go test -v -count=1 -covermode=atomic -tags integration ./internal/repository -run ^TestMongoRespository_StoragePhoneRecord$
func TestMongoRespository_StoragePhoneRecord(t *testing.T) {
	f := setUp()
	defer f.tearDown()

	testCase := []struct {
		testName      string
		phoneOutput   tp.PhoneNumber
		expectedError error
	}{
		{
			testName:      "save record correctly",
			phoneOutput:   tp.TestPhoneValue,
			expectedError: nil,
		},
	}

	for _, test := range testCase {
		t.Run(test.testName, func(t *testing.T) {

			err := f.mongoRepository.StoragePhoneRecord(f.logger, test.phoneOutput, dbName, collectionName)
			if err != nil {
				_id := repository.GetMD5Hash(test.phoneOutput.PhoneInput)
				errMsg := fmt.Sprintf(`write exception: write errors: [E11000 duplicate key error collection: phoneDb.phone-collection index: _id_ dup key: { _id: "%s" }]`, _id)
				test.expectedError = errors.New(errMsg)

				assert.EqualErrorf(t, err, test.expectedError.Error(), "Should be equal!")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// go test -v -count=1 -covermode=atomic -tags integration ./internal/repository -run ^TestMongoRespository_GetPhoneRecord$
func TestMongoRespository_GetPhoneRecord(t *testing.T) {
	f := setUp()
	defer f.tearDown()

	testCases := []struct {
		testName      string
		phoneInput    string
		expectedError error
	}{
		{
			testName:      "get existing phone record",
			phoneInput:    tp.TestPhoneValue.PhoneInput,
			expectedError: nil,
		},
		{
			testName:      "get non-existing phone record",
			phoneInput:    "non-existing-phone",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			data := tp.TestPhoneValue
			if tc.testName == "get existing phone record" {
				err := f.mongoRepository.StoragePhoneRecord(f.logger, data, dbName, collectionName)
				assert.Nil(t, err)
			}

			res, err := f.mongoRepository.GetPhone(f.logger, dbName, collectionName, tc.phoneInput)
			assert.Equal(t, tc.expectedError, err)

			if err == nil && res.ID != "" {
				assert.Contains(t, fmt.Sprint(res.ID), repository.GetMD5Hash(data.PhoneInput))
				assert.Contains(t, fmt.Sprint(res.E164Format), data.PhoneInput)
			} else {
				assert.Empty(t, res)
			}
		})
	}
}

// go test -v -count=1 -covermode=atomic -tags integration ./internal/repository -run ^TestMongoRespository_GetAllPhoneRecords$
func TestMongoRespository_GetAllPhoneRecordsOneRecord(t *testing.T) {
	f := setUp()
	defer f.tearDown()

	data := tp.TestPhoneValue

	err := f.mongoRepository.StoragePhoneRecord(f.logger, data, dbName, collectionName)
	assert.Nil(t, err)

	res, err := f.mongoRepository.GetAllPhoneRecords(f.logger, dbName, collectionName)
	assert.Nil(t, err)

	fmt.Println("Result: ", res)

	assert.Contains(t, fmt.Sprint(res[0].ID), repository.GetMD5Hash(data.PhoneInput))
	assert.Contains(t, fmt.Sprint(res[0].E164Format), data.PhoneInput)

}

// go test -v -count=1 -covermode=atomic -tags integration ./internal/repository -run ^TestMongoRespository_GetAllPhoneRecordsMoreThanOneRecord$
func TestMongoRespository_GetAllPhoneRecordsMoreThanOneRecord(t *testing.T) {
	f := setUp()
	defer f.tearDown()

	data1 := tp.TestPhoneValue
	data2 := tp.TestPhoneValue_2

	err := f.mongoRepository.StoragePhoneRecord(f.logger, data1, dbName, collectionName)
	assert.Nil(t, err)

	err = f.mongoRepository.StoragePhoneRecord(f.logger, data2, dbName, collectionName)
	assert.Nil(t, err)

	res, err := f.mongoRepository.GetAllPhoneRecords(f.logger, dbName, collectionName)
	assert.Nil(t, err)

	res = sortSliceOFStructByField(res, "ID")

	assert.Contains(t, fmt.Sprint(res[0].ID), repository.GetMD5Hash(data1.PhoneInput))
	assert.Contains(t, fmt.Sprint(res[0].E164Format), data1.PhoneInput)

	assert.Contains(t, fmt.Sprint(res[1].ID), repository.GetMD5Hash(data2.PhoneInput))
	assert.Contains(t, fmt.Sprint(res[1].E164Format), data2.PhoneInput)

}

// go test -v -count=1 -covermode=atomic -tags integration ./internal/repository -run ^TestMongoRespository_UpdatePhoneRecord$
func TestMongoRespository_UpdatePhoneRecord(t *testing.T) {
	tests := []struct {
		name                string
		data                tp.PhoneNumber
		newData             map[string]interface{}
		expected            string
		expectedUpdateError error
	}{
		{
			name: "Update existing phone record",
			data: tp.TestPhoneValue,
			newData: map[string]interface{}{
				"_id":        "0dab7e5e343206634713474e42af8fe3",
				"phoneInput": "1211111",
			},
			expected:            "1211111",
			expectedUpdateError: nil,
		},
		{
			name: "Do not update when new data id does not exist",
			data: tp.TestPhoneValue_2,
			newData: map[string]interface{}{
				"_id":        "0dab7e5e343206634713474e42af8111",
				"phoneInput": "12112345",
			},
			expected:            "",
			expectedUpdateError: errors.New("mongo: no documents in result"),
		},
		{
			name: "Do not update when a key in new data does not exist",
			data: tp.TestPhoneValue_3,
			newData: map[string]interface{}{
				"_id":          "95dee1eaa24f8b7912df00f1ef797a63",
				"keyDontExist": "12112365",
			},
			expected:            "",
			expectedUpdateError: errors.New("key 'keyDontExist' does not exist"),
		},
	}

	f := setUp()
	defer f.tearDown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := f.mongoRepository.StoragePhoneRecord(f.logger, tt.data, dbName, collectionName)
			assert.Nil(t, err)

			res, err := f.mongoRepository.GetAllPhoneRecords(f.logger, dbName, collectionName)
			assert.Nil(t, err)
			assert.Contains(t, fmt.Sprint(res), repository.GetMD5Hash(tt.data.PhoneInput))
			assert.Contains(t, fmt.Sprint(res), tt.data.PhoneInput)

			errUpd := f.mongoRepository.UpdatePhoneRecord(f.logger, tt.newData, dbName, collectionName)
			if errUpd == nil {
				res, err = f.mongoRepository.GetAllPhoneRecords(f.logger, dbName, collectionName)
				assert.Nil(t, err)
				assert.Contains(t, fmt.Sprint(res[0].ID), repository.GetMD5Hash(tt.data.PhoneInput))
				assert.Equal(t, fmt.Sprint(res[0].PhoneInput), tt.expected)
			} else {
				assert.EqualError(t, errUpd, tt.expectedUpdateError.Error())
			}

		})
	}
}

func sortSliceOFStructByField(s tp.PhoneNumberResults, f string) tp.PhoneNumberResults {
	// Use sort.Slice with a custom comparison function
	sort.Slice(s, func(i, j int) bool {
		// Use reflection to get the value of the field f for both elements
		valI := reflect.ValueOf(s[i])
		valJ := reflect.ValueOf(s[j])

		// Get the field value using the field name f
		fieldI := valI.FieldByName(f)
		fieldJ := valJ.FieldByName(f)

		// Make sure the field exists and is of the appropriate type (string in this case)
		if !fieldI.IsValid() || !fieldJ.IsValid() || fieldI.Kind() != reflect.String || fieldJ.Kind() != reflect.String {
			panic(fmt.Sprintf("Field %s does not exist or is not a string", f))
		}

		// Perform the comparison
		return fieldI.String() < fieldJ.String()
	})

	return s
}
