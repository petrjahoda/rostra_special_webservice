package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

func UpdateZapsiZdrojFor(workplace SytelineWorkplace) string {
	LogInfo("MAIN", "Updating workplace name: "+workplace.Zapsi_zdroj)
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return workplace.Zapsi_zdroj
	}
	defer db.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplace.Zapsi_zdroj).Find(&zapsiWorkplace)
	LogInfo("MAIN", "Updated to: "+workplace.Zapsi_zdroj+";"+zapsiWorkplace.Name)
	return workplace.Zapsi_zdroj + ";" + zapsiWorkplace.Name
}

func EndOrderInZapsi(userid []string, orderId []string, operationid []string, workplaceid []string, ok []string, nok []string) bool {
	userLogin := strings.Split(userid[0], ";")[0]
	order, suffix := ParseOrder(orderId[0])
	operation := ParseOperation(operationid[0])
	orderName := order + "." + suffix + "-" + operation
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Login = ?", userLogin).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)

	if terminalInputOrder.OID > 0 {
		LogInfo("MAIN", "Closing order "+strconv.Itoa(terminalInputOrder.OID))
		db.Model(&terminalInputOrder).Where("OID = ?", terminalInputOrder.OID).UpdateColumn(TerminalInputOrder{DTE: sql.NullTime{Time: time.Now(), Valid: true}})
		db.Model(&terminalInputOrder).Where("OID = ?", terminalInputOrder.OID).UpdateColumn(TerminalInputOrder{Interval: float32(time.Now().Sub(terminalInputOrder.DTS).Seconds())})
		return true
	} else {
		return false
	}
}

func CalculateAverageCycle(ok []string, nok []string, interval int) (int, int, int) {
	cycle := 0
	okPcs := 0
	nokPcs := 0
	LogInfo("MAIN", "Calculating average cycle for "+ok[0]+" and "+nok[0]+" with interval of "+strconv.Itoa(interval))

	if len(ok) < 0 {
		okPcs, _ = strconv.Atoi(ok[0])
	}
	if len(nok) > 0 {
		nokPcs, _ = strconv.Atoi(ok[0])
	}
	totalPcs := okPcs + nokPcs
	if totalPcs > 0 {
		cycle = interval / totalPcs
	}
	return cycle, okPcs, nokPcs
}

func StartAndCloseOrderInZapsi(userid []string, orderid []string, operationid []string, workplaceid []string, ok []string, nok []string) bool {
	LogInfo("MAIN", "Starting order")
	sytelineOrder := GetOrderFromSyteline(orderid)
	sytelineOperation := GetOperationFromSyteline(orderid, operationid)
	CreateProductInZapsiIfNotExists(sytelineOrder)
	zapsiOrder := CreateOrderInZapsiIfNotExists(sytelineOrder, orderid, operationid, sytelineOperation, workplaceid)
	orderCreated := CreateAndCloseTerminalOrderInZapsi(userid, zapsiOrder, sytelineOperation, workplaceid, ok, nok)
	return orderCreated
}

func SaveNokIntoZapsi(nok []string, noktype []string, workplaceid []string, userid []string) {
	if len(nok) > 0 {
		LogInfo("MAIN", "Saving nok to Zapsi : "+noktype[0])
		CreateFailInZapsiIfNotExists(noktype)
		SaveTerminalInputFail(nok, noktype, workplaceid, userid)
	}
}

func SaveTerminalInputFail(nok []string, noktype []string, workplaceid []string, userid []string) {
	userLogin := strings.Split(userid[0], ";")[0]
	var zapsiFail Fail
	var terminalInputFail TerminalInputFail
	var zapsiWorkplace Workplace
	var zapsiUser User

	pcs, err := strconv.Atoi(nok[0])
	if err != nil {
		LogError("MAIN", "Problem parsing Nok amount when saving terminal input fail")
		return
	}
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", noktype[0]).Find(&zapsiFail)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("Login = ?", userLogin).Find(&zapsiUser)
	terminalInputFail.FailID = zapsiFail.OID
	terminalInputFail.DeviceID = zapsiWorkplace.DeviceID
	terminalInputFail.UserID = zapsiUser.OID
	terminalInputFail.DT = time.Now()
	for i := 0; i < pcs; i++ {
		db.Save(&terminalInputFail)
	}
	return
}

func CreateFailInZapsiIfNotExists(noktype []string) {
	nokTypes := GetNokTypesFromSyteline()
	for _, nokType := range nokTypes {
		if nokType.Nazev == noktype[0] {
			var zapsiFail Fail
			connectionString, dialect := CheckDatabaseType()
			db, err := gorm.Open(dialect, connectionString)

			if err != nil {
				LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
				return
			}
			defer db.Close()
			db.Where("Name = ?", noktype[0]).Find(&zapsiFail)
			if zapsiFail.OID > 0 {
				LogInfo("MAIN", "Fail "+noktype[0]+" already exists")
				return
			}
			LogInfo("MAIN", "Fail "+noktype[0]+" does not exist, creating fail")
			zapsiFail.Name = noktype[0]
			zapsiFail.Barcode = nokType.Kod
			zapsiFail.FailTypeID = 100
			db.Create(&zapsiFail)
			var newZapsiFail Fail
			db.Where("Name = ?", noktype[0]).Find(&newZapsiFail)

		}
	}
}

func CreateAndCloseTerminalOrderInZapsi(userid []string, zapsiOrder Order, sytelineOperation SytelineOperation, workplaceid []string, ok []string, nok []string) bool {
	userLogin := strings.Split(userid[0], ";")[0]
	var zapsiUser User
	var zapsiWorkplace Workplace
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Login = ?", userLogin).Find(&zapsiUser)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	var terminalInputOrder TerminalInputOrder
	defer db.Close()
	parsedCavity, err := strconv.Atoi(sytelineOperation.nasobnost)
	if err != nil {
		LogError("MAIN", "Problem parsing cavity: "+sytelineOperation.nasobnost)
		return false
	}

	if len(nok) > 0 {
		terminalInputOrder.Fail, _ = strconv.Atoi(nok[0])
	} else {
		terminalInputOrder.Fail = 0
	}
	if len(ok) > 0 {
		terminalInputOrder.Count, _ = strconv.Atoi(ok[0])
		terminalInputOrder.Count = terminalInputOrder.Count + terminalInputOrder.Fail
	} else {
		terminalInputOrder.Count = 0 + terminalInputOrder.Fail
	}
	terminalInputOrder.DTS = time.Now()
	terminalInputOrder.DTE = sql.NullTime{Time: time.Now(), Valid: true}
	terminalInputOrder.OrderID = zapsiOrder.OID
	terminalInputOrder.UserID = zapsiUser.OID
	terminalInputOrder.DeviceID = zapsiWorkplace.DeviceID
	terminalInputOrder.Interval = 0
	terminalInputOrder.AverageCycle = 0.0
	terminalInputOrder.WorkerCount = 1
	terminalInputOrder.WorkplaceModeID = 1
	terminalInputOrder.Cavity = parsedCavity
	db.Create(&terminalInputOrder)
	return true
}

func StartOrderInZapsi(userid []string, orderid []string, operationid []string, workplaceid []string, radio []string) bool {
	LogInfo("MAIN", "Starting order "+orderid[0])
	sytelineOrder := GetOrderFromSyteline(orderid)
	sytelineOperation := GetOperationFromSyteline(orderid, operationid)
	CreateProductInZapsiIfNotExists(sytelineOrder)
	zapsiOrder := CreateOrderInZapsiIfNotExists(sytelineOrder, orderid, operationid, sytelineOperation, workplaceid)
	orderCreated := CreateTerminalOrderInZapsi(userid, zapsiOrder, sytelineOperation, workplaceid, radio)
	return orderCreated
}

func CreateTerminalOrderInZapsi(userid []string, zapsiOrder Order, sytelineOperation SytelineOperation, workplaceid []string, radio []string) bool {
	userLogin := strings.Split(userid[0], ";")[0]
	parsedCavity, err := strconv.Atoi(sytelineOperation.nasobnost)
	if err != nil {
		LogError("MAIN", "Problem parsing cavity: "+sytelineOperation.nasobnost)
		return false
	}
	var zapsiUser User
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	var existingTerminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Login = ?", userLogin).Find(&zapsiUser)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceId = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserId is null").Find(&existingTerminalInputOrder)
	if existingTerminalInputOrder.OID > 0 {
		LogInfo("MAIN", "System terminal_input_order exists, just updating")
		db.Model(&terminalInputOrder).Where("DeviceId = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserId is null").Updates(map[string]interface{}{"OrderID": zapsiOrder.OID, "UserID": zapsiUser.OID, "Cavity": parsedCavity})
	} else {
		LogInfo("MAIN", "Creating new terminal_input_order")
		terminalInputOrder.DTS = time.Now()
		terminalInputOrder.OrderID = zapsiOrder.OID
		terminalInputOrder.UserID = zapsiUser.OID
		terminalInputOrder.DeviceID = zapsiWorkplace.DeviceID
		terminalInputOrder.Interval = 0
		terminalInputOrder.Count = 0
		terminalInputOrder.Fail = 0
		terminalInputOrder.AverageCycle = 0.0
		terminalInputOrder.WorkerCount = 1
		terminalInputOrder.WorkplaceModeID = 1
		terminalInputOrder.Cavity = parsedCavity
		terminalInputOrder.Note = radio[0]
		db.Create(&terminalInputOrder)
	}

	return true
}

func CreateOrderInZapsiIfNotExists(sytelineOrder SytelineOrder, orderid []string, operationid []string, sytelineOperation SytelineOperation, workplaceid []string) Order {
	var zapsiOrder Order
	var newOrder Order
	var zapsiProduct Product
	var zapsiWorkplace Workplace
	order, suffix := ParseOrder(orderid[0])
	operation := ParseOperation(operationid[0])
	zapsiOrderName := order + "." + suffix + "-" + operation
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return zapsiOrder
	}
	defer db.Close()
	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
	if zapsiOrder.OID > 0 {
		LogInfo("MAIN", "Order "+zapsiOrder.Name+" already exists")
		return zapsiOrder
	}
	LogInfo("MAIN", "Order "+zapsiOrder.Name+" does not exist, creating order in zapsi")
	db.Where("Name = ?", sytelineOrder.PolozkaVp).Find(&zapsiProduct)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	countRequestedConverted, err := strconv.ParseFloat(sytelineOperation.mn_2_ks, 32)
	if err != nil {
		LogError("MAIN", "Problem parsing count for sytelineOrder: "+sytelineOperation.mn_2_ks)
	}
	newOrder.Name = zapsiOrderName
	newOrder.Barcode = zapsiOrderName
	newOrder.ProductID = zapsiProduct.OID
	newOrder.OrderStatusID = 1
	newOrder.CountRequested = int(countRequestedConverted)
	newOrder.WorkplaceID = zapsiWorkplace.OID
	db.Create(&newOrder)
	db.Where("Name = ?", zapsiOrderName).Find(&zapsiOrder)
	return zapsiOrder
}

func CreateProductInZapsiIfNotExists(order SytelineOrder) {
	var zapsiProduct Product
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Name = ?", order.PolozkaVp).Find(&zapsiProduct)
	if zapsiProduct.OID > 0 {
		LogInfo("MAIN", "Product "+order.PolozkaVp+" already exists")
		return
	}
	LogInfo("MAIN", "Product "+order.PolozkaVp+" does not exist, creating product")
	zapsiProduct.Name = order.PolozkaVp
	zapsiProduct.Barcode = order.PolozkaVp
	zapsiProduct.Cycle = 1
	zapsiProduct.IdleFromTime = 1
	zapsiProduct.ProductGroupID = 1
	zapsiProduct.ProductStatusID = 1
	db.Create(&zapsiProduct)
}
func CheckIfOperatorAmountLessThanInZapsi(userAmount []string, userid []string, orderid []string, operationid []string, workplaceid []string) bool {
	userLogin := strings.Split(userid[0], ";")[0]
	order, suffix := ParseOrder(orderid[0])
	orderName := order + "." + suffix + "-" + operationid[0]
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Login like ?", userLogin).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
	okAmount, err := strconv.Atoi(userAmount[0])
	if err != nil {
		LogError("MAIN", "Problem parsing data from user")
		return false
	}
	if okAmount < terminalInputOrder.Count {
		return true
	}
	return false
}

func CheckIfAnyOpenOrderHasOneOfProducts(workplaceid []string, products []Product) bool {
	var terminalInputOrders []TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	var zapsiWorkplace Workplace
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Find(&terminalInputOrders)
	for _, terminalInputOrder := range terminalInputOrders {
		var zapsiOrder Order
		db.Where("OID = ?", terminalInputOrder.OrderID).Find(&zapsiOrder)
		for _, zapsiProduct := range products {
			if zapsiProduct.OID == zapsiOrder.ProductID {
				return true
			}
		}
	}
	return false
}

func CheckThisOpenOrderInZapsi(userid []string, orderid []string, operationid []string, workplaceid []string) (bool, string) {
	userLogin := strings.Split(userid[0], ";")[0]
	order, suffix := ParseOrder(orderid[0])
	operation := ParseOperation(operationid[0])
	orderName := order + "." + suffix + "-" + operation
	var zapsiUser User
	var zapsiOrder Order
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false, ""
	}
	defer db.Close()
	db.Where("Login = ?", userLogin).Find(&zapsiUser)
	db.Where("Name = ?", orderName).Find(&zapsiOrder)
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserID = ?", zapsiUser.OID).Where("OrderID = ?", zapsiOrder.OID).Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		return true, terminalInputOrder.Note
	}
	return false, ""
}

func CheckAnyOpenOrderInZapsi(workplaceid []string) bool {
	var zapsiWorkplace Workplace
	var terminalInputOrder TerminalInputOrder
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	db.Where("Code = ?", workplaceid[0]).Find(&zapsiWorkplace)
	db.Where("DeviceID = ?", zapsiWorkplace.DeviceID).Where("DTE is null").Where("UserId is not null").Find(&terminalInputOrder)
	if terminalInputOrder.OID > 0 {
		return true
	}
	return false
}

func CheckProductsInZapsi(operation SytelineOperation) []Product {
	var zapsiProducts []Product
	var products []string
	if strings.Contains(operation.seznamm_par_dilu, "|") {
		products = strings.Split(operation.seznamm_par_dilu, "|")
	} else {
		products = append(products, operation.seznamm_par_dilu)
	}

	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return zapsiProducts
	}
	defer db.Close()
	for _, product := range products {
		var zapsiProduct Product
		db.Where("Name = ?", product).Find(&zapsiProduct)
		if zapsiProduct.OID > 0 {
			LogInfo("MAIN", "Product "+product+" already exists")
		} else {
			LogInfo("MAIN", "Product "+product+" does not exist, creating product")
			zapsiProduct.Name = product
			zapsiProduct.Barcode = product
			zapsiProduct.Cycle = 1
			zapsiProduct.IdleFromTime = 1
			zapsiProduct.ProductGroupID = 1
			zapsiProduct.ProductStatusID = 1
			db.Create(&zapsiProduct)
		}
	}
	for _, product := range products {
		var zapsiProduct Product
		db.Where("Name = ?", product).Find(&zapsiProduct)
		zapsiProducts = append(zapsiProducts, zapsiProduct)
	}
	return zapsiProducts
}

func CreateUserInZapsiIfNotExists(user SytelineUser, userid []string) {
	userLogin := userid[0]
	userFirstName := strings.Split(user.Jmeno, ",")[0]
	userSecondName := strings.Split(user.Jmeno, ",")[1]
	var zapsiUser User
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	db.Where("Login LIKE ?", userLogin).Find(&zapsiUser)
	if zapsiUser.OID > 0 {
		LogInfo("MAIN", "User "+user.Jmeno+"already exists")
		return
	}
	LogInfo("MAIN", "User "+user.Jmeno+" does not exist, creating user "+user.Jmeno)
	zapsiUser.Login = userLogin
	zapsiUser.FirstName = userFirstName
	zapsiUser.Name = userSecondName
	zapsiUser.UserRoleID = "1"
	zapsiUser.UserTypeID = "1"
	db.Create(&zapsiUser)
}
