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
	logInfo("MAIN", "Parsing data from page started")
	var data TransferOrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData TransferOrderResponseData
		responseData.Result = "nok"
		responseData.TransferOrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Transfer order started")
	thisOpenOrderInZapsi := CheckThisOpenOrderInZapsi(data.UserId, data.OrderId, data.WorkplaceCode, data.UserInput)
	if thisOpenOrderInZapsi.OID == 0 {
		CreateAndCloseOrderInZapsi(data.UserId, data.OrderId, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType, data.Nasobnost, data.UserInput)
	} else {
		UpdateTerminalInputOrder(thisOpenOrderInZapsi, data.UserId, data.OrderId, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType, data.Nasobnost, data.UserInput)
	}
	transferredToSyteline := TransferOkAndNokToSyteline(data.UserInput, data.OrderInput, data.OperationSelect, data.WorkplaceCode, data.OkCount, data.NokCount, data.NokType)
	if transferredToSyteline {
		var responseData TransferOrderResponseData
		responseData.Result = "ok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Transfer order ended")
		return
	} else {
		logError(data.UserInput, "Problem transferring order to Syteline")
		var responseData TransferOrderResponseData
		responseData.Result = "nok"
		responseData.TransferOrderError = "Problem transferring order to Syteline"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Transfer order ended")
		return
	}
}

func UpdateTerminalInputOrder(thisOpenOrder TerminalInputOrder, userId string, orderId string, workplaceCode string, okCount string, nokCount string, nokType string, nasobnost string, userInput string) {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return
	}
	okAsInt, err := strconv.Atoi(okCount)
	if err != nil {
		logError(userInput, "Updating terminal input order failed, problem parsing ok count: "+err.Error())
		return
	}
	nokAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError(userInput, "Updating terminal input order failed, problem parsing nok count: "+err.Error())
		return
	}
	db.Model(&TerminalInputOrder{}).Where(TerminalInputOrder{OID: thisOpenOrder.OID}).Updates(TerminalInputOrder{
		ExtID:  thisOpenOrder.ExtID + okAsInt,
		ExtNum: thisOpenOrder.ExtNum + float32(nokAsInt),
	})
}

func TransferOkAndNokToSyteline(userInput string, orderInput string, operationSelect string, workplaceCode string, okCount string, nokCount string, nokType string) bool {
	logInfo(userInput, "Transferring ok and nok to Syteline started")
	splittedOrderInput := strings.Split(orderInput, ".")
	if len(splittedOrderInput) < 2 {
		logError(userInput, "Transferring ok and nok to Syteline ended, problem with order, not containing suffix: "+orderInput)
		return false
	}
	order := splittedOrderInput[0]
	suffixAsNumber := splittedOrderInput[1]
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return false
	}
	okAsInt, err := strconv.Atoi(okCount)
	if err != nil {
		logError(userInput, "Transferring ok and nok to Syteline ended, problem parsing ok count: "+err.Error())
		return false
	}
	nokAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError(userInput, "Transferring ok and nok to Syteline ended, problem parsing nok count: "+err.Error())
		return false
	}
	parsedFail := strings.Split(nokType, ";")
	if len(parsedFail) < 2 {
		logError(userInput, "Transferring ok and nok to Syteline ended, problem parsing nok type: "+okCount)
		return false
	}
	failBarcode := parsedFail[0]
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, null, null, ?, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, "5", order, suffixAsNumber, operationSelect, workplaceCode, float64(okAsInt), 0.0, 0.0)
	if nokAsInt > 0 {
		db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
			" VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userInput, "5", order, suffixAsNumber, operationSelect, workplaceCode, 0.0, float64(nokAsInt), 0.0, failBarcode)
	}
	logInfo(userInput, "Transferring ok and nok to Syteline ended")
	return true

}

func CreateAndCloseOrderInZapsi(userId string, orderId string, workplaceCode string, okCount string, nokCount string, nokType string, nasobnost string, userInput string) {
	logInfo(userInput, "Creating and closing terminal input order in Zapsi started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return
	}
	orderIdAsInt, err := strconv.Atoi(orderId)
	if err != nil {
		logError(userInput, "Creating and closing terminal input order in Zapsi ended, problem parsing orderid: "+err.Error())
		return
	}
	userIdAsInt, err := strconv.Atoi(userId)
	if err != nil {
		logError(userInput, "Creating and closing terminal input order in Zapsi ended, problem parsing userid: "+err.Error())
		return
	}
	nasobnostAsInt, err := strconv.Atoi(nasobnost)
	if err != nil {
		logError(userInput, "Creating and closing terminal input order in Zapsi ended, problem parsing nasobnost: "+err.Error())
		return
	}
	okCountAsInt, err := strconv.Atoi(okCount)
	if err != nil {
		logError(userInput, "Creating and closing terminal input order in Zapsi ended, problem parsing ok count: "+err.Error())
		return
	}
	nokCountAsInt, err := strconv.Atoi(nokCount)
	if err != nil {
		logError(userInput, "Creating and closing terminal input order in Zapsi ended, problem parsing nok count: "+err.Error())
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
	terminalInputOrder.ExtID = okCountAsInt
	terminalInputOrder.ExtNum = float32(nokCountAsInt)
	db.Create(&terminalInputOrder)
	if nokCountAsInt > 0 {
		logInfo(userInput, "Saving "+nokType+" fails to Zapsi")
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
	logInfo(userInput, "Creating and closing terminal input order in Zapsi ended")
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
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("Transfer order", "Problem opening database: "+err.Error())
		return 0
	}
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

func CheckThisOpenOrderInZapsi(userId string, orderId string, workplaceCode string, userInput string) TerminalInputOrder {
	logInfo(userInput, "Checking for open terminal input order in Zapsi started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	var terminalInputOrder TerminalInputOrder
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return terminalInputOrder
	}
	var workplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&workplace)
	db.Where("OrderId = ?", orderId).Where("UserId = ?", userId).Where("DeviceId = ?", workplace.DeviceID).Where("Dte is null").Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		logInfo(userInput, "Checking for open terminal input order in Zapsi ended, order found")
		return terminalInputOrder
	}
	logInfo(userInput, "Checking for open terminal input order in Zapsi ended, order not found")
	return terminalInputOrder
}
