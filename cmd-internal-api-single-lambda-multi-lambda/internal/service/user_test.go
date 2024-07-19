package service

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jha-captech/user-microservice/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuit struct {
	suite.Suite
	service *User
	dbMock  sqlmock.Sqlmock
}

func TestTestSuit(t *testing.T) {
	suite.Run(t, new(testSuit))
}

func (s *testSuit) SetupSuite() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	s.dbMock = mock
	s.service = NewUser(db)
}

func (s *testSuit) TearDownSuite() {
	_ = s.service.database.Close()
}

func (s *testSuit) TestListUsers() {
	t := s.T()

	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		expectedReturn []models.User
		expectedError  error
	}{
		"Return slice of users": {
			mockReturn:     mustStructsToRows(users),
			mockReturnErr:  nil,
			expectedReturn: users,
			expectedError:  nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `SELECT * FROM "users"`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.ListUsers()

			assert.Equal(t, tc.expectedError, err, "errors did not match")
			assert.Equal(t, tc.expectedReturn, actualReturn, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestFetchUser() {
	t := s.T()

	user := models.User{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001}

	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		inputID        int
		expectedReturn models.User
		expectedError  error
	}{
		"Return usr by ID": {
			mockReturn:     mustStructsToRows([]models.User{user}),
			mockReturnErr:  nil,
			inputID:        int(user.ID),
			expectedReturn: user,
			expectedError:  nil,
		},
		"User with given ID does not exist": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  sql.ErrNoRows,
			inputID:        0,
			expectedReturn: models.User{},
			expectedError:  fmt.Errorf("[in FetchUser]:, %w", sql.ErrNoRows),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WithArgs(tc.inputID).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.FetchUser(tc.inputID)

			assert.Equal(t, tc.expectedError, err, "errors did not match")
			assert.Equal(t, tc.expectedReturn, actualReturn, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestUpdateUser() {
	t := s.T()

	userIn := models.User{ID: 0, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001}
	userOut := models.User{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001}

	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		inputID        int
		inputUser      models.User
		expectedReturn models.User
		expectedError  error
	}{
		"user updated by ID": {
			mockInputArgs:  []driver.Value{userIn.FirstName, userIn.LastName, userIn.Role, userIn.UserID, userIn.ID},
			mockReturn:     mustStructsToRows([]models.User{userIn}),
			mockReturnErr:  nil,
			inputID:        int(userIn.ID),
			inputUser:      userIn,
			expectedReturn: userOut,
			expectedError:  nil,
		},
		"Error updating user": {
			mockInputArgs:  []driver.Value{},
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  sql.ErrNoRows,
			inputID:        0,
			inputUser:      userIn,
			expectedReturn: models.User{},
			expectedError:  fmt.Errorf("[in FetchUser]:, %w", sql.ErrNoRows),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `
				UPDATE
					"users"
				SET
					"first_name" = $1,
					"last_name" = $2,
					"role" = $3,
					"user_id" = $4
				WHERE
					"id" = $5
			`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.UpdateUser(tc.inputID, tc.inputUser)

			assert.Equal(t, tc.expectedError, err, "errors did not match")
			assert.Equal(t, tc.expectedReturn, actualReturn, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestCreateUser() {
	t := s.T()

	userIn := models.User{ID: 0, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001}

	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		inputUser      models.User
		expectedReturn int
		expectedError  error
	}{
		"create": {
			mockInputArgs:  []driver.Value{userIn.FirstName, userIn.LastName, userIn.Role, userIn.UserID},
			mockReturn:     mustStructsToRows([]struct{ ID int }{{ID: 1}}),
			mockReturnErr:  nil,
			inputUser:      userIn,
			expectedReturn: 1,
			expectedError:  nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `
				INSERT INTO "users" ("first_name", "last_name", "role", "user_id")
					VALUES ($1, $2, $3, $4)
				RETURNING "id"
			`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.CreateUser(tc.inputUser)

			assert.Equal(t, tc.expectedError, err, "errors did not match")
			assert.Equal(t, tc.expectedReturn, actualReturn, "returned data does not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestDeleteUser() {
	t := s.T()

	testCases := map[string]struct {
		mockInputArgs []driver.Value
		mockReturn    driver.Result
		mockReturnErr error
		inputID       int
		expectedError error
	}{
		"create": {
			mockInputArgs: []driver.Value{1},
			mockReturn:    sqlmock.NewResult(1, 1),
			mockReturnErr: nil,
			inputID:       1,
			expectedError: nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `
				DELETE FROM "users"
				WHERE "users"."id" = $1
			`
			s.dbMock.
				ExpectExec(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...).
				WillReturnResult(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			err := s.service.DeleteUser(tc.inputID)

			assert.Equal(t, tc.expectedError, err, "errors did not match")

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

// ━━ HELPERS ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

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
