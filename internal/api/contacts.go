package api

import (
	"net/http"
	"strconv"

	"home-vue-go/internal/database"
	"home-vue-go/internal/ent/contact"

	"github.com/gin-gonic/gin"
)

func GetContacts(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		contacts, err := db.Client.Contact.Query().Order(contact.BySortOrder(), contact.ByID()).All(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(contacts))
		for i, contact := range contacts {
			result[i] = gin.H{
				"id":        contact.ID,
				"type":      contact.Type,
				"icon":      contact.Icon,
				"url":       contact.URL,
				"qrCode":    contact.QrCode,
				"hoverColor": contact.HoverColor,
				"sortOrder": contact.SortOrder,
			}
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateContact(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Type       string `json:"type" binding:"required"`
			Icon       string `json:"icon" binding:"required"`
			URL        string `json:"url"`
			QrCode     string `json:"qrCode"`
			HoverColor string `json:"hoverColor"`
			SortOrder  int    `json:"sortOrder"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证Email类型必须是mailto格式
		if req.Type == "Email" && req.URL != "" {
			if len(req.URL) < 7 || req.URL[:7] != "mailto:" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email URL必须是mailto:格式"})
				return
			}
		}

		ctx := c.Request.Context()
		create := db.Client.Contact.Create().
			SetType(req.Type).
			SetIcon(req.Icon).
			SetSortOrder(req.SortOrder)

		if req.URL != "" {
			create.SetURL(req.URL)
		}
		if req.QrCode != "" {
			create.SetQrCode(req.QrCode)
		}
		if req.HoverColor != "" {
			create.SetHoverColor(req.HoverColor)
		}

		contact, err := create.Save(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        contact.ID,
			"type":      contact.Type,
			"icon":      contact.Icon,
			"url":       contact.URL,
			"qrCode":    contact.QrCode,
			"hoverColor": contact.HoverColor,
			"sortOrder": contact.SortOrder,
		})
	}
}

func UpdateContact(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		var req struct {
			Type       string `json:"type"`
			Icon       string `json:"icon"`
			URL        string `json:"url"`
			QrCode     string `json:"qrCode"`
			HoverColor string `json:"hoverColor"`
			SortOrder  int    `json:"sortOrder"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证Email类型必须是mailto格式
		if req.Type == "Email" && req.URL != "" {
			if len(req.URL) < 7 || req.URL[:7] != "mailto:" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email URL必须是mailto:格式"})
				return
			}
		}

		ctx := c.Request.Context()
		update := db.Client.Contact.UpdateOneID(id)
		if req.Type != "" {
			update.SetType(req.Type)
		}
		if req.Icon != "" {
			update.SetIcon(req.Icon)
		}
		if req.URL != "" {
			update.SetURL(req.URL)
		}
		if req.QrCode != "" {
			update.SetQrCode(req.QrCode)
		}
		if req.HoverColor != "" {
			update.SetHoverColor(req.HoverColor)
		}
		update.SetSortOrder(req.SortOrder)

		contact, err := update.Save(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        contact.ID,
			"type":      contact.Type,
			"icon":      contact.Icon,
			"url":       contact.URL,
			"qrCode":    contact.QrCode,
			"hoverColor": contact.HoverColor,
			"sortOrder": contact.SortOrder,
		})
	}
}

func DeleteContact(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		ctx := c.Request.Context()
		if err := db.Client.Contact.DeleteOneID(id).Exec(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	}
}
