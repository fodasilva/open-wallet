package services

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/util"
)

var (
	FailedGoogleAuthenticationErr = util.NewHTTPError(http.StatusUnauthorized, "google authentication failed")
)
