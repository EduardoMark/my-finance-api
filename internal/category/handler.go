package category

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/my-finance-api/internal/middlewares"
	"github.com/EduardoMark/my-finance-api/pkg/httputils"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	svc   Service
	token *token.TokenManager
}

func NewCategoryHandler(svc Service, token *token.TokenManager) CategoryHandler {
	return CategoryHandler{
		svc:   svc,
		token: token,
	}
}

func (h *CategoryHandler) RegisterCategoryRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(h.token))

		r.Post("/", h.Create)
		r.Get("/{id}", h.GetCategory)
		r.Get("/", h.GetAllCategoriesPerUserId)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middlewares.ContextUserID).(string)

	if !ok || userID == "" {
		httputils.Unauthorized(w)
		return
	}

	data, problems, err := httputils.DecodeValidJson[*CreateCategoryReq](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Create(ctx, userID, data); err != nil {
		httputils.Error(w, r, http.StatusInternalServerError, "handler create: "+err.Error())
		return
	}

	httputils.Created(w)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	userId, ok := ctx.Value(middlewares.ContextUserID).(string)
	if !ok || userId == "" {
		httputils.Unauthorized(w)
		return
	}

	record, err := h.svc.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			_ = httputils.EncodeJson(w, r, http.StatusNotFound, map[string]string{"error": ErrCategoryNotFound.Error()})
			return
		}
		_ = httputils.EncodeJson(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	res := CategoryRes{
		ID:        record.ID.String(),
		Name:      record.Name,
		Type:      string(record.Type),
		UserID:    record.UserID.String(),
		CreatedAt: record.CreatedAt.Time,
		UpdatedAt: record.UpdatedAt.Time,
	}

	_ = httputils.EncodeJson(w, r, http.StatusOK, res)
}

func (h *CategoryHandler) GetAllCategoriesPerUserId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(middlewares.ContextUserID).(string)
	if !ok || userId == "" {
		httputils.Unauthorized(w)
		return
	}

	records, err := h.svc.GetAllCategoriesByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrCategoriesNotFound) {
			_ = httputils.EncodeJson(w, r, http.StatusNotFound, map[string]string{"error": ErrCategoriesNotFound.Error()})
			return
		}
		_ = httputils.EncodeJson(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	res := make([]CategoryRes, len(records))
	for i, record := range records {
		res[i] = CategoryRes{
			ID:        record.ID.String(),
			Name:      record.Name,
			Type:      string(record.Type),
			UserID:    record.UserID.String(),
			CreatedAt: record.CreatedAt.Time,
			UpdatedAt: record.UpdatedAt.Time,
		}
	}

	_ = httputils.EncodeJson(w, r, http.StatusOK, res)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	data, problems, err := httputils.DecodeValidJson[*UpdateCategoryReq](r)
	if err != nil {
		_ = httputils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Update(ctx, id, *data); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			_ = httputils.EncodeJson(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		_ = httputils.EncodeJson(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	httputils.NoContent(w)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			_ = httputils.EncodeJson(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		_ = httputils.EncodeJson(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	httputils.NoContent(w)
}
