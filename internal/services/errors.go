package services

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

var (
	FailedGoogleAuthenticationErr = httputil.NewHTTPError(http.StatusUnauthorized, "google authentication failed")
)
