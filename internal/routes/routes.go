package routes

import (
	ServerConfig "github.com/21satvik/dynamodb-go/config"
	HealthHandler "github.com/21satvik/dynamodb-go/internal/handlers/health"
	ProductHandler "github.com/21satvik/dynamodb-go/internal/handlers/product"
	"github.com/21satvik/dynamodb-go/internal/repository/adapter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig().SetTimeout(ServerConfig.GetConfig().Timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters(repository adapter.Interface) *chi.Mux {
	r.SetConfigsRouters()
	r.RouterHealth(repository)
	r.RouterProduct(repository)

	return r.router
}

func (r *Router) SetConfigsRouters() {
	r.EnableCORS()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecover()
	r.EnableRequestID()
	r.EnableRealIP()
}

func (r *Router) RouterHealth(repository adapter.Interface) {

	handler := HealthHandler.NewHandler(repository)

	r.router.Route("/health", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Put("/{ID}", handler.Put)
		route.Delete("/{ID}", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) RouterProduct(repository adapter.Interface) {

	handler := ProductHandler.NewHandler(repository)

	r.router.Route("/product", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Put("/{ID}", handler.Put)
		route.Delete("/{ID}", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use((middleware.RequestID))
	return r
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}
