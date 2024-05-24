package storage

import (
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

func (s *PostgresStorage) GetLinks(username string, tags []string) ([]string, error) {

	
	rows, err := s.Client.Query("SELECT DISTINCT link, tags FROM links WHERE username = $1 AND tags ?| $2::text[];", username, pq.Array(tags))
	
	if err != nil {
		fmt.Println("Something went wrong when fetching links")
		return nil, errors.New("something went wrong")
	}
	
	defer rows.Close()
	
	var links []string

	for rows.Next() {
		var link string
		var tagsJson []byte
		err := rows.Scan(&link, &tagsJson); 

		if err != nil {
			fmt.Println("Error when scanning link or tagsJson in getlinks")
			return nil, errors.New("something went wrong")
		}

		var tags []string
		err = json.Unmarshal(tagsJson, &tags)

		if err != nil {
			panic(err)
		}
		
		links = append(links, link)
	}

	err = rows.Err()

	if err != nil {		
    	panic(err)
  	}
	
	fmt.Println("\nlinks to send: ", links)
	return links, nil

}


func (s *PostgresStorage) InsertLinkAndTags(username string, link string, tags []string) ([]string, error) {

	potentialDuplicateLink, previousTags, error := s.GetLinksByUrl(username, link)

	if error != nil {
		fmt.Println("Error when calling getlinksbyurl")
		return nil, errors.New("something went wrong")
	}

	fmt.Println("Repeated link: ", potentialDuplicateLink, "previous tags: ", previousTags)

	if potentialDuplicateLink == link{
		mergedTags, error := s.UpdateLinkTags(username, previousTags, tags)

		if error != nil {
			fmt.Println("Something went wrong when updating link tags, error: ", error)
			return nil, errors.New("something went wrong")
		}
		return mergedTags, nil 
	}

	tagsJSON, err := json.Marshal(tags)

	if err != nil {
		fmt.Println("error when marshaling tagsJSON")
		return nil, errors.New("something went wrong, try again later")
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

func (s *PostgresStorage) UpdateLinkTags(username string, previousTags []string, newTags []string) ([]string, error) {

	var mergedTags []string = util.MergeSlices(newTags, previousTags)
	mergedTagsJSON, err := json.Marshal(mergedTags)

	_, queryError:= s.Client.Exec("UPDATE links SET tags = $1 WHERE username = $2", mergedTagsJSON, username)

		if queryError != nil {
			fmt.Println("Error when updating links' tags", err)
			return nil, errors.New("something went wrong")
		}

		return mergedTags, nil
}

func (s *PostgresStorage) GetLinksByUrl(username string, link string) (string, []string, error){
	
	rows, err := s.Client.Query("SELECT link, tags FROM links WHERE username = $1 AND link = $2", username, link)

	if err != nil {
			fmt.Println("Error when trying to query db in getlinksbyurl", err)
			return "", nil, errors.New("something went wrong")
	}

	var repeatedLink string
	var tags []string

	//traverse through the byte like values returned, scan them and save them.
	defer rows.Close()
	for rows.Next() {
		var tagsJSON []byte

		err := rows.Scan(&repeatedLink, &tagsJSON)

		if err != nil {
			fmt.Println("Error when trying to scan links in getlinksbyurl", err)
			return "", nil, errors.New("something went wrong")
		}

		//Unmarshal json tags into array of string tags
		err = json.Unmarshal(tagsJSON, &tags)

		if err != nil {
			fmt.Println("Error when unmarshaling jsontags")
			return "", nil, errors.New("something went wrong")
		}
	}

	return repeatedLink, tags, nil

}

func (s *PostgresStorage) GetUserTags(username string) ([]string, error) {
	rows, err := s.Client.Query("SELECT DISTINCT jsonb_array_elements_text(tags) AS unique_tag FROM links WHERE username = $1;", username)

	if err != nil {
		fmt.Println("Error when retrieving user tags, error: ", err)
		return []string{}, errors.New("something went wrong, try again later")
	}

	defer rows.Close()

	var tags []string

	for rows.Next(){
		var tag string

		err := rows.Scan(&tag)

		if err != nil {
			fmt.Println("Error ocurred when parsing tags rows")
			return []string{}, errors.New("something went wrong, try again later")
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
