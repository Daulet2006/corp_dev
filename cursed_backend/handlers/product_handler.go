// Fixed UpdateProduct in handlers/product.go — bind to struct for JSON tags, validate ownerId if changed
package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"cursed_backend/db"
	"cursed_backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

// GetProducts (аналогично GetPets)
func GetProducts(c *gin.Context) {
	var products []models.Product
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
		if isAuth {
			if role == "manager" || role == "admin" {
				// Полный
			} else {
				query = query.Where("owner_id = 0 OR owner_id = ?", userID)
			}
		} else {
			query = query.Where("owner_id = 0")
		}
	}

	query.Find(&products)
	c.JSON(http.StatusOK, products)
}

// GetProduct
func GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var product models.Product
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0
	if !isAuth && product.OwnerID != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Auth required to view owned product"})
		return
	}
	if product.OwnerID != 0 && product.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// CreateProduct
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := c.GetString("role")
	if role == "user" {
		product.OwnerID = 0 // Force store
	} else {
		if product.OwnerID > 0 {
			var targetUser models.User
			if err := db.GormDB.First(&targetUser, product.OwnerID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
				return
			}
		}
	}

	product.CreatedAt = time.Now()
	if err := db.GormDB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// BuyProduct — с транзакцией: lock + atomic update stock + create owned
func BuyProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.GetUint("user_id")

	tx := db.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var storeProduct models.Product
	// Fixed: Lock without ORDER BY conflict
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&storeProduct, "id = ? AND owner_id = 0 AND stock > 0", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found, not available, or out of stock"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		tx.Rollback()
		return
	}

	// Atomic decrement
	result := tx.Model(&storeProduct).Update("stock", gorm.Expr("stock - ?", 1))
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Stock update failed"})
		tx.Rollback()
		return
	}

	ownedProduct := models.Product{
		Name:        storeProduct.Name,
		Description: storeProduct.Description,
		Price:       storeProduct.Price,
		Stock:       1,
		Category:    storeProduct.Category,
		Brand:       storeProduct.Brand,
		Image:       storeProduct.Image,
		Mass:        storeProduct.Mass,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(&ownedProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Purchase failed"})
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product purchased", "product": ownedProduct})
}

// UpdateProduct — bind to struct for JSON tags, validate ownerId if changed
func UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var product models.Product
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if product.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update store item"})
		return
	}
	if product.OwnerID != 0 && product.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = 0
	// Validate ownerId if changed
	if input.OwnerID != product.OwnerID && input.OwnerID > 0 {
		var targetUser models.User
		if err := db.GormDB.First(&targetUser, input.OwnerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
			return
		}
	}
	if err := db.GormDB.Model(&product).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Refresh with error check
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh product data"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProduct
func DeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var product models.Product
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if product.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete store item"})
		return
	}
	if product.OwnerID != 0 && product.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	if err := db.GormDB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
func MyProducts(c *gin.Context) {
	var products []models.Product
	userID := c.GetUint("user_id")

	if err := db.GormDB.Where("owner_id = ?", userID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    products,
	})
}
