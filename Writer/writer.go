package writer

import (
	config "NeSyCa/Config"
	"bytes"
	"crypto/rand"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const MTU = 65507

var updToken []byte
var httpToken []byte

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
	req, err := http.NewRequest(http.MethodPost, *target, data)
	if err != nil {
		logWithField.Error(err)
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logWithField.Error(err)
		return 0, err
	}
	logWithField.WithField("h", resp.Request.ContentLength).WithField("s", size).Info()
	writeToChannl(&ch, &size)
	return size, nil
}

func init() {
	updToken = make([]byte, MTU)
	rand.Read(updToken)
	httpToken = make([]byte, int64(config.Config().SizeMax))
}
