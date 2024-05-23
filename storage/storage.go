package storage

type Storage interface {
	GetLinks(string, []string) []string
	GetLinksByUrl(string, string) (string, []string, error)
	InsertLinkAndTags(string, string, []string) ([]string, error)
	GetUserTags(string) ([]string, error)
}



