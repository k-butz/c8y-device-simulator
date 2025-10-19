package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/destel/rill"
	"github.com/joho/godotenv"
	"github.com/kbutz/c8y-device-simulator/pkg/app"
	"github.com/spf13/viper"

	_ "github.com/dimiro1/banner/autoload"
	"github.com/reubenmiller/go-c8y/pkg/c8y"
)

func main() {
	configureLogger()
	loadEnvFile()
	loadConfig()
	setDefaultConfigs()

	slog.Info("Current settings", "settings", fmt.Sprintf("%v", viper.AllSettings()))

	// requires env vars: C8Y_HOST, C8Y_TENANT, C8Y_USER, C8Y_PASSWORD
	client := c8y.NewClientFromEnvironment(nil, false)

	// onboard devices in a batch of 3
	ctWorkersRampUp := 3
	rill.ForEach(generateRillChan(0, viper.GetInt("countDevices")), ctWorkersRampUp, func(i int) error {
		device := app.NewDevice(fmt.Sprintf(viper.GetString("deviceIdTemplate"), i), client)

		// Queries device id for serial, creates new device if not existing
		createdMo, err := device.InitC8yDevice()
		if err != nil {
			slog.Error("Error while initializing C8Y Device ID. Skipping this Device", "serial", device.Serial, "err", err)
			return nil
		}

		intervalMs := viper.GetInt("deviceSendingIntervalMs")
		initialWaitTimeMs := i * 100
		// non-blocking routine to start device simulation
		device.Run(intervalMs, initialWaitTimeMs, false)

		slog.Info("Created Device simulation",
			"serial", device.Serial,
			"createdMo", createdMo,
			"c8yDeviceId", device.C8yDeviceId,
			"intervalMs", intervalMs,
			"initialWaitTimeMs", initialWaitTimeMs,
		)

		// having wait time is recommended to flatten data ingestion curve
		time.Sleep(time.Duration(viper.GetInt("deviceAddWaitTimeMs")) * time.Millisecond)
		return nil
	})

	// keep main routine alive
	select {}
}

func generateRillChan(min, max int) <-chan rill.Try[int] {
	slice := []int{}
	for i := min; i < max; i++ {
		slice = append(slice, i)
	}
	return rill.FromSlice(slice, nil)
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "err", err)
	}
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func setDefaultConfigs() {
	viper.SetDefault("countDevices", 10)
	viper.SetDefault("deviceIdTemplate", "dev-serial-%05d")
	viper.SetDefault("deviceSendingIntervalMs", 30000)
	viper.SetDefault("deviceAddWaitTimeMs", 500)
}

func configureLogger() {
	c8y.SilenceLogger()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String(slog.TimeKey, app.ToRFCTimeStamp(a.Value.Time()))
			}
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
