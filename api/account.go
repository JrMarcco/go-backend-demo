package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/token"
	"net/http"
)

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *S) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	payload := ctx.MustGet(payloadKey).(*token.Payload)
	res, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		AccountOwner: payload.Username,
		Currency:     req.Currency,
		Balance:      0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	id, _ := res.LastInsertId()

	ctx.JSON(http.StatusOK, id)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *S) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	account, err := s.store.GetAccount(ctx, sql.NullInt64{Int64: req.ID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	payload := ctx.MustGet(payloadKey).(*token.Payload)
	if account.AccountOwner != payload.Username {
		err := errors.New("account doesn't belong to authorized user")
		ctx.JSON(http.StatusUnauthorized, errorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountReq struct {
	PageNo   int32 `form:"pageNo" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=100"`
}

func (s *S) listAccount(ctx *gin.Context) {
	var req listAccountReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	accounts, err := s.store.ListAccount(ctx, db.ListAccountParams{
		Offset: (req.PageNo - 1) * req.PageSize,
		Limit:  req.PageSize,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
