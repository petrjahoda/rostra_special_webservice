package main

import (
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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
	LogInfo("MAIN", "[user]     : "+userId[0])
	LogInfo("MAIN", "[order]    : "+orderId[0])
	LogInfo("MAIN", "[operation]: "+operationId[0])
	LogInfo("MAIN", "[workplace]: "+workplaceId[0])
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
		orderExists := CheckOrderInZapsi(userId, orderId, operationId, workplaceId)
		if orderExists {
			data.EndOrderButton = ""
			data.TransferOrderButton = ""
			data.UsernameValue = userId[0]
			data.OrderValue = orderId[0]
			data.OperationValue = operationId[0]
			workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
			data.Workplaces = append(data.Workplaces, workplace)
		} else {
			data.StartOrderButton = ""
			data.UsernameValue = userId[0]
			data.OrderValue = orderId[0]
			data.OperationValue = operationId[0]
			workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
			data.Workplaces = append(data.Workplaces, workplace)
		}
	}
	if len(data.Workplaces) == 0 {
		LogInfo("MAIN", "No workplaces, adding null workplace")
		workplace := SytelineWorkplace{Zapsi_zdroj: "", priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
		data.Workplaces = append(data.Workplaces, workplace)
	}
	_ = tmpl.Execute(writer, data)
}

func CheckOrderInZapsi(userId []string, orderId []string, operationid []string, workplaceId []string) bool {
	//TODO: check order here
	return false
}

func CheckOperationInSyteline(userId []string, orderId []string, operationId []string, data *RostraMainPage) (SytelineOperation, []SytelineWorkplace) {
	order, suffix := ParseOrder(orderId[0])
	LogInfo("MAIN", "Checking operation")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
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
		data.Operation = operationId[0]
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.OperationValue = operationId[0]
		data.WorkplaceDisabled = ""
		LogInfo("MAIN", "Operation found: "+data.Operation)
	} else {
		LogInfo("MAIN", "Operation not found for "+orderId[0])
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

func AddWorkplaces(workplaces []SytelineWorkplace) string {
	var workplacesData string
	for _, workplace := range workplaces {
		workplacesData += workplace.Zapsi_zdroj + ","
	}
	return workplacesData
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
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
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
		data.Order = orderId[0]
		data.OrderValue = orderId[0]
		data.UsernameValue = userId[0]
		data.OperationFocus = "autofocus"
		data.OperationDisabled = ""
		LogInfo("MAIN", "Order found: "+data.OrderValue)
	} else {
		LogInfo("MAIN", "Order not found for "+orderId[0]+" for command "+command)
		data.UsernameValue = userId[0]
		data.Order = "Číslo nenalezeno, nebo je neplatné, zadejte prosím znovu"
		data.OrderDisabled = ""
		data.OrderFocus = "autofocus"
	}
	return sytelineOrder
}

func CheckUserInSyteline(userId []string, data *RostraMainPage) SytelineUser {
	LogInfo("MAIN", "Checking user")
	db, err := gorm.Open("mssql", "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test")
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
		data.Username = sytelineUser.Jmeno
		data.UsernameValue = sytelineUser.Jmeno
		data.OrderDisabled = ""
		data.OrderFocus = "autofocus"
		data.Order = "Zadejte prosím číslo zakázky"
		LogInfo("MAIN", "User found: "+data.UsernameValue)
	} else {
		LogInfo("MAIN", "User not found for "+userId[0])
		data.Username = "Číslo nenalezeno, zadejte prosím znovu"
		data.UserDisabled = ""
		data.UserFocus = "autofocus"
	}
	//TODO: CheckUserInZapsi(userId, jmeno)
	return sytelineUser
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
		Version:             "version: " + version,
		Username:            "Zadejte prosím své číslo",
		UsernameValue:       "",
		Order:               "",
		OrderValue:          "",
		Operation:           "",
		OperationValue:      "",
		Workplace:           "",
		OrderDisabled:       "disabled",
		OperationDisabled:   "disabled",
		WorkplaceDisabled:   "disabled",
		UserFocus:           "autofocus",
		StartOrderButton:    "disabled",
		EndOrderButton:      "disabled",
		TransferOrderButton: "disabled",
	}
	if len(data.Workplaces) == 0 {
		LogInfo("MAIN", "No workplaces, adding null workplace")
		workplace := SytelineWorkplace{Zapsi_zdroj: "", priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
		data.Workplaces = append(data.Workplaces, workplace)
	}
	_ = tmpl.Execute(writer, data)
}
