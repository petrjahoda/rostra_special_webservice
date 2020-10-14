package main

//
//import (
//	"github.com/jinzhu/gorm"
//	_ "github.com/jinzhu/gorm"
//	"html/template"
//	"net/http"
//	"strconv"
//	"strings"
//)
//
//func SecondControls(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string, ok []string, nok []string, noktype []string) {
//	logInfo("MAIN", "Starting second controls")
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
//	data := CreateDefaultPage()
//	sytelineOperation := GetOperationFromSyteline(orderid, operationid)
//	sytelineWorkplace := GetWorkplaceFromSyteline(orderid, operationid, workplaceIdSplitted)
//	if sytelineOperation.jen_prenos_mnozstvi == "1" {
//		logInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is one, only transfer will be available")
//		data.Message += "jen_prenos_mnozstvi je 1\n"
//		countFromZapsi := GetCountForActualOpenOrder(workplaceIdSplitted, userid, orderid, operationid)
//		countFromSyteline := GetCountForAllTransferredToSyteline(workplaceIdSplitted, userid, orderid, operationid)
//		countFromUser := GetCountFromUser(ok, nok)
//		countCheck := CheckAmount(ok, sytelineWorkplace, data, userid, orderid, operationid, workplaceid, sytelineOperation)
//		logInfo("MAIN", "[CountZapsi:CountSyteline:CountUser] ["+strconv.Itoa(countFromZapsi)+":"+strconv.Itoa(countFromSyteline)+":"+strconv.Itoa(countFromUser)+"]")
//		if countCheck {
//			EnableTransfer(writer, workplaceid, userid, orderid, operationid, data, ok, nok, noktype, tmpl)
//		} else {
//			data.Message += "V Zapsi je vyrobeno " + strconv.Itoa(countFromZapsi) + " kusu, do Syteline uz je odvedeno " + strconv.Itoa(countFromSyteline) + " kusu, je mozno odvest maximalne " + strconv.Itoa(countFromZapsi-countFromSyteline) + " kusu\n"
//			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
//		}
//
//	} else {
//		logInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is not one, transfer and close will be available")
//		data.Message += "jen_prenos_mnozstvi neni 1\n"
//		countFromZapsi := GetCountForActualOpenOrder(workplaceIdSplitted, userid, orderid, operationid)
//		countFromSyteline := GetCountForAllTransferredToSyteline(workplaceIdSplitted, userid, orderid, operationid)
//		countFromUser := GetCountFromUser(ok, nok)
//		countCheck := CheckAmount(ok, sytelineWorkplace, data, userid, orderid, operationid, workplaceid, sytelineOperation)
//		logInfo("MAIN", "[CountZapsi:CountSyteline:CountUser] ["+strconv.Itoa(countFromZapsi)+":"+strconv.Itoa(countFromSyteline)+":"+strconv.Itoa(countFromUser)+"]")
//		if countCheck {
//			//if countFromUser > (countFromZapsi - countFromSyteline) {
//			//	data.Message += "V Zapsi je vyrobeno " + strconv.Itoa(countFromZapsi) + " kusu, do Syteline uz je odvedeno " + strconv.Itoa(countFromSyteline) + " kusu, je mozno odvest maximalne " + strconv.Itoa(countFromZapsi-countFromSyteline) + " kusu\n"
//			//	EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
//			//} else {
//			if sytelineWorkplace.typ_zdroje_zapsi == "0" {
//				logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
//				data.Message += "typ_zdroje_zapsi je 0\n"
//				if sytelineWorkplace.priznak_mn_1 == "0" || (sytelineWorkplace.priznak_mn_1 == "1" && countFromUser == (countFromZapsi-countFromSyteline)) {
//					logInfo("MAIN", "sytelineWorkplace.priznak_mn_1 equals zero or sytelineWorkplace.priznak_mn_1 equals one with the same amount")
//					data.Message += "priznak_mn_1 je 0 anebo je 1 se stejnym mnozstvim\n"
//					EnableClovekTransferCloseInput(writer, data, userid, orderid, operationid, workplaceid, ok, nok, noktype, tmpl)
//				} else {
//					logInfo("MAIN", "sytelineWorkplace.priznak_mn_1 equals zero or sytelineWorkplace.priznak_mn_1 equals one with the same amount")
//					data.Message += "priznak_mn_1 je 1 s ruznym mnozstvim\n"
//					EnableClovekTransferInput(writer, data, userid, orderid, operationid, workplaceid, ok, nok, noktype, tmpl)
//				}
//
//			} else {
//				logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
//				data.Message += "typ_zdroje_zapsi neni 0\n"
//				if sytelineWorkplace.priznak_mn_1 == "0" || (sytelineWorkplace.priznak_mn_1 == "1" && countFromUser == (countFromZapsi-countFromSyteline)) {
//					logInfo("MAIN", "sytelineWorkplace.priznak_mn_1 equals zero or sytelineWorkplace.priznak_mn_1 equals one with the same amount")
//					data.Message += "priznak_mn_1 je 0 anebo je 1 se stejnym mnozstvim\n"
//					EnableClovekSerizeniStrojTransferCloseInput(writer, data, userid, orderid, operationid, workplaceid, ok, nok, noktype, tmpl)
//				} else {
//					logInfo("MAIN", "sytelineWorkplace.priznak_mn_1 equals zero or sytelineWorkplace.priznak_mn_1 equals one with the same amount")
//					data.Message += "priznak_mn_1 je 1 s ruznym mnozstvim\n"
//					EnableClovekSerizeniStrojTransferInput(writer, data, userid, orderid, operationid, workplaceid, ok, nok, noktype, tmpl)
//				}
//			}
//			//}
//		} else {
//			data.Message += "V Zapsi je vyrobeno " + strconv.Itoa(countFromZapsi) + " kusu, do Syteline uz je odvedeno " + strconv.Itoa(countFromSyteline) + " kusu, je mozno odvest maximalne " + strconv.Itoa(countFromZapsi-countFromSyteline) + " kusu\n"
//			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
//		}
//	}
//}
//
//func EnableClovekTransferCloseInput(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	var nokTypes []SytelineNok
//	nokType := SytelineNok{Nazev: noktype[0]}
//	nokTypes = append(nokTypes, nokType)
//	data.NokTypes = nokTypes
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.UserDisabled = "disabled"
//	data.RadioDisabled = ""
//	data.ClovekDisabled = "checked"
//	data.TransferOrderButton = ""
//	data.EndOrderButton = ""
//	data.OkValue = ok[0]
//	data.NokValue = nok[0]
//	data.DisplayOrder = GetActualDataForUser(userid)
//	_ = tmpl.Execute(*writer, data)
//}
//
//func GetCountFromUser(ok []string, nok []string) int {
//	okAsNumber, err := strconv.Atoi(ok[0])
//	if err != nil {
//		logError("MAIN", "Problem parsing number from: "+ok[0])
//		return 0
//	}
//	nokAsNumber, err := strconv.Atoi(nok[0])
//	if err != nil {
//		logError("MAIN", "Problem parsing number from: "+nok[0])
//		return 0
//	}
//	return okAsNumber + nokAsNumber
//}
//
//func EnableClovekSerizeniStrojTransferInput(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	var nokTypes []SytelineNok
//	nokType := SytelineNok{Nazev: noktype[0]}
//	nokTypes = append(nokTypes, nokType)
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	_, orderNote := CheckOpenOrderForWorkplaceInZapsi(orderid, operationid, workplaceIdSplitted)
//	println("Order note: " + orderNote)
//	if orderNote == "clovek" {
//		data.ClovekDisabled = "checked"
//		data.SerizeniDisabled = "disabled"
//		data.StrojDisabled = "disabled"
//	} else if orderNote == "serizeni" {
//		data.ClovekDisabled = "disabled"
//		data.SerizeniDisabled = "checked"
//		data.StrojDisabled = "disabled"
//	} else if orderNote == "stroj" {
//		data.ClovekDisabled = "disabled"
//		data.SerizeniDisabled = "disabled"
//		data.StrojDisabled = "checked"
//	} else {
//		data.ClovekDisabled = "checked"
//		data.SerizeniDisabled = ""
//		data.StrojDisabled = ""
//	}
//	data.NokTypes = nokTypes
//	data.Workplaces = workplaces
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.UserDisabled = "disabled"
//	data.RadioDisabled = ""
//
//	data.TransferOrderButton = ""
//	data.OkValue = ok[0]
//	data.NokValue = nok[0]
//	data.DisplayOrder = GetActualDataForUser(userid)
//
//	_ = tmpl.Execute(*writer, data)
//}
//
//func CheckOpenOrderForWorkplaceInZapsi(orderid []string, operationid []string, workplaceid []string) (bool, string) {
//	order, suffix := ParseOrder(orderid[0])
//	operation := ParseOperation(operationid[0])
//	orderName := order + "." + suffix + "-" + operation
//	var zapsiOrder Order
//	var zapsiWorkplace Workplace
//	var terminalInputOrder TerminalInputOrder
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//
//	if err != nil {
//		logError("MAIN", "Problem opening database: "+err.Error())
//		return false, ""
//	}
//	defer db.Close()
//	db.Where("Name = ?", orderName).Find(&zapsiOrder)
//	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
//	db.Debug().Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
//	if terminalInputOrder.OID > 0 {
//		return true, terminalInputOrder.Note
//	}
//	return false, ""
//}
//
//func EnableClovekTransferInput(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	var nokTypes []SytelineNok
//	nokType := SytelineNok{Nazev: noktype[0]}
//	nokTypes = append(nokTypes, nokType)
//	data.NokTypes = nokTypes
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.UserDisabled = "disabled"
//	data.RadioDisabled = ""
//	data.ClovekDisabled = "checked"
//	data.TransferOrderButton = ""
//	data.OkValue = ok[0]
//	data.NokValue = nok[0]
//	data.DisplayOrder = GetActualDataForUser(userid)
//	_ = tmpl.Execute(*writer, data)
//}
//
//func EnableClovekSerizeniStrojTransferCloseInput(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string, noktype []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	var nokTypes []SytelineNok
//	nokType := SytelineNok{Nazev: noktype[0]}
//	nokTypes = append(nokTypes, nokType)
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	_, orderNote := CheckOpenOrderForWorkplaceInZapsi(orderid, operationid, workplaceIdSplitted)
//	println("Order note: " + orderNote)
//	if orderNote == "clovek" {
//		data.ClovekDisabled = "checked"
//		data.SerizeniDisabled = "disabled"
//		data.StrojDisabled = "disabled"
//	} else if orderNote == "serizeni" {
//		data.ClovekDisabled = "disabled"
//		data.SerizeniDisabled = "checked"
//		data.StrojDisabled = "disabled"
//	} else if orderNote == "stroj" {
//		data.ClovekDisabled = "disabled"
//		data.SerizeniDisabled = "disabled"
//		data.StrojDisabled = "checked"
//	} else {
//		data.ClovekDisabled = "checked"
//		data.SerizeniDisabled = ""
//		data.StrojDisabled = ""
//	}
//	data.NokTypes = nokTypes
//	data.Workplaces = workplaces
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.UserDisabled = "disabled"
//	data.RadioDisabled = ""
//
//	data.TransferOrderButton = ""
//	data.EndOrderButton = ""
//	data.OkValue = ok[0]
//	data.NokValue = nok[0]
//	data.DisplayOrder = GetActualDataForUser(userid)
//
//	_ = tmpl.Execute(*writer, data)
//}
//
//func FirstControls(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string) {
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	logInfo("MAIN", "Starting first controls")
//	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
//	data := CreateDefaultPage()
//	sytelineOperation := GetOperationFromSyteline(orderid, operationid)
//	sytelineWorkplace := GetWorkplaceFromSyteline(orderid, operationid, workplaceIdSplitted)
//	anyOpenOrderInZapsi := CheckAnyOpenOrderInZapsi(workplaceIdSplitted)
//	if anyOpenOrderInZapsi {
//		logInfo("MAIN", workplaceIdSplitted[0]+" has open order in Zapsi")
//		data.Message += "V Zapsi existuje pro toto pracoviste otevrena zakazka\n"
//		sameOrder, sameUser := CheckUserAndOrderInZapsi(userid, orderid, operationid, workplaceIdSplitted)
//		if sameOrder {
//			if sameUser {
//				logInfo("MAIN", workplaceIdSplitted[0]+" has open this exact order in Zapsi")
//				data.Message += "V Zapsi existuje pro toto pracoviste otevrena zakazka\n"
//				EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
//			} else {
//				if sytelineWorkplace.typ_zdroje_zapsi == "0" {
//					logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
//					data.Message += "typ_zdroje_zapsi je 0\n"
//					EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//				} else {
//					logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
//					data.Message += "typ_zdroje_zapsi neni 0\n"
//					EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//				}
//			}
//		} else {
//			logInfo("MAIN", workplaceid[0]+" has not open this exact order in Zapsi")
//			data.Message += "V Zapsi NEexistuje pro toto pracoviste presne tato otevrena zakazka\n"
//			if sytelineWorkplace.vice_vp == "1" {
//				logInfo("MAIN", "sytelineOperationSource.vice_vp equals one")
//				data.Message += "vice_vp je 1\n"
//				if sytelineWorkplace.typ_zdroje_zapsi == "0" {
//					logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi equals zero")
//					data.Message += "typ_zdroje_zapsi je 0\n"
//					EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//				} else {
//					logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi does not equal zero")
//					data.Message += "typ_zdroje_zapsi neni 0\n"
//					EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//				}
//			} else {
//				logInfo("MAIN", "sytelineOperationSource.vice_vp does not equal one")
//				data.Message += "vice_vp neni 1\n"
//				if sytelineOperation.parovy_dil == "1" {
//					logInfo("MAIN", "sytelineOperation.parovy_dil equals one")
//					data.Message += "parovy dil je 1\n"
//					if len(sytelineOperation.seznamm_par_dilu) > 0 {
//						logInfo("MAIN", "sytelineOperation.seznamm_par_dilu not empty: "+sytelineOperation.seznamm_par_dilu)
//						data.Message += "seznamm_par_dilu osahuje nejaky parovy dil\n"
//						var zapsiProducts = CheckProductsInZapsi(sytelineOperation)
//						anyOpenOrderHasOneOfProducts := CheckIfAnyOpenOrderHasOneOfProducts(workplaceid, zapsiProducts)
//						if anyOpenOrderHasOneOfProducts {
//							logInfo("MAIN", "Some open order in Zapsi contains any of pair parts")
//							data.Message += "Nejaka otevrena zakazka v Zapsi obsahuje nektery z parovych dilu\n"
//							if sytelineWorkplace.typ_zdroje_zapsi == "0" {
//								logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is zero")
//								data.Message += "typ_zdroje_zapsi je 0\n"
//								EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//							} else {
//								logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is not zero")
//								data.Message += "typ_zdroje_zapsi neni 0\n"
//								EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//							}
//						} else {
//							logInfo("MAIN", "No open order in Zapsi contains any of pair parts")
//							CheckOperationInSyteline(writer, userid, orderid, operationid)
//						}
//					} else {
//						logInfo("MAIN", "sytelineOperation.seznamm_par_dilu is empty")
//						data.Message += "seznamm_par_dilu neosahuje zadny parovy dil\n"
//						CheckOperationInSyteline(writer, userid, orderid, operationid)
//					}
//				} else {
//					logInfo("MAIN", "sytelineOperation.parovy_dil does not equal one")
//					data.Message += "parovy dil neni 1\n"
//					CheckOperationInSyteline(writer, userid, orderid, operationid)
//				}
//			}
//		}
//	} else {
//		logInfo("MAIN", workplaceid[0]+" does not have any open order in Zapsi")
//		data.Message += "V Zapsi neexistuje pro toto pracoviste otevrena zakazka\n"
//		if sytelineOperation.jen_prenos_mnozstvi == "1" {
//			logInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is one")
//			data.Message += "jen_prenos_mnozstvi je 1\n"
//			EnableOkNok(writer, workplaceid, userid, orderid, operationid, data, tmpl)
//		} else {
//			logInfo("MAIN", "sytelineOperation.jen_prenos_mnozstvi is not one")
//			data.Message += "jen_prenos_mnozstvi neni 1\n"
//			if sytelineWorkplace.typ_zdroje_zapsi == "0" {
//				logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is zero")
//				data.Message += "typ_zdroje_zapsi je 0\n"
//				EnableClovekStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//			} else {
//				logInfo("MAIN", "sytelineWorkplace.typ_zdroje_zapsi is not zero")
//				data.Message += "typ_zdroje_zapsi neni 0\n"
//				EnableClovekSerizeniStrojStart(writer, data, userid, orderid, operationid, workplaceid, tmpl)
//			}
//		}
//	}
//}
//
//func EnableClovekSerizeniStrojStart(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	workplaceIdSplitted := workplaceid
//	if strings.Contains(workplaceid[0], ";") {
//		workplaceIdSplitted = strings.Split(workplaceid[0], ";")
//	}
//	_, orderNote := CheckOpenOrderForWorkplaceInZapsi(orderid, operationid, workplaceIdSplitted)
//	println("Order note: " + orderNote)
//	if orderNote == "clovek" {
//		data.ClovekDisabled = "checked"
//		data.SerizeniDisabled = "disabled"
//		data.StrojDisabled = "disabled"
//	} else if orderNote == "serizeni" {
//		data.ClovekDisabled = "disabled"
//		data.SerizeniDisabled = "checked"
//		data.StrojDisabled = "disabled"
//	} else if orderNote == "stroj" {
//		data.ClovekDisabled = "disabled"
//		data.SerizeniDisabled = "disabled"
//		data.StrojDisabled = "checked"
//	} else {
//		data.ClovekDisabled = "checked"
//		data.SerizeniDisabled = ""
//		data.StrojDisabled = ""
//	}
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.UserDisabled = "disabled"
//	data.RadioDisabled = ""
//	data.StartOrderButton = ""
//	data.DisplayOrder = GetActualDataForUser(userid)
//	_ = tmpl.Execute(*writer, data)
//}
//
//func EnableClovekStart(writer *http.ResponseWriter, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.UserDisabled = "disabled"
//	data.RadioDisabled = ""
//	data.ClovekDisabled = "checked"
//	data.StartOrderButton = ""
//	data.DisplayOrder = GetActualDataForUser(userid)
//	_ = tmpl.Execute(*writer, data)
//}
//
//func EnableOkNok(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string, data RostraMainPage, tmpl *template.Template) {
//	nokTypes := GetNokTypesFromSyteline()
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.NokTypes = nokTypes
//	data.UserDisabled = "disabled"
//	data.OkDisabled = ""
//	data.NokDisabled = ""
//	data.OkValue = "0"
//	data.NokValue = "0"
//	data.OkFocus = "autofocus"
//	data.DisplayOrder = GetActualDataForUser(userid)
//	_ = tmpl.Execute(*writer, data)
//}
//
//func EnableTransfer(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string, data RostraMainPage, ok []string, nok []string, noktype []string, tmpl *template.Template) {
//	var workplaces []SytelineWorkplace
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceid[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	workplaces = append(workplaces, workplace)
//	var nokTypes []SytelineNok
//	nokType := SytelineNok{Nazev: noktype[0]}
//	nokTypes = append(nokTypes, nokType)
//	data.NokTypes = nokTypes
//	data.Workplaces = workplaces
//	data.UsernameValue = userid[0]
//	data.OrderValue = orderid[0]
//	data.Operation = operationid[0]
//	data.OperationValue = operationid[0]
//	data.DisplayOrder = GetActualDataForUser(userid)
//	data.UserDisabled = "disabled"
//	if len(ok) > 0 {
//		data.OkValue = ok[0]
//	}
//	if len(nok) > 0 {
//		data.NokValue = nok[0]
//	}
//	data.TransferOrderButton = ""
//	_ = tmpl.Execute(*writer, data)
//}
//
//func CheckAmount(inputAmount []string, sytelineWorkplace SytelineWorkplace, data RostraMainPage, userid []string, orderid []string, operationid []string, workplaceid []string, sytelineOperation SytelineOperation) bool {
//	amountCheck := true
//	checkAmount, err := strconv.Atoi(inputAmount[0])
//	if err != nil {
//		logError("MAIN", "Problem parsing data from user")
//		return false
//	}
//	logInfo("MAIN", "Checking amount: "+strconv.Itoa(checkAmount))
//	if sytelineWorkplace.priznak_mn_1 == "1" {
//		logInfo("MAIN", "sytelineWorkplace.priznak_mn_1 is one")
//		data.Message += "priznak_mn_1 je 1\n"
//		amountLessThanInZapsi := CheckIfOperatorAmountLessThanInZapsi(inputAmount, userid, orderid, operationid, workplaceid)
//		if amountLessThanInZapsi {
//			logInfo("MAIN", "ok from user is less than in zapsi")
//			data.Message += "uzivatel zadal mene kusu nez je v zapsi, coz je spravne\n"
//		} else {
//			logInfo("MAIN", "ok from user is more than in zapsi")
//			data.Message += "uzivatel zadal vice kusu nez je v zapsi, coz je spatne\n"
//			amountCheck = false
//		}
//	} else {
//		logInfo("MAIN", "sytelineWorkplace.priznak_mn_1 is not one")
//		data.Message += "priznak_mn_1 neni 1\n"
//	}
//	if sytelineOperation.priznak_mn_2 == "1" {
//		logInfo("MAIN", "sytelineWorkplace.priznak_mn_2 is one")
//		data.Message += "priznak_mn_2 je 1\n"
//		parsedFromSyteline, err := strconv.ParseFloat(sytelineOperation.mn_2_ks, 64)
//		if err != nil {
//			logError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_2_ks+", "+err.Error())
//		}
//		parsedFromInput, err := strconv.Atoi(inputAmount[0])
//		if err != nil {
//			logError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_2_ks+", "+err.Error())
//		}
//		if parsedFromInput <= int(parsedFromSyteline) {
//			logInfo("MAIN", "ok from user is less than mn_2_ks")
//			data.Message += "uzivatel zadal mene kusu nez je v mn_2_ks, coz je spravne\n"
//		} else {
//			logInfo("MAIN", "ok from user is more than mn_2_ks")
//			data.Message += "uzivatel zadal vic kusu nez je v mn_2_ks, coz je spatne\n"
//			amountCheck = false
//		}
//	} else {
//		logInfo("MAIN", "sytelineWorkplace.priznak_mn_2 is not one")
//		data.Message += "priznak_mn_2 neni 1\n"
//	}
//	if sytelineOperation.priznak_mn_3 == "1" {
//		logInfo("MAIN", "sytelineWorkplace.priznak_mn_3 is one")
//		data.Message += "priznak_mn_3 je 1\n"
//		parsedFromSyteline, err := strconv.ParseFloat(sytelineOperation.mn_3_ks, 64)
//		if err != nil {
//			logError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_3_ks+", "+err.Error())
//		}
//		parsedFromInput, err := strconv.Atoi(inputAmount[0])
//		if err != nil {
//			logError("MAIN", "Problem parsing mn_2_ks: "+sytelineOperation.mn_3_ks+", "+err.Error())
//		}
//		if parsedFromInput <= int(parsedFromSyteline) {
//			logInfo("MAIN", "ok from user is less than mn_3_ks")
//			data.Message += "uzivatel zadal mene kusu nez je v mn_3_ks, coz je spravne\n"
//		} else {
//			logInfo("MAIN", "ok from user is more than mn_3_ks")
//			data.Message += "uzivatel zadal vic kusu nez je v mn_3_ks, coz je spatne\n"
//			amountCheck = false
//		}
//
//	} else {
//		logInfo("MAIN", "sytelineWorkplace.priznak_mn_3 is not one")
//		data.Message += "priznak_mn_3 neni 1\n"
//	}
//	logInfo("MAIN", "Returning from amount check: "+strconv.FormatBool(amountCheck))
//	return amountCheck
//}
