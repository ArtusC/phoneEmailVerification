package phonenumberusecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	repository "github.com/ArtusC/phoneEmailVerification/internal/repository"
	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

const (
	dbName         = "phoneDb"
	collectionName = "phone-collection"
)

type phoneNumberUseCase struct {
	ctx         *gin.Context
	logger      zerolog.Logger
	storage     repository.MongoRepository
	api_bdc_key string
}

func NewPhoneUseCases(ctx *gin.Context, logger zerolog.Logger, mongoRepo repository.MongoRepository, api_bdc_key string) *phoneNumberUseCase {
	return &phoneNumberUseCase{
		ctx:         ctx,
		logger:      logger,
		storage:     mongoRepo,
		api_bdc_key: api_bdc_key,
	}
}

func (p *phoneNumberUseCase) CreatePhoneRecord(data t.PhoneNumber) error {
	traceId, _ := p.ctx.Get("CorrelationID")
	p.logger.Info().Str("traceId", traceId.(string)).Msg("[UseCase-CreatePhoneRecord] Starting storage.CreatePhoneRecord")
	err := p.storage.StoragePhoneRecord(data, dbName, collectionName)
	if err != nil {
		p.logger.Panic().Str("traceId", traceId.(string)).Msgf("[UseCase-CreatePhoneRecord] %s", err.Error())
		// panic(fmt.Sprintf("[UseCase-CreatePhoneRecord] %s", err.Error()))
	}

	return nil
}

func (p *phoneNumberUseCase) GetAllPhoneRecords() (t.PhoneNumberResults, error) {
	traceId, _ := p.ctx.Get("CorrelationID")
	p.logger.Info().Str("traceId", traceId.(string)).Msg("[UseCase-GetAllPhoneRecords] Starting storage.GetAllPhoneRecords")
	res, err := p.storage.GetAllPhoneRecords(dbName, collectionName)
	if err != nil {
		p.logger.Panic().Str("traceId", traceId.(string)).Msgf("[UseCase-GetAllPhoneRecords] %s", err.Error())
		// panic(fmt.Sprintf("[UseCase-GetAllPhoneRecords] %s", err.Error()))
	}

	return res, nil
}

func (p *phoneNumberUseCase) GetPhone(phoneNumber string) (t.PhoneNumber, error) {
	traceId, _ := p.ctx.Get("CorrelationID")
	p.logger.Info().Str("traceId", traceId.(string)).Msg("[UseCase-GetPhone] Starting storage.GetPhone")
	res, err := p.storage.GetPhone(dbName, collectionName, phoneNumber)
	if err != nil {
		p.logger.Panic().Str("traceId", traceId.(string)).Msgf("[UseCase-GetAllPhoneRecords] %s", err.Error())
		// panic(fmt.Sprintf("[UseCase-GetPhone] %s", err.Error()))
	}

	return res, nil
}

func (p *phoneNumberUseCase) CollectBigDataCloudApiData(phoneNumber, countryCode, localityLanguage string) (t.PhoneNumber, error) {
	traceId, _ := p.ctx.Get("CorrelationID")
	p.logger.Info().Str("traceId", traceId.(string)).Msg("[UseCase-CollectBigDataCloudApiData] Starting data collection")

	// Collect data from the API
	ctx := context.Background()

	data, err := collectData(ctx, p.logger, traceId.(string), phoneNumber, countryCode, localityLanguage, p.api_bdc_key)
	// data, err := collectData(ctx, "201 867-5309", "us", "en", key)
	if err != nil {
		p.logger.Error().Str("traceId", traceId.(string)).Msgf("[UseCase-CollectBigDataCloudApiData] Error to collect data: %s", err.Error())
		return t.PhoneNumber{}, err
	}

	p.logger.Info().Str("traceId", traceId.(string)).Msgf("[UseCase-CollectBigDataCloudApiData] Data collected: %+v\n", data)

	return data, nil
}

func collectData(ctx context.Context, logger zerolog.Logger, traceId, phoneNumber, countryCode, localityLanguage, apiKey string) (phoneNumberResponse t.PhoneNumber, err error) {

	// Implement logic to collect data from the API
	baseURL := "https://api-bdc.net/data/phone-number-validate"

	logger.Info().Str("traceId", traceId).Msgf("[collectData] baseURL: %s", baseURL)

	// Create a map to store query parameters
	queryParams := url.Values{}
	queryParams.Set("number", phoneNumber)                // "201 867-5309"
	queryParams.Set("countryCode", countryCode)           // "us"
	queryParams.Set("localityLanguage", localityLanguage) // "en"
	queryParams.Set("key", apiKey)                        // "key"

	// Construct the URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	logger.Info().Str("traceId", traceId).Msgf("[collectData] fullURL: %s", fullURL)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		logger.Fatal().Str("traceId", traceId).Msgf("[collectData] error NewRequestWithContext: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Fatal().Str("traceId", traceId).Msgf("[collectData] error DefaultClient.Do: %s", err.Error())
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal().Str("traceId", traceId).Msgf("[collectData] error io.ReadAll: %s", err.Error())
	}

	json.Unmarshal(bodyText, &phoneNumberResponse)

	logger.Info().Str("traceId", traceId).Msg("[collectData] data collected!")
	// fmt.Printf("%+v\n", string())

	phoneNumberResponse.PhoneInput = phoneNumber

	return phoneNumberResponse, nil
}
