package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type WorkplaceInputData struct {
	WorkplaceCode      string
	UserId             string
	OrderInput         string
	OperationSelect    string
	ParovyDil          string
	SeznamParovychDilu string
	JenPrenosMnozstvi  string
	TypZdrojeZapsi     string
	ViceVp             string
}

type WorkplaceResponseData struct {
	Result            string
	OkInput           string
	NokInput          string
	StartButton       string
	EndButton         string
	TransferButton    string
	ClovekSelection   string
	SerizeniSelection string
	StrojSelection    string
	WorkplaceError    string
}

func checkWorkplaceInput(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Check workplace", "Started")
	var data WorkplaceInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check workplace", "Error parsing input: "+err.Error())
		var responseData WorkplaceResponseData
		responseData.Result = "nok"
		responseData.WorkplaceError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check workplace", "Ended with error")
		return
	}
	logInfo("Check workplace", "Data: "+data.WorkplaceCode+", "+data.OrderInput+", "+data.UserId+", "+data.OperationSelect+", "+data.ParovyDil+", "+data.SeznamParovychDilu+", "+data.JenPrenosMnozstvi)
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check workplace", "Problem opening database: "+err.Error())
		var responseData WorkplaceResponseData
		responseData.Result = "nok"
		responseData.WorkplaceError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check workplace", "Ended with error")
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	anyOpenOrderInZapsi := CheckAnyOpenOrderInZapsi(data.WorkplaceCode)
	if anyOpenOrderInZapsi {
		logInfo("Check workplace", data.WorkplaceCode+" has some open order in Zapsi")
		sameOrder, sameUser := CheckSameUserAndSameOrderInZapsi(data.UserId, data.OrderInput, data.OperationSelect, data.WorkplaceCode)
		if sameOrder {
			logInfo("Check workplace", data.WorkplaceCode+" has the same open order in Zapsi")
			if sameUser {
				logInfo("Check workplace", data.WorkplaceCode+"has the same open order in Zapsi with the same user")
				logInfo("Check workplace", "Enabling OK and NOK")
				var responseData WorkplaceResponseData
				responseData.Result = "nok"
				responseData.WorkplaceError = "Problem connecting Syteline database: " + err.Error()
				//TOD: send back nok types
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
			} else {
				if data.TypZdrojeZapsi == "0" {
					logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
					logInfo("Check workplace", "Enabling Clovek and Start")
					var responseData WorkplaceResponseData
					responseData.Result = "ok"
					responseData.OkInput = "false"
					responseData.NokInput = "false"
					responseData.StartButton = "true"
					responseData.EndButton = "false"
					responseData.TransferButton = "false"
					responseData.ClovekSelection = "true"
					responseData.SerizeniSelection = "false"
					responseData.StrojSelection = "false"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					return
				} else {
					logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
					logInfo("Check workplace", "Enabling Clovek Serizeni, Stroj and Start")
					var responseData WorkplaceResponseData
					responseData.Result = "ok"
					responseData.OkInput = "false"
					responseData.NokInput = "false"
					responseData.StartButton = "true"
					responseData.EndButton = "false"
					responseData.TransferButton = "false"
					responseData.ClovekSelection = "true"
					responseData.SerizeniSelection = "true"
					responseData.StrojSelection = "true"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					return
				}
			}
		} else {
			logInfo("Check workplace", data.WorkplaceCode+" has not the same order in Zapsi")
			if data.ViceVp == "1" {
				logInfo("Check workplace", "sytelineOperationSource.Vice_vp equals one")
				if data.TypZdrojeZapsi == "0" {
					logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
					logInfo("Check workplace", "Enabling Clovek and Start")
					var responseData WorkplaceResponseData
					responseData.Result = "ok"
					responseData.OkInput = "false"
					responseData.NokInput = "false"
					responseData.StartButton = "true"
					responseData.EndButton = "false"
					responseData.TransferButton = "false"
					responseData.ClovekSelection = "true"
					responseData.SerizeniSelection = "false"
					responseData.StrojSelection = "false"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					return
				} else {
					logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
					logInfo("Check workplace", "Enabling Clovek Serizeni, Stroj and Start")
					var responseData WorkplaceResponseData
					responseData.Result = "ok"
					responseData.OkInput = "false"
					responseData.NokInput = "false"
					responseData.StartButton = "true"
					responseData.EndButton = "false"
					responseData.TransferButton = "false"
					responseData.ClovekSelection = "true"
					responseData.SerizeniSelection = "true"
					responseData.StrojSelection = "true"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					return
				}
			} else {
				logInfo("Check workplace", "sytelineOperationSource.Vice_vp does not equal one")
				if data.ParovyDil == "1" {
					logInfo("Check workplace", "sytelineOperation.ParovyDil equals one")
					if len(data.SeznamParovychDilu) > 0 {
						logInfo("Check workplace", "sytelineOperation.SeznamParDilu not empty: "+data.SeznamParovychDilu)
						var zapsiProducts = CheckProductsInZapsi(data.SeznamParovychDilu)
						anyOpenOrderHasOneOfProducts := CheckIfAnyOpenOrderHasOneOfProducts(data.WorkplaceCode, zapsiProducts)
						if anyOpenOrderHasOneOfProducts {
							logInfo("Check workplace", "Some open order in Zapsi contains any of pair parts")
							if data.TypZdrojeZapsi == "0" {
								logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi is zero")
								logInfo("Check workplace", "Enabling Clovek and Start")
								var responseData WorkplaceResponseData
								responseData.Result = "ok"
								responseData.OkInput = "false"
								responseData.NokInput = "false"
								responseData.StartButton = "true"
								responseData.EndButton = "false"
								responseData.TransferButton = "false"
								responseData.ClovekSelection = "true"
								responseData.SerizeniSelection = "false"
								responseData.StrojSelection = "false"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								return
							} else {
								logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi is not zero")
								logInfo("Check workplace", "Enabling Clovek Serizeni, Stroj and Start")
								var responseData WorkplaceResponseData
								responseData.Result = "ok"
								responseData.OkInput = "false"
								responseData.NokInput = "false"
								responseData.StartButton = "true"
								responseData.EndButton = "false"
								responseData.TransferButton = "false"
								responseData.ClovekSelection = "true"
								responseData.SerizeniSelection = "true"
								responseData.StrojSelection = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								return
							}
						} else {
							logInfo("Check workplace", "No open order in Zapsi contains any of pair parts")
							var responseData WorkplaceResponseData
							responseData.Result = "nok"
							responseData.WorkplaceError = "Žádná otevřená zakázka neobsahuje žádný z dílu"
							writer.Header().Set("Content-Type", "application/json")
							_ = json.NewEncoder(writer).Encode(responseData)
							return
						}
					} else {
						logInfo("Check workplace", "sytelineOperation.SeznamParDilu is empty")
						var responseData WorkplaceResponseData
						responseData.Result = "nok"
						responseData.WorkplaceError = "Seznam párových dílů je prázdný"
						writer.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(writer).Encode(responseData)
						return
					}
				} else {
					logInfo("Check workplace", "sytelineOperation.ParovyDil does not equal one")
					var responseData WorkplaceResponseData
					responseData.Result = "nok"
					responseData.WorkplaceError = "Parametr párový díl se nerovná 1"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					return
				}
			}
		}
	} else {
		logInfo("Check workplace", data.WorkplaceCode+" does not have any open order in Zapsi")
		if data.JenPrenosMnozstvi == "1" {
			logInfo("Check workplace", "sytelineOperation.JenPrenosMnozstvi is one")
			logInfo("Check workplace", "Enabling OK and NOK")
			var responseData WorkplaceResponseData
			responseData.Result = "ok"
			responseData.OkInput = "true"
			responseData.NokInput = "true"
			responseData.StartButton = "false"
			responseData.EndButton = "false"
			responseData.TransferButton = "false"
			responseData.ClovekSelection = "false"
			responseData.SerizeniSelection = "false"
			responseData.StrojSelection = "false"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			return
		} else {
			logInfo("Check workplace", "sytelineOperation.JenPrenosMnozstvi is not one")
			if data.TypZdrojeZapsi == "0" {
				logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi is zero")
				logInfo("Check workplace", "Enabling Clovek and Start")
				var responseData WorkplaceResponseData
				responseData.Result = "ok"
				responseData.OkInput = "false"
				responseData.NokInput = "false"
				responseData.StartButton = "true"
				responseData.EndButton = "false"
				responseData.TransferButton = "false"
				responseData.ClovekSelection = "true"
				responseData.SerizeniSelection = "false"
				responseData.StrojSelection = "false"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				return
			} else {
				logInfo("Check workplace", "sytelineWorkplace.typ_zdroje_zapsi is not zero")
				logInfo("Check workplace", "Enabling Clovek Serizeni, Stroj and Start")
				var responseData WorkplaceResponseData
				responseData.Result = "ok"
				responseData.OkInput = "false"
				responseData.NokInput = "false"
				responseData.StartButton = "true"
				responseData.EndButton = "false"
				responseData.TransferButton = "false"
				responseData.ClovekSelection = "true"
				responseData.SerizeniSelection = "true"
				responseData.StrojSelection = "true"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				return
			}
		}
	}
}

func CheckIfAnyOpenOrderHasOneOfProducts(workplaceCode string, products []Product) bool {
	var terminalInputOrders []TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DTE is null").Find(&terminalInputOrders)
	for _, terminalInputOrder := range terminalInputOrders {
		var zapsiOrder Order
		db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
		for _, zapsiProduct := range products {
			if zapsiProduct.OID == zapsiOrder.ProductID {
				return true
			}
		}
	}
	return false
}

func CheckProductsInZapsi(seznamParovychDilu string) []Product {
	var zapsiProducts []Product
	var products []string
	if strings.Contains(seznamParovychDilu, "|") {
		products = strings.Split(seznamParovychDilu, "|")
	} else {
		products = append(products, seznamParovychDilu)
	}
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		return zapsiProducts
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	for _, product := range products {
		var zapsiProduct Product
		db.Where("Name = ?", product).Find(&zapsiProduct)
		if zapsiProduct.OID > 0 {
			logInfo("Check workplace", "Product "+product+" already exists")
		} else {
			logInfo("Check workplace", "Product "+product+" does not exist, creating product")
			zapsiProduct.Name = product
			zapsiProduct.Barcode = product
			zapsiProduct.Cycle = 1
			zapsiProduct.IdleFromTime = 1
			zapsiProduct.ProductGroupID = 1
			zapsiProduct.ProductStatusID = 1
			db.Create(&zapsiProduct)
		}
	}
	for _, product := range products {
		var zapsiProduct Product
		db.Where("Name = ?", product).Find(&zapsiProduct)
		zapsiProducts = append(zapsiProducts, zapsiProduct)
	}
	return zapsiProducts
}

func CheckSameUserAndSameOrderInZapsi(userId string, orderInput string, operationInput string, workplaceCode string) (bool, bool) {
	order, suffix := ParseOrder(orderInput)
	operation := ParseOperation(operationInput)
	orderName := order + "." + suffix + "-" + operation
	var zapsiUser User

	var thisOrder TerminalInputOrder
	var thisUser TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		return false, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiOrder Order
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("OrderID = ?", zapsiOrder.OID).Find(&thisOrder)

	db.Where("OID = ?", userId).Find(&zapsiUser)
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Find(&thisUser)
	return thisOrder.OID > 0, thisUser.OID > 0
}

func CheckAnyOpenOrderInZapsi(workplaceCode string) bool {
	var terminalInputOrder TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserId is not null").Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		return true
	}
	return false
}
