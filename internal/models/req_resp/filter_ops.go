package reqresp

import "github.com/google/uuid"

type ShopFilter struct {
	Title  string    // default = ""
	UserID uuid.UUID // default = uuid.Nil
	Offset uint64    // default = 0
	Limit  uint64    // default = 0 // if 0: return all
}

type ProductFilter struct {
	Title      string    // default = ""
	MinCost    uint64    // default = uint64(math.MaxUint64)
	MaxCost    uint64    // default = 0
	ShopID     uuid.UUID // default = uuid.Nil
	CategoryID uuid.UUID // default = uuid.Nil
	Offset     uint64    // default = 0
	Limit      uint64    // default = 0 // if 0: return all
}

type CategoryFilter struct {
	Title  string // default = ""
	Offset uint64 // default = 0
	Limit  uint64 // default = 0 // if 0: return all
}

type PostFilter struct {
	ShopID uuid.UUID // default = uuid.Nil
	Offset uint64    // default = 0
	Limit  uint64    // default = 0 // if 0: return all
}
