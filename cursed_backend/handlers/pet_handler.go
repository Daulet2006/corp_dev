// Fixed UpdatePet in handlers/pet.go — bind to struct for proper JSON tag matching, validate ownerId if changed
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"cursed_backend/db"
	"cursed_backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MyPets(c *gin.Context) {
	var pets []models.Pet
	userID := c.GetUint("user_id")

	if err := db.GormDB.Where("owner_id = ?", userID).Find(&pets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch pets",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    pets,
	})
}

// GetPets (все или фильтр по owner) — public для store, auth для owned
func GetPets(c *gin.Context) {
	var pets []models.Pet
	query := db.GormDB

	ownerIDStr := c.Query("owner_id")

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0

	if ownerIDStr == "me" {
		if isAuth {
			query = query.Where("owner_id = ?", userID)
		}
	} else if ownerIDStr != "" {
		id, _ := strconv.ParseUint(ownerIDStr, 10, 32)
		targetID := uint(id)
		if targetID > 0 && !isAuth {
			c.JSON(http.StatusForbidden, gin.H{"error": "Auth required to view owned items"})
			c.Abort()
			return
		}
		if targetID > 0 && targetID != userID && role != "manager" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view these items"})
			c.Abort()
			return
		}
		query = query.Where("owner_id = ?", targetID)
	} else {
		// Default: все доступные
		if isAuth {
			if role == "manager" || role == "admin" {
				// Полный доступ
			} else {
				query = query.Where("owner_id = 0 OR owner_id = ?", userID)
			}
		} else {
			query = query.Where("owner_id = 0")
		}
	}

	query.Find(&pets)
	c.JSON(http.StatusOK, pets)
}

// GetPet (по ID) — public если owner_id=0, иначе owner/manager/admin
func GetPet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var pet models.Pet
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}
	// Проверка auth
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0
	if !isAuth && pet.OwnerID != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Auth required to view owned pet"})
		return
	}
	if pet.OwnerID != 0 && pet.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this pet"})
		return
	}
	c.JSON(http.StatusOK, pet)
}

// CreatePet (default OwnerID=0 из JSON или force; manager/admin валидируют >0)
func CreatePet(c *gin.Context) {
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := c.GetString("role")
	if role == "user" {
		pet.OwnerID = 0 // Force store для user, ignore JSON
	} else { // manager/admin
		if pet.OwnerID > 0 {
			var targetUser models.User
			if err := db.GormDB.First(&targetUser, pet.OwnerID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
				return
			}
		}
		// Если 0 или не указан — store
	}

	pet.CreatedAt = time.Now()
	if err := db.GormDB.Create(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pet)
}

// BuyPet (assign from store to self) — с транзакцией для atomicity и lock (FOR UPDATE)
func BuyPet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.GetUint("user_id")

	tx := db.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var pet models.Pet
	// Fixed: Lock without ORDER BY conflict
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pet, "id = ? AND owner_id = 0", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pet not found or already owned"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		tx.Rollback()
		return
	}

	if err := tx.Model(&pet).Update("owner_id", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Purchase failed"})
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	// Refresh для response (outside tx)
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh pet data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pet purchased", "pet": pet})
}

// UpdatePet — bind to struct for JSON tags, validate ownerId if changed
func UpdatePet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var pet models.Pet
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if pet.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update store item"})
		return
	}
	if pet.OwnerID != 0 && pet.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	var input models.Pet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = 0
	// Validate ownerId if changed
	if input.OwnerID != pet.OwnerID && input.OwnerID > 0 {
		var targetUser models.User
		if err := db.GormDB.First(&targetUser, input.OwnerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
			return
		}
	}
	if err := db.GormDB.Model(&pet).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Refresh with error check
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh pet data"})
		return
	}
	c.JSON(http.StatusOK, pet)
}

// DeletePet (owner/manager/admin; store only manager/admin)
func DeletePet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var pet models.Pet
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if pet.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete store item"})
		return
	}
	if pet.OwnerID != 0 && pet.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	if err := db.GormDB.Delete(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pet deleted"})
}
