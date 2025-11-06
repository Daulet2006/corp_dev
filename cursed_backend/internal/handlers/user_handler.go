package handlers

import (
	"cursed_backend/internal/db"
	"cursed_backend/internal/dto"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/models"
	"cursed_backend/internal/security"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

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
		Password  string `json:"password" validate:"required,strongpass"`
		Role      string `json:"role" validate:"omitempty,oneof=user manager admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid JSON: " + err.Error()})
		return
	}

	// sanitize inputs (names and email)
	req.FirstName = security.Sanitizer.Sanitize(strings.TrimSpace(req.FirstName))
	req.LastName = security.Sanitizer.Sanitize(strings.TrimSpace(req.LastName))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email)) // canonicalize

	validate := validator.New()
	security.RegisterCommonValidators(validate)
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Validation failed: " + err.Error()})
		return
	}

	// hash password with recommended cost (use 12-14; 14 is expensive)
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
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

	// create - handle unique constraint properly
	if err := db.GormDB.Create(&user).Error; err != nil {
		logger.Log.WithError(err).Warn("User creation failed - likely duplicate email")
		c.JSON(http.StatusConflict, models.APIResponse{Success: false, Message: "User already exists or invalid data"})
		return
	}

	// Do NOT return password in response. Use DTO.
	resp := dto.FromModel(&user)

	// Generate JWT using security package
	token, err := security.GenerateJWT(&user, 1*time.Hour)
	if err != nil {
		logger.Log.WithError(err).Error("Token generation failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to generate token"})
		return
	}
	// set secure cookie — note: requires HTTPS in production
	c.SetCookie("auth_token", token, 3600*24, "/", "", true, true) // Secure=true, HttpOnly=true

	logger.AuditLog("register", user.ID, c.ClientIP(), nil)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    gin.H{"user": resp},
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

	// Generate JWT using security package
	token, err := security.GenerateJWT(&user, 1*time.Hour)
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
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
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
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid JSON: " + err.Error(),
		})
		return
	}

	// Sanitize input (assuming you have a security.Sanitizer)
	req.FirstName = security.Sanitizer.Sanitize(strings.TrimSpace(req.FirstName))
	req.LastName = security.Sanitizer.Sanitize(strings.TrimSpace(req.LastName))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	validate := validator.New()
	security.RegisterCommonValidators(validate)
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// Build only allowed updates
	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "No fields to update",
		})
		return
	}

	// ✅ Correct way — Select takes a slice of strings directly, not variadic
	fields := make([]string, 0, len(updates))
	for k := range updates {
		fields = append(fields, k)
	}

	if err := db.GormDB.Model(&user).Select(fields).Updates(updates).Error; err != nil {
		logger.Log.WithError(err).Error("User update failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Update failed",
		})
		return
	}

	resp := dto.FromModel(&user)
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: resp})
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
	token, err := security.GenerateJWTWithClaims(userID, role, 15*time.Minute) // Short for refresh
	if err != nil {
		logger.Log.WithError(err).Error("Refresh token failed")
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Refresh failed"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: gin.H{"token": token}})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now().Format(time.RFC3339)})
}
