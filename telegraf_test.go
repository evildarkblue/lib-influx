package influxorm

import (
	"fmt"
	"testing"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type testPoint struct {
	point *write.Point
}

func newTestPoint() *testPoint {
	return &testPoint{
		point: write.NewPointWithMeasurement("test").SetTime(time.Now()),
	}
}

func (tp *testPoint) TimePrecision() time.Duration {
	return time.Second
}

func (tp *testPoint) GetPoint() *write.Point {
	return tp.point
}

func (tp *testPoint) AddTag(t, v string) {
	tp.point.AddTag(t, v)
}

func (tp *testPoint) AddField(f string, v any) {
	tp.point.AddField(f, v)
}

func TestGetInfluxLineProtocalString(t *testing.T) {

	// Test case for empty metrics
	t.Run("EmptyMetrics", func(t *testing.T) {
		p1 := newTestPoint()
		p1.AddField("value", 0.95)
		p1.AddTag("name", "nameA")
		p1.AddTag("host", "server1")

		p2 := newTestPoint()
		p2.AddField("value", 0.91)
		p2.AddTag("host", "server2")
		p2.AddTag("name", "nameB")
		metrics := []InfluxPoint{
			p1, p2,
		}
		result := GetInfluxLineProtocalBody(metrics)
		fmt.Println(string(result))
		expected := ""
		if string(result) != expected {
			t.Errorf("Expected: %s, but got: %s", expected, result)
		}
	})
}
