package api

import (
	"github.com/gin-gonic/gin"
	"glossika/service/controller/merchandiseCtrl"
	"glossika/service/internal/config"
	"glossika/service/internal/errorx"
	boMerchandise "glossika/service/internal/model/bo/merchandise"
	"go.uber.org/dig"
	"net/http"
)

func NewMerchandise(pack merchandisePack) {
	m := &merchandise{pack: pack}
	group := pack.Root.Group("merchandise")
	{
		group.GET("recommendation", validateToken(pack.JWT), m.getRecommendation)
	}
}

type merchandisePack struct {
	dig.In

	Root            *gin.RouterGroup
	MerchandiseCtrl merchandiseCtrl.MerchandiseCtrl
	JWT             config.JWT
}

type merchandise struct {
	pack merchandisePack
}

// getRecommendation
//
//	@Summary	Retrieve recommended merchandise items
//	@Tags		merchandise
//	@version	1.0
//	@produce	json
//	@Param		Authorization	header		string	true	"Get token from login response"	default(Bearer {token})
//	@Success	200				{object}	boMerchandise.GetRecommendationReply
//	@Failure	400				{object}	errorx.ErrorResponse
//	@Router		/merchandise/recommendation [GET]
func (api *merchandise) getRecommendation(ctx *gin.Context) {
	recommendation, err := api.pack.MerchandiseCtrl.GetRecommendation(ctx, boMerchandise.GetRecommendationArgs{})
	if err != nil {
		errorx.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, recommendation)
}
