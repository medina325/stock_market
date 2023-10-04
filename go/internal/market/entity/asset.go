package entity

// Asset represents a financial asset.
//
// It includes information about the asset's unique identifier (ID), name, and
// market volume. An asset is a financial instrument or entity that can be
// traded, such as stocks, bonds, or commodities.
type Asset struct {
	ID           string
	Name         string
	MarketVolume int
}

func NewAsset(id string, name string, marketV int) *Asset {
	return &Asset{
		ID:           id,
		Name:         name,
		MarketVolume: marketV,
	}
}
