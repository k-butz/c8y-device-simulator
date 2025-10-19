package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/reubenmiller/go-c8y/pkg/c8y"
	"github.com/tidwall/sjson"
)

type MeasurementCreateArg struct {
	client          *c8y.Client
	seriesFloats    []MeasurementSeriesFloat
	deviceId        string
	measurementType string
	time            time.Time
}

func CreateMeasurement(arg MeasurementCreateArg) (*c8y.Response, *c8y.Measurement, error) {
	json := `{"type":"sim"}`
	json, _ = sjson.Set(json, "source.id", arg.deviceId)
	for _, series := range arg.seriesFloats {
		sjson.Set(json, series.Fragment+"."+series.Series+".value", series.Value)
		if len(series.Unit) > 0 {
			sjson.Set(json, series.Fragment+"."+series.Series+".unit", series.Unit)
		}
	}
	mTime := arg.time
	if mTime.IsZero() {
		mTime = time.Now()
	}
	json, _ = sjson.Set(json, "time", mTime.Format("2006-01-02T15:04:05.000Z07:00"))
	json, _ = sjson.Set(json, "c8y_Temperature.T.value", RandFloat(1, 100, 0))
	json, _ = sjson.Set(json, "c8y_Pressure.P.value", RandFloat(1, 100, 0))

	createdMeasurement := new(c8y.Measurement)
	resp, err := arg.client.SendRequest(context.Background(), c8y.RequestOptions{
		Method:       "POST",
		Path:         "measurement/measurements",
		Body:         json,
		ResponseData: createdMeasurement,
	})
	if err != nil {
		return &c8y.Response{}, createdMeasurement, err
	}
	if resp == nil {
		return &c8y.Response{}, createdMeasurement, errors.New("Platform response is nil")
	}
	if resp.StatusCode() != http.StatusCreated {
		return &c8y.Response{}, createdMeasurement, fmt.Errorf("Received unexpected status code: %s", resp.StatusCode())
	}
	return resp, createdMeasurement, nil
}
