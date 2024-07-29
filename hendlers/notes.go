package hendlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"password/model"
	"password/security"
	"strconv"

	_ "github.com/lib/pq"
)

type NoteResponse struct {
	Text string `json:"text"`
}

type NotesServer struct {
	Counter   int
	NoteModel model.NotesModel
}

func (nt *NotesServer) PassPost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var n model.Note
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		log.Print(err)
		return
	}
	encryptText, err := security.Encrypt(n.Text, n.Password)
	if err != nil {
		fmt.Println("notes, encrypt", err)
	}
	n.Text = encryptText

	hashedPassword, err := security.HashPassword(n.Password)
	if err != nil {
		fmt.Println("hashedPassword", err)
	}
	n.Password = hashedPassword

	nt.Counter++
	n.ID = nt.Counter
	nt.NoteModel.Insert(&n)

	err = json.NewEncoder(w).Encode(nt.Counter)
	if err != nil {
		log.Print(err)
		return
	}

}

func (nt *NotesServer) PassGet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query().Get("id")
	paramsInt, err := strconv.Atoi(params)
	if err != nil {
		log.Print(err)
		return
	}

	note, err := nt.NoteModel.GetByID(paramsInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	password := r.URL.Query().Get("password")
	if security.CheckPasswordHash(password, note.Password) {
		var jsonResponse NoteResponse
		decryptedText, err := security.Decrypt(note.Text, password)
		if err != nil {
			fmt.Println("decryptedText error", err)
		}

		note.Text = decryptedText
		jsonResponse.Text = note.Text
		err = json.NewEncoder(w).Encode(jsonResponse)
		if err != nil {
			log.Print(err)
			return
		}

	} else {
		w.WriteHeader(http.StatusExpectationFailed)
	}
	nt.NoteModel.GetAll()

}
