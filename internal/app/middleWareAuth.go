package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OJOMB/lang-learn-svc/internal/pkg/contextual"
)

func (app *App) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderVal := r.Header.Get("Authorization")
		if authHeaderVal == "" {
			app.logger.Info("no token found in request")
			appErr := newAppErr("no token found in request", http.StatusUnauthorized)
			http.Error(w, appErr.Error(), appErr.code)
			return
		}

		tokenString := strings.TrimPrefix(authHeaderVal, "Bearer ")
		if authHeaderVal == tokenString {
			appErr := newAppErr("auth header value in unexpected format", http.StatusUnauthorized)
			http.Error(w, appErr.Error(), appErr.code)
			return
		}

		subj, err := app.tokenAuth.GetSubject(tokenString)
		if err != nil {
			appErr := newAppErr(fmt.Sprintf("invalid token: %s", err.Error()), http.StatusUnauthorized)
			http.Error(w, appErr.Error(), appErr.code)
			return
		}

		// include the subject in the context as this is needed at the domain level
		r = r.WithContext(contextual.SetSubjectID(r.Context(), subj))
		next.ServeHTTP(w, r)
	})
}
