package influxorm

import (
	"context"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type Task struct {
	Name        string `yaml:"name,omitempty"`
	FromBucket  string `yaml:"from_bucket,omitempty"`
	ToBucket    string `yaml:"to_bucket,omitempty"`
	Measurement string `yaml:"measurement,omitempty"`
	TaskFlux    string `yaml:"task_flux,omitempty"`
}

func (t *Task) IsExist(ctx context.Context, taskApi api.TasksAPI) (bool, error) {
	tasks, err := taskApi.FindTasks(ctx, &api.TaskFilter{
		Name: t.Name,
	})
	if err != nil {
		return false, err
	}
	if len(tasks) > 0 {
		return true, nil
	}
	return false, nil
}

func (t *Task) Create(ctx context.Context, taskApi api.TasksAPI, org *domain.Organization) error {
	// fmt.Println(fmt.Sprintf(t.TaskFlux, t.FromBucket, t.Measurement, t.ToBucket, org.Name))
	_, err := taskApi.CreateTaskByFlux(ctx,
		fmt.Sprintf(t.TaskFlux, t.FromBucket, t.Measurement, t.ToBucket, org.Name),
		*org.Id)
	// _, err := taskApi.CreateTaskWithCron(ctx, t.Name,
	// 	fmt.Sprintf(t.TaskFlux, t.FromBucket, t.Measurement, t.ToBucket, org.Name),
	// 	t.Cron, *org.Id)
	// fmt.Println("aaaa", err)
	return err

}
