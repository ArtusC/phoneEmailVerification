package repository

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
)

type MongoRepository interface {
	StoragePhoneRecord(data t.PhoneNumber, dbName, collectionName string) error
	UpdatePhoneRecord(data map[string]interface{}, dbName, collectionName string) error
	GetAllPhoneRecords(dbName, collectionName string) (t.PhoneNumberResults, error)
	GetPhone(dbName, collectionName, phoneNumber string) (t.PhoneNumber, error)
}
