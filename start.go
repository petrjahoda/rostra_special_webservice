package main

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StartOrderInputData struct {
	WorkplaceCode   string
	UserId          string
	OrderInput      string
	OperationSelect string
	RadioSelect     string
	ProductId       string
	OrderId         string
	Nasobnost       string
	TypZdrojeZapsi  string
	UserInput       string
}

type StartOrderResponseData struct {
	Result          string
	StartOrderError string
}

func startOrder(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data from page started")
	var data StartOrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData StartOrderResponseData
		responseData.Result = "nok"
		responseData.StartOrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Create order started")
	terminalInputOrderCreated, numberOfOpenTerminalInputOrder, deviceId, terminalInputOrderDts := CreateTerminalInputOrderInZapsi(data.UserId, data.WorkplaceCode, data.RadioSelect, data.OrderId, data.Nasobnost, data.UserInput)
	if terminalInputOrderCreated {
		actualTimeDivisor := DownloadActualTimeDivisor(data.WorkplaceCode, data.UserInput)
		if numberOfOpenTerminalInputOrder > actualTimeDivisor {
			logInfo(data.UserInput, "There are more open terminal inoput order than divisor, updating")
			UpdateDeviceWithNew(numberOfOpenTerminalInputOrder, deviceId, data.WorkplaceCode, data.UserInput)
		}
		sytelineOrderCreated := CreateOrderInSyteline(data.TypZdrojeZapsi, data.RadioSelect, data.UserInput, data.OrderInput, data.OperationSelect, data.WorkplaceCode, terminalInputOrderDts)
		if sytelineOrderCreated {
			var responseData StartOrderResponseData
			responseData.Result = "ok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Start order ended")
			return
		} else {
			logError(data.UserInput, "Problem creating order in Syteline")
			var responseData StartOrderResponseData
			responseData.Result = "nok"
			responseData.StartOrderError = "Problem creating order in Syteline"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Start order ended")
			return
		}
	} else {
		logError(data.UserInput, "Problem creating terminal input order")
		var responseData StartOrderResponseData
		responseData.Result = "nok"
		responseData.StartOrderError = "Problem creating order in Zapsi"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Start order ended")
		return
	}
}

func CreateOrderInSyteline(typZdrojeZapsi string, radio string, userInput string, orderInput string, operationSelect string, workplaceCode string, terminalInputOrderDts time.Time) bool {
	logInfo(userInput, "Creating order in Syteline started")
	sytelineOrderStarted := false
	if typZdrojeZapsi == "0" {
		logInfo(userInput, "Typ Zdroje Zapsi is zero")
		sytelineOrderStarted = CreateOrderRecordInSyteline("3", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
	} else {
		logInfo(userInput, "Typ Zdroje Zapsi is not zero")
		switch radio {
		case "clovek":
			{
				logInfo(userInput, "Radio has value of clovek")
				sytelineOrderStarted = CreateOrderRecordInSyteline("3", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
				sytelineOrderStarted = CreateOrderRecordInSyteline("8", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
			}
		case "stroj":
			{
				logInfo(userInput, "Radio has value of stroj")
				sytelineOrderStarted = CreateOrderRecordInSyteline("8", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
			}
		case "serizeni":
			{
				logInfo(userInput, "Radio has value of serizeni")
				sytelineOrderStarted = CreateOrderRecordInSyteline("1", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
			}
		}
	}
	logInfo(userInput, "Creating order in Syteline ended")
	return sytelineOrderStarted
}

func CreateOrderRecordInSyteline(closingNumber string, userInput string, orderInput string, operationSelect string, workplaceCode string, terminalInputOrderDts time.Time) bool {
	logInfo(userInput, "Creating order record in Syteline started")
	splittedOrderInput := strings.Split(orderInput, ".")
	if len(splittedOrderInput) < 2 {
		logError(userInput, "Creating order record in Syteline ended, problem with order, not containing suffix: "+orderInput)
		return false
	}
	order := splittedOrderInput[0]
	suffixAsNumber := splittedOrderInput[1]
	timeToInsert := time.Now()
	if terminalInputOrderDts.Before(time.Now()) {
		timeToInsert = terminalInputOrderDts
	}
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?, null, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, closingNumber, order, suffixAsNumber, operationSelect, workplaceCode, sql.NullTime{Time: timeToInsert, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})
	logInfo(userInput, "Creating order record in Syteline ended")
	return true
}

func UpdateDeviceWithNew(numberOfOpenTerminalInputOrder int, deviceId int, workplaceCode string, userInput string) {
	logInfo(userInput, "Updating number of open orders for: "+workplaceCode+" started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var device Device
	db.Model(&device).Where("OID = ?", deviceId).Update("Setting", numberOfOpenTerminalInputOrder)
	logInfo(userInput, "Updating number of open orders ended")
}

func DownloadActualTimeDivisor(workplaceCode string, userInput string) int {
	logInfo(userInput, "Downloading actual time divisor for workplace: "+workplaceCode+" started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return 1
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	var device Device
	db.Where("OID = ?", zapsiWorkplace.DeviceID).Find(&device)
	timeDivisor, err := strconv.Atoi(device.Setting)
	if err != nil {
		logError(userInput, "Downloading actual time divisor ended, problem parsing data: "+err.Error())
		return 1
	}
	logInfo(userInput, "Downloading actual time divisor ended")
	return timeDivisor

}

func CreateTerminalInputOrderInZapsi(userId string, workplaceCode string, radioSelect string, orderId string, nasobnost string, userInput string) (bool, int, int, time.Time) {
	logInfo(userInput, "Creating terminal input order in Zapsi started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false, 0, 0, time.Now()
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var terminalInputOrder TerminalInputOrder
	var existingTerminalInputOrder TerminalInputOrder
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	db.Where("DeviceId = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserId is null").Find(&existingTerminalInputOrder)
	if existingTerminalInputOrder.OID > 0 {
		logInfo(userInput, "Terminal input order without user found, just updating")
		db.Model(&terminalInputOrder).Where("OID = ?", existingTerminalInputOrder.OID).Updates(map[string]interface{}{"OrderID": orderId, "UserID": userId,
			"Cavity": nasobnost, "Note": radioSelect})
		logInfo(userInput, "Creating terminal input order in Zapsi ended")
		return true, existingTerminalInputOrder.OID, zapsiWorkplace.DeviceID, existingTerminalInputOrder.DTS
	} else {
		logInfo(userInput, "Creating new terminal input order")
		orderIdAsInt, err := strconv.Atoi(orderId)
		if err != nil {
			logError(userInput, "Creating terminal input order ended, problem parsing orderId: "+orderId)
			return false, 0, 0, time.Now()
		}
		nasobnostAsInt, err := strconv.Atoi(nasobnost)
		if err != nil {
			logError(userInput, "Creating terminal input order ended, problem parsing nasobnost: "+nasobnost)
			return false, 0, 0, time.Now()
		}
		userIdAsInt, err := strconv.Atoi(userId)
		if err != nil {
			logError(userInput, "Creating terminal input order ended, problem parsing userid: "+nasobnost)
			return false, 0, 0, time.Now()
		}
		terminalInputOrder.DTS = time.Now()
		terminalInputOrder.OrderID = orderIdAsInt
		terminalInputOrder.UserID = userIdAsInt
		terminalInputOrder.DeviceID = zapsiWorkplace.DeviceID
		terminalInputOrder.Interval = 0
		terminalInputOrder.Count = 0
		terminalInputOrder.Fail = 0
		terminalInputOrder.AverageCycle = 0.0
		terminalInputOrder.WorkerCount = 1
		terminalInputOrder.WorkplaceModeID = 1
		terminalInputOrder.Cavity = nasobnostAsInt
		terminalInputOrder.Note = radioSelect
		db.Create(&terminalInputOrder)
	}
	var actualRunningOrder []TerminalInputOrder
	db.Where("DeviceId = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Find(&actualRunningOrder)
	logInfo(userInput, "Creating terminal input order in Zapsi ended")
	return true, len(actualRunningOrder), zapsiWorkplace.DeviceID, time.Now()
}
