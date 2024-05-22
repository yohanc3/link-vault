package storage

type Storage interface {
	GetLinks(string, []string) []string
	InsertLinkAndTags(string, string, []string) error 
	GetUserTags(string) ([]string, error)
}



