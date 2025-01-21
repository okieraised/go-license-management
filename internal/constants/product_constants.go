package constants

const (
	// ProductDistributionStrategyOpen - Anybody can access releases. No API authentication required.
	ProductDistributionStrategyOpen = "open"
	// ProductDistributionStrategyClosed - Only admins can access releases. API authentication is required.
	ProductDistributionStrategyClosed = "closed"
	// ProductDistributionStrategyLicensed - Only licensed users, with a valid license, can access releases. API authentication is required.
	ProductDistributionStrategyLicensed = "licensed"
)

var ValidProductDistributionStrategyMapper = map[string]bool{
	ProductDistributionStrategyOpen:     true,
	ProductDistributionStrategyClosed:   true,
	ProductDistributionStrategyLicensed: true,
}
