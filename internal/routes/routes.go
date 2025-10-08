package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/aglili/auction-app/internal/handlers/middleware"
	"github.com/aglili/auction-app/internal/provider"
	"github.com/aglili/auction-app/pkg/constants"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

type SessionConfig struct {
	MaxAge     time.Duration
	HttpOnly   bool
	Secure     bool
	SameSite   http.SameSite
	SessionKey string
}

func SetupRoutes(prov *provider.Provider) http.Handler {
	if prov.Config.AppEnv == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	mux := gin.Default()

	store := setupStore(prov)

	mux.Use(sessions.Sessions("auction_store", store))

	v1 := mux.Group("/api/v1")
	v1.GET("/health", prov.HealthHandler.HealthHandler)

	// user routes
	users := v1.Group("/auth")
	users.POST("/sign-up", prov.UserHandler.CreateUserHandler)
	users.POST("/login", prov.UserHandler.LoginUser)

	protectedUser := v1.Group("/users")
	protectedUser.Use(middleware.RequireUserAuth())
	protectedUser.GET("/me", prov.UserHandler.GetUserProfile)
	protectedUser.POST("/logout", prov.UserHandler.Logout)

	auctions := v1.Group("/auctions")
	auctions.Use(middleware.RequireUserAuth())
	auctions.POST("", prov.AuctionHandler.CreateAuctionHandler)
	auctions.GET("/me", prov.AuctionHandler.GetUserAuctions)
	auctions.GET("/:id", prov.AuctionHandler.GetAuction)
	auctions.POST("/:id/bid", prov.BidHandler.CreateBid)
	auctions.GET("/ws", prov.WsHandler.HandleWSConnections)
	auctions.GET("/open", prov.AuctionHandler.GetOpenAuctions)

	return mux
}

func getSessionConfig(isProduction bool) SessionConfig {
	if isProduction {
		return SessionConfig{
			MaxAge:     24 * time.Hour,
			HttpOnly:   true,
			Secure:     true,
			SameSite:   http.SameSiteStrictMode,
			SessionKey: "auction_store",
		}
	}

	return SessionConfig{
		MaxAge:     2 * time.Hour,
		HttpOnly:   true,
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
		SessionKey: "auction_store",
	}
}

func setupStore(prov *provider.Provider) sessions.Store {
	store, err := postgres.NewStore(prov.DB, []byte(prov.Config.SecretKey))
	if err != nil {
		log.Println("failed to create session store")
	}

	config := getSessionConfig(prov.Config.AppEnv == constants.PRODUCTION)

	store.Options(sessions.Options{
		MaxAge:   int(config.MaxAge.Seconds()),
		HttpOnly: config.HttpOnly,
		Secure:   config.Secure,
		Path:     "/",
		SameSite: config.SameSite,
	})

	return store
}
