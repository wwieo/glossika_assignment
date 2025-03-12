package api

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"glossika/service/controller/accountCtrl"
	"glossika/service/internal/errorx"
	boAccount "glossika/service/internal/model/bo/account"
	"go.uber.org/dig"
	"net/http"
	"regexp"
)

func NewAccount(pack accountPack) {
	m := &account{pack: pack}
	group := pack.Root.Group("account")
	{
		group.GET("verify", m.register)
		group.POST("register", m.register)
		group.POST("login", m.login)
	}
}

type accountPack struct {
	dig.In

	Root        *gin.RouterGroup
	AccountCtrl accountCtrl.AccountCtrl
}

type account struct {
	pack accountPack
}

// register
//
//	@Summary	Register an account
//	@Tags		account
//	@version	1.0
//	@produce	json
//	@Param		request	body		api.register.body	true	"Request payload for creating a account"
//	@Success	200		{object}	boAccount.RegisterReply
//	@Failure	400		{object}	errorx.ErrorResponse
//	@Router		/account/register [POST]
func (api *account) register(ctx *gin.Context) {
	type body struct {
		Email    string `valid:"required" json:"email"`
		Password string `valid:"required" json:"password"`
	}
	var reqBody body
	if err := ctx.BindJSON(&reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusInternalServerError, err)
		return
	}
	if _, err := govalidator.ValidateStruct(reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}
	if !govalidator.IsEmail(reqBody.Email) {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.InvalidEmailFormat)
		return
	}
	if !api.validatePassword(reqBody.Password) {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.InvalidPasswordFormat)
		return
	}

	reply, err := api.pack.AccountCtrl.Register(ctx, boAccount.RegisterArgs{
		Account: boAccount.Account{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		},
	})
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, reply)
}

// login
//
//	@Summary	Login
//	@Tags		account
//	@version	1.0
//	@produce	json
//	@Param		request	body		api.login.body	true	"Request payload for login"
//	@Success	200		{object}	boAccount.LoginReply
//	@Failure	400		{object}	errorx.ErrorResponse
//	@Router		/account/login [POST]
func (api *account) login(ctx *gin.Context) {
	type body struct {
		Email    string `valid:"required" json:"email"`
		Password string `valid:"required" json:"password"`
	}
	var reqBody body
	if err := ctx.BindJSON(&reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusInternalServerError, err)
		return
	}
	if _, err := govalidator.ValidateStruct(reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}
	if !govalidator.IsEmail(reqBody.Email) {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.InvalidEmailFormat)
		return
	}
	if !api.validatePassword(reqBody.Password) {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.InvalidPasswordFormat)
		return
	}

	args := boAccount.LoginArgs{
		boAccount.Account{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		},
	}
	reply, err := api.pack.AccountCtrl.Login(ctx, args)
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, reply)
}

func (api *account) validatePassword(password string) bool {
	if len(password) < 6 || len(password) > 16 {
		return false
	}

	upperRegex := `[A-Z]`
	lowerRegex := `[a-z]`
	specialCharRegex := `[()\[\]{}<>+\-*/?,.:;"'_|\~` + "`" + `!@#$%^&=]`

	if !regexp.MustCompile(upperRegex).MatchString(password) {
		return false
	}

	if !regexp.MustCompile(lowerRegex).MatchString(password) {
		return false
	}

	if !regexp.MustCompile(specialCharRegex).MatchString(password) {
		return false
	}

	return true
}
