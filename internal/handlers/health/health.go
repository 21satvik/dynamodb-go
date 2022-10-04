package health

import (
	"errors"
	"net/http"

	"github.com/21satvik/dynamodb-go/internal/handlers"
	"github.com/21satvik/dynamodb-go/internal/repository/adapter"
	HttpStatus "github.com/21satvik/dynamodb-go/utils/http"
)

type HealthHandler struct {
	handlers.Interface
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &HealthHandler{
		Repository: repository,
	}
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {

	if !h.Repository.Health() {
		HttpStatus.StatusInternalServerError(w, r, errors.New("relational database not alive"))
		return
	}
	HttpStatus.StatusOK(w, r, "Service OK")
}

func (h *HealthHandler) Post(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}

func (h *HealthHandler) Put(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}

func (h *HealthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}

func (h *HealthHandler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}
