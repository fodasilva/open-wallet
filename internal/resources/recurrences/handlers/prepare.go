package handlers

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type PrepareOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.RecurrencesUseCases

	UserID string
	Period string
}

func (o *PrepareOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	o.Period = r.PathValue("period")

	return nil
}

func (o *PrepareOptions) Validate() error {
	if len(o.Period) != 6 {
		return httputil.NewHTTPError(http.StatusBadRequest, "invalid period format. Expected YYYYMM.")
	}
	return nil
}

func (o *PrepareOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.R.Context(), "RecurrencesHandler.Prepare")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", o.UserID))

	err := o.UseCases.PrepareRecurrences(tCtx, o.UserID, o.Period)
	if err != nil {
		span.RecordError(err)
		return err
	}

	o.W.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary Prepare recurrences for a period
// @Description Generates entry records for all recurrence templates that don't already have one in the given period.
// @ID v1PrepareRecurrence
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "Period in YYYYMM format (e.g. 202603)"
// @Success 204 "Recurrences prepared"
// @Failure 400 {object} httputil.HTTPError "Bad request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Failure 503 {string} string "Service Unavailable"
// @Router /api/v1/recurrences/{period} [post]
func (api *API) Prepare(w http.ResponseWriter, r *http.Request) {
	cmd := &PrepareOptions{
		UseCases: api.recurrencesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
