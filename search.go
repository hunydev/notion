package notion

type FilterOperation string

const (
	FilterOperationOR  = FilterOperation("or")
	FilterOperationAND = FilterOperation("and")
)

func (operation FilterOperation) String() string {
	return string(operation)
}

type Filter interface {
	Json() JSON
}

type DatabasePropertyFilter struct {
	JSON JSON

	condition Condition
}

func NewFilter(Property string, Condition Condition) Filter {
	filter := &DatabasePropertyFilter{
		JSON: JSON{
			"property": Property,
			Condition.Type(): JSON{
				Condition.Key(): Condition.Value(),
			},
		},
		condition: Condition,
	}

	return filter
}

func (filter *DatabasePropertyFilter) Json() JSON {
	return filter.JSON
}

func (filter *DatabasePropertyFilter) Property() string {
	return filter.JSON.GetString("property")
}

func (filter *DatabasePropertyFilter) Condition() Condition {
	return filter.condition
}

type CompoundFilter struct {
	JSON JSON

	operation FilterOperation
}

func NewCompoundFilter(Operation FilterOperation, Filters ...Filter) Filter {
	switch Operation {
	case FilterOperationOR, FilterOperationAND:
	default:
		return nil
	}

	filter := &CompoundFilter{
		JSON: JSON{
			Operation.String(): []JSON{},
		},
		operation: Operation,
	}

	for _, f := range Filters {
		var v interface{} = f
		switch v.(type) {
		case *DatabasePropertyFilter, DatabasePropertyFilter:
			filter.JSON.Append(Operation.String(), f.Json())
		case *CompoundFilter, CompoundFilter:
			if j, ok := f.Json().GetJSONList(string(FilterOperationOR)); ok {
				for _, jj := range j {
					filter.JSON.Append(Operation.String(), jj)
				}
			} else if j, ok := f.Json().GetJSONList(string(FilterOperationAND)); ok {
				for _, jj := range j {
					filter.JSON.Append(Operation.String(), jj)
				}
			}
		}
	}

	return filter
}

func (filter *CompoundFilter) Json() JSON {
	return filter.JSON
}

func (filter *CompoundFilter) Operation() FilterOperation {
	return filter.operation
}

type Condition interface {
	Type() string
	Key() string
	Value() interface{}
}

type Conditioner struct {
	t string
	k string
	v interface{}
}

func NewCondition(Type, Key string, Value interface{}) Condition {
	return &Conditioner{t: Type, k: Key, v: Value}
}

func (c *Conditioner) Type() string {
	return c.t
}

func (c *Conditioner) Key() string {
	return c.k
}

func (c *Conditioner) Value() interface{} {
	return c.v
}

type ConditionFuncString func(value string) Condition
type ConditionFuncNumber func(value int) Condition
type ConditionFuncBoolean func(value bool) Condition
type ConditionFuncObject func(value interface{}) Condition

type ConditionText struct {
	Equals         ConditionFuncString
	DoesNotEqual   ConditionFuncString
	Contains       ConditionFuncString
	DoesNotContain ConditionFuncString
	StartsWith     ConditionFuncString
	EndsWith       ConditionFuncString
	IsEmpty        ConditionFuncBoolean
	IsNotEmpty     ConditionFuncBoolean
}

var FilterText = &ConditionText{
	Equals:         func(value string) Condition { return NewCondition("text", "equals", value) },
	DoesNotEqual:   func(value string) Condition { return NewCondition("text", "does_not_equal", value) },
	Contains:       func(value string) Condition { return NewCondition("text", "contains", value) },
	DoesNotContain: func(value string) Condition { return NewCondition("text", "does_not_contain", value) },
	StartsWith:     func(value string) Condition { return NewCondition("text", "starts_with", value) },
	EndsWith:       func(value string) Condition { return NewCondition("text", "ends_with", value) },
	IsEmpty:        func(value bool) Condition { return NewCondition("text", "is_empty", true) },
	IsNotEmpty:     func(value bool) Condition { return NewCondition("text", "is_not_empty", true) },
}

type ConditionNumber struct {
	Equals               ConditionFuncNumber
	DoesNotEqual         ConditionFuncNumber
	GreaterThan          ConditionFuncNumber
	LessThan             ConditionFuncNumber
	GreaterThanOrEqualTo ConditionFuncNumber
	LessThanOrEqualTo    ConditionFuncNumber
	IsEmpty              ConditionFuncBoolean
	IsNotEmpty           ConditionFuncBoolean
}

var FilterNumber = &ConditionNumber{
	Equals:               func(value int) Condition { return NewCondition("number", "equals", value) },
	DoesNotEqual:         func(value int) Condition { return NewCondition("number", "does_not_equal", value) },
	GreaterThan:          func(value int) Condition { return NewCondition("number", "greater_than", value) },
	LessThan:             func(value int) Condition { return NewCondition("number", "less_than", value) },
	GreaterThanOrEqualTo: func(value int) Condition { return NewCondition("number", "greater_than_or_equal_to", value) },
	LessThanOrEqualTo:    func(value int) Condition { return NewCondition("number", "less_than_or_equal_to", value) },
	IsEmpty:              func(value bool) Condition { return NewCondition("number", "is_empty", true) },
	IsNotEmpty:           func(value bool) Condition { return NewCondition("number", "is_not_empty", true) },
}

type ConditionCheckbox struct {
	Equals       ConditionFuncBoolean
	DoesNotEqual ConditionFuncBoolean
}

var FilterCheckbox = &ConditionCheckbox{
	Equals:       func(value bool) Condition { return NewCondition("checkbox", "equals", value) },
	DoesNotEqual: func(value bool) Condition { return NewCondition("checkbox", "does_not_equal", value) },
}

type ConditionSelect struct {
	Equals       ConditionFuncString
	DoesNotEqual ConditionFuncString
	IsEmpty      ConditionFuncBoolean
	IsNotEmpty   ConditionFuncBoolean
}

var FilterSelect = &ConditionSelect{
	Equals:       func(value string) Condition { return NewCondition("select", "equals", value) },
	DoesNotEqual: func(value string) Condition { return NewCondition("select", "does_not_equal", value) },
	IsEmpty:      func(value bool) Condition { return NewCondition("select", "is_empty", true) },
	IsNotEmpty:   func(value bool) Condition { return NewCondition("select", "is_not_empty", true) },
}

type ConditionMultiSelect struct {
	Contains       ConditionFuncString
	DoesNotContain ConditionFuncString
	IsEmpty        ConditionFuncBoolean
	IsNotEmpty     ConditionFuncBoolean
}

var FilterMultiSelect = &ConditionMultiSelect{
	Contains:       func(value string) Condition { return NewCondition("multi_select", "contains", value) },
	DoesNotContain: func(value string) Condition { return NewCondition("multi_select", "doeS_not_contain", value) },
	IsEmpty:        func(value bool) Condition { return NewCondition("multi_select", "is_empty", true) },
	IsNotEmpty:     func(value bool) Condition { return NewCondition("multi_select", "is_not_empty", true) },
}

type ConditionDate struct {
	Equals     ConditionFuncString
	Before     ConditionFuncString
	After      ConditionFuncString
	OnOrBefore ConditionFuncString
	IsEmpty    ConditionFuncBoolean
	IsNotEmpty ConditionFuncBoolean
	OnOrAfter  ConditionFuncString
	PastWeek   ConditionFuncObject
	PastMonth  ConditionFuncObject
	PastYear   ConditionFuncObject
	NextWeek   ConditionFuncObject
	NextMonth  ConditionFuncObject
	NextYear   ConditionFuncObject
}

var FilterDate = &ConditionDate{
	Equals:     func(value string) Condition { return NewCondition("date", "equals", value) },
	Before:     func(value string) Condition { return NewCondition("date", "before", value) },
	After:      func(value string) Condition { return NewCondition("date", "after", value) },
	OnOrBefore: func(value string) Condition { return NewCondition("date", "on_or_before", value) },
	IsEmpty:    func(value bool) Condition { return NewCondition("date", "is_empty", true) },
	IsNotEmpty: func(value bool) Condition { return NewCondition("date", "is_not_empty", true) },
	OnOrAfter:  func(value string) Condition { return NewCondition("date", "on_or_after", value) },
	PastWeek:   func(value interface{}) Condition { return NewCondition("date", "past_week", JSON{}) },
	PastMonth:  func(value interface{}) Condition { return NewCondition("date", "past_month", JSON{}) },
	PastYear:   func(value interface{}) Condition { return NewCondition("date", "past_year", JSON{}) },
	NextWeek:   func(value interface{}) Condition { return NewCondition("date", "next_week", JSON{}) },
	NextMonth:  func(value interface{}) Condition { return NewCondition("date", "next_month", JSON{}) },
	NextYear:   func(value interface{}) Condition { return NewCondition("date", "next_year", JSON{}) },
}

type ConditionPeople struct {
	Contains       ConditionFuncString
	DoesNotContain ConditionFuncString
	IsEmpty        ConditionFuncBoolean
	IsNotEmpty     ConditionFuncBoolean
}

var FilterPeople = &ConditionPeople{
	Contains:       func(value string) Condition { return NewCondition("people", "contains", value) },
	DoesNotContain: func(value string) Condition { return NewCondition("people", "does_not_contain", value) },
	IsEmpty:        func(value bool) Condition { return NewCondition("people", "is_empty", true) },
	IsNotEmpty:     func(value bool) Condition { return NewCondition("people", "is_not_empty", true) },
}

type ConditionFiles struct {
	IsEmpty    ConditionFuncBoolean
	IsNotEmpty ConditionFuncBoolean
}

var FilterFiles = &ConditionFiles{
	IsEmpty:    func(value bool) Condition { return NewCondition("files", "is_empty", true) },
	IsNotEmpty: func(value bool) Condition { return NewCondition("files", "is_not_empty", true) },
}

type ConditionFormula struct {
	Text     *ConditionText
	Checkbox *ConditionCheckbox
	Number   *ConditionNumber
	Date     *ConditionDate
}

var FilterFormula = &ConditionFormula{
	Text:     FilterText,
	Checkbox: FilterCheckbox,
	Number:   FilterNumber,
	Date:     FilterDate,
}

type Sort struct {
	Property  string
	Timestamp Timestamp
	Direction Direction
}

type Timestamp string

const (
	CreatedTime    = Timestamp("created_time")
	LastEditedTime = Timestamp("last_edited_time")
)

type Direction string

const (
	Ascending  = Direction("ascending")
	Descending = Direction("descending")
)
