package routing

import (
	"database/sql"
	"net/http"

	"github.com/Daniorocket/RestApi-SQLite/handlers"
	"github.com/Daniorocket/RestApi-SQLite/logger"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {

	handler := handlers.Handler{Db: db}
	initRoutes(handler)
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range initRoutes(handler) {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
