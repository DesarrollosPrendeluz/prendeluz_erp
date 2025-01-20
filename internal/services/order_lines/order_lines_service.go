package services

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"

	"github.com/gin-gonic/gin"
)

type OrderLineService interface {
	OrderLineLabel(id int) (dtos.OrderLineLable, error)
	ReturnDownloadPickingExcel(store_id int) string
	UpdateOrderLineHandler(

		c *gin.Context,
		requestBody dtos.OrdersLinesToUpdatePartially,
		token string,
		failedIds *[]int,
		errorList *[]error,
		callback func(*gin.Context, dtos.LineToUpdate, *models.OrderItem, error, *[]error),
		admin bool)
}
