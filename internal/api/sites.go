package api

import (
	"net/http"
	"strconv"

	"home-vue-go/internal/database"
	"home-vue-go/internal/ent/site"

	"github.com/gin-gonic/gin"
)

func GetSites(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		sites, err := db.Client.Site.Query().Order(site.BySortOrder(), site.ByID()).All(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(sites))
		for i, site := range sites {
			result[i] = gin.H{
				"id":        site.ID,
				"name":      site.Name,
				"url":       site.URL,
				"icon":      site.Icon,
				"sortOrder": site.SortOrder,
			}
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateSite(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name      string `json:"name" binding:"required"`
			URL       string `json:"url" binding:"required"`
			Icon      string `json:"icon" binding:"required"`
			SortOrder int    `json:"sortOrder"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		site, err := db.Client.Site.Create().
			SetName(req.Name).
			SetURL(req.URL).
			SetIcon(req.Icon).
			SetSortOrder(req.SortOrder).
			Save(ctx)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        site.ID,
			"name":      site.Name,
			"url":       site.URL,
			"icon":      site.Icon,
			"sortOrder": site.SortOrder,
		})
	}
}

func UpdateSite(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		var req struct {
			Name      string `json:"name"`
			URL       string `json:"url"`
			Icon      string `json:"icon"`
			SortOrder int    `json:"sortOrder"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		update := db.Client.Site.UpdateOneID(id)
		if req.Name != "" {
			update.SetName(req.Name)
		}
		if req.URL != "" {
			update.SetURL(req.URL)
		}
		if req.Icon != "" {
			update.SetIcon(req.Icon)
		}
		update.SetSortOrder(req.SortOrder)

		site, err := update.Save(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        site.ID,
			"name":      site.Name,
			"url":       site.URL,
			"icon":      site.Icon,
			"sortOrder": site.SortOrder,
		})
	}
}

func DeleteSite(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		ctx := c.Request.Context()
		if err := db.Client.Site.DeleteOneID(id).Exec(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	}
}
