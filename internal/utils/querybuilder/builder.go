package querybuilder

type Condition struct {
	Field    string
	Operator string
	Value    any
}

type Order struct {
	Field string
	Dir   string
}

type Builder struct {
	AndConditions []Condition
	OrGroups      [][]Condition
	Orders        []Order
	LimitValue    *int
	OffsetValue   *int
}

func New() *Builder {
	return &Builder{
		AndConditions: make([]Condition, 0),
		OrGroups:      make([][]Condition, 0),
		Orders:        make([]Order, 0),
		LimitValue:    nil,
		OffsetValue:   nil,
	}
}

func (b *Builder) And(field, operator string, value any) *Builder {
	b.AndConditions = append(b.AndConditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return b
}

type OrBuilder struct {
	conditions []Condition
	builder    *Builder
}

func (ob *OrBuilder) Or(field string, operator string, value any) *OrBuilder {
	ob.conditions = append(ob.conditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return ob
}

func (ob *OrBuilder) EndOr() *Builder {
	ob.builder.OrGroups = append(ob.builder.OrGroups, ob.conditions)
	return ob.builder
}

func (b *Builder) InitOr() *OrBuilder {
	return &OrBuilder{
		builder: b,
	}
}

func (b *Builder) OrderBy(field, direction string) *Builder {
	b.Orders = append(b.Orders, Order{
		Field: field,
		Dir:   direction,
	})
	return b
}

func (b *Builder) Limit(limit int) *Builder {
	b.LimitValue = &limit
	return b
}

func (b *Builder) Offset(offset int) *Builder {
	b.OffsetValue = &offset
	return b
}

func ForCount(b *Builder) *Builder {
	return &Builder{
		AndConditions: b.AndConditions,
		OrGroups:      b.OrGroups,
		Orders:        nil,
		LimitValue:    nil,
		OffsetValue:   nil,
	}
}

func (b *Builder) HasAndFieldOperator(field, operator string) []Condition {
	var results []Condition
	for _, cond := range b.AndConditions {
		if cond.Field == field && cond.Operator == operator {
			results = append(results, cond)
		}
	}
	return results
}

func (b *Builder) HasAndField(field string) []Condition {
	var results []Condition
	for _, cond := range b.AndConditions {
		if cond.Field == field {
			results = append(results, cond)
		}
	}
	return results
}
