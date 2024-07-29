package model

import (
	"database/sql"
	"fmt"
)

type Note struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Password string `json:"password"`
}
type NotesModel struct {
	DB *sql.DB
}

func (nb *NotesModel) Insert(note *Note) {
	query := `insert into notes(id, note, password) values($1,$2,$3)`
	_, err := nb.DB.Exec(query, note.ID, note.Text, note.Password)
	fmt.Println(note.ID, note.Text, note.Password)
	if err != nil {
		fmt.Print(err)

	}

}

func (nb *NotesModel) GetAll() ([]Note, error) {
	query := `select id, note,password from notes`
	row, err := nb.DB.Query(query)
	if err != nil {
		fmt.Println("GetAll, query erorr", err)
		return []Note{}, err
	}
	var notes []Note

	for row.Next() {
		var note Note
		err := row.Scan(&note.ID, &note.Text, &note.Password)
		if err != nil {
			fmt.Println("GetAll, scan erorr", err)
			return []Note{}, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (nb *NotesModel) GetByID(userId int) (Note, error) {
	var note Note
	query := `select id, note,password from notes where id = $1`
	row := nb.DB.QueryRow(query, userId)
	err := row.Scan(&note.ID, &note.Text, &note.Password)
	if err != nil {
		fmt.Println("GetById, scan erorr", err)
		return Note{}, err
	}
	return note, nil

}

func (nb *NotesModel) GetLastRowFromDB() (int, error) {
	var id int
	query := `SELECT coalesce(MAX(id), 0) FROM notes`
	row := nb.DB.QueryRow(query)
	err := row.Scan(&id)
	if err != nil {
		fmt.Println("GetLastRowFromDB, scan erorr", err)
		return 0, err
	}
	return id, nil

}
