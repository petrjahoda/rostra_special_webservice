package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

const version = "2020.2.1.14"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	router := httprouter.New()
	//router.GET("/rostra_main_screen", RostraMainScreen)

	router.GET("/", RostraMainScreen)
	router.GET("/reset", RostraMainScreen)
	router.GET("/data_input", DataInput)

	router.GET("/js/metro.min.js", metrojs)
	router.GET("/css/metro-all.css", metrocss)
	router.GET("/mif/metro.ttf", metrottf)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":80", router)
}

func metrojs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "js/metro.min.js")
}

func metrocss(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "css/metro-all.css")
}

func metrottf(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "mif/metro.ttf")
}
