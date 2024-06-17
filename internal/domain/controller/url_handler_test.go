package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/domain/service"
	"github.com/miladbarzideh/shortify/internal/domain/service/mock"
	infra2 "github.com/miladbarzideh/shortify/internal/infra"
)

const (
	apiURL = "localhost:8513/api/v1/urls/"
)

type URLHandlerTestSuite struct {
	suite.Suite
	mockService *mock.Service
	handler     URLHandler
}

func (suite *URLHandlerTestSuite) SetupTest() {
	suite.mockService = new(mock.Service)
	suite.handler = NewHandler(logrus.New(), &infra2.Config{}, suite.mockService, infra2.NOOPTelemetry)
}

func (suite *URLHandlerTestSuite) TestURLHandler_CreateShortURL_Success() {
	require := suite.Require()
	testCases := []struct {
		input            model.URLData
		expectedResponse model.URLData
		expectedCode     int
	}{
		{
			input:            model.URLData{URL: "https://www.google.com"},
			expectedResponse: model.URLData{URL: apiURL + "R849E"},
			expectedCode:     http.StatusOK,
		},
	}

	for _, tc := range testCases {
		c, rec := newEchoContext(http.MethodPost, "/api/v1/urls", tc.input, "")

		suite.mockService.On("CreateShortURL", testifymock.Anything, testifymock.Anything).Return(tc.expectedResponse.URL, nil)
		err := suite.handler.CreateShortURL()(c)

		require.NoError(err)
		require.Equal(tc.expectedCode, rec.Code)
		var actual model.URLData
		err = json.Unmarshal(rec.Body.Bytes(), &actual)
		require.NoError(err)
		require.Equal(tc.expectedResponse, actual)
	}
}

func (suite *URLHandlerTestSuite) TestURLHandler_CreateShortURL_Failure() {
	require := suite.Require()
	testCases := []struct {
		input        interface{}
		err          error
		expectedCode int
	}{
		{
			input:        "{invalid}",
			expectedCode: http.StatusBadRequest,
		},
		{
			input:        model.URLData{URL: "https/echo.labstack.com/docs/testing"},
			expectedCode: http.StatusBadRequest,
		},
		{
			input:        model.URLData{URL: "https://echo.labstack.com/docs/testing"},
			err:          service.ErrMaxRetriesExceeded,
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			input:        model.URLData{URL: "https://echo.labstack.com/docs/testing"},
			err:          gorm.ErrInvalidData,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		c, _ := newEchoContext(http.MethodPost, "/api/v1/urls", tc.input, "")

		if tc.err != nil {
			suite.mockService.On("CreateShortURL", testifymock.Anything, testifymock.Anything).Return("", tc.err).Once()
		}

		err := suite.handler.CreateShortURL()(c)

		require.Error(err)
		require.IsType(&echo.HTTPError{}, err)
		require.Equal(tc.expectedCode, err.(*echo.HTTPError).Code)
	}
}

func (suite *URLHandlerTestSuite) TestURLHandler_RedirectToLongURL_Success() {
	require := suite.Require()
	expectedCode := http.StatusMovedPermanently
	testCases := []struct {
		input       string
		expectedURL string
	}{
		{
			input:       "R849E",
			expectedURL: "https://www.google.com",
		},
		{
			input:       "L7dRf",
			expectedURL: "https://echo.labstack.com/docs/testing",
		},
	}

	for _, tc := range testCases {
		c, rec := newEchoContext(http.MethodGet, "/api/v1/urls/"+tc.input, nil, tc.input)

		suite.mockService.On("GetLongURL", testifymock.Anything, tc.input).Return(tc.expectedURL, nil)
		err := suite.handler.RedirectToLongURL()(c)

		require.NoError(err)
		require.Equal(expectedCode, rec.Code)
		require.Equal(tc.expectedURL, rec.Header().Get("Location"))
	}
}

func (suite *URLHandlerTestSuite) TestURLHandler_RedirectToLongURL_Failure() {
	require := suite.Require()
	testCases := []struct {
		input        string
		expectedCode int
	}{
		{
			input:        "R849E",
			expectedCode: http.StatusNotFound,
		},
		{
			input:        "=;))",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		c, _ := newEchoContext(http.MethodGet, "/api/v1/urls/"+tc.input, nil, tc.input)

		suite.mockService.On("GetLongURL", testifymock.Anything, tc.input).Return("", service.ErrURLNotFound)
		err := suite.handler.RedirectToLongURL()(c)

		require.Error(err)
		require.Equal(tc.expectedCode, err.(*echo.HTTPError).Code)
	}
}

func TestURLHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerTestSuite))
}

func newEchoContext(method string, endpoint string, body any, param string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, endpoint, nil)
	if body != nil {
		result, _ := json.Marshal(body)
		req = httptest.NewRequest(method, endpoint, strings.NewReader(string(result)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if param != "" {
		c.SetPath("/api/v1/urls/:url")
		c.SetParamNames("url")
		c.SetParamValues(param)
	}

	return c, rec
}
