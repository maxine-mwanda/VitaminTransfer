package controllers

import (
	"html/template"
	"net/http"
	"VitaminTransfer/utils"
)

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	utils.Logger.Info("Success page accessed")
	tmpl, err := template.ParseFiles("templates/success.html")
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