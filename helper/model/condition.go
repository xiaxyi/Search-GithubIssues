package model

type Conditions struct {
	RepoOwner        string
	RepoName         string
	CreatedTimeStart string
	CreatedTimeEnd   string
	KeyWordsInTitle  []string
	ResourceProvider string
}
