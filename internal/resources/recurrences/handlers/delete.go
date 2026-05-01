package handlers

import (
	"net/http"
	"slices"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type DeleteOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.RecurrencesUseCases

	ID     string
	UserID string
	Scope  string
}

func (o *DeleteOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.ID = r.PathValue("id")
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	scope := r.URL.Query().Get("scope")
	if scope == "" {
		scope = "all"
	}
	o.Scope = scope

	return nil
}

func (o *DeleteOptions) Validate() error {
	allowedScopes := []string{"all", "until_current"}
	if !slices.Contains(allowedScopes, o.Scope) {
		return httputil.NewHTTPError(http.StatusBadRequest, "invalid scope. available: all, until_current")
	}
	return nil
}

func (o *DeleteOptions) Run() error {
	err := o.UseCases.DeleteByID(o.R.Context(), o.ID, o.UserID, o.Scope)
	if err != nil {
		return err
	}

	o.W.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary Delete Recurrence By ID
// @Description Delete a recurrence template
// @ID v1DeleteRecurrence
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "recurrence ID"
// @Param scope query string false "Handling of linked transactions: 'all' (default) deletes the recurrence and all related transactions (past/future); 'until_current' preserves past history but removes future recurrences." Enums(all, until_current) default(all)
// @Success 204 "Recurrence deleted"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 404 {object} httputil.HTTPError "Not found"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Router /api/v1/recurrences/{id} [delete]
func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	cmd := &DeleteOptions{
		UseCases: api.recurrencesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
