package influxorm

import (
	"context"
	"errors"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DI interface {
	NewClient() (influxdb2.Client, error)
	GetOrg() string
	GetBucket() string
	NewTelegrafWriter() (TelegrafWriter, error)
	Init(context.Context) error
}

type Config struct {
	Url         string `yaml:"url,omitempty" mapstructure:"url,omitempty"`
	Org         string `yaml:"org,omitempty" mapstructure:"org,omitempty"`
	Token       string `yaml:"token,omitempty" mapstructure:"token,omitempty"`
	Bucket      string `yaml:"bucket,omitempty" mapstructure:"bucket,omitempty"`
	TelegrafUrl string `yaml:"telegraf_url,omitempty" mapstructure:"telegraf_url,omitempty"`
	Timeout     uint   `yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`
}

func (c *Config) NewClient() (influxdb2.Client, error) {
	if c.Url == "" {
		return nil, errors.New("missing url")
	}
	if c.Token == "" {
		return nil, errors.New("missing token")
	}
	if c.Timeout == 0 {
		c.Timeout = 300
	}
	return influxdb2.NewClientWithOptions(c.Url, c.Token,
		influxdb2.DefaultOptions().SetHTTPRequestTimeout(c.Timeout),
	), nil
}

func (c *Config) GetBucket() string {
	return c.Bucket
}

func (c *Config) GetOrg() string {
	return c.Org
}

func (c *Config) NewTelegrafWriter() (TelegrafWriter, error) {
	if c.TelegrafUrl == "" {
		return nil, errors.New("missing telegraf_url")
	}
	return NewTelegraf(c.TelegrafUrl), nil
}

const ENV_INFLUXDB_INIT_FILE = "INFLUXDB_INIT_FILE"

func (c *Config) Init(ctx context.Context) error {
	initFile := os.Getenv(ENV_INFLUXDB_INIT_FILE)

	if initFile == "" {
		return errors.New("missing env INFLUXDB_INIT_FILE")
	}
	initConf, err := loadInitConfig(initFile)
	if err != nil {
		return err
	}

	clt, err := c.NewClient()
	if err != nil {
		return err
	}
	defer clt.Close()

	return initConf.Init(ctx, clt, c.Org)
}
