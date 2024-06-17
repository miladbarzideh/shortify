package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/infra"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type URLRepositoryTestSuite struct {
	suite.Suite
	repo *Repository
	mock sqlmock.Sqlmock
}

func (suite *URLRepositoryTestSuite) SetupTest() {
	require := suite.Require()
	db, mock, err := sqlmock.New()
	require.NoError(err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(err)
	suite.repo = NewRepository(logrus.New(), gormDB, infra.NOOPTelemetry)
	suite.mock = mock
}

func (suite *URLRepositoryTestSuite) TestURLRepository_Create_Success() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				LongURL:   "https://google.com",
				ShortCode: "abcd",
			},
		},
	}

	for i, tc := range testCases {
		suite.mock.ExpectBegin()
		insertQuery := `INSERT INTO "urls" ("long_url","short_code","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`
		suite.mock.ExpectQuery(regexp.QuoteMeta(insertQuery)).
			WithArgs(tc.input.LongURL, tc.input.ShortCode, AnyTime{}, AnyTime{}).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))
		suite.mock.ExpectCommit()
		err := suite.repo.Create(context.TODO(), &tc.input)

		require.NoError(err)
		if err = suite.mock.ExpectationsWereMet(); err != nil {
			suite.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func (suite *URLRepositoryTestSuite) TestURLRepository_Create_FailedInsert_Failure() {
	require := suite.Require()
	testCases := []struct {
		input model.URL
	}{
		{
			input: model.URL{
				LongURL:   "https://google.com",
				ShortCode: "abcd",
			},
		},
	}

	for _, tc := range testCases {
		suite.mock.ExpectBegin()
		insertQuery := `INSERT INTO "urls" ("long_url","short_code","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`
		suite.mock.ExpectQuery(regexp.QuoteMeta(insertQuery)).
			WithArgs(tc.input.LongURL, tc.input.ShortCode, AnyTime{}, AnyTime{}).
			WillReturnError(errors.New("some err"))
		suite.mock.ExpectRollback()
		err := suite.repo.Create(context.TODO(), &tc.input)

		require.Error(err)
		if err = suite.mock.ExpectationsWereMet(); err != nil {
			suite.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func (suite *URLRepositoryTestSuite) TestURLRepository_FindByShortCode_Success() {
	require := suite.Require()
	testCases := []struct {
		input       string
		expectedURL model.URL
	}{
		{
			input: "A5rFt",
			expectedURL: model.URL{
				ID:        1,
				LongURL:   "https://google.com",
				ShortCode: "A5rFt",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		query := `SELECT \* FROM "urls" (.+)`
		rows := sqlmock.NewRows([]string{"id", "long_url", "short_code", "created_at", "updated_at"}).
			AddRow(tc.expectedURL.ID, tc.expectedURL.LongURL, tc.expectedURL.ShortCode, tc.expectedURL.CreatedAt, tc.expectedURL.UpdatedAt)
		suite.mock.ExpectQuery(query).WithArgs(tc.input, 1).WillReturnRows(rows)
		actualUrl, err := suite.repo.FindByShortCode(context.TODO(), tc.input)

		require.NoError(err)
		require.Equal(actualUrl.ShortCode, tc.expectedURL.ShortCode)
		if err = suite.mock.ExpectationsWereMet(); err != nil {
			suite.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func (suite *URLRepositoryTestSuite) TestURLRepository_FindByShortCode_Failure() {
	require := suite.Require()
	testCases := []struct {
		input string
	}{
		{
			input: "A5rFt",
		},
	}

	for _, tc := range testCases {
		query := `SELECT \* FROM "urls" (.+)`
		suite.mock.ExpectQuery(query).WithArgs(tc.input, 1).WillReturnError(gorm.ErrRecordNotFound)
		_, err := suite.repo.FindByShortCode(context.TODO(), tc.input)

		require.Error(err)
		require.Equal(gorm.ErrRecordNotFound, err)
		if err = suite.mock.ExpectationsWereMet(); err != nil {
			suite.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestURLRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(URLRepositoryTestSuite))
}
