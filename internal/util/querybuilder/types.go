package querybuilder

type Metadata struct {
	NextPage   bool `json:"next_page" binding:"required"`
	Page       int  `json:"page" binding:"required"`
	PerPage    int  `json:"per_page" binding:"required"`
	TotalItems int  `json:"total_items" binding:"required"`
	TotalPages int  `json:"total_pages" binding:"required"`
}
