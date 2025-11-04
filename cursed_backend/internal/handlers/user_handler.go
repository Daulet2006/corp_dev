package handlers

import (
	"cursed_backend/internal/db"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET is required")
	}
}

func GetUsers(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Admin only"})
		return
	}

	var users []models.User
	if err := db.GormDB.Find(&users).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to fetch users")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: users})
}

func GetUser(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid user ID"})
		return
	}
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Admin only"})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: user})
}

func Register(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	var req struct {
		FirstName string `json:"firstName" validate:"required,min=2,max=50"`
		LastName  string `json:"lastName" validate:"required,min=2,max=50"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8,strongpass"`
		Role      string `json:"role" validate:"omitempty,oneof=user manager admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Validation failed: " + err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		logger.Log.WithError(err).Error("Password hashing failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to hash password"})
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashed),
		Role:      models.RoleUser,
	}
	if req.Role != "" {
		user.Role = models.Role(req.Role)
	}

	if err := db.GormDB.Create(&user).Error; err != nil {
		logger.Log.WithError(err).Warn("User creation failed - duplicate email")
		c.JSON(http.StatusConflict, models.APIResponse{Success: false, Message: "User already exists"})
		return
	}

	token, err := generateJWT(&user)
	if err != nil {
		logger.Log.WithError(err).Error("Token generation failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to generate token"})
		return
	}

	logger.AuditLog("register", user.ID, c.ClientIP(), nil)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    gin.H{"token": token, "user": user},
	})
}

func Login(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	var user models.User
	if err := db.GormDB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		logger.AuditLog("login_fail", 0, c.ClientIP(), err)
		c.JSON(http.StatusUnauthorized, models.APIResponse{Success: false, Message: "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.AuditLog("login_fail", user.ID, c.ClientIP(), err)
		c.JSON(http.StatusUnauthorized, models.APIResponse{Success: false, Message: "Invalid credentials"})
		return
	}

	if user.Blocked {
		logger.AuditLog("login_blocked", user.ID, c.ClientIP(), nil)
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "User blocked"})
		return
	}

	token, err := generateJWT(&user)
	if err != nil {
		logger.Log.WithError(err).Error("Token generation failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to generate token"})
		return
	}

	logger.AuditLog("login", user.ID, c.ClientIP(), nil)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"token": token, "user": user},
	})
}

func UpdateUser(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	userID := c.GetUint("user_id")
	var req struct {
		FirstName string `json:"firstName" validate:"omitempty,min=2,max=50"`
		LastName  string `json:"lastName" validate:"omitempty,min=2,max=50"`
		Email     string `json:"email" validate:"omitempty,email"`
		Image     string `json:"image" validate:"omitempty,url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "User not found"})
		return
	}

	updates := map[string]interface{}{
		"FirstName": req.FirstName,
		"LastName":  req.LastName,
		"Email":     req.Email,
		"Image":     req.Image,
	}
	if err := db.GormDB.Model(&user).Updates(updates).Error; err != nil {
		logger.Log.WithError(err).Error("User update failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Update failed"})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: user})
}

func BlockUser(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid user ID"})
		return
	}
	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "User not found"})
		return
	}

	user.Blocked = true
	if err := db.GormDB.Save(&user).Error; err != nil {
		logger.Log.WithError(err).Error("Block user failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Block failed"})
		return
	}

	logger.AuditLog("block_user", user.ID, c.ClientIP(), nil)
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User blocked", Data: user})
}

func UnblockUser(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid user ID"})
		return
	}
	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "User not found"})
		return
	}

	user.Blocked = false
	if err := db.GormDB.Save(&user).Error; err != nil {
		logger.Log.WithError(err).Error("Unblock user failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Unblock failed"})
		return
	}

	logger.AuditLog("unblock_user", user.ID, c.ClientIP(), nil)
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User unblocked", Data: user})
}

func ChangeRole(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid user ID"})
		return
	}
	var req struct {
		Role string `json:"role" validate:"required,oneof=user manager admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	newRole := models.Role(req.Role)
	if !newRole.IsValid() {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid role"})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "User not found"})
		return
	}

	user.Role = newRole
	if err := db.GormDB.Save(&user).Error; err != nil {
		logger.Log.WithError(err).Error("Role change failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Role change failed"})
		return
	}

	logger.AuditLog("change_role", user.ID, c.ClientIP(), nil)
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Role changed", Data: user})
}

func RefreshToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	token, err := generateJWTWithClaims(userID, role, 15*time.Minute) // Short for refresh
	if err != nil {
		logger.Log.WithError(err).Error("Refresh token failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Refresh failed"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: gin.H{"token": token}})
}

func generateJWT(user *models.User) (string, error) {
	return generateJWTWithClaims(user.ID, user.Role.String(), 24*time.Hour)
}

func generateJWTWithClaims(userID uint, role string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(exp).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now().Format(time.RFC3339)})
}
