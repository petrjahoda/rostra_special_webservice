package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CountInputData struct {
	WorkplaceCode      string
	UserId             string
	UserInput          string
	OrderInput         string
	OperationSelect    string
	ParovyDil          string
	SeznamParovychDilu string
	JenPrenosMnozstvi  string
	TypZdrojeZapsi     string
	ViceVp             string
	PriznakMn1         string
	PriznakMn2         string
	PriznakMn3         string
	Mn2Ks              string
	Mn3Ks              string
	OkCount            string
	NokCount           string
}

type CountResponseData struct {
	Result      string
	Transfer    string
	End         string
	Clovek      string
	Stroj       string
	Serizeni    string
	CountError  string
	RostraError string
}

func checkCountInput(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check count", "Started")
	var data CountInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check count", "Error parsing input: "+err.Error())
		var responseData CountResponseData
		responseData.Result = "nok"
		responseData.CountError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check count", "Ended with error")
		return
	}
	logInfo("Check workplace", "Data: "+data.WorkplaceCode+", "+data.OrderInput+", "+data.UserId+", "+data.UserInput+", "+data.OperationSelect+
		", "+data.ParovyDil+", "+data.SeznamParovychDilu+", "+data.JenPrenosMnozstvi+",  "+data.TypZdrojeZapsi+", "+data.ViceVp+
		", "+data.PriznakMn1+", "+data.PriznakMn2+", "+data.PriznakMn3+", "+data.Mn2Ks+", "+data.Mn3Ks)

	countFromZapsi, terminalInputOrder := GetCountForActualOpenOrder(data.WorkplaceCode, data.UserId, data.OrderInput, data.OperationSelect)
	countFromSyteline := GetCountForAllTransferredToSyteline(data.OrderInput, data.OperationSelect, terminalInputOrder, data.UserInput)
	logInfo("MAIN", "[CountZapsi:CountSyteline:CountUserOK:CountUserNOK] ["+strconv.Itoa(countFromZapsi)+":"+strconv.Itoa(countFromSyteline)+":"+data.OkCount+":"+data.NokCount+"]")
	countOkFromUser, err := strconv.Atoi(data.OkCount)
	if err != nil {
		logError("Check count", "Problem parsing count from user: "+err.Error())
		var responseData CountResponseData
		responseData.Result = "nok"
		responseData.CountError = "Problem parsing ok count from user: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check count", "Ended with error")
		return
	}
	countNokFromUser, err := strconv.Atoi(data.NokCount)
	if err != nil {
		logError("Check count", "Problem parsing nok count from user: "+err.Error())
		var responseData CountResponseData
		responseData.Result = "nok"
		responseData.CountError = "Problem parsing count from user: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check count", "Ended with error")
		return
	}
	okCountAsInt, err := strconv.Atoi(data.OkCount)
	if err != nil {
		logError("Check count", "Problem parsing ok count: "+err.Error())
	}
	nokCountAsInt, err := strconv.Atoi(data.NokCount)
	if err != nil {
		logError("Check count", "Problem parsing nok count: "+err.Error())
	}
	totalCount := okCountAsInt + nokCountAsInt
	switch data.JenPrenosMnozstvi {
	case "1":
		{
			logInfo("Check count", "sytelineOperation.JenPrenosMnozstvi is one, only transfer will be available")
			if data.PriznakMn1 == "1" {
				logInfo("Check count", "Priznak Mn1 is one")
				if (countOkFromUser + countNokFromUser) <= (countFromZapsi - countFromSyteline) {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is less or equal than difference between transferred ["+strconv.Itoa(countFromSyteline)+"] and actual count from Zapsi ["+strconv.Itoa(countFromZapsi)+"]")
				} else {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is more or equal than difference between transferred ["+strconv.Itoa(countFromSyteline)+"] and actual count from Zapsi ["+strconv.Itoa(countFromZapsi)+"]")
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.Transfer = "false"
					responseData.End = "false"
					responseData.RostraError = "Nelze provést přenos množství " + strconv.Itoa(totalCount) + "ks, v Zapsi je vyrobeno " + strconv.Itoa(countFromZapsi) + "ks, do Syteline již přeneseno " + strconv.Itoa(countFromSyteline) + "ks , je možno přenést maximálně " + strconv.Itoa(countFromZapsi-countFromSyteline) + "ks"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
			} else {
				logInfo("Check count", "Priznak Mn1 is not one")
			}
			if data.PriznakMn2 == "1" {
				logInfo("Check count", "Priznak Mn2 is one")
				mnozstvi2, err := strconv.Atoi(data.Mn2Ks)
				if err != nil {
					logError("Check count", "Problem parsing mn2 count: "+err.Error())
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.CountError = "Problem parsing mn2: " + err.Error()
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
				if (countOkFromUser + countNokFromUser) < mnozstvi2 {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is less or equal than mn2 ["+data.Mn2Ks+"]")
				} else {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is more or equal than mn2 ["+data.Mn2Ks+"]")
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.Transfer = "false"
					responseData.End = "false"
					responseData.RostraError = "Nelze provést přenos množství " + strconv.Itoa(totalCount) + "ks, z předchozí operace bylo předáno " + data.Mn2Ks + "ks"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
			}
			if data.PriznakMn3 == "1" {
				logInfo("Check count", "Priznak Mn3 is one")
				mnozstvi3, err := strconv.Atoi(data.Mn3Ks)
				if err != nil {
					logError("Check count", "Problem parsing mn3 count: "+err.Error())
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.CountError = "Problem parsing mn3: " + err.Error()
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
				if (countOkFromUser + countNokFromUser) < mnozstvi3 {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is less or equal than mn3 ["+data.Mn3Ks+"]")
				} else {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is more or equal than mn3 ["+data.Mn3Ks+"]")
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.Transfer = "false"
					responseData.End = "false"
					responseData.RostraError = "Nelze provést přenos množství " + strconv.Itoa(totalCount) + "ks, do operace bylo vydáno " + data.Mn3Ks + "ks"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
			}
			var responseData CountResponseData
			responseData.Result = "ok"
			responseData.Transfer = "true"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("Check count", "Ended successfully")
			return
		}
	default:
		{
			logInfo("Check count", "sytelineOperation.JenPrenosMnozstvi IS NOT one, transfer and close will be available")
			if data.PriznakMn1 == "1" {
				logInfo("Check count", "Priznak Mn1 is one")
				if (countOkFromUser + countNokFromUser) <= (countFromZapsi - countFromSyteline) {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is less or equal than difference between transferred ["+strconv.Itoa(countFromSyteline)+"] and actual count from Zapsi ["+strconv.Itoa(countFromZapsi)+"]")
				} else {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is more or equal than difference between transferred ["+strconv.Itoa(countFromSyteline)+"] and actual count from Zapsi ["+strconv.Itoa(countFromZapsi)+"]")
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.Transfer = "false"
					responseData.End = "false"
					responseData.RostraError = "Nelze provést přenos množství " + strconv.Itoa(totalCount) + "ks, v Zapsi je vyrobeno " + strconv.Itoa(countFromZapsi) + "ks, do Syteline již přeneseno " + strconv.Itoa(countFromSyteline) + "ks , je možno přenést maximálně " + strconv.Itoa(countFromZapsi-countFromSyteline) + "ks"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
			} else {
				logInfo("Check count", "Priznak Mn1 is not one")
			}
			if data.PriznakMn2 == "1" {
				logInfo("Check count", "Priznak Mn2 is one")
				mnozstvi2, err := strconv.Atoi(data.Mn2Ks)
				if err != nil {
					logError("Check count", "Problem parsing mn2 count: "+err.Error())
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.CountError = "Problem parsing mn2: " + err.Error()
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
				if (countOkFromUser + countNokFromUser) < mnozstvi2 {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is less or equal than mn2 ["+data.Mn2Ks+"]")
				} else {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is more or equal than mn2 ["+data.Mn2Ks+"]")
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.Transfer = "false"
					responseData.End = "false"
					responseData.RostraError = "Nelze provést přenos množství " + strconv.Itoa(totalCount) + "ks, z předchozí operace bylo předáno " + data.Mn2Ks + "ks"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
			}
			if data.PriznakMn3 == "1" {
				logInfo("Check count", "Priznak Mn3 is one")
				mnozstvi3, err := strconv.Atoi(data.Mn3Ks)
				if err != nil {
					logError("Check count", "Problem parsing mn3 count: "+err.Error())
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.CountError = "Problem parsing mn3: " + err.Error()
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
				if (countOkFromUser + countNokFromUser) < mnozstvi3 {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is less or equal than mn3 ["+data.Mn3Ks+"]")
				} else {
					logInfo("Check count", "OK and NOK ["+data.OkCount+":"+data.NokCount+"] from user is more or equal than mn3 ["+data.Mn3Ks+"]")
					var responseData CountResponseData
					responseData.Result = "nok"
					responseData.Transfer = "false"
					responseData.End = "false"
					responseData.RostraError = "Nelze provést přenos množství " + strconv.Itoa(totalCount) + "ks, do operace bylo vydáno " + data.Mn3Ks + "ks"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended with error")
					return
				}
			}
			if data.TypZdrojeZapsi == "0" {
				logInfo("Check count", "Typ Zdroje Zapsi is zero")
				if data.PriznakMn1 == "0" {
					logInfo("Check count", "Priznak mn1 is zero, enabling transfer and end")
					var responseData CountResponseData
					responseData.Result = "ok"
					responseData.Transfer = "true"
					responseData.End = "true"
					responseData.Clovek = "true"
					writer.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(writer).Encode(responseData)
					logInfo("Check count", "Ended successfully")
					return
				} else {
					logInfo("Check count", "Priznak mn1 is not zero")
					if (countOkFromUser + countNokFromUser) == (countFromZapsi - countFromSyteline) {
						logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling transfer and end")
						var responseData CountResponseData
						responseData.Result = "ok"
						responseData.Transfer = "true"
						responseData.End = "true"
						responseData.Clovek = "true"
						writer.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(writer).Encode(responseData)
						logInfo("Check count", "Ended successfully")
						return
					} else {
						logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling just transfer")
						var responseData CountResponseData
						responseData.Result = "ok"
						responseData.Transfer = "true"
						responseData.End = "false"
						responseData.Clovek = "true"
						writer.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(writer).Encode(responseData)
						logInfo("Check count", "Ended successfully")
						return
					}
				}

			} else {
				logInfo("Check count", "Typ Zdroje Zapsi IS NOT zero, actual terminal inpout order note: "+terminalInputOrder.Note)
				switch terminalInputOrder.Note {
				case "clovek":
					{
						if data.PriznakMn1 == "0" {
							logInfo("Check count", "Priznak mn1 is zero, enabling transfer and end")
							var responseData CountResponseData
							responseData.Result = "ok"
							responseData.Transfer = "true"
							responseData.End = "true"
							responseData.Clovek = "true"
							writer.Header().Set("Content-Type", "application/json")
							_ = json.NewEncoder(writer).Encode(responseData)
							logInfo("Check count", "Ended successfully")
							return
						} else {
							logInfo("Check count", "Priznak mn1 is not zero")
							if (countOkFromUser + countNokFromUser) == (countFromZapsi - countFromSyteline) {
								logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling transfer and end")
								var responseData CountResponseData
								responseData.Result = "ok"
								responseData.Transfer = "true"
								responseData.End = "true"
								responseData.Clovek = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								logInfo("Check count", "Ended successfully")
								return
							} else {
								logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling just transfer")
								var responseData CountResponseData
								responseData.Result = "ok"
								responseData.Transfer = "true"
								responseData.End = "false"
								responseData.Clovek = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								logInfo("Check count", "Ended successfully")
								return
							}
						}
					}
				case "stroj":
					{
						if data.PriznakMn1 == "0" {
							logInfo("Check count", "Priznak mn1 is zero, enabling transfer and end")
							var responseData CountResponseData
							responseData.Result = "ok"
							responseData.Transfer = "true"
							responseData.End = "true"
							responseData.Stroj = "true"
							writer.Header().Set("Content-Type", "application/json")
							_ = json.NewEncoder(writer).Encode(responseData)
							logInfo("Check count", "Ended successfully")
							return
						} else {
							logInfo("Check count", "Priznak mn1 is not zero")
							if (countOkFromUser + countNokFromUser) == (countFromZapsi - countFromSyteline) {
								logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling transfer and end")
								var responseData CountResponseData
								responseData.Result = "ok"
								responseData.Transfer = "true"
								responseData.End = "true"
								responseData.Stroj = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								logInfo("Check count", "Ended successfully")
								return
							} else {
								logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling just transfer")
								var responseData CountResponseData
								responseData.Result = "ok"
								responseData.Transfer = "true"
								responseData.End = "false"
								responseData.Stroj = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								logInfo("Check count", "Ended successfully")
								return
							}
						}
					}
				case "serizeni":
					{
						if data.PriznakMn1 == "0" {
							logInfo("Check count", "Priznak mn1 is zero, enabling transfer and end")
							var responseData CountResponseData
							responseData.Result = "ok"
							responseData.Transfer = "true"
							responseData.End = "true"
							responseData.Serizeni = "true"
							writer.Header().Set("Content-Type", "application/json")
							_ = json.NewEncoder(writer).Encode(responseData)
							logInfo("Check count", "Ended successfully")
							return
						} else {
							logInfo("Check count", "Priznak mn1 is not zero")
							if (countOkFromUser + countNokFromUser) == (countFromZapsi - countFromSyteline) {
								logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling transfer and end")
								var responseData CountResponseData
								responseData.Result = "ok"
								responseData.Transfer = "true"
								responseData.End = "true"
								responseData.Serizeni = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								logInfo("Check count", "Ended successfully")
								return
							} else {
								logInfo("Check count", "OK and NOK is equal to difference between Zapsi and Syteline, enabling just transfer")
								var responseData CountResponseData
								responseData.Result = "ok"
								responseData.Transfer = "true"
								responseData.End = "false"
								responseData.Serizeni = "true"
								writer.Header().Set("Content-Type", "application/json")
								_ = json.NewEncoder(writer).Encode(responseData)
								logInfo("Check count", "Ended successfully")
								return
							}
						}
					}
				}
			}
		}
	}
}

func GetCountForAllTransferredToSyteline(orderInput string, operationInput string, terminalInputOrder TerminalInputOrder, userInput string) int {
	transferredTotalThisOrder := 0
	order, suffix := ParseOrder(orderInput)
	operation := ParseOperation(operationInput)
	orderName := order + "." + suffix + "-" + operation
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check count", "Problem opening database: "+err.Error())
		return 0
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiTransThisOrder []zapsi_trans
	db.Raw("SELECT * FROM [zapsi_trans]  WHERE (job = '" + orderName + "') AND (qty_complete is not null) AND (trans_date > '" + terminalInputOrder.DTS.Format("2006-01-02 15:04:05") + "') AND (emp_num = '" + userInput + "')").Find(&zapsiTransThisOrder)
	logInfo("Check count", "Checking "+strconv.Itoa(len(zapsiTransThisOrder))+" transferred orders for "+orderName)
	for _, thisTrans := range zapsiTransThisOrder {
		transferredTotalThisOrder += int(thisTrans.QtyComplete)
	}
	return transferredTotalThisOrder
}

func GetCountForActualOpenOrder(workplaceCode string, userId string, orderInput string, operationInput string) (int, TerminalInputOrder) {
	order, suffix := ParseOrder(orderInput)
	operation := ParseOperation(operationInput)
	orderName := order + "." + suffix + "-" + operation

	var thisOrder TerminalInputOrder
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check operation", "Problem opening database: "+err.Error())
		return 0, thisOrder
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceCode).Find(&zapsiWorkplace)
	var zapsiOrder Order
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("OrderID = ?", zapsiOrder.OID).Where("UserID = ?", userId).Find(&thisOrder)
	return thisOrder.Count, thisOrder
}
