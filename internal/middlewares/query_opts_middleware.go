package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/utils"
)

func QueryOptsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, _ := ctx.GetQuery("page")
		perPage, _ := ctx.GetQuery("per_page")
		orderBy, _ := ctx.GetQuery("order_by")
		pageNum, err := strconv.Atoi(page)

		if err != nil {
			pageNum = 1
		}

		perPageNum, err := strconv.Atoi(perPage)

		if err != nil {
			perPageNum = 10
		}

		queryOpts := utils.QueryOpts()
		queryOpts.Offset((pageNum - 1) * perPageNum)
		queryOpts.Limit(perPageNum + 1)

		if orderBy != "" {
			splittedByComma := strings.Split(orderBy, ",")

			for _, field := range splittedByComma {
				field = strings.TrimSpace(field)

				if field == "" {
					continue
				}

				splittedByColon := strings.Split(field, ":")
				fieldName := strings.TrimSpace(splittedByColon[0])

				if fieldName == "" {
					apiErr := utils.NewHTTPError(http.StatusBadRequest, "invalid order_by: empty field name")
					ctx.JSON(apiErr.StatusCode, apiErr)
					ctx.Abort()
					return
				}

				direction := "asc"
				if len(splittedByColon) > 1 {
					direction = strings.ToLower(strings.TrimSpace(splittedByColon[1]))

					if direction != "asc" && direction != "desc" {
						apiErr := utils.NewHTTPError(http.StatusBadRequest, "invalid order: must be 'asc' or 'desc'")
						ctx.JSON(apiErr.StatusCode, apiErr)
						ctx.Abort()
						return
					}
				}

				queryOpts.OrderBy(fieldName, direction)
			}
		}

		filter, _ := ctx.GetQuery("filter")

		if filter != "" {
			queryOpts, err = parseFilter(filter, queryOpts)
			apiErr := utils.GetApiErr(err)
			if err != nil {
				ctx.JSON(apiErr.StatusCode, apiErr)
				ctx.Abort()
				return
			}
		}

		ctx.Set("page", pageNum)
		ctx.Set("per_page", perPageNum)
		ctx.Set("query_opts", queryOpts)
		ctx.Next()
	}
}

func parseFilter(filter string, query *utils.QueryOptsBuilder) (*utils.QueryOptsBuilder, error) {
	splitted := strings.Split(filter, " and ")

	allowedOperators := map[string]bool{
		"eq":   true,
		"ne":   true,
		"gt":   true,
		"gte":  true,
		"lt":   true,
		"lte":  true,
		"like": true,
	}

	for _, filter := range splitted {
		isOrGroup := strings.HasPrefix(filter, "(") && strings.HasSuffix(filter, ")")

		if isOrGroup {
			filter = strings.TrimPrefix(filter, "(")
			filter = strings.TrimSuffix(filter, ")")

			splittedOrGroup := strings.Split(filter, " or ")

			orQuery := query.InitOr()

			for _, filter := range splittedOrGroup {
				splittedFilter, err := splitFilter(filter)

				if err != nil {
					return nil, utils.NewHTTPError(http.StatusBadRequest, err.Error())
				}
				if len(splittedFilter) != 3 {
					return nil, utils.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("filter sintax error at '%s'", filter))
				}

				if !allowedOperators[splittedFilter[1]] {
					return nil, utils.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("operator not '%s' not allowed at '%s'", splittedFilter[1], filter))
				}

				condValue, err := parseFilterValue(splittedFilter[2])

				if err != nil {
					return nil, utils.NewHTTPError(http.StatusBadRequest, err.Error())
				}

				orQuery.Or(splittedFilter[0], splittedFilter[1], condValue)
			}

			orQuery.EndOr()
		} else {
			splittedFilter, err := splitFilter(filter)

			if err != nil {
				return nil, utils.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			if len(splittedFilter) != 3 {
				return nil, utils.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("filter sintax error at '%s'", filter))
			}

			if !allowedOperators[splittedFilter[1]] {
				return nil, utils.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("operator not '%s' not allowed at '%s'", splittedFilter[1], filter))
			}

			condValue, err := parseFilterValue(splittedFilter[2])

			if err != nil {
				return nil, utils.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			query.And(splittedFilter[0], splittedFilter[1], condValue)
		}
	}

	return query, nil
}

func parseFilterValue(value string) (any, error) {
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return strings.Trim(value, "'"), nil
	}

	if value == "true" || value == "false" {
		return value == "true", nil
	}

	if value == "null" {
		return nil, nil
	}

	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	return nil, fmt.Errorf("invalid value %s", value)
}

func splitFilter(filter string) ([]string, error) {
	splitted := strings.Split(filter, " ")

	if len(splitted) < 3 {
		return nil, fmt.Errorf("filter sintax error at '%s'", filter)
	}

	return []string{splitted[0], splitted[1], strings.Join(splitted[2:], " ")}, nil
}
