// handlers/user.go — added GetUsers and GetUser for admin
package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"cursed_backend/db"
	"cursed_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("your-secret-key")
	}
}

// GetUsers (admin only)
func GetUsers(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	var users []models.User
	if err := db.GormDB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser (admin only)
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Register (без изменений, role=user по умолчанию)
func Register(c *gin.Context) {
	var req struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"` // Игнорируем, всегда user
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashed),
		Role:      models.RoleUser,
	}

	if err := db.GormDB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	token, _ := generateJWT(&user)
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

// Login (без изменений)
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.GormDB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Blocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "User blocked"})
		return
	}

	token, _ := generateJWT(&user)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// UpdateUser (свои данные, any role)
func UpdateUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Image     string `json:"image"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	db.GormDB.Model(&user).Updates(models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Image:     req.Image,
	})

	c.JSON(http.StatusOK, user)
}

// BlockUser (admin only, используем Param("id")) — можно добавить tx если нужно, но single update OK
func BlockUser(c *gin.Context) {
	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Blocked = true
	db.GormDB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User blocked", "user": user})
}

// UnblockUser (admin only)
func UnblockUser(c *gin.Context) {
	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Blocked = false
	db.GormDB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User unblocked", "user": user})
}

// ChangeRole (admin only: PUT /admin/users/:id/role { "role": "manager" })
func ChangeRole(c *gin.Context) {
	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRole := models.Role(req.Role)
	if !newRole.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var user models.User
	if err := db.GormDB.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = newRole
	db.GormDB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Role changed", "user": user})
}

func generateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}
