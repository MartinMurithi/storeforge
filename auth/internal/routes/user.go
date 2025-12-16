package routes

import (
	"github.com/MartinMurithi/storeforge/auth/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	UserHandler *handler.UserHandler
}

func NewUserRouter(r *gin.Engine, handler *handler.UserHandler) *UserRouter {
	return &UserRouter{UserHandler: handler}
}

func UserRoutes(r *gin.Engine, h *handler.UserHandler) {
	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("/users")
	{
		public.POST("/register", h.RegisterUser)
		public.POST("/login", h.LoginUser)
	}

	// Protected routes, revisit this later
	// protected := v1.Group("/users")
	// protected.Use(middleware.AuthMiddleware(authService))
	// {
	//     protected.GET("/", h.ListAllUsers)        // admin-only in handler logic
	//     protected.GET("/:id", h.GetUserById)     // self or admin
	//     protected.PUT("/:id", h.UpdateUser)      // self or admin
	//     protected.DELETE("/:id", h.DeleteUser)  // admin-only
	// }
}
