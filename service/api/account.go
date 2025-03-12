package api

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"meme_coin_api/service/controller/accountCtrl"
	"meme_coin_api/service/internal/errorx"
	boAccount "meme_coin_api/service/internal/model/bo/account"
	"net/http"
	"strconv"
)

func NewAccount(pack accountPack) {
	m := &account{pack: pack}
	group := pack.Root.Group("account")
	{
		group.POST("register", m.createAccount)
		group.GET("verify", m.getAccount)
		group.POST("login", m.pokeAccount)
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

// getAccount
//
//	@Summary	Retrieve data for a meme coin
//	@Tags		meme_coin
//	@version	1.0
//	@produce	json
//	@Param		id				path		int	true	"id = ?"
//	@Success	200				{object}	boAccount.GetReply
//	@Failure	400				{object}	errorx.ErrorResponse
//	@Router		/meme_coin/{id}	[GET]
func (api *account) getAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}

	args := boAccount.GetArgs{
		ID: int64(id),
	}
	result, err := api.pack.AccountCtrl.Get(ctx, args)
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// createAccount
//
//	@Summary	Create a meme coin
//	@Tags		meme_coin
//	@version	1.0
//	@produce	json
//	@Param		request	body		api.createAccount.body	true	"Request payload for creating a meme coin"
//	@Success	200		{object}	boAccount.CreateReply
//	@Failure	400		{object}	errorx.ErrorResponse
//	@Router		/meme_coin [POST]
func (api *account) createAccount(ctx *gin.Context) {
	type body struct {
		Name        string `valid:"required" json:"name"`
		Description string `valid:"-" json:"description"`
	}
	var reqBody body
	if err := ctx.BindJSON(&reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusInternalServerError, err)
		return
	}
	if _, err := govalidator.ValidateStruct(reqBody); err != nil || govalidator.HasWhitespaceOnly(reqBody.Name) {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}

	args := boAccount.CreateArgs{
		Name:        reqBody.Name,
		Description: reqBody.Description,
	}
	reply, err := api.pack.AccountCtrl.Create(ctx, args)
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, reply)
}

// updateAccount
//
//	@Summary	Update a meme coin
//	@Tags		meme_coin
//	@version	1.0
//	@produce	json
//	@Param		id		path	int						true	"id = ?"
//	@Param		request	body	api.updateAccount.body	true	"Payload to update the description of a meme coin"
//	@Success	200
//	@Failure	400	{object}	errorx.ErrorResponse
//	@Router		/meme_coin/{id} [PATCH]
func (api *account) updateAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}
	type body struct {
		Description string `valid:"-" json:"description"`
	}
	var reqBody body
	if err = ctx.BindJSON(&reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusInternalServerError, err)
		return
	}
	if _, err = govalidator.ValidateStruct(reqBody); err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}

	args := boAccount.UpdateArgs{
		ID:          int64(id),
		Description: reqBody.Description,
	}
	if err = api.pack.AccountCtrl.Update(ctx, args); err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.Status(http.StatusOK)
}

// deleteAccount
//
//	@Summary	Delete a meme coin
//	@Tags		meme_coin
//	@version	1.0
//	@produce	json
//	@Param		id	path	int	true	"id = ?"
//	@Success	200
//	@Failure	400				{object}	errorx.ErrorResponse
//	@Router		/meme_coin/{id}	[DELETE]
func (api *account) deleteAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}

	args := boAccount.DeleteArgs{
		ID: int64(id),
	}
	if err = api.pack.AccountCtrl.Delete(ctx, args); err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.Status(http.StatusOK)
}

// pokeAccount
//
//	@Summary	Poke a meme coin
//	@Tags		meme_coin
//	@version	1.0
//	@produce	json
//	@Param		id	path	int	true	"id = ?"
//	@Success	200
//	@Failure	400						{object}	errorx.ErrorResponse
//	@Router		/meme_coin/{id}/poke	[POST]
func (api *account) pokeAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, errorx.ParamKeyRequired)
		return
	}

	args := boAccount.IncreasePopularityScoreArgs{
		ID: int64(id),
	}
	if err = api.pack.AccountCtrl.IncreasePopularityScore(ctx, args); err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.Status(http.StatusOK)
}
