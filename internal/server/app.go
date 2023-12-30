package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/YnMann/chat_backend/internal/auth"

	// a = auth
	ahttp "github.com/YnMann/chat_backend/internal/auth/delivery/http"
	amongo "github.com/YnMann/chat_backend/internal/auth/repository/mongo"
	ausecase "github.com/YnMann/chat_backend/internal/auth/usecase"
)

// APIServer
type App struct {
	logger *logrus.Logger

	httpServer *http.Server
	authUC     auth.UseCase
	// store  *store.Store
}

// New
func NewApp() *App {
	db := initDb()

	userRepo := amongo.NewUserRepository(db, viper.GetString("mongo.user_collection"))

	return &App{
		logger: logrus.New(),
		authUC: ausecase.NewAuthUseCase(
			userRepo,
			viper.GetString("auth.hash_salt"),
			[]byte(viper.GetString("auth.signing_key")),
			viper.GetDuration("auth.token_ttl"),
		),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	r := gin.Default()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Maintenance of static files for the front of a folder 'web/build'
	// this will provide access to the files along the way /public
	r.Static("/public", "./web/build")

	r.NoRoute(func(c *gin.Context) {
		c.File("./web/build/index.html")
	})

	// Set up http handlers
	// SignUp/SignIn endpoints
	ahttp.RegisterHTTPEndpoints(r, a.authUC)

	// Set up http handlers
	// SignUp/SignIn endpoints
	ahttp.RegisterHTTPEndpoints(r, a.authUC)

	// API endpoints
	// authMiddleware := ahttp.NewAuthMiddleware(a.authUC)
	// api := router.Group("/api", authMiddleware)

	// HTTP server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and server: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDb() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
