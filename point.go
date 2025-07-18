package influxorm

import (
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxPoint interface {
	GetPoint() *write.Point
	TimePrecision() time.Duration
}
