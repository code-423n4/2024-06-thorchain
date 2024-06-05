package common

// RoundSignificantFigures rounds an unsigned 64-bit integer to the specified number of
// significant digits. It takes the number and significant digits as arguments.
func RoundSignificantFigures(number uint64, significantDigits int64) uint64 {
	if number == 0 {
		return 0
	}

	// calculate the magnitude of the number
	magnitude := int64(0)
	temp := number
	for temp != 0 {
		temp /= 10
		magnitude++
	}

	// if the magnitude is less than significant digits, return the number
	if magnitude <= significantDigits {
		return number
	}

	// determine the divisor based on the number of significant digits to keep
	divisor := uint64(1)
	for i := int64(0); i < magnitude-significantDigits; i++ {
		divisor *= 10
	}

	// round the number to the desired significant digits
	return number / divisor * divisor
}
