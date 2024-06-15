package handler

import (
	"net/http"

	"github.com/shahbaz275817/prismo/internal/responder"
	"github.com/shahbaz275817/prismo/internal/wrappers"
)

func PingHandler() http.HandlerFunc {
	return wrappers.DefaultWrapper(func(w http.ResponseWriter, r *http.Request) error {
		responder.WriteAnyResponse(r.Context(), w, map[string]string{
			"message": "pong",
		})
		return nil
	})
}
