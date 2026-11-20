// Package pgnschema resolves wire rules shared by the generator and the runtime.
package pgnschema

import "math"

// UnknownValue returns the field's unavailable raw tick, including the
// signed 64-bit sentinel omitted by the pinned GNSS position definition.
func UnknownValue(variant string, order int, declared *int64) *int64 {
	if declared == nil && variant == "GnssPositionData" && order >= 4 && order <= 6 {
		return new(int64(math.MaxInt64))
	}
	return declared
}
