package app

type MeasurementSeriesFloat struct {
	Fragment string
	Series   string
	Unit     string
	Value    float64
}

func NewMeasurementSeriesFloatRandomized(fragment, series, unit string, min, max, jitter float64) MeasurementSeriesFloat {
	return MeasurementSeriesFloat{
		Fragment: fragment,
		Series:   series,
		Unit:     unit,
		Value:    RandFloat(min, max, jitter),
	}
}
