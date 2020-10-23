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

type EndOrderInputData struct {
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
	RadioSelect     string
	TypZdrojeZapsi  string
}

type EndOrderResponseData struct {
	Result        string
	EndOrderError string
}

func endOrder(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data from page started")
	var data EndOrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData EndOrderResponseData
		responseData.Result = "nok"
		responseData.EndOrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Closing order started")
	actualTimeDivisor := DownloadActualTimeDivisor(data.WorkplaceCode, data.UserInput)
	actualTerminalInputOrder := DownloadActualTerminalInputOrder(data.UserId, data.WorkplaceCode, data.OrderId, data.UserInput)
	logInfo(data.UserInput, "Actual running terminal input order id: "+strconv.Itoa(actualTerminalInputOrder.OID))
	sytelineOrderEnded := CloseOrderInSyteline(data.UserInput, data.OrderInput, data.OperationSelect, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType, data.RadioSelect, actualTimeDivisor, data.TypZdrojeZapsi, actualTerminalInputOrder.DTS)
	if !sytelineOrderEnded {
		logError(data.UserInput, "Order not closed in syteline ")
		var responseData EndOrderResponseData
		responseData.Result = "nok"
		responseData.EndOrderError = "Order not closed in Syteline"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Closing order ended")
		return
	}
	CheckActualDivisorFor(data.WorkplaceCode, data.UserInput)
	zapsiOrderEnded := CloseOrderInZapsi(data.UserId, data.WorkplaceCode, data.NokCount, data.NokType, actualTerminalInputOrder, data.UserInput)
	if !zapsiOrderEnded {
		logError(data.UserInput, "Order not ended in Zapsi")
		var responseData EndOrderResponseData
		responseData.Result = "nok"
		responseData.EndOrderError = "Order not ended in Zapsi"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Closing order ended")
		return
	}
	var responseData EndOrderResponseData
	responseData.Result = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(data.UserInput, "Closing order ended")
	return
}

func CloseOrderInZapsi(userId string, workplaceCode string, nokCount string, nokType string, actualterminalInputOrder TerminalInputOrder, userInput string) bool {
	logInfo(userInput, "Closing terminal input order in Zapsi started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	userIdAsInt, err := strconv.Atoi(userId)
	if err != nil {
		logError(userInput, "Closing terminal input order in Zapsi ended, problem parsing userid: "+err.Error())
		return false
	}
	nokCountAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError(userInput, "Closing terminal input order in Zapsi ended, problem parsing  nok count: "+err.Error())
		return false
	}
	var terminalInputOrder TerminalInputOrder
	db.Model(&terminalInputOrder).Where("OID = ?", actualterminalInputOrder.OID).Updates(map[string]interface{}{"DTE": time.Now(), "Interval": float32(time.Now().Sub(terminalInputOrder.DTS).Seconds())})
	if nokCountAsInt > 0 {
		logInfo(userInput, "Saving "+nokType+" fails to Zapsi")
		var workplace Workplace
		db.Where("Code = ?", workplaceCode).Find(&workplace)
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
	logInfo(userInput, "Closing terminal input order in Zapsi ended")
	return true
}

func CheckActualDivisorFor(workplaceCode string, userInput string) {
	logInfo(userInput, "Checking actual divisor for workplace started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	var terminalInputOrder []TerminalInputOrder
	db.Where("DTE is null").Where("DeviceID = ?", workplace.DeviceID).Find(&terminalInputOrder)
	if len(terminalInputOrder) == 1 {
		logInfo(userInput, "Actual open terminal input order in Zapsi is just one, setting divisor to initial value of 1")
		var device Device
		db.Model(&device).Where("OID = ?", workplace.DeviceID).Update("Setting", "1")
	}
	logInfo(userInput, "Checking actual divisor for workplace ended")
}

func DownloadActualTerminalInputOrder(userId string, workplaceCode string, orderId string, userInput string) TerminalInputOrder {
	logInfo(userInput, "Downloading actual open terminal input order from Zapsi started")
	var terminalInputOrder TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return terminalInputOrder
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	db.Where("OrderId = ?", orderId).Where("DeviceId = ?", workplace.DeviceID).Where("DTE is null").Where("UserId = ?", userId).Find(&terminalInputOrder)
	logInfo(userInput, "Downloading actual open terminal input order from Zapsi ended")
	return terminalInputOrder
}

func CloseOrderInSyteline(userInput string, orderInput string, operationSelect string, workplaceCode string, okCount string, nokCount string, nokType string, radioSelect string, actualTimeDivisor int, typZdrojeZapsi string, dts time.Time) bool {
	logInfo(userInput, "Closing order in Syteline started")
	transferredToSyteline := false
	if typZdrojeZapsi == "0" {
		logInfo(userInput, "Typ Zdroje Zapsi is zero")
		transferredToSyteline = TransferOkAndNokToSyteline(userInput, orderInput, operationSelect, workplaceCode, okCount, nokCount, nokType)
		transferredToSyteline = CloseOrderRecordInSyteline("4", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
	} else {
		logInfo(userInput, "Typ Zdroje Zapsi is not zero")
		transferredToSyteline = TransferOkAndNokToSyteline(userInput, orderInput, operationSelect, workplaceCode, okCount, nokCount, nokType)
		switch radioSelect {
		case "clovek":
			{
				logInfo(userInput, "Radio has value of clovek")
				transferredToSyteline = CloseOrderRecordInSyteline("9", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
				transferredToSyteline = CloseOrderRecordInSyteline("4", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
			}
		case "stroj":
			{
				logInfo(userInput, "Radio has value of stroj")
				transferredToSyteline = CloseOrderRecordInSyteline("9", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
			}
		case "serizeni":
			{
				logInfo(userInput, "Radio has value of serizeni")
				transferredToSyteline = CloseOrderRecordInSyteline("2", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
			}
		}

	}
	logInfo(userInput, "Closing order in Syteline ended")
	return transferredToSyteline
}

func CloseOrderRecordInSyteline(closingNumber string, userInput string, orderInput string, operationSelect string, workplaceCode string, timeDivisor int, dts time.Time) bool {
	logInfo(userInput, "Closing order record in Syteline started")
	splittedOrderInput := strings.Split(orderInput, ".")
	if len(splittedOrderInput) < 2 {
		logError(userInput, "Closing order record in Syteline ended, problem with order, not containing suffix: "+orderInput)
		return false
	}
	order := splittedOrderInput[0]
	suffixAsNumber := splittedOrderInput[1]
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code, time_divisor)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?, null, null, ?);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, closingNumber, order, suffixAsNumber, operationSelect, workplaceCode, sql.NullTime{Time: dts, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, timeDivisor)
	logInfo(userInput, "Closing order record in Syteline ended")
	return true
}
