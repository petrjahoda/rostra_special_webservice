package main

import (
	"html/template"
	"net/http"
)

//
//import (
//	"github.com/jinzhu/gorm"
//	"strings"
//)

func FirstControls(writer *http.ResponseWriter, workplaceid []string, userid []string, orderid []string, operationid []string) {
	LogInfo("MAIN", "Starting controls")
	tmpl := template.Must(template.ParseFiles("html/rostra.html"))
	data := CreateDefaultPage()
	data.Message = "Starting controls"
	var anyOpenOrderInZapsi bool = CheckAnyOpenOrderInZapsi(workplaceid)
	if anyOpenOrderInZapsi {
		var thisOpenOrderInZapsi bool = CheckThisOpenOrderInZapsi(workplaceid, orderid, operationid, userid)
		if thisOpenOrderInZapsi {
			// Zobraz vstup pro OK a NOK kusy
		} else {
			var jeViceVpEqualOne bool = CheckJeViceVpEqualOne(orderid, operationid)
			if jeViceVpEqualOne {
				var jeTypZdrojeZapsiEqualZero bool = CheckJeTypZdrojeZapsiEqualZero(orderid, operationid)
				if jeTypZdrojeZapsiEqualZero {
					// zobraz volbu clovek, zobraz tlacitko zahajeni
				} else {
					// zobraz clovek, stroj a serizeni a tlacitko zahajeni
				}
			} else {
				var jeParovyDilEqualOne bool = CheckJeParovyDilEqualOne(orderid, operationid)
				if jeParovyDilEqualOne {
					var seznamParDiluIsNotEmpty bool = CheckSeznamParDiluIsNotEmpty(orderid, operationid)
					if seznamParDiluIsNotEmpty {
						var zapsiProduct Product = CheckProductInZapsiIfNotExists(orderid, operationid)
						var anyOpenOrderHasThisProduct bool = CheckIfAnyOpenOrderHasThisProduct(workplaceid, zapsiProduct)
						if anyOpenOrderHasThisProduct {
							var typZdrojeZapsiEqualZero = CheckTypZdrojeZapsiEqualZero(orderid, operationid, workplaceid)
							if typZdrojeZapsiEqualZero {
								// zobraz volbu clovek, zobraz tlacitko zahajeni
							} else {
								//zobraz clovek, stroj a serizeni a tlacitko zahajeni
							}
						} else {
							// vrat na vyber pracoviste s informaci o chybe
						}
					} else {
						// vrat na vyber pracoviste s informaci o chybe
					}
				} else {
					// vrat na vyber pracoviste s informaci o chybe
				}
			}
		}
	} else {
		var jenPrenosMnozstviEqualOne bool = CheckJenPrenosMnozstviEqualOne(orderid, operationid)
		if jenPrenosMnozstviEqualOne {
			// Zobraz vstup pro OK a NOK kusy
		} else {
			var jeTypZdrojeEqualZero bool = CheckJeTypZdrojeEqualZero(orderid, operationid, workplaceid)
			if jeTypZdrojeEqualZero {
				// zobraz volbu clovek a tlacitko zahajeni
			} else {
				// zobraz clovek, stroj a serizeni a tlacitko zahajeni
			}
		}
	}
	_ = tmpl.Execute(*writer, data)
}

//
//func MakeFirstControls(workplaceId []string, userId []string, orderId []string, operationId []string) {
//	anyOrderExists := CheckAnyOrderInZapsi(workplaceId)
//	if anyOrderExists {
//		LogInfo("MAIN", "Some open order in Zapsi already exists")
//		data.Message = "Some open order in Zapsi already exists"
//		thisOrderIsOpen := CheckThisOrderInZapsi(userId, orderId, operationId, workplaceId)
//		if thisOrderIsOpen {
//			LogInfo("MAIN", "This order in Zapsi already exists, enabling end and transfer button")
//			data.Message = "This order in Zapsi already exists, enabling end and transfer button"
//			EnableOkAndNokInput(workplaceId, userId, orderId, operationId, data)
//			//EnableTransferAndEndButton(workplaceId, userId, orderId, operationId, data)
//		} else {
//			LogInfo("MAIN", "This order in Zapsi not exists, checking for vice_vp")
//			data.Message = "This order in Zapsi not exists, checking for vice_vp"
//			sytelineOperation, sytelineWorkplaces := CheckOperationInSyteline(userId, orderId, operationId)
//			for _, workplace := range sytelineWorkplaces {
//				if workplace.Zapsi_zdroj == workplaceId[0] {
//					if workplace.vice_vp == "1" {
//						LogInfo("MAIN", workplace.Zapsi_zdroj+" has parameter vice_vp == 1, enabling start button")
//						data.Message = workplace.Zapsi_zdroj + " has parameter vice_vp == 1, enabling start button"
//						EnableStartButton(workplaceId, userId, orderId, operationId, data)
//						break
//					} else {
//						LogInfo("MAIN", workplace.Zapsi_zdroj+" without parameter vice_vp == 1")
//						data.Message = workplace.Zapsi_zdroj + " without parameter vice_vp == 1"
//						if sytelineOperation.parovy_dil == "1" {
//							LogInfo("MAIN", workplace.Zapsi_zdroj+" has parameter parovy_dil == 1")
//							data.Message = workplace.Zapsi_zdroj + " has parameter parovy_dil == 1"
//							if len(sytelineOperation.seznamm_par_dilu) > 1 {
//								productName := GetProductNameForOpenOrder(workplaceId)
//								LogInfo("MAIN", workplace.Zapsi_zdroj+" has parameter seznamm_par_dilu > 1")
//								data.Message = workplace.Zapsi_zdroj + " has parameter seznamm_par_dilu > 1"
//								if strings.Contains(sytelineOperation.seznamm_par_dilu, productName) {
//									LogInfo("MAIN", "Products are matching, enabling start button")
//									data.Message = "Products are matching, enabling start button"
//									EnableStartButton(workplaceId, userId, orderId, operationId, data)
//								} else {
//									LogInfo("MAIN", "Products not matching: ["+sytelineOperation.seznamm_par_dilu+"] ["+productName+"]")
//									data.Message = "Products not matching: [" + sytelineOperation.seznamm_par_dilu + "] [" + productName + "]"
//									EnableWorkplaceSelect(userId, orderId, operationId, data)
//								}
//							} else {
//								LogInfo("MAIN", workplace.Zapsi_zdroj+" has parameter seznamm_par_dilu == 0")
//								data.Message = workplace.Zapsi_zdroj + " has parameter seznamm_par_dilu == 0"
//								EnableWorkplaceSelect(userId, orderId, operationId, data)
//							}
//						} else {
//							LogInfo("MAIN", workplace.Zapsi_zdroj+" without parameter parovy_dil == 1")
//							data.Message = workplace.Zapsi_zdroj + " without parameter parovy_dil == 1"
//							EnableWorkplaceSelect(userId, orderId, operationId, data)
//						}
//						break
//					}
//				}
//			}
//		}
//	} else {
//		LogInfo("MAIN", "No open order in Zapsi exists")
//		data.Message = "No open order in Zapsi exists"
//		sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId)
//		if sytelineOperation.jen_prenos_mnozstvi == "1" {
//			LogInfo("MAIN", "Operation has only data transfer, enabling transfer button")
//			data.Message = "Operation has only data transfer, enabling transfer button"
//			EnableOkAndNokInput(workplaceId, userId, orderId, operationId, data)
//			//EnableTransferButton(workplaceId, userId, orderId, operationId, data)
//		} else {
//			LogInfo("MAIN", "Enabling start button")
//			data.Message = "Enabling start button"
//			EnableStartButton(workplaceId, userId, orderId, operationId, data)
//		}
//	}
//}
//
//func EnableOkAndNokInput(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage) {
//	var nokTypes []SytelineNok
//	db, err := gorm.Open("mssql", SytelineConnection)
//	defer db.Close()
//
//	command := "declare @JePlatny ListYesNoType, @Kod ReasonCodeType = NULL exec [rostra_exports_test].dbo.ZapsiKodyDuvoduZmetkuSp @Kod= @Kod, @JePlatny = @JePlatny output select JePlatny = @JePlatny"
//	rows, err := db.Raw(command).Rows()
//	if err != nil {
//		LogError("MAIN", "Error: "+err.Error())
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var nokType SytelineNok
//		err = rows.Scan(&nokType.Kod, &nokType.Nazev)
//		nokTypes = append(nokTypes, nokType)
//		if err != nil {
//			LogError("MAIN", "Error: "+err.Error())
//		}
//	}
//	data.NokTypes = nokTypes
//	data.OkDisabled = ""
//	data.NokDisabled = ""
//	data.OkFocus = "autofocus"
//	data.UsernameValue = userId[0]
//	data.OrderValue = orderId[0]
//	data.OperationValue = operationId[0]
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	data.Workplaces = append(data.Workplaces, workplace)
//}
//
//func GetProductNameForOpenOrder(workplaceId []string) string {
//	var zapsiWorkplace Workplace
//	var zapsiOrder Order
//	var zapsiProduct Product
//	var terminalInputOrder TerminalInputOrder
//	connectionString, dialect := CheckDatabaseType()
//	db, err := gorm.Open(dialect, connectionString)
//	if err != nil {
//		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
//		return ""
//	}
//	defer db.Close()
//	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
//	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Find(&terminalInputOrder)
//	db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
//	db.Where("OID = ?", zapsiOrder.ProductID).Find(&zapsiProduct)
//	if len(zapsiProduct.Name) > 0 {
//		LogInfo("MAIN", "Found product "+zapsiProduct.Name)
//		return zapsiProduct.Name
//	} else {
//		LogInfo("MAIN", "Product not found")
//		return ""
//	}
//}
//
//func EnableWorkplaceSelect(userId []string, orderId []string, operationId []string, data *RostraMainPage) {
//	order, suffix := ParseOrder(orderId[0])
//	db, err := gorm.Open("mssql", SytelineConnection)
//	var sytelineWorkplaces []SytelineWorkplace
//	if err != nil {
//		LogError("MAIN", "Error opening db: "+err.Error())
//		data.UsernameValue = userId[0]
//		data.OrderValue = orderId[0]
//		data.Operation = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
//		data.OperationDisabled = ""
//		data.OperationFocus = "autofocus"
//		return
//	}
//	defer db.Close()
//	command := "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationId[0] + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
//	workplaceRows, err := db.Raw(command).Rows()
//	if err != nil {
//		LogError("MAIN", "Error: "+err.Error())
//	}
//	defer workplaceRows.Close()
//	for workplaceRows.Next() {
//		var sytelineWorkplace SytelineWorkplace
//		err = workplaceRows.Scan(&sytelineWorkplace.Zapsi_zdroj, &sytelineWorkplace.priznak_mn_1, &sytelineWorkplace.vice_vp, &sytelineWorkplace.SL_prac, &sytelineWorkplace.auto_prevod_mnozstvi, &sytelineWorkplace.mnozstvi_auto_prevodu)
//		sytelineWorkplaces = append(sytelineWorkplaces, sytelineWorkplace)
//		if err != nil {
//			LogError("MAIN", "Error: "+err.Error())
//		}
//	}
//	data.WorkplaceDisabled = ""
//	data.Workplaces = sytelineWorkplaces
//	data.WorkplaceFocus = "autofocus"
//}
//
//func EnableTransferButton(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage, ok []string, nok []string, noktype []string) {
//	LogInfo("MAIN", "Operation has only data transfer, enabling transfer button")
//	data.TransferOrderButton = ""
//	data.UsernameValue = userId[0]
//	data.OrderValue = orderId[0]
//	data.OperationValue = operationId[0]
//	data.Ok = ok[0]
//	data.Nok = nok[0]
//	data.OkValue = ok[0]
//	data.NokValue = nok[0]
//	data.RadioDisabled = ""
//	nokType := SytelineNok{Nazev: noktype[0]}
//	data.NokTypes = append(data.NokTypes, nokType)
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	data.Workplaces = append(data.Workplaces, workplace)
//}
//
//func EnableStartButton(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage) {
//	data.StartOrderButton = ""
//	data.UsernameValue = userId[0]
//	data.OrderValue = orderId[0]
//	data.OperationValue = operationId[0]
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	data.Workplaces = append(data.Workplaces, workplace)
//}
//
//func EnableTransferAndEndButton(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage, ok []string, nok []string, noktype []string) {
//	LogInfo("MAIN", "This order in Zapsi already exists, enabling end and transfer button")
//	data.EndOrderButton = ""
//	data.TransferOrderButton = ""
//	data.UsernameValue = userId[0]
//	data.OrderValue = orderId[0]
//	data.OperationValue = operationId[0]
//	data.Ok = ok[0]
//	data.Nok = nok[0]
//	data.OkValue = ok[0]
//	data.NokValue = nok[0]
//	data.RadioDisabled = ""
//	nokType := SytelineNok{Nazev: noktype[0]}
//	data.NokTypes = append(data.NokTypes, nokType)
//	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
//	data.Workplaces = append(data.Workplaces, workplace)
//}
