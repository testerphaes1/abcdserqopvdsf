package repos

// FilterOp filter operation
type FilterOp string

const (
	// FilterOpEq equal
	FilterOpEq = "="
	// FilterOpNeq not equal
	FilterOpNeq = "!="
	// FilterOpGt greater than
	FilterOpGt = ">"
	// FilterOpGte greater than equal
	FilterOpGte = ">="
	// FilterOpLt less than
	FilterOpLt = "<"
	// FilterOpLte less than equal
	FilterOpLte = "<="
	// FilterOpIn within
	FilterOpIn = "IN"
)

// Filter represent a filter
type Filter struct {
	Field string
	Op    FilterOp
	Value interface{}
}

// Filters filters
type Filters []Filter
