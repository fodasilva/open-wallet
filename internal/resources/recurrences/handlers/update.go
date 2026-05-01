package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type UpdateOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.RecurrencesUseCases

	ID         string
	UserID     string
	PassedKeys []string
	Body       UpdateRecurrenceRequest
}

func (o *UpdateOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.ID = r.PathValue("id")
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	keys, err := util.GetJSONKeys(r)
	if err != nil {
		return util.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	o.PassedKeys = keys

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return util.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *UpdateOptions) Validate() error {
	if len(o.PassedKeys) == 0 {
		return util.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
	}

	nonNullable := map[string]interface{}{
		"name":         o.Body.Name,
		"amount":       o.Body.Amount,
		"day_of_month": o.Body.DayOfMonth,
		"start_period": o.Body.StartPeriod,
	}

	for _, key := range o.PassedKeys {
		if val, ok := nonNullable[key]; ok && (val == nil || reflect.ValueOf(val).IsNil()) {
			return util.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s cannot be null", key))
		}
	}

	if slices.Contains(o.PassedKeys, "name") && o.Body.Name != nil && len(*o.Body.Name) == 0 {
		return util.NewHTTPError(http.StatusBadRequest, "name cannot be empty")
	}

	return nil
}

func (o *UpdateOptions) Run() error {
	payload := repository.UpdateRecurrenceDTO{
		Name:        util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "name"), Value: o.Body.Name},
		CategoryID:  util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "category_id"), Value: o.Body.CategoryID},
		Note:        util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "note"), Value: o.Body.Note},
		Amount:      util.OptionalNullable[float64]{Set: slices.Contains(o.PassedKeys, "amount"), Value: o.Body.Amount},
		DayOfMonth:  util.OptionalNullable[int]{Set: slices.Contains(o.PassedKeys, "day_of_month"), Value: o.Body.DayOfMonth},
		StartPeriod: util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "start_period"), Value: o.Body.StartPeriod},
		EndPeriod:   util.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "end_period"), Value: o.Body.EndPeriod},
	}

	rec, err := o.UseCases.Update(o.R.Context(), o.ID, o.UserID, payload)
	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusOK, util.ResponseData[UpdateRecurrenceResponseData]{
		Data: UpdateRecurrenceResponseData{
			Recurrence: MapRecurrenceResource(rec),
		},
	})
	return nil
}

// @Summary Update a recurrence
// @Description Update a recurrence template
// @ID v1UpdateRecurrence
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "recurrence ID"
// @Param body body UpdateRecurrenceRequest true "Recurrence payload"
// @Success 200 {object} util.ResponseData[UpdateRecurrenceResponseData] "Recurrence updated"
// @Failure 400 {object} util.HTTPError "Bad request"
// @Failure 401 {object} util.HTTPError "Unauthorized"
// @Failure 500 {object} util.HTTPError "Internal server error"
// @Router /api/v1/recurrences/{id} [patch]
func (api *API) Update(w http.ResponseWriter, r *http.Request) {
	cmd := &UpdateOptions{
		UseCases: api.recurrencesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
