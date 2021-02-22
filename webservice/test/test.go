package test 

import (
    "net/http"
    "encoding/json"
    "fmt"
    "log"
    "context"
    // "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
	"myvaksin/webservice/db"
)

type Identity struct {
    Ident string    `json:"ident"`
}

func TestGetPeople(conn *pgxpool.Pool, ident string) error {
    row := conn.QueryRow(context.Background(), 
        `select kkm.people.name
          from kkm.people             
          where ident=$1`,
        ident)
    var name string    
    err := row.Scan(&name)
    if err != nil {
        return err
    }  
	fmt.Printf("[TestGetPeople] Name for ident %s: %s\n", 
				ident, name)  
    return nil
}

func TestGetPeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[TestGetPeopleHandler] Request received")

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%+v\n", identity)    

	db.CheckDbConn()
	err = TestGetPeople(db.Conn, identity.Ident)
	if err != nil {
		log.Print(err)
	}
}