//go:build unit
// +build unit

package mock

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

type MockPhoneNumberUseCase struct {
	mock.Mock
}

func (mock *MockPhoneNumberUseCase) CollectBigDataCloudApiData(log zerolog.Logger, phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error) {
	args := mock.Called(log, phoneNumber, countryCode, localityLanguage)
	return args.Get(0).(t.PhoneNumber), args.Error(1)
}

func (mock *MockPhoneNumberUseCase) CreatePhoneRecord(log zerolog.Logger, phoneRecord t.PhoneNumber) error {
	args := mock.Called(log, phoneRecord)
	return args.Error(0)
}

func (mock *MockPhoneNumberUseCase) GetAllPhoneRecords(log zerolog.Logger) (t.PhoneNumberResults, error) {
	args := mock.Called(log)
	return args.Get(0).(t.PhoneNumberResults), args.Error(1)
}

func (mock *MockPhoneNumberUseCase) GetPhone(log zerolog.Logger, phoneNumber string) (t.PhoneNumber, error) {
	args := mock.Called(log, phoneNumber)
	return args.Get(0).(t.PhoneNumber), args.Error(1)
}
