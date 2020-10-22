package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strings"
)

type HomePageData struct {
	Version string
}

func home(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	logInfo(ipAddress[0], "Sending home page")
	var data HomePageData
	data.Version = version
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/home.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(ipAddress[0], "Home page sent")
}
