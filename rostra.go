package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RostraMainPage struct {
	Version             string
	Username            string
	Order               string
	Operation           string
	Workplace           string
	UsernameValue       string
	OrderValue          string
	OperationValue      string
	OkValue             string
	NokValue            string
	Workplaces          []SytelineWorkplace
	UserDisabled        string
	OrderDisabled       string
	OperationDisabled   string
	WorkplaceDisabled   string
	UserFocus           string
	OrderFocus          string
	OperationFocus      string
	WorkplaceFocus      string
	StartOrderButton    string
	EndOrderButton      string
	TransferOrderButton string
	Message             string
	Ok                  string
	OkDisabled          string
	OkFocus             string
	Nok                 string
	NokTypes            []SytelineNok
	NokDisabled         string
	NokFocus            string
	RadioDisabled       string
}

const (
	checkUser int = iota
	checkOrder
	checkOperation
	checkWorkplace
	checkOk
	checkNok
	checkBoth
)

func DataInput(writer http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Checking data input")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	_ = r.ParseForm()
	userId := r.Form["userid"]
	orderId := r.Form["orderid"]
	operationId := r.Form["operationid"]
	workplaceId := r.Form["workplaceid"]
	startorder := r.Form["startorder"]
	endorder := r.Form["endorder"]
	transferorder := r.Form["transferorder"]
	noktype := r.Form["noktype"]
	nok := r.Form["nok"]
	ok := r.Form["ok"]
	LogInfo("MAIN", "[OK:NOK]  ["+ok[0]+":"+nok[0]+"]")
	LogInfo("MAIN", "[user:"+userId[0]+"] [order:"+orderId[0]+"] [operation:"+operationId[0]+"] [workplace:"+workplaceId[0]+"]")
	data := RostraMainPage{
		Version:             "version: " + version,
		Username:            "Zadejte prosím své číslo",
		Order:               "",
		Operation:           "",
		Workplace:           "",
		Ok:                  "",
		Nok:                 "",
		OkValue:             "",
		NokValue:            "",
		UserDisabled:        "disabled",
		OrderDisabled:       "disabled",
		OperationDisabled:   "disabled",
		WorkplaceDisabled:   "disabled",
		OkDisabled:          "disabled",
		NokDisabled:         "disabled",
		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
		RadioDisabled:       "disabled",
	}
	if len(startorder) == 1 {
		StartOrderInZapsi(&data, userId, orderId, operationId, workplaceId)
	} else if len(endorder) == 1 {
		LogInfo("MAIN", "Ending order")
		data.Message = "Ending order"
		if len(nok) > 0 && len(ok) > 0 {
			SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
			EndOrderInZapsi(orderId, operationId, userId, workplaceId)
			SaveNokIntoSyteline(nok, noktype)
			SaveOkIntoSyteline(ok)
		} else if len(nok) > 0 {
			SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
			SaveNokIntoSyteline(nok, noktype)
		} else if len(ok) > 0 {
			SaveOkIntoSyteline(ok)
		}
		data.Username = "Zadejte prosím své číslo"
		data.UsernameValue = ""
		data.UserDisabled = ""
		data.UserFocus = "autofocus"

	} else if len(transferorder) == 1 {
		LogInfo("MAIN", "Transferring order")
		data.Message = "Transferring order"
		if len(nok) > 0 && len(ok) > 0 {
			LogInfo("MAIN", "Saving both")
			SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
			SaveNokIntoSyteline(nok, noktype)
			SaveOkIntoSyteline(ok)
		} else if len(nok) > 0 {
			LogInfo("MAIN", "Saving just nok")
			SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
			SaveNokIntoSyteline(nok, noktype)
		} else if len(ok) > 0 {
			LogInfo("MAIN", "Saving just ok")
			SaveOkIntoSyteline(ok)
		}
		data.EndOrderButton = "disabled"
		data.TransferOrderButton = "disabled"
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.OperationValue = operationId[0]
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
		data.NokTypes = nokTypes
		workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
		data.Workplaces = append(data.Workplaces, workplace)
		data.Ok = ""
		data.Nok = ""
		data.OkFocus = "autofocus"
		data.OkDisabled = ""
		data.NokDisabled = ""
	} else {
		inputStep := CheckInputStep(ok, nok, orderId, operationId, workplaceId)
		switch inputStep {
		case checkUser:
			CheckUserInSyteline(userId, &data)
		case checkOrder:
			CheckOrderInSyteline(userId, orderId, &data)
		case checkOperation:
			CheckOperationInSyteline(userId, orderId, operationId, &data)
		case checkWorkplace:
			MakeFirstControls(workplaceId, userId, orderId, operationId, &data)
		case checkBoth:
			CheckOk(userId, orderId, operationId, &data, workplaceId, ok, nok, noktype)
			CheckNok(userId, orderId, operationId, &data, workplaceId, ok, nok, noktype)
		case checkOk:
			CheckOk(userId, orderId, operationId, &data, workplaceId, ok, nok, noktype)
		case checkNok:
			CheckNok(userId, orderId, operationId, &data, workplaceId, ok, nok, noktype)
		}
	}

	if len(data.Workplaces) == 0 {
		LogInfo("MAIN", "No workplaces, adding null workplace")
		workplace := SytelineWorkplace{Zapsi_zdroj: "", priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
		data.Workplaces = append(data.Workplaces, workplace)
	}
	if len(data.NokTypes) == 0 {
		LogInfo("MAIN", "No workplaces, adding null workplace")
		nokType := SytelineNok{Kod: "", Nazev: ""}
		data.NokTypes = append(data.NokTypes, nokType)
	}
	_ = tmpl.Execute(writer, data)
}

func SaveNokIntoZapsi(nok []string, noktype []string, workplaceId []string, userId []string) {
	CreateFailInZapsiIfNotExists(noktype)
	SaveTerminalInputFail(nok, noktype, workplaceId, userId)
}

func SaveTerminalInputFail(nok []string, noktype []string, workplaceId []string, userId []string) {
	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
	var splittedUserName []string
	if strings.Contains(trimmedUserName, ",") {
		splittedUserName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+userId[0])
		splittedUserName = append(splittedUserName, trimmedUserName)
		splittedUserName = append(splittedUserName, trimmedUserName)
	}
	var zapsiFail Fail
	var terminalInputFail TerminalInputFail
	var zapsiWorkplace Workplace
	var zapsiUser User

	pcs, err := strconv.Atoi(nok[0])
	if err != nil {
		LogError("MAIN", "Problem parsing Nok amount when saving terminal input fail")
		return
	}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", noktype[0]).Find(&zapsiFail)
	db.Where("Code = ?", workplaceId[0]).Find(&zapsiWorkplace)
	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
	terminalInputFail.FailID = zapsiFail.OID
	terminalInputFail.DeviceID = zapsiWorkplace.DeviceID
	terminalInputFail.UserID = zapsiUser.OID
	terminalInputFail.DT = time.Now()
	for i := 0; i < pcs; i++ {
		db.Save(&terminalInputFail)
	}
	return
}

func CreateFailInZapsiIfNotExists(noktype []string) {
	var zapsiFail Fail
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", noktype[0]).Find(&zapsiFail)
	if zapsiFail.OID > 0 {
		LogInfo("MAIN", "Fail "+noktype[0]+" already exists")
		return
	}
	LogInfo("MAIN", "Fail "+noktype[0]+" does not exist, creating fail")
	zapsiFail.Name = noktype[0]
	zapsiFail.FailTypeID = 100
	db.Create(&zapsiFail)
	var newZapsiFail Fail
	db.Where("Name = ?", noktype[0]).Find(&newZapsiFail)
	return
}

func SaveOkIntoSyteline(ok []string) {

}

func SaveNokIntoSyteline(nok []string, noktype []string) {
	LogInfo("MAIN", "Saving NOK to Syteline")
	//quantity, err := strconv.Atoi(nok[0])
	//if err != nil {
	//	LogError("MAIN", "Problem parsing count when saving to Syteline")
	//}
	db, err := gorm.Open("mssql", SytelineConnection)
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		return
	}
	defer db.Close()
	db.Exec("set nocount, ansi_nulls, quoted_identifier, arithabort, xact_abort on \nbegin tran\n" +
		" \ninsert into zapsi_trans" +
		"\n( trans_date , emp_num , trans_type , job , suffix , oper_num , wc , qty_complete , qty_scrapped , lot , start_date_time , end_date_time , complete_op , shift , reason_code , time_divisor)" +
		"\nvalues" +
		"\n( '20200430' , N' 500001' , 5 , N'3VP0014981' , 0, 10, N'HLIQ12' , 100.0 , 0.0 , '3VP0014981' , NULL , NULL , 0 , NULL , NULL  , NULL)" +
		"\nif @@error <> 0\nbegin\n  rollback\n  raiserror (N'Chyba zápisu do zapsi_trans', 16, 1)" +
		"\n \nend\nelse\nbegin\n  commit\nend")
}

func CheckNok(userId []string, orderId []string, operationId []string, data *RostraMainPage, workplaceId []string, ok []string, nok []string, noktype []string) {
	LogInfo("MAIN", "Checking NOK: ["+ok[0]+":"+nok[0]+"]")
	sytelineOperation, sytelineWorkplaces := CheckOperationInSyteline(userId, orderId, operationId, data)
	mn1 := CheckForMn1(workplaceId, sytelineWorkplaces)
	mn2 := sytelineOperation.priznak_mn_2 == "1"
	mn3 := sytelineOperation.priznak_mn_3 == "1"
	LogInfo("MAIN", "Priznak mn_1: "+strconv.FormatBool(mn1))
	LogInfo("MAIN", "Priznak mn_2: "+strconv.FormatBool(mn2))
	LogInfo("MAIN", "Priznak mn_3: "+strconv.FormatBool(mn3))
	checkedOk, operatorAmountLessThanInZapsi := CheckIfOperatorAmountLessThanInZapsi(nok, userId, orderId, operationId, workplaceId)
	LogInfo("MAIN", "Operator inserted less amount than in Zapsi: "+strconv.FormatBool(operatorAmountLessThanInZapsi))
	if checkedOk {
		if mn1 && !operatorAmountLessThanInZapsi {
			LogInfo("MAIN", "Mn1 and more amount than in Zapsi, displaying error")
			//TODO: displayError
		}
		if mn2 && !operatorAmountLessThanInZapsi {
			LogInfo("MAIN", "Mn2 and more amount than in Zapsi, displaying error")
			//TODO: displayError
		}
		if mn3 && !operatorAmountLessThanInZapsi {
			LogInfo("MAIN", "Mn3 and more amount than in Zapsi, displaying error")
			//TODO: displayError
		}
		anyOrderExists := CheckAnyOrderInZapsi(workplaceId)
		if anyOrderExists {
			thisOrderIsOpen := CheckThisOrderInZapsi(userId, orderId, operationId, workplaceId)
			if thisOrderIsOpen {
				EnableTransferAndEndButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
			}
		} else {
			sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId, data)
			if sytelineOperation.jen_prenos_mnozstvi == "1" {

				EnableTransferButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
			}
		}
	} else {
		LogError("MAIN", "Problem checking data from Zapsi")
	}
}

func CheckOk(userId []string, orderId []string, operationId []string, data *RostraMainPage, workplaceId []string, ok []string, nok []string, noktype []string) {
	LogInfo("MAIN", "Checking OK")
	sytelineOperation, sytelineWorkplaces := CheckOperationInSyteline(userId, orderId, operationId, data)
	mn1 := CheckForMn1(workplaceId, sytelineWorkplaces)
	mn2 := sytelineOperation.priznak_mn_2 == "1"
	mn3 := sytelineOperation.priznak_mn_3 == "1"
	LogInfo("MAIN", "Priznak mn_1: "+strconv.FormatBool(mn1))
	LogInfo("MAIN", "Priznak mn_2: "+strconv.FormatBool(mn2))
	LogInfo("MAIN", "Priznak mn_3: "+strconv.FormatBool(mn3))
	checkedOk, operatorAmountLessThanInZapsi := CheckIfOperatorAmountLessThanInZapsi(ok, userId, orderId, operationId, workplaceId)
	LogInfo("MAIN", "Operator inserted less amount than in Zapsi: "+strconv.FormatBool(operatorAmountLessThanInZapsi))
	if checkedOk {
		if mn1 && !operatorAmountLessThanInZapsi {
			LogInfo("MAIN", "Mn1 and more amount than in Zapsi, displaying error")
			//TODO: displayError
		}
		if mn2 && !operatorAmountLessThanInZapsi {
			LogInfo("MAIN", "Mn2 and more amount than in Zapsi, displaying error")
			//TODO: displayError
		}
		if mn3 && !operatorAmountLessThanInZapsi {
			LogInfo("MAIN", "Mn3 and more amount than in Zapsi, displaying error")
			//TODO: displayError
		}
		anyOrderExists := CheckAnyOrderInZapsi(workplaceId)
		if anyOrderExists {
			thisOrderIsOpen := CheckThisOrderInZapsi(userId, orderId, operationId, workplaceId)
			if thisOrderIsOpen {
				EnableTransferAndEndButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
			}
		} else {
			sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId, data)
			if sytelineOperation.jen_prenos_mnozstvi == "1" {
				EnableTransferButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
			}
		}
	} else {
		LogError("MAIN", "Problem checking data from Zapsi")
	}
}

func CheckIfOperatorAmountLessThanInZapsi(userAmount []string, userId []string, orderId []string, operationId []string, workplaceId []string) (bool, bool) {
	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
	var splittedUserName []string
	if strings.Contains(trimmedUserName, ",") {
		splittedUserName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+userId[0])
		splittedUserName = append(splittedUserName, trimmedUserName)
		splittedUserName = append(splittedUserName, trimmedUserName)
	}
	order, suffix := ParseOrder(orderId[0])
	orderName := order + "." + suffix + "-" + operationId[0]
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false, false
	}
	defer db.Close()
	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
	okAmount, err := strconv.Atoi(userAmount[0])
	if err != nil {
		LogError("MAIN", "Problem parsing data from user")
		return false, false
	}
	if okAmount < terminalInputOrder.Count {
		return true, true
	}
	return true, false
}

func CheckForMn1(workplaceId []string, workplaces []SytelineWorkplace) bool {
	for _, workplace := range workplaces {
		if workplace.Zapsi_zdroj == workplaceId[0] {
			if workplace.priznak_mn_1 == "1" {
				return true
			}
		}
	}
	return false
}

func StartOrderInZapsi(data *RostraMainPage, userId []string, orderId []string, operationId []string, workplaceId []string) {
	LogInfo("MAIN", "Starting order")
	data.Message = "Starting order"
	sytelineOrder := CheckOrderInSyteline(userId, orderId, data)
	sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId, data)
	CreateProductInZapsiIfNotExists(sytelineOrder)
	zapsiOrder := CreateOrderInZapsiIfNotExists(sytelineOrder, orderId, operationId, sytelineOperation, workplaceId)
	CreateTerminalOrderInZapsi(userId, zapsiOrder, sytelineOperation, workplaceId)
	data.Username = "Zadejte prosím své číslo"
	data.Order = ""
	data.OrderValue = ""
	data.Operation = ""
	data.OperationValue = ""
	data.UsernameValue = ""
	data.UserDisabled = ""
	data.OrderDisabled = "disabled"
	data.OperationDisabled = "disabled"
	data.WorkplaceDisabled = "disabled"
	data.Workplaces = []SytelineWorkplace{}
	data.UserFocus = "autofocus"
}

func CheckAnyOrderInZapsi(workplaceId []string) bool {
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		return true
	}
	return false

}

func EndOrderInZapsi(orderId []string, operationId []string, userId []string, workplaceId []string) {
	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
	var splittedUserName []string
	if strings.Contains(trimmedUserName, ",") {
		splittedUserName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+userId[0])
		splittedUserName = append(splittedUserName, trimmedUserName)
		splittedUserName = append(splittedUserName, trimmedUserName)
	}
	order, suffix := ParseOrder(orderId[0])
	orderName := order + "." + suffix + "-" + operationId[0]
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceId[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
	//TODO: Get OK and NOK pcs and  count average cycle
	if terminalInputOrder.OID > 0 {
		LogInfo("MAIN", "Closing order "+strconv.Itoa(terminalInputOrder.OID))
		db.Model(&terminalInputOrder).Where("OID = ?", terminalInputOrder.OID).UpdateColumn(TerminalInputOrder{DTE: sql.NullTime{Time: time.Now(), Valid: true}})
		db.Model(&terminalInputOrder).Where("OID = ?", terminalInputOrder.OID).UpdateColumn(TerminalInputOrder{Interval: float32(time.Now().Sub(terminalInputOrder.DTS).Seconds())})
	}
}

func CreateTerminalOrderInZapsi(userId []string, order Order, operation SytelineOperation, workplaceId []string) {
	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
	var splittedUserName []string
	if strings.Contains(trimmedUserName, ",") {
		splittedUserName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+userId[0])
		splittedUserName = append(splittedUserName, trimmedUserName)
		splittedUserName = append(splittedUserName, trimmedUserName)
	}
	var zapsiUser User
	var zapsiWorkplace Workplace
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
	var terminalInputOrder TerminalInputOrder
	defer db.Close()
	parsedCavity, err := strconv.Atoi(operation.nasobnost)
	if err != nil {
		LogError("MAIN", "Problem parsing cavity: "+operation.nasobnost)
	}
	terminalInputOrder.DTS = time.Now()
	terminalInputOrder.OrderID = order.OID
	terminalInputOrder.UserID = zapsiUser.OID
	terminalInputOrder.DeviceID = zapsiWorkplace.DeviceID
	terminalInputOrder.Interval = 0
	terminalInputOrder.Count = 0
	terminalInputOrder.Fail = 0
	terminalInputOrder.AverageCycle = 0.0
	terminalInputOrder.WorkerCount = 1
	terminalInputOrder.WorkplaceModeID = 1
	terminalInputOrder.Cavity = parsedCavity
	db.Create(&terminalInputOrder)
}

func CreateProductInZapsiIfNotExists(order SytelineOrder) {
	var zapsiProduct Product
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", order.PolozkaVp).Find(&zapsiProduct)
	if zapsiProduct.OID > 0 {
		LogInfo("MAIN", "Product "+order.PolozkaVp+" already exists")
		return
	}
	LogInfo("MAIN", "Product "+order.PolozkaVp+" does not exist, creating product")
	zapsiProduct.Name = order.PolozkaVp
	zapsiProduct.Barcode = order.PolozkaVp
	zapsiProduct.Cycle = 1
	zapsiProduct.IdleFromTime = 1
	zapsiProduct.ProductGroupID = 1
	zapsiProduct.ProductStatusID = 1
	db.Create(&zapsiProduct)
}

func CreateOrderInZapsiIfNotExists(sytelineOrder SytelineOrder, orderId []string, operationId []string, operation SytelineOperation, workplaceId []string) Order {
	var zapsiOrder Order
	var newOrder Order
	var zapsiProduct Product
	var zapsiWorkplace Workplace
	order, suffix := ParseOrder(orderId[0])
	zapsiOrderName := order + "." + suffix + "-" + operationId[0]
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return zapsiOrder
	}
	defer db.Close()
	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
	if zapsiOrder.OID > 0 {
		LogInfo("MAIN", "Order "+zapsiOrder.Name+" already exists")
		return zapsiOrder
	}
	LogInfo("MAIN", "Order "+zapsiOrder.Name+" does not exist, creating order in zapsi")
	db.Where("Name = ?", sytelineOrder.PolozkaVp).Find(&zapsiProduct)
	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
	countRequestedConverted, err := strconv.ParseFloat(operation.mn_2_ks, 32)
	if err != nil {
		LogError("MAIN", "Problem parsing count for sytelineOrder: "+operation.mn_2_ks)
	}
	newOrder.Name = zapsiOrderName
	newOrder.Barcode = zapsiOrderName
	newOrder.ProductID = zapsiProduct.OID
	newOrder.OrderStatusID = 1
	newOrder.CountRequested = int(countRequestedConverted)
	newOrder.WorkplaceID = zapsiWorkplace.OID
	db.Create(&newOrder)
	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
	return zapsiOrder
}

func CreateUserInZapsiIfNotExists(user SytelineUser, userId []string) {
	trimmedUserName := strings.ReplaceAll(user.Jmeno, " ", "")
	var splittedUserName []string
	if strings.Contains(trimmedUserName, ",") {
		splittedUserName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+user.Jmeno)
		splittedUserName = append(splittedUserName, trimmedUserName)
		splittedUserName = append(splittedUserName, trimmedUserName)
	}
	var zapsiUser User
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name LIKE ?", splittedUserName[0]).Where("FirstName LIKE ?", splittedUserName[1]).Find(&zapsiUser)
	if zapsiUser.OID > 0 {
		LogInfo("MAIN", "User "+user.Jmeno+"already exists")
		return
	}
	LogInfo("MAIN", "User "+user.Jmeno+" does not exist, creating user "+user.Jmeno)
	zapsiUser.Login = userId[0]
	zapsiUser.Name = splittedUserName[0]
	zapsiUser.FirstName = splittedUserName[1]
	zapsiUser.UserRoleID = "1"
	zapsiUser.UserTypeID = "1"
	db.Create(&zapsiUser)
}

func CheckThisOrderInZapsi(userId []string, orderId []string, operationid []string, workplaceId []string) bool {
	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
	var splittedUserName []string
	if strings.Contains(trimmedUserName, ",") {
		splittedUserName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+userId[0])
		splittedUserName = append(splittedUserName, trimmedUserName)
		splittedUserName = append(splittedUserName, trimmedUserName)
	}
	order, suffix := ParseOrder(orderId[0])
	orderName := order + "." + suffix + "-" + operationid[0]
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		return true
	}
	return false
}

func CheckOperationInSyteline(userId []string, orderId []string, operationId []string, data *RostraMainPage) (SytelineOperation, []SytelineWorkplace) {
	order, suffix := ParseOrder(orderId[0])
	LogInfo("MAIN", "Checking operation")
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOperation SytelineOperation
	var sytelineWorkplaces []SytelineWorkplace
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.OperationDisabled = ""
		data.OperationFocus = "autofocus"
		return sytelineOperation, sytelineWorkplaces
	}
	defer db.Close()
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationId[0] + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
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
		data.Message = "Operation found: " + operationId[0]
		data.Operation = operationId[0]
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.OperationValue = operationId[0]
		data.WorkplaceDisabled = ""
	} else {
		LogInfo("MAIN", "Operation not found for "+operationId[0])
		data.Message = "Operation not found for " + operationId[0]
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Operace nenalezena, zadejte prosím znovu"
		data.OperationDisabled = ""
		data.OperationFocus = "autofocus"
		return sytelineOperation, sytelineWorkplaces
	}

	command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationId[0] + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
	workplaceRows, err := db.Raw(command).Rows()
	if err != nil {
		LogError("MAIN", "Error: "+err.Error())
	}
	defer workplaceRows.Close()
	for workplaceRows.Next() {
		var sytelineWorkplace SytelineWorkplace
		err = workplaceRows.Scan(&sytelineWorkplace.Zapsi_zdroj, &sytelineWorkplace.priznak_mn_1, &sytelineWorkplace.vice_vp, &sytelineWorkplace.SL_prac, &sytelineWorkplace.auto_prevod_mnozstvi, &sytelineWorkplace.mnozstvi_auto_prevodu)
		sytelineWorkplaces = append(sytelineWorkplaces, sytelineWorkplace)
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
	}
	if len(sytelineWorkplaces) > 0 {
		data.Workplaces = sytelineWorkplaces
		LogInfo("MAIN", "Workplaces found: "+strconv.Itoa(len(sytelineWorkplaces)))
		data.WorkplaceFocus = "autofocus"
	} else {
		LogInfo("MAIN", "Workplaces not found for "+orderId[0])
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Pracoviště nenalezeny, zadejte prosím znovu"
		data.OperationDisabled = ""
		data.OperationFocus = "autofocus"
	}
	return sytelineOperation, sytelineWorkplaces
}

func ParseOrder(orderId string) (string, string) {
	if strings.Contains(orderId, "-") {
		splittedOrder := strings.Split(orderId, "-")
		return splittedOrder[0], splittedOrder[1]
	} else if strings.Contains(orderId, ".") {
		splittedOrder := strings.Split(orderId, ".")
		return splittedOrder[0], splittedOrder[1]
	}
	return orderId, "0"
}

func CheckOrderInSyteline(userId []string, orderId []string, data *RostraMainPage) SytelineOrder {
	order, suffix := ParseOrder(orderId[0])
	LogInfo("MAIN", "Checking order")
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOrder SytelineOrder
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.Order = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.OrderFocus = "autofocus"
		data.OrderDisabled = ""
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
	if len(sytelineOrder.CisloVp) > 0 {
		LogInfo("MAIN", "Order found: "+orderId[0])
		data.Message = "Order found: " + orderId[0]
		data.Order = orderId[0]
		data.OrderValue = orderId[0]
		data.UsernameValue = userId[0]
		data.OperationFocus = "autofocus"
		data.OperationDisabled = ""
	} else {
		LogInfo("MAIN", "Order not found for "+orderId[0]+" for command "+command)
		data.Message = "Order not found for " + orderId[0] + " for command " + command
		data.UsernameValue = userId[0]
		data.Order = "Číslo nenalezeno, nebo je neplatné, zadejte prosím znovu"
		data.OrderDisabled = ""
		data.OrderFocus = "autofocus"
	}
	return sytelineOrder
}

func CheckUserInSyteline(userId []string, data *RostraMainPage) SytelineUser {
	LogInfo("MAIN", "Checking user")
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineUser SytelineUser
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.Username = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
		return sytelineUser
	}
	defer db.Close()
	command := "declare @Zamestnanec EmpNumType, @JePlatny ListYesNoType, @Jmeno NameType Exec [rostra_exports_test].dbo.ZapsiKontrolaZamSp @Zamestnanec = N'" + userId[0] + "', @JePlatny = @JePlatny output, @Jmeno = @Jmeno output select JePlatny = @JePlatny, Jmeno = @Jmeno"
	row := db.Raw(command).Row()
	err = row.Scan(&sytelineUser.JePlatny, &sytelineUser.Jmeno)
	if sytelineUser.JePlatny == "1" {
		LogInfo("MAIN", "User found: "+userId[0])
		data.Message = "User found: " + userId[0]
		data.Username = sytelineUser.Jmeno
		data.UsernameValue = sytelineUser.Jmeno
		data.OrderDisabled = ""
		data.OrderFocus = "autofocus"
		data.Order = "Zadejte prosím číslo zakázky"
	} else {
		LogInfo("MAIN", "User not found: "+userId[0])
		data.Message = "User not found for " + userId[0]
		data.Username = "Číslo nenalezeno, zadejte prosím znovu"
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
		return sytelineUser
	}
	CreateUserInZapsiIfNotExists(sytelineUser, userId)
	return sytelineUser
}

func CheckInputStep(ok []string, nok []string, orderId []string, operationId []string, workplaceId []string) interface{} {
	if len(ok[0]) > 0 && len(nok[0]) > 0 {
		return checkBoth
	} else if len(ok[0]) > 0 {
		return checkOk
	} else if len(nok[0]) > 0 {
		return checkNok
	} else if orderId[0] == "" && operationId[0] == "" && workplaceId[0] == "" {
		return checkUser
	} else if operationId[0] == "" && workplaceId[0] == "" {
		return checkOrder
	} else if workplaceId[0] == "" {
		return checkOperation
	}
	return checkWorkplace
}

func RostraMainScreen(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Displaying main screen")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := RostraMainPage{
		Version:        "version: " + version,
		Username:       "Zadejte prosím své číslo",
		UsernameValue:  "",
		Order:          "",
		OrderValue:     "",
		Operation:      "",
		OperationValue: "",
		Workplace:      "",
		Ok:             "",
		Nok:            "",
		OkValue:        "",
		NokValue:       "",

		OrderDisabled:     "disabled",
		OperationDisabled: "disabled",
		WorkplaceDisabled: "disabled",
		OkDisabled:        "disabled",
		NokDisabled:       "disabled",

		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
		UserFocus:           "autofocus",
		RadioDisabled:       "disabled",
	}
	if len(data.Workplaces) == 0 {
		LogInfo("MAIN", "No workplaces, adding null workplace")
		workplace := SytelineWorkplace{Zapsi_zdroj: "", priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
		data.Workplaces = append(data.Workplaces, workplace)
	}

	if len(data.NokTypes) == 0 {
		LogInfo("MAIN", "No noktypes, adding null noktype")
		nokType := SytelineNok{Kod: "", Nazev: ""}
		data.NokTypes = append(data.NokTypes, nokType)
	}
	_ = tmpl.Execute(writer, data)
}
