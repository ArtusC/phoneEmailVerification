//go:build unit
// +build unit

package mock

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/stretchr/testify/mock"
)

type MockMongoRepository struct {
	mock.Mock
}

func (mock *MockMongoRepository) StoragePhoneRecord(data map[string]interface{}, dbName, collectionName string) error {
	args := mock.Called(data, dbName, collectionName)
	return args.Error(0)
}

func (mock *MockMongoRepository) UpdateePhoneRecord(data map[string]interface{}, dbName, collectionName string) error {
	args := mock.Called(data, dbName, collectionName)
	return args.Error(0)
}

func (mock *MockMongoRepository) GetAllPhoneRecords(dbName, collectionName string) (results map[string]interface{}, err error) {
	args := mock.Called(dbName, collectionName)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (mock *MockMongoRepository) GetPhone(dbName, collectionName, phoneNumber string) (result t.PhoneNumber, err error) {
	args := mock.Called(dbName, collectionName, phoneNumber)
	return args.Get(0).(t.PhoneNumber), args.Error(1)
}
