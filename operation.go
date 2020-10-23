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

type OperationInputData struct {
	OperationSelect string
	OrderInput      string
	ProductId       string
	UserInput       string
}

type OperationResponseData struct {
	Result             string
	OperationInput     string
	OperationError     string
	Workplaces         []SytelineWorkplace
	ParovyDil          string
	SeznamParovychDilu string
	JenPrenosMnozstvi  string
	PriznakMn2         string
	Mn2Ks              string
	PriznakMn3         string
	Mn3Ks              string
	PriznakNasobnost   string
	Nasobnost          string
	OrderId            string
}

func checkOperationInput(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data from page started")
	var data OperationInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationSelect
		responseData.OperationError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo(data.UserInput, "Data parsed, checking operation in syteline started")
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(data.UserInput, "Problem opening database: "+err.Error())
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationSelect
		responseData.OperationError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking operation in syteline ended")
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	order, suffix := ParseOrder(data.OrderInput, data.UserInput)
	operation := ParseOperation(data.OperationSelect, data.UserInput)
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny;\n"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		logError(data.UserInput, "Error reading data from Syteline: "+err.Error())
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationSelect
		responseData.OperationError = "Problem reading data from Syteline: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking operation in syteline ended")
		return
	}

	defer rows.Close()
	var sytelineOperation SytelineOperation
	var sytelineWorkplaces []SytelineWorkplace
	var updatedSytelineWorkplaces []SytelineWorkplace
	for rows.Next() {
		err = rows.Scan(&sytelineOperation.Pracoviste, &sytelineOperation.PracovistePopis, &sytelineOperation.UvolnenoOp, &sytelineOperation.PriznakMn2, &sytelineOperation.Mn2Ks, &sytelineOperation.PriznakMn3, &sytelineOperation.Mn3Ks, &sytelineOperation.JenPrenosMnozstvi, &sytelineOperation.PriznakNasobnost, &sytelineOperation.Nasobnost, &sytelineOperation.ParovyDil, &sytelineOperation.SeznamParDilu)
		if err != nil {
			logError(data.UserInput, "Error reading data from Syteline: "+err.Error())
			var responseData OperationResponseData
			responseData.Result = "nok"
			responseData.OperationInput = data.OperationSelect
			responseData.OperationError = "Problem reading data from Syteline: " + err.Error()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Checking operation in syteline ended")
			return
		}
	}
	if len(sytelineOperation.Pracoviste) > 0 {
		logInfo(data.UserInput, "Operation "+data.OperationSelect+" found")
		command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace;\n"
		workplaceRows, err := db.Raw(command).Rows()
		if err != nil {
			logError(data.UserInput, "Error reading data from Syteline: "+err.Error())
			var responseData OperationResponseData
			responseData.Result = "nok"
			responseData.OperationInput = data.OperationSelect
			responseData.OperationError = "Problem reading data from Syteline: " + err.Error()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Checking operation in syteline ended")
			return
		}
		defer workplaceRows.Close()
		for workplaceRows.Next() {
			var sytelineWorkplace SytelineWorkplace
			err = workplaceRows.Scan(&sytelineWorkplace.ZapsiZdroj, &sytelineWorkplace.PriznakMn1, &sytelineWorkplace.ViceVp, &sytelineWorkplace.SlPrac, &sytelineWorkplace.TypZdrojeZapsi, &sytelineWorkplace.AutoPrevodMnozstvi, &sytelineWorkplace.MnozstviAutoPrevodu)
			sytelineWorkplaces = append(sytelineWorkplaces, sytelineWorkplace)
			if err != nil {
				logError(data.UserInput, "Error reading data from Syteline: "+err.Error())
				var responseData OperationResponseData
				responseData.Result = "nok"
				responseData.OperationInput = data.OperationSelect
				responseData.OperationError = "Problem reading data from Syteline: " + err.Error()
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				logInfo(data.UserInput, "Checking operation in syteline ended")
				return
			}
		}
		for _, sytelineWorkplace := range sytelineWorkplaces {
			sytelineWorkplace.ZapsiZdroj = UpdateZapsiZdrojFor(sytelineWorkplace, data.UserInput)
			updatedSytelineWorkplaces = append(updatedSytelineWorkplaces, sytelineWorkplace)
		}
		if len(updatedSytelineWorkplaces) > 0 {
			logInfo(data.UserInput, "Found workplaces: "+strconv.Itoa(len(updatedSytelineWorkplaces)))
			orderId := checkOrderInZapsi(data.OrderInput, data.OperationSelect, data.ProductId, sytelineOperation.Mn2Ks, data.UserInput)
			var responseData OperationResponseData
			responseData.Result = "ok"
			responseData.OperationInput = data.OperationSelect
			responseData.OperationError = "everything ok"
			responseData.ParovyDil = sytelineOperation.ParovyDil
			responseData.SeznamParovychDilu = sytelineOperation.SeznamParDilu.String
			responseData.JenPrenosMnozstvi = sytelineOperation.JenPrenosMnozstvi
			responseData.PriznakMn2 = sytelineOperation.PriznakMn2
			responseData.OrderId = strconv.Itoa(orderId)
			if strings.Contains(sytelineOperation.Mn2Ks, ".") {
				responseData.Mn2Ks = sytelineOperation.Mn2Ks[:strings.Index(sytelineOperation.Mn2Ks, ".")]
			} else {
				responseData.Mn2Ks = sytelineOperation.Mn2Ks
			}
			if strings.Contains(sytelineOperation.Mn3Ks, ".") {
				responseData.Mn3Ks = sytelineOperation.Mn3Ks[:strings.Index(sytelineOperation.Mn3Ks, ".")]
			} else {
				responseData.Mn3Ks = sytelineOperation.Mn3Ks
			}
			responseData.PriznakMn3 = sytelineOperation.PriznakMn3
			responseData.PriznakNasobnost = sytelineOperation.PriznakNasobnost
			responseData.Nasobnost = sytelineOperation.Nasobnost
			responseData.Workplaces = updatedSytelineWorkplaces
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Checking operation in syteline ended")
			return
		} else {
			logInfo(data.UserInput, "Workplaces not found for "+data.OperationSelect)
			var responseData OperationResponseData
			responseData.Result = "nok"
			responseData.OperationInput = data.OperationSelect
			responseData.OperationError = "No workplace found for " + data.OperationSelect
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo(data.UserInput, "Checking operation in syteline ended")
			return
		}
	} else {
		logInfo(data.UserInput, "Operation "+data.OperationSelect+" not found")
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationSelect
		responseData.OperationError = "Operation not found for " + data.OperationSelect
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(data.UserInput, "Checking operation in syteline ended")
		return
	}
}

func checkOrderInZapsi(orderInput string, operationSelect string, productId string, mn2Ks string, userInput string) int {
	logInfo(userInput, "Checking order in Zapsi started")
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return 0
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	zapsiOrderName := orderInput + "-" + operationSelect
	var zapsiOrder Order
	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
	if zapsiOrder.OID > 0 {
		logInfo(userInput, "Checking order in Zapsi ended, order "+zapsiOrder.Name+" already exists")
		return zapsiOrder.OID
	}
	logInfo(userInput, "Order "+zapsiOrder.Name+" does not exist, creating order in zapsi")
	productIdAsInt, err := strconv.Atoi(productId)
	if err != nil {
		logError(userInput, "Checking order in Zapsi ended, problem parsing productId: "+productId)
		return 0
	}
	countAsInt, err := strconv.Atoi(mn2Ks)
	if err != nil {
		logError(userInput, "Checking order in Zapsi ended, problem parsing mn2Ks: "+productId)
		return 0
	}
	var newOrder Order
	newOrder.Name = zapsiOrderName
	newOrder.Barcode = zapsiOrderName
	newOrder.ProductID = productIdAsInt
	newOrder.OrderStatusID = 1
	newOrder.CountRequested = countAsInt
	db.Create(&newOrder)
	var returnOrder Order
	db.Where("Name = ?", zapsiOrderName).Find(&returnOrder)
	logInfo(userInput, "Checking order in Zapsi ended")
	return returnOrder.OID
}

func ParseOperation(operationid string, userInput string) string {
	logInfo(userInput, "Parsing operation started")
	if strings.Contains(operationid, ";") {
		parsedOperation := strings.Split(operationid, ";")
		logInfo(userInput, "Parsing operation ended")
		return parsedOperation[0]
	}
	logInfo(userInput, "Parsing operation ended")
	return operationid
}

func UpdateZapsiZdrojFor(workplace SytelineWorkplace, userInput string) string {
	logInfo(userInput, "Updating ZapsiZdroj for workplace"+workplace.ZapsiZdroj)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(userInput, "Problem opening database: "+err.Error())
		return ""
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplace.ZapsiZdroj).Find(&zapsiWorkplace)
	logInfo(userInput, "Zapsi Zdroj updated to: "+workplace.ZapsiZdroj+";"+zapsiWorkplace.Name)
	return workplace.ZapsiZdroj + ";" + zapsiWorkplace.Name
}
