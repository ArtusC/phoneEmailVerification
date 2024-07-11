package repository

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
)

type MongoRepository interface {
	StoragePhoneRecord(log zerolog.Logger, data t.PhoneNumber, dbName, collectionName string) error
	UpdatePhoneRecord(log zerolog.Logger, data map[string]interface{}, dbName, collectionName string) error
	GetAllPhoneRecords(log zerolog.Logger, dbName, collectionName string) (t.PhoneNumberResults, error)
	GetPhone(log zerolog.Logger, dbName, collectionName, phoneNumber string) (t.PhoneNumber, error)
}
