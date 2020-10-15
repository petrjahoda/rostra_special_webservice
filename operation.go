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
	OperationInput string
	OrderInput     string
}

type OperationResponseData struct {
	Result         string
	OperationInput string
	OperationError string
	Workplaces     []SytelineWorkplace
}

func checkOperationInput(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check operation", "Started")
	var data OperationInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check operation", "Error parsing input: "+err.Error())
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationInput
		responseData.OperationError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check operation", "Ended with error")
		return
	}
	logInfo("Check operation", "Data: "+data.OperationInput+", "+data.OrderInput)
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationInput
		responseData.OperationError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check operation", "Ended with error")
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	order, suffix := ParseOrder(data.OrderInput)
	operation := ParseOperation(data.OperationInput)

	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny;\n"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		logError("Check operation", "Error: "+err.Error())
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationInput
		responseData.OperationError = "Problem connecting getting data: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check operation", "Ended with error")
		return
	}

	defer rows.Close()
	var sytelineOperation SytelineOperation
	var sytelineWorkplaces []SytelineWorkplace
	var updatedSytelineWorkplaces []SytelineWorkplace
	for rows.Next() {
		err = rows.Scan(&sytelineOperation.pracoviste, &sytelineOperation.pracoviste_popis, &sytelineOperation.uvolneno_op, &sytelineOperation.priznak_mn_2, &sytelineOperation.mn_2_ks, &sytelineOperation.priznak_mn_3, &sytelineOperation.mn_3_ks, &sytelineOperation.jen_prenos_mnozstvi, &sytelineOperation.priznak_nasobnost, &sytelineOperation.nasobnost, &sytelineOperation.parovy_dil, &sytelineOperation.seznamm_par_dilu)
		if err != nil {
			logError("Check operation", "Error: "+err.Error())
			var responseData OperationResponseData
			responseData.Result = "nok"
			responseData.OperationInput = data.OperationInput
			responseData.OperationError = "Problem connecting getting data: " + err.Error()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Check operation", "Ended with error")
			return
		}
	}
	if len(sytelineOperation.pracoviste) > 0 {
		logInfo("Check operation", "Operation found: "+data.OperationInput)
		command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace;\n"
		workplaceRows, err := db.Raw(command).Rows()
		if err != nil {
			logError("Check operation", "Error: "+err.Error())
			var responseData OperationResponseData
			responseData.Result = "nok"
			responseData.OperationInput = data.OperationInput
			responseData.OperationError = "Problem connecting getting data: " + err.Error()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Check operation", "Ended with error")
			return
		}
		defer workplaceRows.Close()
		for workplaceRows.Next() {
			var sytelineWorkplace SytelineWorkplace
			err = workplaceRows.Scan(&sytelineWorkplace.Zapsi_zdroj, &sytelineWorkplace.priznak_mn_1, &sytelineWorkplace.vice_vp, &sytelineWorkplace.SL_prac, &sytelineWorkplace.typ_zdroje_zapsi, &sytelineWorkplace.auto_prevod_mnozstvi, &sytelineWorkplace.mnozstvi_auto_prevodu)
			sytelineWorkplaces = append(sytelineWorkplaces, sytelineWorkplace)
			if err != nil {
				logError("Check operation", "Error: "+err.Error())
				var responseData OperationResponseData
				responseData.Result = "nok"
				responseData.OperationInput = data.OperationInput
				responseData.OperationError = "Problem connecting getting data: " + err.Error()
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				logInfo("Check operation", "Ended with error")
				return
			}
		}
		for _, sytelineWorkplace := range sytelineWorkplaces {
			sytelineWorkplace.Zapsi_zdroj = UpdateZapsiZdrojFor(sytelineWorkplace)
			updatedSytelineWorkplaces = append(updatedSytelineWorkplaces, sytelineWorkplace)
		}
		if len(updatedSytelineWorkplaces) > 0 {
			logInfo("Check operation", "Workplaces found: "+strconv.Itoa(len(updatedSytelineWorkplaces)))
			var responseData OperationResponseData
			responseData.Result = "ok"
			responseData.OperationInput = data.OperationInput
			responseData.Workplaces = updatedSytelineWorkplaces
			responseData.OperationError = "everything ok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Check operation", "Ended successfully")
			return
		} else {
			logInfo("Check operation", "Workplaces not found for "+data.OperationInput)
			var responseData OperationResponseData
			responseData.Result = "nok"
			responseData.OperationInput = data.OperationInput
			responseData.OperationError = "Workplaces not found"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Check operation", "Ended with error")
			return
		}
	} else {
		logInfo("Check operation", "Operation not found for "+data.OperationInput)
		var responseData OperationResponseData
		responseData.Result = "nok"
		responseData.OperationInput = data.OperationInput
		responseData.OperationError = "Operation not found"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check operation", "Ended with error")
		return
	}
}

func ParseOperation(operationid string) string {
	if strings.Contains(operationid, ";") {
		parsedOperation := strings.Split(operationid, ";")
		return parsedOperation[0]
	}
	return operationid
}

func UpdateZapsiZdrojFor(workplace SytelineWorkplace) string {
	logInfo("Check operation", "Updating workplace name: "+workplace.Zapsi_zdroj)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		return ""
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplace.Zapsi_zdroj).Find(&zapsiWorkplace)
	logInfo("Check operation", "Updated to: "+workplace.Zapsi_zdroj+";"+zapsiWorkplace.Name)
	return workplace.Zapsi_zdroj + ";" + zapsiWorkplace.Name
}
