package phonenumberusecase

import (
	"errors"

	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	t "github.com/ArtusC/phoneEmailVerification/types"
)

const (
	dbName         = "phoneDb"
	collectionName = "phone-collection"
)

type PhoneNumberUseCase struct {
	storage repository.MongoRepository
}

func NewPhoneUseCases(mongoRepo repository.MongoRepository) PhoneNumberUseCase {
	return PhoneNumberUseCase{
		storage: mongoRepo,
	}
}

func (p PhoneNumberUseCase) CreatePhoneRecord(data t.PhoneNumber) error {

	if data.PhoneInput == "" {
		return errors.New("the phone number is required")
	}

	err := p.storage.StoragePhoneRecord(data, dbName, collectionName)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (p PhoneNumberUseCase) GetPhoneRecords() (t.PhoneNumberResults, error) {

	res, err := p.storage.GetPhoneRecords(dbName, collectionName)
	if err != nil {
		panic(err.Error())
	}

	return res, nil
}
