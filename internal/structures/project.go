package structures

type Project struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	RepositoryUrl string `json:"repository_url"`
	BannerUrl     string `json:"banner_url"`
}
