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
	json := `{}`
	json, _ = sjson.Set(json, "source.id", arg.deviceId)
	json, _ = sjson.Set(json, "type", arg.measurementType)
	mTime := arg.time
	if mTime.IsZero() {
		mTime = time.Now()
	}
	json, _ = sjson.Set(json, "time", ToRFCTimeStamp(mTime))
	for _, series := range arg.seriesFloats {
		json, _ = sjson.Set(json, series.Fragment+"."+series.Series+".value", series.Value)
		if len(series.Unit) > 0 {
			json, _ = sjson.Set(json, series.Fragment+"."+series.Series+".unit", series.Unit)
		}
	}

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
		return &c8y.Response{}, createdMeasurement, fmt.Errorf("Received unexpected status code: %d", resp.StatusCode())
	}
	return resp, createdMeasurement, nil
}
