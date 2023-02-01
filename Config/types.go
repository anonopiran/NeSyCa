package config

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type LogLevelType string
type UDPAddressType string
type SettingsType struct {
	UDPTargets    []UDPAddressType `env:"UDP_TARGETS" env-required:"true"`
	SizeMin       float32          `env:"SIZE_MIN" env-default:"500000"`
	SizeMax       float32          `env:"SIZE_MAX" env-default:"2000000"`
	RateLambda    float64          `env:"RATE_LAMBDA" env-default:"1"`
	StatsInterval int              `env:"STATS_INTERVAL" env-default:"30"`
	LogLevel      LogLevelType     `env:"LOG_LEVEL" env-default:"warning"`
}

func (f *UDPAddressType) AsString() string {
	return string(*f)
}
func (f *LogLevelType) SetValue(s string) error {
	LogWithRaw := log.WithField("value", s)
	ll, err := log.ParseLevel(s)
	if err != nil {
		LogWithRaw.Error(err)
		return err
	}
	log.SetLevel(ll)
	*f = LogLevelType(s)
	return nil
}
func (f *UDPAddressType) SetValue(s string) error {
	// LogWithRaw := log.WithField("value", s)
	_, err := net.Dial("udp", s)
	if err != nil {
		return err
	}
	*f = UDPAddressType(s)
	return nil
}
