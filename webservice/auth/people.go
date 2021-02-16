package auth 

import (
    "net/http"
    "encoding/json"
    // "time"
    "fmt"
    // "log"
    // "os"
    "context"
    "github.com/jackc/pgx"
	"myvaksin/webservice/db"
)

const (
    DATE_ISO =  "2006-01-02"
)

type People struct {
	Name string 	`json:"name"`
    Ident string    `json:"ident"`
	Pwd string 		`json:"pwd"`
}

type SignUpHttpRespCode struct {
	SignUpRespCode string	`json:"signUpRespCode"`	
}

type SignInHttpRespCode struct {
	SignInRespCode string	`json:"signInRespCode"`
}

func SignUpPeople(conn *pgx.Conn, people People) error {
	sql :=
	    `insert into kkm.people
		(
			name, ident, password
		)
		values
		(
			$1, $2, $3
		)`
	
	_, err := conn.Exec(context.Background(), sql,
		people.Name, people.Ident, people.Pwd)
	if err != nil {
		return err
	}
	return nil 
}

func SignUpPeopleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	if (r.Method == "OPTIONS") { return }
    fmt.Println("[SignUpPeopleHandler] request received")
        
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        db.LogErrAndSendBadReqStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", people)

	db.CheckDbConn()
    err = SignUpPeople(db.Conn, people)
    if err != nil {
        db.LogErrAndSendBadReqStatus(w, err)
        return
    }  

	signUpRespCode := SignUpHttpRespCode {
		SignUpRespCode: "0",
	}
	signUpRespJson, err := json.MarshalIndent(signUpRespCode, "", "\t")
	if err != nil {
        db.LogErrAndSendBadReqStatus(w, err)
        return
    } 
	fmt.Fprintf(w, "%s", signUpRespJson)
}

func Bind(conn *pgx.Conn, people People) (string, error) {
	sql := 
		`select name from kkm.people
		 where ident=$1 and password=$2`

	row := conn.QueryRow(context.Background(), sql,
				people.Ident, people.Pwd)

	var dummy string				
	err := row.Scan(&dummy)				
	if err != nil {
	    if err == pgx.ErrNoRows {
		    return "0", nil 
		}
		return "0", err
	}  	   
	return "1", nil
}

func BindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	if (r.Method == "OPTIONS") { return }
    fmt.Println("[BindHandler] request received")
    
	// Decode    
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        db.LogErrAndSendBadReqStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", people)

	// Bind
	var bindResult string
	db.CheckDbConn()
    bindResult, err = Bind(db.Conn, people)
    if err != nil {
        db.LogErrAndSendBadReqStatus(w, err)
        return
    }  
	fmt.Printf("Bind status: %s\n", bindResult)

	// Encode
	signInRespCode := SignInHttpRespCode {
		SignInRespCode: bindResult,
	}
	signInRespJson, err := json.MarshalIndent(signInRespCode, "", "\t")
	if err != nil {
        db.LogErrAndSendBadReqStatus(w, err)
        return
    } 
	fmt.Fprintf(w, "%s", signInRespJson)
}