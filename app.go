package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	ID   uint32 `json:"id"`
	Note string `json:"note"`
}

type Session struct {
	UserEmail string
}

type Database struct {
	Users    map[string]User
	Sessions map[string]Session
	Notes    map[string][]Note
	NoteID   uint32
	sync.Mutex
}

var db = &Database{
	Users:    make(map[string]User),
	Sessions: make(map[string]Session),
	Notes:    make(map[string][]Note),
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	db.Lock()
	defer db.Unlock()

	// Check if the user already exists
	if _, exists := db.Users[user.Email]; exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Store the user in the database
	db.Users[user.Email] = user

	w.WriteHeader(http.StatusOK)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	db.Lock()
	defer db.Unlock()

	// Check if the user exists
	user, exists := db.Users[login.Email]
	if !exists || user.Password != login.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a new session
	sessionID := fmt.Sprintf("session_%d", db.NoteID)
	db.Sessions[sessionID] = Session{UserEmail: user.Email}

	response := map[string]string{"sid": sessionID}
	jsonResponse(w, http.StatusOK, response)
}

func listNotesHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sid")

	db.Lock()
	defer db.Unlock()

	// Check if the session is valid
	session, exists := db.Sessions[sessionID]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve notes for the user associated with the session
	userNotes := db.Notes[session.UserEmail]

	response := map[string][]Note{"notes": userNotes}
	jsonResponse(w, http.StatusOK, response)
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note struct {
		SID  string `json:"sid"`
		Note string `json:"note"`
	}
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	db.Lock()
	defer db.Unlock()

	// Check if the session is valid
	session, exists := db.Sessions[note.SID]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create a new note
	newNote := Note{ID: db.NoteID, Note: note.Note}
	db.NoteID++
	db.Notes[session.UserEmail] = append(db.Notes[session.UserEmail], newNote)

	response := map[string]uint32{"id": newNote.ID}
	jsonResponse(w, http.StatusOK, response)
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	var deleteNote struct {
		SID string `json:"sid"`
		ID  uint32 `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&deleteNote)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	db.Lock()
	defer db.Unlock()

	// Check if the session is valid
	session, exists := db.Sessions[deleteNote.SID]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find and delete the note
	userNotes := db.Notes[session.UserEmail]
	var updatedNotes []Note
	for _, n := range userNotes {
		if n.ID != deleteNote.ID {
			updatedNotes = append(updatedNotes, n)
		}
	}
	db.Notes[session.UserEmail] = updatedNotes

	w.WriteHeader(http.StatusOK)
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listNotesHandler(w, r)
		case http.MethodPost:
			createNoteHandler(w, r)
		case http.MethodDelete:
			deleteNoteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
