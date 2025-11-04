package handlers

import (
	"cursed_backend/internal/db"
	"cursed_backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Database not available"})
		return
	}

	var userCount, petCount, productCount, ownedPets, ownedProducts, storePets, storeProducts int64

	tx := db.GormDB.Begin()
	defer tx.Rollback()

	if err := tx.Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch user count"})
		return
	}
	if err := tx.Model(&models.Pet{}).Count(&petCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch pet count"})
		return
	}
	if err := tx.Model(&models.Product{}).Count(&productCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch product count"})
		return
	}
	if err := tx.Model(&models.Pet{}).Where("owner_id > 0").Count(&ownedPets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch owned pets"})
		return
	}
	if err := tx.Model(&models.Product{}).Where("owner_id > 0").Count(&ownedProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch owned products"})
		return
	}
	if err := tx.Model(&models.Pet{}).Where("owner_id = 0").Count(&storePets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch store pets"})
		return
	}
	if err := tx.Model(&models.Product{}).Where("owner_id = 0").Count(&storeProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Failed to fetch store products"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Transaction failed"})
		return
	}

	stats := models.APIResponse{
		Success: true,
		Data: gin.H{
			"users":         userCount,
			"totalPets":     petCount,
			"ownedPets":     ownedPets,
			"storePets":     storePets,
			"totalProducts": productCount,
			"ownedProducts": ownedProducts,
			"storeProducts": storeProducts,
		},
	}
	c.JSON(http.StatusOK, stats)
}
