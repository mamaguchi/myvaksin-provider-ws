package data

import (
    "net/http"
    "encoding/json"
    "time"
    "fmt"
    "os"
    "log"
    "context"
    // "strings"
    // "strconv"
    "github.com/jackc/pgx"
)

const (
    DATE_ISO =  "2006-01-02"
)

var conn *pgx.Conn

func init() {
    var err error
    conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to make connection to database: %v\n", err)
		os.Exit(1)
    }     
}

func Close() {
    conn.Close(context.Background())
}

type Peoples struct {
    Peoples []People    `json:"peoples"`
}

type People struct {
    Ident string          `json:"ident"`
    Name string           `json:"name"`
    Dob time.Time         `json:"dob"`
    Tel string            `json:"tel"`
    Address string        `json:"address"`
    Race string           `json:"race"`
    Nationality string    `json:"nationality"`
    Edu_lvl string        `json:"edu_lvl"`
    Occupation string     `json:"occupation"`
    Comorbids []int       `json:"comorbids"`
    Support_vac bool      `json:"support_vac"`
}

func GetPeoples(conn *pgx.Conn) ([]byte, error) {
    var peoples Peoples
    rows, _ := conn.Query(context.Background(), 
        "select ident, name, dob, tel, address, race, nationality, edu_lvl, occupation, comorbids, support_vac from kkm.people")

    for rows.Next() {
        var ident string
        var name string
        var dob time.Time
        var tel string
        var address string
        var race string
        var nationality string
        var edu_lvl string
        var occupation string
        var comorbids []int
        var support_vac bool
        err := rows.Scan(&ident, &name, &dob, &tel, &address, &race,
                        &nationality, &edu_lvl, &occupation, &comorbids, 
                        &support_vac)
        if err != nil {
            return nil, err
        }
        people := People{
            Ident: ident,
            Name: name,
            Dob: dob,
            Tel: tel,
            Address: address,
            Race: race,
            Nationality: nationality,
            Edu_lvl: edu_lvl,
            Occupation: occupation,
            Comorbids: comorbids,
            Support_vac: support_vac,
        }
        peoples.Peoples = append(peoples.Peoples, people)
    }
    output, err := json.MarshalIndent(peoples, "", "\t")
        
    return output, err
}

func GetPeople(conn *pgx.Conn, ident string) ([]byte, error) {
    row := conn.QueryRow(context.Background(), 
        `select name, dob, tel, address, race, nationality, 
           edu_lvl, occupation, comorbids, support_vac 
         from kkm.people 
         where ident=$1`,
        ident)
    
    var name string
    var dob time.Time
    var tel string
    var address string
    var race string
    var nationality string
    var edu_lvl string
    var occupation string
    var comorbids []int
    var support_vac bool
    err := row.Scan(&name, &dob, &tel, &address, &race,
                    &nationality, &edu_lvl, &occupation, &comorbids, 
                    &support_vac)
    if err != nil {
        return nil, err
    }
    people := People{
        Ident: ident,
        Name: name,
        Dob: dob,
        Tel: tel,
        Address: address,
        Race: race,
        Nationality: nationality,
        Edu_lvl: edu_lvl,
        Occupation: occupation,
        Comorbids: comorbids,
        Support_vac: support_vac,
    }    
    output, err := json.MarshalIndent(people, "", "\t")
        
    return output, err
}


func GetPeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type") 
    if (r.Method == "OPTIONS") { return }
    r.ParseForm()
    fmt.Println("[GetPeopleHandler] Request form data received")
    fmt.Println(r.Form)

    peopleIdent := r.Form["ident"][0]
    peopleJson, err := GetPeople(conn, peopleIdent)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        } else {
            log.Print(err)
        }
    }
    fmt.Fprintf(w, "%s", peopleJson)
}

func UpdatePeople(conn *pgx.Conn, people People) error {
    sql := `update kkm.people 
            set name=$1, dob=$2, tel=$3, address=$4, race=$5,
              nationality=$6, edu_lvl=$7, occupation=$8, comorbids=$9, support_vac=$10 
            where ident=$11`   

    _, err := conn.Exec(context.Background(), sql,
        people.Name, people.Dob, people.Tel, people.Address, 
        people.Race, people.Nationality, people.Edu_lvl, 
        people.Occupation, people.Comorbids, people.Support_vac,
        people.Ident)
    if err != nil {
        return err
    }    
    return nil
}

func UpdatePeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method == "OPTIONS") { return }
    r.ParseForm()
    fmt.Println("[UpdatePeopleHandler] Request form data received")

    /* MORE-EFFICIENT-JSON_DECODING-METHOD */
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%v\n", people)

    /* LESS-EFFICIENT-JSON_DECODING-METHOD */
    // var people People
    // err := json.Unmarshal([]byte(input), &people)
    // if err != nil {
    //     log.Print(err)
    //     w.WriteHeader(500)
    //     fmt.Fprintf(w, "Internal server error! Unable to read http json input")
    //     return
    // }

    err = UpdatePeople(conn, people)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }   
}

func AddPeople(conn *pgx.Conn, people People) error {
    sql := `insert into kkm.people
            (ident, name, dob, tel, address, race, nationality,
            edu_lvl, occupation, comorbids, support_vac)
            values 
            ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
    
    _, err := conn.Exec(context.Background(), sql, 
        people.Ident, people.Name, people.Dob, people.Tel, people.Address, 
        people.Race, people.Nationality, people.Edu_lvl, 
        people.Occupation, people.Comorbids, people.Support_vac)
    if err != nil {
        return err
    }
    return nil
}

func AddPeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method == "OPTIONS") { return }
    r.ParseForm()
    fmt.Println("[AddPeopleHandler] Request form data received")
    
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%v\n", people)

    err = AddPeople(conn, people)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }
}

func DeletePeople(conn *pgx.Conn, identity Identity) error {
    sql := `delete from kkm.people 
            where ident=$1`

    _, err := conn.Exec(context.Background(), sql, identity.Ident)
    if err != nil {
        return err
    }
    return nil
}

type Identity struct {
    Ident string    `json:"ident"`
}

func DeletePeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[DeletePeopleHandler] Request form data received")

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%v\n", identity)

    err = DeletePeople(conn, identity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[TestHandler] Request form data received")

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%+v\n", identity)    
}


