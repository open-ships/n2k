package pgn

// physicalJSONValue is a derived measurement, never an input to the wire codec.
// Keep the unit label even when empty (a dimensionless scaled quantity).
type physicalJSONValue struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

func physicalJSONMeasurement(unit string, getter func() (float64, bool)) *physicalJSONValue {
	value, ok := getter()
	if !ok {
		return nil
	}
	return &physicalJSONValue{Value: value, Unit: unit}
}
