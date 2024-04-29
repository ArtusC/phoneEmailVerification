//go:build unit
// +build unit

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtusC/phoneEmailVerification/mock"
	ty "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

// PARA O TESTE UNITÁRIO DO API
// package repository_test

type fixture struct {
	phoneUseCaseMock *mock.PhoneNumberUseCaseMock
	api              Api
	router           *gin.Engine
}

func setUp() *fixture {
	phoneUseCaseMock := &mock.PhoneNumberUseCaseMock{}

	newApi := NewApi(phoneUseCaseMock)

	r := gin.Default()
	r.GET("/healthz", newApi.Healthz)
	r.GET("/api/getPhones", newApi.GetPhoneRecordsExported)
	r.POST("/api/insertPhoneNumber/:numberToInsert", newApi.CreatePhoneRecordExported)

	return &fixture{
		phoneUseCaseMock: phoneUseCaseMock,
		api:              newApi,
		router:           r,
	}

}

func TestHealthzAPI(t *testing.T) {
	f := setUp()

	req, err := http.NewRequest("GET", "/healthz", nil)
	r := httptest.NewRecorder()

	f.router.ServeHTTP(r, req)

	assert.Exactly(t, http.StatusOK, r.Code, "success")
	assert.Nil(t, err)

}

func TestGetPhoneRecords(t *testing.T) {
	f := setUp()

	f.phoneUseCaseMock.On("GetPhoneRecordsExported").Return(ty.TestPhoneValue, nil)

	req, err := http.NewRequest("GET", "/api/getPhones", nil)
	r := httptest.NewRecorder()
	f.router.ServeHTTP(r, req)

	var result ty.PhoneNumber
	err = json.Unmarshal(r.Body.Bytes(), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	assert.Equal(t, result, ty.TestPhoneValue)
	assert.Exactly(t, http.StatusOK, http.StateClosed, "success")
	assert.Nil(t, err)

}

func TestCreatePhoneRecords(t *testing.T) {
	f := setUp()

	f.phoneUseCaseMock.On("CreatePhoneRecord").Return(nil)

	req, err := http.NewRequest("POST", "/api/insertPhoneNumber/12018675309", nil)
	fmt.Println("req: ", req)
	r := httptest.NewRecorder()
	f.router.ServeHTTP(r, req)

	var result map[string]interface{}
	err = json.Unmarshal(r.Body.Bytes(), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	assert.Exactly(t, http.StatusOK, r.Code, "success")
	assert.Nil(t, err)

}
