package controllers

import (
	"net/http"
	"strconv"

	"errors"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/itemlocationrepo"
	"prendeluz/erp/internal/repositories/storelocationrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
	ItemStockLocationService "prendeluz/erp/internal/services/item_stock_locations"

	"github.com/gin-gonic/gin"
)

func GetItemStockLocation(c *gin.Context) {
	var err error
	var data []models.ItemLocation
	var datum *models.ItemLocation
	var recount int64

	repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	main_sku := c.Query("main_sku")
	store_id, _ := strconv.Atoi(c.DefaultQuery("store_id", "0"))
	storeLocation, _ := strconv.Atoi(c.DefaultQuery("item_stock_location", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if main_sku != "" && store_id != 0 {
		data, err = repo.FindByItemsAndStore(main_sku, uint64(store_id), pageSize, page)

	} else if main_sku != "" {
		data, err = repo.FindByItem(main_sku, pageSize, page)
	} else {
		if storeLocation != 0 {
			datum, err = repo.FindByID(uint64(storeLocation))
			if datum != nil {
				data = append(data, *datum)
			}
			recount = 1
		} else {
			data, err = repo.FindAll(pageSize, page)
			recount, _ = repo.CountAll()

		}

	}

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data, "recount": recount}})

}

func DropItemStockLocation(c *gin.Context) {
	location_id, _ := strconv.Atoi(c.DefaultQuery("location_id", "0"))
	var error error
	error = errors.New("Se ha enviado id = 0")
	if location_id != 0 {
		error = ItemStockLocationService.NewItemStockLocationService().DropItemLocation(uint64(location_id))
	}
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": error.Error()}})

	} else {
		c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Se ha eliminado el registro de id ": location_id}})

	}

}

func PostItemStockLocation(c *gin.Context) {
	var requestBody dtos.ItemStockLocationCreateReq
	var idArray []uint64

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
		model := models.ItemLocation{
			ItemMainSku:     dataItem.ItemMainSku,
			StoreLocationID: dataItem.StoreLocationID,
			Stock:           dataItem.Stock,
		}
		repo.Create(&model)
		idArray = append(idArray, model.ID)
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created", "CreatedIds": idArray}})

}

func PatchItemStockLocation(c *gin.Context) {
	var requestBody dtos.ItemStockLocationUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if requestObject.ItemMainSku != nil {
			model.ItemMainSku = *requestObject.ItemMainSku
		}
		if requestObject.StoreLocationID != nil {
			model.StoreLocationID = *requestObject.StoreLocationID
		}
		if requestObject.Stock != nil {
			model.Stock = *requestObject.Stock
		}
		error := repo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}

// n las request hay que enviarle el stock de articulo nuevo el total es decir si tenemos 100 y le restas 4 le mandamos 96
func StockChanges(c *gin.Context) {

	var requestBody dtos.ItemStockLocationStockChangeRequest
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	repoLoc := storelocationrepo.NewStoreLocationRepository(db.DB)
	repoStock := storestockrepo.NewStoreStockRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		loc, err1 := repoLoc.FindByID(model.StoreLocationID)
		stock, err2 := repoStock.FindByItemAndStore(model.ItemMainSku, strconv.FormatUint(loc.StoreID, 10))
		if err != nil || err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stock.Amount = ((stock.Amount - int64(model.Stock)) + int64(requestObject.Stock))
		model.Stock = requestObject.Stock
		if requestObject.Stock >= 0 {
			error := repo.Update(model)
			error2 := repoStock.Update(&stock)
			if error != nil && error2 != nil {
				errorList = append(errorList, error)
				errorList = append(errorList, error2)
			}

		} else {

			errorList = append(errorList, errors.New("Stock can't be negative"))
		}

	}
	if len(errorList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errorList}})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}

func StockMovements(c *gin.Context) {

	var requestBody dtos.ItemStockLocationStockMovementRequest
	var errorList []error
	repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	repoLoc := storelocationrepo.NewStoreLocationRepository(db.DB)
	repoStock := storestockrepo.NewStoreStockRepository(db.DB)

	stockMov := func(sku string, location uint64, stockVariant int64) error {
		model, err := repo.FindByItemsAndLocation(sku, location)
		loc, err1 := repoLoc.FindByID(model.StoreLocationID)
		stock, err2 := repoStock.FindByItemAndStore(model.ItemMainSku, strconv.FormatUint(loc.StoreID, 10))
		if err != nil || err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return err
		}

		stock.Amount = (stock.Amount + stockVariant)
		model.Stock = model.Stock + int(stockVariant)
		if model.Stock < 0 {
			return errors.New("Stock can't be negative")

		}
		error := repo.Update(&model)
		error2 := repoStock.Update(&stock)
		if error != nil || error2 != nil {
			errorList = append(errorList, error)
			errorList = append(errorList, error2)
		}
		return nil
	}

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo

	for _, requestObject := range requestBody.Data {
		if requestObject.StoreLocationID1 != requestObject.StoreLocationID2 {
			errMov := stockMov(requestObject.MainSku, requestObject.StoreLocationID1, -int64(requestObject.Stock))
			if errMov == nil {
				errMov2 := stockMov(requestObject.MainSku, requestObject.StoreLocationID2, int64(requestObject.Stock))
				if errMov2 != nil {
					errorList = append(errorList, errMov2)
				}
			} else {
				errorList = append(errorList, errMov)
			}

		} else {
			errorList = append(errorList, errors.New("The locations are the same"))
		}

	}
	if len(errorList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errorList}})
		return

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}
