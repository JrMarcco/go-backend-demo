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
		ctx.JSON(http.StatusBadRequest, ErrorResp(err))
		return
	}

	hashedPasswd, err := util.HashPasswd(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
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
				ctx.JSON(http.StatusForbidden, ErrorResp(mysqlErr))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	id, _ := res.LastInsertId()
	ctx.JSON(http.StatusOK, id)
}

type getUserReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type userResp struct {
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
		ctx.JSON(http.StatusBadRequest, ErrorResp(err))
		return
	}

	user, err := s.store.GetUser(ctx, sql.NullInt64{Int64: req.ID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResp(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	c, err := copier.NewRefCopier[db.User, userResp]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	resp, err := c.Copy(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResp struct {
	AccessToken string    `json:"accessToken"`
	User        *userResp `json:"user"`
}

func (s *Server) login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResp(err))
		return
	}

	user, err := s.store.FindUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	if err = util.CheckPasswd(req.Password, user.HashedPasswd); err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResp(err))
		return
	}

	token, err := s.GenerateToken(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	c, err := copier.NewRefCopier[db.User, userResp]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	ur, err := c.Copy(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK,
		loginResp{
			AccessToken: token,
			User:        ur,
		},
	)
}
