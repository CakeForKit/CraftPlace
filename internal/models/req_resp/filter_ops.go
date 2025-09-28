package reqresp

import "github.com/google/uuid"

type ShopFilter struct {
	Title  string    // default = ""
	UserID uuid.UUID // default = uuid.Nil
}

type ProductFilter struct {
	Title      string    // default = ""
	MaxCost    uint64    // default = 0
	MinCost    uint64    // default = uint64(math.MaxUint64)
	ShopID     uuid.UUID // default = uuid.Nil
	CategoryID uuid.UUID // default = uuid.Nil
}

type CategoryFilter struct {
	Title string // default = ""
}
