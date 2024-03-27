package api

import "github.com/gin-gonic/gin"

func (api *Api) routes() *gin.Engine {
	router := gin.Default()

	router.POST("/api/:insertPhoneNumber", api.createPhoneRecord)
	router.GET("/api/getPhones", api.getPhoneRecords)

	return router
}
