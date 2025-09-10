package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // blank import to register the driver
	"log"
)

var db *sql.DB

// type url_model struct {
// 	ID int
// 	long_url string
// 	short_url string
// }
func init(){
	var err error
	db, err =  sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	sql_command := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		long_url TEXT NOT NULL,
		short_url TEXT NOT NULL UNIQUE
	);`
	 _, err = db.Exec(sql_command)
	if err != nil {
	log.Fatalf("Error creating table: %q: %s\n", err, sql_command) // Log an error if table creation fails
	}
	
}

func UpdateDB(longUrl string, shortUrl string){
	_, err := db.Exec("INSERT INTO urls (long_url, short_url) VALUES (?, ?)", longUrl, shortUrl)
	if err != nil {
		log.Fatal(err)
	}

}

func Fetch(shortUrl string) string {
	var longUrl string
	err := db.QueryRow("SELECT long_url FROM urls WHERE short_url = ?", shortUrl).Scan(&longUrl)
	if err != nil {
		log.Fatal(err)
	}
	return longUrl
}
