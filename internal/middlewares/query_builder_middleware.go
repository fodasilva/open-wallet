package middlewares

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func QueryBuilderMiddleware(config querybuilder.ParseConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			results, err := querybuilder.ParseRequest(
				r.URL.Query().Get("filter"),
				r.URL.Query().Get("page"),
				r.URL.Query().Get("per_page"),
				r.URL.Query().Get("order_by"),
				config,
			)
			if err != nil {
				apiErr := util.NewHTTPError(http.StatusBadRequest, err.Error())
				httputil.JSON(w, apiErr.StatusCode, apiErr)
				return
			}

			ctx := context.WithValue(r.Context(), util.ContextKeyPage, results.Page)
			ctx = context.WithValue(ctx, util.ContextKeyPerPage, results.PerPage)
			ctx = querybuilder.WithBuilder(ctx, results.Builder)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
