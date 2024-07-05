//go:build unit
// +build unit

package mock

import (
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/stretchr/testify/mock"
)

type MockPhoneNumberUseCase struct {
	mock.Mock
}

func (mock *MockPhoneNumberUseCase) CollectBigDataCloudApiData(phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error) {
	args := mock.Called(phoneNumber, countryCode, localityLanguage)
	return args.Get(0).(t.PhoneNumber), args.Error(1)
}

func (mock *MockPhoneNumberUseCase) CreatePhoneRecord(phoneRecord t.PhoneNumber) error {
	args := mock.Called(phoneRecord)
	return args.Error(0)
}

func (mock *MockPhoneNumberUseCase) GetAllPhoneRecords() (t.PhoneNumberResults, error) {
	args := mock.Called()
	return args.Get(0).(t.PhoneNumberResults), args.Error(1)
}

func (mock *MockPhoneNumberUseCase) GetPhone(phoneNumber string) (t.PhoneNumber, error) {
	args := mock.Called(phoneNumber)
	return args.Get(0).(t.PhoneNumber), args.Error(1)
}
