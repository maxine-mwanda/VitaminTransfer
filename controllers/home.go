package handlers

import (
	"html/template"
	"net/http"
	"vitamin-transfer/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	utils.Logger.Info("Home page accessed")
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		utils.Logger.Error("Error parsing template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		utils.Logger.Error("Error executing template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}