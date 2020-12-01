package main

import (
	"github.com/goodsign/monday"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"net/http"
	"os"
	"time"
)

const version = "2020.4.3.1"
const serviceName = "Rostra Special Web Service"
const serviceDescription = "Rostra Special Web Service"
const zapsiDatabaseConnection = "zapsi_uzivatel:zapsi@tcp(zapsidatabase:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
const sytelineDatabaseConnection = "sqlserver://zapsi:Zapsi_8513@192.168.1.26?database=rostra_exports"

type program struct{}

func (p *program) Start(s service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started at "+time.Now().String())
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
	timer := sse.New()

	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/mif/*filepath", http.Dir("mif"))

	router.GET("/", home)
	router.Handler("GET", "/time", timer)

	router.POST("/check_user_input", checkUserInput)
	router.POST("/check_order_input", checkOrderInput)
	router.POST("/check_operation_input", checkOperationInput)
	router.POST("/check_workplace_input", checkWorkplaceInput)
	router.POST("/check_count_input", checkCountInput)
	router.POST("/start_order", startOrder)
	router.POST("/transfer_order", transferOrder)
	router.POST("/end_order", endOrder)

	go streamTime(timer)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
}

func streamTime(streamer *sse.Streamer) {
	logInfo("SSE", "Streaming time process started")
	for {
		streamer.SendString("", "time", monday.Format(time.Now(), "Monday, 2. January 2006, 15:04:05", monday.LocaleCsCZ))
		time.Sleep(1 * time.Second)
	}
}
