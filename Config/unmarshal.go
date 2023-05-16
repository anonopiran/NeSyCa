package config

import (
	"fmt"
	net "net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (f *TCPAddressType) UnmarshalText(text []byte) error {
	val := string(text)
	tcpAddr, err := net.ResolveTCPAddr("tcp", val)
	if err != nil {
		return err
	}
	*f = TCPAddressType{Plain: val, Address: *tcpAddr}
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
func (f *SecondDuration) UnmarshalText(text []byte) error {
	val, err := parseDur(text)
	if err != nil {
		return err
	}
	*f = SecondDuration(val * time.Second)
	return nil
}
func (f *MiliSecondDuration) UnmarshalText(text []byte) error {
	val, err := parseDur(text)
	if err != nil {
		return err
	}
	*f = MiliSecondDuration(val * time.Millisecond)
	return nil
}
