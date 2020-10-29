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

type OrderInputData struct {
	OrderInput string
	UserInput  string
}

type OrderResponseData struct {
	Result               string
	OrderInput           string
	OrderName            string
	OrderError           string
	PriznakSeriovaVyroba string
	ProductId            string
	Operations           []OperationList
}

func checkOrderInput(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data from page started")
	var data OrderInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Problem parsing data: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Data parsed, checking order in syteline started")
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(data.UserInput, "Problem opening database: "+err.Error())
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking order in syteline ended")
		return
	}
	order, suffix := ParseOrder(data.OrderInput, data.UserInput)
	command := "declare @JePlatny ListYesNoType, @VP Infobar = N'" + order + "." + suffix + "' exec [rostra_exports_test].dbo.ZapsiKontrolaVPSp @VP= @VP, @JePlatny = @JePlatny output select JePlatny = @JePlatny;\n"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		logError(data.UserInput, "Error: "+err.Error())
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Problem getting data from syteline: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking order in syteline ended")
		return
	}
	defer rows.Close()
	var sytelineOrder SytelineOrder
	for rows.Next() {
		err = rows.Scan(&sytelineOrder.CisloVp, &sytelineOrder.SuffixVp, &sytelineOrder.PolozkaVp, &sytelineOrder.PopisPolVp, &sytelineOrder.PriznakSeriovaVyroba)
		if err != nil {
			logError(data.UserInput, "Error: "+err.Error())
		}
	}
	if len(sytelineOrder.CisloVp) > 0 {
		logInfo(data.UserInput, "Order found: "+data.OrderInput+", getting list of operations ")
		command := "declare @CisloVP JobType, @PriponaVP SuffixType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + " exec [rostra_exports_test].dbo.ZapsiOperaceVpSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP;\n"
		rows, err := db.Debug().Raw(command).Rows()
		if err != nil {
			logError(data.UserInput, "Error: "+err.Error())
			var responseData OrderResponseData
			responseData.Result = "nok"
			responseData.OrderInput = data.OrderInput
			responseData.OrderError = "Problem getting data from syteline: " + err.Error()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Checking order in syteline ended")
			return
		}
		var operationList []OperationList
		for rows.Next() {
			var operation OperationList
			err = rows.Scan(&operation.Operace, &operation.Pracoviste, &operation.PracovistePopis)
			if err != nil {
				logError(data.UserInput, "Error: "+err.Error())
			}
			logInfo(data.UserInput, operation.Operace+"-"+operation.Pracoviste+"-"+operation.PracovistePopis)
			operationList = append(operationList, operation)
		}
		logInfo(data.UserInput, "Scanned: "+strconv.Itoa(len(operationList))+" operations")
		productId := checkProductInZapsi(sytelineOrder.PolozkaVp, data.UserInput)
		var responseData OrderResponseData
		responseData.Result = "ok"
		responseData.OrderInput = order + "." + suffix
		responseData.OrderName = sytelineOrder.PolozkaVp + " " + sytelineOrder.PopisPolVp
		responseData.Operations = operationList
		responseData.ProductId = strconv.Itoa(productId)
		responseData.PriznakSeriovaVyroba = sytelineOrder.PriznakSeriovaVyroba
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking order in syteline ended")
		return
	} else {
		logInfo(data.UserInput, "Order not found for "+data.OrderInput+" for command "+command)
		var responseData OrderResponseData
		responseData.Result = "nok"
		responseData.OrderInput = data.OrderInput
		responseData.OrderError = "Výrobní příkaz " + data.OrderInput + " neexistuje, zopakujte zadání"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking order in syteline ended")
		return
	}
}

func checkProductInZapsi(polozkaVp string, userInput string) int {
	logInfo(userInput, "Checking product in Zapsi started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return 0
	}
	var zapsiProduct Product
	db.Where("Name = ?", polozkaVp).Find(&zapsiProduct)
	if zapsiProduct.OID > 0 {
		logInfo(userInput, "Checking product in Zapsi ended, product "+polozkaVp+" already exists")
		return zapsiProduct.OID
	}
	logInfo(userInput, "Product "+polozkaVp+" does not exist, creating product")
	zapsiProduct.Name = polozkaVp
	zapsiProduct.Barcode = polozkaVp
	zapsiProduct.Cycle = 1
	zapsiProduct.IdleFromTime = 1
	zapsiProduct.ProductGroupID = 1
	zapsiProduct.ProductStatusID = 1
	db.Create(&zapsiProduct)
	var newZapsiProduct Product
	db.Where("Name = ?", polozkaVp).Find(&newZapsiProduct)
	logInfo(userInput, "Checking product in Zapsi ended")
	return newZapsiProduct.OID
}

func ParseOrder(orderId string, userInput string) (string, string) {
	logInfo(userInput, "Parsing order started")
	if strings.Contains(orderId, ";") {
		splitted := strings.Split(orderId, ";")
		if strings.Contains(splitted[0], "-") {
			splittedOrder := strings.Split(splitted[0], "-")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError(userInput, "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			logInfo(userInput, "Parsing order ended")
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		} else if strings.Contains(splitted[0], ".") {
			splittedOrder := strings.Split(splitted[0], ".")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError(userInput, "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			logInfo(userInput, "Parsing order ended")
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		}
	} else {
		if strings.Contains(orderId, "-") {
			splittedOrder := strings.Split(orderId, "-")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError(userInput, "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			logInfo(userInput, "Parsing order ended")
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		} else if strings.Contains(orderId, ".") {
			splittedOrder := strings.Split(orderId, ".")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				logError(userInput, "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			logInfo(userInput, "Parsing order ended")
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		}
	}
	return orderId, "0"
}
