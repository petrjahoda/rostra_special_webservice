package main

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type OperationList struct {
	Operce          string `gorm:"column:operce"`
	Pracoviste      string `gorm:"column:pracoviste"`
	PracovistePopis string `gorm:"column:pracoviste_popis"`
}

type zapsi_trans struct {
	Trans_num       float64      `gorm:"column:trans_num"`
	Posted          int          `gorm:"column:posted"`
	Trans_date      sql.NullTime `gorm:"column:trans_date"`
	Emp_num         string       `gorm:"column:emp_num"`
	Trans_type      string       `gorm:"column:trans_type"`
	Job             string       `gorm:"column:job"`
	Suffix          int          `gorm:"column:suffix"`
	Oper_num        int          `gorm:"column:oper_num"`
	Wc              string       `gorm:"column:wc"`
	Qty_complete    float64      `gorm:"column:qty_complete"`
	Qty_scrapped    float64      `gorm:"column:qty_scrapped"`
	Lot             string       `gorm:"column:lot"`
	Start_date_time sql.NullTime `gorm:"column:start_date_time"`
	End_date_time   sql.NullTime `gorm:"column:end_date_time"`
	Complete_op     int          `gorm:"column:complete_op"`
	Shift           string       `gorm:"column:shift"`
	Reason_code     string       `gorm:"column:reacon_code"`
	Time_divisor    float64      `gorm:"column:time_divisor"`
}

func (zapsi_trans) TableName() string {
	return "zapsi_trans"
}

type SytelineUser struct {
	JePlatny string         `gorm:"column:JePlatny"`
	Jmeno    sql.NullString `gorm:"column:Jmeno"`
	Chyba    sql.NullString `gorm:"column:Chyba"`
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
	seznamm_par_dilu    sql.NullString
}

type SytelineWorkplace struct {
	Zapsi_zdroj           string
	priznak_mn_1          string
	vice_vp               string
	SL_prac               string
	typ_zdroje_zapsi      string
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

type TerminalInputFail struct {
	OID      int       `gorm:"column:OID"`
	DT       time.Time `gorm:"column:DT"`
	FailID   int       `gorm:"column:FailID"`
	UserID   int       `gorm:"column:UserID"`
	DeviceID int       `gorm:"column:DeviceID"`
	Note     string    `gorm:"column:Note"`
}

func (TerminalInputFail) TableName() string {
	return "terminal_input_fail"
}

type Fail struct {
	OID        int    `gorm:"column:OID"`
	Name       string `gorm:"column:Name"`
	Barcode    string `gorm:"column:Barcode"`
	FailTypeID int    `gorm:"column:FailTypeID"`
}

func (Fail) TableName() string {
	return "fail"
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

type Device struct {
	OID     int    `gorm:"column:OID"`
	Setting string `gorm:"column:Setting"`
}

func (Device) TableName() string {
	return "device"
}
