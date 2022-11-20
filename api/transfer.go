package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/token"
	"net/http"
)

type createTransferReq struct {
	FromID   int64  `json:"fromID" binding:"required,min=1"`
	ToID     int64  `json:"toID" binding:"required,min=1"`
	Amount   int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency""`
}

func (s *S) createTransfer(ctx *gin.Context) {
	var req createTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	fromAccount, valid := s.validAccount(ctx, req.FromID, req.Currency)
	if !valid {
		return
	}

	payload := ctx.MustGet(payloadKey).(*token.Payload)
	if fromAccount.AccountOwner != payload.Username {
		err := errors.New("from account doesn't belong to authorized user")
		ctx.JSON(http.StatusUnauthorized, errorResp(err))
		return
	}

	_, valid = s.validAccount(ctx, req.ToID, req.Currency)
	if !valid {
		return
	}

	res, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromID: req.FromID,
		ToID:   req.ToID,
		Amount: req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *S) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, sql.NullInt64{Int64: accountID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return account, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, errorResp(
			fmt.Errorf("account [%d] mismatch: %s vs %s", accountID, account.Currency, currency),
		))
		return account, false
	}

	return account, true
}
