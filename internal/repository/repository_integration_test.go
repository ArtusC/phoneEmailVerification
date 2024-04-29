//go:build unit || integration
// +build unit integration

package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	tp "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type fixture struct {
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

	//TODO: change the docker-compose that will up the database (to not sobrescribe), dont put login/password
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("MongoDBTest connection established!")

	if mongoSession, err = client.StartSession(); err != nil {
		panic(err.Error())
	}

	repository := repository.NewMongoRepository(mongoSession)

	return &fixture{
		mongoSession:    mongoSession,
		mongoRepository: repository,
	}
}

func (f *fixture) tearDown() {
	fmt.Println("MongoDBTest connection finished!")
	defer f.mongoSession.EndSession(context.Background())
	d := f.mongoSession.Client().Database(dbName)
	d.Drop(context.TODO())
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
			phoneOutput:   TestPhoneValue,
			expectedError: nil,
		},
	}

	for _, test := range testCase {
		t.Run(test.testName, func(t *testing.T) {

			err := f.mongoRepository.StoragePhoneRecord(test.phoneOutput, dbName, collectionName)
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

// go test -v -count=1 -covermode=atomic -tags integration ./internal/repository -run ^TestMongoRespository_GetPhoneRecords$
func TestMongoRespository_GetPhoneRecords(t *testing.T) {
	f := setUp()
	defer f.tearDown()

	data := TestPhoneValue

	err := f.mongoRepository.StoragePhoneRecord(data, dbName, collectionName)
	assert.Nil(t, err)

	res, err := f.mongoRepository.GetPhoneRecords(dbName, collectionName)
	assert.Nil(t, err)

	fmt.Println("Result: ", res)

	assert.Contains(t, fmt.Sprint(res[0].ID), repository.GetMD5Hash(data.PhoneInput))
	assert.Contains(t, fmt.Sprint(res[0].E164Format), data.PhoneInput)

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
			data: TestPhoneValue,
			newData: map[string]interface{}{
				"_id":        "0dab7e5e343206634713474e42af8fe3",
				"phoneInput": "1211111",
			},
			expected:            "1211111",
			expectedUpdateError: nil,
		},
		{
			name: "Do not update when new data id does not exist",
			data: TestPhoneValue_2,
			newData: map[string]interface{}{
				"_id":        "0dab7e5e343206634713474e42af8111",
				"phoneInput": "12112345",
			},
			expected:            "",
			expectedUpdateError: errors.New("mongo: no documents in result"),
		},
		{
			name: "Do not update when a key in new data does not exist",
			data: TestPhoneValue_3,
			newData: map[string]interface{}{
				"_id":          "95dee1eaa24f8b7912df00f1ef797a63",
				"keyDontExist": "12112365",
			},
			expected:            "",
			expectedUpdateError: errors.New("key 'keyDontExist' does not exist in the existing document"),
		},
	}

	f := setUp()
	defer f.tearDown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := f.mongoRepository.StoragePhoneRecord(tt.data, dbName, collectionName)
			assert.Nil(t, err)

			res, err := f.mongoRepository.GetPhoneRecords(dbName, collectionName)
			assert.Nil(t, err)
			assert.Contains(t, fmt.Sprint(res), repository.GetMD5Hash(tt.data.PhoneInput))
			assert.Contains(t, fmt.Sprint(res), tt.data.PhoneInput)

			errUpd := f.mongoRepository.UpdatePhoneRecord(tt.newData, dbName, collectionName)
			fmt.Println("Error: ", errUpd)
			if errUpd == nil {
				res, err = f.mongoRepository.GetPhoneRecords(dbName, collectionName)
				assert.Nil(t, err)
				assert.Contains(t, fmt.Sprint(res[0].ID), repository.GetMD5Hash(tt.data.PhoneInput))
				assert.Equal(t, fmt.Sprint(res[0].PhoneInput), tt.expected)
			} else {
				assert.EqualError(t, errUpd, tt.expectedUpdateError.Error())
			}

		})
	}
}
