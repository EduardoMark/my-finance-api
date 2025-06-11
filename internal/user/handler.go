package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EduardoMark/my-finance-api/pkg/hash"
	"github.com/EduardoMark/my-finance-api/pkg/httpResponse"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	svc   Service
	token *token.TokenManager
}

func NewUserHandler(svc Service, token *token.TokenManager) *UserHandler {
	return &UserHandler{
		svc:   svc,
		token: token,
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/users/login", h.Login)
	r.Post("/users/signup", h.Signup)
	r.Get("/users/{id}", h.GetUser)
	r.Get("/users", h.GetAllUser)
	r.Put("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpResponse.Error(w, http.StatusBadRequest, "error when decoding body:"+err.Error())
		return
	}
	defer r.Body.Close()

	record, err := h.svc.GetUserByEmail(ctx, body.Email)
	if err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, "invalid credentials")
		return
	}

	if err := hash.ComparePassword(body.Password, record.Password); err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, "invalid credentials")
		return
	}

	token, err := h.token.GenerateToken(record.Name, record.Email)
	if err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, "error on generate token:"+err.Error())
		return
	}

	resp := UserLoginResponse{Token: token}

	httpResponse.SendJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpResponse.Error(w, http.StatusBadRequest, "error when decoding body:"+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.svc.Create(ctx, body); err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, "error on create user: "+err.Error())
		return
	}

	httpResponse.Created(w)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	record, err := h.svc.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpResponse.Error(w, http.StatusNotFound, "user not found")
			return
		}

		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := UserResponse{
		ID:        record.ID,
		Name:      record.Name,
		Email:     record.Email,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}

	httpResponse.SendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	records, err := h.svc.GetAllUsers(ctx)
	if err != nil {
		if errors.Is(err, ErrNoUsersFound) {
			httpResponse.Error(w, http.StatusNotFound, "no users found")
			return
		}
	}

	response := make([]UserResponse, len(records))
	for i, record := range records {
		response[i] = UserResponse{
			ID:        record.ID,
			Name:      record.Name,
			Email:     record.Email,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}
	}

	httpResponse.SendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var body UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpResponse.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	defer r.Body.Close()

	if err := h.svc.Update(ctx, id, body); err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, "error when updating user: "+err.Error())
		return
	}

	httpResponse.NoContent(w)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(ctx, id); err != nil {
		httpResponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpResponse.NoContent(w)
}
