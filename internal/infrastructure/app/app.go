package app

import (
	"fmt"
	"os"
	"os/signal"
	"portfolyo/internal/infrastructure/cache"
	"portfolyo/internal/infrastructure/config"
	"portfolyo/internal/infrastructure/database"
	"portfolyo/internal/infrastructure/errorsx"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/uptrace/bun"
)

type App struct {
	FiberApp *fiber.App
	DB       *bun.DB
	Cfg      *config.Config
	Redis    *cache.RedisClient
}

type IRouter interface {
	RegisterRouter(app *App)
}

func New(router IRouter) *App {
	cfg, err := config.Setup()
	if err != nil {
		panic(err)
	}
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: errorsx.ErrorHandler,
	})

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // http://localhost:5173	http://---IP---:5173
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	db := database.New(cfg.Database)

	redisClient := cache.NewRedisClient(
		cfg.Redis.Host + ":" + cfg.Redis.Port,
	)

	app := &App{
		FiberApp: fiberApp,
		DB:       db,
		Cfg:      cfg,
		Redis:    redisClient,
	}

	router.RegisterRouter(app)

	return app
}

func (a *App) Start() {
	go func() {
		err := a.FiberApp.Listen(fmt.Sprintf(":%v", a.Cfg.Server.Port))
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}
