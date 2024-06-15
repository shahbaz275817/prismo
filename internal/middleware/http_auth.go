package middleware

import (
	"net/http"

	"github.com/shahbaz275817/prismo/internal/config"
	"github.com/shahbaz275817/prismo/internal/responder"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

func WithHTTPAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		cfgAuthUsername := config.Auth().Username
		cfgAuthPassword := config.Auth().Password
		username, password, ok := req.BasicAuth()

		if !ok {
			responder.WriteError(wr, req, errors.NewUnauthorizedError("Not Authorized", nil))
			return
		}

		isValid := (username == cfgAuthUsername) && (password == cfgAuthPassword)
		if !isValid {
			responder.WriteError(wr, req, errors.NewUnauthorizedError("Wrong username/password", nil))
			return
		}
		next.ServeHTTP(wr, req)
	})
}
