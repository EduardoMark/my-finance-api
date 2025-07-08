package account

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMark/my-finance-api/internal/middlewares"
	"github.com/EduardoMark/my-finance-api/pkg/httpResponse"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	svc   Service
	token *token.TokenManager
}

func NewAccountHandler(svc Service, token *token.TokenManager) AccountHandler {
	return AccountHandler{
		svc:   svc,
		token: token,
	}
}

func (h *AccountHandler) RegisterAccountRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(h.token))

		r.Post("/", h.Create)
		r.Get("/{id}", h.GetAccount)
		r.Get("/", h.GetAllAccountsPerUser)
		r.Post("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok {
		httpResponse.Unauthorized(w)
		return
	}

	var body AccountCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpResponse.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.svc.Create(ctx, userID, body); err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpResponse.Created(w)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	_, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok {
		httpResponse.Unauthorized(w)
		return
	}

	record, err := h.svc.GetAccount(ctx, id)
	if err != nil {
		if err == ErrAccountNotFound {
			httpResponse.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := AccountResponse{
		ID:        record.ID.String(),
		UserID:    record.UserID.String(),
		Name:      record.Name,
		Type:      record.Type,
		Balance:   record.Balance.Float64,
		CreatedAt: record.CreatedAt.Time,
		UpdatedAt: record.UpdatedAt.Time,
	}

	httpResponse.SendJSON(w, http.StatusOK, response)
}

func (h *AccountHandler) GetAllAccountsPerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(middlewares.ContextUserID).(string)
	if !ok {
		httpResponse.Unauthorized(w)
		return
	}

	records, err := h.svc.GetAllAccountsByUserID(ctx, userId)
	if err != nil {
		if err == ErrNoAccountsFound {
			httpResponse.NotFound(w)
			return
		}
		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]AccountResponse, len(records))
	for i, record := range records {
		response[i] = AccountResponse{
			ID:        record.ID.String(),
			UserID:    record.UserID.String(),
			Name:      record.Name,
			Type:      record.Type,
			Balance:   record.Balance.Float64,
			CreatedAt: record.CreatedAt.Time,
			UpdatedAt: record.UpdatedAt.Time,
		}
	}

	httpResponse.SendJSON(w, http.StatusOK, response)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var body UpdateAccountReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpResponse.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.svc.UpdateAccount(ctx, id, body); err != nil {
		if err == ErrAccountNotFound {
			httpResponse.Error(w, http.StatusBadRequest, ErrAccountNotFound.Error())
			return
		}
		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpResponse.NoContent(w)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(ctx, id); err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpResponse.NoContent(w)
}
