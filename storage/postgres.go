package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/yohanc3/link-vault/error"
	. "github.com/yohanc3/link-vault/logger"
	"github.com/yohanc3/link-vault/util"

	"github.com/lib/pq"
)

type PostgresStorage struct {
	Client *sql.DB
}

func NewPostgresDb(psqlInfo string) *PostgresStorage {

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		StorageLogger.Panic().Msg("Error when establishing connection with db." + err.Error())	
	}

	fmt.Println("Successfully started a supabase client!!")

	return &PostgresStorage{
		Client: db,
	}
}

func (s *PostgresStorage) GetLinks(username string, tags []string) ([]string, error) {

	
	rows, err := s.Client.Query("SELECT DISTINCT link, tags FROM links WHERE username = $1 AND tags ?| $2::text[];", username, pq.Array(tags))
	
	if err != nil {
		StorageLogger.Error().Msg("error when fetching links from db. " + err.Error())
		return nil, GenericError
	}
	
	defer rows.Close()
	
	var links []string

	for rows.Next() {
		var link string
		var tagsJson []byte
		err := rows.Scan(&link, &tagsJson); 

		if err != nil {
			StorageLogger.Error().Msg("error when scanning link or tagsJson. " + err.Error())
			return nil, GenericError
		}

		var tags []string
		err = json.Unmarshal(tagsJson, &tags)

		if err != nil {
			StorageLogger.Error().Str("json to unmarshal", "["+strings.Join(tags, ",")+"]").Msg("error when unmarshaling json")
			return nil, GenericError
		}
		
		links = append(links, link)
	}

	err = rows.Err()

	if err != nil {		
			StorageLogger.Error().Msg("error when iterating through rows." + err.Error())
			return nil, GenericError
  	}
	
	fmt.Println("\nlinks to send: ", links)
	return links, nil

}


func (s *PostgresStorage) InsertLinkAndTags(username string, link string, tags []string) ([]string, error) {

	potentialDuplicateLink, previousTags, err := s.GetLinksByUrl(username, link)

	if err != nil {
		return nil, GenericError 
	}

	GeneralLogger.Trace().Str("given link", link).Str("repeated link", potentialDuplicateLink).Str("previous tags", "["+strings.Join(previousTags, ",")+"]").Str("username", username).Msg("")

	if potentialDuplicateLink == link{
		mergedTags, err := s.UpdateLinkTags(username, link, previousTags, tags)

		if err != nil {
			StorageLogger.Error().Str("username", username).Msg("error when updating link tags. " + err.Error())
			return nil, GenericError 
		}
		return mergedTags, nil 
	}

	tagsJSON, err := json.Marshal(tags)

	if err != nil {
		StorageLogger.Error().Msg("error when marshaling tagsJSON. " + err.Error())
		return nil, GenericError
	}

	_, err = s.Client.Exec("INSERT INTO links (username, link, tags) VALUES ($1, $2, $3)", username, link, tagsJSON)

	if err != nil {
		StorageLogger.Error().Msg("error when inserting links. " + err.Error())
		return nil, GenericError
	}

	return nil, nil

}

func (s *PostgresStorage) UpdateLinkTags(username string, link string, previousTags []string, newTags []string) ([]string, error) {

	var mergedTags []string = util.MergeSlices(newTags, previousTags)
	mergedTagsJSON, err := json.Marshal(mergedTags)

	if err != nil {
		StorageLogger.Error().Msg("error when marshaling tags" + err.Error())
		return nil, GenericError 
	}

	_, err = s.Client.Exec("UPDATE links SET tags = $1 WHERE username = $2 AND link = $3", mergedTagsJSON, username, link)

		if err != nil {
			StorageLogger.Error().Msg("error when updating links' tags" + err.Error())
			return nil, GenericError 
		}

		return mergedTags, nil
}

func (s *PostgresStorage) GetLinksByUrl(username string, link string) (string, []string, error){
	
	rows, err := s.Client.Query("SELECT link, tags FROM links WHERE username = $1 AND link = $2", username, link)

	if err != nil {
		return "", nil, GenericError 
	}

	var repeatedLink string
	var tags []string

	//traverse through the byte like values returned, scan them and save them.
	defer rows.Close()
	for rows.Next() {
		var tagsJSON []byte

		err := rows.Scan(&repeatedLink, &tagsJSON)

		if err != nil {
			StorageLogger.Error().Msg("error when scanning links. " + err.Error())
			return "", nil, GenericError 
		}

		//Unmarshal json tags into array of string tags
		err = json.Unmarshal(tagsJSON, &tags)

		if err != nil {
			StorageLogger.Error().Msg("Error when unmarshaling jsontags. " + err.Error())
			return "", nil, GenericError 
		}
	}

	fmt.Println("repeated link: ", repeatedLink)

	return repeatedLink, tags, nil

}

func (s *PostgresStorage) GetUserTags(username string) ([]string, error) {
	rows, err := s.Client.Query("SELECT DISTINCT jsonb_array_elements_text(tags) AS unique_tag FROM links WHERE username = $1;", username)

	if err != nil {
		StorageLogger.Error().Msg("error when retrieving user tags. " + err.Error())
		return []string{}, GenericError
	}

	defer rows.Close()

	var tags []string

	for rows.Next(){
		var tag string

		err := rows.Scan(&tag)

		if err != nil {
			StorageLogger.Error().Msg("error when scanning rows. " + err.Error())
			return []string{}, GenericError
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
