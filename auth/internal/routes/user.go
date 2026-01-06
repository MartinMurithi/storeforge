package routes

import (
	"github.com/MartinMurithi/storeforge/auth/internal/handler"
	// "github.com/MartinMurithi/storeforge/auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler, authMiddleware gin.HandlerFunc) {
	// group endpoints based on api version
	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("/auth")

	{
		public.POST("/register", h.RegisterUser)
		public.POST("/login", h.LoginUser)
	}

	// Protected routes, revisit this later
	protected := v1.Group("/auth")
	protected.Use(authMiddleware)
	{
		protected.GET("/users", h.FetchAllUsers) // admin-only in handler logic
		protected.GET("/user", h.FetchUserById)   // self or admin
		//     protected.PUT("/:id", h.UpdateUser)      // self or admin
		//     protected.DELETE("/:id", h.DeleteUser)  // admin-only
	}
}
