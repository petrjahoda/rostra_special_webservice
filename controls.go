package main

import (
	"github.com/jinzhu/gorm"
	"strings"
)

func MakeControls(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	anyOrderExists := CheckAnyOrderInZapsi(workplaceId)
	if anyOrderExists {
		LogInfo("MAIN", "Some open order in Zapsi already exists")
		thisOrderIsOpen := CheckThisOrderInZapsi(userId, orderId, operationId, workplaceId)
		if thisOrderIsOpen {
			LogInfo("MAIN", "This order in Zapsi already exists, enabling end and transfer button")
			EnableTransferAndEndButton(workplaceId, userId, orderId, operationId, data)
		} else {
			LogInfo("MAIN", "This order in Zapsi not exists, checking for vice_vp")
			sytelineOperation, sytelineWorkplaces := CheckOperationInSyteline(userId, orderId, operationId, data)
			for _, workplace := range sytelineWorkplaces {
				if workplace.Zapsi_zdroj == workplaceId[0] {
					if workplace.vice_vp == "1" {
						LogInfo("MAIN", workplace.Zapsi_zdroj+"has parameter vice_vp == 1, enabling start button")
						EnableStartButton(workplaceId, userId, orderId, operationId, data)
						break
					} else {
						LogInfo("MAIN", workplace.Zapsi_zdroj+"without parameter vice_vp == 1")
						if sytelineOperation.parovy_dil == "1" {
							LogInfo("MAIN", workplace.Zapsi_zdroj+"has parameter parovy_dil == 1")
							if len(sytelineOperation.seznamm_par_dilu) > 1 {
								productName := GetProductNameForOpenOrder(workplaceId)
								LogInfo("MAIN", workplace.Zapsi_zdroj+"has parameter seznamm_par_dilu > 1")
								if strings.Contains(sytelineOperation.seznamm_par_dilu, productName) {
									LogInfo("MAIN", "Products are matching, enabling start button")
									EnableStartButton(workplaceId, userId, orderId, operationId, data)
								} else {
									LogInfo("MAIN", "Products not matching")
									EnableWorkplaceSelect(userId, orderId, operationId, data)
								}
							} else {
								LogInfo("MAIN", workplace.Zapsi_zdroj+"has parameter seznamm_par_dilu == 0")
								EnableWorkplaceSelect(userId, orderId, operationId, data)
							}
						} else {
							LogInfo("MAIN", workplace.Zapsi_zdroj+"without parameter parovy_dil == 1")
							EnableWorkplaceSelect(userId, orderId, operationId, data)
						}
						break
					}
				}
			}
		}
	} else {
		LogInfo("MAIN", "No open order in Zapsi exists")
		sytelineOperation, _ := CheckOperationInSyteline(userId, orderId, operationId, data)
		if sytelineOperation.jen_prenos_mnozstvi == "1" {
			LogInfo("MAIN", "Operation has only data transfer, enabling transfer button")
			EnableTransferButton(workplaceId, userId, orderId, operationId, data)
		} else {
			LogInfo("MAIN", "Enabling start button")
			EnableStartButton(workplaceId, userId, orderId, operationId, data)

		}
	}
}

func GetProductNameForOpenOrder(workplaceId []string) string {
	var zapsiWorkplace Workplace
	var zapsiOrder Order
	var zapsiProduct Product
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return ""
	}
	defer db.Close()
	db.Where("Code = ?", workplaceId).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Find(&terminalInputOrder)
	db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
	db.Where("OID = ?", zapsiOrder.ProductID).Find(&zapsiProduct)
	if len(zapsiProduct.Name) > 0 {
		LogInfo("MAIN", "Found product "+zapsiProduct.Name)
		return zapsiProduct.Name
	} else {
		LogInfo("MAIN", "Product not found")
		return ""
	}
}

func EnableWorkplaceSelect(userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	order, suffix := ParseOrder(orderId[0])
	db, err := gorm.Open("mssql", SytelineConnection)
	var sytelineWorkplaces []SytelineWorkplace
	if err != nil {
		LogError("MAIN", "Error opening db: "+err.Error())
		data.UsernameValue = userId[0]
		data.OrderValue = orderId[0]
		data.Operation = "Problém při komunikaci se Syteline, kontaktujte prosím IT"
		data.OperationDisabled = ""
		data.OperationFocus = "autofocus"
		return
	}
	defer db.Close()
	command := "declare @CisloVP JobType, @PriponaVP SuffixType, @Operace OperNumType select   @CisloVP = N'" + order + "', @PriponaVP = " + suffix + ", @Operace = " + operationId[0] + " exec dbo.ZapsiZdrojeOperaceSp @CisloVP = @CisloVP, @PriponaVp = @PriponaVP , @Operace = @Operace"
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
	data.WorkplaceDisabled = ""
	data.Workplaces = sytelineWorkplaces
	data.WorkplaceFocus = "autofocus"
}

func EnableTransferButton(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	data.TransferOrderButton = ""
	data.UsernameValue = userId[0]
	data.OrderValue = orderId[0]
	data.OperationValue = operationId[0]
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	data.Workplaces = append(data.Workplaces, workplace)
}

func EnableStartButton(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	data.StartOrderButton = ""
	data.UsernameValue = userId[0]
	data.OrderValue = orderId[0]
	data.OperationValue = operationId[0]
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	data.Workplaces = append(data.Workplaces, workplace)
}

func EnableTransferAndEndButton(workplaceId []string, userId []string, orderId []string, operationId []string, data *RostraMainPage) {
	data.EndOrderButton = ""
	data.TransferOrderButton = ""
	data.UsernameValue = userId[0]
	data.OrderValue = orderId[0]
	data.OperationValue = operationId[0]
	workplace := SytelineWorkplace{Zapsi_zdroj: workplaceId[0], priznak_mn_1: "", vice_vp: "", SL_prac: "", auto_prevod_mnozstvi: "", mnozstvi_auto_prevodu: ""}
	data.Workplaces = append(data.Workplaces, workplace)
}
