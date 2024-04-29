package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Api struct {
	phoneUseCases PhoneNumberUseCase
	router        *gin.Engine
}

func (api Api) StartServer() error {
	// log := api.log.Method("main").CreateTraceID()
	// log.Debug("api server running on http://localhost:80\n")

	if err := api.router.Run(":8080"); err != nil {
		// log.Error(err.Error())
		return err
	}
	return nil
}

func NewApi(phoneUseCases PhoneNumberUseCase) Api {
	api := Api{
		phoneUseCases: phoneUseCases,
	}

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
	// log := api.log.Method("healthz").CreateTraceID()
	// log.Debug("receive a simple request")

	send(ctx, http.StatusOK, nil)
}

func (api Api) createPhoneRecord(ctx *gin.Context) {

	// TEST WITH UNITY TEST FIRST: CHANGE THE STRUCT THAT WILL GO TO THE CreatePhoneRecord
	// TODO: verify number already exists (if yes, use update route) -> collect data from API -> send to storage

	phoneNumber := ctx.Param("insertPhoneNumber")

	err := api.phoneUseCases.CreatePhoneRecord(phoneNumber)
	if err != nil {
		panic(err.Error())
	}
}

// TODO: create updatePhoneRecord route

func (api Api) getPhoneRecords(ctx *gin.Context) {

	phones, err := api.phoneUseCases.GetPhoneRecords()
	if err != nil {
		panic(err.Error())
	}

	send(ctx, http.StatusOK, phones)
}
