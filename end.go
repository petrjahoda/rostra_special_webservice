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
	logInfo("End order", "Started")
	var data EndOrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("End order", "Error parsing input: "+err.Error())
		var responseData EndOrderResponseData
		responseData.Result = "nok"
		responseData.EndOrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("End order", "Ended with error")
		return
	}
	actualTimeDivisor := GetActualTimeDivisor(data.WorkplaceCode)
	actualTerminalInputOrder := GetActualTerminalInputOrder(data.UserId, data.WorkplaceCode, data.OrderId)
	logInfo("End order", "Actual running order: "+strconv.Itoa(actualTerminalInputOrder.OID))
	sytelineOrderEnded := EndOrderInSyteline(data.UserInput, data.OrderInput, data.OperationSelect, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType, data.RadioSelect, actualTimeDivisor, data.TypZdrojeZapsi, actualTerminalInputOrder.DTS)
	if !sytelineOrderEnded {
		logError("End order", "Order not ended in syteline ")
		var responseData EndOrderResponseData
		responseData.Result = "nok"
		responseData.EndOrderError = "Order not ended in Syteline"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("End order", "Ended with error")
		return
	}
	CheckActualDivisorFor(data.WorkplaceCode)
	zapsiOrderEnded := EndOrderInZapsi(data.UserId, data.WorkplaceCode, data.NokCount, data.NokType, actualTerminalInputOrder)
	if !zapsiOrderEnded {
		logError("End order", "Order not ended in Zapsi")
		var responseData EndOrderResponseData
		responseData.Result = "nok"
		responseData.EndOrderError = "Order not ended in Zapsi"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("End order", "Ended with error")
		return
	}
	logInfo("End order", "Order ended in Zapsi")
	var responseData EndOrderResponseData
	responseData.Result = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("End order", "Ended successfully")
	return
}

func EndOrderInZapsi(userId string, workplaceCode string, nokCount string, nokType string, actualterminalInputOrder TerminalInputOrder) bool {
	logInfo("End order", "Closing terminal input order in Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("End order", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	userIdAsInt, err := strconv.Atoi(userId)
	if err != nil {
		logError("Transfer order", "Problem parsing user id: "+userId)
		return false
	}
	nokCountAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError("Transfer order", "Problem parsing nok count: "+nokCount)
		return false
	}
	var terminalInputOrder TerminalInputOrder
	db.Model(&terminalInputOrder).Where("OID = ?", actualterminalInputOrder.OID).Updates(map[string]interface{}{"DTE": time.Now(), "Interval": float32(time.Now().Sub(terminalInputOrder.DTS).Seconds())})
	if nokCountAsInt > 0 {
		logInfo("Transfer order", "Saving "+nokType+" fails to Zapsi")
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
	return true
}

func CheckActualDivisorFor(workplaceCode string) {
	logInfo("End order", "Getting actual divisor for workplace")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("End order", "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	var terminalInputOrder []TerminalInputOrder
	db.Where("DTE is null").Where("DeviceID = ?", workplace.DeviceID).Find(&terminalInputOrder)
	if len(terminalInputOrder) == 1 {
		logInfo("End order", "Actual running order in Zapsi: 1.... setting divisor to initial value of 1")
		var device Device
		db.Model(&device).Where("OID = ?", workplace.DeviceID).Update("Setting", "1")
	}
}

func GetActualTerminalInputOrder(userId string, workplaceCode string, orderId string) TerminalInputOrder {
	logInfo("End order", "Getting actual open terminal input order from Zapsi")
	var terminalInputOrder TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Transfer order", "Problem opening database: "+err.Error())
		return terminalInputOrder
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	db.Where("OrderId = ?", orderId).Where("DeviceId = ?", workplace.DeviceID).Where("DTE is null").Where("UserId = ?", userId).Find(&terminalInputOrder)
	return terminalInputOrder
}

func EndOrderInSyteline(userInput string, orderInput string, operationSelect string, workplaceCode string, okCount string, nokCount string, nokType string, radioSelect string, actualTimeDivisor int, typZdrojeZapsi string, dts time.Time) bool {
	transferredToSyteline := false
	if typZdrojeZapsi == "0" {
		transferredToSyteline = TransferOkAndNokToSyteline(userInput, orderInput, operationSelect, workplaceCode, okCount, nokCount, nokType)
		transferredToSyteline = CloseOrderRecordInSyteline("4", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
	} else {
		transferredToSyteline = TransferOkAndNokToSyteline(userInput, orderInput, operationSelect, workplaceCode, okCount, nokCount, nokType)
		switch radioSelect {
		case "clovek":
			{
				transferredToSyteline = CloseOrderRecordInSyteline("9", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
				transferredToSyteline = CloseOrderRecordInSyteline("4", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
			}
		case "stroj":
			{
				transferredToSyteline = CloseOrderRecordInSyteline("9", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
			}
		case "serizeni":
			{
				transferredToSyteline = CloseOrderRecordInSyteline("2", userInput, orderInput, operationSelect, workplaceCode, actualTimeDivisor, dts)
			}
		}

	}
	return transferredToSyteline
}

func CloseOrderRecordInSyteline(closingNumber string, userInput string, orderInput string, operationSelect string, workplaceCode string, timeDivisor int, dts time.Time) bool {
	logInfo("End order", "Closing order in Syteline")
	splittedOrderInput := strings.Split(orderInput, ".")
	if len(splittedOrderInput) < 2 {
		logError("End order", "Problem with order, not containing suffix: "+orderInput)
		return false
	}
	order := splittedOrderInput[0]
	suffixAsNumber := splittedOrderInput[1]
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Start order", "Problem opening database: "+err.Error())
		return false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code, time_divisor)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?, null, null, ?);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, closingNumber, order, suffixAsNumber, operationSelect, workplaceCode, sql.NullTime{Time: dts, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, timeDivisor)
	return true
}
