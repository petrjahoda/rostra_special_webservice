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
	logInfo("Start order", "Started")
	var data StartOrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Start order", "Error parsing input: "+err.Error())
		var responseData StartOrderResponseData
		responseData.Result = "nok"
		responseData.StartOrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Start order", "Ended with error")
		return
	}
	logInfo("Start order", "Data: workplaceCode: "+data.WorkplaceCode+", userId: "+data.UserId+",  order: "+data.OrderInput+", operation: "+data.OperationSelect+", radio: "+data.RadioSelect+", productId: "+data.ProductId+", orderId: "+data.OrderId)
	terminalInputOrderCreated, numberOfOpenTerminalInputOrder, deviceId, terminalInputOrderDts := StartTerminalInputOrderInZapsi(data.UserId, data.WorkplaceCode, data.RadioSelect, data.OrderId, data.Nasobnost)
	if terminalInputOrderCreated {
		actualTimeDivisor := GetActualTimeDivisor(data.WorkplaceCode)
		if numberOfOpenTerminalInputOrder > actualTimeDivisor {
			UpdateDeviceWithNew(numberOfOpenTerminalInputOrder, deviceId, data.WorkplaceCode)
		}
		sytelineOrderCreated := StartOrderInSyteline(data.TypZdrojeZapsi, data.RadioSelect, data.UserInput, data.OrderInput, data.OperationSelect, data.WorkplaceCode, terminalInputOrderDts)
		if sytelineOrderCreated {
			var responseData StartOrderResponseData
			responseData.Result = "ok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Start order", "Ended successfully")
			return
		} else {
			logError("Start order", "Problem creating order in Syteline")
			var responseData StartOrderResponseData
			responseData.Result = "nok"
			responseData.StartOrderError = "Problem creating order in Syteline"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Start order", "Ended with error")
			return
		}
	} else {
		logError("Start order", "Problem creating terminal input order")
		var responseData StartOrderResponseData
		responseData.Result = "nok"
		responseData.StartOrderError = "Problem creating order in Zapsi"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Start order", "Ended with error")
		return
	}
}

func StartOrderInSyteline(typZdrojeZapsi string, radio string, userInput string, orderInput string, operationSelect string, workplaceCode string, terminalInputOrderDts time.Time) bool {
	sytelineOrderStarted := false
	if typZdrojeZapsi == "0" {
		sytelineOrderStarted = StartOrderRecordInSyteline("3", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
	} else {
		switch radio {
		case "clovek":
			{
				sytelineOrderStarted = StartOrderRecordInSyteline("3", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
				sytelineOrderStarted = StartOrderRecordInSyteline("8", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
			}
		case "stroj":
			{
				sytelineOrderStarted = StartOrderRecordInSyteline("8", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
			}
		case "serizeni":
			{
				sytelineOrderStarted = StartOrderRecordInSyteline("1", userInput, orderInput, operationSelect, workplaceCode, terminalInputOrderDts)
			}
		}
	}
	return sytelineOrderStarted
}

func StartOrderRecordInSyteline(closingNumber string, userInput string, orderInput string, operationSelect string, workplaceCode string, terminalInputOrderDts time.Time) bool {
	logInfo("Start order", "Creating order in Syteline")
	splittedOrderInput := strings.Split(orderInput, ".")
	if len(splittedOrderInput) < 2 {
		logError("Start order", "Problem with order, not containing suffix: "+orderInput)
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
		logError("Start order", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?, null, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, closingNumber, order, suffixAsNumber, operationSelect, workplaceCode, sql.NullTime{Time: timeToInsert, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})
	return true
}

func UpdateDeviceWithNew(numberOfOpenTerminalInputOrder int, deviceId int, workplaceCode string) {
	logInfo("Start order", "Updating number of open orders for: "+workplaceCode)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Start order", "Problem opening database: "+err.Error())
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var device Device
	db.Model(&device).Where("OID = ?", deviceId).Update("Setting", numberOfOpenTerminalInputOrder)
}

func GetActualTimeDivisor(workplaceCode string) int {
	logInfo("Start order", "Getting actual time divisor for workplace: "+workplaceCode)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Start order", "Problem opening database: "+err.Error())
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
		logError("Start order", "Problem parsing time_divisor from device")
		return 1
	}
	return timeDivisor

}

func StartTerminalInputOrderInZapsi(userId string, workplaceCode string, radioSelect string, orderId string, nasobnost string) (bool, int, int, time.Time) {
	logInfo("Start order", "Creating terminal input order in Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Start order", "Problem opening database: "+err.Error())
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
		logInfo("Start order", "System terminal_input_order without user exists, just updating")
		db.Model(&terminalInputOrder).Where("OID = ?", existingTerminalInputOrder.OID).Updates(map[string]interface{}{"OrderID": orderId, "UserID": userId,
			"Cavity": nasobnost, "Note": radioSelect})
		return true, existingTerminalInputOrder.OID, zapsiWorkplace.DeviceID, existingTerminalInputOrder.DTS
	} else {
		logInfo("Start order", "Creating new terminal_input_order")
		orderIdAsInt, err := strconv.Atoi(orderId)
		if err != nil {
			logError("Start order", "Problem parsing orderId: "+orderId)
			return false, 0, 0, time.Now()
		}
		nasobnostAsInt, err := strconv.Atoi(nasobnost)
		if err != nil {
			logError("Start order", "Problem parsing nasobnost: "+nasobnost)
			return false, 0, 0, time.Now()
		}
		userIdAsInt, err := strconv.Atoi(userId)
		if err != nil {
			logError("Start order", "Problem parsing userid: "+nasobnost)
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
	return true, len(actualRunningOrder), zapsiWorkplace.DeviceID, time.Now()
}
