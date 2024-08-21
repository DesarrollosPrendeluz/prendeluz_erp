package config

import (
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"

	"github.com/google/wire"
)

func InitializeController() error {

	wire.Build(
		itemsrepo.NewItemRepository,
		orderitemrepo.NewOrderItemRepository,
	)
	return nil
}
