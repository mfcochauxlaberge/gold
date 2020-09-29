package gold

// ComparisonError represents a difference between a given output and the
// content of a file.
type ComparisonError struct{}

func (e ComparisonError) Error() string {
	return "output and file are different"
}
