package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/jrmarcco/go-backend-demo/db/mock"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func (a *apiTestSuite) buildGetAccountReq(userID int64) *http.Request {
	t := a.T()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/account/get/%d", userID), nil)
	require.NoError(t, err)

	return req
}

func (a *apiTestSuite) TestGetAccountApi() {
	t := a.T()

	account := randAccount()

	tcs := []struct {
		name      string
		req       *http.Request
		username  string
		buildStub func(store *mockdb.MockStore)
		checkResp func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "Normal Case",
			req:      a.buildGetAccountReq(account.ID.Int64),
			username: account.AccountOwner,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var respAccount db.Account
				err = json.Unmarshal(data, &respAccount)
				require.NoError(t, err)
				require.Equal(t, account, respAccount)
			},
		},
		{
			name:     "Invalid ID Case",
			req:      a.buildGetAccountReq(0),
			username: account.AccountOwner,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "NotFound Case",
			req:      a.buildGetAccountReq(account.ID.Int64),
			username: account.AccountOwner,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError Case",
			req:      a.buildGetAccountReq(account.ID.Int64),
			username: account.AccountOwner,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "Unauthorized User Case",
			req:      a.buildGetAccountReq(account.ID.Int64),
			username: "unauthorized user",
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "Unauthorized Case",
			req:      a.buildGetAccountReq(0),
			username: "",
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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
			a.setupTestServer(store)
			recorder := httptest.NewRecorder()

			a.setAuthorization(tc.req, tc.username, time.Minute)

			a.s.router.ServeHTTP(recorder, tc.req)
			tc.checkResp(t, recorder)
		})
	}
}

func randAccount() db.Account {
	return db.Account{
		ID: sql.NullInt64{
			Int64: util.RandomInt64(1, 100),
			Valid: true,
		},
		AccountOwner: util.RandomString(6),
		Balance:      util.RandomInt64(1000, 10000),
		Currency:     "RMB",
	}
}
