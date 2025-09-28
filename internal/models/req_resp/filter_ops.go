package reqresp

type ShopFilter struct {
	Title string
}

type ProductFilter struct {
	Title   string
	MaxCost uint64
	MinCost uint64
}

type CategoryFilter struct {
	Title string
}
