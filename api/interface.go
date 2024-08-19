package api

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
)

type PhoneNumberUseCase interface {
	CollectBigDataCloudApiData(log zerolog.Logger, phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error)
	InsertPhone(log zerolog.Logger, data t.PhoneNumber) error
	UpsertPhone(log zerolog.Logger, data t.PhoneNumber) error
	GetAllPhones(log zerolog.Logger) (t.PhoneNumberResults, error)
	GetPhone(log zerolog.Logger, phoneNumber string) (t.PhoneNumber, error)
}
