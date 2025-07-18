package influxorm

import (
	"context"
	"fmt"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gopkg.in/yaml.v2"
)

type Init interface {
	Init(ctx context.Context) error
}

type InitConfig struct {
	Buckets []*BucketRetention `yaml:"buckets,omitempty"`
	Tasks   []*Task            `yaml:"tasks,omitempty"`
}

func loadInitConfig(path string) (*InitConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config InitConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *InitConfig) Init(ctx context.Context, clt influxdb2.Client, orgName string) error {
	org, err := clt.OrganizationsAPI().FindOrganizationByName(ctx, orgName)
	if err != nil {
		return err
	}
	for _, bucket := range c.Buckets {
		if exist, err := bucket.IsExist(ctx, clt.BucketsAPI()); err != nil {
			return err
		} else if !exist {
			if err := bucket.Create(ctx, clt.BucketsAPI(), org); err != nil {
				return err
			}
		}
		fmt.Println("bucket created", bucket.Name)
	}

	for _, task := range c.Tasks {
		if exist, err := task.IsExist(ctx, clt.TasksAPI()); err != nil {
			return err
		} else if !exist {
			if err := task.Create(ctx, clt.TasksAPI(), org); err != nil {
				return err
			}
		}
	}
	return nil
}
