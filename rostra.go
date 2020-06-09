package main

import (
	"github.com/julienschmidt/httprouter"
	"strconv"

	//"golang.org/x/crypto/openpgp/packet"
	"html/template"
	"net/http"
)

type RostraMainPage struct {
	Version   string
	Username  string
	Order     string
	Operation string
	Workplace string
	Ok        string
	Nok       string
	Message   string

	UsernameValue  string
	OrderValue     string
	OperationValue string
	OkValue        string
	NokValue       string

	StartOrderButton    string
	EndOrderButton      string
	TransferOrderButton string

	UserFocus      string
	OrderFocus     string
	OperationFocus string
	WorkplaceFocus string
	OkFocus        string
	NokFocus       string

	UserDisabled      string
	OrderDisabled     string
	OperationDisabled string
	WorkplaceDisabled string
	OkDisabled        string
	NokDisabled       string
	RadioDisabled     string

	NokTypes   []SytelineNok
	Workplaces []SytelineWorkplace
}

const (
	checkUserStep int = iota
	checkOrderStep
	checkOperationStep
	checkWorkplaceStep
	checkAmountStep
	startOrderStep
	transferOrderStep
	endOrderStep
)

func DataInput(writer http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Checking data input")
	_ = r.ParseForm()
	userid := r.Form["userid"]
	orderid := r.Form["orderid"]
	operationid := r.Form["operationid"]
	workplaceid := r.Form["workplaceid"]
	startorder := r.Form["startorder"]
	endorder := r.Form["endorder"]
	transferorder := r.Form["transferorder"]
	noktype := r.Form["noktype"]
	nok := r.Form["nok"]
	ok := r.Form["ok"]
	radio := r.Form["syteline"]
	LogInfo("MAIN", "[USER:"+userid[0]+"] [ORDER:"+orderid[0]+"] [OPERATION:"+operationid[0]+"] [WORKPLACE:"+workplaceid[0]+"]")
	LogInfo("MAIN", "[OK:NOK:TYPE]  ["+ok[0]+":"+nok[0]+":"+noktype[0]+"]")
	LogInfo("MAIN", "[RADIO:"+strconv.Itoa(len(radio))+"]")
	LogInfo("MAIN", "[START:TRANSFER:END] ["+strconv.Itoa(len(startorder))+":"+strconv.Itoa(len(transferorder))+":"+strconv.Itoa(len(endorder))+"]")
	inputStep := CheckInputStep(orderid, operationid, workplaceid, startorder, transferorder, endorder, ok, nok)
	switch inputStep {
	case checkUserStep:
		CheckUserInSyteline(&writer, userid)
	case checkOrderStep:
		CheckOrderInSyteline(&writer, userid, orderid)
	case checkOperationStep:
		CheckOperationInSyteline(&writer, userid, orderid, operationid)
	case checkWorkplaceStep:
		FirstControls(&writer, workplaceid, userid, orderid, operationid)
	case checkAmountStep:
		//SecondControls(&writer, workplaceid, userid, orderid, operationid, ok, nok, noktype)
	case startOrderStep:
		//StartOrder(userid, orderid, operationid, workplaceid)
	case transferOrderStep:
		//CheckOk(userid, orderid, operationid, workplaceid, ok, nok, noktype)
		//CheckNok(userid, orderid, operationid, workplaceid, ok, nok, noktype)
		//TransferOrder(userid, orderid, operationid, workplaceid, ok, nok, noktype, radio)
	case endOrderStep:
		//CheckOk(userid, orderid, operationid, workplaceid, ok, nok, noktype)
		//CheckNok(userid, orderid, operationid, workplaceid, ok, nok, noktype)
		//EndOrder(userid, orderid, operationid, workplaceid, ok, nok, noktype, radio)
	}
	//if len(startorder) == 1 {
	//	StartOrderInZapsi(&data, userId, orderId, operationId, workplaceId)
	//} else if len(endorder) == 1 {
	//	LogInfo("MAIN", "Ending order")
	//	data.Message = "Ending order"
	//	if len(nok) > 0 && len(ok) > 0 {
	//		SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
	//		EndOrderInZapsi(orderId, operationId, userId, workplaceId)
	//		SaveNokIntoSyteline(nok, noktype)
	//		SaveOkIntoSyteline(ok)
	//	} else if len(nok) > 0 {
	//		SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
	//		SaveNokIntoSyteline(nok, noktype)
	//	} else if len(ok) > 0 {
	//		SaveOkIntoSyteline(ok)
	//	}
	//	data.Username = "Zadejte prosím své číslo"
	//	data.UsernameValue = ""
	//	data.UserDisabled = ""
	//	data.UserFocus = "autofocus"
	//
	//} else if len(transferorder) == 1 {
	//	LogInfo("MAIN", "Transferring order")
	//	data.Message = "Transferring order"
	//	if len(nok) > 0 && len(ok) > 0 {
	//		LogInfo("MAIN", "Saving both")
	//		SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
	//		SaveNokIntoSyteline(nok, noktype)
	//		SaveOkIntoSyteline(ok)
	//	} else if len(nok) > 0 {
	//		LogInfo("MAIN", "Saving just nok")
	//		SaveNokIntoZapsi(nok, noktype, workplaceId, userId)
	//		SaveNokIntoSyteline(nok, noktype)
	//	} else if len(ok) > 0 {
	//		LogInfo("MAIN", "Saving just ok")
	//		SaveOkIntoSyteline(ok)
	//	}
	//	data.EndOrderButton = "disabled"
	//	data.TransferOrderButton = "disabled"
	//	data.UsernameValue = userId[0]
	//	data.OrderValue = orderId[0]
	//	data.OperationValue = operationId[0]
	//	var nokTypes []SytelineNok
	//	db, err := gorm.Open("mssql", SytelineConnection)
	//
	//	command := "declare @JePlatny ListYesNoType, @Kod ReasonCodeType = NULL exec [rostra_exports_test].dbo.ZapsiKodyDuvoduZmetkuSp @Kod= @Kod, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	//	rows, err := db.Raw(command).Rows()
	//	if err != nil {
	//		LogError("MAIN", "Error: "+err.Error())
	//	}
	//	defer rows.Close()
	//	for rows.Next() {
	//		var nokType SytelineNok
	//		err = rows.Scan(&nokType.Kod, &nokType.Nazev)
	//		nokTypes = append(nokTypes, nokType)
	//		if err != nil {
	//			LogError("MAIN", "Error: "+err.Error())
	//		}
	//	}
	//	db.Close()
	//	data.NokTypes = nokTypes
	//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	//	data.Workplaces = append(data.Workplaces, workplace)
	//	data.Ok = ""
	//	data.Nok = ""
	//	data.OkFocus = "autofocus"
	//	data.OkDisabled = ""
	//	data.NokDisabled = ""
	//}

}

//
//func SaveNokIntoZapsi(nok []string, noktype []string, workplaceId []string, userId []string) {
//	CreateFailInZapsiIfNotExists(noktype)
//	SaveTerminalInputFail(nok, noktype, workplaceId, userId)
//}
//
//func SaveTerminalInputFail(nok []string, noktype []string, workplaceId []string, userId []string) {
//	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
//	var splittedUserName []string
//	if strings.Contains(trimmedUserName, ",") {
//		splittedUserName = strings.Split(trimmedUserName, ",")
//	} else {
//		LogError("MAIN", "Bad username format: "+userId[0])
//		splittedUserName = append(splittedUserName, trimmedUserName)
//		splittedUserName = append(splittedUserName, trimmedUserName)
//	}
//	var zapsiFail Fail
//	var terminalInputFail TerminalInputFail
//	var zapsiWorkplace Workplace
//	var zapsiUser User
//
//	pcs, err := strconv.Atoi(nok[0])
//	if err != nil {
//		LogError("MAIN", "Problem parsing Nok amount when saving terminal input fail")
//		return
//	}
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return
//	}
//	defer db.Close()
//	db.Where("Name = ?", noktype[0]).Find(&zapsiFail)
//	db.Where("Code = ?", workplaceId[0]).Find(&zapsiWorkplace)
//	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
//	terminalInputFail.FailID = zapsiFail.OID
//	terminalInputFail.DeviceID = zapsiWorkplace.DeviceID
//	terminalInputFail.UserID = zapsiUser.OID
//	terminalInputFail.DT = time.Now()
//	for i := 0; i < pcs; i++ {
//		db.Save(&terminalInputFail)
//	}
//	return
//}
//
//func CreateFailInZapsiIfNotExists(noktype []string) {
//	var zapsiFail Fail
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return
//	}
//	defer db.Close()
//	db.Where("Name = ?", noktype[0]).Find(&zapsiFail)
//	if zapsiFail.OID > 0 {
//		LogInfo("MAIN", "Fail "+noktype[0]+" already exists")
//		return
//	}
//	LogInfo("MAIN", "Fail "+noktype[0]+" does not exist, creating fail")
//	zapsiFail.Name = noktype[0]
//	zapsiFail.FailTypeID = 100
//	db.Create(&zapsiFail)
//	var newZapsiFail Fail
//	db.Where("Name = ?", noktype[0]).Find(&newZapsiFail)
//	return
//}
//
//func SaveOkIntoSyteline(ok []string) {
//
//}
//
//func SaveNokIntoSyteline(nok []string, noktype []string) {
//	LogInfo("MAIN", "Saving NOK to Syteline")
//	//quantity, err := strconv.Atoi(nok[0])
//	//if err != nil {
//	//	LogError("MAIN", "Problem parsing count when saving to Syteline")
//	//}
//	db, err := gorm.Open("mssql", SytelineConnection)
//	if err != nil {
//		LogError("MAIN", "Error opening db: "+err.Error())
//		return
//	}
//	defer db.Close()
//	db.Exec("set nocount, ansi_nulls, quoted_identifier, arithabort, xact_abort on \nbegin tran\n" +
//		" \ninsert into zapsi_trans" +
//		"\n( trans_date , emp_num , trans_type , job , suffix , oper_num , wc , qty_complete , qty_scrapped , lot , start_date_time , end_date_time , complete_op , shift , reason_code , time_divisor)" +
//		"\nvalues" +
//		"\n( '20200430' , N' 500001' , 5 , N'3VP0014981' , 0, 10, N'HLIQ12' , 100.0 , 0.0 , '3VP0014981' , NULL , NULL , 0 , NULL , NULL  , NULL)" +
//		"\nif @@error <> 0\nbegin\n  rollback\n  raiserror (N'Chyba zápisu do zapsi_trans', 16, 1)" +
//		"\n \nend\nelse\nbegin\n  commit\nend")
//}
//
//func CheckNok(userId []string, orderId []string, operationId []string, workplaceId []string, ok []string, nok []string, noktype []string) {
//	LogInfo("MAIN", "Checking NOK: ["+ok[0]+":"+nok[0]+"]")
//	sytelineOperation, sytelineWorkplaces := CheckOperationInSyteline(userId, orderId, operationId)
//	mn1 := CheckForMn1(workplaceId, sytelineWorkplaces)
//	mn2 := sytelineOperation.priznak_mn_2 == "1"
//	mn3 := sytelineOperation.priznak_mn_3 == "1"
//	LogInfo("MAIN", "Priznak mn_1: "+strconv.FormatBool(mn1))
//	LogInfo("MAIN", "Priznak mn_2: "+strconv.FormatBool(mn2))
//	LogInfo("MAIN", "Priznak mn_3: "+strconv.FormatBool(mn3))
//	checkedOk, operatorAmountLessThanInZapsi := CheckIfOperatorAmountLessThanInZapsi(nok, userId, orderId, operationId, workplaceId)
//	LogInfo("MAIN", "Operator inserted less amount than in Zapsi: "+strconv.FormatBool(operatorAmountLessThanInZapsi))
//	if checkedOk {
//		if mn1 && !operatorAmountLessThanInZapsi {
//			LogInfo("MAIN", "Mn1 and more amount than in Zapsi, displaying error")
//			//TODO: displayError
//		}
//		if mn2 && !operatorAmountLessThanInZapsi {
//			LogInfo("MAIN", "Mn2 and more amount than in Zapsi, displaying error")
//			//TODO: displayError
//		}
//		if mn3 && !operatorAmountLessThanInZapsi {
//			LogInfo("MAIN", "Mn3 and more amount than in Zapsi, displaying error")
//			//TODO: displayError
//		}
//		anyOrderExists := CheckAnyOrderInZapsi(workplaceId)
//		if anyOrderExists {
//			thisOrderIsOpen := CheckThisOrderInZapsi(userId, orderId, operationId, workplaceId)
//			if thisOrderIsOpen {
//				EnableTransferAndEndButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
//			}
//		} else {
//			sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId)
//			if sytelineOperation.jen_prenos_mnozstvi == "1" {
//
//				EnableTransferButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
//			}
//		}
//	} else {
//		LogError("MAIN", "Problem checking data from Zapsi")
//	}
//}
//
//func CheckOk(userId []string, orderId []string, operationId []string, workplaceId []string, ok []string, nok []string, noktype []string) {
//	LogInfo("MAIN", "Checking OK")
//	sytelineOperation, sytelineWorkplaces := CheckOperationInSyteline(userId, orderId, operationId)
//	mn1 := CheckForMn1(workplaceId, sytelineWorkplaces)
//	mn2 := sytelineOperation.priznak_mn_2 == "1"
//	mn3 := sytelineOperation.priznak_mn_3 == "1"
//	LogInfo("MAIN", "Priznak mn_1: "+strconv.FormatBool(mn1))
//	LogInfo("MAIN", "Priznak mn_2: "+strconv.FormatBool(mn2))
//	LogInfo("MAIN", "Priznak mn_3: "+strconv.FormatBool(mn3))
//	checkedOk, operatorAmountLessThanInZapsi := CheckIfOperatorAmountLessThanInZapsi(ok, userId, orderId, operationId, workplaceId)
//	LogInfo("MAIN", "Operator inserted less amount than in Zapsi: "+strconv.FormatBool(operatorAmountLessThanInZapsi))
//	if checkedOk {
//		if mn1 && !operatorAmountLessThanInZapsi {
//			LogInfo("MAIN", "Mn1 and more amount than in Zapsi, displaying error")
//			//TODO: displayError
//		}
//		if mn2 && !operatorAmountLessThanInZapsi {
//			LogInfo("MAIN", "Mn2 and more amount than in Zapsi, displaying error")
//			//TODO: displayError
//		}
//		if mn3 && !operatorAmountLessThanInZapsi {
//			LogInfo("MAIN", "Mn3 and more amount than in Zapsi, displaying error")
//			//TODO: displayError
//		}
//		anyOrderExists := CheckAnyOrderInZapsi(workplaceId)
//		if anyOrderExists {
//			thisOrderIsOpen := CheckThisOrderInZapsi(userId, orderId, operationId, workplaceId)
//			if thisOrderIsOpen {
//				EnableTransferAndEndButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
//			}
//		} else {
//			sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId)
//			if sytelineOperation.jen_prenos_mnozstvi == "1" {
//				EnableTransferButton(workplaceId, userId, orderId, operationId, data, ok, nok, noktype)
//			}
//		}
//	} else {
//		LogError("MAIN", "Problem checking data from Zapsi")
//	}
//}
//
//func CheckIfOperatorAmountLessThanInZapsi(userAmount []string, userId []string, orderId []string, operationId []string, workplaceId []string) (bool, bool) {
//	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
//	var splittedUserName []string
//	if strings.Contains(trimmedUserName, ",") {
//		splittedUserName = strings.Split(trimmedUserName, ",")
//	} else {
//		LogError("MAIN", "Bad username format: "+userId[0])
//		splittedUserName = append(splittedUserName, trimmedUserName)
//		splittedUserName = append(splittedUserName, trimmedUserName)
//	}
//	order, suffix := ParseOrder(orderId[0])
//	orderName := order + "." + suffix + "-" + operationId[0]
//	var zapsiUser User
//	var zapsiOrder Order
//	var zapsiWorkplace Workplace
//	var terminalInputOrder TerminalInputOrder
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return false, false
//	}
//	defer db.Close()
//	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
//	db.Where("Name = ?", orderName).Find(&zapsiOrder)
//	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
//	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
//	okAmount, err := strconv.Atoi(userAmount[0])
//	if err != nil {
//		LogError("MAIN", "Problem parsing data from user")
//		return false, false
//	}
//	if okAmount < terminalInputOrder.Count {
//		return true, true
//	}
//	return true, false
//}
//
//func CheckForMn1(workplaceId []string, workplaces []SytelineWorkplace) bool {
//	for _, workplace := range workplaces {
//		if workplace.Zapsi_zdroj == workplaceId[0] {
//			if workplace.priznak_mn_1 == "1" {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//func StartOrderInZapsi(data *RostraMainPage, userId []string, orderId []string, operationId []string, workplaceId []string) {
//	LogInfo("MAIN", "Starting order")
//	data.Message = "Starting order"
//	sytelineOrder := CheckOrderInSyteline(userId, orderId)
//	sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId)
//	CreateProductInZapsiIfNotExists(sytelineOrder)
//	zapsiOrder := CreateOrderInZapsiIfNotExists(sytelineOrder, orderId, operationId, sytelineOperation, workplaceId)
//	CreateTerminalOrderInZapsi(userId, zapsiOrder, sytelineOperation, workplaceId)
//	data.Username = "Zadejte prosím své číslo"
//	data.Order = ""
//	data.OrderValue = ""
//	data.Operation = ""
//	data.OperationValue = ""
//	data.UsernameValue = ""
//	data.UserDisabled = ""
//	data.OrderDisabled = "disabled"
//	data.OperationDisabled = "disabled"
//	data.WorkplaceDisabled = "disabled"
//	data.Workplaces = []SytelineWorkplace{}
//	data.UserFocus = "autofocus"
//}
//
//func CheckAnyOrderInZapsi(workplaceId []string) bool {
//	var zapsiWorkplace Workplace
//	var terminalInputOrder TerminalInputOrder
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return false
//	}
//	defer db.Close()
//	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
//	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Find(&terminalInputOrder)
//	if terminalInputOrder.OID > 0 {
//		return true
//	}
//	return false
//
//}
//
//func EndOrderInZapsi(orderId []string, operationId []string, userId []string, workplaceId []string) {
//	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
//	var splittedUserName []string
//	if strings.Contains(trimmedUserName, ",") {
//		splittedUserName = strings.Split(trimmedUserName, ",")
//	} else {
//		LogError("MAIN", "Bad username format: "+userId[0])
//		splittedUserName = append(splittedUserName, trimmedUserName)
//		splittedUserName = append(splittedUserName, trimmedUserName)
//	}
//	order, suffix := ParseOrder(orderId[0])
//	orderName := order + "." + suffix + "-" + operationId[0]
//	var zapsiUser User
//	var zapsiOrder Order
//	var zapsiWorkplace Workplace
//	var terminalInputOrder TerminalInputOrder
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return
//	}
//	defer db.Close()
//	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
//	db.Where("Name = ?", orderName).Find(&zapsiOrder)
//	db.Where("Code = ?", workplaceId[0]).Find(&zapsiWorkplace)
//	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
//	//TODO: Get OK and NOK pcs and  count average cycle
//	if terminalInputOrder.OID > 0 {
//		LogInfo("MAIN", "Closing order "+strconv.Itoa(terminalInputOrder.OID))
//		db.Model(&terminalInputOrder).Where("OID = ?", terminalInputOrder.OID).UpdateColumn(TerminalInputOrder{DTE: sql.NullTime{Time: time.Now(), Valid: true}})
//		db.Model(&terminalInputOrder).Where("OID = ?", terminalInputOrder.OID).UpdateColumn(TerminalInputOrder{Interval: float32(time.Now().Sub(terminalInputOrder.DTS).Seconds())})
//	}
//}
//
//func CreateTerminalOrderInZapsi(userId []string, order Order, operation SytelineOperation, workplaceId []string) {
//	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
//	var splittedUserName []string
//	if strings.Contains(trimmedUserName, ",") {
//		splittedUserName = strings.Split(trimmedUserName, ",")
//	} else {
//		LogError("MAIN", "Bad username format: "+userId[0])
//		splittedUserName = append(splittedUserName, trimmedUserName)
//		splittedUserName = append(splittedUserName, trimmedUserName)
//	}
//	var zapsiUser User
//	var zapsiWorkplace Workplace
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return
//	}
//	defer db.Close()
//	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
//	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
//	var terminalInputOrder TerminalInputOrder
//	defer db.Close()
//	parsedCavity, err := strconv.Atoi(operation.nasobnost)
//	if err != nil {
//		LogError("MAIN", "Problem parsing cavity: "+operation.nasobnost)
//	}
//	terminalInputOrder.DTS = time.Now()
//	terminalInputOrder.OrderID = order.OID
//	terminalInputOrder.UserID = zapsiUser.OID
//	terminalInputOrder.DeviceID = zapsiWorkplace.DeviceID
//	terminalInputOrder.Interval = 0
//	terminalInputOrder.Count = 0
//	terminalInputOrder.Fail = 0
//	terminalInputOrder.AverageCycle = 0.0
//	terminalInputOrder.WorkerCount = 1
//	terminalInputOrder.WorkplaceModeID = 1
//	terminalInputOrder.Cavity = parsedCavity
//	db.Create(&terminalInputOrder)
//}
//
//func CreateProductInZapsiIfNotExists(order SytelineOrder) {
//	var zapsiProduct Product
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return
//	}
//	defer db.Close()
//	db.Where("Name = ?", order.PolozkaVp).Find(&zapsiProduct)
//	if zapsiProduct.OID > 0 {
//		LogInfo("MAIN", "Product "+order.PolozkaVp+" already exists")
//		return
//	}
//	LogInfo("MAIN", "Product "+order.PolozkaVp+" does not exist, creating product")
//	zapsiProduct.Name = order.PolozkaVp
//	zapsiProduct.Barcode = order.PolozkaVp
//	zapsiProduct.Cycle = 1
//	zapsiProduct.IdleFromTime = 1
//	zapsiProduct.ProductGroupID = 1
//	zapsiProduct.ProductStatusID = 1
//	db.Create(&zapsiProduct)
//}
//
//func CreateOrderInZapsiIfNotExists(sytelineOrder SytelineOrder, orderId []string, operationId []string, operation SytelineOperation, workplaceId []string) Order {
//	var zapsiOrder Order
//	var newOrder Order
//	var zapsiProduct Product
//	var zapsiWorkplace Workplace
//	order, suffix := ParseOrder(orderId[0])
//	zapsiOrderName := order + "." + suffix + "-" + operationId[0]
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return zapsiOrder
//	}
//	defer db.Close()
//	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
//	if zapsiOrder.OID > 0 {
//		LogInfo("MAIN", "Order "+zapsiOrder.Name+" already exists")
//		return zapsiOrder
//	}
//	LogInfo("MAIN", "Order "+zapsiOrder.Name+" does not exist, creating order in zapsi")
//	db.Where("Name = ?", sytelineOrder.PolozkaVp).Find(&zapsiProduct)
//	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
//	countRequestedConverted, err := strconv.ParseFloat(operation.mn_2_ks, 32)
//	if err != nil {
//		LogError("MAIN", "Problem parsing count for sytelineOrder: "+operation.mn_2_ks)
//	}
//	newOrder.Name = zapsiOrderName
//	newOrder.Barcode = zapsiOrderName
//	newOrder.ProductID = zapsiProduct.OID
//	newOrder.OrderStatusID = 1
//	newOrder.CountRequested = int(countRequestedConverted)
//	newOrder.WorkplaceID = zapsiWorkplace.OID
//	db.Create(&newOrder)
//	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
//	return zapsiOrder
//}
//
//func CheckThisOrderInZapsi(userId []string, orderId []string, operationid []string, workplaceId []string) bool {
//	trimmedUserName := strings.ReplaceAll(userId[0], " ", "")
//	var splittedUserName []string
//	if strings.Contains(trimmedUserName, ",") {
//		splittedUserName = strings.Split(trimmedUserName, ",")
//	} else {
//		LogError("MAIN", "Bad username format: "+userId[0])
//		splittedUserName = append(splittedUserName, trimmedUserName)
//		splittedUserName = append(splittedUserName, trimmedUserName)
//	}
//	order, suffix := ParseOrder(orderId[0])
//	orderName := order + "." + suffix + "-" + operationid[0]
//	var zapsiUser User
//	var zapsiOrder Order
//	var zapsiWorkplace Workplace
//	var terminalInputOrder TerminalInputOrder
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return false
//	}
//	defer db.Close()
//	db.Where("Name = ?", splittedUserName[0]).Where("FirstName = ?", splittedUserName[1]).Find(&zapsiUser)
//	db.Where("Name = ?", orderName).Find(&zapsiOrder)
//	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
//	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
//	if terminalInputOrder.OID > 0 {
//		return true
//	}
//	return false
//}
//

func CheckInputStep(orderId []string, operationId []string, workplaceId []string, startorder []string, transferorder []string, endorder []string, ok []string, nok []string) interface{} {
	if len(startorder) == 1 {
		return startOrderStep
	} else if len(transferorder) == 1 {
		return transferOrderStep
	} else if len(endorder) == 1 {
		return endOrderStep
	} else if orderId[0] == "" && operationId[0] == "" && workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
		return checkUserStep
	} else if operationId[0] == "" && workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
		return checkOrderStep
	} else if workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
		return checkOperationStep
	} else if ok[0] == "" && nok[0] == "" {
		return checkWorkplaceStep
	}
	return checkAmountStep
}

func RostraMainScreen(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Displaying main screen")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	_ = tmpl.Execute(writer, data)
}

func CreateDefaultPage() RostraMainPage {
	data := RostraMainPage{
		Version:             "version: " + version,
		Username:            "Zadejte prosím své číslo",
		UsernameValue:       "",
		Order:               "",
		OrderValue:          "",
		Operation:           "",
		OperationValue:      "",
		Workplace:           "",
		Ok:                  "",
		Nok:                 "",
		OkValue:             "",
		NokValue:            "",
		OrderDisabled:       "disabled",
		OperationDisabled:   "disabled",
		WorkplaceDisabled:   "disabled",
		OkDisabled:          "disabled",
		NokDisabled:         "disabled",
		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
		UserFocus:           "autofocus",
		RadioDisabled:       "disabled",
	}
	if len(data.NokTypes) == 0 {
		nokType := SytelineNok{Kod: "", Nazev: ""}
		data.NokTypes = append(data.NokTypes, nokType)
	}
	if len(data.Workplaces) == 0 {
		workplace := SytelineWorkplace{Zapsi_zdroj: "", priznak_mn_1: "", vice_vp: "", SL_prac: "", typ_zdroje_zapsi: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
		data.Workplaces = append(data.Workplaces, workplace)
	}
	return data
}
