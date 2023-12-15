package models

import (
	"database/sql"
	"log"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	var LastInsertId int

	sqlStatement := `INSERT INTO snippets (title,content,created,expires) VALUES($1,$2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP + $3 * INTERVAL '1 DAY') RETURNING id`
	err := sm.DB.QueryRow(sqlStatement, title, content, expires).Scan(&LastInsertId)

	if err != nil {
		log.Println(err.Error())
	}

	id := LastInsertId

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*Snippet, error) {
	row := &Snippet{}

	sqlStatement := `SELECT id,title,content,created,expires FROM snippets WHERE expires > CURRENT_TIMESTAMP AND id=$1`
	err := sm.DB.QueryRow(sqlStatement, id).Scan(&row.ID, &row.Title, &row.Content, &row.Created, &row.Expires)

	if err != nil {
		log.Println(err.Error())
	}

	return row, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {

	latestSnippets := []*Snippet{}

	sqlStatement := `SELECT id, title, content,created, expires FROM snippets WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10 `
	rows, err := sm.DB.Query(sqlStatement)

	if err != nil {
		log.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		row := &Snippet{}
		err := rows.Scan(&row.ID, &row.Title, &row.Content, &row.Created, &row.Expires)
		if err != nil {
			return nil, err
		}
		latestSnippets = append(latestSnippets, row)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return latestSnippets, nil
}
