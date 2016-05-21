package metar

// InHgTohPa converts "inch of mercury" to "hectopascal"
func InHgTohPa(inHg float64) float64 {
	return inHg * 33.8638866667
}

// KtsToMs converts "knots" to "meters per second"
func KtsToMs(kts float64) float64 {
	return kts * 0.514444
}

// KtsToBft converts "knots" to "bft"
func KtsToBft(kts float64) int {
	switch {
	case kts < 1:
		return 0
	case kts < 4:
		return 1
	case kts < 7:
		return 2
	case kts < 11:
		return 3
	case kts < 16:
		return 4
	case kts < 22:
		return 5
	case kts < 28:
		return 6
	case kts < 34:
		return 7
	case kts < 41:
		return 8
	case kts < 48:
		return 9
	case kts < 56:
		return 10
	case kts < 64:
		return 11
	default:
		return 12
	}
}

// StatMileToKm converts "statute miles" to "kilometers"
func StatMileToKm(sm float64) float64 {
	return sm * 1.60934
}

// MbTohPa converts "millibar" to "hectopascal"
func MbTohPa(mb float64) float64 {
	return mb * 0.1
}
