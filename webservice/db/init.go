package db 

import (
    "net/http"
    "fmt"
    "log"
    "os"
    "context"
    // "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
)

/*
    ============================================
    Ver 1 - PostgreSQL Connection Initialization
    ============================================
*/
// func init() {
//     var err error
//     conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
//     if err != nil {
//         fmt.Fprintf(os.Stderr, "Unable to make connection to database: %v\n", err)
// 		os.Exit(1)
//     }     
// }
// func Close() {
//     conn.Close(context.Background())
// }

// var Conn *pgx.Conn
var Conn *pgxpool.Pool 

func Open() {
    var err error
    // Conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    Conn, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to make connection to database: %v\n", err)
		os.Exit(1)
    }     
}

func Close() {
    CheckDbConn()
    // Conn.Close(context.Background())
    Conn.Close()
}

func CheckDbConn() {
    if Conn == nil {
        fmt.Fprint(
            os.Stderr, 
            "DB connection is not initialized yet. Please initialize DB connection first with Open()\n")
		os.Exit(1)
    }
}

func LogErrAndSendBadReqStatus(w http.ResponseWriter, err error) {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusBadRequest) //Http status code: 400
}

func SendUnauthorizedStatus(w http.ResponseWriter) {
    http.Error(w, "Unauthorized Access", http.StatusUnauthorized) //Http status code: 401
}

func LogErrAndSendInternalServerErrorStatus(w http.ResponseWriter, err error) {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusInternalServerError) //Http status code: 500
}

