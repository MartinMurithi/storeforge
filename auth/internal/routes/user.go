package routes

import (
	"github.com/MartinMurithi/storeforge/auth/internal/handler"
	// "github.com/MartinMurithi/storeforge/auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	// group endpoints based on api version
	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("/users")

	{
		public.POST("/register", h.RegisterUser)
		public.POST("/login", h.LoginUser)
	}

	// Protected routes, revisit this later
	// protected := v1.Group("/users")
	// protected.Use(middleware.AuthMiddleware(publicKey, "storeforge-api", "auth.storeforge.io"))
	// {
	//     protected.GET("/", h.ListAllUsers)        // admin-only in handler logic
	//     protected.GET("/:id", h.GetUserById)     // self or admin
	//     protected.PUT("/:id", h.UpdateUser)      // self or admin
	//     protected.DELETE("/:id", h.DeleteUser)  // admin-only
	// }
}
