package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Shirts struct {
	ID    int    `json:"id"`
	Brand string `json:"brand"`
	Size  string `json:"size"`
}

var db *sql.DB

// 1. READ ALL (GET /shirts)
func getShirts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, brand, size FROM shirts ORDER BY id ASC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	shirts := []Shirts{}
	for rows.Next() {
		var s Shirts
		if err := rows.Scan(&s.ID, &s.Brand, &s.Size); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		shirts = append(shirts, s)
	}

	json.NewEncoder(w).Encode(shirts)
}

// 2. READ BY ID (GET /shirts/{id})
func getShirtByID(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	var s Shirts
	query := "SELECT id, brand, size FROM shirts WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&s.ID, &s.Brand, &s.Size)
	if err == sql.ErrNoRows {
		http.Error(w, "Data baju tidak ditemukan", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// 3. CREATE (POST /shirts)
func createShirt(w http.ResponseWriter, r *http.Request) {
	var newShirt Shirts

	if err := json.NewDecoder(r.Body).Decode(&newShirt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO shirts (brand, size) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, newShirt.Brand, newShirt.Size).Scan(&newShirt.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newShirt)
}

// 4. UPDATE (PUT /shirts/{id})
func updateShirt(w http.ResponseWriter, r *http.Request, id int) {
	var updatedShirt Shirts
	if err := json.NewDecoder(r.Body).Decode(&updatedShirt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE shirts SET brand = $1, size = $2 WHERE id = $3"
	res, err := db.Exec(query, updatedShirt.Brand, updatedShirt.Size, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data baju dengan ID tersebut tidak ditemukan", http.StatusNotFound)
		return
	}

	updatedShirt.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedShirt)
}

// 5. DELETE (DELETE /shirts/{id})
func deleteShirt(w http.ResponseWriter, r *http.Request, id int) {
	query := "DELETE FROM shirts WHERE id = $1"
	res, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data baju dengan ID tersebut tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Baju dengan ID %d berhasil dihapus", id),
	})
}

// Handler untuk rute /shirts
func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getShirts(w, r)
	case "POST":
		createShirt(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler untuk rute /shirts/{id}
func shirtByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/shirts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID harus berupa angka", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getShirtByID(w, r, id)
	case "PUT":
		updateShirt(w, r, id)
	case "DELETE":
		deleteShirt(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// SESUAIKAN port dan password dengan pengaturan PostgreSQL Anda!
	connStr := "host=127.0.0.1 port=5433 user=postgres password=jawakontolanjing dbname=db_baju sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}
	fmt.Println("Berhasil terhubung ke PostgreSQL (db_baju)!")

	http.HandleFunc("/shirts", mainHandler)
	http.HandleFunc("/shirts/", shirtByIDHandler)

	fmt.Println("Server is running on port 8081...")
	http.ListenAndServe(":8081", nil)
}
