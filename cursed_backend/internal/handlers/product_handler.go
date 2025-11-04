// handlers/products.go

package handlers

import (
	"cursed_backend/internal/db"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

func MyProducts(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	var products []models.Product
	userID := c.GetUint("user_id")

	if err := db.GormDB.Where("owner_id = ?", userID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    products,
	})
}

func GetProducts(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	var products []models.Product
	query := db.GormDB

	ownerIDStr := c.Query("owner_id")

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0

	if ownerIDStr == "me" {
		if isAuth {
			query = query.Where("owner_id = ?", userID)
		} else {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Auth required to view owned items",
			})
			c.Abort()
			return
		}
	} else if ownerIDStr != "" {
		id, _ := strconv.ParseUint(ownerIDStr, 10, 32)
		targetID := uint(id)
		if targetID > 0 && !isAuth {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Auth required to view owned items",
			})
			c.Abort()
			return
		}
		if targetID > 0 && targetID != userID && role != "manager" && role != "admin" {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Not authorized to view these items",
			})
			c.Abort()
			return
		}
		query = query.Where("owner_id = ?", targetID)
	} else {
		if isAuth {
			if role == "manager" || role == "admin" {
				// All products
			} else {
				query = query.Where("owner_id = 0 OR owner_id = ?", userID)
			}
		} else {
			query = query.Where("owner_id = 0")
		}
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    products,
	})
}

func GetProduct(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var product models.Product
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
		})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	isAuth := userID > 0
	if !isAuth && product.OwnerID != 0 {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Auth required to view owned product",
		})
		return
	}
	if product.OwnerID != 0 && product.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Not your product",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    product,
	})
}

func CreateProduct(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	err := models.ValidateProduct(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
		return
	}
	if product.OwnerID > 0 {
		var targetUser models.User
		if err := db.GormDB.First(&targetUser, product.OwnerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid owner_id",
			})
			return
		}
	}

	product.CreatedAt = time.Now()
	if err := db.GormDB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    product,
	})
}

func BuyProduct(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.GetUint("user_id")
	logger.Log.WithFields(logrus.Fields{"user_id": userID, "action": "buy_product"}).Info("BuyProduct called") // Fixed log
	tx := db.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var storeProduct models.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&storeProduct, "id = ? AND owner_id = 0 AND stock > 0", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Product not found, not available, or out of stock",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Database error: " + err.Error(),
			})
		}
		tx.Rollback()
		return
	}

	result := tx.Model(&storeProduct).Update("stock", gorm.Expr("stock - ?", 1))
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Stock update failed: " + result.Error.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Purchase failed: " + err.Error(),
		})
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Transaction commit failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product purchased",
		Data:    ownedProduct,
	})
}

func UpdateProduct(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var product models.Product
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
		})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if product.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Not authorized to update store item",
		})
		return
	}
	if product.OwnerID != 0 && product.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Not authorized",
		})
		return
	}
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	input.ID = 0
	if input.OwnerID != product.OwnerID && input.OwnerID > 0 {
		var targetUser models.User
		if err := db.GormDB.First(&targetUser, input.OwnerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid owner_id",
			})
			return
		}
	}
	if err := db.GormDB.Model(&product).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to refresh product data: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    product,
	})
}

func DeleteProduct(c *gin.Context) {
	if db.GormDB == nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database not available",
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var product models.Product
	if err := db.GormDB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Product not found",
		})
		return
	}
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if product.OwnerID == 0 && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Not authorized to delete store item",
		})
		return
	}
	if product.OwnerID != 0 && product.OwnerID != userID && role != "manager" && role != "admin" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Not authorized",
		})
		return
	}
	if err := db.GormDB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product deleted",
	})
}
