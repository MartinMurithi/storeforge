package router

import (
	"github.com/MartinMurithi/storeforge/gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, tenantHandler *handlers.TenantHandler, rbacHandler *handlers.RbacHandler, authMiddleware gin.HandlerFunc) *gin.Engine {
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

		user := api.Group("/users")
		user.Use(authMiddleware)
		{
			user.GET("/me", userHandler.GetCurrentUser)
			user.GET("/", userHandler.FetchAll) // Admin (paginated)
			// user.GET("/:id", userHandler.FetchByID)
			user.PATCH("/me", userHandler.UpdateMe)
		}

		stores := api.Group("/stores/new")
		stores.Use(authMiddleware)
		{
			stores.POST("/", tenantHandler.CreateTenant)
		}

		roles := api.Group("/roles")
		roles.Use(authMiddleware)
		{
			roles.POST("/", rbacHandler.CreateRole)
		}
	}

	return r
}
