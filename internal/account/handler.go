package account

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/my-finance-api/internal/middlewares"
	"github.com/EduardoMark/my-finance-api/pkg/httputils"
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
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	data, problems, err := httputils.DecodeValidJson[*AccountCreateRequest](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Create(ctx, userID, *data); err != nil {
		httputils.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.Created(w)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	record, err := h.svc.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			httputils.Error(w, r, http.StatusBadRequest, "account not found")
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, err.Error())
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

	httputils.EncodeJson(w, r, http.StatusOK, response)
}

func (h *AccountHandler) GetAllAccountsPerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(middlewares.ContextUserID).(string)
	if !ok || userId == "" {
		httputils.Unauthorized(w)
		return
	}

	records, err := h.svc.GetAllAccountsByUserID(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrNoAccountsFound) {
			httputils.Error(w, r, http.StatusBadRequest, "account not found")
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, err.Error())
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

	httputils.EncodeJson(w, r, http.StatusOK, response)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	data, problems, err := httputils.DecodeValidJson[*AccountUpdateAccountReq](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	if err := h.svc.UpdateAccount(ctx, id, *data); err != nil {
		if err == ErrAccountNotFound {
			httputils.Error(w, r, http.StatusBadRequest, ErrAccountNotFound.Error())
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.NoContent(w)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	userId, ok := ctx.Value(middlewares.ContextUserID).(string)
	if !ok || userId == "" {
		httputils.Unauthorized(w)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			httputils.Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.NoContent(w)
}
