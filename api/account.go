package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "go-backend-demo/db/sqlc"
	"net/http"
)

type createAccountReq struct {
	AccountOwner string `json:"accountOwner" binding:"required"`
	Currency     string `json:"currency" binding:"required,oneof=RMB USD"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	res, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		AccountOwner: req.AccountOwner,
		Currency:     req.Currency,
		Balance:      0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, id)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, account)
}

type listAccountReq struct {
	PageNo   int32 `form:"pageNo" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=100"`
}

func (s *Server) listAccount(ctx *gin.Context) {
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
