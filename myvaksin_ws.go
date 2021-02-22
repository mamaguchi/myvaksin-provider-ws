package main

import (
	"net/http"
	"myvaksin/webservice/db"
	"myvaksin/webservice/test"
	"myvaksin/webservice/data"
	"myvaksin/webservice/auth"
)

func main() {
	/* INIT DATABASE CONNECTION */
	// defer data.Close()	
	db.Open()
	defer db.Close()

	/* HANDLER FUNC */
	// Test
	http.HandleFunc("/test", test.TestGetPeopleHandler)
	// Auth
	http.HandleFunc("/signup", auth.SignUpPeopleHandler)
	http.HandleFunc("/signin", auth.BindHandler)
	// People
	http.HandleFunc("/people/search", data.SearchPeopleHandler)
	http.HandleFunc("/people/create", data.CreateNewPeopleHandler)
	http.HandleFunc("/people/get", data.GetPeopleHandler)
	http.HandleFunc("/people/update", data.UpdatePeopleHandler)
	http.HandleFunc("/people/delete", data.DeletePeopleHandler)
	http.HandleFunc("/vacrec/create", data.CreateNewVacRecHandler)
	http.HandleFunc("/vacrec/update", data.UpdateVacRecHandler)
	http.HandleFunc("/vacrec/delete", data.DeleteVacRecHandler)

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}