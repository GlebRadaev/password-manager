// Package swagger provides Swagger UI integration for API documentation.
package swagger

import (
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterSwaggerUI adds Swagger endpoints to an existing http.Handler.
func RegisterSwaggerUI(mux http.Handler, appName string) http.Handler {
	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)

	swaggerFile := strings.ToLower(appName) + ".swagger.json"
	swaggerPath := "/swagger/" + swaggerFile

	httpMux.HandleFunc(swaggerPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	httpMux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL(swaggerFile)))

	return httpMux
}
