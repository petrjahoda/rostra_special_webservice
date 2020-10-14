package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type OrderInputData struct {
	OrderInput string
}

type OrderResponseData struct {
	Result     string
	OrderInput string
	OrderName  string
	OrderId    string
	OrderError string
}

func checkOrderInput(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check order", "Started")
	var data OrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check order", "Error parsing input: "+err.Error())
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check order", "Ended with error")
		return
	}
	logInfo("Check order", "Data: "+data.OrderInput)
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check user", "Problem opening database: "+err.Error())
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended with error")
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	order, suffix := ParseOrder(data.OrderInput)
	command := "declare @JePlatny ListYesNoType, @VP Infobar = N'" + order + "." + suffix + "' exec [rostra_exports_test].dbo.ZapsiKontrolaVPSp @VP= @VP, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		logError("Check order", "Error: "+err.Error())
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Problem getting data from syteline: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended with error")
		return
	}
	defer rows.Close()
	var sytelineOrder SytelineOrder
	for rows.Next() {
		err = rows.Scan(&sytelineOrder.CisloVp, &sytelineOrder.SuffixVp, &sytelineOrder.PolozkaVp, &sytelineOrder.PopisPolVp, &sytelineOrder.priznak_seriova_vyroba)
		if err != nil {
			logError("Check order", "Error: "+err.Error())
		}
	}
	if len(sytelineOrder.CisloVp) > 0 {
		logInfo("Check order", "Order found: "+data.OrderInput)
		var responseData OrderResponseData
		responseData.Result = "ok"
		responseData.OrderInput = order + "." + suffix
		responseData.OrderName = sytelineOrder.PolozkaVp + " " + sytelineOrder.PopisPolVp
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	} else {
		logInfo("Check order", "Order not found for "+data.OrderInput+" for command "+command)
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Výrobní příkaz " + data.OrderInput + " neexistuje, zopakujte zadání"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended with error")
		return
	}
}

func ParseOrder(orderId string) (string, string) {
	if strings.Contains(orderId, ";") {
		splitted := strings.Split(orderId, ";")
		if strings.Contains(splitted[0], "-") {
			splittedOrder := strings.Split(splitted[0], "-")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError("MAIN", "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		} else if strings.Contains(splitted[0], ".") {
			splittedOrder := strings.Split(splitted[0], ".")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError("MAIN", "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		}
	} else {
		if strings.Contains(orderId, "-") {
			splittedOrder := strings.Split(orderId, "-")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError("MAIN", "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		} else if strings.Contains(orderId, ".") {
			splittedOrder := strings.Split(orderId, ".")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError("MAIN", "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		}
	}
	return orderId, "0"
}
