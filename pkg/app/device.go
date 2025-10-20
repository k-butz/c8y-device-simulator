package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/reubenmiller/go-c8y/pkg/c8y"
)

type Device struct {
	Serial      string
	client      *c8y.Client
	C8yDeviceId string
	currentTick int
}

func NewDevice(id string, client *c8y.Client) *Device {
	return &Device{
		Serial:      id,
		client:      client,
		currentTick: 1,
	}
}

func (device *Device) Run(intervalMs, initialWaitTimeMs int, firstExecuteOnSchedule bool) {
	go func() {
		time.Sleep(time.Duration(initialWaitTimeMs) * time.Millisecond)
		// setup ticker in provided interval
		interval := time.Duration(intervalMs) * time.Millisecond
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// routine stops after 24 hours
		done := time.After(24 * time.Hour)

		// each of these functions will be executed in each interval
		fns := collectFunctions(device)

		tickFunction := func() {
			device.incrementTick()
			for _, fn := range fns {
				fn()
			}
		}

		if firstExecuteOnSchedule {
			tickFunction()
		}

		for {
			select {
			case <-ticker.C:
				tickFunction()
			case <-done:
				slog.Info("Simulation finished", "serial", device.Serial)
				return
			}
		}
	}()
}

// here is where you define the executed functions for each interval/tick
func collectFunctions(device *Device) []func() {
	res := []func(){}

	// function to create measurements
	res = append(res, func() {
		resp, m, err := CreateMeasurement(MeasurementCreateArg{
			client:          device.client,
			deviceId:        device.C8yDeviceId,
			measurementType: "simulated",
			seriesFloats: []MeasurementSeriesFloat{
				// creates a random float between 0 and 100
				NewMeasurementSeriesFloatRandomized("c8y_Temperature", "T", "", 0, 100, 0),
				// creates a random float between 95 and 105 (target value = 100, jitter = +- 5%)
				NewMeasurementSeriesFloatRandomized("c8y_Pressure", "P", "", 100, 100, 5),
			},
		})
		if err != nil {
			slog.Error("Error while creating Measurement", "serial", device.Serial, "err", err)
		}
		slog.Info("Created measurement",
			"serial", device.Serial,
			"c8yDeviceId", m.Item.Get("source.id"),
			"measurementTs", m.Item.Get("time"),
			"statusCode", resp.Response.StatusCode,
			"durationMs", resp.Duration().Milliseconds(),
			"host", resp.Response.Request.Host)
	})

	// this could be another function that is executed only every second interval
	// res = append(res, func(){
	// 	if device.currentTick % 2 == 0 {
	// 		// do something useful
	// 	}
	// })

	return res
}

func (d *Device) incrementTick() {
	if d.currentTick == 100 {
		d.currentTick = 1
		return
	}
	d.currentTick += 1
}

func (d *Device) InitC8yDevice() (bool, error) {
	identity, _, err := d.client.Identity.GetExternalID(context.Background(), "c8y_Serial", d.Serial)
	if err == nil && identity != nil && len(identity.ExternalID) > 0 {
		d.C8yDeviceId = identity.ManagedObject.ID
		return false, nil
	}

	mo, _, err := d.client.Inventory.Create(context.Background(), &c8y.Device{
		ManagedObject: c8y.ManagedObject{
			Name: d.Serial,
			Type: "simulatedDevice",
		},
		DeviceFragment: c8y.DeviceFragment{},
	})
	if err != nil {
		return true, err
	}
	if mo == nil || len(mo.ID) == 0 {
		return true, fmt.Errorf("Newly created Managed Object for Serial %s is null", d.Serial)
	}
	newExtId, _, err := d.client.Identity.Create(context.Background(), mo.ID, "c8y_Serial", d.Serial)
	if err != nil {
		return true, err
	}
	if newExtId == nil || len(newExtId.ExternalID) == 0 {
		return true, fmt.Errorf("Created Managed Object with ID %s, but failed to created external ID for it", mo.ID)
	}

	d.C8yDeviceId = mo.ID

	return true, nil
}
