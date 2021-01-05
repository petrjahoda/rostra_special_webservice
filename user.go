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

type UserInputData struct {
	UserInput string
}

type UserResponseData struct {
	Result    string
	UserInput string
	UserName  string
	UserId    string
	UserError string
	TableData []Table
}

type Table struct {
	TerminalInputOrderId               string
	OrderCode                          string
	OrderName                          string
	ProductName                        string
	SytelineWorkplace                  string
	OrderStart                         string
	OrderRequestedTotal                string
	TotalProducedCount                 string
	TerminalInputOrderProducedCount    string
	TotalTransferredCount              string
	TotalNokTransferredCount           string
	TerminalInputOrderTransferredCount string
	WaitingForTransferCount            string
	TotalNokCount                      string
}

func checkUserInput(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data from page started")
	var data UserInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData UserResponseData
		responseData.Result = "nok"
		responseData.UserInput = data.UserInput
		responseData.UserError = "Problem parsing data: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Data parsed, checking user in syteline started for user "+data.UserInput)
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(data.UserInput, "Problem opening database: "+err.Error())
		var responseData UserResponseData
		responseData.Result = "nok"
		responseData.UserInput = data.UserInput
		responseData.UserError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking user in syteline ended")
		return
	}
	var sytelineUser SytelineUser
	command := "declare @Zamestnanec EmpNumType, @JePlatny ListYesNoType, @Jmeno NameType, @Chyba Infobar  Exec [rostra_exports].dbo.ZapsiKontrolaZamSp @Zamestnanec = N'" + data.UserInput + "', @JePlatny = @JePlatny output, @Jmeno = @Jmeno output, @Chyba = @Chyba output select JePlatny = @JePlatny, Jmeno = @Jmeno, Chyba = @Chyba;\n"
	db.Raw(command).Scan(&sytelineUser)
	if sytelineUser.JePlatny == "1" {
		logInfo(data.UserInput, "User found")
		userId := checkUserInZapsi(sytelineUser, data.UserInput)
		tableData := checkOpenOrderInZapsi(userId, data.UserInput)
		var responseData UserResponseData
		responseData.Result = "ok"
		responseData.UserInput = data.UserInput
		responseData.UserId = strconv.Itoa(userId)
		responseData.UserName = sytelineUser.Jmeno.String
		responseData.TableData = tableData
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking user in syteline ended")
		return
	} else {
		logInfo(data.UserInput, "User not found: "+sytelineUser.Chyba.String)
		var responseData UserResponseData
		responseData.Result = "nok"
		responseData.UserInput = data.UserInput
		responseData.UserError = sytelineUser.Chyba.String
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking user in syteline ended")
		return
	}
}

func checkOpenOrderInZapsi(userId int, userInput string) []Table {
	logInfo(userInput, "Checking open orders in Zapsi")
	var dataTable []Table
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return dataTable
	}
	var terminalInputOrder []TerminalInputOrder
	db.Where("DTE is null").Where("UserId = ?", userId).Find(&terminalInputOrder)
	var user User
	db.Where("OID = ?", userId).Find(&user)
	logInfo(userInput, "Found "+strconv.Itoa(len(terminalInputOrder))+" open orders in Zapsi")
	for _, terminalInputOrder := range terminalInputOrder {
		oneTableData := DownloadDataForOrder(terminalInputOrder, user, userInput)
		dataTable = append(dataTable, oneTableData)
	}
	logInfo(userInput, "Updated "+strconv.Itoa(len(dataTable))+" open orders")
	return dataTable
}

func DownloadDataForOrder(terminalInputOrder TerminalInputOrder, user User, userInput string) Table {
	logInfo(userInput, "Downloading data for order with id: "+strconv.Itoa(terminalInputOrder.OID)+" started")
	var oneTableData Table
	if terminalInputOrder.Note == "clovek" {
		oneTableData.OrderCode = "PC"
	} else if terminalInputOrder.Note == "stroj" {
		oneTableData.OrderCode = "PS"
	} else if terminalInputOrder.Note == "serizeni" {
		oneTableData.OrderCode = "SE"
	} else {
		oneTableData.OrderCode = "N/A"
	}
	logInfo(userInput, "Downloading data from Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return oneTableData
	}
	var order Order
	db.Where("OID = ?", terminalInputOrder.OrderID).Find(&order)
	oneTableData.OrderName = order.Name
	oneTableData.OrderRequestedTotal = strconv.Itoa(order.CountRequested)
	var product Product
	db.Where("OID = ?", order.ProductID).Find(&product)
	oneTableData.ProductName = product.Name
	var device Device
	db.Where("OID = ?", terminalInputOrder.DeviceID).Find(&device)
	var workplace Workplace
	db.Where("DeviceID = ?", device.OID).Find(&workplace)
	oneTableData.SytelineWorkplace = workplace.Code + ";" + workplace.Name
	oneTableData.OrderStart = terminalInputOrder.DTS.Format("02.01.2006 15:04:05")
	oneTableData.TerminalInputOrderProducedCount = strconv.Itoa(terminalInputOrder.Count)
	var terminalInputOrders []TerminalInputOrder
	db.Where("OrderId = ?", terminalInputOrder.OrderID).Find(&terminalInputOrders)
	totalCount := 0
	for _, inputOrder := range terminalInputOrders {
		totalCount += inputOrder.Count
	}
	oneTableData.TotalProducedCount = strconv.Itoa(totalCount)

	logInfo(userInput, "Downloading data from Syteline")
	dbSyteline, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return oneTableData
	}
	sqlDBSyteline, err := dbSyteline.DB()
	defer sqlDBSyteline.Close()
	parsedOrderName := order.Name
	if strings.Contains(order.Name, ".") {
		parsedOrderName = strings.Split(order.Name, ".")[0]
	}
	var zapsiTrans []zapsi_trans
	dbSyteline.Where("job = ?", parsedOrderName).Where("qty_complete is not null").Find(&zapsiTrans)
	transferredTotal := 0
	transferredNok := 0
	for _, oneTrans := range zapsiTrans {
		transferredTotal += int(oneTrans.QtyComplete)
		transferredNok += int(oneTrans.QtyScrapped)
	}
	oneTableData.TotalTransferredCount = strconv.Itoa(transferredTotal)
	oneTableData.TotalNokTransferredCount = strconv.Itoa(transferredNok)
	var zapsiTransThisOrder []zapsi_trans
	dbSyteline.Raw("SELECT * FROM [zapsi_trans]  WHERE (job = '" + parsedOrderName + "') AND (qty_complete is not null) AND (trans_date >= '" + terminalInputOrder.DTS.Format("2006-01-02 15:04:05") + "') AND (emp_num = '" + user.Login + "')").Find(&zapsiTransThisOrder)
	transferredTotalThisOrder := 0
	for _, thisTrans := range zapsiTransThisOrder {
		transferredTotalThisOrder += int(thisTrans.QtyComplete)
	}
	oneTableData.TerminalInputOrderTransferredCount = strconv.Itoa(transferredTotalThisOrder)
	forSave := terminalInputOrder.Count - transferredTotalThisOrder
	oneTableData.WaitingForTransferCount = strconv.Itoa(forSave)
	logInfo(userInput, "Downloading data for order with id: "+strconv.Itoa(terminalInputOrder.OID)+" ended")
	return oneTableData
}

func checkUserInZapsi(user SytelineUser, userInput string) int {
	logInfo(userInput, "Checking user in Zapsi started")
	userFirstName := strings.Split(user.Jmeno.String, ",")[0]
	userSecondName := strings.Split(user.Jmeno.String, ",")[1]
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return 0
	}
	var zapsiUser User
	db.Where("Login LIKE ?", userInput).Find(&zapsiUser)
	if zapsiUser.OID > 0 {
		logInfo(userInput, "User already exists")
		return zapsiUser.OID
	}
	logInfo(userInput, "User does not exist, creating")
	zapsiUser.Login = userInput
	zapsiUser.FirstName = userFirstName
	zapsiUser.Name = userSecondName
	zapsiUser.UserRoleID = "1"
	zapsiUser.UserTypeID = "1"
	db.Create(&zapsiUser)
	var newUser User
	db.Where("Login LIKE ?", userInput).Find(&newUser)
	logInfo(userInput, "Checking user in Zapsi ended")
	return newUser.OID
}
