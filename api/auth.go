package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/natanaelrusli/go-bank/db/sqlc"
	"github.com/natanaelrusli/go-bank/token"
	"github.com/natanaelrusli/go-bank/util"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	var user db.User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	config, err := util.LoadConfig(".")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := tokenMaker.CreateToken(user.Username, time.Minute*15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginResponse{
		Token: token,
	}

	ctx.JSON(http.StatusOK, res)
}
