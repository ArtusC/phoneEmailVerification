package api

import (
	"github.com/gin-gonic/gin"
)

func (a Api) Healthz(c *gin.Context) {
	a.healthz(c)
}

func (a Api) GetAllPhones(c *gin.Context) {
	a.getAllPhones(c)
}

func (a Api) GetPhone(c *gin.Context) {
	a.getPhone(c)
}

func (a Api) InsertPhone(c *gin.Context) {
	a.insertPhone(c)
}

func (a Api) UpsertPhone(c *gin.Context) {
	a.upsertPhone(c)
}
