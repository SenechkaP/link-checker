package link

type LinksRequest struct {
	Links []string `json:"links"`
}

type LinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

type LinksListRequest struct {
	LinksList []int `json:"links_list"`
}
