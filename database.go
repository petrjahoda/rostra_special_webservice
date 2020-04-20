package main

import (
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type SytelineUser struct {
	JePlatny string
	Jmeno    string
}

type SytelineOrder struct {
	CisloVp                string
	SuffixVp               string
	PolozkaVp              string
	PopisPolVp             string
	priznak_seriova_vyroba string
}

type SytelineOperation struct {
	pracoviste          string
	pracoviste_popis    string
	uvolneno_op         string
	priznak_mn_2        string
	mn_2_ks             string
	priznak_mn_3        string
	mn_3_ks             string
	jen_prenos_mnozstvi string
	priznak_nasobnost   string
	nasobnost           string
	parovy_dil          string
	seznamm_par_dilu    string
}

type SytelineWorkplace struct {
	Zapsi_zdroj           string
	priznak_mn_1          string
	vice_vp               string
	SL_prac               string
	auto_prevod_mnozstvi  string
	mnozstvi_auto_prevodu string
}

type TerminalInputOrder struct {
	OID      int       `gorm:"column:OID"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE; default:null"`
	Interval float32   `gorm:"column:Interval"`
	OrderID  int       `gorm:"column:OrderID"`
	UserID   int       `gorm:"column:UserID"`
	DeviceID int       `gorm:"column:DeviceID"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}

type TerminalInputIdle struct {
	OID      int       `gorm:"column:OID"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE; default:null"`
	IdleID   int       `gorm:"column:IdleID"`
	UserID   int       `gorm:"column:UserID"`
	Interval float32   `gorm:"column:Interval"`
	DeviceID int       `gorm:"column:DeviceID"`
}

func (TerminalInputIdle) TableName() string {
	return "terminal_input_idle"
}

type User struct {
	OID   int    `gorm:"column:OID"`
	Login string `gorm:"column:Login"`
	Name  string `gorm:"column:Name"`
}

func (User) TableName() string {
	return "user"
}

type Order struct {
	OID     int    `gorm:"column:OID"`
	Name    string `gorm:"column:Name"`
	Barcode string `gorm:"column:Barcode"`
}

func (Order) TableName() string {
	return "order"
}

type Workplace struct {
	OID                 int    `gorm:"column:OID"`
	Name                string `gorm:"column:Name"`
	WorkplaceDivisionId int    `gorm:"column:WorkplaceDivisionID"`
	DeviceID            int    `gorm:"column:DeviceID"`
}

func (Workplace) TableName() string {
	return "workplace"
}

type WorkplaceState struct {
	OID         int       `gorm:"column:OID"`
	StateID     int       `gorm:"column:StateID"`
	WorkplaceID int       `gorm:"column:WorkplaceID"`
	DTS         time.Time `gorm:"column:DTS"`
}

func (WorkplaceState) TableName() string {
	return "workplace_state"
}

func CheckDatabaseType() (string, string) {
	var connectionString string
	var dialect string
	if DatabaseType == "postgres" {
		connectionString = "host=" + DatabaseIpAddress + " sslmode=disable port=" + DatabasePort + " user=" + DatabaseLogin + " dbname=" + DatabaseName + " password=" + DatabasePassword
		dialect = "postgres"
	} else if DatabaseType == "mysql" {
		connectionString = DatabaseLogin + ":" + DatabasePassword + "@tcp(" + DatabaseIpAddress + ":" + DatabasePort + ")/" + DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
		dialect = "mysql"
	}
	return connectionString, dialect
}
