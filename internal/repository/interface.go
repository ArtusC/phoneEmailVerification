package repository

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
)

type MongoRepository interface {
	StoragePhoneRecord(data t.PhoneNumber, dbName string, collectionName string) error
	UpdatePhoneRecord(data map[string]interface{}, dbName string, collectionName string) error
	GetPhoneRecords(dbName string, collectionName string) (t.PhoneNumberResults, error)
}
