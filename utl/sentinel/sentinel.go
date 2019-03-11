// Package sentinel renames the empty struct to make it more convenient to use as a sentinel value.
package sentinel

// SNT is a sentinel type: it doesn't carry any data but is used only as a signal.
type SNT struct{}

// S returns a sentinel instance.
func S() SNT {
	return SNT{}
}
