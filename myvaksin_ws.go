package main

import (
	"net/http"
	"myvaksin/webservice/data"
)

func main() {
	defer data.Close()

	/* HANDLER FUNC */
	http.HandleFunc("/test", data.TestHandler)
	http.HandleFunc("/people/search", data.SearchPeopleHandler)
	http.HandleFunc("/people/get", data.GetPeopleHandler)
	http.HandleFunc("/people/update", data.UpdatePeopleHandler)
	http.HandleFunc("/people/add", data.AddPeopleHandler)
	http.HandleFunc("/people/delete", data.DeletePeopleHandler)
	http.HandleFunc("/vacrec/insert", data.InsertNewVacRecHandler)
	http.HandleFunc("/vacrec/update", data.UpdateVacRecHandler)

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}