package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type PrepareOptions struct {
	Ctx      *gin.Context
	UseCases usecases.RecurrencesUseCases

	UserID string
	Period string
}

func (o *PrepareOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.Period = ctx.Param("period")

	return nil
}

func (o *PrepareOptions) Validate() error {
	if len(o.Period) != 6 {
		return utils.NewHTTPError(http.StatusBadRequest, "invalid period format. Expected YYYYMM.")
	}
	return nil
}

func (o *PrepareOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.Ctx.Request.Context(), "RecurrencesHandler.Prepare")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", o.UserID))

	err := o.UseCases.PrepareRecurrences(tCtx, o.UserID, o.Period)
	if err != nil {
		span.RecordError(err)
		return err
	}

	o.Ctx.Status(http.StatusNoContent)
	return nil
}

// @Summary Prepare recurrences for a period
// @Description Generates entry records for all recurrence templates that don't already have one in the given period.
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "Period in YYYYMM format (e.g. 202603)"
// @Success 204 "Recurrences prepared"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences/{period} [post]
func (api *API) Prepare(ctx *gin.Context) {
	cmd := &PrepareOptions{
		UseCases: api.recurrencesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
