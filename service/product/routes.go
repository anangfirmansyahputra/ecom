package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anangfirmansyahputra/ecom/service/auth"
	"github.com/anangfirmansyahputra/ecom/types"
	"github.com/anangfirmansyahputra/ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/products/{productID}", h.handleGetProductByID).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", auth.WithJWTAuth(h.handleDeleteProduct, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse form: %v", err))
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to get image: %v", err))
		return
	}

	defer file.Close()

	imagePath, err := utils.SaveUploadedFile(file, header)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	payload := types.ProductPayload{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       imagePath, // Path gambar yang disimpan
		Price:       utils.ParseFloat(r.FormValue("price")),
		Quantity:    utils.ParseInt(r.FormValue("quantity")),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if err := h.store.CreateProduct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	if err := h.store.DeleteProduct(productID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
