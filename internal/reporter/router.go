package reporter

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes(router *gin.Engine) {
	fileGroup := router.Group("gstReport")
	fileGroup.POST("/monthly", h.GenerateMonthlyGSTReport)
}
