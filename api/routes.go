package api

import (
	"github.com/gin-gonic/gin"
)

func (api *Api) routes() *gin.Engine {
	router := gin.Default()

	router.Use(CorrelationIDMiddleware())

	router.GET("/healthz", api.healthz)

	// /api/phoneNumber/12018675309/countryCode/us/localityLanguage/en
	router.POST("/api/phoneNumber/:numberToSearch/countryCode/:countryCodeToSearch/localityLanguage/:localityLanguageToSearch", ValidateStoragePhoneRoute(), api.createPhoneRecord)
	router.GET("/api/getAllPhones", api.getAllPhoneRecords)
	router.GET("/api/getPhone/:numberToSearch", api.getPhone)

	return router
}
