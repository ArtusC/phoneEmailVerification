package api

import (
	"github.com/gin-gonic/gin"
)

func (a Api) Healthz(c *gin.Context) {
	a.healthz(c)
}

func (a Api) GetPhoneRecordsExported(c *gin.Context) {
	a.getPhoneRecords(c)
}

func (a Api) CreatePhoneRecordExported(c *gin.Context) {
	a.createPhoneRecord(c)
}
