package entity

// Investor represents information about an individual investor.
type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}

// NewInvestor creates a new Investor instance with the specified ID.
//
// It initializes a new Investor struct with the provided ID and an empty list
// of asset positions. The returned pointer points to the newly created Investor.
//
// Parameters:
//   - id: The unique identifier for the investor.
//
// Returns:
//   - *Investor: A pointer to the newly created Investor instance.
func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

// AddAssetPosition appends a new asset position to the investor's list of positions.
//
// It adds the provided asset position to the list of positions held by the investor.
//
// Parameters:
//   - assetPosition: A pointer to the InvestorAssetPosition to add to the list.
func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

// UpdateAssetPosition updates an investor's position in a specific asset.
//
// It first attempts to find the asset position for the given asset ID.
// If a position for the asset does not exist, it creates
// a new position with the provided asset ID and shares count and adds it to the
// investor's list of asset positions.
// If a position for the asset already exists,
// it adds the provided shares count to the existing position's shares count.
//
// Parameters:
//   - assetID: The unique identifier of the asset to update.
//   - sharesCount: The number of shares (or "cotas") to add or subtract from
//     the investor's position in the asset.
func (i *Investor) UpdateAssetPosition(assetID string, sharesCount int) {
	// Attempt to find the asset position for the given asset ID.
	assetPosition := i.GetAssetPosition(assetID)

	if assetPosition == nil {
		// If a position doesn't exist, create a new one and add it to the list.
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, sharesCount))
		return
	}

	// If a position already exists, update the shares count.
	assetPosition.Shares += sharesCount
}

// GetAssetPosition retrieves the asset position for a specific asset by its ID.
//
// It searches through the investor's list of asset positions and returns the
// first position that matches the provided asset ID. If no matching position
// is found, it returns nil.
//
// Parameters:
//   - assetID: The unique identifier of the asset to retrieve.
//
// Returns:
//   - *InvestorAssetPosition: A pointer to the InvestorAssetPosition if found,
//     or nil if no matching position is found.
func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}
	return nil
}

// InvestorAssetPosition represents an investor's position in a specific asset.
//
// It includes information about the asset's unique identifier (AssetID) and the
// number of shares (or "cotas") the investor holds in that asset.
type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  shares,
	}
}
