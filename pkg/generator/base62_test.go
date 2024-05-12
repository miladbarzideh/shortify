package generator

import (
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
	generator Generator
	patches   *gomonkey.Patches
}

func (suite *GeneratorTestSuite) SetupTest() {
	suite.patches = gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Date(2024, time.May, 11, 19, 47, 0, 0, time.UTC)
	})
	suite.generator = NewGenerator()
}

func (suite *GeneratorTestSuite) TearDownTest() {
	suite.patches.Reset()
}

func (suite *GeneratorTestSuite) TestGenerator_GenerateShortURLCode_Success() {
	require := suite.Require()
	testCases := []struct {
		input    int
		expected string
	}{
		{
			input:    5,
			expected: "M7OBS",
		},
		{
			input:    7,
			expected: "Hn1OSj1",
		},
	}

	for _, tc := range testCases {
		actual := suite.generator.GenerateShortURLCode(tc.input)

		require.Equal(tc.input, len(actual))
		require.Equal(tc.expected, actual)
	}
}

func (suite *GeneratorTestSuite) TestGenerator_IsValidBase62_Success() {
	require := suite.Require()
	testCases := []struct {
		input          string
		expectedResult bool
	}{
		{input: "abc", expectedResult: true},
		{input: "/;;p", expectedResult: false},
		{input: "Abc4e", expectedResult: true},
		{input: "gool=", expectedResult: false},
		{input: "90#l", expectedResult: false},
	}

	for _, tc := range testCases {
		actualResult := IsValidBase62(tc.input)

		require.Equal(tc.expectedResult, actualResult)
	}
}

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
