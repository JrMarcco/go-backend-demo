package api

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	mockdb "github.com/jrmarcco/go-backend-demo/db/mock"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/util"
	"testing"
	"time"
)

func TestGetAccountApi(t *testing.T) {

	account := randAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
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
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
