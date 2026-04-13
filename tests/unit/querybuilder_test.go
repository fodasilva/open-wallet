package unit

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func TestQueryBuilder_Fluent(t *testing.T) {
	t.Run("should build fluently", func(t *testing.T) {
		b := querybuilder.New().
			And("name", "eq", "John").
			OrderBy("created_at", "desc").
			Limit(10).
			Offset(0)

		assert.Len(t, b.AndConditions, 1)
		assert.Equal(t, "John", b.AndConditions[0].Value)
		assert.Equal(t, "name", b.AndConditions[0].Field)
		assert.Equal(t, "desc", b.Orders[0].Dir)
		assert.Equal(t, 10, *b.LimitValue)
	})

	t.Run("should build with OR groups", func(t *testing.T) {
		b := querybuilder.New().
			InitOr().
			Or("status", "eq", "pending").
			Or("status", "eq", "failed").
			EndOr().
			And("active", "eq", true)

		assert.Len(t, b.OrGroups, 1)
		assert.Len(t, b.OrGroups[0], 2)
		assert.Equal(t, "pending", b.OrGroups[0][0].Value)
		assert.Equal(t, "active", b.AndConditions[0].Field)
	})
}

func TestQueryBuilder_ParseRequest(t *testing.T) {
	tests := []struct {
		name        string
		filter      string
		page        string
		perPage     string
		orderBy     string
		wantPage    int
		wantPerPage int
		wantErr     bool
		validate    func(*testing.T, *querybuilder.Results)
	}{
		{
			name:        "full request",
			filter:      "name eq 'John Doe' and active eq true",
			orderBy:     "created_at:desc",
			page:        "2",
			perPage:     "15",
			wantPage:    2,
			wantPerPage: 15,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 2)
				assert.Equal(t, "John Doe", r.Builder.AndConditions[0].Value)
				assert.Equal(t, true, r.Builder.AndConditions[1].Value)
				assert.Equal(t, "created_at", r.Builder.Orders[0].Field)
				assert.Equal(t, "desc", r.Builder.Orders[0].Dir)
			},
		},
		{
			name:        "empty parameters with defaults",
			filter:      "",
			page:        "",
			perPage:     "",
			orderBy:     "",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Empty(t, r.Builder.AndConditions)
				assert.Empty(t, r.Builder.Orders)
			},
		},
		{
			name:        "different value types",
			filter:      "age gt 25 and verified eq false and description eq null",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 3)
				assert.Equal(t, float64(25), r.Builder.AndConditions[0].Value)
				assert.Equal(t, false, r.Builder.AndConditions[1].Value)
				assert.Nil(t, r.Builder.AndConditions[2].Value)
			},
		},
		{
			name:        "keywords inside quotes",
			filter:      "description eq 'This and that' and name eq 'Or or something'",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 2)
				assert.Equal(t, "This and that", r.Builder.AndConditions[0].Value)
				assert.Equal(t, "Or or something", r.Builder.AndConditions[1].Value)
			},
		},
		{
			name:        "OR groups",
			filter:      "(status eq 'pending' or status eq 'failed') and active eq true",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.OrGroups, 1)
				assert.Len(t, r.Builder.OrGroups[0], 2)
				assert.Equal(t, "pending", r.Builder.OrGroups[0][0].Value)
				assert.Equal(t, "failed", r.Builder.OrGroups[0][1].Value)
				assert.Equal(t, true, r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "all operators parsing",
			filter:      "a ne 1 and b gt 2 and c gte 3 and d lt 4 and e lte 5 and f like 'foo'",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 6)
				assert.Equal(t, "ne", r.Builder.AndConditions[0].Operator)
				assert.Equal(t, "gt", r.Builder.AndConditions[1].Operator)
				assert.Equal(t, "gte", r.Builder.AndConditions[2].Operator)
				assert.Equal(t, "lt", r.Builder.AndConditions[3].Operator)
				assert.Equal(t, "lte", r.Builder.AndConditions[4].Operator)
				assert.Equal(t, "like", r.Builder.AndConditions[5].Operator)
			},
		},
		{
			name:        "in operator with numbers",
			filter:      "age in (25, 30, 35)",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "in", r.Builder.AndConditions[0].Operator)
				assert.Equal(t, []any{float64(25), float64(30), float64(35)}, r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "in operator with strings",
			filter:      "name in ('Alice', 'Bob')",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "in", r.Builder.AndConditions[0].Operator)
				assert.Equal(t, []any{"Alice", "Bob"}, r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "in operator with booleans",
			filter:      "is_active in (true, false)",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "in", r.Builder.AndConditions[0].Operator)
				assert.Equal(t, []any{true, false}, r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "in operator with null",
			filter:      "name in ('Alice', null)",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "in", r.Builder.AndConditions[0].Operator)
				assert.Equal(t, []any{"Alice", nil}, r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "multiple OR groups",
			filter:      "(status eq 'pending' or status eq 'failed') and (type eq 1 or type eq 2)",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.OrGroups, 2)
				assert.Len(t, r.Builder.OrGroups[0], 2)
				assert.Len(t, r.Builder.OrGroups[1], 2)
			},
		},
		{
			name:        "extra spaces between parts",
			filter:      "name   eq   'John'",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "name", r.Builder.AndConditions[0].Field)
				assert.Equal(t, "eq", r.Builder.AndConditions[0].Operator)
				assert.Equal(t, "John", r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "quoted boolean should be string",
			filter:      "active eq 'true'",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.IsType(t, "", r.Builder.AndConditions[0].Value)
				assert.Equal(t, "true", r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "quoted value with leading/trailing spaces",
			filter:      "name eq '  John Doe  '",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "  John Doe  ", r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "empty string value",
			filter:      "name eq ''",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "", r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "string with only space",
			filter:      "name eq ' '",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, " ", r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "escaped single quotes (double single quotes)",
			filter:      "name eq 'O''Brien'",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.AndConditions, 1)
				assert.Equal(t, "O'Brien", r.Builder.AndConditions[0].Value)
			},
		},
		{
			name:        "null inside OR group",
			filter:      "(category eq null or category eq 'Misc')",
			page:        "1",
			perPage:     "10",
			wantPage:    1,
			wantPerPage: 10,
			validate: func(t *testing.T, r *querybuilder.Results) {
				assert.Len(t, r.Builder.OrGroups, 1)
				assert.Nil(t, r.Builder.OrGroups[0][0].Value)
				assert.Equal(t, "Misc", r.Builder.OrGroups[0][1].Value)
			},
		},
		{
			name:    "invalid syntax - missing parts",
			filter:  "name eq",
			wantErr: true,
		},
		{
			name:    "invalid operator",
			filter:  "name invalidop 'value'",
			wantErr: true,
		},
		{
			name:    "unclosed parenthesis",
			filter:  "(age gt 10",
			wantErr: true,
		},
		{
			name:    "empty parenthesis",
			filter:  "()",
			wantErr: true,
		},
		{
			name:    "duplicated logic",
			filter:  "age gt 10 and and name eq 'John'",
			wantErr: true,
		},
		{
			name:    "malformed in list - empty item",
			filter:  "age in (1,,3)",
			wantErr: true,
		},
		{
			name:    "invalid sort direction",
			orderBy: "name:sideways",
			wantErr: true,
		},
		{
			name:    "empty field in sort",
			orderBy: ":asc",
			wantErr: true,
		},
	}

	permissiveConfig := &querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"id":             {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"name":           {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"age":            {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"is_active":      {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"active":         {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"type":           {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"created_at":     {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"color":          {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"amount":         {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"reference_date": {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"verified":       {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"description":    {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"status":         {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"category":       {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"a":              {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"b":              {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"c":              {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"d":              {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"e":              {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
			"f":              {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "like", "in"}},
		},
		AllowedSortFields: []string{"id", "name", "age", "is_active", "active", "created_at", "color", "amount", "reference_date", "verified", "description", "status", "category", "a", "b", "c", "d", "e", "f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := querybuilder.ParseRequest(tt.filter, tt.page, tt.perPage, tt.orderBy, *permissiveConfig)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantPage, results.Page)
			assert.Equal(t, tt.wantPerPage, results.PerPage)
			if tt.validate != nil {
				tt.validate(t, results)
			}
		})
	}
}

func TestQueryBuilder_ParseRequest_WithConfig(t *testing.T) {
	config := &querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"name": {
				AllowedOperators: []string{"eq", "like"},
			},
			"status": {
				AllowedOperators: []string{"eq"},
			},
		},
		AllowedSortFields: []string{"name", "created_at"},
	}

	tests := []struct {
		name    string
		filter  string
		orderBy string
		wantErr bool
	}{
		{
			name:   "permitted field and operator",
			filter: "name eq 'John'",
		},
		{
			name:    "disallowed field",
			filter:  "age eq 25",
			wantErr: true,
		},
		{
			name:    "disallowed operator for a permitted field",
			filter:  "status like 'pen%'",
			wantErr: true,
		},
		{
			name:    "permitted sort field",
			orderBy: "name:desc",
		},
		{
			name:    "multiple sorts - one disallowed",
			orderBy: "name:asc, id:desc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := querybuilder.ParseRequest(tt.filter, "1", "10", tt.orderBy, *config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQueryBuilder_ToSquirrel(t *testing.T) {
	tests := []struct {
		name     string
		build    func() *querybuilder.Builder
		wantSQL  string
		wantArgs []interface{}
	}{
		{
			name: "basic select with where, order and limit",
			build: func() *querybuilder.Builder {
				return querybuilder.New().
					And("name", "eq", "John").
					OrderBy("id", "asc").
					Limit(10)
			},
			wantSQL:  "SELECT * FROM users WHERE name = ? ORDER BY id ASC LIMIT 10",
			wantArgs: []interface{}{"John"},
		},
		{
			name: "all operators SQL",
			build: func() *querybuilder.Builder {
				return querybuilder.New().
					And("a", "ne", 1).
					And("b", "gt", 2).
					And("c", "gte", 3).
					And("d", "lt", 4).
					And("e", "lte", 5)
			},
			wantSQL:  "SELECT * FROM users WHERE a <> ? AND b > ? AND c >= ? AND d < ? AND e <= ?",
			wantArgs: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name: "like operator with upper case",
			build: func() *querybuilder.Builder {
				return querybuilder.New().And("title", "like", "rock")
			},
			wantSQL:  "SELECT * FROM music WHERE upper(title) LIKE ?",
			wantArgs: []interface{}{"%ROCK%"},
		},
		{
			name: "OR conditions",
			build: func() *querybuilder.Builder {
				return querybuilder.New().
					InitOr().
					Or("a", "eq", 1).
					Or("b", "eq", 2).
					EndOr()
			},
			wantSQL:  "SELECT * FROM t WHERE (a = ? OR b = ?)",
			wantArgs: []interface{}{1, 2},
		},
		{
			name: "IN operator SQL",
			build: func() *querybuilder.Builder {
				return querybuilder.New().And("id", "in", []any{1, 2, 3})
			},
			wantSQL:  "SELECT * FROM users WHERE id IN (?,?,?)",
			wantArgs: []interface{}{1, 2, 3},
		},
		{
			name: "IN operator with NULL SQL",
			build: func() *querybuilder.Builder {
				return querybuilder.New().And("id", "in", []any{1, nil})
			},
			wantSQL:  "SELECT * FROM users WHERE (id IN (?) OR id IS NULL)",
			wantArgs: []interface{}{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.build()
			// Extract table name from wantSQL for simulation
			table := "users"
			if tt.name == "like operator with upper case" {
				table = "music"
			} else if tt.name == "OR conditions" {
				table = "t"
			}

			query := squirrel.Select("*").From(table)
			query = querybuilder.ToSquirrel(query, b)

			sql, args, err := query.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSQL, sql)
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}
