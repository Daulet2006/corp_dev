package router

import (
	"cursed_backend/internal/handlers"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/metrics"
	"cursed_backend/internal/middleware"
	"expvar"
	"net/http"
	"strings"
	"time"

	"cursed_backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var rateLimiters = make(map[string]*rate.Limiter)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to set trusted proxies")
		return nil
	}

	// Security headers
	r.Use(middleware.SecurityHeaders())

	// Logging
	r.Use(requestLogger())

	// CORS
	origins := strings.Split(cfg.CORSOrigins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Metrics
	r.Use(metrics.Middleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Rate limit global (configurable, skips sensitive paths)
	r.Use(globalRateLimit(cfg))

	// Home
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API!"})
	})

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)
		public.GET("/pets", handlers.GetPets)
		public.GET("/products", handlers.GetProducts)
		public.GET("/stats", handlers.GetStats)
		public.GET("/health", handlers.HealthCheck)

		// MOVED: CSRF token endpoint to public (no auth needed for initial fetch)
		public.GET("/csrf-token", middleware.CSRFToken())
	}

	// Global error handler
	r.Use(middleware.ErrorHandler())

	// Protected routes
	protected := r.Group("/api").Use(middleware.JWTAuth(), middleware.CSRF())
	{
		protected.POST("/refresh", handlers.RefreshToken)
		protected.PUT("/user", handlers.UpdateUser)
		protected.GET("/my/pets", handlers.MyPets)
		protected.GET("/my/products", handlers.MyProducts)
		protected.POST("/pets/:id/buy", handlers.BuyPet)
		protected.POST("/products/:id/buy", handlers.BuyProduct)
		protected.DELETE("/pets/:id", handlers.DeletePet)
		protected.DELETE("/products/:id", handlers.DeleteProduct)
	}

	// Manager routes
	manager := r.Group("/api").Use(middleware.JWTAuth(), middleware.RoleOr("manager", "admin"), middleware.CSRF())
	{
		manager.POST("/pets", handlers.CreatePet)
		manager.GET("/pets/:id", handlers.GetPet)
		manager.PUT("/pets/:id", handlers.UpdatePet)
		manager.POST("/products", handlers.CreateProduct)
		manager.GET("/products/:id", handlers.GetProduct)
		manager.PUT("/products/:id", handlers.UpdateProduct)
	}

	// Admin routes
	admin := r.Group("/api/admin").Use(middleware.JWTAuth(), middleware.RoleAuth("admin"), middleware.CSRF())
	{
		admin.GET("/users", handlers.GetUsers)
		admin.GET("/users/:id", handlers.GetUser)
		admin.POST("/users/:id/block", handlers.BlockUser)
		admin.POST("/users/:id/unblock", handlers.UnblockUser)
		admin.PUT("/users/:id/role", handlers.ChangeRole)
		admin.GET("/debug/vars", gin.WrapH(expvar.Handler())) // Protected
	}

	return r
}

// Global rate limit middleware (now takes cfg, skips paths, configurable burst)
func globalRateLimit(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limit for health/debug/public auth/CSRF (prevents loops)
		skipPaths := []string{"/health", "/metrics", "/debug/vars", "/login", "/register", "/csrf-token"}
		path := c.Request.URL.Path
		shouldSkip := false
		for _, skip := range skipPaths {
			if strings.HasSuffix(path, skip) {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			c.Next()
			return
		}

		ip := c.ClientIP()
		if _, ok := rateLimiters[ip]; !ok {
			// Dev: Higher limit; Prod: Strict
			burst := 80
			if cfg.Env == "dev" {
				burst = 100
			}
			rateLimiters[ip] = rate.NewLimiter(rate.Every(time.Minute), burst)
		}
		if !rateLimiters[ip].Allow() {
			logger.Log.WithFields(logrus.Fields{
				"ip":   ip,
				"path": path,
			}).Warn("Rate limit exceeded")
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Request logger with logrus
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userID := c.GetUint("user_id")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		fields := logrus.Fields{
			"method":  c.Request.Method,
			"path":    path,
			"ip":      ip,
			"user_id": userID,
			"status":  status,
			"latency": latency.String(),
			"size":    c.Writer.Size(),
		}
		if status >= 500 {
			logger.Log.WithFields(fields).Error("HTTP request failed")
		} else if status >= 400 {
			logger.Log.WithFields(fields).Warn("HTTP request warning")
		} else {
			logger.Log.WithFields(fields).Info("HTTP request")
		}
	}
}
