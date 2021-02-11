package main

import (
	"net/http"
	"myvaksin/webservice/data"
)

func main() {
	defer data.Close()

	/* HANDLER FUNC */
	http.HandleFunc("/people", data.GetPeopleHandler)
	http.HandleFunc("/people/update", data.UpdatePeopleHandler)
	http.HandleFunc("/people/add", data.AddPeopleHandler)
	http.HandleFunc("/people/delete", data.DeletePeopleHandler)
	http.HandleFunc("/test", data.TestHandler)

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}