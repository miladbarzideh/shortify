package repository

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/miladbarzideh/shortify/infra"
	"github.com/miladbarzideh/shortify/internal/model"
)

type URLCacheRepositoryTestSuite struct {
	suite.Suite
	cacheRepo URLCacheRepository
	cacheMock redismock.ClientMock
}

func (suite *URLCacheRepositoryTestSuite) SetupTest() {
	db, mock := redismock.NewClientMock()
	tracer := infra.NOOPTelemetry.TraceProvider.Tracer("")
	suite.cacheRepo = NewCacheRepository(logrus.New(), db, tracer)
	suite.cacheMock = mock
}

func (suite *URLCacheRepositoryTestSuite) TestURLCacheRepository_Set_Success() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				ID:        1,
				LongURL:   "https://google.com",
				ShortCode: "A5rFt",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		value, _ := json.Marshal(tc.input)
		suite.cacheMock.ExpectSet(suite.cacheRepo.BuildKeyWithPrefix(tc.input.ShortCode), value, 24*time.Hour).SetVal(string(value))
		err := suite.cacheRepo.Set(context.TODO(), tc.input)

		require.Nil(err)
	}
}

func (suite *URLCacheRepositoryTestSuite) TestURLCacheRepository_Set_Failure() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				ID:        1,
				LongURL:   "https://google.com",
				ShortCode: "A5rFt",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		value, _ := json.Marshal(tc.input)
		suite.cacheMock.ExpectSet(suite.cacheRepo.BuildKeyWithPrefix(tc.input.ShortCode), value, 24*time.Hour).SetErr(errors.New("FAIL"))
		err := suite.cacheRepo.Set(context.TODO(), tc.input)

		require.NotNil(err)
	}
}

func (suite *URLCacheRepositoryTestSuite) TestURLCacheRepository_Get_Success() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				ID:        1,
				LongURL:   "https://google.com",
				ShortCode: "A5rFt",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		value, _ := json.Marshal(tc.input)
		suite.cacheMock.ExpectGet(suite.cacheRepo.BuildKeyWithPrefix(tc.input.ShortCode)).SetVal(string(value))
		actualURL, err := suite.cacheRepo.Get(context.TODO(), tc.input.ShortCode)

		require.Nil(err)
		require.Equal(actualURL.LongURL, tc.input.LongURL)
	}
}

func (suite *URLCacheRepositoryTestSuite) TestURLCacheRepository_Get_Failure() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				ID:        1,
				LongURL:   "https://google.com",
				ShortCode: "A5rFt",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		suite.cacheMock.ExpectGet(suite.cacheRepo.BuildKeyWithPrefix(tc.input.ShortCode)).SetErr(errors.New("nil"))
		_, err := suite.cacheRepo.Get(context.TODO(), tc.input.ShortCode)

		require.NotNil(err)
	}
}

func TestCacheRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(URLCacheRepositoryTestSuite))
}
