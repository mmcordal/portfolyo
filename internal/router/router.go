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
	redis := a.Redis

	// Repos
	ur := repository.NewUserRepository(db)
	uar := repository.NewUserAssetsRepository(db)
	tr := repository.NewTransactionRepository(db)
	rr := repository.NewReminderRepository(db)

	// Services
	as := service.NewAuthService(ur)
	ks := service.NewKurService(redis)
	uas := service.NewUserAssetsService(uar, ks, tr)
	rs := service.NewReminderService(rr)

	// Handlers
	ah := handler.NewAuthHandler(as)
	uah := handler.NewUserAssetsHandler(uas)
	rh := handler.NewReminderHandler(rs)

	v1 := app.Group("/api/v1")

	router.Post(v1, "/Register", ah.Register)
	router.Post(v1, "/Login", ah.Login)

	v1.Use(middleware.JWTMiddleware())

	// User
	router.Get(v1, "/user/profile", ah.GetUserProfile)
	router.Put(v1, "/user/update", ah.UpdateUser)
	router.Delete(v1, "/user/delete", ah.DeleteUser)

	// UserAssets
	router.Post(v1, "/user/your-assets", uah.UserAssetAdd)
	router.Get(v1, "/user/your-assets", uah.GetUserAssets)
	router.Get(v1, "/user/your-assets/pdf", uah.GetUserAssetsPDF)
	router.Get(v1, "/user/transaction", uah.GetAllTransaction)
	router.Get(v1, "/user/transaction/pdf", uah.GetTransactionPDF)
	router.Get(v1, "/user/transaction/excel", uah.GetAllTransactionExcel)

	// Reminder
	router.Post(v1, "/user/reminder", rh.Create)
	router.Get(v1, "/user/reminder", rh.GetAll)
	router.Delete(v1, "/user/reminder", rh.Delete)
}
