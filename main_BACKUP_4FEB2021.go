package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"github.com/jackc/pgx/v4"
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

    err = listPeople()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to list people: %v\n", err)
	os.Exit(1)
    }
}

func listPeople() error {
    rows, _ := conn.Query(context.Background(), "select ident, name, dob from kkm.people")

    for rows.Next() {
        var ident string
	var name string
	var dob string
	err := rows.Scan(&ident, &name, &dob)
	if err != nil {
            return err
	}
        fmt.Printf("ID: %s, Name: %s, DOB: %s", ident, name, dob)
    }

    return rows.Err()
}


