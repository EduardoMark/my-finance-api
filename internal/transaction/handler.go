package transaction

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/my-finance-api/internal/middlewares"
	"github.com/EduardoMark/my-finance-api/pkg/httputils"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	svc   Service
	token *token.TokenManager
}

func NewTransactionHandler(svc Service, token *token.TokenManager) TransactionHandler {
	return TransactionHandler{
		svc:   svc,
		token: token,
	}
}

func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/transactions", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(h.token))

		r.Post("/", h.Create)
		r.Get("/", h.GetAllTransactions)
		r.Get("/{id}", h.GetTransaction)
		r.Put("/{id}", h.UpdateTransaction)
		r.Delete("/{id}", h.DeleteTransaction)
	})
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	data, problems, err := httputils.DecodeValidJson[*TransactionCreateRequest](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	if err := h.svc.Create(ctx, userID, *data); err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			httputils.NotFound(w)
			return
		}
		httputils.Error(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		httputils.Error(w, r, http.StatusBadRequest, "transaction ID is required")
		return
	}

	transaction, err := h.svc.GetTransaction(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			httputils.NotFound(w)
			return
		}
		httputils.Error(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = httputils.EncodeJson(w, r, http.StatusOK, transaction)
}

func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	filters := &TransactionFilters{}

	if accountID := r.URL.Query().Get("account_id"); accountID != "" {
		filters.AccountID = &accountID
	}

	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		filters.CategoryID = &categoryID
	}

	if transactionType := r.URL.Query().Get("type"); transactionType != "" {
		filters.Type = &transactionType
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		filters.StartDate = &startDate
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		filters.EndDate = &endDate
	}

	transactions, err := h.svc.GetAllTransactions(ctx, userID, filters)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			_ = httputils.EncodeJson(w, r, http.StatusOK, []*TransactionResponse{})
			return
		}
		httputils.Error(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	_ = httputils.EncodeJson(w, r, http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		httputils.Error(w, r, http.StatusBadRequest, "transaction ID is required")
		return
	}

	data, problems, err := httputils.DecodeValidJson[*TransactionUpdateRequest](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	if err := h.svc.UpdateTransaction(ctx, id, *data); err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			httputils.NotFound(w)
			return
		}
		httputils.Error(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		httputils.Error(w, r, http.StatusBadRequest, "transaction ID is required")
		return
	}

	if err := h.svc.DeleteTransaction(ctx, id); err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			httputils.NotFound(w)
			return
		}
		httputils.Error(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
