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
	http.HandleFunc("/people/create", data.CreateNewPeopleHandler)
	http.HandleFunc("/people/get", data.GetPeopleHandler)
	http.HandleFunc("/people/update", data.UpdatePeopleHandler)
	http.HandleFunc("/people/add", data.AddPeopleHandler)
	http.HandleFunc("/people/delete", data.DeletePeopleHandler)
	http.HandleFunc("/vacrec/create", data.CreateNewVacRecHandler)
	http.HandleFunc("/vacrec/update", data.UpdateVacRecHandler)
	http.HandleFunc("/vacrec/delete", data.DeleteVacRecHandler)

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}