package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type UserInputData struct {
	UserInput string
}

type UserResponseData struct {
	Result    string
	UserInput string
	UserName  string
	UserId    string
	UserError string
}

func checkUserInput(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check user", "Started")
	var data UserInputData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check user", "Error parsing input: "+err.Error())
		var responseData UserResponseData
		responseData.Result = "nok"
		responseData.UserInput = data.UserInput
		responseData.UserError = "Problem parsing input: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended with error")
		return
	}
	logInfo("Check user", "Data: "+data.UserInput)
	db, err := gorm.Open(sqlserver.Open(sytelineDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check user", "Problem opening database: "+err.Error())
		var responseData UserResponseData
		responseData.Result = "nok"
		responseData.UserInput = data.UserInput
		responseData.UserError = "Problem connecting Syteline database: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended with error")
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var sytelineUser SytelineUser
	command := "declare @Zamestnanec EmpNumType, @JePlatny ListYesNoType, @Jmeno NameType, @Chyba Infobar  Exec [rostra_exports_test].dbo.ZapsiKontrolaZamSp @Zamestnanec = N'" + data.UserInput + "', @JePlatny = @JePlatny output, @Jmeno = @Jmeno output, @Chyba = @Chyba output select JePlatny = @JePlatny, Jmeno = @Jmeno, Chyba = @Chyba;\n"
	db.Raw(command).Scan(&sytelineUser)
	if sytelineUser.JePlatny == "1" {
		logInfo("Check user", "User found: "+data.UserInput)
		userId := CreateUserInZapsiIfNotExists(sytelineUser, data.UserInput)
		var responseData UserResponseData
		responseData.Result = "ok"
		responseData.UserInput = data.UserInput
		responseData.UserId = strconv.Itoa(userId)
		responseData.UserName = sytelineUser.Jmeno.String
		responseData.UserError = sytelineUser.Chyba.String
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended successfully")
		return
	} else {
		logInfo("Check user", "User not found: "+sytelineUser.Chyba.String)
		var responseData UserResponseData
		responseData.Result = "nok"
		responseData.UserInput = data.UserInput
		responseData.UserError = sytelineUser.Chyba.String
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check user", "Ended successfully with no user found")
		return
	}
}

func CreateUserInZapsiIfNotExists(user SytelineUser, input string) int {
	logInfo("Check user", "Checking user in Zapsi")
	userFirstName := strings.Split(user.Jmeno.String, ",")[0]
	userSecondName := strings.Split(user.Jmeno.String, ",")[1]
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("Check user", "Problem opening database: "+err.Error())
		return 0
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiUser User
	db.Where("Login LIKE ?", input).Find(&zapsiUser)
	if zapsiUser.OID > 0 {
		logInfo("Check user", "User "+user.Jmeno.String+"already exists")
		return zapsiUser.OID
	}
	logInfo("Check user", "User "+user.Jmeno.String+" does not exist, creating user "+user.Jmeno.String)
	zapsiUser.Login = input
	zapsiUser.FirstName = userFirstName
	zapsiUser.Name = userSecondName
	zapsiUser.UserRoleID = "1"
	zapsiUser.UserTypeID = "1"
	db.Create(&zapsiUser)
	logInfo("Check user", "User "+user.Jmeno.String+" created")
	var newUser User
	db.Where("Login LIKE ?", input).Find(&newUser)
	return newUser.OID
}
