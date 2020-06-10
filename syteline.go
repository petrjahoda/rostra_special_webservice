package main

import (
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func EndOrderInSyteline(userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) bool {
	// TODO: complete
	return false
}

func TransferOrderInSyteline(userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) bool {
	// TODO: complete
	return false
}

func StartOrderInSyteline(userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) bool {
	// TODO: complete
	return false
}

func GetNokTypesFromSyteline() []SytelineNok {
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
	return nokTypes
}
func GetWorkplaceFromSyteline(orderid []string, operationid []string, workplaceid []string) SytelineWorkplace {
	LogInfo("MAIN", "Checking operation")
	order, suffix := ParseOrder(orderid[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOperation SytelineOperation
	var sytelineWorkplace SytelineWorkplace
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		return sytelineWorkplace
	}
	defer db.Close()
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationid[0] + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
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
		LogInfo("MAIN", "Operation found: "+operationid[0])
		command = "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationid[0] + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
		workplaceRows, err := db.Raw(command).Rows()
		if err != nil {
			LogError("MAIN", "Error: "+err.Error())
		}
		defer workplaceRows.Close()
		for workplaceRows.Next() {
			var sytelineWorkplaceScanned SytelineWorkplace
			err = workplaceRows.Scan(&sytelineWorkplaceScanned.Zapsi_zdroj, &sytelineWorkplaceScanned.priznak_mn_1, &sytelineWorkplaceScanned.vice_vp, &sytelineWorkplaceScanned.SL_prac, &sytelineWorkplaceScanned.typ_zdroje_zapsi, &sytelineWorkplaceScanned.auto_prevod_mnozstvi, &sytelineWorkplaceScanned.mnozstvi_auto_prevodu)
			if err != nil {
				LogError("MAIN", "Error: "+err.Error())
			}
			if sytelineWorkplaceScanned.Zapsi_zdroj == workplaceid[0] {
				sytelineWorkplace = sytelineWorkplaceScanned
			}
		}
	}
	return sytelineWorkplace
}

func GetOperationFromSyteline(orderid []string, operationid []string) SytelineOperation {
	LogInfo("MAIN", "Getting operation from syteline")
	order, suffix := ParseOrder(orderid[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOperation SytelineOperation
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		return sytelineOperation
	}
	defer db.Close()
	command := "declare @JePlatny ListYesNoType, @CisloVP JobType, @PriponaVP  SuffixType, @Operace OperNumType select @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationid[0] + " exec [rostra_exports_test].dbo.ZapsiKontrolaOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP, @Operace = @Operace, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
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
	return sytelineOperation
}

func GetOrderFromSyteline(orderid []string) SytelineOrder {
	LogInfo("MAIN", "Getting order from Syteline")
	order, suffix := ParseOrder(orderid[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineOrder SytelineOrder
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
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
	return sytelineOrder
}

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
			CreateUserInZapsiIfNotExists(sytelineUser, userId)
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

		LogInfo("MAIN", "Sending page for user check")
		_ = tmpl.Execute(*writer, data)
	}
}
