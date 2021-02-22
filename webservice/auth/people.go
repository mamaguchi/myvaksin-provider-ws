package auth 

import (
    "net/http"
    "encoding/json"
    "fmt"
    "context"
    "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
	"myvaksin/webservice/db"
)

const (
	DEFAULT_PEOPLE_PWD = "myvaksin"
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
	Name string 			`json:"name"`
	Ident string 			`json:"ident"`
	Role string 			`json:"role"`
}

func SignUpPeople(conn *pgxpool.Pool, people People) (string, error) {
	sqlSelect := 
		`select name, password from kkm.people
		 where ident=$1`

	row := conn.QueryRow(context.Background(), sqlSelect,
				people.Ident)
	var name string				
	var password string
	err := row.Scan(&name, &password)				
	if err != nil {
		// People Ident doesn't exist, 
		// so can sign up a new account.
	    if err == pgx.ErrNoRows { 
			sqlInsert :=
				`insert into kkm.people
				(
					name, ident, password, role
				)
				values
				(
					$1, $2, $3, $4
				)`
			
			_, err = conn.Exec(context.Background(), sqlInsert,
				people.Name, people.Ident, people.Pwd, "receiver")
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

	// A People profile has been created by myVaksin provider
	// before the account has been registered.
	// 'myvaksin' is the default password inserted when 
	// creating a new People profile for a person who 
	// has not registered an account.
	// (Got profile, No account)
	if password == "myvaksin" {
		sqlUpdate := `update kkm.people
				  set name=$1, password=$2
				  where ident=$3`
		_, err = conn.Exec(context.Background(), sqlUpdate,
			people.Name, people.Pwd, people.Ident)
		if err != nil {
			return "", err
		}
		return "1", nil
	} else {
		// People Ident already exists in the table, 
		// and the password is not the default one.
		// So unable to sign up a new account.
		// (Got account, Got profile)
		return "0", nil
	}


	
}

func SignUpPeopleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")	
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
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
func Bind(conn *pgxpool.Pool, people People) (bool, error, string, string) {
	sql := 
		`select name, role from kkm.people
		 where ident=$1 and password=$2`

	row := conn.QueryRow(context.Background(), sql,
				people.Ident, people.Pwd)

	var name string	
	var role string
	err := row.Scan(&name, &role)				
	if err != nil {
	    if err == pgx.ErrNoRows {
		    return false, nil, "", "" 
		}
		return false, err, "", ""
	}  	   
	return true, nil, name, role
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
	var name string 
	var role string
	db.CheckDbConn()
    bindResult, err, name, role = Bind(db.Conn, people)
    if err != nil {
		db.LogErrAndSendBadReqStatus(w, err)
        return
	}  
	fmt.Printf("Bind status: %v\n", bindResult)
	if !bindResult {
        db.SendUnauthorizedStatus(w)
        return
    }  
	tokenString, err := NewTokenHMAC(people.Ident)
	if err != nil {
		db.LogErrAndSendInternalServerErrorStatus(w, err)
        return
	}

	// Encode	
	authResult := SignInAuthResult{
		Token: tokenString,
		Name: name,
		Ident: people.Ident,
		Role: role,
	}
	authResultJson, err := json.MarshalIndent(&authResult, "", "")
	if err != nil {
		db.LogErrAndSendInternalServerErrorStatus(w, err)
        return
	}
	fmt.Fprintf(w, "%s", authResultJson)
}