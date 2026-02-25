package router

import (
	"github.com/MartinMurithi/storeforge/gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, authMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.Default()

	// Global Middleware (Optional: CORS, Logging)
	// r.Use(middleware.CORS())

	// API Versioning
	api := r.Group("/api/v1")
	{
		// --- PUBLIC ROUTES ---
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// --- PROTECTED ROUTES ---
		// Apply the AuthMiddleware to this group only
		user := api.Group("/users")
		user.Use(authMiddleware)
		{
			user.GET("/me", userHandler.GetCurrentUser)  // Get logged-in user profile
			user.GET("/", userHandler.FetchAll) // Admin or listing (paginated)
			// user.GET("/:id", userHandler.FetchByID)
			user.PATCH("/me", userHandler.UpdateMe)
		}
	}

	return r
}
