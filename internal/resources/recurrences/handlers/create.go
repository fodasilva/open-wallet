package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateOptions struct {
	Ctx      *gin.Context
	UseCases usecases.RecurrencesUseCases

	UserID     string
	PassedKeys []string
	Body       CreateRecurrenceRequest
}

func (o *CreateOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")

	keys, err := utils.GetJSONKeys(ctx)
	if err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	o.PassedKeys = keys

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *CreateOptions) Validate() error {
	return nil
}

func (o *CreateOptions) Run() error {
	payload := repository.CreateRecurrenceDTO{
		UserID:      o.UserID,
		Name:        o.Body.Name,
		CategoryID:  utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "category_id"), Value: o.Body.CategoryID},
		Note:        utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "note"), Value: o.Body.Note},
		Amount:      o.Body.Amount,
		DayOfMonth:  o.Body.DayOfMonth,
		StartPeriod: o.Body.StartPeriod,
		EndPeriod:   utils.OptionalNullable[string]{Set: slices.Contains(o.PassedKeys, "end_period"), Value: o.Body.EndPeriod},
	}

	res, err := o.UseCases.Create(payload)
	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusCreated, CreateRecurrenceResponse{
		Data: CreateRecurrenceResponseData{
			Recurrence: res,
		},
	})

	return nil
}

// @Summary Create a recurrence
// @Description Create a new recurrence template
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateRecurrenceRequest true "Recurrence payload"
// @Success 201 {object} CreateRecurrenceResponse "Recurrence created"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences [post]
func (api *API) Create(ctx *gin.Context) {
	cmd := &CreateOptions{
		UseCases: api.recurrencesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
