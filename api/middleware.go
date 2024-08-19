package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateStoragePhoneRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		phoneToInsert := ctx.Param("numberToSearch")
		countryCode := ctx.Param("countryCodeToSearch")
		localityLanguage := ctx.Param("localityLanguageToSearch")

		paramsNotNull := phoneToInsert != "" && countryCode != "" && localityLanguage != ""
		if !paramsNotNull {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "this API must be only 3 fields: phoneToInsert, countryCode and localityLanguage",
			})
			return
		}

		if !isNumeric(phoneToInsert) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "phone number must contain only numbers",
			})
			return
		}
		ctx.Next()
	}
}

func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
