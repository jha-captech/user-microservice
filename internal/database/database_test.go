package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"

	"user-microservice/internal/database/entity"
	"user-microservice/internal/testutil"
)

// TEST SETUP

type databaseSuit struct {
	suite.Suite
	session Database
	db      *sql.DB
	dbMock  sqlmock.Sqlmock
}

func TestDatabaseSuit(t *testing.T) {
	suite.Run(t, new(databaseSuit))
}

func (s *databaseSuit) SetupSuite() {
	db, mock, _ := sqlmock.New()

	dbSession := MustNewDatabase(
		postgres.New(
			postgres.Config{
				Conn:       db,
				DriverName: "postgres",
			},
		),
		WithRetryCount(5),
		WithAutoMigrate(false),
	)

	s.db = db
	s.session = dbSession
	s.dbMock = mock
}

func (s *databaseSuit) TearDownSuite() {
	_ = s.db.Close()
}

// TESTS

func (s *databaseSuit) TestListUsers() {
	t := s.T()

	testCases := map[string]struct {
		expectedReturn []entity.User
		expectedError  error
	}{
		"Good": {
			expectedReturn: []entity.User{
				testutil.NewUser(),
				testutil.NewUser(),
				testutil.NewUser(),
			},
			expectedError: nil,
		},
		// "Bad": {}, // TODO: add bad case for test
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rows := mustStructsToRows(tc.expectedReturn)

			query := regexp.QuoteMeta(`SELECT * FROM "users"`)
			s.dbMock.
				ExpectQuery(query).
				WillReturnRows(rows)

			users, err := s.session.ListUsers()

			assert.Equal(t, tc.expectedError, err, "error in ListUsers")
			assert.Equal(t, tc.expectedReturn, users, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err, "mock expectations not met")
		})
	}
}

func (s *databaseSuit) TestFetchUser() {
	t := s.T()

	testUser := testutil.NewUser(testutil.WithID(1))

	testCases := map[string]struct {
		mockReturnRows *sqlmock.Rows
		expectedReturn entity.User
		userID         int
		expectedError  error
	}{
		"Good": {
			mockReturnRows: mustStructsToRows([]entity.User{testUser}),
			expectedReturn: testUser,
			userID:         int(testUser.ID),
			expectedError:  nil,
		},
		"Bad": {
			mockReturnRows: mustStructToEmptyRow(entity.User{}),
			expectedReturn: entity.User{},
			userID:         2,
			expectedError:  nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			query := regexp.QuoteMeta(
				`SELECT * FROM "users" WHERE ID = $1 ORDER BY "users"."id" LIMIT $2`,
			)
			s.dbMock.
				ExpectQuery(query).
				WithArgs(tc.userID, 1).
				WillReturnRows(tc.mockReturnRows)

			users, err := s.session.FetchUser(tc.userID)

			assert.Equal(t, tc.expectedError, err, "error in FetchUser")
			assert.Equal(t, tc.expectedReturn, users, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err, "mock expectations not met")
		})
	}
}

// HELPERS

// structSliceToSQLMockRows converts a slice of structs to sqlmock.Rows using reflect.
// It can also be used when only a single struct is needed by wrapping in a slice.
func mustStructsToRows[T any](slice []T) *sqlmock.Rows {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected a slice but got %T", slice))
	}

	if v.Len() == 0 {
		panic("slice is empty")
	}

	elemType := reflect.TypeOf(slice).Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected a slice of structs but got a slice of %v", elemType.Kind()))
	}

	numFields := elemType.NumField()
	columns := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		colName := elemType.Field(i).Name
		colNameSnake := toSnake(colName)
		columns[i] = colNameSnake
	}

	rows := sqlmock.NewRows(columns)

	for i := 0; i < v.Len(); i++ {
		var values []driver.Value
		elem := v.Index(i)
		for j := 0; j < elem.NumField(); j++ {
			values = append(values, elem.Field(j).Interface())
		}
		rows.AddRow(values...)
	}

	return rows
}

// mustStructToEmptyRow converts a struct into an *sqlmock.Rows object with headers but no rows.
func mustStructToEmptyRow[T any](obj T) *sqlmock.Rows {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected a struct but got %T", obj))
	}

	elemType := v.Type()
	numFields := v.NumField()
	columns := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		colName := elemType.Field(i).Name
		colNameSnake := toSnake(colName)
		columns[i] = colNameSnake
	}

	return sqlmock.NewRows(columns)
}

// toSnake converts PascalCase to snake_case with special handling for abbreviations
func toSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		if (i != 0 || i == l-1) &&
			((i > 0 && rune(camel[i-1]) >= 'a') || (i < l-1 && rune(camel[i+1]) >= 'a')) {
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}
