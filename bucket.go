package influxorm

import (
	"context"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type BucketRetention struct {
	Name      string
	Retention time.Duration
}

func (c *BucketRetention) IsExist(ctx context.Context, bucketApi api.BucketsAPI) (bool, error) {
	bucket, err := bucketApi.FindBucketByName(ctx, c.Name)
	if err == nil && bucket != nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "not found") {
		return false, nil
	}
	return false, err
}

func (c *BucketRetention) Create(ctx context.Context, bucketApi api.BucketsAPI, org *domain.Organization) error {
	_, err := bucketApi.CreateBucketWithName(ctx, org, c.Name, domain.RetentionRule{
		EverySeconds: int64(c.Retention.Seconds()),
	})
	return err
}
