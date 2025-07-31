package api

import (
	"github.com/EduardoMark/my-finance-api/internal/account"
	"github.com/EduardoMark/my-finance-api/internal/category"
	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/internal/transaction"
	"github.com/EduardoMark/my-finance-api/internal/user"
	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	User        *user.UserHandler
	Account     *account.AccountHandler
	Category    *category.CategoryHandler
	Transaction *transaction.TransactionHandler
}

type Api struct {
	Router  *chi.Mux
	Cfg     *config.Env
	Db      *db.Queries
	Token   *token.TokenManager
	Handler *Handler
}

func NewApi(cfg *config.Env, db *db.Queries, token *token.TokenManager) *Api {
	return &Api{
		Router: chi.NewRouter(),
		Cfg:    cfg,
		Db:     db,
		Token:  token,
	}
}

func (api *Api) SetupApi() {
	userRepo := user.NewUserRepository(api.Db)
	userSvc := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc, api.Token)

	accRepo := account.NewAccountRepo(api.Db)
	accSvc := account.NewAccountService(accRepo)
	accHandler := account.NewAccountHandler(accSvc, api.Token)

	ctRepo := category.NewCategoryRepository(api.Db)
	ctSvc := category.NewCategoryService(ctRepo)
	ctHandler := category.NewCategoryHandler(ctSvc, api.Token)

	transRepo := transaction.NewTransactionRepo(api.Db)
	transSvc := transaction.NewTransactionService(transRepo)
	transHandler := transaction.NewTransactionHandler(transSvc, api.Token)

	api.Handler = &Handler{
		User:        userHandler,
		Account:     &accHandler,
		Category:    &ctHandler,
		Transaction: &transHandler,
	}
}
