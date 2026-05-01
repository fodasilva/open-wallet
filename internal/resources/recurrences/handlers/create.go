package handlers

import (
	"net/http"
	"slices"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type CreateOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.RecurrencesUseCases

	UserID     string
	PassedKeys []string
	Body       CreateRecurrenceRequest
}

func (o *CreateOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	keys, err := httputil.GetJSONKeys(r)
	if err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	o.PassedKeys = keys

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *CreateOptions) Validate() error {
	if len(o.Body.Name) == 0 {
		return httputil.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if o.Body.Amount == 0 {
		return httputil.NewHTTPError(http.StatusBadRequest, "amount is required")
	}
	return nil
}

func (o *CreateOptions) Run() error {
	payload := repository.CreateRecurrenceDTO{
		UserID:      o.UserID,
		Name:        o.Body.Name,
		CategoryID:  util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "category_id"), Value: o.Body.CategoryID},
		Note:        util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "note"), Value: o.Body.Note},
		Amount:      o.Body.Amount,
		DayOfMonth:  o.Body.DayOfMonth,
		StartPeriod: o.Body.StartPeriod,
		EndPeriod:   util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "end_period"), Value: o.Body.EndPeriod},
	}

	res, err := o.UseCases.Create(o.R.Context(), payload)
	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusCreated, util.ResponseData[CreateRecurrenceResponseData]{
		Data: CreateRecurrenceResponseData{
			Recurrence: MapRecurrenceResource(res),
		},
	})

	return nil
}

// @Summary Create a recurrence
// @Description Create a new recurrence template
// @ID v1CreateRecurrence
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateRecurrenceRequest true "Recurrence payload"
// @Success 201 {object} util.ResponseData[CreateRecurrenceResponseData] "Recurrence created"
// @Failure 400 {object} httputil.HTTPError "Bad request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Router /api/v1/recurrences [post]
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	cmd := &CreateOptions{
		UseCases: api.recurrencesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
