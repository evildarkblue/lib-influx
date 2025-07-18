package influxorm

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type telegraf struct {
	url string
}

type TelegrafWriter interface {
	Write(metrics []InfluxPoint) error
}

func NewTelegraf(url string) TelegrafWriter {
	return &telegraf{
		url: url,
	}
}

func GetInfluxLineProtocalBody(metrics []InfluxPoint) []byte {
	builder := &strings.Builder{}
	for _, metric := range metrics {
		write.PointToLineProtocolBuffer(
			metric.GetPoint().SortTags().SortFields(),
			builder, metric.TimePrecision())
	}

	return []byte(builder.String())
}

func (t *telegraf) Write(metrics []InfluxPoint) error {
	data := GetInfluxLineProtocalBody(metrics)
	r, err := http.NewRequest("POST", t.url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return errors.New(string(bodyBytes))
}
