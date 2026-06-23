package http

import (
	"grps-go-redis-psql/internal/order"
	"grps-go-redis-psql/internal/user"
	"net/http"
	"strings"
)

type Router struct {
	userHandler  *user.Handler
	orderHandler *order.Handler
}

func NewRouter(userHandler *user.Handler, orderHandler *order.Handler) *Router {
	return &Router{
		userHandler:  userHandler,
		orderHandler: orderHandler}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {

	case req.URL.Path == "/health":
		health(w, req)

	case req.URL.Path == "/users":
		r.userHandler.Users(w, req)

	case strings.HasPrefix(req.URL.Path, "/users/"):
		r.userHandler.GetByID(w, req)

	case req.URL.Path == "/orders":
		r.orderHandler.Orders(w, req)

	default:
		http.NotFound(w, req)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
