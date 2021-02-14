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

type SqlInputVars struct {
    Ident string            `json:"ident"`
    Name string             `json:"name"`
    DobInterval DobInterval `json:"dobInterval"`
    Race string             `json:"race"`
    Nationality string      `json:"nationality"`
    State string            `json:"state"`
    District string         `json:"district"`
    Locality string         `json:"locality"`
    SqlOpt string           `json:"sqlOpt"`
}

type DobInterval struct {
    MinDate string      `json:"minDate"`
    MaxDate string      `json:"maxDate"`
}

type Peoples struct {
    Peoples []People    `json:"peoples"`
}

type People struct {
    Ident string          `json:"ident"`
    Name string           `json:"name"`
    Gender string         `json:"gender"`
    // Dob time.Time         `json:"dob"`
    Dob string            `json:"dob"`
    Nationality string    `json:"nationality"`
    Race string           `json:"race"`
    Tel string            `json:"tel"`
    Email string          `json:"email"`
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
    // Fa bool                  `json:"fa"`
    // Fdd time.Time            `json:"fdd"`
    // Sdd time.Time            `json:"sdd"`
    Fa string                `json:"fa"`
    Fdd string               `json:"fdd"`
    Sdd string               `json:"sdd"`
    AefiClass string         `json:"aefiClass"`
    AefiReaction []string    `json:"aefiReaction"`
    Remarks string           `json:"remarks"`
}

type PeopleProfile struct {
    People People                           `json:"people"`   
    VaccinationRecords []VaccinationRecord  `json:"vaccinationRecords"` 
}

type VacRecUpsert struct {
    Ident string             `json:"ident"`   
    VacRec VaccinationRecord `json:"vacRec"` 
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

func GetPeoples(conn *pgx.Conn) ([]byte, error) {
    var peoples Peoples
    rows, _ := conn.Query(context.Background(), 
        "select ident, name, dob::text, tel, address, race, nationality, eduLvl, occupation, comorbids, supportVac from kkm.people")

    for rows.Next() {
        var ident string
        var name string
        // var dob time.Time
        var dob string
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
        `select kkm.people.name, kkm.people.gender, kkm.people.dob::text, 
        kkm.people.nationality, kkm.people.race, kkm.people.tel, 
        kkm.people.address, kkm.people.postalcode, kkm.people.locality, 
        kkm.people.district, kkm.people.state, kkm.people.eduLvl, 
        kkm.people.occupation, kkm.people.comorbids, kkm.people.supportVac,
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
    // var dob time.Time
    var dob string
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
        `select people.name, people.gender, people.dob::text, 
         people.nationality, people.race, people.tel, people.email,
         people.address, people.postalcode, people.locality, 
         people.district, people.state, people.eduLvl, 
         people.occupation, people.comorbids, people.supportVac,
         vaccine.brand, vaccine.type, vaccine.against, 
         vaccine.raoa,
         vaccination.vaccination, vaccination.firstAdm::text, 
         coalesce(vaccination.firstDoseDt::text, '') as firstDoseDt, 
         coalesce(vaccination.secondDoseDt::text, '') as secondDoseDt, 
         vaccination.aefiClass, vaccination.aefiReaction, vaccination.remarks
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
        // var aoa string
        // var fa bool 
        // var fdd time.Time 
        // var sdd time.Time 
        var fa string 
        var fdd string
        var sdd string 
        var aefiClass string 
        var aefiReaction []string 
        var remarks string 

        if firstRecord {
            // People
            var name string
            var gender string
            // var dob time.Time
            var dob string
            var nationality string
            var race string
            var tel string
            var email string 
            var address string
            var postalCode string 
            var locality string 
            var district string 
            var state string 
            var eduLvl string
            var occupation string
            var comorbids []int
            var supportVac bool

            err = rows.Scan(&name, &gender, &dob, &nationality, &race, &tel, 
                &email, &address, &postalCode, &locality, &district, &state, 
                &eduLvl, &occupation, &comorbids, &supportVac, 
                &brand, &vacType, &against, &raoa, 
                &vaccination, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)
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
                Email: email,
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
                nil, nil, nil, nil, nil, nil, nil, nil, nil,
                &brand, &vacType, &against, &raoa, 
                &vaccination, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)                      
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
            // Aoa: aoa,
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
    fmt.Println("[GetPeopleHandler] request received")
    // r.ParseForm()
    // fmt.Println(r.Form)
    // peopleIdent := r.Form["ident"][0]
    // fmt.Printf("%s\n", peopleIdent)    

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // peopleProfJson, err := GetPeople(conn, identity.Ident)
    peopleProfJson, err := GetPeopleProfile(conn, identity.Ident)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        }
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }
    fmt.Printf("%s\n", peopleProfJson)
    fmt.Fprintf(w, "%s", peopleProfJson)
}

func SearchPeople(conn *pgx.Conn, sqlInputVars SqlInputVars) ([]byte, error) {
    sqlOpt1 := 
        `select people.ident, people.name, people.dob, people.race, 
           people.nationality, people.locality, people.district, people.state 
         from kkm.people
         where ident=$1`

    // NOTE: pgx (Golang PostgreSQL driver) does not support the term
    //       'timestamp' before the date string in the sql, or else it
    //       will cause syntax error.
    //       Using 'timestamp' term before date string is supported 
    //       but optional in native psql command.
    sqlOpt2 := 
         `select people.ident, people.name, people.dob, people.race,
            people.nationality, people.locality, people.district, people.state
          from kkm.people
          where dob between $1 and $2`

    sqlOpt3 := 
        `select people.ident, people.name, people.dob, people.race,
           people.nationality, people.locality, people.district, people.state
         from kkm.people
         where name ilike $1
           and race::text ilike $2
           and nationality::text ilike $3
           and state::text ilike $4
           and district ilike $5
           and locality ilike $6`
    
    var rows pgx.Rows 
    var err error
    if sqlInputVars.SqlOpt == "1" {
        rows, err = conn.Query(context.Background(), sqlOpt1, 
          sqlInputVars.Ident)        
    } else if sqlInputVars.SqlOpt == "2" {
        rows, err = conn.Query(context.Background(), sqlOpt2, 
          sqlInputVars.DobInterval.MinDate,
          sqlInputVars.DobInterval.MaxDate)   
    } else if sqlInputVars.SqlOpt == "3" {
        rows, err = conn.Query(context.Background(), sqlOpt3, 
          sqlInputVars.Name,
          sqlInputVars.Race,
          sqlInputVars.Nationality,
          sqlInputVars.State,
          sqlInputVars.District,
          sqlInputVars.Locality)
    }
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

    var sqlInputVars SqlInputVars
    err := json.NewDecoder(r.Body).Decode(&sqlInputVars)
    if err != nil {
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%+v\n", sqlInputVars)

    SearchPeopleResultJson, err := SearchPeople(conn, sqlInputVars)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        }
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }
    fmt.Printf("JSON Output\n%s\n", SearchPeopleResultJson)
    fmt.Fprintf(w, "%s", SearchPeopleResultJson)
}

func UpdatePeople(conn *pgx.Conn, people People) error {
    // sql := `update kkm.people 
    //         set name=$1, dob=$2, tel=$3, address=$4, race=$5,
    //           nationality=$6, eduLvl=$7, occupation=$8, comorbids=$9, supportVac=$10 
    //         where ident=$11`   

    // _, err := conn.Exec(context.Background(), sql,
    //     people.Name, people.Dob, people.Tel, people.Address, 
    //     people.Race, people.Nationality, people.EduLvl, 
    //     people.Occupation, people.Comorbids, people.SupportVac,
    //     people.Ident)

    sql := 
        `update kkm.people 
           set name=$1, gender=$2, dob=$3, nationality=$4, race=$5, 
             tel=$6, email=$7, address=$8, postalCode=$9, locality=$10,
             district=$11, state=$12, eduLvl=$13, occupation=$14, 
             comorbids=$15, supportVac=$16 
           where ident=$17`   

    _, err := conn.Exec(context.Background(), sql,
        people.Name, people.Gender, people.Dob, people.Nationality, 
        people.Race, people.Tel, people.Email, people.Address, people.PostalCode, 
        people.Locality, people.District, people.State, people.EduLvl, 
        people.Occupation, people.Comorbids, people.SupportVac, people.Ident)
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
    fmt.Println("[UpdatePeopleHandler] request received")
    
    /* LESS-EFFICIENT-JSON_DECODING-METHOD (Produces intermediate byte slice)
       var people People
       err := json.Unmarshal([]byte(input), &people) */

    /* MORE-EFFICIENT-JSON_DECODING-METHOD (No intermediate byte slice) */
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%+v\n", people)

    err = UpdatePeople(conn, people)
    if err != nil {
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }   
}

func AddPeople(conn *pgx.Conn, people People) error {
    sql := `insert into kkm.people
            (ident, name, dob, tel, address, race, nationality,
            eduLvl, occupation, comorbids, supportVac)
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

// type VaccinationRecord struct {
//     Vaccination string       `json:"vaccination"`
//     VaccineBrand string      `json:"vaccineBrand"`
//     VaccineType string       `json:"vaccineType"`
//     VaccineAgainst string    `json:"vaccineAgainst"`
//     VaccineRaoa string       `json:"vaccineRaoa"`
//     Aoa string               `json:"aoa"`
//     Fa bool                  `json:"fa"`
//     Fdd time.Time            `json:"fdd"`
//     Sdd time.Time            `json:"sdd"`
//     AefiClass string         `json:"aefiClass"`
//     AefiReaction []string    `json:"aefiReaction"`
//     Remarks string           `json:"remarks"`
// }
// `insert into kkm.vaccination
// (vaccine, people, vaccination, aoa, first_adm,
//   first_dose_dt, second_dose_dt, aefi_class, 
//   aefi_reaction, remarks)`
func InsertNewVacRec(conn *pgx.Conn, vru VacRecUpsert) error {                  
    var err error
    if vru.VacRec.Fdd == "" {
        sql := 
            `insert into kkm.vaccination
            (
                vaccine, people, vaccination, firstAdm,  
                secondDoseDt, aefiClass, aefiReaction, remarks
            )
            select vac.id, $1, $2, $3, $4, $5, $6, $7
            from kkm.vaccine vac
            where vac.brand=$8` 

        _, err = conn.Exec(context.Background(), sql, 
        vru.Ident, vru.VacRec.Vaccination, vru.VacRec.Fa,
        vru.VacRec.Sdd, vru.VacRec.AefiClass,
        vru.VacRec.AefiReaction, vru.VacRec.Remarks,
        vru.VacRec.VaccineBrand)
    } else if vru.VacRec.Sdd == "" {
        sql := 
            `insert into kkm.vaccination
            (
                vaccine, people, vaccination, firstAdm, firstDoseDt, 
                aefiClass, aefiReaction, remarks
            )
            select vac.id, $1, $2, $3, $4, $5, $6, $7
            from kkm.vaccine vac
            where vac.brand=$8`

        _, err = conn.Exec(context.Background(), sql, 
        vru.Ident, vru.VacRec.Vaccination, vru.VacRec.Fa,
        vru.VacRec.Fdd, vru.VacRec.AefiClass,
        vru.VacRec.AefiReaction, vru.VacRec.Remarks,
        vru.VacRec.VaccineBrand)
    } else if vru.VacRec.Fdd == "" && vru.VacRec.Sdd == "" {
        sql := 
            `insert into kkm.vaccination
            (
                vaccine, people, vaccination, firstAdm,  
                aefiClass, aefiReaction, remarks
            )
            select vac.id, $1, $2, $3, $4, $5, $6
            from kkm.vaccine vac
            where vac.brand=$7` 
        _, err = conn.Exec(context.Background(), sql, 
        vru.Ident, vru.VacRec.Vaccination, vru.VacRec.Fa,
        vru.VacRec.AefiClass,
        vru.VacRec.AefiReaction, vru.VacRec.Remarks,
        vru.VacRec.VaccineBrand)
    } else {
        sql := 
            `insert into kkm.vaccination
            (
                vaccine, people, vaccination, firstAdm, firstDoseDt, 
                secondDoseDt, aefiClass, aefiReaction, remarks
            )
            select vac.id, $1, $2, $3, $4, $5, $6, $7, $8
            from kkm.vaccine vac
            where vac.brand=$9` 
        _, err = conn.Exec(context.Background(), sql, 
        vru.Ident, vru.VacRec.Vaccination, vru.VacRec.Fa,
        vru.VacRec.Fdd, vru.VacRec.Sdd, vru.VacRec.AefiClass,
        vru.VacRec.AefiReaction, vru.VacRec.Remarks,
        vru.VacRec.VaccineBrand)
    }    
    if err != nil {
        return err
    }
    return nil
}

func InsertNewVacRecHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "authorization")
    w.Header().Set("Access-Control-Allow-Headers", "content-type")
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[InsertNewVacRecHandler] request received")
        
    var vru VacRecUpsert
    err := json.NewDecoder(r.Body).Decode(&vru)
    if err != nil {
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Printf("%+v\n", vru)

    err = InsertNewVacRec(conn, vru)
    if err != nil {
        log.Print(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }   
}

 

