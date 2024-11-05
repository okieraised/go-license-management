package constants

const (
	ProductDistributionStrategyOpen     = "open"
	ProductDistributionStrategyClosed   = "closed"
	ProductDistributionStrategyLicensed = "licensed"
)

var ValidProductDistributionStrategyMapper = map[string]bool{
	ProductDistributionStrategyOpen:     true,
	ProductDistributionStrategyClosed:   true,
	ProductDistributionStrategyLicensed: true,
}
