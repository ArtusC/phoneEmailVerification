package api

import "github.com/gin-gonic/gin"

func (api *Api) routes() *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", api.healthz)

	router.POST("/api/insertPhoneNumber/:numberToInsert", api.createPhoneRecord)
	router.GET("/api/getPhones", api.getPhoneRecords)

	return router
}
