package main

import (
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strings"
)

func CheckOrderInSyteline(writer *http.ResponseWriter, userId []string, orderId []string) {
	LogInfo("MAIN", "Checking order")
	order, suffix := ParseOrder(orderId[0])
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOrder SytelineOrder
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.Order = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.Message = err.Error()
		data.OrderFocus = "autofocus"
		data.OrderDisabled = ""
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
			data.Order = orderId[0]
			data.OrderValue = orderId[0]
			data.UsernameValue = userId[0]
			data.OperationFocus = "autofocus"
			data.OperationDisabled = ""
		} else {
			LogInfo("MAIN", "Order not found for "+orderId[0]+" for command "+command)
			data.UsernameValue = userId[0]
			data.Order = "Číslo nenalezeno, nebo je neplatné, zadejte prosím znovu"
			data.OrderDisabled = ""
			data.OrderFocus = "autofocus"
		}
		LogInfo("MAIN", "Sending page for order check")
		_ = tmpl.Execute(*writer, data)
	}
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
			data.Username = userId[0] + ";" + sytelineUser.Jmeno
			data.UsernameValue = userId[0] + ";" + sytelineUser.Jmeno
			data.OrderDisabled = ""
			data.OrderFocus = "autofocus"
			data.Order = "Zadejte prosím číslo zakázky"
		} else {
			LogInfo("MAIN", userId[0]+", user not found")
			data.Username = "Číslo nenalezeno, zadejte prosím znovu"
			data.UserDisabled = ""
			data.UserFocus = "autofocus"
		}
		CreateUserInZapsiIfNotExists(sytelineUser, userId)
		LogInfo("MAIN", "Sending page for user check")
		_ = tmpl.Execute(*writer, data)
	}
}

func CreateUserInZapsiIfNotExists(user SytelineUser, userId []string) {
	trimmedUserName := strings.ReplaceAll(user.Jmeno, " ", "")
	var userName []string
	if strings.Contains(trimmedUserName, ",") {
		userName = strings.Split(trimmedUserName, ",")
	} else {
		LogError("MAIN", "Bad username format: "+user.Jmeno)
		userName = append(userName, trimmedUserName)
		userName = append(userName, trimmedUserName)
	}
	var zapsiUser User
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name LIKE ?", userName[0]).Where("FirstName LIKE ?", userName[1]).Find(&zapsiUser)
	if zapsiUser.OID > 0 {
		LogInfo("MAIN", "User "+user.Jmeno+"already exists")
		return
	}
	LogInfo("MAIN", "User "+user.Jmeno+" does not exist, creating user "+user.Jmeno)
	zapsiUser.Login = userId[0]
	zapsiUser.Name = userName[0]
	zapsiUser.FirstName = userName[1]
	zapsiUser.UserRoleID = "1"
	zapsiUser.UserTypeID = "1"
	db.Create(&zapsiUser)
}
