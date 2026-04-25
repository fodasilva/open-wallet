package querybuilder

type Metadata struct {
	NextPage   bool `json:"next_page"`
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
}
