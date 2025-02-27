package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func PingSystem(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the list of tables
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Tables in database:")

	// Loop through results and print table names
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		fmt.Println("- " + tableName)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// func InsertTable(w http.ResponseWriter, r *http.Request) {
// 	db, err := sql.Open("sqlite", "test.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	createTableSQL := `
// 	CREATE TABLE IF NOT EXISTS green (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		title TEXT NOT NULL,
// 		description TEXT,
// 		is_complete INTEGER NOT NULL DEFAULT 0
// 	);`

// 	// Execute the query
// 	_, err = db.Exec(createTableSQL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func setupSystemRouter() http.Handler {
	mux := http.NewServeMux()

	// Add your handlers to the mux
	mux.Handle("/", http.HandlerFunc(PingSystem))
	mux.Handle("/get-todos", http.HandlerFunc(GetTodos))
	// mux.Handle("/insert-table", http.HandlerFunc(InsertTable))

	// Wrap the mux with the CORS middleware
	return addCORSHeaders(mux)
}

func addCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust the origin as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		handler.ServeHTTP(w, r)
	})
}

func main() {

	fmt.Println("Hello worlds!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	systemRouter := setupSystemRouter()

	systemServer := &http.Server{
		Addr:    port,
		Handler: systemRouter,
	}

	go func() {
		fmt.Println("Starting server interface on port 9090")
		if err := systemServer.ListenAndServe(); err != nil {
			fmt.Println("Server 2 error:", err)
		}
	}()

	select {}
}
