package util 

import 
(
	"net/http"
	"log"
)


func SendBadReqStatus(w http.ResponseWriter, err error) {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusBadRequest) //Http status code: 400
}

func SendUnauthorizedStatus(w http.ResponseWriter) {
    http.Error(w, "Unauthorized Access", http.StatusUnauthorized) //Http status code: 401
}

func SendInternalServerErrorStatus(w http.ResponseWriter, err error) {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusInternalServerError) //Http status code: 500
}