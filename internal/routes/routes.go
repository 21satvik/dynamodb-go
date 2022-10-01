package routes

import "github.com/go-chi/chi"

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig().SetTimeout(ServerConfig.GetConfig().GetTimeout()),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters() *chi.Mux {}

func (r *Router) SetConfigsRouters() {}

func RouterHealth() {

}

func RouterProduct() {

}

func EnableTimeout() {

}

func EnableCORS() {

}

func EnableRecover() {

}

func EnableRequestID() {

}

func EnableRealIP() {

}
