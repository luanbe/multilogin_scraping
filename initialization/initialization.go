package initialization

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"multilogin_scraping/app/delivery/api"
	"multilogin_scraping/helper"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"multilogin_scraping/app/delivery"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/helper/database"
)

type LoggerImpl struct {
	LogFilePath string
}

// TODO: add logger later
func InitDb() (*gorm.DB, error) {
	db, err := database.NewConnectionDB(
		viper.GetString("database.driver"),
		viper.GetString("database.dbname"),
		viper.GetString("database.host"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetInt("database.port"),
	)
	if err != nil {
		return nil, err
	}

	//run drop table to refresh data.
	// db.Migrator().DropTable(&entity.User{})

	// Define auto migration here
	_ = db.AutoMigrate(
		&entity.User{},
		&entity.ZillowMaindb3Address{},
		&entity.ZillowDetail{},
		&entity.ZillowPublicTaxHistory{},
		&entity.ZillowPriceHistory{},
	)

	//seedingPredefined(db, logger)

	return db, nil
}

func InitLogger(InitFields map[string]interface{}, LogFile string) *zap.Logger {

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	if InitFields != nil {
		cfg.InitialFields = InitFields
	}
	switch viper.GetString("crawler.log_level") {
	case "error":
		cfg.Level.SetLevel(zap.ErrorLevel)
	case "info":
		cfg.Level.SetLevel(zap.InfoLevel)
	case "debug":
		cfg.Level.SetLevel(zap.DebugLevel)
	default:
		cfg.Level.SetLevel(zap.InfoLevel)
	}

	// viper.GetString("crawler.log_level")
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout", "stderr"}
	if LogFile != "" {
		logFolder := "./logs"
		if _, err := os.Stat(logFolder); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(logFolder, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		cfg.OutputPaths = append(cfg.OutputPaths, fmt.Sprint(logFolder, "/", LogFile))
		//cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, fmt.Sprint(logFolder, "/", LogFile))
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	return logger
}

// IntSessionManager function create session manager and int configuration
// Docs: https://pkg.go.dev/github.com/alexedwards/scs#section-readme
func IntSessionManager() *scs.SessionManager {
	// Use Gob Register to apply User type to session manager
	gob.Register(entity.User{})

	var sessionManager *scs.SessionManager
	// Initialize a new session manager and configure the session lifetime.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	//sessionManager.Cookie.Name = "session_id"
	//sessionManager.Cookie.Domain = "example.com"
	//sessionManager.Cookie.HttpOnly = true
	//sessionManager.Cookie.Path = "/example/"
	//sessionManager.Cookie.Persist = true
	//sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true
	return sessionManager
}

// TODO: Add logger later
func InitRouting(
	db *gorm.DB,
	sessionManager *scs.SessionManager,
	rabbitMQ helper.RabbitMQBroker,
	redis helper.RedisCache,
) *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/", fontEndRouter(db))
	r.Mount("/api", apiRouter(db, rabbitMQ, redis))
	r.Mount("/admin", adminRouter(db, sessionManager))

	return r
}

func fontEndRouter(db *gorm.DB) http.Handler {
	// Service registry
	userService := registry.RegisterUserService(db)

	r := chi.NewRouter()
	index := delivery.NewIndexDelivery()
	r.Mount("/", index.Routes())

	user := delivery.NewUserDelivery(userService)
	r.Mount("/users", user.Routes())
	return r
}

func apiRouter(db *gorm.DB, rabbitMQ helper.RabbitMQBroker, redis helper.RedisCache) http.Handler {
	r := chi.NewRouter()
	crawler := api.CrawlerDelivery{RabbitMQ: rabbitMQ, Redis: redis}
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/crawler", crawler.Routes())
	return r
}

func adminRouter(db *gorm.DB, sessionManager *scs.SessionManager) http.Handler {
	// Services registry
	userService := registry.RegisterUserService(db)

	// Middlewares registry
	adminMiddleware := registry.RegisterAdminMiddleware(sessionManager)

	// Deliveries registry
	indexAdminDelivery := registry.RegisterIndexAdminDelivery(userService, sessionManager)
	userAdminDelivery := registry.RegisterUserAdminDelivery(userService, sessionManager)

	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "text/html; charset=utf-8"))
	r.Mount("/", indexAdminDelivery.Routes(adminMiddleware))
	r.With(adminMiddleware.UserAuth).Mount("/users", userAdminDelivery.Routes())

	return r
}
