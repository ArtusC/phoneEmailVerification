package api

import (
	"github.com/gin-gonic/gin"
)

func (a Api) Healthz(c *gin.Context) {
	a.healthz(c)
}

func (a Api) GetAllPhoneRecords(c *gin.Context) {
	a.getAllPhoneRecords(c)
}

func (a Api) GetPhone(c *gin.Context) {
	a.getPhone(c)
}

func (a Api) CreatePhoneRecord(c *gin.Context) {
	a.createPhoneRecord(c)
}
