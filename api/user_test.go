package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	mockdb "github.com/jrmarcco/go-backend-demo/db/mock"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type createUserMatcher struct {
	arg    db.CreateUserParams
	passwd string
}

func (m createUserMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPasswd(m.passwd, arg.HashedPasswd)
	if err != nil {
		return false
	}

	m.arg.HashedPasswd = arg.HashedPasswd
	return reflect.DeepEqual(m.arg, arg)
}

func (m createUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", m.arg, m.passwd)
}

func eqCreateUserParams(arg db.CreateUserParams, passwd string) gomock.Matcher {
	return createUserMatcher{
		arg:    arg,
		passwd: passwd,
	}
}

type mockSqlRes struct {
	affectedRows int64
	insertId     int64
}

func (res *mockSqlRes) LastInsertId() (int64, error) {
	return res.insertId, nil
}

func (res *mockSqlRes) RowsAffected() (int64, error) {
	return res.affectedRows, nil
}

func newMockSqlRes(affectedRows, insertId int64) *mockSqlRes {
	return &mockSqlRes{
		affectedRows: affectedRows,
		insertId:     insertId,
	}
}

func (a *apiTestSuite) TestCreateUserApi() {
	t := a.T()

	user, passwd := randomUser(t)

	tcs := []struct {
		name      string
		arg       gin.H
		buildStub func(store *mockdb.MockStore)
		checkResp func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Invalid Username Case",
			arg: gin.H{
				"username": "Foo#Bar",
				"password": passwd,
				"email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Email Case",
			arg: gin.H{
				"username": user.Username,
				"password": passwd,
				"email":    "invalid_email",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Short Password Case",
			arg: gin.H{
				"username": user.Username,
				"password": "123",
				"email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Normal Case",
			arg: gin.H{
				"username": user.Username,
				"password": passwd,
				"email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:     user.Username,
					HashedPasswd: passwd,
					Email:        user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, passwd)).
					Times(1).
					Return(newMockSqlRes(1, user.ID.Int64), nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				id, err := strconv.ParseInt(string(data), 10, 64)
				require.NoError(t, err)
				require.Equal(t, user.ID.Int64, id)
			},
		},
		{
			name: "Duplicate Username Case",
			arg: gin.H{
				"username": user.Username,
				"password": passwd,
				"email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:     user.Username,
					HashedPasswd: passwd,
					Email:        user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, passwd)).
					Times(1).
					Return(
						newMockSqlRes(0, 0),
						&mysql.MySQLError{Number: 1062, Message: "mock message"},
					)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Internal Error Case",
			arg: gin.H{
				"username": user.Username,
				"password": passwd,
				"email":    user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:     user.Username,
					HashedPasswd: passwd,
					Email:        user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, passwd)).
					Times(1).
					Return(newMockSqlRes(0, 0), sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stub
			tc.buildStub(store)

			// start server and send request
			server := a.newTestServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.arg)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/user/add", bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, req)
			tc.checkResp(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(8)

	hashed, err := util.HashPasswd(password)
	require.NoError(t, err)

	return db.User{
		ID: sql.NullInt64{
			Int64: util.RandomInt64(1, 1000),
			Valid: true,
		},
		Username:     util.RandomString(6),
		HashedPasswd: hashed,
		Email:        fmt.Sprintf("%s@email.com", util.RandomString(6)),
	}, password
}

func (a *apiTestSuite) TestLogin() {

}
