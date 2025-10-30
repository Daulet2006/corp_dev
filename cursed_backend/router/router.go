package router

import (
	"cursed_backend/handlers"
	"cursed_backend/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// === Логирование запросов ===
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		fmt.Printf("[%s] %s %s (%v)\n", c.Request.Method, c.Request.URL.Path, c.ClientIP(), time.Since(start))
	})

	// === CORS ===
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// === Prometheus ===
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// === Home ===
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API!"})
	})

	// === Public Routes ===
	public := r.Group("/api")
	{
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)
		public.GET("/pets", handlers.GetPets)
		public.GET("/products", handlers.GetProducts)
		public.GET("/stats", handlers.GetStats)
	}

	// === Protected Routes ===
	protected := r.Group("/api").Use(middleware.JWTAuth())
	{
		// Общие действия для всех авторизованных пользователей
		protected.PUT("/user", handlers.UpdateUser)
		protected.GET("/my/pets", handlers.MyPets)
		protected.GET("/my/products", handlers.MyProducts)

		// Покупка — доступна всем авторизованным (user/manager/admin)
		protected.POST("/pets/:id/buy", handlers.BuyPet)
		protected.POST("/products/:id/buy", handlers.BuyProduct)
	}

	// === Manager/Admin Routes (CRUD питомцев и товаров) ===
	manager := r.Group("/api").Use(middleware.JWTAuth(), middleware.RoleOr("manager", "admin"))
	{
		// Pets CRUD
		manager.POST("/pets", handlers.CreatePet)
		manager.GET("/pets/:id", handlers.GetPet)
		manager.PUT("/pets/:id", handlers.UpdatePet)
		manager.DELETE("/pets/:id", handlers.DeletePet)

		// Products CRUD
		manager.POST("/products", handlers.CreateProduct)
		manager.GET("/products/:id", handlers.GetProduct)
		manager.PUT("/products/:id", handlers.UpdateProduct)
		manager.DELETE("/products/:id", handlers.DeleteProduct)
	}

	// === Admin-only Routes ===
	admin := r.Group("/api/admin").Use(middleware.JWTAuth(), middleware.RoleAuth("admin"))
	{
		admin.GET("/users", handlers.GetUsers)
		admin.GET("/users/:id", handlers.GetUser)
		admin.POST("/users/:id/block", handlers.BlockUser)
		admin.POST("/users/:id/unblock", handlers.UnblockUser)
		admin.PUT("/users/:id/role", handlers.ChangeRole)
	}

	return r
}
