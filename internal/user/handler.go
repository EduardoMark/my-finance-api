package user

import (
	"encoding/json"
	"errors"
	"net/http"

	httpresponse "github.com/EduardoMark/my-finance-api/pkg/httpResponse"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	svc Service
}

func NewUserHandler(svc Service) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/users/signup", h.Signup)
	r.Get("/users/{id}", h.GetUser)
	r.Get("/users", h.GetAllUser)
	r.Put("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error when decoding body:"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.Create(ctx, body); err != nil {
		http.Error(w, "error on create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	httpresponse.Created(w)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	record, err := h.svc.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httpresponse.Error(w, http.StatusNotFound, "user not found")
			return
		}

		httpresponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := UserResponse{
		ID:        record.ID,
		Name:      record.Name,
		Email:     record.Email,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}

	httpresponse.SendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	records, err := h.svc.GetAllUsers(ctx)
	if err != nil {
		if errors.Is(err, ErrNoUsersFound) {
			httpresponse.Error(w, http.StatusNotFound, "no users found")
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

	httpresponse.SendJSON(w, http.StatusOK, response)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var body UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.svc.Update(ctx, id, body); err != nil {
		httpresponse.Error(w, http.StatusInternalServerError, "error when updating user: "+err.Error())
		return
	}

	httpresponse.NoContent(w)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(ctx, id); err != nil {
		httpresponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.NoContent(w)
}
