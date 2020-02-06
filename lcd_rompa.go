package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type LcdWorkplaces struct {
	LcdWorkplaces []LcdWorkplace
	Version       string
}
type LcdWorkplace struct {
	Name       string
	User       string
	StateColor string
	Duration   time.Duration
	InforData  string
}

func RostraMainScreen(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Displaying LCD for Rompa")
	tmpl := template.Must(template.ParseFiles("html/rostra_main_screen.html"))
	lcdWorkplaces := LcdWorkplaces{}
	lcdWorkplaces.Version = "version: " + version
	_ = tmpl.Execute(writer, lcdWorkplaces)
}

func StartOrder(writer http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Displaying LCD for Rompa")
	tmpl := template.Must(template.ParseFiles("html/start_order.html"))
	lcdWorkplaces := LcdWorkplaces{}
	lcdWorkplaces.Version = "version: " + version
	_ = tmpl.Execute(writer, lcdWorkplaces)
}

