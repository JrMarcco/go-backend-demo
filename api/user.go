package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	db "github.com/jrmarcco/go-backend-demo/db/sqlc"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/jrmarcco/go-backend-demo/util/copier"
	"net/http"
	"time"
)

type createUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	hashedPasswd, err := util.HashPasswd(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	res, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		HashedPasswd: hashedPasswd,
	})

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case uint16(1062):
				ctx.JSON(http.StatusForbidden, errorResp(mysqlErr))
				return
			}
		}
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

type getUserReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserResp struct {
	ID                sql.NullInt64 `json:"id"`
	Username          string        `json:"username"`
	Email             string        `json:"email"`
	PasswordChangedAt time.Time     `json:"passwordChangedAt"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResp(err))
		return
	}

	user, err := s.store.GetUser(ctx, sql.NullInt64{Int64: req.ID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResp(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	c, err := copier.NewRefCopier[db.User, getUserResp]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	resp, err := c.Copy(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
