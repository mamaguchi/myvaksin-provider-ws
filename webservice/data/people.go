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

type Identity struct {
    Ident string    `json:"ident"`
}

type Peoples struct {
    Peoples []People    `json:"peoples"`
}

type People struct {
    Ident string          `json:"ident"`
    Name string           `json:"name"`
    Gender string         `json:"gender"`
    Dob time.Time         `json:"dob"`
    Nationality string    `json:"nationality"`
    Race string           `json:"race"`
    Tel string            `json:"tel"`
    Address string        `json:"address"`  
    PostalCode string     `json:"postalCode"` 
    Locality string       `json:"locality"`
    District string       `json:"district"`
    State string          `json:"state"` 
    EduLvl string         `json:"eduLvl"`
    Occupation string     `json:"occupation"`
    Comorbids []int       `json:"comorbids"`
    SupportVac bool       `json:"supportVac"`
}

type Vaccine struct {
    Brand string          `json:"brand"`
    Type string           `json:"type"`
    Against string        `json:"against"`
    Raoa string           `json:"raoa"`
}

type Vaccination struct {
    Vaccination string    `json:"vaccination"`
    Aoa string            `json:"aoa"`
    Fa bool               `json:"fa"`
    Fdd time.Time         `json:"fdd"`
    Sdd time.Time         `json:"sdd"`
    AefiClass string      `json:"aefiClass"`
    AefiReaction []string `json:"aefiReaction"`
    Remarks string        `json:"remarks"`
}

type PeoplePage struct {
    People People            `json:"people"`
    Vaccine Vaccine          `json:"vaccine"`
    Vaccination Vaccination  `json:"vaccination"`
}

type VaccinationRecord struct {
    Vaccination string       `json:"vaccination"`
    VaccineBrand string      `json:"vaccineBrand"`
    VaccineType string       `json:"vaccineType"`
    VaccineAgainst string    `json:"vaccineAgainst"`
    VaccineRaoa string       `json:"vaccineRaoa"`
    Aoa string               `json:"aoa"`
    Fa bool                  `json:"fa"`
    Fdd time.Time            `json:"fdd"`
    Sdd time.Time            `json:"sdd"`
    AefiClass string         `json:"aefiClass"`
    AefiReaction []string    `json:"aefiReaction"`
    Remarks string           `json:"remarks"`
}

type PeopleProfile struct {
    People People                           `json:"people"`   
    VaccinationRecords []VaccinationRecord  `json:"vaccinationRecords"` 
}

type PeopleSearchResult struct {
    Ident string             `json:"ident"`
    Name string              `json:"name"`
    Dob time.Time            `json:"dob"`
    Race string              `json:"race"`
    Nationality string       `json:"nationality"`
    Locality string          `json:"locality"`
    District string          `json:"district"`
    State string             `json:"state"`
}

type PeopleSearch struct {
    SearchResults []PeopleSearchResult    `json:"peopleSearchResults"`
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
        var eduLvl string
        var occupation string
        var comorbids []int
        var supportVac bool
        err := rows.Scan(&ident, &name, &dob, &tel, &address, &race,
                        &nationality, &eduLvl, &occupation, &comorbids, 
                        &supportVac)
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
            EduLvl: eduLvl,
            Occupation: occupation,
            Comorbids: comorbids,
            SupportVac: supportVac,
        }
        peoples.Peoples = append(peoples.Peoples, people)
    }
    output, err := json.MarshalIndent(peoples, "", "\t")
        
    return output, err
}

func GetPeople(conn *pgx.Conn, ident string) ([]byte, error) {
    row := conn.QueryRow(context.Background(), 
        `select kkm.people.name, kkm.people.gender, kkm.people.dob, 
        kkm.people.nationality, kkm.people.race, kkm.people.tel, 
        kkm.people.address, kkm.people.postalcode, kkm.people.locality, 
        kkm.people.district, kkm.people.state, kkm.people.edu_lvl, 
        kkm.people.occupation, kkm.people.comorbids, kkm.people.support_vac,
        kkm.vaccine.brand, kkm.vaccine.type, kkm.vaccine.against, 
        kkm.vaccine.raoa,
        kkm.vaccination.vaccination, kkm.vaccination.aoa, 
        kkm.vaccination.first_adm, kkm.vaccination.first_dose_dt,
        kkm.vaccination.second_dose_dt, kkm.vaccination.aefi_class,
        kkm.vaccination.aefi_reaction
          from kkm.people 
            join kkm.vaccination 
              on kkm.people.ident = kkm.vaccination.people
            join kkm.vaccine
              on kkm.vaccination.vaccine = kkm.vaccine.id
          where ident=$1`,
        ident)
    // People
    var name string
    var gender string
    var dob time.Time
    var nationality string
    var race string
    var tel string
    var address string
    var postalCode string 
    var locality string 
    var district string 
    var state string 
    var eduLvl string
    var occupation string
    var comorbids []int
    var supportVac bool
    // Vaccine
    var brand string 
    var vacType string 
    var against string 
    var raoa string 
    // Vaccination
    var vaccination string  
    var aoa string
    var fa bool 
    var fdd time.Time 
    var sdd time.Time 
    var aefiClass string 
    var aefiReaction string 
    err := row.Scan(&name, &gender, &dob, &nationality, &race, &tel, &address,
                &postalCode, &locality, &district, &state, &eduLvl, &occupation, 
                &comorbids, &supportVac, 
                &brand, &vacType, &against, &raoa, 
                &vaccination, &aoa, &fa, &fdd, &sdd, &aefiClass, &aefiReaction)
    if err != nil {
        return nil, err
    }
    people := People{
        Ident: ident,
        Name: name,
        Gender: gender,
        Dob: dob,
        Nationality: nationality,
        Race: race,
        Tel: tel,
        Address: address,
        PostalCode: postalCode,
        Locality: locality,
        District: district,
        State: state,
        EduLvl: eduLvl,
        Occupation: occupation,
        Comorbids: comorbids,
        SupportVac: supportVac,
    } 
    vaccine := Vaccine{
        Brand: brand,
        Type: vacType,
        Against: against,
        Raoa: raoa,
    }   
    vaccinationStrct := Vaccination{
        Vaccination: vaccination,
        Aoa: aoa,
        Fa: fa,
        Fdd: fdd,
        Sdd: sdd,
        AefiClass: aefiClass,
        AefiReaction: []string{ aefiReaction },
    }
    peoplePage := PeoplePage{
        People: people,
        Vaccine: vaccine,
        Vaccination: vaccinationStrct,
    }
    outputJson, err := json.MarshalIndent(peoplePage, "", "\t")
        
    return outputJson, err
}

func GetPeopleProfile(conn *pgx.Conn, ident string) ([]byte, error) {
    rows, err := conn.Query(context.Background(), 
        `select people.name, people.gender, people.dob, 
         people.nationality, people.race, people.tel, 
         people.address, people.postalcode, people.locality, 
         people.district, people.state, people.edu_lvl, 
         people.occupation, people.comorbids, people.support_vac,
         vaccine.brand, vaccine.type, vaccine.against, 
         vaccine.raoa,
         vaccination.vaccination, vaccination.aoa, 
         vaccination.first_adm, vaccination.first_dose_dt,
         vaccination.second_dose_dt, vaccination.aefi_class,
         vaccination.aefi_reaction, vaccination.remarks
           from kkm.people 
             join kkm.vaccination 
               on kkm.people.ident = kkm.vaccination.people
             join kkm.vaccine
               on kkm.vaccination.vaccine = kkm.vaccine.id
           where ident=$1`,
        ident)
    if err != nil {
        return nil, err
    }
    var peopleProfile PeopleProfile
    firstRecord := true

    for rows.Next() {
        // Vaccine
        var brand string 
        var vacType string 
        var against string 
        var raoa string 
        // Vaccination
        var vaccination string  
        var aoa string
        var fa bool 
        var fdd time.Time 
        var sdd time.Time 
        var aefiClass string 
        var aefiReaction []string 
        var remarks string 

        if firstRecord {
            // People
            var name string
            var gender string
            var dob time.Time
            var nationality string
            var race string
            var tel string
            var address string
            var postalCode string 
            var locality string 
            var district string 
            var state string 
            var eduLvl string
            var occupation string
            var comorbids []int
            var supportVac bool

            err = rows.Scan(&name, &gender, &dob, &nationality, &race, &tel, &address,
                &postalCode, &locality, &district, &state, &eduLvl, &occupation, 
                &comorbids, &supportVac, 
                &brand, &vacType, &against, &raoa, 
                &vaccination, &aoa, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)
            if err != nil {
                return nil, err
            }
            peopleProfile.People = People{
                Ident: ident,
                Name: name,
                Gender: gender,
                Dob: dob,
                Nationality: nationality,
                Race: race,
                Tel: tel,
                Address: address,
                PostalCode: postalCode,
                Locality: locality,
                District: district,
                State: state,
                EduLvl: eduLvl,
                Occupation: occupation,
                Comorbids: comorbids,
                SupportVac: supportVac,
            }
            firstRecord = false                                     
        } else {
            err = rows.Scan(nil, nil, nil, nil, nil, nil, nil,
                nil, nil, nil, nil, nil, nil, nil, nil, 
                &brand, &vacType, &against, &raoa, 
                &vaccination, &aoa, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)                      
            if err != nil {
                return nil, err
            }
        }
        vaccinationRecord := VaccinationRecord{
            Vaccination: vaccination,
            VaccineBrand: brand,
            VaccineType: vacType,
            VaccineAgainst: against,
            VaccineRaoa: raoa,
            Aoa: aoa,
            Fa: fa,
            Fdd: fdd,
            Sdd: sdd,
            AefiClass: aefiClass,
            AefiReaction: aefiReaction,
            Remarks: remarks,
        }
        peopleProfile.VaccinationRecords = append(
            peopleProfile.VaccinationRecords,
            vaccinationRecord,
        )
    }

    outputJson, err := json.MarshalIndent(peopleProfile, "", "\t")        
    return outputJson, err
}

func GetPeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type") 
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[GetPeopleHandler] Request form data received")
    // r.ParseForm()
    // fmt.Println(r.Form)
    // peopleIdent := r.Form["ident"][0]
    // fmt.Printf("%s\n", peopleIdent)    

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // peopleProfJson, err := GetPeople(conn, identity.Ident)
    peopleProfJson, err := GetPeopleProfile(conn, identity.Ident)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        }
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }
    fmt.Fprintf(w, "%s", peopleProfJson)
}

func SearchPeople(conn *pgx.Conn, ident string) ([]byte, error) {
    sql1 := `select people.ident, people.name, people.dob, people.race, 
               people.nationality, people.locality, people.district, people.state 
             from kkm.people
             where ident=$1`
    
    rows, err := conn.Query(context.Background(), sql1, ident)
    if err != nil {
        return nil, err 
    }

    var peopleSearch PeopleSearch
    for rows.Next() {
        var ident string 
        var name string 
        var dob time.Time
        var race string 
        var nationality string 
        var locality string 
        var district string 
        var state string 

        err = rows.Scan(&ident, &name, &dob, &race, &nationality, 
                        &locality, &district, &state) 
        if err != nil {
            return nil, err 
        }                   
        peopleSearchResult := PeopleSearchResult{
            Ident: ident,
            Name: name,
            Dob: dob,
            Race: race,
            Nationality: nationality,
            Locality: locality,
            District: district,
            State: state,
        }     
        peopleSearch.SearchResults = append(
            peopleSearch.SearchResults,
            peopleSearchResult)
    }

    outputJson, err := json.MarshalIndent(peopleSearch, "", "\t")
    return outputJson, err
}

func SearchPeopleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method =="OPTIONS") {return}
    fmt.Println("[SearchPeopleHandler] request received")

    var identity Identity 
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    SearchPeopleResultJson, err := SearchPeople(conn, identity.Ident)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        }
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }
    fmt.Fprintf(w, "%s", SearchPeopleResultJson)
}

func UpdatePeople(conn *pgx.Conn, people People) error {
    sql := `update kkm.people 
            set name=$1, dob=$2, tel=$3, address=$4, race=$5,
              nationality=$6, edu_lvl=$7, occupation=$8, comorbids=$9, support_vac=$10 
            where ident=$11`   

    _, err := conn.Exec(context.Background(), sql,
        people.Name, people.Dob, people.Tel, people.Address, 
        people.Race, people.Nationality, people.EduLvl, 
        people.Occupation, people.Comorbids, people.SupportVac,
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
        people.Race, people.Nationality, people.EduLvl, 
        people.Occupation, people.Comorbids, people.SupportVac)
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

 

