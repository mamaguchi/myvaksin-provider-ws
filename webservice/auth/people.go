package auth 

import (
    "net/http"
    "encoding/json"
    "fmt"
    "context"
    "github.com/jackc/pgx"
	"myvaksin/webservice/db"
)

type People struct {
	Name string 	`json:"name"`
    Ident string    `json:"ident"`
	Pwd string 		`json:"pwd"`
}

type SignUpHttpRespCode struct {
	SignUpRespCode string	`json:"signUpRespCode"`	
}

// This struct is for debugging during development
type SignInHttpRespCode struct {
	SignInRespCode string	`json:"signInRespCode"`
}

type SignInAuthResult struct {
	Token string			`json:"token"`
}

func SignUpPeople(conn *pgx.Conn, people People) (string, error) {
	sqlSelect := 
		`select name from kkm.people
		 where ident=$1`

	row := conn.QueryRow(context.Background(), sqlSelect,
				people.Ident)
	var dummy string				
	err := row.Scan(&dummy)				
	if err != nil {
		// People Ident doesn't exist, 
		// so can sign up a new account.
	    if err == pgx.ErrNoRows { 
			sqlInsert :=
				`insert into kkm.people
				(
					name, ident, password
				)
				values
				(
					$1, $2, $3
				)`
			
			_, err = conn.Exec(context.Background(), sqlInsert,
				people.Name, people.Ident, people.Pwd)
			if err != nil {
				// New account create failed.
				return "", err
			}
			// New account created successfully.
			return "1", nil
		} 
		// Other unknown error during database scan.
		return "", err
	} 
	// People Ident already exists in the table, 
	// so unable to sign up a new account.
	return "0", nil


	// sql :=
	//     `insert into kkm.people
	// 	(
	// 		name, ident, password
	// 	)
	// 	values
	// 	(
	// 		$1, $2, $3
	// 	)`
	
	// _, err := conn.Exec(context.Background(), sql,
	// 	people.Name, people.Ident, people.Pwd)
	// if err != nil {
	// 	return err
	// }
	// return nil 
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
    signUpResult, err := SignUpPeople(db.Conn, people)
    if err != nil {
        db.LogErrAndSendInternalServerErrorStatus(w, err)
        return
    }  

	signUpRespCode := SignUpHttpRespCode {
		SignUpRespCode: signUpResult,
	}
	signUpRespJson, err := json.MarshalIndent(signUpRespCode, "", "")
	if err != nil {
        db.LogErrAndSendInternalServerErrorStatus(w, err)
        return
    } 
	fmt.Fprintf(w, "%s", signUpRespJson)
}

// Bind == SignIn 
func Bind(conn *pgx.Conn, people People) (bool, error) {
	sql := 
		`select name from kkm.people
		 where ident=$1 and password=$2`

	row := conn.QueryRow(context.Background(), sql,
				people.Ident, people.Pwd)

	var dummy string				
	err := row.Scan(&dummy)				
	if err != nil {
	    if err == pgx.ErrNoRows {
		    return false, nil 
		}
		return false, err
	}  	   
	return true, nil
}

func BindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")	
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
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
	var bindResult bool
	db.CheckDbConn()
    bindResult, err = Bind(db.Conn, people)
    if err != nil {
		db.LogErrAndSendBadReqStatus(w, err)
        return
	}  
	fmt.Printf("Bind status: %v\n", bindResult)
	if !bindResult {
        db.LogErrAndSendUnauthorizedStatus(w, err)
        return
    }  
	tokenString, err := NewTokenHMAC(people.Ident)
	if err != nil {
		db.LogErrAndSendInternalServerErrorStatus(w, err)
        return
	}

	// Encode
	// signInRespCode := SignInHttpRespCode {
	// 	SignInRespCode: bindResult,
	// }
	// signInRespJson, err := json.MarshalIndent(signInRespCode, "", "\t")
	// if err != nil {
    //     db.LogErrAndSendBadReqStatus(w, err)
    //     return
    // } 
	// fmt.Fprintf(w, "%s", signInRespJson)
	authResult := SignInAuthResult{
		Token: tokenString,
	}
	authResultJson, err := json.MarshalIndent(&authResult, "", "")
	if err != nil {
		db.LogErrAndSendInternalServerErrorStatus(w, err)
        return
	}
	fmt.Fprintf(w, "%s", authResultJson)
}