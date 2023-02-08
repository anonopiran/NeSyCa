package writer

import (
	"crypto/rand"
	"net"

	log "github.com/sirupsen/logrus"
)

const MTU = 65507

var token = make([]byte, MTU)

func writeToChannl(ch *chan<- int, val *int) {
	*ch <- *val
}
func Write(target *string, size int, ch chan<- int) (int, error) {
	logWithField := log.WithField("target", *target)
	conn, err := net.Dial("udp", *target)
	if err != nil {
		logWithField.Error(err)
		return 0, err
	}
	wrote := 0
	defer writeToChannl(&ch, &wrote)
	defer conn.Close()
	for size > 0 {
		w, err := conn.Write(token)
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

func init() {
	rand.Read(token)
}
