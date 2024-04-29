//go:build unit
// +build unit

package mock

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/stretchr/testify/mock"
)

type PhoneNumberUseCaseMock struct {
	mock.Mock
}

func (mock *PhoneNumberUseCaseMock) CreatePhoneRecord(phoneNumber string) error {
	args := mock.Called(phoneNumber)
	return args.Error(0)
}

func (mock *PhoneNumberUseCaseMock) GetPhoneRecords() (t.PhoneNumberResults, error) {
	args := mock.Called()
	return args.Get(0).(t.PhoneNumberResults), args.Error(1)
}
