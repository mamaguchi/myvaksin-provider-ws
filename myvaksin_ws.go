package main

import (
	"net/http"
	"myvaksin/webservice/data"
)

func main() {
	defer data.Close()

	/* HANDLER FUNC */
	http.HandleFunc("/peoples", data.GetPeopleHandler)
	http.HandleFunc("/update/people", data.UpdatePeopleHandler)
	http.HandleFunc("/add/people", data.AddPeopleHandler)
	http.HandleFunc("/delete/people", data.DeletePeopleHandler)
	http.HandleFunc("/test", data.TestHandler)

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}