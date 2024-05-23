package storage

import (
	// "context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"wisdom-hoard/util"

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

func (s *PostgresStorage) GetLinks(username string, tags []string) ([]string) {

	
	rows, err := s.Client.Query("SELECT DISTINCT link, tags FROM links WHERE username = $1 AND tags ?| $2::text[];", username, pq.Array(tags))
	
	if err != nil {
		fmt.Println("Something went wrong when fetching links")
		return nil
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


func (s *PostgresStorage) InsertLinkAndTags(username string, link string, tags []string) ([]string, error) {

	tagsJSON, err := json.Marshal(tags)

	if err != nil {
		fmt.Println("error when marshaling tagsJSON")
		return nil, errors.New("something went wrong, try again later")
	}

	repeatedLink, previousTags, error := s.GetLinksByUrl(username, link)

	fmt.Println("REpeated link: ", repeatedLink, "previous tags: ", previousTags)

	if error != nil {
		fmt.Println("Error when calling getlinksbyurl")
		return nil, errors.New("something went wrong")
	}

	var	isLinkRepeated bool = false

	if repeatedLink == link{
		isLinkRepeated = true
	}

	var mergedTags []string = util.MergeSlices(tags, previousTags)
	mergedTagsJSON, err := json.Marshal(mergedTags)

	if err != nil {
		return nil, errors.New("something went wrong")
	}

	if isLinkRepeated{
		_, err := s.Client.Exec("UPDATE links SET tags = $1 WHERE username = $2", mergedTagsJSON, username)

		if err != nil {
			fmt.Println("Error when updating links' tags", err)
			return nil, errors.New("something went wrong")
		}

		return mergedTags, nil
	}

	_, err = s.Client.Exec("INSERT INTO links (username, link, tags) VALUES ($1, $2, $3)", username, link, tagsJSON)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key"){
			return nil, errors.New("you can't save the same url twice")
		}
		return nil, errors.New("something went wrong, try again later")
	}

	return nil, nil

}

func (s *PostgresStorage) GetLinksByUrl(username string, link string) (string, []string, error){
	rows, err := s.Client.Query("SELECT links, tags FROM links WHERE username = $1 AND link = $2", username, link)

	if err != nil {
			fmt.Println("Error when trying to query db in getlinksbyurl", err)
			return "", nil, errors.New("something went wrong")
	}
	var repeatedLink string
	var newTags []string

	defer rows.Close()
	for rows.Next() {
		var item string 
		var jsontags []byte

		err := rows.Scan(&item, &jsontags)

		if err != nil {
			fmt.Println("Error when trying to scan links in getlinksbyurl", err)
			return "", nil, errors.New("something went wrong")
		}

		var tags []string

		err = json.Unmarshal(jsontags, &tags)

		if err != nil {
			fmt.Println("Error when unmarshaling jsontags")
			return "", nil, errors.New("something went wrong")
		}

				// Split the string by commas
		parts := strings.Split(item, ",")

		// Extract the second element (index 1) after trimming any leading/trailing whitespace
		repeatedLink = strings.TrimSpace(parts[1])
		newTags = tags

	}

	return repeatedLink, newTags, nil


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
