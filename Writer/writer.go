package writer

import (
	"crypto/rand"
	"fmt"
	"net"
	"time"

	"context"

	utils "NeSyCa/Utils"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

func WriteTCP(target net.TCPAddr, lenght int, rateLimit int, batchSize int) error {
	conn, err := net.Dial("tcp", target.String())
	if err != nil {
		return err
	}
	limiter := rate.NewLimiter(rate.Limit(rateLimit), rateLimit)
	sTime := time.Now().UnixMicro()
	count, err := doWriteTCP(conn, *limiter, lenght, batchSize)
	if err != nil {
		logrus.Error(err)
	}
	dur := time.Now().UnixMicro() - sTime
	var rate float32
	if dur > 0 {
		rate = float32(count) / float32(dur) * 1e6
	} else {
		rate = 0
	}
	logrus.WithField("target", target).
		WithField("size", fmt.Sprintf("%s/%s", utils.ByteCountIEC(count), utils.ByteCountIEC(lenght))).
		WithField("rate", fmt.Sprintf("%s/%s", utils.ByteCountIEC(int(rate)), utils.ByteCountIEC(rateLimit))).
		WithField("batch size", utils.ByteCountIEC(batchSize)).
		Info()
	return err
}
func doWriteTCP(conn net.Conn, limiter rate.Limiter, lenght int, batchSize int) (int, error) {
	count := 0
	data := make([]byte, batchSize)
	rand.Read(data)
	for {
		if err := limiter.WaitN(context.Background(), batchSize); err != nil {
			return count, err
		}
		c_, err := conn.Write(data)
		if err != nil {
			return count, err
		}
		logrus.Debugf("wrote %d batch to %s", c_, conn.RemoteAddr())
		count += c_
		if count >= lenght {
			break
		}
	}
	return count, nil
}
