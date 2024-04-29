package api

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
)

type PhoneNumberUseCase interface {
	CreatePhoneRecord(phoneNumber string) error
	GetPhoneRecords() (t.PhoneNumberResults, error)
}
