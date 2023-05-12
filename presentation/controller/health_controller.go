package controller

import (
	"goserver/utils/gsrender"
	"net/http"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	gsrender.Status(w, http.StatusOK)
}