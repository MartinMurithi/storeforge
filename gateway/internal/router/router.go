package router

import (
	"github.com/MartinMurithi/storeforge/gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, tenantHandler *handlers.TenantHandler, rbacHandler *handlers.RbacHandler, productHandler *handlers.ProductHandler, authMiddleware gin.HandlerFunc) *gin.Engine {
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

		stores := api.Group("/stores")
		stores.Use(authMiddleware)
		{
			stores.POST("/new", tenantHandler.CreateTenant)
			stores.GET("/:id", tenantHandler.GetTenantContext)
			stores.PATCH("/:id", tenantHandler.UpdateTenant)

			// All product routes are nested under /stores
			products := stores.Group("/:id/products")
			{
				products.POST("", productHandler.CreateProduct) // POST /stores/:tenantID/products
				// products.GET("", productHandler.GetProductsByTenant)                  // GET /stores/:tenantID/products
				// products.PATCH("/:productID", productHandler.UpdateProductWithImages) // PATCH /stores/:tenantID/products/:productID
				// products.DELETE("/:productID", productHandler.SoftDeleteProduct)      // DELETE /stores/:tenantID/products/:productID
			}
		}

		roles := api.Group("/roles")
		roles.Use(authMiddleware)
		{
			roles.POST("/", rbacHandler.CreateRole)
			roles.GET("/:id", rbacHandler.GetRoleById)
			roles.PATCH("/:id", rbacHandler.UpdateRole)
		}

	}

	return r
}
