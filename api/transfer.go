package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"net/http"
)

type createTransferReq struct {
	FromID   int64  `json:"fromID" binding:"required,min=1"`
	ToID     int64  `json:"toID" binding:"required,min=1"`
	Amount   int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency""`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req createTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	if !s.validAccount(ctx, req.FromID, req.Currency) || !s.validAccount(ctx, req.ToID, req.Currency) {
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

func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, sql.NullInt64{Int64: accountID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, errorResp(
			fmt.Errorf("account [%d] mismatch: %s vs %s", accountID, account.Currency, currency),
		))
		return false
	}

	return true
}
