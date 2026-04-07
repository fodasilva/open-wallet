package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type UpdateOptions struct {
	Ctx      *gin.Context
	UseCases usecases.RecurrencesUseCases

	ID         string
	UserID     string
	PassedKeys []string
	Body       UpdateRecurrenceRequest
}

func (o *UpdateOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.ID = ctx.Param("id")
	o.UserID = ctx.GetString("user_id")

	keys, err := utils.GetJSONKeys(ctx)
	if err != nil {
		return err
	}
	o.PassedKeys = keys

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return err
	}

	return nil
}

func (o *UpdateOptions) Validate() error {
	if len(o.PassedKeys) == 0 {
		return utils.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
	}

	nonNullable := map[string]interface{}{
		"name":         o.Body.Name,
		"amount":       o.Body.Amount,
		"day_of_month": o.Body.DayOfMonth,
		"start_period": o.Body.StartPeriod,
	}

	for _, key := range o.PassedKeys {
		if val, ok := nonNullable[key]; ok && (val == nil || reflect.ValueOf(val).IsNil()) {
			return utils.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s cannot be null", key))
		}
	}

	return nil
}

func (o *UpdateOptions) Run() error {
	payload := repository.UpdateRecurrenceDTO{
		Name:        utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "name"), Value: o.Body.Name},
		CategoryID:  utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "category_id"), Value: o.Body.CategoryID},
		Note:        utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "note"), Value: o.Body.Note},
		Amount:      utils.OptionalNullable[float64]{Set: slices.Contains(o.PassedKeys, "amount"), Value: o.Body.Amount},
		DayOfMonth:  utils.OptionalNullable[int]{Set: slices.Contains(o.PassedKeys, "day_of_month"), Value: o.Body.DayOfMonth},
		StartPeriod: utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "start_period"), Value: o.Body.StartPeriod},
		EndPeriod:   utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "end_period"), Value: o.Body.EndPeriod},
	}

	rec, err := o.UseCases.Update(o.ID, o.UserID, payload)
	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusOK, UpdateRecurrenceResponse{
		Data: UpdateRecurrenceResponseData{
			Recurrence: rec,
		},
	})
	return nil
}

// @Summary Update a recurrence
// @Description Update a recurrence template
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "recurrence ID"
// @Param body body UpdateRecurrenceRequest true "Recurrence payload"
// @Success 200 {object} UpdateRecurrenceResponse "Recurrence updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences/{id} [patch]
func (api *API) Update(ctx *gin.Context) {
	cmd := &UpdateOptions{
		UseCases: api.recurrencesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
