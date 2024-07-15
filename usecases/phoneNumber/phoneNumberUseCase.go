package phonenumberusecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
)

const (
	dbName         = "phoneDb"
	collectionName = "phone-collection"
)

type PhoneNumberUseCase struct {
	logger      zerolog.Logger
	storage     repository.MongoRepository
	api_bdc_key string
}

func NewPhoneUseCases(logger zerolog.Logger, mongoRepo repository.MongoRepository, api_bdc_key string) *PhoneNumberUseCase {
	return &PhoneNumberUseCase{
		logger:      logger,
		storage:     mongoRepo,
		api_bdc_key: api_bdc_key,
	}
}

func (p *PhoneNumberUseCase) CreatePhoneRecord(log zerolog.Logger, data t.PhoneNumber) error {
	p.logger.Info().Msg("[UseCase-CreatePhoneRecord] Starting storage.CreatePhoneRecord")
	err := p.storage.StoragePhoneRecord(p.logger, data, dbName, collectionName)
	if err != nil {
		p.logger.Panic().Msgf("[UseCase-CreatePhoneRecord] %s", err.Error())
	}

	return nil
}

func (p *PhoneNumberUseCase) GetAllPhoneRecords(log zerolog.Logger) (t.PhoneNumberResults, error) {
	p.logger.Info().Msg("[UseCase-GetAllPhoneRecords] Starting storage.GetAllPhoneRecords")
	res, err := p.storage.GetAllPhoneRecords(p.logger, dbName, collectionName)
	if err != nil {
		p.logger.Panic().Msgf("[UseCase-GetAllPhoneRecords] %s", err.Error())
	}

	return res, nil
}

func (p *PhoneNumberUseCase) GetPhone(log zerolog.Logger, phoneNumber string) (t.PhoneNumber, error) {
	p.logger.Info().Msg("[UseCase-GetPhone] Starting storage.GetPhone")
	res, err := p.storage.GetPhone(p.logger, dbName, collectionName, phoneNumber)
	if err != nil {
		p.logger.Panic().Msgf("[UseCase-GetAllPhoneRecords] %s", err.Error())
	}

	return res, nil
}

func (p *PhoneNumberUseCase) CollectBigDataCloudApiData(log zerolog.Logger, phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error) {
	p.logger.Info().Msg("[UseCase-CollectBigDataCloudApiData] Starting data collection")

	// Collect data from the API
	data, err := collectData(p.logger, phoneNumber, countryCode, localityLanguage, p.api_bdc_key)
	// data, err := collectData(ctx, "201 867-5309", "us", "en", key)
	if err != nil {
		p.logger.Error().Msgf("[UseCase-CollectBigDataCloudApiData] Error to collect data: %s", err.Error())
		return t.PhoneNumber{}, err
	}

	p.logger.Info().Msgf("[UseCase-CollectBigDataCloudApiData] Data collected: %+v\n", data)

	return data, nil
}

func collectData(logger zerolog.Logger, phoneNumber, countryCode, localityLanguage, apiKey string) (phoneNumberResponse t.PhoneNumber, err error) {

	// Implement logic to collect data from the API
	baseURL := "https://api-bdc.net/data/phone-number-validate"

	logger.Info().Msgf("[collectData] baseURL: %s", baseURL)

	// Create a map to store query parameters
	queryParams := url.Values{}
	queryParams.Set("number", phoneNumber)                // "201 867-5309"
	queryParams.Set("countryCode", countryCode)           // "us"
	queryParams.Set("localityLanguage", localityLanguage) // "en"
	queryParams.Set("key", apiKey)                        // "key"

	// Construct the URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	logger.Info().Msgf("[collectData] fullURL: %s", fullURL)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		logger.Fatal().Msgf("[collectData] error NewRequestWithContext: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Fatal().Msgf("[collectData] error DefaultClient.Do: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		logger.Fatal().Msgf("[collectData] error with status response: %s", resp.Status)
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal().Msgf("[collectData] error io.ReadAll: %s", err.Error())
	}

	json.Unmarshal(bodyText, &phoneNumberResponse)

	logger.Info().Msg("[collectData] data collected!")
	// fmt.Printf("%+v\n", string())

	phoneNumberResponse.PhoneInput = phoneNumber

	return phoneNumberResponse, nil
}
