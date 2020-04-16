package main

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

type RostraMainPage struct {
	Version   string
	Username  string
	Order     string
	Operation string
	Workplace string

	UsernameValue  string
	OrderValue     string
	OperationValue string
	WorkplaceValue string

	UserDisabled      string
	OrderDisabled     string
	OperationDisabled string
	WorkplaceDisabled string

	UserFocus      string
	OrderFocus     string
	OperationFocus string
	WorkplaceFocus string

	StartOrderButton    string
	EndOrderButton      string
	TransferOrderButton string
}

const (
	checkUser int = iota
	checkOrder
	checkOperation
	checkWorkplace
)

func DataInput(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Checking data input")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	_ = r.ParseForm()
	userId := r.Form["userid"]
	orderId := r.Form["orderid"]
	operationId := r.Form["operationid"]
	workplaceId := r.Form["workplaceid"]
	LogInfo("MAIN", "user: "+userId[0])
	LogInfo("MAIN", "order: "+orderId[0])
	LogInfo("MAIN", "operation: "+operationId[0])
	LogInfo("MAIN", "workplace: "+workplaceId[0])
	data := RostraMainPage{
		Version:             "version: " + version,
		Username:            "Zadejte prosím své číslo",
		Order:               "",
		Operation:           "",
		Workplace:           "",
		UserDisabled:        "disabled",
		OrderDisabled:       "disabled",
		OperationDisabled:   "disabled",
		WorkplaceDisabled:   "disabled",
		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
	}
	inputStep := CheckInputStep(orderId, operationId, workplaceId)
	switch inputStep {
	case checkUser:
		CheckUserInSyteline(userId, &data)
	case checkOrder:
		CheckOrderInSyteline(userId, orderId, &data)
	case checkOperation:
		CheckOperationInSyteline(userId, orderId, operationId, &data)
	case checkWorkplace:
		//GetWorkplacesFromSyteline(userId, orderId, operationId, &data)
	}

	//TODO: Show start-stop-update buttons
	_ = tmpl.Execute(writer, data)
}

func GetWorkplacesFromSyteline(userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	//LogInfo("MAIN", "Checking workplaces")
	//db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	//if err != nil {
	//	LogError("MAIN", "Error opening db: "+err.Error())
	//	data.UsernameValue = userId[0]
	//	data.OrderValue = orderId[0]
	//	data.OperationValue = operationId[0]
	//	data.Workplace = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
	//	data.WorkplaceDisabled = ""
	//	return
	//}
	//defer db.Close()
	//var jePlatny string
	//command := "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select @CisloVP = N'" + orderId[0] + "', @PriponaVP = 0, @Operace = " + operationId[0] + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
	//row := db.Raw(command).Row()
	//err = row.Scan(&jePlatny)
	//if jePlatny == "1" {
	//	data.Workplace =
	//	data.UsernameValue = userId[0]
	//	data.OrderValue = orderId[0]
	//	data.OperationValue = operationId[0]
	//	data.WorkplaceValue =
	//	data.WorkplaceDisabled = ""
	//} else {
	//	LogInfo("MAIN", "Workplaces not found for "+orderId[0])
	//	data.UsernameValue = userId[0]
	//	data.OrderValue = orderId[0]
	//	data.OperationValue = operationId[0]
	//	data.Workplace = "Pracoviště nenalezeny"
	//	data.OperationDisabled = ""
	//}
}

func CheckOperationInSyteline(userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	LogInfo("MAIN", "Checking operation")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.OperationDisabled = ""
		return
	}
	defer db.Close()
	var jePlatny string
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + orderId[0] + "', @PriponaVP = 0, @Operace = " + operationId[0] + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	row := db.Raw(command).Row()
	err = row.Scan(&jePlatny)
	if jePlatny == "1" {
		data.Operation = operationId[0]
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.OperationValue = operationId[0]
		data.OperationDisabled = ""
	} else {
		LogInfo("MAIN", "Operation not found for "+orderId[0])
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Číslo nenalezeno, zadejte prosím znovu"
		data.OperationDisabled = ""
	}
}

func CheckOrderInSyteline(userId []string, orderId []string, data *RostraMainPage) {
	LogInfo("MAIN", "Checking order")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.Order = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.OrderDisabled = ""
		return
	}
	defer db.Close()
	var jePlatny string
	command := "declare   @JePlatny ListYesNoType, @VP Infobar = N'" + orderId[0] + "' exec [rostra_exports_test].dbo.ZapsiKontrolaVPSp @VP= @VP, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	row := db.Raw(command).Row()
	err = row.Scan(&jePlatny)
	if jePlatny == "1" {
		data.Order = orderId[0]
		data.OrderValue = orderId[0]
		data.UsernameValue = userId[0]
		data.OperationDisabled = ""
	} else {
		LogInfo("MAIN", "Order not found for "+orderId[0])
		data.UsernameValue = userId[0]
		data.Order = "Číslo nenalezeno, nebo je neplatné, zadejte prosím znovu"
		data.OrderDisabled = ""
	}
}

func CheckUserInSyteline(userId []string, data *RostraMainPage) {
	LogInfo("MAIN", "Checking user")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.Username = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
		return
	}
	defer db.Close()
	var jePlatny string
	var jmeno string
	command := "declare   @Zamestnanec EmpNumType, @JePlatny ListYesNoType, @Jmeno NameType Exec [rostra_exports_test].dbo.ZapsiKontrolaZamSp @Zamestnanec = N'" + userId[0] + "', @JePlatny = @JePlatny output, @Jmeno = @Jmeno output select JePlatny = @JePlatny, Jmeno = @Jmeno"
	row := db.Raw(command).Row()
	err = row.Scan(&jePlatny, &jmeno)
	if jePlatny == "1" {
		data.Username = jmeno
		data.UsernameValue = jmeno
		data.OrderDisabled = ""
		data.OrderFocus = "autofocus"
		data.Order = "Zadejte prosím číslo zakázky"
	} else {
		LogInfo("MAIN", "User not found for "+userId[0])
		data.Username = "Číslo nenalezeno, zadejte prosím znovu"
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
	}
	//TODO: CheckUserInZapsi(userId, jmeno)
}

func CheckInputStep(orderId []string, operationId []string, workplaceId []string) interface{} {
	if orderId[0] == "" && operationId[0] == "" && workplaceId[0] == "" {
		return checkUser
	} else if operationId[0] == "" && workplaceId[0] == "" {
		return checkOrder
	} else if workplaceId[0] == "" {
		return checkOperation
	}
	return checkWorkplace

}

func RostraMainScreen(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Displaying main screen")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := RostraMainPage{
		Version: "version: " + version,

		Username:       "Zadejte prosím své číslo",
		UsernameValue:  "",
		Order:          "",
		OrderValue:     "",
		Operation:      "",
		OperationValue: "",
		Workplace:      "",
		WorkplaceValue: "",

		OrderDisabled:     "disabled",
		OperationDisabled: "disabled",
		WorkplaceDisabled: "disabled",

		UserFocus: "autofocus",

		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
	}

	_ = tmpl.Execute(writer, data)
}
