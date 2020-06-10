package samples

import (
	"fmt"
	"strings"
)

type sample struct {
	Name     string
	Interval string
	Buckets  []string
}

var DefaultSamples = []sample{
	{
		Name:     "short",
		Interval: "30s",
		Buckets:  []string{"5m", "30m", "1h"},
	},
	{
		Name:     "medium",
		Interval: "2m",
		Buckets:  []string{"2h", "6h"},
	},
	{
		Name:     "daily",
		Interval: "5m",
		Buckets:  []string{"1d", "3d"},
	},
}

var disabletBucketsForTickets = []string{"3d", "1d", "2h"}

func IsTicketSample(sample string) bool {
	for _, bucketSample := range disabletBucketsForTickets {
		if bucketSample == sample {
			return true
		}
	}

	return false
}

func ValidateSample(sample string) error {
	validSamples := []string{}
	for _, defaultSample := range DefaultSamples {
		for _, bucket := range defaultSample.Buckets {
			if bucket == sample {
				return nil
			}

			validSamples = append(validSamples, bucket)
		}
	}
	return fmt.Errorf("Sample %s is not a valid sample, valid samples: %s", sample, strings.Join(validSamples, ","))
}
