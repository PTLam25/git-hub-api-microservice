package github

// Представления модели для запроса и
//ответа при создания репозиторий через GitHub API

type CreateRepoRequest struct { // Представления Request для создания репозиторий
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Private     bool   `json:"private"`
	HasIssues   bool   `json:"has_issues"`
	HasProjects bool   `json:"has_projects"`
	HasWiki     bool   `json:"has_wiki"`
}

type CreateRepoResponse struct { // Представления Response от GitHub при создания репозиторий
	Id         int64          `json:"id"`
	Name       string         `json:"name"`
	FullName   string         `json:"full_name"`
	Owner      RepoOwner      `json:"owner"`
	Permission RepoPermission `json:"permission"`
}

type RepoOwner struct {
	Id      int64  `json:"id"`
	Login   string `json:"login"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

type RepoPermission struct {
	IsAdmin bool `json:"is_admin"`
	HasPull bool `json:"has_pull"`
	HasPush bool `json:"has_push"`
}
