package phonenumberusecase

import (
	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
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

func (p PhoneNumberUseCase) CreatePhoneRecord(phoneNumber string) error {

	data := map[string]interface{}{
		"phoneNumber": phoneNumber,
	}

	err := p.storage.StoragePhoneRecord(data, dbName, collectionName)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func (p PhoneNumberUseCase) GetPhoneRecords() (map[string]interface{}, error) {

	res, err := p.storage.GetPhoneRecords(dbName, collectionName)
	if err != nil {
		panic(err.Error())
	}

	return res, nil
}
