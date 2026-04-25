package querybuilder

func BuildMetadata(page, perPage, totalItems int) Metadata {
	totalPages := 0
	if perPage > 0 {
		totalPages = (totalItems + perPage - 1) / perPage
	}

	hasNext := (page * perPage) < totalItems

	return Metadata{
		Page:       page,
		PerPage:    perPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
		NextPage:   hasNext,
	}
}
