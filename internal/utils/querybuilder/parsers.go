package querybuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type Results struct {
	Builder *Builder
	Page    int
	PerPage int
}

type FieldConfig struct {
	AllowedOperators []string
}

type ParseConfig struct {
	AllowedFields     map[string]FieldConfig
	AllowedSortFields []string
}

func ParseRequest(filter, pageStr, perPageStr, orderBy string, config ParseConfig) (*Results, error) {
	pageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		pageNum = 1
	}

	perPageNum, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPageNum = 10
	}

	builder := New()
	builder.Offset((pageNum - 1) * perPageNum)
	builder.Limit(perPageNum + 1)

	if orderBy != "" {
		splittedByComma := strings.Split(orderBy, ",")

		for _, field := range splittedByComma {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}

			splittedByColon := strings.Split(field, ":")
			fieldName := strings.TrimSpace(splittedByColon[0])

			if fieldName == "" {
				return nil, fmt.Errorf("invalid order_by: empty field name")
			}

			// Validate sort field
			if len(config.AllowedSortFields) == 0 {
				return nil, fmt.Errorf("no fields allowed for sorting")
			}
			found := false
			for _, f := range config.AllowedSortFields {
				if f == fieldName {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("field '%s' not allowed for sorting", fieldName)
			}

			direction := "asc"
			if len(splittedByColon) > 1 {
				direction = strings.ToLower(strings.TrimSpace(splittedByColon[1]))
				if direction != "asc" && direction != "desc" {
					return nil, fmt.Errorf("invalid order: must be 'asc' or 'desc'")
				}
			}

			builder.OrderBy(fieldName, direction)
		}
	}

	if filter != "" {
		b, err := parseFilter(filter, builder, config)
		if err != nil {
			return nil, err
		}
		builder = b
	}

	return &Results{
		Builder: builder,
		Page:    pageNum,
		PerPage: perPageNum,
	}, nil
}

func parseFilter(filter string, b *Builder, config ParseConfig) (*Builder, error) {
	splitted := splitByDelimiterOutsideQuotesAndParens(filter, " and ")
	if splitted == nil {
		return nil, fmt.Errorf("malformed filter: unclosed parenthesis in '%s'", filter)
	}

	allowedOperators := map[string]bool{
		"eq":   true,
		"ne":   true,
		"gt":   true,
		"gte":  true,
		"lt":   true,
		"lte":  true,
		"like": true,
		"in":   true,
	}

	for _, filterPart := range splitted {
		filterPart = strings.TrimSpace(filterPart)
		if filterPart == "" {
			return nil, fmt.Errorf("malformed filter: empty condition or redundant 'and' in '%s'", filter)
		}
		isOrGroup := strings.HasPrefix(filterPart, "(") && strings.HasSuffix(filterPart, ")")

		if isOrGroup {
			trimmedFilter := filterPart[1 : len(filterPart)-1]
			if strings.TrimSpace(trimmedFilter) == "" {
				return nil, fmt.Errorf("malformed filter: empty parenthesis '()' in '%s'", filter)
			}

			splittedOrGroup := splitByDelimiterOutsideQuotesAndParens(trimmedFilter, " or ")
			if splittedOrGroup == nil {
				return nil, fmt.Errorf("malformed filter: unclosed parenthesis in or group '%s'", trimmedFilter)
			}

			orQuery := b.InitOr()
			for _, filterSubPart := range splittedOrGroup {
				splittedFilter, err := splitFilter(filterSubPart)
				if err != nil {
					return nil, err
				}

				if len(splittedFilter) != 3 {
					return nil, fmt.Errorf("malformed filter: expected 'field operator value' but got '%s' in or group", filterSubPart)
				}

				field := splittedFilter[0]
				operator := splittedFilter[1]
				valueRaw := splittedFilter[2]

				// Validate field and operator
				if len(config.AllowedFields) == 0 {
					return nil, fmt.Errorf("no fields allowed for filtering")
				}
				fieldConfig, allowed := config.AllowedFields[field]
				if !allowed {
					return nil, fmt.Errorf("field '%s' not allowed for filtering", field)
				}
				if len(fieldConfig.AllowedOperators) > 0 {
					found := false
					for _, op := range fieldConfig.AllowedOperators {
						if op == operator {
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("operator '%s' not allowed for field '%s'", operator, field)
					}
				}

				if !allowedOperators[operator] {
					return nil, fmt.Errorf("operator '%s' not allowed at '%s'", operator, filterSubPart)
				}

				condValue, err := parseFilterValue(valueRaw)
				if err != nil {
					return nil, err
				}

				orQuery.Or(field, operator, condValue)
			}

			orQuery.EndOr()
		} else {
			splittedFilter, err := splitFilter(filterPart)
			if err != nil {
				return nil, err
			}

			if len(splittedFilter) != 3 {
				return nil, fmt.Errorf("malformed filter: expected 'field operator value' but got '%s'", filterPart)
			}

			field := splittedFilter[0]
			operator := splittedFilter[1]
			valueRaw := splittedFilter[2]

			// Validate field and operator
			if len(config.AllowedFields) == 0 {
				return nil, fmt.Errorf("no fields allowed for filtering")
			}
			fieldConfig, allowed := config.AllowedFields[field]
			if !allowed {
				return nil, fmt.Errorf("field '%s' not allowed for filtering", field)
			}
			if len(fieldConfig.AllowedOperators) > 0 {
				found := false
				for _, op := range fieldConfig.AllowedOperators {
					if op == operator {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("operator '%s' not allowed for field '%s'", operator, field)
				}
			}

			if !allowedOperators[operator] {
				return nil, fmt.Errorf("operator '%s' not allowed at '%s'", operator, filterPart)
			}

			condValue, err := parseFilterValue(valueRaw)
			if err != nil {
				return nil, err
			}

			b.And(field, operator, condValue)
		}
	}

	return b, nil
}

func parseFilterValue(value string) (any, error) {
	if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
		trimmed := value[1 : len(value)-1]
		if strings.TrimSpace(trimmed) == "" {
			return nil, fmt.Errorf("malformed list: empty list '()'")
		}

		elements := splitByDelimiterOutsideQuotesAndParens(trimmed, ",")
		if elements == nil {
			return nil, fmt.Errorf("malformed list: unclosed parenthesis in '%s'", value)
		}
		var slice []any
		for _, p := range elements {
			p = strings.TrimSpace(p)
			if p == "" {
				return nil, fmt.Errorf("malformed list: empty item")
			}
			v, err := parseFilterValue(p)
			if err != nil {
				return nil, err
			}
			slice = append(slice, v)
		}
		return slice, nil
	}

	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		trimmed := strings.TrimPrefix(value, "'")
		trimmed = strings.TrimSuffix(trimmed, "'")
		return strings.ReplaceAll(trimmed, "''", "'"), nil
	}

	if value == "true" || value == "false" {
		return value == "true", nil
	}

	if value == "null" {
		return nil, nil
	}

	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	return nil, fmt.Errorf("invalid value %s", value)
}

func splitFilter(filter string) ([]string, error) {
	firstSpace := strings.Index(filter, " ")
	if firstSpace == -1 {
		return nil, fmt.Errorf("filter syntax error (missing field separator) at '%s'", filter)
	}

	field := filter[:firstSpace]
	remaining := strings.TrimSpace(filter[firstSpace:])

	secondSpace := strings.Index(remaining, " ")
	if secondSpace == -1 {
		return nil, fmt.Errorf("filter syntax error (missing operator separator) at '%s'", filter)
	}

	operator := remaining[:secondSpace]
	value := strings.TrimSpace(remaining[secondSpace:])

	return []string{field, operator, value}, nil
}

func splitByDelimiterOutsideQuotesAndParens(input, delimiter string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	parenLevel := 0

	for i := 0; i < len(input); i++ {
		char := input[i]
		if char == '\'' {
			inQuotes = !inQuotes
		} else if char == '(' && !inQuotes {
			parenLevel++
		} else if char == ')' && !inQuotes {
			parenLevel--
		}

		if !inQuotes && parenLevel == 0 && strings.HasPrefix(input[i:], delimiter) {
			result = append(result, strings.TrimSpace(current.String()))
			current.Reset()
			i += len(delimiter) - 1
			continue
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	if parenLevel != 0 {
		return nil // Indicate unclosed parens by returning nil so caller knows it failed
	}
	return result
}
