package service

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/infra"
	"github.com/miladbarzideh/shortify/internal/model"
	"github.com/miladbarzideh/shortify/internal/repository/mock"
	genMock "github.com/miladbarzideh/shortify/pkg/generator/mock"
)

type URLServiceTestSuite struct {
	suite.Suite
	service       URLService
	mockRepo      *mock.Repository
	mockCacheRepo *mock.CacheRepository
	mockGen       *genMock.Generator
}

func (suite *URLServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mock.Repository)
	suite.mockCacheRepo = new(mock.CacheRepository)
	suite.mockGen = new(genMock.Generator)
	cfg := infra.Config{}
	cfg.Server.Address = "localhost:8513"
	cfg.Shortener.CodeLength = 7
	suite.service = NewService(logrus.New(), &cfg, suite.mockRepo, suite.mockCacheRepo, suite.mockGen)
}

func (suite *URLServiceTestSuite) TestURLService_CreateShortURL_Success() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "http://google.com",
			expectedURL: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "gclmd",
			},
		},
	}

	for _, tc := range testCases {
		suite.mockGen.On("GenerateShortURLCode", testifyMock.Anything).Return(tc.expectedURL.ShortCode)
		suite.mockRepo.On("Create", tc.input, tc.expectedURL.ShortCode).Return(tc.expectedURL, nil).Once()
		url, err := suite.service.CreateShortURL(tc.input)

		require.NoError(err)
		require.NotEmpty(url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_CreateShortURL_Retry_Success() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "http://google.com",
			expectedURL: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "gclmd",
			},
		},
	}

	for _, tc := range testCases {
		suite.mockGen.On("GenerateShortURLCode", testifyMock.Anything).Return(tc.expectedURL.ShortCode)
		suite.mockRepo.On("Create", tc.input, tc.expectedURL.ShortCode).Return(model.URL{}, gorm.ErrDuplicatedKey)
		suite.mockRepo.On("Create", tc.input, tc.expectedURL.ShortCode).Return(tc.expectedURL, nil)
		url, err := suite.service.CreateShortURL(tc.input)

		require.NoError(err)
		require.NotEmpty(url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_CreateShortURL_Failure() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "http://google.com",
			expectedURL: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "gclmd",
			},
		},
	}

	for _, tc := range testCases {
		suite.mockGen.On("GenerateShortURLCode", testifyMock.Anything).Return(tc.expectedURL.ShortCode)
		suite.mockRepo.On("Create", tc.input, tc.expectedURL.ShortCode).Return(model.URL{}, gorm.ErrInvalidData).Once()
		url, err := suite.service.CreateShortURL(tc.input)

		require.Error(err)
		require.Empty(url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_GetLongURL_ReadFromCache_Success() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "G2ogLe",
			expectedURL: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "G2ogLe",
				ID:        1,
			},
		},
	}
	ctx := context.TODO()

	for _, tc := range testCases {
		suite.mockCacheRepo.On("Get", ctx, tc.input).Return(tc.expectedURL, nil).Once()
		suite.mockRepo.On("FindByShortCode", testifyMock.AnythingOfType("*model.URL")).Times(0)
		url, err := suite.service.GetLongURL(ctx, tc.input)

		require.NoError(err)
		require.Equal(tc.expectedURL.LongURL, url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_GetLongURL_ReadFromDb_Success() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "G2ogLe",
			expectedURL: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "G2ogLe",
				ID:        1,
			},
		},
	}
	ctx := context.TODO()

	for _, tc := range testCases {
		suite.mockCacheRepo.On("Get", ctx, tc.input).Return(model.URL{}, redis.Nil).Once()
		suite.mockRepo.On("FindByShortCode", tc.input).Return(tc.expectedURL, nil).Once()
		suite.mockCacheRepo.On("Set", ctx, tc.expectedURL).Return(nil)
		url, err := suite.service.GetLongURL(ctx, tc.input)

		require.NoError(err)
		require.Equal(tc.expectedURL.LongURL, url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_GetLongURL_Failure() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "G2ogLe",
			expectedURL: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "G2ogLe",
				ID:        1,
			},
		},
	}
	ctx := context.TODO()

	for _, tc := range testCases {
		suite.mockCacheRepo.On("Get", ctx, tc.input).Return(model.URL{}, redis.Nil).Once()
		suite.mockRepo.On("FindByShortCode", tc.input).Return(model.URL{}, gorm.ErrRecordNotFound).Once()
		url, err := suite.service.GetLongURL(ctx, tc.input)

		require.Error(err)
		require.Empty(url)
	}
}

func TestURLServiceTestSuite(t *testing.T) {
	suite.Run(t, new(URLServiceTestSuite))
}
