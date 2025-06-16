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
	})
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("user_id").(string)
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

	httpResponse.NoContent(w)
}
