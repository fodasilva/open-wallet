package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func QueryBuilderMiddleware(config querybuilder.ParseConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results, err := querybuilder.ParseRequest(
			ctx.Query("filter"),
			ctx.Query("page"),
			ctx.Query("per_page"),
			ctx.Query("order_by"),
			config,
		)
		if err != nil {
			apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
			ctx.JSON(apiErr.StatusCode, apiErr)
			ctx.Abort()
			return
		}

		ctx.Set("page", results.Page)
		ctx.Set("per_page", results.PerPage)
		ctx.Set("query_builder", results.Builder)
		ctx.Next()
	}
}
