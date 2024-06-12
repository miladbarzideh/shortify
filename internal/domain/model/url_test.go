package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type URLTestSuite struct {
	suite.Suite
}

func (suite *URLTestSuite) TestURL_Validate_Success() {
	require := suite.Require()
	testCases := []struct {
		input          string
		expectedResult bool
	}{
		{
			input:          "http://google.com",
			expectedResult: true,
		},
		{
			input:          "https://google.com",
			expectedResult: true,
		},
		{
			input:          "https://echo.labstack.com/docs/binding",
			expectedResult: true,
		},
		{
			input:          "http/google.com",
			expectedResult: false,
		},
		{
			input:          "google.com",
			expectedResult: false,
		},
		{
			input:          "/api/v1/urls",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		url := URLData{URL: tc.input}
		actualResult := url.Validate()

		require.Equal(tc.expectedResult, actualResult)
	}
}

func TestURLTestSuite(t *testing.T) {
	suite.Run(t, new(URLTestSuite))
}
