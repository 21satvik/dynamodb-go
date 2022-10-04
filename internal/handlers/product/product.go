package product

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/21satvik/dynamodb-go/internal/controllers/product"
	EntityProduct "github.com/21satvik/dynamodb-go/internal/entities/product"
	"github.com/21satvik/dynamodb-go/internal/handlers"
	"github.com/21satvik/dynamodb-go/internal/repository/adapter"
	Rules "github.com/21satvik/dynamodb-go/internal/rules"
	RulesProduct "github.com/21satvik/dynamodb-go/internal/rules/product"
	HttpStatus "github.com/21satvik/dynamodb-go/utils/http"
)

type ProductHandler struct {
	handlers.Interface
	Controller product.Interface
	Rules      Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &ProductHandler{
		Controller: product.NewController(repository),
		Rules:      RulesProduct.NewRules(),
	}
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *ProductHandler) getOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}
	response, err := h.Controller.ListOne(ID)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOK(w, r, response)
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOK(w, r, response)
}

func (h *ProductHandler) Post(w http.ResponseWriter, r *http.Request) {

	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(productBody)

	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"id": ID.String()})
}

func (h *ProductHandler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	productBody, err := h.getBodyAndValidate(r, ID)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, productBody); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusNoContent(w, r)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	if err := h.Controller.Remove(ID); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusNoContent(w, r)
}

func (h *ProductHandler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}

func (h *ProductHandler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityProduct.Product, error) {

	productBody := &EntityProduct.Product{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityProduct.Product{}, errors.New("body is required")
	}

	productParsed, err := EntityProduct.InterfaceToModel(body)
	if err != nil {
		return &EntityProduct.Product{}, errors.New("error in converting body to model")
	}

	setDefaultValues(productParsed, ID)
	return &productParsed, h.Rules.Validate(productParsed)
}

func setDefaultValues(product *EntityProduct.Product, ID uuid.UUID) {

	product.UpdatedAt = time.Now()

	if ID == uuid.Nil {
		product.ID == uuid.New()
		product.CreatedAt = time.Now()
	} else {
		product.ID = ID
	}
}
