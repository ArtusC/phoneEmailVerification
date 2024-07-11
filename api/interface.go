package api

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
)

type PhoneNumberUseCase interface {
	CollectBigDataCloudApiData(log zerolog.Logger, phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error)
	CreatePhoneRecord(log zerolog.Logger, data t.PhoneNumber) error
	GetAllPhoneRecords(log zerolog.Logger) (t.PhoneNumberResults, error)
	GetPhone(log zerolog.Logger, phoneNumber string) (t.PhoneNumber, error)
}
