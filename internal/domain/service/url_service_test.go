package service

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/domain/repository/mock"
	"github.com/miladbarzideh/shortify/internal/infra"
	genMock "github.com/miladbarzideh/shortify/pkg/generator/mock"
	wpMock "github.com/miladbarzideh/shortify/pkg/worker/mock"
)

type URLServiceTestSuite struct {
	suite.Suite
	service       URLService
	mockRepo      *mock.Repository
	mockCacheRepo *mock.CacheRepository
	mockGen       *genMock.Generator
	mockWP        *wpMock.Pool
}

func (suite *URLServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mock.Repository)
	suite.mockCacheRepo = new(mock.CacheRepository)
	suite.mockGen = new(genMock.Generator)
	cfg := infra.Config{}
	cfg.Server.Address = "localhost:8513"
	cfg.Shortener.CodeLength = 7
	suite.service = NewService(logrus.New(), &cfg, suite.mockRepo, suite.mockCacheRepo, suite.mockGen, infra.NOOPTelemetry)
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
		suite.mockGen.On("GenerateShortURLCode").Return(tc.expectedURL.ShortCode)
		suite.mockRepo.On("Create", context.TODO(), testifyMock.Anything).Return(nil)
		url, err := suite.service.CreateShortURL(context.TODO(), tc.input)

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
		suite.mockGen.On("GenerateShortURLCode").Return(tc.expectedURL.ShortCode)
		suite.mockRepo.On("Create", context.TODO(), testifyMock.Anything).Return(gorm.ErrInvalidData)
		url, err := suite.service.CreateShortURL(context.TODO(), tc.input)

		require.Error(err)
		require.Empty(url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_CreateShortURL_DoRetry_Success() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "shrtcd",
			},
		},
	}

	for _, tc := range testCases {
		suite.mockGen.On("GenerateShortURLCode").Return(tc.input.ShortCode)
		suite.mockRepo.On("Create", context.TODO(), testifyMock.Anything).Return(gorm.ErrDuplicatedKey).Once()
		suite.mockRepo.On("Create", context.TODO(), testifyMock.Anything).Return(nil).Once()
		url, err := suite.service.CreateShortURL(context.TODO(), tc.input.LongURL)

		require.NoError(err)
		require.NotEmpty(url)
	}
}

func (suite *URLServiceTestSuite) TestURLService_CreateShortURL_DoRetry_Failure() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				LongURL:   "http://google.com",
				ShortCode: "shrtcd",
			},
		},
	}

	for _, tc := range testCases {
		suite.mockGen.On("GenerateShortURLCode").Return(tc.input.ShortCode)
		suite.mockRepo.On("Create", context.TODO(), testifyMock.Anything).Return(gorm.ErrDuplicatedKey).Once()
		suite.mockRepo.On("Create", context.TODO(), testifyMock.Anything).Return(gorm.ErrInvalidData).Once()
		url, err := suite.service.CreateShortURL(context.TODO(), tc.input.LongURL)

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

	for _, tc := range testCases {
		suite.mockCacheRepo.On("Get", context.TODO(), tc.input).Return(&tc.expectedURL, nil).Once()
		suite.mockRepo.On("FindByShortCode", context.TODO(), testifyMock.Anything).Times(0)
		url, err := suite.service.GetLongURL(context.TODO(), tc.input)

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

	for _, tc := range testCases {
		suite.mockCacheRepo.On("Get", context.TODO(), tc.input).Return(nil, redis.Nil).Once()
		suite.mockRepo.On("FindByShortCode", context.TODO(), tc.input).Return(&tc.expectedURL, nil).Once()
		suite.mockCacheRepo.On("Set", context.TODO(), &tc.expectedURL).Return(nil)
		url, err := suite.service.GetLongURL(context.TODO(), tc.input)

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

	for _, tc := range testCases {
		suite.mockCacheRepo.On("Get", context.TODO(), tc.input).Return(nil, redis.Nil).Once()
		suite.mockRepo.On("FindByShortCode", context.TODO(), tc.input).Return(nil, gorm.ErrRecordNotFound).Once()
		url, err := suite.service.GetLongURL(context.TODO(), tc.input)

		require.Error(err)
		require.Empty(url)
	}
}

func TestURLServiceTestSuite(t *testing.T) {
	suite.Run(t, new(URLServiceTestSuite))
}
