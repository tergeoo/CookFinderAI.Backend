package mdw

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

/*
все мидлвары в 1
как временное решение со стандартным рутером http.ServeMux, у которого нет удобного способа добавить мидлвары
но так же ничего не мешает использовать его с другими рутерами
*/

func Common[T any](f func(w http.ResponseWriter, r *http.Request) (T, error)) http.Handler {
	md := RealIP(middleware.Recoverer(JSON(f)))
	md = CorsAllowAll(md)
	md = middleware.Logger(md)

	return md
}
