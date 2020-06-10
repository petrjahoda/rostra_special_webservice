package main

import (
	"html/template"
	"net/http"
	"strconv"
)

func SecondControls(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string, ok []string, nok []string, noktype []string) {
	LogInfo("MAIN", "Starting second controls")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	sytelineOperation := GetOperationFromSyteline(orderid, operationid)
	sytelineWorkplace := GetWorkplaceFromSyteline(orderid, operationid, workplaceid)
	if sytelineOperation.jen_prenos_mnozstvi == "1" {
		LogInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is one, only transfer will be available")
		data.Message += "jen_prenos_mnozstvi je 1\n"
		okCheck := CheckAmount(ok, sytelineWorkplace, data, userid, orderid, operationid, workplaceid, sytelineOperation)
		nokCheck := CheckAmount(nok, sytelineWorkplace, data, userid, orderid, operationid, workplaceid, sytelineOperation)
		if okCheck == false || nokCheck == false {
			data.Message += "uzivatel zadal vice kusu, nez je odvedeno v zapsi\n"
			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
		} else {
			EnableTransfer(writer, workplaceid, userid, orderid, operationid, data, ok, nok, tmpl)
		}
	} else {
		LogInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is not one, transfer and close will be available")
		data.Message += "jen_prenos_mnozstvi neni 1\n"
		okCheck := CheckAmount(ok, sytelineWorkplace, data, userid, orderid, operationid, workplaceid, sytelineOperation)
		nokCheck := CheckAmount(nok, sytelineWorkplace, data, userid, orderid, operationid, workplaceid, sytelineOperation)
		if okCheck == false || nokCheck == false {
			data.Message += "uzivatel zadal vice kusu, nez je odvedeno v zapsi\n"
			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
		} else {
			if sytelineWorkplace.typ_zdroje_zapsi == "0" {
				LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
				data.Message += "typ_zdroje_zapsi je 0\n"
				EnableClovekTransferCloseInput(writer, data, userid, orderid, operationid, workplaceid, ok, nok, noktype, tmpl)

			} else {
				LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
				data.Message += "typ_zdroje_zapsi neni 0\n"
				EnableClovekSerizeniStrojTransferCloseInput(writer, data, userid, orderid, operationid, workplaceid, ok, nok, noktype, tmpl)
			}
		}
	}
}

func EnableClovekSerizeniStrojTransferCloseInput(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, tmpl *template.Template) {
	var workplaces []SytelineWorkplace
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	workplaces = append(workplaces, workplace)
	var nokTypes []SytelineNok
	nokType := SytelineNok{Nazev: noktype[0]}
	nokTypes = append(nokTypes, nokType)
	data.Workplaces = workplaces
	data.Workplaces = workplaces
	data.UsernameValue = userid[0]
	data.OrderValue = orderid[0]
	data.Operation = operationid[0]
	data.OperationValue = operationid[0]
	data.UserDisabled = "disabled"
	data.RadioDisabled = ""
	data.ClovekDisabled = "checked"
	data.SerizeniDisabled = ""
	data.StrojDisabled = ""
	data.TransferOrderButton = ""
	data.EndOrderButton = ""
	data.OkValue = ok[0]
	data.NokValue = nok[0]
	_ = tmpl.Execute(*writer, data)
}

func EnableClovekTransferCloseInput(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, tmpl *template.Template) {
	var workplaces []SytelineWorkplace
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	workplaces = append(workplaces, workplace)
	var nokTypes []SytelineNok
	nokType := SytelineNok{Nazev: noktype[0]}
	nokTypes = append(nokTypes, nokType)
	data.Workplaces = workplaces
	data.NokTypes = nokTypes
	data.UsernameValue = userid[0]
	data.OrderValue = orderid[0]
	data.Operation = operationid[0]
	data.OperationValue = operationid[0]
	data.UserDisabled = "disabled"
	data.RadioDisabled = ""
	data.ClovekDisabled = "checked"
	data.TransferOrderButton = ""
	data.OkValue = ok[0]
	data.NokValue = nok[0]
	_ = tmpl.Execute(*writer, data)
}

func CheckAmount(inputAmount []string, sytelineWorkplace SytelineWorkplace, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, sytelineOperation SytelineOperation) bool {
	amountCheck := true
	checkAmount, err := strconv.Atoi(inputAmount[0])
	if err != nil {
		LogError("MAIN", "Problem parsing data from user")
		return false
	}
	LogInfo("MAIN", "Checking amount: "+strconv.Itoa(checkAmount))
	if sytelineWorkplace.priznak_mn_1 == "1" {
		LogInfo("MAIN", "sytelineWorkplace.priznak_mn_1 is one")
		data.Message += "priznak_mn_1 je 1\n"
		amountLessThanInZapsi := CheckIfOperatorAmountLessThanInZapsi(inputAmount, userid, orderid, operationid, workplaceid)
		if amountLessThanInZapsi {
			LogInfo("MAIN", "ok from user is less than in zapsi")
			data.Message += "uzivatel zadal mene kusu nez je v zapsi, coz je spravne\n"
		} else {
			LogInfo("MAIN", "ok from user is more than in zapsi")
			data.Message += "uzivatel zadal vice kusu nez je v zapsi, coz je spatne\n"
			amountCheck = false
		}
	} else {
		LogInfo("MAIN", "sytelineWorkplace.priznak_mn_1 is not one")
		data.Message += "priznak_mn_1 neni 1\n"
	}
	if sytelineOperation.priznak_mn_2 == "1" {
		LogInfo("MAIN", "sytelineWorkplace.priznak_mn_2 is one")
		data.Message += "priznak_mn_2 je 1\n"
		parsedFromSyteline, err := strconv.ParseFloat(sytelineOperation.mn_2_ks, 64)
		if err != nil {
			LogError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_2_ks+", "+err.Error())
		}
		parsedFromInput, err := strconv.Atoi(inputAmount[0])
		if err != nil {
			LogError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_2_ks+", "+err.Error())
		}
		if parsedFromInput < int(parsedFromSyteline) {
			LogInfo("MAIN", "ok from user is less than mn_2_ks")
			data.Message += "uzivatel zadal mene kusu nez je v mn_2_ks, coz je spravne\n"
		} else {
			LogInfo("MAIN", "ok from user is more than mn_2_ks")
			data.Message += "uzivatel zadal vic kusu nez je v mn_2_ks, coz je spatne\n"
			amountCheck = false
		}
	} else {
		LogInfo("MAIN", "sytelineWorkplace.priznak_mn_2 is not one")
		data.Message += "priznak_mn_2 neni 1\n"
	}
	if sytelineOperation.priznak_mn_3 == "1" {
		LogInfo("MAIN", "sytelineWorkplace.priznak_mn_3 is one")
		data.Message += "priznak_mn_3 je 1\n"
		parsedFromSyteline, err := strconv.ParseFloat(sytelineOperation.mn_3_ks, 64)
		if err != nil {
			LogError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_3_ks+", "+err.Error())
		}
		parsedFromInput, err := strconv.Atoi(inputAmount[0])
		if err != nil {
			LogError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_3_ks+", "+err.Error())
		}
		if parsedFromInput < int(parsedFromSyteline) {
			LogInfo("MAIN", "ok from user is less than mn_3_ks")
			data.Message += "uzivatel zadal mene kusu nez je v mn_3_ks, coz je spravne\n"
		} else {
			LogInfo("MAIN", "ok from user is more than mn_3_ks")
			data.Message += "uzivatel zadal vic kusu nez je v mn_3_ks, coz je spatne\n"
			amountCheck = false
		}

	} else {
		LogInfo("MAIN", "sytelineWorkplace.priznak_mn_3 is not one")
		data.Message += "priznak_mn_3 neni 1\n"
	}
	return amountCheck
}

func FirstControls(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string) {
	LogInfo("MAIN", "Starting first controls")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	sytelineOperation := GetOperationFromSyteline(orderid, operationid)
	sytelineWorkplace := GetWorkplaceFromSyteline(orderid, operationid, workplaceid)
	anyOpenOrderInZapsi := CheckAnyOpenOrderInZapsi(workplaceid)
	if anyOpenOrderInZapsi {
		LogInfo("MAIN", workplaceid[0]+" has open order in Zapsi")
		data.Message += "V Zapsi existuje pro toto pracoviste otevrena zakazka\n"
		thisOpenOrderInZapsi := CheckThisOpenOrderInZapsi(userid, orderid, operationid, workplaceid)
		if thisOpenOrderInZapsi {
			LogInfo("MAIN", workplaceid[0]+" has open this exact order in Zapsi")
			data.Message += "V Zapsi existuje pro toto pracoviste otevrena zakazka\n"
			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
		} else {
			LogInfo("MAIN", workplaceid[0]+" has not open this exact order in Zapsi")
			data.Message += "V Zapsi NEexistuje pro toto pracoviste presne tato otevrena zakazka\n"
			if sytelineWorkplace.vice_vp == "1" {
				LogInfo("MAIN", "sytelineOperationSource.vice_vp equals one")
				data.Message += "vice_vp je 1\n"
				if sytelineWorkplace.typ_zdroje_zapsi == "0" {
					LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
					data.Message += "typ_zdroje_zapsi je 0\n"
					EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
				} else {
					LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
					data.Message += "typ_zdroje_zapsi neni 0\n"
					EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
				}
			} else {
				LogInfo("MAIN", "sytelineOperationSource.vice_vp does not equal one")
				data.Message += "vice_vp neni 1\n"
				if sytelineOperation.parovy_dil == "1" {
					LogInfo("MAIN", "sytelineOperation.parovy_dil equals one")
					data.Message += "parovy dil je 1\n"
					if len(sytelineOperation.seznamm_par_dilu) > 0 {
						LogInfo("MAIN", "sytelineOperation.seznamm_par_dilu not empty")
						data.Message += "seznamm_par_dilu osahuje nejaky parovy dil\n"
						var zapsiProducts = CheckProductsInZapsi(sytelineOperation)
						anyOpenOrderHasOneOfProducts := CheckIfAnyOpenOrderHasOneOfProducts(workplaceid, zapsiProducts)
						if anyOpenOrderHasOneOfProducts {
							LogInfo("MAIN", "Some open order in Zapsi contains any of pair parts")
							data.Message += "Nejaka otevrena zakazka v Zapsi obsahuje nektery z parovych dilu\n"
							if sytelineWorkplace.typ_zdroje_zapsi == "0" {
								LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is zero")
								data.Message += "typ_zdroje_zapsi je 0\n"
								EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
							} else {
								LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is not zero")
								data.Message += "typ_zdroje_zapsi neni 0\n"
								EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
							}
						} else {
							LogInfo("MAIN", "No open order in Zapsi contains any of pair parts")
							data.Message += "Zadna otevrena zakazka v Zapsi neobsahuje zadny z parovych dilu\n"
							CheckOperationInSyteline(writer, userid, orderid, operationid)
						}
					} else {
						LogInfo("MAIN", "sytelineOperation.seznamm_par_dilu is empty")
						data.Message += "seznamm_par_dilu neosahuje zadny parovy dil\n"
						CheckOperationInSyteline(writer, userid, orderid, operationid)
					}
				} else {
					LogInfo("MAIN", "sytelineOperation.parovy_dil does not equal one")
					data.Message += "parovy dil neni 1\n"
					CheckOperationInSyteline(writer, userid, orderid, operationid)
				}
			}
		}
	} else {
		LogInfo("MAIN", workplaceid[0]+" does not have any open order in Zapsi")
		data.Message += "V Zapsi neexistuje pro toto pracoviste otevrena zakazka\n"
		if sytelineOperation.jen_prenos_mnozstvi == "1" {
			LogInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is one")
			data.Message += "jen_prenos_mnozstvi je 1\n"
			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
		} else {
			LogInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is not one")
			data.Message += "jen_prenos_mnozstvi neni 1\n"
			if sytelineWorkplace.typ_zdroje_zapsi == "0" {
				LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is zero")
				data.Message += "typ_zdroje_zapsi je 0\n"
				EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
			} else {
				LogInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is not zero")
				data.Message += "typ_zdroje_zapsi neni 0\n"
				EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
			}
		}
	}
}

func EnableClovekSerizeniStrojStart(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, tmpl *template.Template) {
	var workplaces []SytelineWorkplace
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	workplaces = append(workplaces, workplace)
	data.Workplaces = workplaces
	data.UsernameValue = userid[0]
	data.OrderValue = orderid[0]
	data.Operation = operationid[0]
	data.OperationValue = operationid[0]
	data.UserDisabled = "disabled"
	data.RadioDisabled = ""
	data.ClovekDisabled = "checked"
	data.SerizeniDisabled = ""
	data.StrojDisabled = ""
	data.StartOrderButton = ""
	_ = tmpl.Execute(*writer, data)
}

func EnableClovekStart(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, tmpl *template.Template) {
	var workplaces []SytelineWorkplace
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	workplaces = append(workplaces, workplace)
	data.Workplaces = workplaces
	data.UsernameValue = userid[0]
	data.OrderValue = orderid[0]
	data.Operation = operationid[0]
	data.OperationValue = operationid[0]
	data.UserDisabled = "disabled"
	data.RadioDisabled = ""
	data.ClovekDisabled = "checked"
	data.StartOrderButton = ""
	_ = tmpl.Execute(*writer, data)
}

func EnableOkNok(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string, data RostraMainPage, tmpl *template.Template) {
	nokTypes := GetNokTypesFromSyteline()
	var workplaces []SytelineWorkplace
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	workplaces = append(workplaces, workplace)
	data.Workplaces = workplaces
	data.UsernameValue = userid[0]
	data.OrderValue = orderid[0]
	data.Operation = operationid[0]
	data.OperationValue = operationid[0]
	data.NokTypes = nokTypes
	data.UserDisabled = "disabled"
	data.OkDisabled = ""
	data.NokDisabled = ""
	data.OkValue = "0"
	data.NokValue = "0"
	data.OkFocus = "autofocus"
	_ = tmpl.Execute(*writer, data)
}

func EnableTransfer(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string, data RostraMainPage, ok []string, nok []string, tmpl *template.Template) {
	nokTypes := GetNokTypesFromSyteline()
	var workplaces []SytelineWorkplace
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	workplaces = append(workplaces, workplace)
	data.Workplaces = workplaces
	data.UsernameValue = userid[0]
	data.OrderValue = orderid[0]
	data.Operation = operationid[0]
	data.OperationValue = operationid[0]
	data.NokTypes = nokTypes
	data.UserDisabled = "disabled"
	if len(ok) > 0 {
		data.OkValue = ok[0]
	}
	if len(nok) > 0 {
		data.NokValue = nok[0]
	}
	data.TransferOrderButton = ""
	_ = tmpl.Execute(*writer, data)
}
