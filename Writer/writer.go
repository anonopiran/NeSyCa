package writer

import (
	config "NeSyCa/Config"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const MTU = 65507

var updToken []byte
var httpToken []byte
var httpClient *http.Client

func writeToChannl(ch *chan<- int, val *int) {
	*ch <- *val
}
func WriteUDP(target *string, size int, ch chan<- int) (int, error) {
	logWithField := log.WithField("target", *target).WithField("type", "udp")
	conn, err := net.Dial("udp", *target)
	if err != nil {
		logWithField.Error(err)
		return 0, err
	}
	wrote := 0
	defer writeToChannl(&ch, &wrote)
	defer conn.Close()
	for size > 0 {
		w, err := conn.Write(updToken)
		if err != nil {
			logWithField.Error(err)
			return wrote, err
		}
		size -= MTU
		wrote += w
		logWithField.Debug("MTU batch written")
	}
	return wrote, nil
}

func WriteHTTP(target *string, size int, ch chan<- int) (int, error) {
	logWithField := log.WithField("target", *target).WithField("type", "http")
	data := bytes.NewReader(httpToken[:size])
	_, err := httpClient.Post(*target, "application/x-binary", data)
	if err != nil {
		logWithField.Error(err)
		return 0, err
	}
	writeToChannl(&ch, &size)
	return size, nil
}

func init() {
	updToken = make([]byte, MTU)
	rand.Read(updToken)
	httpToken = make([]byte, int64(config.Config().SizeMax))
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: tr}
}
