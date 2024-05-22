package storage

import (
	// "context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)


type PostgresStorage struct {
	Client *sql.DB
}

func NewPostgresDb(psqlInfo string) *PostgresStorage {

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully started a supabase client!!")

	return &PostgresStorage{
		Client: db,
	}
}

func (s *PostgresStorage) GetLinks(username string, tags []string) []string {

	
	rows, err := s.Client.Query("SELECT link, tags FROM links WHERE username = $1 AND tags ?| $2::text[];", username, pq.Array(tags))
	
	if err != nil {
		panic(err)
	}
	
	defer rows.Close()
	
	var results []string

	for rows.Next() {
		var link string
		var tagsJson []byte
		err := rows.Scan(&link, &tagsJson); 

		if err != nil {
				panic(err)
		}

		var tags []string
		err = json.Unmarshal(tagsJson, &tags)

		if err != nil {
			panic(err)
		}
		
		results = append(results, link)
	}

	err = rows.Err()

	if err != nil {
    panic(err)
  }

	fmt.Println("\nlinks to send: ", results)
	return results
}


func (s *PostgresStorage) InsertLinkAndTags(username string, link string, tags []string) error {

	tagsJSON, err := json.Marshal(tags)

	if err != nil {
		fmt.Println("error when marshaling tagsJSON")
		return errors.New("something went wrong, try again later")
	}

	_, err = s.Client.Exec("INSERT INTO links (username, link, tags) VALUES ($1, $2, $3)", username, link, tagsJSON)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key"){
			return errors.New("you can't save the same url twice")
		}
		return errors.New("something went wrong, try again later")
	}

	return nil

}

func (s *PostgresStorage) GetUserTags(username string) ([]string, error) {
	rows, err := s.Client.Query("SELECT DISTINCT jsonb_array_elements_text(tags) AS unique_tag FROM links WHERE username = $1;", username)

	if err != nil {
		fmt.Println("Error when retrieving user tags, error: ", err)
		return []string{}, errors.New("something went wrong, try again later")
	}

	defer rows.Close()

	var results []string

	for rows.Next(){
		var tag string

		err := rows.Scan(&tag)

		if err != nil {
			fmt.Println("Error ocurred when parsing tags rows")
			return []string{}, errors.New("something went wrong, try again later")
		}

		results = append(results, tag)

	}

	return results, nil
}