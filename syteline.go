package main

import (
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func CheckOperationInSyteline(writer *http.ResponseWriter, userId []string, orderId []string, operationId []string) {
	LogInfo("MAIN", "Checking operation")
	order, suffix := ParseOrder(orderId[0])
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
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
		data.Username = "disabled"
		data.Message = err.Error()
		LogInfo("MAIN", "Sending error page for order check")
		_ = tmpl.Execute(*writer, data)
	} else {
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
			command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationId[0] + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
			LogInfo("MAIN", command)
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
			}
			if len(sytelineWorkplaces) > 0 {
				data.Workplaces = sytelineWorkplaces
				LogInfo("MAIN", "Workplaces found: "+strconv.Itoa(len(sytelineWorkplaces)))
				data.UsernameValue = userId[0]
				data.OrderValue = orderId[0]
				data.Operation = operationId[0]
				data.OperationValue = operationId[0]
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
	// DEBUG
	//data.Order = orderId[0]
	//data.OrderValue = orderId[0]
	//data.UsernameValue = userId[0]
	//data.OperationFocus = "autofocus"
	//data.OperationDisabled = ""
	//data.UserDisabled = "disabled"
	//_ = tmpl.Execute(*writer, data)

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
			data.Order = orderId[0]
			data.OrderValue = orderId[0]
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

	// DEBUG
	//data.Username = "1234;Petr Jahoda"
	//data.UsernameValue = "1234;Petr Jahoda"
	//data.UserDisabled = "disabled"
	//data.OrderDisabled = ""
	//data.OrderFocus = "autofocus"
	//data.Order = "Zadejte prosím číslo zakázky"
	//_ = tmpl.Execute(*writer, data)

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
