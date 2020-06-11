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
	ClovekDisabled    string
	StrojDisabled     string
	SerizeniDisabled  string

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
		SecondControls(&writer, workplaceid, userid, orderid, operationid, ok, nok, noktype)
	case startOrderStep:
		StartOrderButton(&writer, userid, orderid, operationid, workplaceid, radio)
	case transferOrderStep:
		TransferOrderButton(&writer, userid, orderid, operationid, workplaceid, ok, nok, noktype, radio)
	case endOrderStep:
		EndOrderButton(&writer, userid, orderid, operationid, workplaceid, ok, nok, noktype, radio)
	}
}

func EndOrderButton(writer *http.ResponseWriter, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, radio []string) {
	LogInfo("MAIN", "Ending order")
	sytelineOrderEnded := EndOrderInSyteline(userid, orderid, operationid, workplaceid, ok, nok, noktype, radio)
	zapsiOrderEnded := EndOrderInZapsi(userid, orderid, operationid, workplaceid, ok, nok)
	SaveNokIntoZapsi(nok, noktype, workplaceid, userid)
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	LogInfo("MAIN", "Order in Zapsi closed "+strconv.FormatBool(zapsiOrderEnded))
	data.Message += "Order in zapsi closed: " + strconv.FormatBool(zapsiOrderEnded) + "\n"
	LogInfo("MAIN", "Order in Syteline closed "+strconv.FormatBool(sytelineOrderEnded))
	data.Message += "Order in syteline closed: " + strconv.FormatBool(sytelineOrderEnded) + "\n"
	_ = tmpl.Execute(*writer, data)
}

func TransferOrderButton(writer *http.ResponseWriter, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, radio []string) {
	LogInfo("MAIN", "Transferring order")
	zapsiOrderCreated := StartAndCloseOrderInZapsi(userid, orderid, operationid, workplaceid, ok, nok)
	sytelineOkAndNokTransferred := TransferOkAndNokToSyteline(userid, orderid, operationid, workplaceid, ok, nok, noktype)
	SaveNokIntoZapsi(nok, noktype, workplaceid, userid)
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	LogInfo("MAIN", "Order in Zapsi transfered "+strconv.FormatBool(zapsiOrderCreated))
	data.Message += "Order in zapsi transfered: " + strconv.FormatBool(zapsiOrderCreated) + "\n"
	LogInfo("MAIN", "Ok and NOK to Syteline transfered "+strconv.FormatBool(sytelineOkAndNokTransferred))
	data.Message += "Ok and NOK to Syteline transfered: " + strconv.FormatBool(sytelineOkAndNokTransferred) + "\n"
	_ = tmpl.Execute(*writer, data)
}

func StartOrderButton(writer *http.ResponseWriter, userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) {
	LogInfo("MAIN", "Starting order")
	zapsiOrderCreated := StartOrderInZapsi(userid, orderid, operationid, workplaceid)
	sytelineOrderCreated := StartOrderInSyteline(userid, orderid, operationid, workplaceid, radio)
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	LogInfo("MAIN", "Order in Zapsi created "+strconv.FormatBool(zapsiOrderCreated))
	data.Message += "Order in zapsi created: " + strconv.FormatBool(zapsiOrderCreated) + "\n"
	LogInfo("MAIN", "Order in Syteline created "+strconv.FormatBool(sytelineOrderCreated))
	data.Message += "Order in syteline created: " + strconv.FormatBool(sytelineOrderCreated) + "\n"
	_ = tmpl.Execute(*writer, data)
}

func CheckInputStep(orderId []string, operationId []string, workplaceId []string, startorder []string, transferorder []string, endorder []string, ok []string, nok []string) interface{} {
	if len(startorder) == 1 {
		LogInfo("MAIN", "Start order step")
		return startOrderStep
	} else if len(transferorder) == 1 {
		LogInfo("MAIN", "Transfer order step")
		return transferOrderStep
	} else if len(endorder) == 1 {
		LogInfo("MAIN", "End order step")
		return endOrderStep
	} else if orderId[0] == "" && operationId[0] == "" && workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
		LogInfo("MAIN", "Check user step")
		return checkUserStep
	} else if operationId[0] == "" && workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
		LogInfo("MAIN", "Check order step")
		return checkOrderStep
	} else if workplaceId[0] == "" && ok[0] == "" && nok[0] == "" {
		LogInfo("MAIN", "Check operation step")
		return checkOperationStep
	} else if ok[0] == "" && nok[0] == "" {
		LogInfo("MAIN", "Check workplace step")
		return checkWorkplaceStep
	}
	LogInfo("MAIN", "Check amount step")
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
		ClovekDisabled:      "disabled",
		SerizeniDisabled:    "disabled",
		StrojDisabled:       "disabled",
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
