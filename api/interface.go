package api

type PhoneNumberUseCase interface {
	CreatePhoneRecord(phoneNumber string) error
	GetPhoneRecords() (map[string]interface{}, error)
}
