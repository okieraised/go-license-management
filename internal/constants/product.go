package constants

const (
	// ProductDistributionStrategyOpen - Anybody can access releases. No API authentication required, so this is a great option for rendering releases on a public downloads page, open source projects, or freemium products.
	ProductDistributionStrategyOpen = "open"
	// ProductDistributionStrategyClosed - Only admins can access releases. Download links will need to be generated server-side. API authentication is required.
	ProductDistributionStrategyClosed = "closed"
	// ProductDistributionStrategyLicensed - Only licensed users, with a valid license, can access releases and release artifacts. API authentication is required.
	ProductDistributionStrategyLicensed = "licensed"
)

var ValidProductDistributionStrategyMapper = map[string]bool{
	ProductDistributionStrategyOpen:     true,
	ProductDistributionStrategyClosed:   true,
	ProductDistributionStrategyLicensed: true,
}
