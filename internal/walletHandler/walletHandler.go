package walletHandler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"main/internal/request"
	"main/internal/walletService"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type WalletHandler struct {
	service *walletService.WalletService
}

func NewWalletHandler(serv *walletService.WalletService) *WalletHandler {
	return &WalletHandler{service: serv}
}

func (h *WalletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	const op = "internal.walletHandler.GetWalletBalance"

	receivedID := chi.URLParam(r, "WALLET_UUID")
	if receivedID == "" {
		log.Println(fmt.Errorf("%s: empty id", op))
		render.JSON(w, r, "id not found")
		return
	}

	balance, err := h.service.GetWalletBalance(receivedID)
	if err != nil {
		log.Println(fmt.Errorf("%s: %w", op, err))
		render.JSON(w, r, "failed to get balance")
		return
	}

	render.JSON(w, r, balance)
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	const op = "internal.walletHandler.CreateWallet"

	var requestData request.RequestCreateWallet
	err := render.DecodeJSON(r.Body, &requestData)
	if errors.Is(err, io.EOF) {
		log.Println(fmt.Errorf("%s: render body is empty", op))
		render.JSON(w, r, "empty request")
		return
	}
	if err != nil {
		log.Println(fmt.Errorf("%s: failed to decode request body: %w", op, err))
		render.JSON(w, r, "failde to decode request")
		return
	}

	createdId, err := h.service.CreateWallet(requestData.Balance)
	if err != nil {
		log.Println(fmt.Errorf("%s: failed to create wallet: %w", op, err))
		render.JSON(w, r, "failde to create wallet")
		return
	}
	render.JSON(w, r, createdId)
}

func (h *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	const op = "internal.walletHandler.UpdateWallet"

	var requestData request.RequestUpdateWallet
	err := render.DecodeJSON(r.Body, &requestData)
	if errors.Is(err, io.EOF) {
		log.Println(fmt.Errorf("%s: render body is empty", op))
		render.JSON(w, r, "empty request")
		return
	}
	if err != nil {
		log.Println(fmt.Errorf("%s: failed to decode request body: %w", op, err))
		render.JSON(w, r, "failde to decode request")
		return
	}

	updatedId, err := h.service.UpdateBalance(requestData.ID, requestData.Amount, requestData.TypeOperation)
	if err != nil {
		log.Println(fmt.Errorf("%s: failed to update wallet: %w", op, err))
		render.JSON(w, r, "failde to update wallet")
		return
	}
	render.JSON(w, r, updatedId)
}
