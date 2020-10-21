package main

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type OperationList struct {
	Operace         string `gorm:"column:operce"`
	Pracoviste      string `gorm:"column:pracoviste"`
	PracovistePopis string `gorm:"column:pracoviste_popis"`
}

type zapsi_trans struct {
	TransNum      float64      `gorm:"column:trans_num"`
	Posted        int          `gorm:"column:posted"`
	TransDate     sql.NullTime `gorm:"column:trans_date"`
	EmpNum        string       `gorm:"column:emp_num"`
	TransType     string       `gorm:"column:trans_type"`
	Job           string       `gorm:"column:job"`
	Suffix        int          `gorm:"column:suffix"`
	OperNum       int          `gorm:"column:oper_num"`
	Wc            string       `gorm:"column:wc"`
	QtyComplete   float64      `gorm:"column:qty_complete"`
	QtyScrapped   float64      `gorm:"column:qty_scrapped"`
	Lot           string       `gorm:"column:lot"`
	StartDateTime sql.NullTime `gorm:"column:start_date_time"`
	EndDateTime   sql.NullTime `gorm:"column:end_date_time"`
	CompleteOp    int          `gorm:"column:complete_op"`
	Shift         string       `gorm:"column:shift"`
	ReasonCode    string       `gorm:"column:reacon_code"`
	TimeDivisor   float64      `gorm:"column:time_divisor"`
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
	CisloVp              string `gorm:"column:CisloVp"`
	SuffixVp             string `gorm:"column:SuffixVp"`
	PolozkaVp            string `gorm:"column:PolozkaVp"`
	PopisPolVp           string `gorm:"column:PopisPolVp"`
	PriznakSeriovaVyroba string `gorm:"column:priznak_seriova_vyroba"`
}

type SytelineOperation struct {
	Pracoviste        string         `gorm:"column:pracoviste"`
	PracovistePopis   string         `gorm:"column:pracoviste_popis"`
	UvolnenoOp        string         `gorm:"column:uvolneno_op"`
	PriznakMn2        string         `gorm:"column:priznak_mn_2"`
	Mn2Ks             string         `gorm:"column:mn_2_ks"`
	PriznakMn3        string         `gorm:"column:priznak_mn_3"`
	Mn3Ks             string         `gorm:"column:mn_3_ks"`
	JenPrenosMnozstvi string         `gorm:"column:jen_prenos_mnozstvi"`
	PriznakNasobnost  string         `gorm:"column:priznak_nasobnost"`
	Nasobnost         string         `gorm:"column:nasobnost"`
	ParovyDil         string         `gorm:"column:parovy_dil"`
	SeznamParDilu     sql.NullString `gorm:"column:seznamm_par_dilu"`
}

type SytelineWorkplace struct {
	ZapsiZdroj          string `gorm:"column:zapsi_zdroj"`
	PriznakMn1          string `gorm:"column:priznak_mn_1"`
	ViceVp              string `gorm:"column:vice_vp"`
	SlPrac              string `gorm:"column:SL_prac"`
	TypZdrojeZapsi      string `gorm:"column:typ_zdroje_zapsi"`
	AutoPrevodMnozstvi  string `gorm:"column:auto_prevod_mnozstvi"`
	MnozstviAutoPrevodu string `gorm:"column:auto_prevod_mnozstvi"`
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
	Code                string `gorm:"column:Code"`
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
