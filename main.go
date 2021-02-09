package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"errors"
	"strconv"
	"strings"
	"github.com/jackc/pgx"
	"myVaksin/data"
)

const (
    DATE_ISO =  "2006-01-02"
)

var conn *pgx.Conn

func main() {
    var err error
    conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to make connection to database: %v\n", err)
		os.Exit(1)
    }
    defer conn.Close(context.Background())
  
    if len(os.Args) == 1 {
        fmt.Println("Arguments needed to call this program!")
		printUsage()
		os.Exit(1)
    }

    switch os.Args[1] {

    case "listPeople":
        err = listPeople()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to list people: %v\n", err)
	    os.Exit(1)
        }

    case "updatePeople":
		if len(os.Args) != 5 {
			fmt.Println("Incorrect number of arguments!")
			printUsage()
		}
	
        ident := os.Args[2]
        key := os.Args[3]
        value := os.Args[4]

		err = updatePeople(ident, key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to update people: %v\n", err)
			os.Exit(1)
		}

    case "getPeoples":
        output, err := data.GetPeoples(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get list of people: %v\n", err)
		}
		fmt.Printf("%s", output)

	}
}

func listPeople() error {
    rows, _ := conn.Query(context.Background(), "select ident, name, dob, tel, address, race, nationality, edu_lvl, occupation, comorbids, support_vac from kkm.people")

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
		return err
	}
	fmt.Printf("Ident: %v, Name: %v, DOB: %v, Tel: %v, Add: %v, " +
		"Race: %v, Nationality: %v, EduLvl: %v, Occupation: %v, " + 
		"Comorbids: %v, Support_vac: %v\n\n", 
		ident, name, dob, tel, address, race, nationality, edu_lvl, 
		occupation, comorbids, support_vac)
    }

    return rows.Err()
}

func updatePeople(ident string, key string, value string) error {
    sql := fmt.Sprintf("update kkm.people set %s=$1 where ident=$2", key)

    switch key {
    
    case "dob":
	    dob, _ := time.Parse(DATE_ISO, value)
        _, err := conn.Exec(context.Background(), sql, dob, ident)
        return err

    case "comorbids":
	var comorbids []int64
	strArr := strings.Split(value, ",")
	for _, s := range strArr {
	    n, err := strconv.ParseInt(s, 10, 32)
	    if err != nil {
			return err
	    }
            comorbids = append(comorbids, n)
	}
        _, err := conn.Exec(context.Background(), sql, comorbids, ident)
	return err

    case "support_vac":
	support_vac, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
        _, err = conn.Exec(context.Background(), sql, support_vac, ident)
	return err

    default:
        _, err := conn.Exec(context.Background(), sql, value, ident)
        return err
    }

    return errors.New("No matching update key!")
}

func printUsage() {
    fmt.Print(`
    Postgresql CRUD Demo Program

    Usage:

    go run main.go listPeople
    go run main.go updatePeople <ident> <key> <value>`)
}
