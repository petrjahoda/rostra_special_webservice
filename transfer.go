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

type TransferOrderInputData struct {
	WorkplaceCode   string
	UserId          string
	OrderInput      string
	OperationSelect string
	OrderId         string
	Nasobnost       string
	UserInput       string
	OkCount         string
	NokCount        string
	NokType         string
}

type TransferOrderResponseData struct {
	Result             string
	TransferOrderError string
}

func transferOrder(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Transfer order", "Started")
	var data TransferOrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Transfer order", "Error parsing input: "+err.Error())
		var responseData TransferOrderResponseData
		responseData.Result = "nok"
		responseData.TransferOrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Transfer order", "Ended with error")
		return
	}
	thisOpenOrderInZapsi := CheckThisOpenOrderInZapsi(data.UserId, data.OrderId, data.WorkplaceCode)
	if !thisOpenOrderInZapsi {
		StartWithCloseOrderInZapsi(data.UserId, data.OrderId, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType, data.Nasobnost)
	}
	transferredToSyteline := TransferOkAndNokToSyteline(data.UserInput, data.OrderInput, data.OperationSelect, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType)
	if transferredToSyteline {
		var responseData TransferOrderResponseData
		responseData.Result = "ok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Transfer order", "Ended successfully")
		return
	} else {
		logError("Transfer order", "Problem transferring order to Syteline")
		var responseData TransferOrderResponseData
		responseData.Result = "nok"
		responseData.TransferOrderError = "Problem transferring order to Syteline"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
}

func TransferOkAndNokToSyteline(userInput string, orderInput string, operationSelect string, workplaceCode string, okCount string, nokCount string, nokType string) bool {
	logInfo("Transfer order", "Transferring ok and nok to Syteline")
	splittedOrderInput := strings.Split(orderInput, ".")
	if len(splittedOrderInput) < 2 {
		logError("Transfer order", "Problem with order, not containing suffix: "+orderInput)
		return false
	}
	order := splittedOrderInput[0]
	suffixAsNumber := splittedOrderInput[1]
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Transfer order", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	okAsInt, err := strconv.Atoi(okCount)
	if err != nil {
		logError("Transfer order", "Problem parsing ok count: "+okCount)
		return false
	}
	nokAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError("Transfer order", "Problem parsing nok count: "+nokCount)
		return false
	}
	parsedFail := strings.Split(nokType, ";")
	if len(parsedFail) < 2 {
		logError("Transfer order", "Problem parsing fail, does not contain ; "+nokType)
		return false
	}
	failBarcode := parsedFail[0]
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, null, null, ?, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, "5", order, suffixAsNumber, operationSelect, workplaceCode, float64(okAsInt), 0.0, 0.0)
	if nokAsInt > 0 {
		db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
			" VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, "5", order, suffixAsNumber, operationSelect, workplaceCode, 0.0, float64(nokAsInt), 0.0, failBarcode)
	}
	return true

}

func StartWithCloseOrderInZapsi(userId string, orderId string, workplaceCode string, okCount string, nokCount string, nokType string, nasobnost string) {
	logInfo("Transfer order", "Starting and closing terminal input order in Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Transfer order", "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	orderIdAsInt, err := strconv.Atoi(orderId)
	if err != nil {
		logError("Transfer order", "Problem parsing order id: "+orderId)
		return
	}
	userIdAsInt, err := strconv.Atoi(userId)
	if err != nil {
		logError("Transfer order", "Problem parsing user id: "+userId)
		return
	}
	nasobnostAsInt, err := strconv.Atoi(nasobnost)
	if err != nil {
		logError("Transfer order", "Problem parsing user id: "+nasobnost)
		return
	}
	okCountAsInt, err := strconv.Atoi(okCount)
	if err != nil {
		logError("Transfer order", "Problem parsing ok count: "+okCount)
		return
	}
	nokCountAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError("Transfer order", "Problem parsing nok count: "+nokCount)
		return
	}
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	var terminalInputOrder TerminalInputOrder
	terminalInputOrder.DTS = time.Now()
	terminalInputOrder.DTE = sql.NullTime{Time: time.Now(), Valid: true}
	terminalInputOrder.OrderID = orderIdAsInt
	terminalInputOrder.UserID = userIdAsInt
	terminalInputOrder.DeviceID = workplace.DeviceID
	terminalInputOrder.Interval = 0
	terminalInputOrder.AverageCycle = 0.0
	terminalInputOrder.WorkerCount = 1
	terminalInputOrder.WorkplaceModeID = 1
	terminalInputOrder.Cavity = nasobnostAsInt
	terminalInputOrder.Count = okCountAsInt + nokCountAsInt
	terminalInputOrder.Fail = nokCountAsInt
	db.Create(&terminalInputOrder)
	if nokCountAsInt > 0 {
		logInfo("Transfer order", "Saving "+nokType+" fails to Zapsi")
		failId := CheckFailInZapsi(nokType)
		for i := 0; i < nokCountAsInt; i++ {
			var terminalInputFail TerminalInputFail
			terminalInputFail.DT = time.Now()
			terminalInputFail.FailID = failId
			terminalInputFail.UserID = userIdAsInt
			terminalInputFail.DeviceID = workplace.DeviceID
			terminalInputFail.Note = ""
			db.Create(&terminalInputFail)
		}
	}
}

func CheckFailInZapsi(nokType string) int {
	logInfo("Transfer Order", "Checking fail in Zapsi "+nokType)
	parsedFail := strings.Split(nokType, ";")
	if len(parsedFail) < 2 {
		logError("Transfer order", "Problem parsing fail, does not contain ; "+nokType)
		return 0
	}
	failBarcode := parsedFail[0]
	failName := parsedFail[1]
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Transfer order", "Problem opening database: "+err.Error())
		return 0
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var fail Fail
	db.Where("Barcode = ?", failBarcode).Find(&fail)
	if fail.OID > 0 {
		logInfo("Transfer order", "Found fail: "+fail.Name)
		return fail.OID
	}
	var newFail Fail
	newFail.Name = failName
	newFail.Barcode = failBarcode
	newFail.FailTypeID = 100
	db.Save(&newFail)
	var checkFail Fail
	db.Where("Barcode = ?", failBarcode).Find(&checkFail)
	return checkFail.OID
}

func CheckThisOpenOrderInZapsi(userId string, orderId string, workplaceCode string) bool {
	logInfo("Transfer order", "Checking for open order in Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Transfer order", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	var terminalInputOrder TerminalInputOrder
	db.Where("OrderId = ?", orderId).Where("UserId = ?", userId).Where("DeviceId = ?", workplace.DeviceID).Where("Dte is null").Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		return true
	}
	return false
}
