package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	UserInput          string
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
	RostraError       string
	NokTypes          []SytelineNok
}

func checkWorkplaceInput(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data from page started")
	var data WorkplaceInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData WorkplaceResponseData
		responseData.Result = "nok"
		responseData.WorkplaceError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Data parsed, checking workplace in Syteline started")
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(data.UserInput, "Problem opening database: "+err.Error())
		var responseData WorkplaceResponseData
		responseData.Result = "nok"
		responseData.WorkplaceError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking workplace in Syteline ended")
		return
	}
	anyOpenOrderInZapsi := CheckAnyOpenOrderInZapsi(data.WorkplaceCode, data.UserInput)
	if anyOpenOrderInZapsi {
		logInfo(data.UserInput, data.WorkplaceCode+" has an open order in Zapsi")
		sameOrder, sameUser := CheckSameUserAndSameOrderInZapsi(data.UserId, data.OrderInput, data.OperationSelect, data.WorkplaceCode, data.UserInput)
		if sameOrder {
			logInfo(data.UserInput, data.WorkplaceCode+" has the same open order in Zapsi")
			if sameUser {
				logInfo(data.UserInput, data.WorkplaceCode+" has the same user in Zapsi")
				var responseData WorkplaceResponseData
				nokTypes := GetNokTypesFromSyteline(data.UserInput)
				responseData.NokTypes = nokTypes
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
				logInfo(data.UserInput, "Checking workplace in Syteline ended")
				return
			} else {
				logInfo(data.UserInput, data.WorkplaceCode+" has not the same user in Zapsi")
				if data.TypZdrojeZapsi == "0" {
					logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi equals zero, enabling Clovek and Start")
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
					logInfo(data.UserInput, "Checking workplace in Syteline ended")
					return
				} else {
					logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi does not equal zero, enabling Clovek, Serizeni, Stroj and Start")
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
					logInfo(data.UserInput, "Checking workplace in Syteline ended")
					return
				}
			}
		} else {
			logInfo(data.UserInput, data.WorkplaceCode+" has not the same open order in Zapsi")
			if data.ViceVp == "1" {
				logInfo(data.UserInput, "sytelineOperationSource.Vice_vp equals one")
				if data.TypZdrojeZapsi == "0" {
					logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi equals zero, enabling Clovek and Start")
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
					logInfo(data.UserInput, "Checking workplace in Syteline ended")
					return
				} else {
					logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi does not equal zero, enabling Clovek, Serizeni, Stroj and Start")
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
					logInfo(data.UserInput, "Checking workplace in Syteline ended")
					return
				}
			} else {
				logInfo(data.UserInput, "sytelineOperationSource.Vice_vp does not equals one")
				if data.ParovyDil == "1" {
					logInfo(data.UserInput, "sytelineOperation.ParovyDil equals one")
					if len(data.SeznamParovychDilu) > 0 {
						logInfo(data.UserInput, "sytelineOperation.SeznamParDilu not empty: "+data.SeznamParovychDilu)
						var zapsiProducts = CheckProductsInZapsi(data.SeznamParovychDilu, data.UserInput)
						anyOpenOrderHasOneOfProducts := CheckIfAnyOpenOrderHasOneOfProducts(data.WorkplaceCode, zapsiProducts, data.UserInput)
						if anyOpenOrderHasOneOfProducts {
							logInfo(data.UserInput, "Some open order in Zapsi contains any of pair parts")
							if data.TypZdrojeZapsi == "0" {
								logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi is zero, enabling Clovek and Start")
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
								logInfo(data.UserInput, "Checking workplace in Syteline ended")
								return
							} else {
								logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi is not zero, enabling Clovek, Serizeni, Stroj and Start")
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
								logInfo(data.UserInput, "Checking workplace in Syteline ended")
								return
							}
						} else {
							logInfo(data.UserInput, "No open order in Zapsi contains any of pair parts")
							var responseData WorkplaceResponseData
							responseData.Result = "nok"
							responseData.WorkplaceError = "Žádná otevřená zakázka neobsahuje žádný z dílu"
							writer.Header().Set("Content-Type", "application/json")
							_ = json.NewEncoder(writer).Encode(responseData)
							logInfo(data.UserInput, "Checking workplace in Syteline ended")
							return
						}
					} else {
						logInfo(data.UserInput, "sytelineOperation.SeznamParDilu is empty")
						var responseData WorkplaceResponseData
						responseData.Result = "nok"
						responseData.WorkplaceError = "Seznam párových dílů je prázdný"
						writer.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(writer).Encode(responseData)
						logInfo(data.UserInput, "Checking workplace in Syteline ended")
						return
					}
				} else {
					logInfo(data.UserInput, "sytelineOperation.ParovyDil does not equals one")
					var responseData WorkplaceResponseData
					responseData.Result = "nok"
					responseData.WorkplaceError = "Parametr párový díl se nerovná 1"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo(data.UserInput, "Checking workplace in Syteline ended")
					return
				}
			}
		}
	} else {
		logInfo(data.UserInput, data.WorkplaceCode+" has not any open order in Zapsi")
		if data.JenPrenosMnozstvi == "1" {
			logInfo(data.UserInput, "sytelineOperation.JenPrenosMnozstvi is one, enabling ok and nok")
			var responseData WorkplaceResponseData
			nokTypes := GetNokTypesFromSyteline(data.UserInput)
			responseData.NokTypes = nokTypes
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
			logInfo(data.UserInput, "Checking workplace in Syteline ended")
			return
		} else {
			logInfo(data.UserInput, "sytelineOperation.JenPrenosMnozstvi is not one")
			if data.TypZdrojeZapsi == "0" {
				logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi is zero, enabling Clovek and Start")
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
				logInfo(data.UserInput, "Checking workplace in Syteline ended")
				return
			} else {
				logInfo(data.UserInput, "sytelineWorkplace.typ_zdroje_zapsi is not zero, enabling Clovek, Serizeni, Stroj and Start")
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
				logInfo(data.UserInput, "Checking workplace in Syteline ended")
				return
			}
		}
	}
}

func GetNokTypesFromSyteline(userInput string) []SytelineNok {
	logInfo(userInput, "Downloading nok types from Syteline started")
	var nokTypes []SytelineNok
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return nokTypes
	}
	command := "declare @JePlatny ListYesNoType, @Kod ReasonCodeType = NULL exec [rostra_exports_test].dbo.ZapsiKodyDuvoduZmetkuSp @Kod= @Kod, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		logError(userInput, "Error downloading data from Syteline: "+err.Error())
		return nokTypes
	}
	defer rows.Close()
	for rows.Next() {
		var nokType SytelineNok
		err = rows.Scan(&nokType.Kod, &nokType.Nazev)
		nokTypes = append(nokTypes, nokType)
		if err != nil {
			logError(userInput, "Error downloading data from Syteline: "+err.Error())
			return nokTypes
		}
	}
	logInfo(userInput, "Downloading nok types from Syteline ended with "+strconv.Itoa(len(nokTypes))+" noktypes")
	return nokTypes
}

func CheckIfAnyOpenOrderHasOneOfProducts(workplaceCode string, products []Product, userInput string) bool {
	logInfo(userInput, "Checking for any open order for any products started")
	var terminalInputOrders []TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false
	}
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DTE is null").Find(&terminalInputOrders)
	for _, terminalInputOrder := range terminalInputOrders {
		var zapsiOrder Order
		db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
		for _, zapsiProduct := range products {
			if zapsiProduct.OID == zapsiOrder.ProductID {
				logInfo(userInput, "Checking for any open order for any products ended, result found")
				return true
			}
		}
	}
	logInfo(userInput, "Checking for any open order for any products ended, no result found")
	return false
}

func CheckProductsInZapsi(seznamParovychDilu string, userInput string) []Product {
	logInfo(userInput, "Checking products in Zapsi started")
	var zapsiProducts []Product
	var products []string
	if strings.Contains(seznamParovychDilu, "|") {
		products = strings.Split(seznamParovychDilu, "|")
	} else {
		products = append(products, seznamParovychDilu)
	}
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return zapsiProducts
	}
	for _, product := range products {
		var zapsiProduct Product
		db.Where("Name = ?", product).Find(&zapsiProduct)
		if zapsiProduct.OID > 0 {
			logInfo(userInput, "Product "+product+" already exists")
		} else {
			logInfo(userInput, "Product "+product+" does not exist, creating product")
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
	logInfo(userInput, "Checking products in Zapsi ended")
	return zapsiProducts
}

func CheckSameUserAndSameOrderInZapsi(userId string, orderInput string, operationInput string, workplaceCode string, userInput string) (bool, bool) {
	logInfo(userInput, "Checking for same order and same user started")
	order, suffix := ParseOrder(orderInput, userInput)
	operation := ParseOperation(operationInput, userInput)
	orderName := order + "." + suffix + "-" + operation
	var zapsiUser User
	var thisOrder TerminalInputOrder
	var thisUser TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false, false
	}
	var zapsiOrder Order
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("OrderID = ?", zapsiOrder.OID).Find(&thisOrder)
	db.Where("OID = ?", userId).Find(&zapsiUser)
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Find(&thisUser)
	logInfo(userInput, "Checking for same order and same user ended")
	return thisOrder.OID > 0, thisUser.OID > 0
}

func CheckAnyOpenOrderInZapsi(workplaceCode string, userInput string) bool {
	logInfo(userInput, "Checking for any open order in Zapsi for workplace "+workplaceCode+" started")
	var terminalInputOrder TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false
	}
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserId is not null").Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		logInfo(userInput, "Checking for any open order in Zapsi ended, open order found")
		return true
	}
	logInfo(userInput, "Checking for any open order in Zapsi ended, open order not found")
	return false
}
