package order

import (
	"net/http"

	"github.com/anangfirmansyahp5/ecom/services/auth"
	"github.com/anangfirmansyahp5/ecom/types"
	"github.com/anangfirmansyahp5/ecom/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.OrderStore
	userStore types.UserStore
}

func NewHandler(store types.OrderStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders", auth.WithJWTAuth(h.handleCreateOrder, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload types.CreateOrderPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.store.CreateOrder(ctx, types.CreateOrderPayload{
		Items: payload.Items,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "create order success")
}
