package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type DeleteOptions struct {
	Ctx      *gin.Context
	UseCases usecases.RecurrencesUseCases

	ID     string
	UserID string
	Scope  string
}

func (o *DeleteOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.ID = ctx.Param("id")
	o.UserID = ctx.GetString("user_id")
	o.Scope = ctx.DefaultQuery("scope", "all")

	return nil
}

func (o *DeleteOptions) Validate() error {
	allowedScopes := []string{"all", "until_current"}
	if !slices.Contains(allowedScopes, o.Scope) {
		return utils.NewHTTPError(http.StatusBadRequest, "invalid scope. available: all, until_current")
	}
	return nil
}

func (o *DeleteOptions) Run() error {
	err := o.UseCases.DeleteByID(o.ID, o.UserID, o.Scope)
	if err != nil {
		return err
	}

	o.Ctx.Status(http.StatusNoContent)
	return nil
}

// @Summary Delete Recurrence By ID
// @Description Delete a recurrence template
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "recurrence ID"
// @Param scope query string false "Handling of linked transactions: 'all' (default) deletes the recurrence and all related transactions (past/future); 'until_current' preserves past history but removes future recurrences." Enums(all, until_current) default(all)
// @Success 204 "Recurrence deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences/{id} [delete]
func (api *API) Delete(ctx *gin.Context) {
	cmd := &DeleteOptions{
		UseCases: api.recurrencesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
