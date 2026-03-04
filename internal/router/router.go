package router

import (
	"portfolyo/internal/handler"
	"portfolyo/internal/infrastructure/app"
	"portfolyo/internal/infrastructure/router"
	"portfolyo/internal/middleware"
	"portfolyo/internal/repository"
	"portfolyo/internal/service"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (Router) RegisterRouter(a *app.App) {
	app := a.FiberApp
	db := a.DB
	ks := a.Ks

	// Repos
	ur := repository.NewUserRepository(db)
	uar := repository.NewUserAssetsRepository(db)
	tr := repository.NewTransactionRepository(db)
	rr := repository.NewReminderRepository(db)

	// Services
	as := service.NewAuthService(ur)
	uas := service.NewUserAssetsService(uar, ks)
	ts := service.NewTransactionService(tr, ks, uar)
	rs := service.NewReminderService(rr)

	// Handlers
	ah := handler.NewAuthHandler(as)
	uah := handler.NewUserAssetsHandler(uas)
	th := handler.NewTransactionHandler(ts)
	rh := handler.NewReminderHandler(rs)

	v1 := app.Group("/api/v1")

	auth := v1.Group("/auth")
	router.Post(auth, "/register", ah.Register)
	router.Post(auth, "/login", ah.Login)

	v1.Use(middleware.JWTMiddleware())

	users := v1.Group("/users")
	router.Get(users, "/me", ah.GetUserProfile)
	router.Put(users, "/me", ah.UpdateUser)
	router.Delete(users, "/me", ah.DeleteUser)

	assets := v1.Group("/assets")
	router.Get(assets, "/all", uah.GetUserAssets)
	router.Get(assets, "/pdf", uah.GetUserAssetsPDF)
	router.Get(assets, "/:asset", uah.GetUserAsset)

	transactions := v1.Group("/transactions")
	router.Post(transactions, "/", th.AddTransaction)
	router.Get(transactions, "/all", th.GetAllTransaction)
	router.Get(transactions, "/excel", th.GetAllTransactionExcel)
	router.Get(transactions, "/pdf/:tx_id", th.GetTransactionPDF)
	router.Get(transactions, "/:asset", th.GetAllTransactionByAsset)

	reminders := v1.Group("/reminders")
	router.Post(reminders, "/", rh.Create)
	router.Get(reminders, "/", rh.GetAll)
	router.Delete(reminders, "/:id", rh.Delete)
}
