package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"net/http"
	"time"
)

const version = "2020.3.1.14"
const programName = "Rostra Special Web Service"
const programDescription = "Allow users to start and end orders"
const deleteLogsAfter = 240 * time.Hour

type program struct{}

func (p *program) Start(s service.Service) error {
	LogInfo("MAIN", "Starting "+programName+" on "+s.Platform())
	go p.run()
	return nil
}

func (p *program) run() {
	LogDirectoryFileCheck("MAIN")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	router := httprouter.New()
	router.GET("/", RostraMainScreen)
	router.GET("/data_input", DataInput)
	router.GET("/js/metro.min.js", metrojs)
	router.GET("/css/metro-all.css", metrocss)
	router.GET("/mif/metro.ttf", metrottf)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":80", router)
}

func main() {
	serviceConfig := &service.Config{
		Name:        programName,
		DisplayName: programName,
		Description: programDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		LogError("MAIN", err.Error())
	}
	err = s.Run()
	if err != nil {
		LogError("MAIN", "Problem starting "+serviceConfig.Name)
	}
}

func (p *program) Stop(s service.Service) error {
	LogInfo("MAIN", "Stopped on platform "+s.Platform())
	return nil
}

func metrojs(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "js/metro.min.js")
}

func metrocss(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "css/metro-all.css")
}

func metrottf(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "mif/metro.ttf")
}
