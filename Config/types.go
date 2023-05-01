package config

import (
	"fmt"
	net "net"
	"strings"

	"strconv"

	log "github.com/sirupsen/logrus"
)

type LogLevelType string
type TCPAddressType net.TCPAddr
type SizeType int
type SettingsType struct {
	TCPTargets   []TCPAddressType `koanf:"tcp_targets" validate:"required,dive,required,tcp4_addr"`
	RateLimitMin SizeType         `koanf:"rate_limit_min" validate:"required"`
	RateLimitMax SizeType         `koanf:"rate_limit_max" validate:"required"`
	BatchSizeMin SizeType         `koanf:"batch_size_min"`
	BatchSizeMax SizeType         `koanf:"batch_size_max"`
	FileSizeMin  SizeType         `koanf:"file_size_min"`
	FileSizeMax  SizeType         `koanf:"file_size_max"`
	LogLevel     LogLevelType     `koanf:"log_level"`
	Timeout      int              `koanf:"timeout"`
}

func (f *TCPAddressType) UnmarshalText(text []byte) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", string(text))
	if err != nil {
		return err
	}
	*f = TCPAddressType(*tcpAddr)
	return nil
}
func (f *LogLevelType) UnmarshalText(text []byte) error {
	ll, err := log.ParseLevel(string(text))
	if err != nil {
		return err
	}
	log.SetLevel(ll)
	*f = LogLevelType(text)
	return nil
}
func (f *SizeType) UnmarshalText(text []byte) error {
	val := string(text)
	unit := strings.ToLower(strings.Split(val, "")[len(val)-1])
	v_, err := strconv.Atoi(strings.Join(strings.Split(val, "")[:len(val)-1], ""))
	if err != nil {
		return err
	}
	var prd int
	switch unit {
	case "b":
		prd = 1
	case "k":
		prd = 1 << 10
	case "m":
		prd = 1 << 20
	case "g":
		prd = 1 << 30
	default:
		return fmt.Errorf("%s unit not understood", unit)
	}
	*f = SizeType(prd * v_)
	return nil
}
