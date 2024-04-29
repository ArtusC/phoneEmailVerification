//go:build unit
// +build unit

package mock

import "github.com/stretchr/testify/mock"

type MongoRepositoryMock struct {
	mock.Mock
}

func (mock *MongoRepositoryMock) StoragePhoneRecord(data map[string]interface{}, dbName string, collectionName string) error {
	args := mock.Called(data, dbName, collectionName)
	return args.Error(0)
}

func (mock *MongoRepositoryMock) UpdateePhoneRecord(data map[string]interface{}, dbName string, collectionName string) error {
	args := mock.Called(data, dbName, collectionName)
	return args.Error(0)
}

func (mock *MongoRepositoryMock) GetPhoneRecords(dbName string, collectionName string) (results map[string]interface{}, err error) {
	args := mock.Called(dbName, collectionName)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
