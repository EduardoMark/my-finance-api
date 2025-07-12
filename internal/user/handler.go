package user

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/my-finance-api/internal/middlewares"
	"github.com/EduardoMark/my-finance-api/pkg/httputils"
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

	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(h.token))

		r.Get("/{id}", h.GetUser)
		r.Get("/", h.GetAllUsers)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, problems, err := httputils.DecodeValidJson[*UserLoginRequest](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	token, err := h.svc.Login(ctx, h.token, *data)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httputils.Error(w, r, http.StatusBadRequest, map[string]string{"error": "invalid credential"})
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := UserLoginResponse{Token: token}

	_ = httputils.EncodeJson(w, r, http.StatusOK, resp)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, problems, err := httputils.DecodeValidJson[*UserCreateRequest](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Create(ctx, *data); err != nil {
		if errors.Is(err, ErrDuplicatedCredential) {
			httputils.Error(w, r, http.StatusBadRequest, map[string]string{"error": "user already exist"})
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	httputils.Created(w)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	record, err := h.svc.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httputils.Error(w, r, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	response := UserResponse{
		ID:        record.ID.String(),
		Name:      record.Name,
		Email:     record.Email,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}

	httputils.EncodeJson(w, r, http.StatusOK, response)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	records, err := h.svc.GetAllUsers(ctx)
	if err != nil {
		if errors.Is(err, ErrNoUsersFound) {
			httputils.Error(w, r, http.StatusNotFound, map[string]string{"error": "no users found"})
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := make([]UserResponse, len(records))
	for i, record := range records {
		response[i] = UserResponse{
			ID:        record.ID.String(),
			Name:      record.Name,
			Email:     record.Email,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}
	}

	httputils.EncodeJson(w, r, http.StatusOK, response)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	data, problems, err := httputils.DecodeValidJson[*UserUpdateRequest](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Update(ctx, id, *data); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httputils.Error(w, r, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}

		if errors.Is(err, ErrDuplicatedCredential) {
			httputils.Error(w, r, http.StatusBadRequest, map[string]string{"error": "user already exist"})
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	httputils.NoContent(w)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			httputils.Error(w, r, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}

		httputils.Error(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	httputils.NoContent(w)
}
