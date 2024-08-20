package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Api struct {
	logger        zerolog.Logger
	phoneUseCases PhoneNumberUseCase
	router        *gin.Engine
}

var (
	port = ":8080"
)

func (api Api) StartServer() error {
	api.logger.Info().Msgf("api server running on http://localhost%s", port)
	if err := api.router.Run(port); err != nil {
		api.logger.Error().Msgf("error to start server: %s", err.Error())
		return err
	}
	return nil
}

func NewApi(logger zerolog.Logger, phoneUseCases PhoneNumberUseCase) *Api {
	api := &Api{
		logger:        logger,
		phoneUseCases: phoneUseCases,
	}

	api.logger.Info().Msg("instantiating routes")
	api.router = api.routes()

	return api
}

func send(ctx *gin.Context, code int, val interface{}) {
	ctx.Header("Access-Control-Allow-Methods", "GET, PATCH, POST")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
	ctx.JSON(code, val)
}

func (api Api) healthz(ctx *gin.Context) {
	api.logger.Info().Msg("[API-healthz] received a simple request")
	send(ctx, http.StatusOK, nil)
}

func (api Api) insertPhone(ctx *gin.Context) {
	api.logger.Info().Msg("[API-createPhoneRecord] starting")

	phoneNumber := ctx.Param("numberToSearch")
	countryCode := ctx.Param("countryCodeToSearch")
	localityLanguage := ctx.Param("localityLanguageToSearch")

	api.logger.Info().Msgf("[API-createPhoneRecord] phoneNumber: %s\ncountryCode: %s\nlocalityLanguage: %s\n", phoneNumber, countryCode, localityLanguage)

	phone, err := api.phoneUseCases.GetPhone(api.logger, phoneNumber)
	if err != nil {
		api.logger.Error().Msgf("[API-createPhoneRecord] error to check the phone in db: %s\n", err.Error())
		send(ctx, http.StatusBadRequest, nil)
		return
	}

	if phone.PhoneInput != "" {
		api.logger.Error().Msgf("[API-createPhoneRecord] error: phone number %s already exists\n", phoneNumber)
		send(ctx, http.StatusBadRequest, nil)
		return
	}

	phoneData, err := api.phoneUseCases.CollectBigDataCloudApiData(api.logger, phoneNumber, countryCode, localityLanguage)
	if err != nil {
		api.logger.Error().Msgf("[API-createPhoneRecord] error to collect the phone data: %s\n", err.Error())
		send(ctx, http.StatusBadRequest, nil)
		return
	}

	err = api.phoneUseCases.InsertPhone(api.logger, phoneData)
	if err != nil {
		api.logger.Error().Msgf("[API-createPhoneRecord] error to storage the phone data: %s\n", err.Error())
		send(ctx, http.StatusBadRequest, nil)
		return
	}

	api.logger.Info().Msgf("[API-createPhoneRecord] phone number %s collected and storaged on mongo", phoneNumber)
	send(ctx, http.StatusCreated, nil)
}

func (api Api) getAllPhones(ctx *gin.Context) {
	api.logger.Info().Msg("[API-getAllPhones] starting")

	phones, err := api.phoneUseCases.GetAllPhones(api.logger)
	if err != nil {
		api.logger.Panic().Msgf("[API-getAllPhones] error to get all phone data: %s\n", err.Error())
	}

	api.logger.Info().Msg("[API-getAllPhones] got all phone data")
	send(ctx, http.StatusOK, phones)
}

func (api Api) getPhone(ctx *gin.Context) {
	api.logger.Info().Msg("[API-getPhone] starting")

	phoneNumber := ctx.Param("numberToSearch")
	api.logger.Info().Msgf("[API-getPhone] phoneNumber: %s\n", phoneNumber)

	phone, err := api.phoneUseCases.GetPhone(api.logger, phoneNumber)
	if err != nil {
		api.logger.Panic().Msgf("[API-getPhone] error to get phone data: %s\n", err.Error())
	}

	api.logger.Info().Msgf("[API-getPhone] got phone data %s\n", phoneNumber)
	send(ctx, http.StatusOK, phone)
}

func (api Api) upsertPhone(ctx *gin.Context) {
	api.logger.Info().Msg("[API-upsertPhone] starting")

	phoneNumber := ctx.Param("numberToSearch")
	countryCode := ctx.Param("countryCodeToSearch")
	localityLanguage := ctx.Param("localityLanguageToSearch")

	api.logger.Info().Msgf("[API-upsertPhone] phoneNumber: %s\ncountryCode: %s\nlocalityLanguage: %s\n", phoneNumber, countryCode, localityLanguage)

	phoneData, err := api.phoneUseCases.CollectBigDataCloudApiData(api.logger, phoneNumber, countryCode, localityLanguage)
	if err != nil {
		api.logger.Error().Msgf("[API-upsertPhone] error to collect the phone data: %s\n", err.Error())
		send(ctx, http.StatusBadRequest, nil)
		return
	}

	err = api.phoneUseCases.UpsertPhone(api.logger, phoneData)
	if err != nil {
		api.logger.Error().Msgf("[API-upsertPhone] error to storage the phone data: %s\n", err.Error())
		send(ctx, http.StatusBadRequest, nil)
		return
	}

	api.logger.Info().Msgf("[API-upsertPhone] phone number %s collected and storaged on mongo", phoneNumber)
	send(ctx, http.StatusCreated, nil)
}
