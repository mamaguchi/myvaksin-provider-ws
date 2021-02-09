package main

import (
	"net/http"
	// "encoding/json"
	// "fmt"
	"myVaksin/webservice/data"
)

func main() {
	defer data.Close()

	/* HANDLER FUNC */
	http.HandleFunc("/peoples", data.GetPeopleHandler)
	http.HandleFunc("/update/people", data.UpdatePeopleHandler)

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}