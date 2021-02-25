package data

import (
    "net/http"
    "encoding/json"
    "time"
    "strconv"
    "strings"
    "fmt"
    "log"
    "context"
    "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
    "myvaksin/webservice/db"
    "myvaksin/webservice/auth"
    "myvaksin/webservice/util"
)

const (
    DATE_ISO =  "2006-01-02"
)

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

// type Peoples struct {
//     Peoples []People    `json:"peoples"`
// }

type People struct {
    Ident string          `json:"ident"`
    Name string           `json:"name"`
    Gender string         `json:"gender"`
    Dob string            `json:"dob"` //kiv change to time.Time type
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
    ProfilePicData string `json:"profilePicData"`
    Role string           `json:"role"` 
}

// type Vaccine struct {
//     Brand string          `json:"brand"`
//     Type string           `json:"type"`
//     Against string        `json:"against"`
//     Raoa string           `json:"raoa"`
// }

// type Vaccination struct {
//     Vaccination string    `json:"vaccination"`
//     Aoa string            `json:"aoa"`
//     Fa bool               `json:"fa"`
//     Fdd time.Time         `json:"fdd"`
//     Sdd time.Time         `json:"sdd"`
//     AefiClass string      `json:"aefiClass"`
//     AefiReaction []string `json:"aefiReaction"`
//     Remarks string        `json:"remarks"`
// }

type VaccinationRecord struct {
    VaccinationId int64      `json:"vaccinationId"`
    Vaccination string       `json:"vaccination"`
    VaccineBrand string      `json:"vaccineBrand"`
    VaccineType string       `json:"vaccineType"`
    VaccineAgainst string    `json:"vaccineAgainst"`
    VaccineRaoa string       `json:"vaccineRaoa"`    
    Fa string                `json:"fa"`
    Fdd string               `json:"fdd"` //kiv change to time.Time type
    Sdd string               `json:"sdd"` //kiv change to time.Time type
    AefiClass string         `json:"aefiClass"`
    AefiReaction []string    `json:"aefiReaction"`
    Remarks string           `json:"remarks"`
}

type VaccinationRecord2 struct {
    VaccinationId int64      `json:"vaccinationId"`
    Vaccination string       `json:"vaccination"`
    FdVaccineBrand string    `json:"fdVaccineBrand"`
    SdVaccineBrand string    `json:"sdVaccineBrand"`
    FdTCA string             `json:"fdTCA"` //kiv change to time.Time type
    FdGiven string           `json:"fdGiven"` //kiv change to time.Time type
    SdTCA string             `json:"sdTCA"` //kiv change to time.Time type
    SdGiven string           `json:"sdGiven"` //kiv change to time.Time type
    FdAefiClass string       `json:"fdAefiClass"`
    SdAefiClass string       `json:"sdAefiClass"`
    FdAefiReaction []string  `json:"fdAefiReaction"`
    SdAefiReaction []string  `json:"sdAefiReaction"`
    FdRemarks string         `json:"fdRemarks"`
    SdRemarks string         `json:"sdRemarks"`
}

type PeopleProfile struct {
    People People                           `json:"people"`   
    VaccinationRecords []VaccinationRecord  `json:"vaccinationRecords"` 
}

type PeopleProfile2 struct {
    People People                           `json:"people"`   
    VaccinationRecords []VaccinationRecord2  `json:"vaccinationRecords"` 
}

type VacRecUpsert struct {
    Ident string             `json:"ident"`//People's Ident  
    VacRec VaccinationRecord `json:"vacRec"` 
}

type VacRecUpsert2 struct {
    Ident string              `json:"ident"`//People's Ident  
    VacRec VaccinationRecord2 `json:"vacRec"` 
}

type VacRecDelete struct {
    VaccinationId int64      `json:"vaccinationId`
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
    Vaccination string       `json:"vaccination"`
    VaccineBrand string      `json:"vaccineBrand"`
    NumDose string           `json:"numDose"`
    DoseTaken string         `json:"doseTaken"`    
}

type PeopleSearchResult2 struct {
    Ident string             `json:"ident"`
    Name string              `json:"name"`
    Dob time.Time            `json:"dob"`
    Race string              `json:"race"`
    Nationality string       `json:"nationality"`
    Locality string          `json:"locality"`
    District string          `json:"district"`
    State string             `json:"state"`
    Vaccination string       `json:"vaccination"`
    FdVaccineBrand string    `json:"fdVaccineBrand"`
    SdVaccineBrand string    `json:"sdVaccineBrand"`
    NumDose string           `json:"numDose"`
    DoseTaken string         `json:"doseTaken"`    
}

type PeopleSearch struct {
    SearchResults []PeopleSearchResult    `json:"peopleSearchResults"`
}

type PeopleSearch2 struct {
    SearchResults []PeopleSearchResult2    `json:"peopleSearchResults"`
}


func TestHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
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

func SearchPeople(conn *pgxpool.Pool, sqlInputVars SqlInputVars) ([]byte, error) {    
    sqlOpt1 := 
        `select people.ident, people.name, people.dob, people.race,
            people.nationality, people.locality, people.district, people.state,
            coalesce(vaccine.brand, '') as brand,
            coalesce(vaccine.numdose, 0) as numdose,
            coalesce(vaccination.vaccination, '') as vaccination,
            coalesce(vaccination.firstdosedt::text, '') as firstdosedt, 
            coalesce(vaccination.seconddosedt::text, '') as seconddosedt, 
            coalesce(vaccination.thirddosedt::text, '') as thirddosedt            
            from kkm.people 
                left join kkm.vaccination 
                on kkm.people.ident = kkm.vaccination.people
                left join kkm.vaccine
                on kkm.vaccination.vaccine = kkm.vaccine.id
            where ident=$1`

    // NOTE: pgx (Golang PostgreSQL driver) does not support the term
    //       'timestamp' before the date string in the sql, or else it
    //       will cause syntax error.
    //       Using 'timestamp' term before date string is supported 
    //       but optional in native psql command.
    sqlOpt2 := 
         `select people.ident, people.name, people.dob, people.race,
            people.nationality, people.locality, people.district, people.state,
            coalesce(vaccine.brand, '') as brand,
            coalesce(vaccine.numdose, 0) as numdose,
            coalesce(vaccination.vaccination, '') as vaccination,
            coalesce(vaccination.firstdosedt::text, '') as firstdosedt, 
            coalesce(vaccination.seconddosedt::text, '') as seconddosedt, 
            coalesce(vaccination.thirddosedt::text, '') as thirddosedt            
            from kkm.people 
                left join kkm.vaccination 
                on kkm.people.ident = kkm.vaccination.people
                left join kkm.vaccine
                on kkm.vaccination.vaccine = kkm.vaccine.id
            where dob between $1 and $2`

    sqlOpt3 := 
        `select people.ident, people.name, people.dob, people.race,
            people.nationality, people.locality, people.district, people.state,
            coalesce(vaccine.brand, '') as brand,
            coalesce(vaccine.numdose, 0) as numdose,
            coalesce(vaccination.vaccination, '') as vaccination,
            coalesce(vaccination.firstdosedt::text, '') as firstdosedt, 
            coalesce(vaccination.seconddosedt::text, '') as seconddosedt, 
            coalesce(vaccination.thirddosedt::text, '') as thirddosedt            
            from kkm.people 
                left join kkm.vaccination 
                on kkm.people.ident = kkm.vaccination.people
                left join kkm.vaccine
                on kkm.vaccination.vaccine = kkm.vaccine.id
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
        var vaccineBrand string 
        var numDose int 
        var vaccination string 
        var fdd string 
        var sdd string 
        var tdd string 

        err = rows.Scan(&ident, &name, &dob, &race, &nationality, 
                        &locality, &district, &state, 
                        &vaccineBrand, &numDose, 
                        &vaccination, &fdd, &sdd, &tdd) 
        if err != nil {
            // TODO: Kiv to add a return which return 
            // a code that signals 0 search results.
            return nil, err 
        }                   

        numDoseStr := strconv.Itoa(numDose)
        doseTaken := 0
        if len(vaccination) != 0 {
            if len(fdd) != 0 {
                doseTaken++
            }
            if len(sdd) != 0 {
                doseTaken++
            }
            if len(tdd) != 0 {
                doseTaken++
            }
        }
        doseTakenStr := strconv.Itoa(doseTaken)

        peopleSearchResult := PeopleSearchResult{
            Ident: ident,
            Name: name,
            Dob: dob,
            Race: race,
            Nationality: nationality,
            Locality: locality,
            District: district,
            State: state,            
            Vaccination: vaccination,
            VaccineBrand: vaccineBrand,
            NumDose: numDoseStr,
            DoseTaken: doseTakenStr,            
        }     
        peopleSearch.SearchResults = append(
            peopleSearch.SearchResults,
            peopleSearchResult)
    }

    outputJson, err := json.MarshalIndent(peopleSearch, "", "")
    return outputJson, err
}

func SearchPeople2(conn *pgxpool.Pool, sqlInputVars SqlInputVars) ([]byte, error) {    
    sqlOpt1 := 
        `select p.ident, p.name, p.dob, p.race,
            p.nationality, p.locality, p.district, p.state,
            coalesce(fdv.brand, '') as fdbrand,
            coalesce(sdv.brand, '') as sdbrand,
            coalesce(fdv.numdose, 0) as fdnumdose,
            coalesce(v.vaccination, '') as vaccination,
            coalesce(v.fdgiven::text, '') as fdgiven, 
            coalesce(v.sdgiven::text, '') as sdgiven
            from kkm.people p
                left join kkm.vaccination v
                on p.ident = v.people
                left join kkm.vaccine fdv
                on v.fdvaccine = fdv.id
                left join kkm.vaccine sdv
                on v.sdvaccine = sdv.id
            where p.ident=$1`

    // NOTE: pgx (Golang PostgreSQL driver) does not support the term
    //       'timestamp' before the date string in the sql, or else it
    //       will cause syntax error.
    //       Using 'timestamp' term before date string is supported 
    //       but optional in native psql command.
    sqlOpt2 := 
         `select p.ident, p.name, p.dob, p.race,
            p.nationality, p.locality, p.district, p.state,
            coalesce(fdv.brand, '') as fdbrand,
            coalesce(sdv.brand, '') as sdbrand,
            coalesce(fdv.numdose, 0) as fdnumdose,
            coalesce(v.vaccination, '') as vaccination,
            coalesce(v.fdgiven::text, '') as fdgiven, 
            coalesce(v.sdgiven::text, '') as sdgiven
            from kkm.people p
                left join kkm.vaccination v
                on p.ident = v.people
                left join kkm.vaccine fdv
                on v.fdvaccine = fdv.id
                left join kkm.vaccine sdv
                on v.sdvaccine = sdv.id
            where dob between $1 and $2`

    sqlOpt3 := 
        `select p.ident, p.name, p.dob, p.race,
            p.nationality, p.locality, p.district, p.state,
            coalesce(fdv.brand, '') as fdbrand,
            coalesce(sdv.brand, '') as sdbrand,
            coalesce(fdv.numdose, 0) as fdnumdose,
            coalesce(v.vaccination, '') as vaccination,
            coalesce(v.fdgiven::text, '') as fdgiven, 
            coalesce(v.sdgiven::text, '') as sdgiven
            from kkm.people p
                left join kkm.vaccination v
                on p.ident = v.people
                left join kkm.vaccine fdv
                on v.fdvaccine = fdv.id
                left join kkm.vaccine sdv
                on v.sdvaccine = sdv.id
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

    var peopleSearch PeopleSearch2
    for rows.Next() {
        var ident string 
        var name string 
        var dob time.Time
        var race string 
        var nationality string 
        var locality string 
        var district string 
        var state string
        var fdVaccineBrand string 
        var sdVaccineBrand string 
        var fdNumDose int 
        var vaccination string 
        var fdGiven string 
        var sdGiven string 

        err = rows.Scan(&ident, &name, &dob, &race, &nationality, 
                        &locality, &district, &state, 
                        &fdVaccineBrand, &sdVaccineBrand, 
                        &fdNumDose, &vaccination, 
                        &fdGiven, &sdGiven) 
        if err != nil {
            // TODO: Kiv to add a return which return 
            // a code that signals 0 search results.
            return nil, err 
        }                   

        numDoseStr := strconv.Itoa(fdNumDose)
        doseTaken := 0
        if len(vaccination) != 0 {
            if len(fdGiven) != 0 {
                doseTaken++
            }
            if len(sdGiven) != 0 {
                doseTaken++
            }           
        }
        doseTakenStr := strconv.Itoa(doseTaken)

        peopleSearchResult := PeopleSearchResult2{
            Ident: ident,
            Name: name,
            Dob: dob,
            Race: race,
            Nationality: nationality,
            Locality: locality,
            District: district,
            State: state,            
            Vaccination: vaccination,
            FdVaccineBrand: fdVaccineBrand,
            SdVaccineBrand: sdVaccineBrand,
            NumDose: numDoseStr,
            DoseTaken: doseTakenStr,            
        }     
        peopleSearch.SearchResults = append(
            peopleSearch.SearchResults,
            peopleSearchResult)
    }

    outputJson, err := json.MarshalIndent(peopleSearch, "", "")
    return outputJson, err
}

func SearchPeopleHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method =="OPTIONS") {return}
    fmt.Println("[SearchPeopleHandler] request received")   

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }

    var sqlInputVars SqlInputVars
    err := json.NewDecoder(r.Body).Decode(&sqlInputVars)
    if err != nil {        
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", sqlInputVars)

    db.CheckDbConn()
    SearchPeopleResultJson, err := SearchPeople2(db.Conn, sqlInputVars)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        }
        util.SendInternalServerErrorStatus(w, err)
        return 
    }
    fmt.Printf("JSON Output\n%s\n", SearchPeopleResultJson)
    fmt.Fprintf(w, "%s", SearchPeopleResultJson)
}

func GetCovidVacRec(conn *pgxpool.Pool, ident string) ([]byte, error) {   
    row := conn.QueryRow(context.Background(), 
        `select 
            vaccine.brand, 
            vaccine.type, 
            vaccine.against, 
            vaccine.raoa, 
            vaccination.id, 
            vaccination.firstAdm::text,                     
            coalesce(vaccination.firstDoseDt::text, '') as firstDoseDt, 
            coalesce(vaccination.secondDoseDt::text, '') as secondDoseDt, 
            vaccination.aefiClass::text, 
            coalesce(vaccination.aefiReaction, '{}') as aefiReaction, 
            vaccination.remarks
         from kkm.vaccination 
             join kkm.vaccine
               on kkm.vaccination.vaccine = kkm.vaccine.id
           where vaccination.people=$1 
             and vaccination.vaccination=$2`,
        ident, "COVID-19")
   
    // Vaccine
    var brand string 
    var vacType string 
    var against string 
    var raoa string 
    // Vaccination
    var vaccinationId int64
    var fa string 
    var fdd string
    var sdd string 
    var aefiClass string 
    var aefiReaction []string 
    var remarks string 

    err := row.Scan(&brand, &vacType, &against, &raoa, 
        &vaccinationId, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)
    if err != nil {
        return nil, err
    }

    vaccinationRecord := VaccinationRecord{
        VaccinationId: vaccinationId,
        Vaccination: "COVID-19",
        VaccineBrand: brand,
        VaccineType: vacType,
        VaccineAgainst: against,
        VaccineRaoa: raoa,
        Fa: fa,
        Fdd: fdd,
        Sdd: sdd,
        AefiClass: aefiClass,
        AefiReaction: aefiReaction,
        Remarks: remarks,
    }   
    outputJson, err := json.MarshalIndent(vaccinationRecord, "", "")        
    return outputJson, err
}

func GetCovidVacRec2(conn *pgxpool.Pool, ident string) ([]byte, error) {   
    row := conn.QueryRow(context.Background(), 
        `select 
            coalesce(fdv.brand, '') as fdBrand,
            coalesce(sdv.brand, '') as sdBrand,
            v.id, 
            coalesce(v.fdtca::text, '') as fdTCA, 
            coalesce(v.fdgiven::text, '') as fdGiven, 
            coalesce(v.sdtca::text, '') as sdTCA, 
            coalesce(v.sdgiven::text, '') as sdGiven, 
            coalesce(v.fdaeficlass::text, '') as fdAefiClass, 
            coalesce(v.sdaeficlass::text, '') as sdAefiClass,
            coalesce(v.fdaefireaction, '{}') as fdAefiReaction, 
            coalesce(v.sdaefireaction, '{}') as sdAefiReaction, 
            coalesce(v.fdremarks, '') as fdRemarks,
            coalesce(v.sdremarks, '') as sdRemarks
         from kkm.vaccination v
             join kkm.vaccine fdv
               on v.fdvaccine = fdv.id
             join kkm.vaccine sdv
               on v.sdvaccine = sdv.id
           where v.people=$1 
             and v.vaccination=$2`,
        ident, "COVID-19")
   
    // Vaccine
    var fdBrand string 
    var sdBrand string 
    // Vaccination
    var vaccinationId int64
    var fdTCA string
    var fdGiven string
    var sdTCA string
    var sdGiven string
    var fdAefiClass string 
    var sdAefiClass string 
    var fdAefiReaction []string 
    var sdAefiReaction []string 
    var fdRemarks string 
    var sdRemarks string 

    err := row.Scan(&fdBrand, &sdBrand, &vaccinationId, 
        &fdTCA, &fdGiven, &sdTCA, &sdGiven, 
        &fdAefiClass, &sdAefiClass, &fdAefiReaction, &sdAefiReaction,
        &fdRemarks, &sdRemarks)
    if err != nil {
        return nil, err
    }

    vaccinationRecord := VaccinationRecord2{
        VaccinationId: vaccinationId,
        Vaccination: "COVID-19",
        FdVaccineBrand: fdBrand,
        SdVaccineBrand: sdBrand,
        FdTCA: fdTCA,
        FdGiven: fdGiven,
        SdTCA: sdTCA,
        SdGiven: sdGiven,
        FdAefiClass: fdAefiClass,
        SdAefiClass: sdAefiClass,
        FdAefiReaction: fdAefiReaction,
        SdAefiReaction: sdAefiReaction,
        FdRemarks: fdRemarks,
        SdRemarks: sdRemarks,
    }   
    outputJson, err := json.MarshalIndent(vaccinationRecord, "", "")        
    return outputJson, err
}

func GetCovidVacRecHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[GetCovidVacRecHandler] request received")    

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }   

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    
    db.CheckDbConn()
    covidVacRecJson, err := GetCovidVacRec2(db.Conn, identity.Ident)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("Vaccination record entry not found in database")
        }
        util.SendInternalServerErrorStatus(w, err)
        return 
    }
    fmt.Printf("%s\n", covidVacRecJson)
    fmt.Fprintf(w, "%s", covidVacRecJson)
}

func CreateNewPeople(conn *pgxpool.Pool, people People) (string, error) {
    sqlSelect := 
		`select name from kkm.people
		 where ident=$1`

	row := conn.QueryRow(context.Background(), sqlSelect,
				people.Ident)
	var name string				
	err := row.Scan(&name)				
	if err != nil {
		// People Ident doesn't exist, 
		// so can create a new People profile.
	    if err == pgx.ErrNoRows { 
			if people.ProfilePicData == "" {
                sql :=
                    `insert into kkm.people
                    (
                        ident, name, gender, dob, nationality, race,
                        tel, email, address, postalCode, locality,
                        district, state, eduLvl, occupation, comorbids, 
                        supportvac, password, role
                    )
                    values
                    (
                        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
                        $11, $12, $13, $14, $15, $16, $17, $18, $19
                    )`
        
                _, err = conn.Exec(context.Background(), sql, 
                    people.Ident, people.Name, people.Gender, people.Dob, 
                    people.Nationality, people.Race, people.Tel, people.Email, 
                    people.Address, people.PostalCode, people.Locality, 
                    people.District, people.State, people.EduLvl, 
                    people.Occupation, people.Comorbids, people.SupportVac, 
                    auth.DEFAULT_PEOPLE_PWD, people.Role)
            } else {    
                sql :=
                    `insert into kkm.people
                    (
                        ident, name, gender, dob, nationality, race,
                        tel, email, address, postalCode, locality,
                        district, state, eduLvl, occupation, comorbids, 
                        supportvac, password, profilepic, role
                    )
                    values
                    (
                        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
                        $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
                    )`
        
                    _, err = conn.Exec(context.Background(), sql, 
                        people.Ident, people.Name, people.Gender, people.Dob, 
                        people.Nationality, people.Race, people.Tel, people.Email, 
                        people.Address, people.PostalCode, people.Locality, 
                        people.District, people.State, people.EduLvl, 
                        people.Occupation, people.Comorbids, people.SupportVac, 
                        auth.DEFAULT_PEOPLE_PWD, people.ProfilePicData, 
                        people.Role)
            }
            // New People profile create failed
            if err != nil {
                return "", err
            }
            // New People profile created successfully
            return "1", nil
		} 
		// Other unknown error during database scan.
		return "", err
	} 

    // People profile already exists in database,
    // so no new profile created. 
    return "0", nil   
}

type CreateNewPeopleHttpRespCode struct {
	CreateNewPeopleRespCode string	`json:"createNewPeopleRespCode"`	
}

func CreateNewPeopleHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[CreateNewPeopleHandler] request received")

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }

    // DECODING
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", people)
    
    // CREATE NEW PEOPLE
    db.CheckDbConn()
    createNewPeopleResult, err := CreateNewPeople(db.Conn, people)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }   
    
    // RESULT
    createNewPeopleRespCode := CreateNewPeopleHttpRespCode {
		CreateNewPeopleRespCode: createNewPeopleResult,
	}
	createNewPeopleRespJson, err := json.MarshalIndent(createNewPeopleRespCode, "", "")
	if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    } 
	fmt.Fprintf(w, "%s", createNewPeopleRespJson)
}

func GetPeopleProfile(conn *pgxpool.Pool, ident string) ([]byte, error) {   
    rows, err := conn.Query(context.Background(), 
        `select people.name, people.gender, people.dob::text, 
         people.nationality, people.race, people.tel, people.email,
         people.address, people.postalcode, people.locality, 
         people.district, people.state, people.eduLvl, 
         people.occupation, people.comorbids, people.supportVac, 
         coalesce(people.profilepic, '') as profilepic, people.role,
         coalesce(vaccine.brand, '') as brand, 
         coalesce(vaccine.type, '') as type, 
         coalesce(vaccine.against, '') as against, 
         coalesce(vaccine.raoa, '') as raoa, 
         coalesce( vaccination.id, -1) as id, 
         coalesce(vaccination.vaccination, '') as vaccination, 
         coalesce(vaccination.firstAdm::text, '') as firstAdm,                    
         coalesce(vaccination.firstDoseDt::text, '') as firstDoseDt, 
         coalesce(vaccination.secondDoseDt::text, '') as secondDoseDt, 
         coalesce(vaccination.aefiClass::text, '') as aefiClass, 
         coalesce(vaccination.aefiReaction, '{}') as aefiReaction, 
         coalesce(vaccination.remarks, '') as remarks
         from kkm.people 
             left join kkm.vaccination 
               on kkm.people.ident = kkm.vaccination.people
             left join kkm.vaccine
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
        var vaccinationId int64
        var vaccination string          
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
            var profilePicData string 
            var role string

            err = rows.Scan(&name, &gender, &dob, &nationality, &race, &tel, 
                &email, &address, &postalCode, &locality, &district, &state, 
                &eduLvl, &occupation, &comorbids, &supportVac, &profilePicData, 
                &role,
                &brand, &vacType, &against, &raoa, 
                &vaccinationId, &vaccination, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)
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
                ProfilePicData: profilePicData,
                Role: role,
            }
            firstRecord = false                                     
        } else {
            err = rows.Scan(nil, nil, nil, nil, nil, nil, nil,
                nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
                &brand, &vacType, &against, &raoa, 
                &vaccinationId, &vaccination, &fa, &fdd, &sdd, &aefiClass, &aefiReaction, &remarks)                      
            if err != nil {
                return nil, err
            }
        }
        if vaccinationId != -1 {        
            vaccinationRecord := VaccinationRecord{
                VaccinationId: vaccinationId,
                Vaccination: vaccination,
                VaccineBrand: brand,
                VaccineType: vacType,
                VaccineAgainst: against,
                VaccineRaoa: raoa,
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
    }

    outputJson, err := json.MarshalIndent(peopleProfile, "", "\t")        
    return outputJson, err
}

func GetPeopleProfile2(conn *pgxpool.Pool, ident string) ([]byte, error) {   
    rows, err := conn.Query(context.Background(), 
        `select p.name, p.gender, p.dob::text, 
         p.nationality, p.race, p.tel, p.email,
         p.address, p.postalcode, p.locality, 
         p.district, p.state, p.eduLvl, 
         p.occupation, p.comorbids, p.supportVac, 
         coalesce(p.profilepic, '') as profilepic, p.role,
         coalesce(fdv.brand, '') as fdBrand, 
         coalesce(sdv.brand, '') as sdBrand,           
         coalesce(v.id, -1) as id, 
         coalesce(v.fdtca::text, '') as fdTCA, 
         coalesce(v.fdgiven::text, '') as fdGiven, 
         coalesce(v.sdtca::text, '') as sdTCA, 
         coalesce(v.sdgiven::text, '') as sdGiven, 
         coalesce(v.fdaeficlass::text, '') as fdAefiClass, 
         coalesce(v.sdaeficlass::text, '') as sdAefiClass,
         coalesce(v.fdaefireaction, '{}') as fdAefiReaction, 
         coalesce(v.sdaefireaction, '{}') as sdAefiReaction, 
         coalesce(v.fdremarks, '') as fdRemarks,
         coalesce(v.sdremarks, '') as sdRemarks
         from kkm.people p
             left join kkm.vaccination v
               on p.ident = v.people
             left join kkm.vaccine fdv
               on v.vaccine = fdv.id
             left join kkm.vaccine sdv
               on v.vaccine = sdv.id
         where p.ident=$1`,
        ident)
    if err != nil {
        return nil, err
    }
    // var peopleProfile PeopleProfile
    var peopleProfile PeopleProfile2
    firstRecord := true

    for rows.Next() {
        // Vaccine
        var fdBrand string 
        var sdBrand string 
        // Vaccination
        var vaccinationId int64
        var vaccination string          
        var fdTCA string
        var fdGiven string
        var sdTCA string
        var sdGiven string
        var fdAefiClass string 
        var sdAefiClass string 
        var fdAefiReaction []string 
        var sdAefiReaction []string 
        var fdRemarks string 
        var sdRemarks string 

        if firstRecord {
            // People
            var name string
            var gender string
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
            var profilePicData string 
            var role string

            err = rows.Scan(&name, &gender, &dob, &nationality, &race, &tel, 
                &email, &address, &postalCode, &locality, &district, &state, 
                &eduLvl, &occupation, &comorbids, &supportVac, &profilePicData, 
                &role,
                &fdBrand, &sdBrand, &vaccinationId, 
                &fdTCA, &fdGiven, &sdTCA, &sdGiven, 
                &fdAefiClass, &sdAefiClass, &fdAefiReaction, &sdAefiReaction,
                &fdRemarks, &sdRemarks)
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
                ProfilePicData: profilePicData,
                Role: role,
            }
            firstRecord = false                                     
        } else {
            err = rows.Scan(nil, nil, nil, nil, nil, nil, nil,
                nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
                &fdBrand, &sdBrand, &vaccinationId, 
                &fdTCA, &fdGiven, &sdTCA, &sdGiven, 
                &fdAefiClass, &sdAefiClass, &fdAefiReaction, &sdAefiReaction,
                &fdRemarks, &sdRemarks)                      
            if err != nil {
                return nil, err
            }
        }
        if vaccinationId != -1 {        
            vaccinationRecord := VaccinationRecord2{
                VaccinationId: vaccinationId,
                Vaccination: vaccination,
                FdVaccineBrand: fdBrand,
                SdVaccineBrand: sdBrand,
                FdTCA: fdTCA,
                FdGiven: fdGiven,
                SdTCA: sdTCA,
                SdGiven: sdGiven,
                FdAefiClass: fdAefiClass,
                SdAefiClass: sdAefiClass,
                FdAefiReaction: fdAefiReaction,
                SdAefiReaction: sdAefiReaction,
                FdRemarks: fdRemarks,
                SdRemarks: sdRemarks,
            }
            peopleProfile.VaccinationRecords = append(
                peopleProfile.VaccinationRecords,
                vaccinationRecord,
            )
        }
    }

    outputJson, err := json.MarshalIndent(peopleProfile, "", "\t")        
    return outputJson, err
}

func GetPeopleHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[GetPeopleHandler] request received")    

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }   

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    
    db.CheckDbConn()
    peopleProfJson, err := GetPeopleProfile2(db.Conn, identity.Ident)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Print("People entry not found in database")
        }
        util.SendInternalServerErrorStatus(w, err)
        return 
    }
    fmt.Printf("%s\n", peopleProfJson)
    fmt.Fprintf(w, "%s", peopleProfJson)
}

func UpdatePeople(conn *pgxpool.Pool, people People) error {  
    var err error

    if people.ProfilePicData == "" {
        sql := 
        `update kkm.people 
           set name=$1, gender=$2, dob=$3, nationality=$4, race=$5, 
             tel=$6, email=$7, address=$8, postalCode=$9, locality=$10,
             district=$11, state=$12, eduLvl=$13, occupation=$14, 
             comorbids=$15, supportVac=$16, role=$17 
           where ident=$18`   

        _, err = conn.Exec(context.Background(), sql,
            people.Name, people.Gender, people.Dob, people.Nationality, 
            people.Race, people.Tel, people.Email, people.Address, people.PostalCode, 
            people.Locality, people.District, people.State, people.EduLvl, 
            people.Occupation, people.Comorbids, people.SupportVac, people.Role,
            people.Ident)
    } else {
        sql := 
        `update kkm.people 
           set name=$1, gender=$2, dob=$3, nationality=$4, race=$5, 
             tel=$6, email=$7, address=$8, postalCode=$9, locality=$10,
             district=$11, state=$12, eduLvl=$13, occupation=$14, 
             comorbids=$15, supportVac=$16, profilepic=$17, role=$18 
           where ident=$19`   

        _, err = conn.Exec(context.Background(), sql,
            people.Name, people.Gender, people.Dob, people.Nationality, 
            people.Race, people.Tel, people.Email, people.Address, people.PostalCode, 
            people.Locality, people.District, people.State, people.EduLvl, 
            people.Occupation, people.Comorbids, people.SupportVac, people.ProfilePicData,
            people.Role, people.Ident)
    }    
    if err != nil {
        return err
    }    
    return nil
}

func UpdatePeopleHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[UpdatePeopleHandler] request received")

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }    
    
    /* LESS-EFFICIENT-JSON_DECODING-METHOD (Produces intermediate byte slice)
       var people People
       err := json.Unmarshal([]byte(input), &people) */

    /* MORE-EFFICIENT-JSON_DECODING-METHOD (No intermediate byte slice) */
    var people People
    err := json.NewDecoder(r.Body).Decode(&people)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", people)

    db.CheckDbConn()
    err = UpdatePeople(db.Conn, people)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }   
}

func DeletePeople(conn *pgxpool.Pool, identity Identity) error {
    sql := `delete from kkm.people 
            where ident=$1`

    _, err := conn.Exec(context.Background(), sql, identity.Ident)
    if err != nil {
        return err
    }
    return nil
}

func DeletePeopleHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[DeletePeopleHandler] Request form data received")

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }

    var identity Identity
    err := json.NewDecoder(r.Body).Decode(&identity)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%v\n", identity)

    db.CheckDbConn()
    err = DeletePeople(db.Conn, identity)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
}

func CreateNewVacRec(conn *pgxpool.Pool, vru VacRecUpsert) error {                  
    var err error
    if vru.VacRec.Fdd == "" {
        sql := 
            `insert into kkm.vaccination
            (
                vaccine, people, vaccination, firstadm,  
                seconddosedt, aeficlass, aefireaction, remarks
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
                vaccine, people, vaccination, firstadm, firstdosedt, 
                aeficlass, aefireaction, remarks
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
                vaccine, people, vaccination, firstadm,  
                aeficlass, aefireaction, remarks
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
                vaccine, people, vaccination, firstadm, firstdosedt, 
                seconddosedt, aeficlass, aefireaction, remarks
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

func CreateNewVacRec2(conn *pgxpool.Pool, vru VacRecUpsert2) error {                  
    var err error
    // if vru.VacRec.Fdd == "" {
    //     sql := 
    //         `insert into kkm.vaccination
    //         (
    //             fdvaccine, people, vaccination,   
    //             fdtca, fdgiven, fdaeficlass, fdaefireaction, fdremarks, 
    //             sdtca, sdgiven, sdaeficlass, sdaefireaction, sdremarks, 
    //             sdvaccine
    //         )
    //         select fdv.id, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
    //         from kkm.vaccine fdv
    //         where fdv.brand=$13
    //         union all
    //         select sdv.id
    //         from kkm.vaccine sdv
    //         where sdv.brand=$14` 

    //     _, err = conn.Exec(context.Background(), sql, 
    //     vru.Ident, vru.VacRec.Vaccination, 
    //     vru.VacRec.FdTCA, vru.VacRec.FdGiven, vru.VacRec.FdAefiClass, vru.VacRec.FdAefiReaction,
    //     vru.VacRec.FdRemarks,
    //     vru.VacRec.SdTCA, vru.VacRec.SdGiven, vru.VacRec.SdAefiClass, vru.VacRec.SdAefiReaction,
    //     vru.VacRec.SdRemarks,
    //     vru.VacRec.FdVaccineBrand,
    //     vru.VacRec.SdVaccineBrand)
    // } else if vru.VacRec.Sdd == "" {

    if ((vru.VacRec.FdTCA != "" && vru.VacRec.FdGiven != "") && (vru.VacRec.SdTCA == "" || vru.VacRec.SdGiven == "")) {
        sql := 
            `insert into kkm.vaccination
            (
                fdvaccine, people, vaccination,   
                fdtca, fdgiven, fdaeficlass::text, fdaefireaction, fdremarks               
            )
            select fdv.id, $1, $2, $3, $4, $5, $6, $7
            from kkm.vaccine fdv
            where fdv.brand=$8`

        _, err = conn.Exec(context.Background(), sql, 
        vru.Ident, vru.VacRec.Vaccination, 
        vru.VacRec.FdTCA, vru.VacRec.FdGiven, vru.VacRec.FdAefiClass, vru.VacRec.FdAefiReaction,
        vru.VacRec.FdRemarks,
        vru.VacRec.FdVaccineBrand)
    } else if vru.VacRec.FdGiven == "" {
        sql := 
            `insert into kkm.vaccination
            (
                fdvaccine, people, vaccination,   
                fdtca, fdaeficlass, fdaefireaction, fdremarks
            )
            select fdv.id, $1, $2, $3, $4, $5, $6
            from kkm.vaccine fdv
            where fdv.brand=$7`  

        _, err = conn.Exec(context.Background(), sql, 
            vru.Ident, vru.VacRec.Vaccination, 
            vru.VacRec.FdTCA, vru.VacRec.FdAefiClass, vru.VacRec.FdAefiReaction,
            vru.VacRec.FdRemarks,
            vru.VacRec.FdVaccineBrand)
    } else {
        sql := 
            `insert into kkm.vaccination
            (
                fdvaccine, people, vaccination,   
                fdtca, fdgiven, fdaeficlass, fdaefireaction, fdremarks, 
                sdtca, sdgiven, sdaeficlass, sdaefireaction, sdremarks, 
                sdvaccine
            )
            select fdv.id, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
            from kkm.vaccine fdv
            where fdv.brand=$13
            union all
            select sdv.id
            from kkm.vaccine sdv
            where sdv.brand=$14`  

        _, err = conn.Exec(context.Background(), sql, 
            vru.Ident, vru.VacRec.Vaccination, 
            vru.VacRec.FdTCA, vru.VacRec.FdGiven, vru.VacRec.FdAefiClass, vru.VacRec.FdAefiReaction,
            vru.VacRec.FdRemarks,
            vru.VacRec.SdTCA, vru.VacRec.SdGiven, vru.VacRec.SdAefiClass, vru.VacRec.SdAefiReaction,
            vru.VacRec.SdRemarks,
            vru.VacRec.FdVaccineBrand,
            vru.VacRec.SdVaccineBrand)
    }    
    if err != nil {
        return err
    }
    return nil
}

func CreateNewVacRecHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[CreateNewVacRecHandler] request received")

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }
        
    // var vru VacRecUpsert
    var vru VacRecUpsert2
    err := json.NewDecoder(r.Body).Decode(&vru)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", vru)

    db.CheckDbConn()
    err = CreateNewVacRec2(db.Conn, vru)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }   
}

func UpdateVacRec(conn *pgxpool.Pool, vru VacRecUpsert) error {
    var err error
    if vru.VacRec.Fdd == "" {
        sql := 
            `update kkm.vaccination
            set vaccine=subq.id, firstadm=$1,
                seconddosedt=$2, aeficlass=$3, aefireaction=$4,
                remarks=$5
            from (select vac.id 
                from kkm.vaccine vac
                where vac.brand=$6) as subq
            where kkm.vaccination.id=$7`

        _, err = conn.Exec(context.Background(), sql, 
            vru.VacRec.Fa, vru.VacRec.Sdd,
            vru.VacRec.AefiClass, vru.VacRec.AefiReaction,
            vru.VacRec.Remarks, vru.VacRec.VaccineBrand, 
            vru.VacRec.VaccinationId)
    } else if vru.VacRec.Sdd == "" {
        sql := 
            `update kkm.vaccination
            set vaccine=subq.id, firstadm=$1, firstdosedt=$2,
                aeficlass=$3, aefireaction=$4,
                remarks=$5
            from (select vac.id 
                from kkm.vaccine vac
                where vac.brand=$6) as subq
            where kkm.vaccination.id=$7`

        _, err = conn.Exec(context.Background(), sql, 
            vru.VacRec.Fa, vru.VacRec.Fdd,
            vru.VacRec.AefiClass, vru.VacRec.AefiReaction,
            vru.VacRec.Remarks, vru.VacRec.VaccineBrand, 
            vru.VacRec.VaccinationId)

    } else if vru.VacRec.Fdd == "" && vru.VacRec.Sdd == "" {
        sql := 
            `update kkm.vaccination
            set vaccine=subq.id, firstadm=$1,
                aeficlass=$2, aefireaction=$3,
                remarks=$4
            from (select vac.id 
                from kkm.vaccine vac
                where vac.brand=$5) as subq
            where kkm.vaccination.id=$6`

        _, err = conn.Exec(context.Background(), sql, 
            vru.VacRec.Fa,
            vru.VacRec.AefiClass, vru.VacRec.AefiReaction,
            vru.VacRec.Remarks, vru.VacRec.VaccineBrand, 
            vru.VacRec.VaccinationId)
        
    } else {
        sql := 
            `update kkm.vaccination
            set vaccine=subq.id, firstadm=$1, firstdosedt=$2,
                seconddosedt=$3, aeficlass=$4, aefireaction=$5,
                remarks=$6
            from (select vac.id 
                from kkm.vaccine vac
                where vac.brand=$7) as subq
            where kkm.vaccination.id=$8`

        _, err = conn.Exec(context.Background(), sql, 
            vru.VacRec.Fa, vru.VacRec.Fdd, vru.VacRec.Sdd,
            vru.VacRec.AefiClass, vru.VacRec.AefiReaction,
            vru.VacRec.Remarks, vru.VacRec.VaccineBrand, 
            vru.VacRec.VaccinationId)
    }
    if err != nil {
        return err
    }
    return nil
}

func UpdateVacRec2(conn *pgxpool.Pool, vru VacRecUpsert2) error {
    var err error
    // if vru.VacRec.Fdd == "" {
    //     sql := 
    //         `update kkm.vaccination
    //         set vaccine=subq.id, firstadm=$1,
    //             seconddosedt=$2, aeficlass=$3, aefireaction=$4,
    //             remarks=$5
    //         from (select vac.id 
    //             from kkm.vaccine vac
    //             where vac.brand=$6) as subq
    //         where kkm.vaccination.id=$7`

    //     _, err = conn.Exec(context.Background(), sql, 
    //         vru.VacRec.Fa, vru.VacRec.Sdd,
    //         vru.VacRec.AefiClass, vru.VacRec.AefiReaction,
    //         vru.VacRec.Remarks, vru.VacRec.VaccineBrand, 
    //         vru.VacRec.VaccinationId)
    // } else if vru.VacRec.Sdd == "" {

    if vru.VacRec.SdTCA == "" || vru.VacRec.SdGiven == "" {
        sql := 
            `update kkm.vaccination
            set fdvaccine=subqFd.id, 
                fdtca=$1, fdgiven=$2, fdaeficlass=$3, fdaefireaction=$4,
                fdremarks=$5               
            from (select vac.id 
                from kkm.vaccine vac
                where vac.brand=$6) as subqFd
            where kkm.vaccination.id=$7`

        _, err = conn.Exec(context.Background(), sql, 
            vru.VacRec.FdTCA, vru.VacRec.FdGiven, vru.VacRec.FdAefiClass, vru.VacRec.FdAefiReaction,
            vru.VacRec.FdRemarks, 
            vru.VacRec.FdVaccineBrand, vru.VacRec.VaccinationId)

    } else {
        sql := 
        `update kkm.vaccination
        set fdvaccine=subq.fdId, sdvaccine=subq.sdId
            fdtca=$1, fdgiven=$2, fdaeficlass=$3, fdaefireaction=$4,
            fdremarks=$5,
            sdtca=$6, sdgiven=$7, sdaeficlass=$8, sdaefireaction=$9,
            sdremarks=$10                
        from (select fdVac.id as fdId
            from kkm.vaccine fdVac
            where fdVac.brand=$11
            union all
            select sdVac.id as sdId
            from kkm.vaccine sdVac
            where sdVac.brand=$12
            ) as subq
        where kkm.vaccination.id=$13`

        _, err = conn.Exec(context.Background(), sql, 
            vru.VacRec.FdTCA, vru.VacRec.FdGiven, vru.VacRec.FdAefiClass, vru.VacRec.FdAefiReaction,
            vru.VacRec.FdRemarks,
            vru.VacRec.SdTCA, vru.VacRec.SdGiven, vru.VacRec.SdAefiClass, vru.VacRec.SdAefiReaction,
            vru.VacRec.SdRemarks, 
            vru.VacRec.FdVaccineBrand, vru.VacRec.SdVaccineBrand, vru.VacRec.VaccinationId)
    }
    if err != nil {
        return err
    }
    return nil
}
 
func UpdateVacRecHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[UpdateVacRecHandler] request received")

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }
        
    // var vru VacRecUpsert
    var vru VacRecUpsert2
    err := json.NewDecoder(r.Body).Decode(&vru)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", vru)

    db.CheckDbConn()
    err = UpdateVacRec2(db.Conn, vru)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }   
}

func DeleteVacRec(conn *pgxpool.Pool, vacRecId int64) error {
    sql := `delete from kkm.vaccination where id=$1`

    _, err := conn.Exec(context.Background(), sql, vacRecId)
    if err != nil {
        return err
    }
    return nil
}

func DeleteVacRecHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[DeleteVacRecHandler] request received")

    // VERIFY AUTH TOKEN
    authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    if !auth.VerifyTokenHMAC(authToken) {
        util.SendUnauthorizedStatus(w)
        return
    }
        
    var vrd VacRecDelete
    err := json.NewDecoder(r.Body).Decode(&vrd)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    fmt.Printf("%+v\n", vrd)

    db.CheckDbConn()
    err = DeleteVacRec(db.Conn, vrd.VaccinationId)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }   
}