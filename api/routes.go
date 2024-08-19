package api

import (
	"github.com/gin-gonic/gin"
)

func (api *Api) routes() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", api.healthz)

	// /api/phoneNumber/12018675309/countryCode/us/localityLanguage/en
	router.POST("/api/phoneNumber/:numberToSearch/countryCode/:countryCodeToSearch/localityLanguage/:localityLanguageToSearch", ValidateStoragePhoneRoute(), api.insertPhone)
	router.PUT("/api/phoneNumber/:numberToSearch/countryCode/:countryCodeToSearch/localityLanguage/:localityLanguageToSearch", ValidateStoragePhoneRoute(), api.upsertPhone)
	router.GET("/api/getAllPhones", api.getAllPhones)
	router.GET("/api/getPhone/:numberToSearch", api.getPhone)

	return router
}
