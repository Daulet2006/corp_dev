// handlers/stats.go — улучшено: добавлены store counts явно (WHERE owner_id=0)
package handlers

import (
	"cursed_backend/db"
	"cursed_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStats (добавь counts по owners)
func GetStats(c *gin.Context) {
	var userCount, petCount, productCount, ownedPets, ownedProducts, storePets, storeProducts int64

	db.GormDB.Model(&models.User{}).Count(&userCount)
	db.GormDB.Model(&models.Pet{}).Count(&petCount)
	db.GormDB.Model(&models.Product{}).Count(&productCount)
	db.GormDB.Model(&models.Pet{}).Where("owner_id > 0").Count(&ownedPets)
	db.GormDB.Model(&models.Product{}).Where("owner_id > 0").Count(&ownedProducts)
	db.GormDB.Model(&models.Pet{}).Where("owner_id = 0").Count(&storePets)
	db.GormDB.Model(&models.Product{}).Where("owner_id = 0").Count(&storeProducts)

	stats := gin.H{
		"users":         userCount,
		"totalPets":     petCount,
		"ownedPets":     ownedPets,
		"storePets":     storePets,
		"totalProducts": productCount,
		"ownedProducts": ownedProducts,
		"storeProducts": storeProducts,
	}
	c.JSON(http.StatusOK, stats)
}
