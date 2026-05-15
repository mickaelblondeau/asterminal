package config

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Lat, Lon, Zoom float64
	TimeOffset     int64
	Track          string
}

func getArgValue(option string, arg string) (float64, bool) {
	if strings.HasPrefix(arg, "--"+option+"=") {
		if _, value, found := strings.Cut(arg, "="); found {
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				return v, true
			}
		}
	}

	return 0, false
}

func getArgTimestampValue(option string, arg string) (int64, bool) {
	if strings.HasPrefix(arg, "--"+option+"=") {
		if _, value, found := strings.Cut(arg, "="); found {
			t, err := time.Parse("2006-01-02T15:04:05", value)
			if err == nil {
				return t.Unix(), true
			}

			t, err = time.Parse("2006-01-02", value)
			if err == nil {
				return t.Unix(), true
			}
		}
	}

	return 0, false
}

func getArgStringValue(option string, arg string) (string, bool) {
	if strings.HasPrefix(arg, "--"+option+"=") {
		if _, value, found := strings.Cut(arg, "="); found {
			return value, true
		}
	}

	return "", false
}

func GetConfig() Config {
	config := Config{}
	args := os.Args[1:]

	for _, arg := range args {
		if value, found := getArgValue("lat", arg); found {
			config.Lat = math.Max(-90, math.Min(90, value))
		}

		if value, found := getArgValue("lon", arg); found {
			config.Lon = math.Max(-180, math.Min(180, value))
		}

		if value, found := getArgValue("zoom", arg); found {
			config.Zoom = math.Max(0, math.Min(100, value))
		}

		if value, found := getArgStringValue("track", arg); found {
			config.Track = strings.ToLower(value)
		}

		if value, found := getArgTimestampValue("date", arg); found {
			config.TimeOffset = value - time.Now().Unix()
		}
	}

	return config
}
