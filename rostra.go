package main

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

type RostraMainPage struct {
	Version        string
	Username       string
	Order          string
	Operation      string
	WorkplaceGroup string
	Workplace      string

	UsernameValue       string
	OrderValue          string
	OperationValue      string
	WorkplaceGroupValue string
	WorkplaceValue      string

	UserDisabled           string
	OrderDisabled          string
	OperationDisabled      string
	WorkplaceGroupDisabled string
	WorkplaceDisabled      string

	UserFocus           string
	OrderFocus          string
	OperationFocus      string
	WorkplaceGroupFocus string
	WorkplaceFocus      string

	StartOrderButton    string
	EndOrderButton      string
	TransferOrderButton string
}

const (
	checkUser int = iota
	checkOrder
	checkOperation
	checkWorkplaceGroup
	checkWorkplace
)

func DataInput(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Checking data input")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	_ = r.ParseForm()
	userId := r.Form["userid"]
	orderId := r.Form["orderid"]
	operationId := r.Form["operationid"]
	workplaceGroupId := r.Form["workplacegroupid"]
	workplaceId := r.Form["workplaceid"]
	LogInfo("MAIN", "user: "+userId[0])
	LogInfo("MAIN", "order: "+orderId[0])
	LogInfo("MAIN", "operation: "+operationId[0])
	LogInfo("MAIN", "workplacegroup: "+workplaceGroupId[0])
	LogInfo("MAIN", "workplace: "+workplaceId[0])
	data := RostraMainPage{
		Version:                "version: " + version,
		Username:               "Zadejte prosím své číslo",
		Order:                  "Zadejte prosím číslo zakázky",
		Operation:              "Zadejte prosím číslo operace",
		WorkplaceGroup:         "Zadejte prosím číslo skupiny pracoviště",
		Workplace:              "Zadejte prosím číslo pracoviště",
		UserDisabled:           "disabled",
		OrderDisabled:          "disabled",
		OperationDisabled:      "disabled",
		WorkplaceGroupDisabled: "disabled",
		WorkplaceDisabled:      "disabled",
		StartOrderButton:       "disabled",
		EndOrderButton:         "disabled",
		TransferOrderButton:    "disabled",
	}
	inputStep := CheckInputStep(userId, orderId, operationId, workplaceGroupId, workplaceId)
	switch inputStep {
	case checkUser:
		CheckUserInSyteline(userId, &data)
	case checkOrder:
		CheckOrderInSyteline(userId, orderId, &data)
	case checkOperation:
		{
			//TODO: Check Operation
			LogInfo("MAIN", "Checking operation")
			data.Operation = "Operace"
			data.OperationValue = "Operace"
			data.UsernameValue = userId[0]
			data.OrderValue = orderId[0]
			data.WorkplaceGroupDisabled = ""
		}
	case checkWorkplaceGroup:
		{
			//TODO: Check WorkplaceGroup
			LogInfo("MAIN", "Checking workplacegroup")
			data.WorkplaceGroup = "Skupina"
			data.WorkplaceGroupValue = "Skupina"
			data.UsernameValue = userId[0]
			data.OrderValue = orderId[0]
			data.OperationValue = operationId[0]
			data.WorkplaceDisabled = ""
		}
	case checkWorkplace:
		{
			//TODO: Check Workplace
			LogInfo("MAIN", "Checking workplace")
			data.Workplace = "Pracoviste"
			data.WorkplaceValue = "Pracoviste"
			data.UsernameValue = userId[0]
			data.OrderValue = orderId[0]
			data.OperationValue = operationId[0]
			data.WorkplaceGroupValue = workplaceGroupId[0]
			data.StartOrderButton = ""
		}
	}

	//TODO: Check start-stop-update
	//if inputStep == checkWorkplace {
	//	LogInfo("MAIN", "Checking workplace")
	//
	//	//db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	//	//if err != nil {
	//	//	println("Error opening db: " + err.Error())
	//	//}
	//	//defer db.Close()
	//	//var jePlatny string
	//	//command := "declare   @JePlatny   ListYesNoType, @CisloVP    JobType, @PriponaVP  SuffixType, @Operace    OperNumType select   @CisloVP = N'"+workplaceId[0]+"', @PriponaVP = 0, @Operace = 10 exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace , @JePlatny = @JePlatny output select JePlatny = @JePlatny"
	//	//row := db.Raw(command).Row()
	//	//err = row.Scan(&jePlatny)
	//
	//	jePlatny := "1"
	//	if jePlatny == "1" {
	//		data.Workplace = workplaceId[0]
	//		data.WorkplaceValue = workplaceId[0]
	//		data.Order = orderId[0]
	//		data.OrderValue = orderId[0]
	//		data.UsernameValue = userId[0]
	//		data.WorkplaceDisabled = ""
	//	} else {
	//		data.Username = "Číslo nenalezeno, zadejte prosím znovu"
	//		LogInfo("MAIN", "Workplace not found for "+workplaceId[0])
	//
	//	}
	//
	//}
	_ = tmpl.Execute(writer, data)
}

func CheckOrderInSyteline(orderId []string, userId []string, data *RostraMainPage) {
	LogInfo("MAIN", "Checking order")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
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
		data.Order = "Číslo nenalezeno, zadejte prosím znovu"
		data.OrderDisabled = ""
	}
	//TODO: CheckOrderInZapsi(orderId)
}

func CheckUserInSyteline(userId []string, data *RostraMainPage) {
	//data.Username = "Jahoda"
	//data.UsernameValue = "Jahoda"
	//data.OrderDisabled = ""
	//data.OrderFocus = "autofocus"
	LogInfo("MAIN", "Checking user")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
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
	} else {
		LogInfo("MAIN", "User not found for "+userId[0])
		data.Username = "Číslo nenalezeno, zadejte prosím znovu"
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
	}
	//TODO: CheckUserInZapsi(userId, jmeno)
}

func CheckInputStep(userId []string, orderId []string, operationId []string, workplaceGroupId []string, workplaceId []string) interface{} {
	if orderId[0] == "" && operationId[0] == "" && workplaceGroupId[0] == "" && workplaceId[0] == "" {
		return checkUser
	} else if operationId[0] == "" && workplaceGroupId[0] == "" && workplaceId[0] == "" {
		return checkOrder
	} else if workplaceGroupId[0] == "" && workplaceId[0] == "" {
		return checkOperation
	} else if workplaceId[0] == "" {
		return checkWorkplaceGroup
	}
	return checkWorkplace

}

func OrderInput(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Checking Order")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	_ = r.ParseForm()
	userId := r.Form["userid"]
	println(len(userId))
	orderId := r.Form["order"]
	println(len(orderId))
	data := RostraMainPage{
		Version:             "version: " + version,
		Username:            userId[0],
		Order:               "Zadejte prosím číslo zakázky",
		Workplace:           "Zadejte prosím číslo pracoviště",
		UserDisabled:        "disabled",
		OrderDisabled:       "disabled",
		WorkplaceDisabled:   "disabled",
		StartOrderButton:    "disabled",
		EndOrderButton:      " disabled",
		TransferOrderButton: "disabled",
	}
	println(4)
	if len(orderId) > 0 {
		LogInfo("MAIN", "Downloading order for user")
		//db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
		//if err != nil {
		//	println("Error opening db: " + err.Error())
		//}
		//defer db.Close()
		//var jePlatny string
		//command := "declare   @JePlatny ListYesNoType, @VP Infobar = N'" + orderId[0] + "' exec [rostra_exports_test].dbo.ZapsiKontrolaVPSp @VP= @VP, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
		//row := db.Raw(command).Row()
		//err = row.Scan(&jePlatny)
		jePlatny := "1"
		if jePlatny == "1" {
			data.Order = orderId[0]
			data.WorkplaceDisabled = ""
		} else {
			data.Username = "Číslo nenalezeno, zadejte prosím znovu"
		}
	}
	_ = tmpl.Execute(writer, data)
}

func RostraMainScreen(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Displaying main screen")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := RostraMainPage{
		Version: "version: " + version,

		Username:            "Zadejte prosím své číslo",
		UsernameValue:       "",
		Order:               "Zadejte prosím číslo zakázky",
		OrderValue:          "",
		Operation:           "Zadejte prosím číslo operace",
		OperationValue:      "",
		WorkplaceGroup:      "Zadejte prosím číslo skupiny pracoviště",
		WorkplaceGroupValue: "",
		Workplace:           "Zadejte prosím číslo pracoviště",
		WorkplaceValue:      "",

		OrderDisabled:          "disabled",
		OperationDisabled:      "disabled",
		WorkplaceGroupDisabled: "disabled",
		WorkplaceDisabled:      "disabled",

		UserFocus: "autofocus",

		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
	}

	_ = tmpl.Execute(writer, data)
}
