package util

const (
	RMB = "RMB"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case RMB, USD:
		return true
	}
	return false
}
