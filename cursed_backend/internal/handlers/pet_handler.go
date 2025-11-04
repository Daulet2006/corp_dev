package handlers

import (
	"cursed_backend/internal/db"
	"cursed_backend/internal/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MyPets(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	var pets []models.Pet
	userID := c.GetUint("user_id")

	if err := db.GormDB.Where("owner_id = ?", userID).Find(&pets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch pets"})
		return
	}

	// Sanitize descriptions
	for i := range pets {
		pets[i].Description = bluemonday.UGCPolicy().Sanitize(pets[i].Description)
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: pets})
}

func GetPets(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	var pets []models.Pet
	query := db.GormDB

	ownerIDStr := c.Query("owner_id")
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0

	if ownerIDStr == "me" {
		if isAuth {
			query = query.Where("owner_id = ?", userID)
		} else {
			c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Auth required to view owned items"})
			c.Abort()
			return
		}
	} else if ownerIDStr != "" {
		id, _ := strconv.ParseUint(ownerIDStr, 10, 32)
		targetID := uint(id)
		if targetID > 0 && !isAuth {
			c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Auth required to view owned items"})
			c.Abort()
			return
		}
		if targetID > 0 && targetID != userID && role != "manager" && role != "admin" {
			c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Not authorized to view these items"})
			c.Abort()
			return
		}
		query = query.Where("owner_id = ?", targetID)
	} else {
		if isAuth {
			if role == "manager" || role == "admin" {
				// All pets
			} else {
				query = query.Where("owner_id = 0 OR owner_id = ?", userID)
			}
		} else {
			query = query.Where("owner_id = 0")
		}
	}

	if err := query.Find(&pets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch pets"})
		return
	}

	// Sanitize
	for i := range pets {
		pets[i].Description = bluemonday.UGCPolicy().Sanitize(pets[i].Description)
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: pets})
}

func GetPet(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var pet models.Pet
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "Pet not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0
	if !isAuth && pet.OwnerID != 0 {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Auth required to view owned pet"})
		return
	}
	if pet.OwnerID != 0 && pet.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Not authorized to view this pet"})
		return
	}

	pet.Description = bluemonday.UGCPolicy().Sanitize(pet.Description)

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: pet})
}

func CreatePet(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	if err := models.ValidatePet(&pet); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Validation failed: " + err.Error()})
		return
	}

	pet.Description = bluemonday.UGCPolicy().Sanitize(pet.Description)

	if pet.OwnerID > 0 {
		var targetUser models.User
		if err := db.GormDB.First(&targetUser, pet.OwnerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid owner_id"})
			return
		}
	}

	pet.CreatedAt = time.Now()
	if err := db.GormDB.Create(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Creation failed"})
		return
	}
	c.JSON(http.StatusCreated, models.APIResponse{Success: true, Data: pet})
}

func BuyPet(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.GetUint("user_id")
	tx := db.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var pet models.Pet
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pet, "id = ? AND owner_id = 0", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Pet not found or already owned"})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database error"})
		}
		tx.Rollback()
		return
	}

	if err := tx.Model(&pet).Update("owner_id", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Purchase failed"})
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Transaction commit failed"})
		return
	}

	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to refresh pet data"})
		return
	}

	pet.Description = bluemonday.UGCPolicy().Sanitize(pet.Description)
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Pet purchased", Data: pet})
}

func UpdatePet(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var pet models.Pet
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "Pet not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if pet.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Not authorized to update store item"})
		return
	}
	if pet.OwnerID != 0 && pet.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Not authorized"})
		return
	}
	var input models.Pet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	input.ID = 0
	input.Description = bluemonday.UGCPolicy().Sanitize(input.Description)
	if input.OwnerID != pet.OwnerID && input.OwnerID > 0 {
		var targetUser models.User
		if err := db.GormDB.First(&targetUser, input.OwnerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid owner_id"})
			return
		}
	}
	if err := db.GormDB.Model(&pet).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Update failed"})
		return
	}
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to refresh pet data"})
		return
	}
	pet.Description = bluemonday.UGCPolicy().Sanitize(pet.Description)
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: pet})
}

func DeletePet(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var pet models.Pet
	if err := db.GormDB.First(&pet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{Success: false, Message: "Pet not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if pet.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Not authorized to delete store item"})
		return
	}
	if pet.OwnerID != 0 && pet.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{Success: false, Message: "Not authorized"})
		return
	}
	if err := db.GormDB.Delete(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Pet deleted"})
}
