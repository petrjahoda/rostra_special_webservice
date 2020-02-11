package main

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

type RostraMainPage struct {
	Version             string
	Username            string
	Order               string
	Workplace           string
	UsernameValue       string
	OrderValue          string
	WorkplaceValue      string
	OrderDisabled       string
	UserDisabled        string
	WorkplaceDisabled   string
	StartOrderButton    string
	EndOrderButton      string
	TransferOrderButton string
}

const (
	checkUser int = iota
	checkOrder
	checkWorkplace
)

func DataInput(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Checking data input")
	tmpl := template.Must(template.ParseFiles("html/rostra_main_screen.html"))
	_ = r.ParseForm()
	userId := r.Form["userid"]
	orderId := r.Form["order"]
	workplaceId := r.Form["workplace"]
	data := RostraMainPage{
		Version:             "version: " + version,
		Username:            "Zadejte prosím své číslo",
		Order:               "Zadejte prosím číslo zakázky",
		Workplace:           "Zadejte prosím číslo pracoviště",
		UserDisabled:        "",
		OrderDisabled:       "disabled",
		WorkplaceDisabled:   "disabled",
		StartOrderButton:    "disabled",
		EndOrderButton:      " disabled",
		TransferOrderButton: "disabled",
	}
	inputStep := CheckInputStep(userId, orderId, workplaceId)

	if inputStep == checkUser {
		LogInfo("MAIN", "Checking user")
		db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
		if err != nil {
			println("Error opening db: " + err.Error())
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
			data.UserDisabled = "disabled"
		} else {
			data.Username = "Číslo nenalezeno, zadejte prosím znovu"
			LogInfo("MAIN", "User not found for "+userId[0])
		}
	}
	if inputStep == checkOrder {
		LogInfo("MAIN", "Checking order")
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
			data.OrderValue = orderId[0]
			data.UsernameValue = userId[0]
			data.WorkplaceDisabled = ""
		} else {
			data.Order = "Číslo nenalezeno, zadejte prosím znovu"
			LogInfo("MAIN", "Order not found for "+orderId[0])
		}

	}
	if inputStep == checkWorkplace {
		LogInfo("MAIN", "Checking workplace")

		//db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
		//if err != nil {
		//	println("Error opening db: " + err.Error())
		//}
		//defer db.Close()
		//var jePlatny string
		//command := "declare   @JePlatny   ListYesNoType, @CisloVP    JobType, @PriponaVP  SuffixType, @Operace    OperNumType select   @CisloVP = N'"+workplaceId[0]+"', @PriponaVP = 0, @Operace = 10 exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace , @JePlatny = @JePlatny output select JePlatny = @JePlatny"
		//row := db.Raw(command).Row()
		//err = row.Scan(&jePlatny)

		jePlatny := "1"
		if jePlatny == "1" {
			data.Workplace = workplaceId[0]
			data.WorkplaceValue = workplaceId[0]
			data.Order = orderId[0]
			data.OrderValue = orderId[0]
			data.UsernameValue = userId[0]
			data.WorkplaceDisabled = ""
		} else {
			data.Username = "Číslo nenalezeno, zadejte prosím znovu"
			LogInfo("MAIN", "Workplace not found for "+workplaceId[0])

		}

	}
	_ = tmpl.Execute(writer, data)
}

func CheckInputStep(userId []string, orderId []string, workplaceId []string) interface{} {
	if orderId[0] == "" && workplaceId[0] == "" {
		return checkUser
	} else if userId[0] == "" && workplaceId[0] == "" {
		return checkOrder
	}
	return checkWorkplace

}

func OrderInput(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Checking Order")
	tmpl := template.Must(template.ParseFiles("html/rostra_main_screen.html"))
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
	tmpl := template.Must(template.ParseFiles("html/rostra_main_screen.html"))
	data := RostraMainPage{
		Version:             "version: " + version,
		Username:            "Zadejte prosím své číslo",
		UsernameValue:       "",
		Order:               "Zadejte prosím číslo zakázky",
		OrderValue:          "",
		Workplace:           "Zadejte prosím číslo pracoviště",
		WorkplaceValue:      "",
		OrderDisabled:       "disabled",
		WorkplaceDisabled:   "disabled",
		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
	}

	_ = tmpl.Execute(writer, data)
}
