package main

//
//import (
//	"github.com/julienschmidt/httprouter"
//	"html/template"
//	"net/http"
//	"strconv"
//	"strings"
//)
//
//type DisplayOrder struct {
//	OrderCode                   string
//	OrderName                   string
//	ProductName                 string
//	WorkplaceName               string
//	OrderStart                  string
//	OrderRequestedTotal         string
//	OrderCountTotal             string
//	OrderCountActual            string
//	OrderSendToSytelineTotal    string
//	OrderSendToSytelineActual   string
//	ForSave                     string
//	OrderSendToSytelineNokTotal string
//}
//type RostraMainPage struct {
//	Version   string
//	Username  string
//	Order     string
//	Operation string
//	Workplaces string
//	Ok        string
//	Nok       string
//	Message   string
//
//	UsernameValue  string
//	OrderValue     string
//	OperationValue string
//	OkValue        string
//	NokValue       string
//
//	StartOrderButton    string
//	EndOrderButton      string
//	TransferOrderButton string
//
//	UserFocus      string
//	OrderFocus     string
//	OperationFocus string
//	WorkplaceFocus string
//	OkFocus        string
//	NokFocus       string
//
//	UserDisabled      string
//	OrderDisabled     string
//	OperationDisabled string
//	WorkplaceDisabled string
//	OkDisabled        string
//	NokDisabled       string
//	RadioDisabled     string
//	ClovekDisabled    string
//	StrojDisabled     string
//	SerizeniDisabled  string
//
//	NokTypes     []SytelineNok
//	Workplaces   []SytelineWorkplace
//	DisplayOrder []DisplayOrder
//}
//
//const (
//	checkUserStep int = iota
//	checkOrderStep
//	checkOperationStep
//	checkWorkplaceStep
//	checkAmountStep
//	startOrderStep
//	transferOrderStep
//	endOrderStep
//)
//
//func DataInput(writer http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	logInfo("MAIN", "Checking data input")
//	_ = r.ParseForm()
//	userid := r.Form["userid"]
//	orderid := r.Form["orderid"]
//	operationid := r.Form["operationid"]
//	workplaceid := r.Form["workplaceid"]
//	startorder := r.Form["startorder"]
//	endorder := r.Form["endorder"]
//	transferorder := r.Form["transferorder"]
//	noktype := r.Form["noktype"]
//	nok := r.Form["nok"]
//	ok := r.Form["ok"]
//	radio := r.Form["syteline"]
//	logInfo("MAIN", "[USER:"+userid[0]+"] [ORDER:"+orderid[0]+"] [OPERATION:"+operationid[0]+"] [WORKPLACE:"+workplaceid[0]+"]")
//	logInfo("MAIN", "[OK:NOK:TYPE]  ["+ok[0]+":"+nok[0]+":"+noktype[0]+"]")
//	logInfo("MAIN", "[RADIO:"+strconv.Itoa(len(radio))+"]")
//	logInfo("MAIN", "[START:TRANSFER:END] ["+strconv.Itoa(len(startorder))+":"+strconv.Itoa(len(transferorder))+":"+strconv.Itoa(len(endorder))+"]")
//	inputStep := CheckInputStep(orderid, operationid, workplaceid, startorder, transferorder, endorder, ok, nok)
//	switch inputStep {
//	case checkUserStep:
//		CheckUserInSyteline(&writer, userid)
//	case checkOrderStep:
//		CheckOrderInSyteline(&writer, userid, orderid)
//	case checkOperationStep:
//		CheckOperationInSyteline(&writer, userid, orderid, operationid)
//	case checkWorkplaceStep:
//		FirstControls(&writer, workplaceid, userid, orderid, operationid)
//	case checkAmountStep:
//		SecondControls(&writer, workplaceid, userid, orderid, operationid, ok, nok, noktype)
//	case startOrderStep:
//		StartOrderButton(&writer, userid, orderid, operationid, workplaceid, radio)
//	case transferOrderStep:
//		TransferOrderButton(&writer, userid, orderid, operationid, workplaceid, ok, nok, noktype)
//	case endOrderStep:
//		EndOrderButton(&writer, userid, orderid, operationid, workplaceid, ok, nok, noktype, radio)
//	}
//}
//
//func EndOrderButton(writer *http.ResponseWriter, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, radio []string) {
//	logInfo("MAIN", "Ending order")
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//
//	actualTimeDivisor := GetActualTimeDivisor(workplaceIdSplitted)
//	actualZapsiOpenorders := GetActualZapsiOpenFor(workplaceIdSplitted)
//	sytelineOrderEnded := EndOrderInSyteline(userid, orderid, operationid, workplaceIdSplitted, ok, nok, noktype, radio, actualTimeDivisor)
//	if actualZapsiOpenorders == 1 {
//		actualTimeDivisor = 1
//	}
//	UpdateDeviceWithNew(actualTimeDivisor, workplaceIdSplitted)
//	zapsiOrderEnded := EndOrderInZapsi(userid, orderid, operationid, workplaceIdSplitted, ok, nok)
//	SaveNokIntoZapsi(nok, noktype, workplaceIdSplitted, userid)
//	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
//	data := CreateDefaultPage()
//	logInfo("MAIN", "Order in Zapsi closed "+strconv.FormatBool(zapsiOrderEnded))
//	data.Message += "Order in zapsi closed: " + strconv.FormatBool(zapsiOrderEnded) + "\n"
//	logInfo("MAIN", "Order in Syteline closed "+strconv.FormatBool(sytelineOrderEnded))
//	data.Message += "Order in syteline closed: " + strconv.FormatBool(sytelineOrderEnded) + "\n"
//	_ = tmpl.Execute(*writer, data)
//}
//
//func TransferOrderButton(writer *http.ResponseWriter, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string) {
//	logInfo("MAIN", "Transferring order")
//	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
//	data := CreateDefaultPage()
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	thisOpenOrderInZapsi, _ := CheckThisOpenOrderInZapsi(userid, orderid, operationid, workplaceIdSplitted)
//	if !thisOpenOrderInZapsi {
//		zapsiOrderCreated := StartAndCloseOrderInZapsi(userid, orderid, operationid, workplaceIdSplitted, ok, nok)
//		logInfo("MAIN", "Order in Zapsi transfered "+strconv.FormatBool(zapsiOrderCreated))
//		data.Message += "Order in zapsi transfered: " + strconv.FormatBool(zapsiOrderCreated) + "\n"
//
//	}
//	sytelineOkAndNokTransferred := TransferOkAndNokToSyteline(userid, orderid, operationid, workplaceIdSplitted, ok, nok, noktype)
//	SaveNokIntoZapsi(nok, noktype, workplaceIdSplitted, userid)
//	logInfo("MAIN", "Ok and NOK to Syteline transfered "+strconv.FormatBool(sytelineOkAndNokTransferred))
//	data.Message += "Ok and NOK to Syteline transfered: " + strconv.FormatBool(sytelineOkAndNokTransferred) + "\n"
//	_ = tmpl.Execute(*writer, data)
//}
//
//func StartOrderButton(writer *http.ResponseWriter, userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) {
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	logInfo("MAIN", "Starting order for "+workplaceIdSplitted[0]+" and "+radio[0])
//	zapsiOrderCreated := StartOrderInZapsi(userid, orderid, operationid, workplaceIdSplitted, radio)
//	actualTimeDivisor := GetActualTimeDivisor(workplaceIdSplitted)
//	actualZapsiOpenorders := GetActualZapsiOpenFor(workplaceIdSplitted)
//	if actualZapsiOpenorders > actualTimeDivisor {
//		UpdateDeviceWithNew(actualZapsiOpenorders, workplaceIdSplitted)
//	}
//	sytelineOrderCreated := StartOrderInSyteline(userid, orderid, operationid, workplaceIdSplitted, radio)
//	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
//	data := CreateDefaultPage()
//	logInfo("MAIN", "Order in Zapsi created "+strconv.FormatBool(zapsiOrderCreated))
//	data.Message += "Order in zapsi created: " + strconv.FormatBool(zapsiOrderCreated) + "\n"
//	logInfo("MAIN", "Order in Syteline created "+strconv.FormatBool(sytelineOrderCreated))
//	data.Message += "Order in syteline created: " + strconv.FormatBool(sytelineOrderCreated) + "\n"
//	_ = tmpl.Execute(*writer, data)
//}
//
//func CheckInputStep(orderId []string, operationId []string, workplaceId []string, startorder []string, transferorder []string, endorder []string, ok []string, nok []string) interface{} {
//	if len(startorder) == 1 {
//		logInfo("MAIN", "Start order step")
//		return startOrderStep
//	} else if len(transferorder) == 1 {
//		logInfo("MAIN", "Transfer order step")
//		return transferOrderStep
//	} else if len(endorder) == 1 {
//		logInfo("MAIN", "End order step")
//		return endOrderStep
//	} else if orderId[0] == "" && operationId[0] == "" && workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
//		logInfo("MAIN", "Check user step")
//		return checkUserStep
//	} else if operationId[0] == "" && workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
//		logInfo("MAIN", "Check order step")
//		return checkOrderStep
//	} else if workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
//		logInfo("MAIN", "Check operation step")
//		return checkOperationStep
//	} else if ok[0] == "" && nok[0] == "" {
//		logInfo("MAIN", "Check workplace step")
//		return checkWorkplaceStep
//	}
//	logInfo("MAIN", "Check amount step")
//	return checkAmountStep
//}
//
//func RostraMainScreen(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
//	logInfo("MAIN", "Displaying main screen")
//	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
//	data := CreateDefaultPage()
//	_ = tmpl.Execute(writer, data)
//}
//
//func CreateDefaultPage() RostraMainPage {
//	data := RostraMainPage{
//		Version:             "version: " + version,
//		Username:            "Zadejte prosím své číslo",
//		UsernameValue:       "",
//		Order:               "",
//		OrderValue:          "",
//		Operation:           "",
//		OperationValue:      "",
//		Workplaces:           "",
//		Ok:                  "",
//		Nok:                 "",
//		OkValue:             "",
//		NokValue:            "",
//		OrderDisabled:       "disabled",
//		OperationDisabled:   "disabled",
//		WorkplaceDisabled:   "disabled",
//		OkDisabled:          "disabled",
//		NokDisabled:         "disabled",
//		StartOrderButton:    "disabled",
//		EndOrderButton:      "disabled",
//		TransferOrderButton: "disabled",
//		UserFocus:           "autofocus",
//		RadioDisabled:       "disabled",
//		ClovekDisabled:      "disabled",
//		SerizeniDisabled:    "disabled",
//		StrojDisabled:       "disabled",
//	}
//	if len(data.NokTypes) == 0 {
//		nokType := SytelineNok{Kod: "", Nazev: ""}
//		data.NokTypes = append(data.NokTypes, nokType)
//	}
//	if len(data.Workplaces) == 0 {
//		workplace := SytelineWorkplace{Zapsi_zdroj: "", priznak_mn_1: "", Vice_vp: "", SL_prac: "", typ_zdroje_zapsi: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//		data.Workplaces = append(data.Workplaces, workplace)
//	}
//	if len(data.DisplayOrder) == 0 {
//		displayOrder := DisplayOrder{
//			OrderCode:                   "",
//			OrderName:                   "",
//			ProductName:                 "",
//			WorkplaceName:               "",
//			OrderStart:                  "",
//			OrderRequestedTotal:         "",
//			OrderCountTotal:             "",
//			OrderCountActual:            "",
//			OrderSendToSytelineTotal:    "",
//			OrderSendToSytelineActual:   "",
//			ForSave:                     "",
//			OrderSendToSytelineNokTotal: "",
//		}
//		data.DisplayOrder = append(data.DisplayOrder, displayOrder)
//	}
//	return data
//}
