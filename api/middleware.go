package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			fmt.Println("Phone number is not numeric")
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

// CorrelationIDMiddleware adds a correlation ID to the context and response headers
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}

		// Add the correlation ID to the context
		c.Set("CorrelationID", correlationID)

		// Add the correlation ID to the response headers
		c.Writer.Header().Set("X-Correlation-ID", correlationID)

		// Proceed to the next middleware/handler
		c.Next()
	}
}

// Generate a unique correlation ID
func generateCorrelationID() string {
	return uuid.New().String()
}
