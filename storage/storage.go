package storage

type Storage interface {
	//Takes in: username, array of tags. Returns: array of strings links and error
	GetLinks(string, []string) ([]string, error)

	//Takes in: string username, string link. Returns: string repeatedLinks, array of strings tags, error
	GetLinksByUrl(string, string) (string, []string, error)

	//Takes in: string username, string link, array of strings tags. Returns: mergedTags (if any), error
	InsertLinkAndTags(string, string, []string) ([]string, error)

	//Takes in: string username, array of strings previousTags, array of strings newTags. Returns: mergedTags and error
	UpdateLinkTags(string, string, []string, []string) ([]string, error)

	//Takes in: string username. Returns: array of strings tags and error
	GetUserTags(string) ([]string, error)
}



