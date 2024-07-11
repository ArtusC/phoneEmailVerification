//go:build unit
// +build unit

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ArtusC/phoneEmailVerification/mock"
	ty "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

type fixture struct {
	phoneUseCaseMock *mock.MockPhoneNumberUseCase
	api              Api
	router           *gin.Engine
	logger           zerolog.Logger
}

var traceID gin.HandlerFunc

func setUp() *fixture {

	log := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	phoneUseCaseMock := &mock.MockPhoneNumberUseCase{}

	newApi := NewApi(log, phoneUseCaseMock)

	r := gin.Default()
	r.GET("/healthz", newApi.Healthz)
	r.GET("/api/getAllPhones", newApi.GetAllPhoneRecords)
	r.GET("/api/getPhone/:numberToSearch", newApi.GetPhone)
	r.POST("/api/phoneNumber/:numberToSearch/countryCode/:countryCodeToSearch/localityLanguage/:localityLanguageToSearch", ValidateStoragePhoneRoute(), newApi.CreatePhoneRecord)

	return &fixture{
		phoneUseCaseMock: phoneUseCaseMock,
		api:              *newApi,
		router:           r,
		logger:           log,
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

func TestGetPhone(t *testing.T) {
	f := setUp()

	data := ty.TestPhoneValue

	f.phoneUseCaseMock.On("GetPhone", f.api.logger, "12018675309").Return(data, nil)

	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/getPhone/12018675309", nil)
	f.router.ServeHTTP(r, req)

	var result ty.PhoneNumber
	err = json.Unmarshal(r.Body.Bytes(), &result)

	assert.Equal(t, result, data)
	assert.Exactly(t, http.StatusOK, r.Code, "success")
	assert.Nil(t, err)

}

func TestGetAllPhoneRecords_OneRecord(t *testing.T) {
	f := setUp()

	data := ty.PhoneNumberResults{ty.TestPhoneValue}

	f.phoneUseCaseMock.On("GetAllPhoneRecords", f.api.logger).Return(data, nil)

	req, err := http.NewRequest("GET", "/api/getAllPhones", nil)
	r := httptest.NewRecorder()
	f.router.ServeHTTP(r, req)

	var result ty.PhoneNumberResults
	err = json.Unmarshal(r.Body.Bytes(), &result)

	assert.Equal(t, result, data)
	assert.Exactly(t, http.StatusOK, r.Code, "success")
	assert.Nil(t, err)

}

func TestGetAllPhoneRecords_TwoRecords(t *testing.T) {
	f := setUp()

	data := ty.PhoneNumberResults{ty.TestPhoneValue, ty.TestPhoneValue_2}

	f.phoneUseCaseMock.On("GetAllPhoneRecords", f.api.logger).Return(data, nil)

	req, err := http.NewRequest("GET", "/api/getAllPhones", nil)
	r := httptest.NewRecorder()
	f.router.ServeHTTP(r, req)

	var result ty.PhoneNumberResults
	err = json.Unmarshal(r.Body.Bytes(), &result)

	assert.Equal(t, result, data)
	assert.Exactly(t, http.StatusOK, r.Code, "success")
	assert.Nil(t, err)

}

func TestCreatePhoneRecords(t *testing.T) {
	testCases := []struct {
		name               string
		number             string
		countryCode        string
		localityLanguage   string
		statusExpected     int
		errorExpected      error
		expetedGetData     ty.PhoneNumber
		expetedStorageData ty.PhoneNumber
	}{
		{
			name:               "Valid phone number",
			number:             "12018675309",
			countryCode:        "us",
			localityLanguage:   "en",
			statusExpected:     http.StatusCreated,
			errorExpected:      nil,
			expetedGetData:     ty.TestPhoneValue_Null,
			expetedStorageData: ty.TestPhoneValue,
		},
		{
			name:             "Not valid phone number",
			number:           "a12018675309",
			countryCode:      "us",
			localityLanguage: "en",
			statusExpected:   http.StatusBadRequest,
			errorExpected:    nil,
			expetedGetData:   ty.TestPhoneValue_Null,
		},
		{
			name:             "Phone already exists",
			number:           "12018675355",
			countryCode:      "us",
			localityLanguage: "en",
			statusExpected:   http.StatusConflict,
			errorExpected:    nil,
			expetedGetData:   ty.TestPhoneValue_3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := setUp()

			os.Setenv("API_BDC_KEY", "test_value")

			f.phoneUseCaseMock.On("GetPhone", f.api.logger, tc.number).Return(tc.expetedGetData, tc.errorExpected)
			f.phoneUseCaseMock.On("CollectBigDataCloudApiData", f.api.logger, tc.number, tc.countryCode, tc.localityLanguage).Return(tc.expetedStorageData, tc.errorExpected)
			f.phoneUseCaseMock.On("CreatePhoneRecord", f.api.logger, tc.expetedStorageData).Return(nil)

			url := fmt.Sprintf("/api/phoneNumber/%s/countryCode/%s/localityLanguage/%s", tc.number, tc.countryCode, tc.localityLanguage)

			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			r := httptest.NewRecorder()
			f.router.ServeHTTP(r, req)

			assert.Exactly(t, tc.statusExpected, r.Code, "success")
			assert.Nil(t, err)
		})
	}
}

// TODO: create updatePhoneRecord test
