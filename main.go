package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"net/http"
	"os"
)

const version = "2020.4.1.19"
const serviceName = "Rostra Special Web Service"
const serviceDescription = "Rostra Special Web Service"
const zapsiDatabaseConnection = "zapsi_uzivatel:zapsi@tcp(zapsidatabase:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
const sytelineDatabaseConnection = "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports_test"

type program struct{}

func (p *program) Start(s service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) run() {
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/mif/*filepath", http.Dir("mif"))

	router.GET("/", home)

	router.POST("/check_user_input", checkUserInput)
	router.POST("/check_order_input", checkOrderInput)
	router.POST("/check_operation_input", checkOperationInput)
	router.POST("/check_workplace_input", checkWorkplaceInput)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
}
