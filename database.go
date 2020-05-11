package main

import (
	"database/sql"
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

type SytelineNok struct {
	Kod   string
	Nazev string
}

type TerminalInputOrder struct {
	OID             int           `gorm:"column:OID"`
	DTS             time.Time     `gorm:"column:DTS"`
	DTE             sql.NullTime  `gorm:"column:DTE; default:null"`
	OrderID         int           `gorm:"column:OrderID"`
	UserID          int           `gorm:"column:UserID"`
	DeviceID        int           `gorm:"column:DeviceID"`
	Interval        float32       `gorm:"column:Interval"`
	Count           int           `gorm:"column:Count"`
	Fail            int           `gorm:"column:Fail"`
	AverageCycle    float32       `gorm:"column:AverageCycle"`
	WorkerCount     int           `gorm:"column:WorkerCount"`
	WorkplaceModeID int           `gorm:"column:WorkplaceModeID"`
	Note            string        `gorm:"column:Note"`
	WorkshiftID     sql.NullInt32 `gorm:"column:WorkshiftID"`
	Cavity          int           `gorm:"column:Cavity"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}

type User struct {
	OID        int    `gorm:"column:OID"`
	Login      string `gorm:"column:Login"`
	Name       string `gorm:"column:Name"`
	FirstName  string `gorm:"column:FirstName"`
	UserTypeID string `gorm:"column:UserTypeID"`
	UserRoleID string `gorm:"column:UserRoleID"`
}

func (User) TableName() string {
	return "user"
}

type Order struct {
	OID            int    `gorm:"column:OID"`
	Name           string `gorm:"column:Name"`
	Barcode        string `gorm:"column:Barcode"`
	ProductID      int    `gorm:"column:ProductID"`
	OrderStatusID  int    `gorm:"column:OrderStatusID"`
	CountRequested int    `gorm:"column:CountRequested"`
	WorkplaceID    int    `gorm:"column:WorkplaceID"`
	Cavity         int    `gorm:"column:Cavity"`
}

func (Order) TableName() string {
	return "order"
}

type Product struct {
	OID             int     `gorm:"column:OID"`
	Name            string  `gorm:"column:Name"`
	Barcode         string  `gorm:"column:Barcode"`
	Cycle           float32 `gorm:"column:Cycle"`
	IdleFromTime    int     `gorm:"column:IdleFromTime"`
	ProductStatusID int     `gorm:"column:ProductStatusID"`
	Deleted         int     `gorm:"column:Deleted"`
	ProductGroupID  int     `gorm:"column:ProductGroupID"`
}

func (Product) TableName() string {
	return "product"
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
