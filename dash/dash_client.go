package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"io"
	"sber-tech.com/synapse/synai-gen-go/dashboard/autoscaler/api"
	"sber-tech.com/synapse/synai-gen-go/dashboard/autoscaler/model"
	"time"
)

type DashboardAutoscalerClient struct {
	initClientFunc func() (api.SynaiDashboardAutoscalerClient, io.Closer, error)
}

func NewDashboardAutoscalerClient() *DashboardAutoscalerClient {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	return &DashboardAutoscalerClient{
		initClientFunc: func() (api.SynaiDashboardAutoscalerClient, io.Closer, error) {
			connection, err := grpc.Dial("localhost:5656", opts...)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Не удалось подключиться к dashboard")
			}
			return api.NewSynaiDashboardAutoscalerClient(connection), connection, nil
		}}
}

func (dashCon *DashboardAutoscalerClient) ChangeMode(ctx context.Context, request *model.ChangeModeRequest) error {
	client, conn, err := dashCon.initClientFunc()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := client.ChangeMode(ctx, request)
	if err != nil {
		return errors.Wrap(err, "Произошла ошибка при получении ответа от dashboard по запросу DeleteFromStore")
	}
	if res != nil {
		fmt.Printf("Получен ответ на запрос ChangeMode от dashboard со статус-кодом: %v", res.Status.String())
	}
	return nil
}

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*300)
	defer cancelFunc()
	client := NewDashboardAutoscalerClient()
	req := &model.ChangeModeRequest{
		Namespace: "ci01994970-edevgen-synai-dev",
		WorkMode:  model.ScalerWorkMode_RECOMMENDATION,
	}
	if err := client.ChangeMode(ctx, req); err != nil {
		fmt.Println(err)
	}
}
