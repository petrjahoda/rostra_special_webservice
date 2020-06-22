package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func EndOrderInSyteline(userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, radio []string) bool {
	sytelineWorkplace := GetWorkplaceFromSyteline(orderid, operationid, workplaceid)
	sytelineOkAndNokTransferred := false
	sytelineOrderClosed := false
	if sytelineWorkplace.typ_zdroje_zapsi == "0" {
		sytelineOkAndNokTransferred = TransferOkAndNokToSyteline(userid, orderid, operationid, workplaceid, ok, nok, noktype)
		sytelineOrderClosed = CloseOrderRecordInSyteline("4", userid, orderid, operationid, workplaceid)
	} else {
		sytelineOkAndNokTransferred = TransferOkAndNokToSyteline(userid, orderid, operationid, workplaceid, ok, nok, noktype)
		if radio[0] == "clovek" {
			sytelineOrderClosed = CloseOrderRecordInSyteline("9", userid, orderid, operationid, workplaceid)
			sytelineOrderClosed = CloseOrderRecordInSyteline("4", userid, orderid, operationid, workplaceid)
		} else if radio[0] == "stroj" {
			sytelineOrderClosed = CloseOrderRecordInSyteline("9", userid, orderid, operationid, workplaceid)
		} else if radio[0] == "serizeni" {
			sytelineOrderClosed = CloseOrderRecordInSyteline("2", userid, orderid, operationid, workplaceid)
		}
	}
	return sytelineOkAndNokTransferred || sytelineOrderClosed
}

func CloseOrderRecordInSyteline(closingNumber string, userid []string, orderid []string, operationid []string, workplaceid []string) bool {
	order, suffix := ParseOrder(orderid[0])
	operation := ParseOperation(operationid[0])
	suffixAsNumber, _ := strconv.Atoi(suffix)
	operationAsNumber, _ := strconv.Atoi(operation)
	userCode := strings.Split(userid[0], ";")[0]
	okTransferred := TransferCloseOrderToSyteline(closingNumber, userid, orderid, operationid, userCode, order, suffixAsNumber, operationAsNumber, workplaceid)
	return okTransferred
}

func TransferCloseOrderToSyteline(closingNumber string, userid []string, orderid []string, operationid []string, userCode string, order string, suffixAsNumber int, operationAsNumber int, workplaceid []string) bool {
	terminalInputOrder := GetActualOpenOrderForWorkplaces(userid, orderid, operationid, order, workplaceid)

	db, err := gorm.Open("mssql", SytelineConnection)
	if err != nil {
		LogError("MAIN", "Problem with Syteline: "+err.Error())
		return false
	}
	defer db.Close()
	LogInfo("MAIN", "Closing order in Syteline")
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?, null, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userCode, closingNumber, order, suffixAsNumber, operationAsNumber, workplaceid[0], sql.NullTime{Time: terminalInputOrder.DTS, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})
	return true
}

func GetActualOpenOrderForWorkplaces(userid []string, orderid []string, operationid []string, order string, workplaceid []string) TerminalInputOrder {
	userLogin := strings.Split(userid[0], ";")[0]
	order, suffix := ParseOrder(orderid[0])
	operation := ParseOperation(operationid[0])
	orderName := order + "." + suffix + "-" + operation
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return terminalInputOrder
	}
	defer db.Close()
	db.Where("Login = ?", userLogin).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
	return terminalInputOrder
}

func TransferOkAndNokToSyteline(userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string) bool {
	order, suffix := ParseOrder(orderid[0])
	suffixAsNumber, _ := strconv.Atoi(suffix)
	operationAsNumber, _ := strconv.Atoi(operationid[0])
	userCode := strings.Split(userid[0], ";")[0]
	okTransferred := TransferOkRecordToSyteline(ok, userCode, order, suffixAsNumber, operationAsNumber, workplaceid)
	nokTransferred := TransferNokRecordToSyteline(nok, noktype, userCode, order, suffixAsNumber, operationAsNumber, workplaceid)
	return okTransferred || !nokTransferred
}

func TransferNokRecordToSyteline(nok []string, noktype []string, userCode string, order string, suffixAsNumber int, operationAsNumber int, workplaceid []string) bool {
	db, err := gorm.Open("mssql", SytelineConnection)
	if err != nil {
		LogError("MAIN", "Problem with Syteline: "+err.Error())
		return false
	}
	defer db.Close()
	nokAsNumber, nokErr := strconv.Atoi(nok[0])
	if nokErr != nil {
		LogError("MAIN", "Problem parsing nok: "+err.Error())
		return false
	} else if nokAsNumber > 0 {
		LogInfo("MAIN", "Saving NOK to Syteline")
		nokTypes := GetNokTypesFromSyteline()
		failCodeToInsert := "0"
		for _, nokType := range nokTypes {
			if nokType.Nazev == noktype[0] {
				failCodeToInsert = nokType.Kod
				continue
			}
		}
		db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
			" VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userCode, "5", order, suffixAsNumber, operationAsNumber, workplaceid[0], 0.0, float64(nokAsNumber), 0.0, failCodeToInsert)
	}
	return true
}

func TransferOkRecordToSyteline(ok []string, userCode string, order string, suffixAsNumber int, operationAsNumber int, workplaceid []string) bool {
	db, err := gorm.Open("mssql", SytelineConnection)
	if err != nil {
		LogError("MAIN", "Problem with Syteline: "+err.Error())
		return false
	}
	defer db.Close()
	okAsNumber, okErr := strconv.Atoi(ok[0])
	if okErr != nil && okAsNumber > 0 {
		LogError("MAIN", "Problem parsing ok: "+err.Error())
		return false
	} else if okAsNumber > 0 {
		LogInfo("MAIN", "Saving OK to Syteline")
		db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
			" VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, null, null, ?, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userCode, "5", order, suffixAsNumber, operationAsNumber, workplaceid[0], float64(okAsNumber), 0.0, 0.0)
	}
	return true
}

func StartOrderInSyteline(userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) bool {
	sytelineWorkplace := GetWorkplaceFromSyteline(orderid, operationid, workplaceid)
	sytelineOrderStarted := false
	if sytelineWorkplace.typ_zdroje_zapsi == "0" {
		sytelineOrderStarted = StartOrderRecordInSyteline("3", userid, orderid, operationid, workplaceid)
	} else {
		if radio[0] == "clovek" {
			sytelineOrderStarted = StartOrderRecordInSyteline("3", userid, orderid, operationid, workplaceid)
			sytelineOrderStarted = StartOrderRecordInSyteline("8", userid, orderid, operationid, workplaceid)
		} else if radio[0] == "stroj" {
			sytelineOrderStarted = StartOrderRecordInSyteline("8", userid, orderid, operationid, workplaceid)
		} else if radio[0] == "serizeni" {
			sytelineOrderStarted = StartOrderRecordInSyteline("1", userid, orderid, operationid, workplaceid)
		}
	}
	return sytelineOrderStarted
}

func StartOrderRecordInSyteline(closingNumber string, userid []string, orderid []string, operationid []string, workplaceid []string) bool {
	order, suffix := ParseOrder(orderid[0])
	suffixAsNumber, _ := strconv.Atoi(suffix)
	operationAsNumber, _ := strconv.Atoi(operationid[0])
	userCode := strings.Split(userid[0], ";")[0]
	terminalInputOrder := GetActualOpenOrderForWorkplaces(userid, orderid, operationid, order, workplaceid)
	timeToInsert := time.Now()
	if terminalInputOrder.DTS.Before(time.Now()) && terminalInputOrder.OID > 0 {
		timeToInsert = terminalInputOrder.DTS
	}
	db, err := gorm.Open("mssql", SytelineConnection)
	if err != nil {
		LogError("MAIN", "Problem with Syteline: "+err.Error())
		return false
	}
	defer db.Close()
	LogInfo("MAIN", "Starting order in Syteline")
	db.Exec("SET ANSI_WARNINGS OFF;INSERT INTO rostra_exports_test.dbo.zapsi_trans (trans_date, emp_num, trans_type, job, suffix, oper_num, wc, qty_complete, qty_scrapped, start_date_time, end_date_time, complete_op, reason_code)"+
		" VALUES ( ?, ?, ?, ?, ?, ?, ?, null, null, ?, ?, null, null);SET ANSI_WARNINGS ON;", sql.NullTime{Time: time.Now(), Valid: true}, userCode, closingNumber, order, suffixAsNumber, operationAsNumber, workplaceid[0], sql.NullTime{Time: timeToInsert, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})
	return true
}

func GetNokTypesFromSyteline() []SytelineNok {
	var nokTypes []SytelineNok
	db, err := gorm.Open("mssql", SytelineConnection)
	command := "declare @JePlatny ListYesNoType, @Kod ReasonCodeType = NULL exec [rostra_exports_test].dbo.ZapsiKodyDuvoduZmetkuSp @Kod= @Kod, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		LogError("MAIN", "Error: "+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var nokType SytelineNok
		err = rows.Scan(&nokType.Kod, &nokType.Nazev)
		nokTypes = append(nokTypes, nokType)
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
	}
	db.Close()
	return nokTypes
}
func GetWorkplaceFromSyteline(orderid []string, operationid []string, workplaceid []string) SytelineWorkplace {
	LogInfo("MAIN", "Checking operation")
	order, suffix := ParseOrder(orderid[0])
	operation := ParseOperation(operationid[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOperation SytelineOperation
	var sytelineWorkplace SytelineWorkplace
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		return sytelineWorkplace
	}
	defer db.Close()
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		LogError("MAIN", "Error: "+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sytelineOperation.pracoviste, &sytelineOperation.pracoviste_popis, &sytelineOperation.uvolneno_op, &sytelineOperation.priznak_mn_2, &sytelineOperation.mn_2_ks, &sytelineOperation.priznak_mn_3, &sytelineOperation.mn_3_ks, &sytelineOperation.jen_prenos_mnozstvi, &sytelineOperation.priznak_nasobnost, &sytelineOperation.nasobnost, &sytelineOperation.parovy_dil, &sytelineOperation.seznamm_par_dilu)
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
	}
	if len(sytelineOperation.pracoviste) > 0 {
		LogInfo("MAIN", "Operation found: "+operationid[0])
		command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
		workplaceRows, err := db.Raw(command).Rows()
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
		defer workplaceRows.Close()
		for workplaceRows.Next() {
			var sytelineWorkplaceScanned SytelineWorkplace
			err = workplaceRows.Scan(&sytelineWorkplaceScanned.Zapsi_zdroj, &sytelineWorkplaceScanned.priznak_mn_1, &sytelineWorkplaceScanned.vice_vp, &sytelineWorkplaceScanned.SL_prac, &sytelineWorkplaceScanned.typ_zdroje_zapsi, &sytelineWorkplaceScanned.auto_prevod_mnozstvi, &sytelineWorkplaceScanned.mnozstvi_auto_prevodu)
			if err != nil {
				LogError("MAIN", "Error: "+err.Error())
			}
			if sytelineWorkplaceScanned.Zapsi_zdroj == workplaceid[0] {
				sytelineWorkplace = sytelineWorkplaceScanned
			}
		}
	}
	return sytelineWorkplace
}

func GetOperationFromSyteline(orderid []string, operationid []string) SytelineOperation {
	LogInfo("MAIN", "Getting operation from syteline")
	order, suffix := ParseOrder(orderid[0])
	operation := ParseOperation(operationid[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOperation SytelineOperation
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		return sytelineOperation
	}
	defer db.Close()
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		LogError("MAIN", "Error: "+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sytelineOperation.pracoviste, &sytelineOperation.pracoviste_popis, &sytelineOperation.uvolneno_op, &sytelineOperation.priznak_mn_2, &sytelineOperation.mn_2_ks, &sytelineOperation.priznak_mn_3, &sytelineOperation.mn_3_ks, &sytelineOperation.jen_prenos_mnozstvi, &sytelineOperation.priznak_nasobnost, &sytelineOperation.nasobnost, &sytelineOperation.parovy_dil, &sytelineOperation.seznamm_par_dilu)
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
	}
	return sytelineOperation
}

func ParseOperation(operationid string) string {
	if strings.Contains(operationid, ";") {
		parsedOperation := strings.Split(operationid, ";")
		return parsedOperation[0]
	}
	return operationid
}

func GetOrderFromSyteline(orderid []string) SytelineOrder {
	LogInfo("MAIN", "Getting order from Syteline")
	order, suffix := ParseOrder(orderid[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOrder SytelineOrder
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		return sytelineOrder
	}
	defer db.Close()
	command := "declare @JePlatny ListYesNoType, @VP Infobar = N'" + order + "." + suffix + "' exec [rostra_exports_test].dbo.ZapsiKontrolaVPSp @VP= @VP, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	rows, err := db.Raw(command).Rows()
	if err != nil {
		LogError("MAIN", "Error: "+err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sytelineOrder.CisloVp, &sytelineOrder.SuffixVp, &sytelineOrder.PolozkaVp, &sytelineOrder.PopisPolVp, &sytelineOrder.priznak_seriova_vyroba)
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
	}
	return sytelineOrder
}

func CheckOperationInSyteline(writer *http.ResponseWriter, userId []string, orderId []string, operationId []string) {
	LogInfo("MAIN", "Checking operation")
	order, suffix := ParseOrder(orderId[0])
	operation := ParseOperation(operationId[0])
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOperation SytelineOperation
	var sytelineWorkplaces []SytelineWorkplace
	var updatedSytelineWorkplaces []SytelineWorkplace
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.OperationDisabled = ""
		data.OperationFocus = "autofocus"
		data.Username = "disabled"
		data.Message = err.Error()
		LogInfo("MAIN", "Sending error page for order check")
		_ = tmpl.Execute(*writer, data)
	} else {
		defer db.Close()
		command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
		rows, err := db.Raw(command).Rows()
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&sytelineOperation.pracoviste, &sytelineOperation.pracoviste_popis, &sytelineOperation.uvolneno_op, &sytelineOperation.priznak_mn_2, &sytelineOperation.mn_2_ks, &sytelineOperation.priznak_mn_3, &sytelineOperation.mn_3_ks, &sytelineOperation.jen_prenos_mnozstvi, &sytelineOperation.priznak_nasobnost, &sytelineOperation.nasobnost, &sytelineOperation.parovy_dil, &sytelineOperation.seznamm_par_dilu)
			if err != nil {
				LogError("MAIN", "Error: "+err.Error())
			}
		}
		if len(sytelineOperation.pracoviste) > 0 {
			LogInfo("MAIN", "Operation found: "+operationId[0])
			command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operation + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
			workplaceRows, err := db.Raw(command).Rows()
			if err != nil {
				LogError("MAIN", "Error: "+err.Error())
			}
			defer workplaceRows.Close()
			for workplaceRows.Next() {
				var sytelineWorkplace SytelineWorkplace
				err = workplaceRows.Scan(&sytelineWorkplace.Zapsi_zdroj, &sytelineWorkplace.priznak_mn_1, &sytelineWorkplace.vice_vp, &sytelineWorkplace.SL_prac, &sytelineWorkplace.typ_zdroje_zapsi, &sytelineWorkplace.auto_prevod_mnozstvi, &sytelineWorkplace.mnozstvi_auto_prevodu)
				sytelineWorkplaces = append(sytelineWorkplaces, sytelineWorkplace)
				if err != nil {
					LogError("MAIN", "Error: "+err.Error())
				}
				for _, sytelineWorkplace := range sytelineWorkplaces {
					sytelineWorkplace.Zapsi_zdroj = UpdateZapsiZdrojFor(sytelineWorkplace)
					updatedSytelineWorkplaces = append(updatedSytelineWorkplaces, sytelineWorkplace)
				}

			}
			if len(updatedSytelineWorkplaces) > 0 {
				data.Workplaces = updatedSytelineWorkplaces
				LogInfo("MAIN", "Workplaces found: "+strconv.Itoa(len(updatedSytelineWorkplaces)))
				data.UsernameValue = userId[0]
				data.OrderValue = orderId[0]
				if strings.Contains(operationId[0], ";") {
					data.Operation = operationId[0]
					data.OperationValue = operationId[0]
				} else {
					data.Operation = operationId[0] + ";" + sytelineOperation.pracoviste + "-" + sytelineOperation.pracoviste_popis
					data.OperationValue = operationId[0] + ";" + sytelineOperation.pracoviste + "-" + sytelineOperation.pracoviste_popis
				}
				data.WorkplaceDisabled = ""
				data.UserDisabled = "disabled"
				data.WorkplaceFocus = "autofocus"
			} else {
				LogInfo("MAIN", "Workplaces not found for "+orderId[0])
				data.UsernameValue = userId[0]
				data.OrderValue = orderId[0]
				data.Operation = "Pracoviště nenalezeny, zadejte prosím znovu"
				data.OperationDisabled = ""
				data.UserDisabled = "disabled"
				data.OperationFocus = "autofocus"
			}
			_ = tmpl.Execute(*writer, data)
		} else {
			LogInfo("MAIN", "Operation not found for "+operationId[0])
			data.Message = "Operation not found for " + operationId[0]
			data.UsernameValue = userId[0]
			data.OrderValue = orderId[0]
			data.Operation = "Operace nenalezena, zadejte prosím znovu"
			data.OperationDisabled = ""
			data.UserDisabled = "disabled"
			data.OperationFocus = "autofocus"
			_ = tmpl.Execute(*writer, data)
		}
	}
}

func CheckOrderInSyteline(writer *http.ResponseWriter, userId []string, orderId []string) {
	LogInfo("MAIN", "Checking order")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	order, suffix := ParseOrder(orderId[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOrder SytelineOrder
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.Order = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.Message = err.Error()
		data.UserFocus = ""
		data.OrderFocus = "autofocus"
		data.OrderDisabled = ""
		data.UserDisabled = "disabled"
		LogInfo("MAIN", "Sending error page for order check")
		_ = tmpl.Execute(*writer, data)
	} else {
		defer db.Close()
		command := "declare @JePlatny ListYesNoType, @VP Infobar = N'" + order + "." + suffix + "' exec [rostra_exports_test].dbo.ZapsiKontrolaVPSp @VP= @VP, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
		rows, err := db.Raw(command).Rows()
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&sytelineOrder.CisloVp, &sytelineOrder.SuffixVp, &sytelineOrder.PolozkaVp, &sytelineOrder.PopisPolVp, &sytelineOrder.priznak_seriova_vyroba)
			if err != nil {
				LogError("MAIN", "Error: "+err.Error())
			}
		}
		if len(sytelineOrder.CisloVp) > 0 {
			LogInfo("MAIN", "Order found: "+orderId[0])
			data.Order = order + "." + suffix + ";" + sytelineOrder.PolozkaVp + " " + sytelineOrder.PopisPolVp
			data.OrderValue = order + "." + suffix + ";" + sytelineOrder.PolozkaVp + "-" + sytelineOrder.PopisPolVp
			data.UsernameValue = userId[0]
			data.UserFocus = ""
			data.OperationFocus = "autofocus"
			data.UserDisabled = "disabled"
			data.OperationDisabled = ""
		} else {
			LogInfo("MAIN", "Order not found for "+orderId[0]+" for command "+command)
			data.UsernameValue = userId[0]
			data.Order = "Číslo nenalezeno, nebo je neplatné, zadejte prosím znovu"
			data.OrderDisabled = ""
			data.UserDisabled = "disabled"
			data.UserFocus = ""
			data.OrderFocus = "autofocus"
		}
		LogInfo("MAIN", "Sending page for order check")
		_ = tmpl.Execute(*writer, data)
	}
}

func ParseOrder(orderId string) (string, string) {
	if strings.Contains(orderId, ";") {
		splitted := strings.Split(orderId, ";")
		if strings.Contains(splitted[0], "-") {
			splittedOrder := strings.Split(splitted[0], "-")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				LogError("MAIN", "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		} else if strings.Contains(splitted[0], ".") {
			splittedOrder := strings.Split(splitted[0], ".")
			suffixAsNumber, err := strconv.Atoi(splittedOrder[1])
			if err != nil {
				LogError("MAIN", "Problem converting suffix: "+splittedOrder[1])
				return splittedOrder[0], splittedOrder[1]
			}
			return splittedOrder[0], strconv.Itoa(suffixAsNumber)
		}
	}
	return orderId, "0"
}

func CheckUserInSyteline(writer *http.ResponseWriter, userId []string) {
	LogInfo("MAIN", "Checking user")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineUser SytelineUser
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.Username = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.Message = err.Error()
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
		LogInfo("MAIN", "Sending error page for user check")
		_ = tmpl.Execute(*writer, data)
	} else {
		defer db.Close()
		command := "declare @Zamestnanec EmpNumType, @JePlatny ListYesNoType, @Jmeno NameType Exec [rostra_exports_test].dbo.ZapsiKontrolaZamSp @Zamestnanec = N'" + userId[0] + "', @JePlatny = @JePlatny output, @Jmeno = @Jmeno output select JePlatny = @JePlatny, Jmeno = @Jmeno"
		row := db.Raw(command).Row()
		err = row.Scan(&sytelineUser.JePlatny, &sytelineUser.Jmeno)
		if sytelineUser.JePlatny == "1" {
			LogInfo("MAIN", "User found: "+userId[0])
			CreateUserInZapsiIfNotExists(sytelineUser, userId)
			data.Username = userId[0] + ";" + sytelineUser.Jmeno
			data.UsernameValue = userId[0] + ";" + sytelineUser.Jmeno
			data.OrderDisabled = ""
			data.UserFocus = ""
			data.OrderFocus = "autofocus"
			data.UserDisabled = "disabled"
			data.Order = "Zadejte prosím číslo zakázky"
		} else {
			LogInfo("MAIN", userId[0]+", user not found")
			data.Username = "Číslo nenalezeno, zadejte prosím znovu"
			data.UserDisabled = ""
			data.UserFocus = "autofocus"
		}

		LogInfo("MAIN", "Sending page for user check")
		_ = tmpl.Execute(*writer, data)
	}
}
