package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anangfirmansyahp5/ecom/services/auth"
	"github.com/anangfirmansyahp5/ecom/types"
	"github.com/anangfirmansyahp5/ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", h.handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/products/{productID}", auth.WithJWTAuth(h.handleUpdateProduct, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/products/{productID}", auth.WithJWTAuth(h.handleDeleteProduct, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.Response{
		Success: true,
		Message: "get products success",
		Data:    products,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.Response{
		Success: true,
		Message: "create product success",
		Data:    nil,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productID"))
		return
	}

	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid productID"))
		return
	}

	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.Response{
		Success: true,
		Message: "get product success",
		Data:    product,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productID"))
		return
	}

	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid productID"))
		return
	}

	_, err = h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err = h.store.UpdateProduct(productID, types.Product{
		ID:          productID,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.Response{
		Success: true,
		Message: "product updated success",
		Data:    nil,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing productID"))
		return
	}

	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid productID"))
	}

	_, err = h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.DeleteProduct(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := types.Response{
		Data:    nil,
		Message: "delete product success",
		Success: true,
	}

	utils.WriteJSON(w, http.StatusNoContent, response)
}
