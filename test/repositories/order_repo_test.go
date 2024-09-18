package repositories

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/repositories/orderrepo"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObtainSupplierOrder(t *testing.T) {
	var expectedType []dtos.SupplierOrders
	repo := orderrepo.NewOrderRepository(db.DB)
	data, err := repo.GetSupplierOrders(nil)
	assert.NoError(t, err)
	assert.IsType(t, expectedType, data)

}
