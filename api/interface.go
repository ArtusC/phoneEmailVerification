package api

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
)

type PhoneNumberUseCase interface {
	CollectBigDataCloudApiData(phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error)
	CreatePhoneRecord(data t.PhoneNumber) error
	GetAllPhoneRecords() (t.PhoneNumberResults, error)
	GetPhone(phoneNumber string) (t.PhoneNumber, error)
}
