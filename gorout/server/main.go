package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sber-tech.com/synapse/synai-gen-go/dashboard/autoscaler/api"
	"sber-tech.com/synapse/synai-gen-go/dashboard/autoscaler/model"
	"sync"
	"syscall"
)

type DashboardAutoscalerServer struct {
	api.UnimplementedSynaiDashboardAutoscalerServer
}

func NewDashboardAutoscalerServer() *DashboardAutoscalerServer {
	return &DashboardAutoscalerServer{}
}

func (dashCon *DashboardAutoscalerServer) InitServer() {
	lis, err := net.Listen("tcp", ":5656")
	if err != nil {
		fmt.Println(err, "Ошибка запуска gRPC сервера")
	}
	s := grpc.NewServer()
	api.RegisterSynaiDashboardAutoscalerServer(s, dashCon)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//done := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		sign := <-stop
		fmt.Println("Останавливаем GRPC сервер", "Sign", sign)
		s.GracefulStop()
		wg.Done()
	}()

	//go func() {
	//	sig := <-stop
	//	fmt.Println()
	//	fmt.Println(sig)
	//	done <- true
	//}()

	go func() {
		if err := s.Serve(lis); err != nil {
			fmt.Println(err, "Ошибка в работе gRPC сервера")
		}
	}()

	//fmt.Println("awaiting signal")
	//<-done
	//fmt.Println("done")
	fmt.Println("awaiting signal")
	wg.Wait()

	//defer func() {
	//	dashCon.log.Info("Останавливаем GRPC сервер")
	//	s.GracefulStop()
	//}()
	//<-stop

	//defer func() {
	//	s.GracefulStop()
	//}()

	// block until either OS signal, or server fatal error
	//select {
	//case sign := <-stop:
	//	dashCon.log.Info("Останавливаем GRPC сервер", "Sign", sign)
	//}
}

func (dashCon *DashboardAutoscalerServer) InitServerNew() {
	lis, err := net.Listen("tcp", ":5657")
	if err != nil {
		fmt.Println(err, "Ошибка запуска gRPC сервера1")
	}
	s := grpc.NewServer()
	api.RegisterSynaiDashboardAutoscalerServer(s, dashCon)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//done := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		sign := <-stop
		fmt.Println("Останавливаем GRPC сервер1", "Sign", sign)
		s.GracefulStop()
		wg.Done()
	}()

	//go func() {
	//	sig := <-stop
	//	fmt.Println()
	//	fmt.Println(sig)
	//	done <- true
	//}()

	go func() {
		if err := s.Serve(lis); err != nil {
			fmt.Println(err, "Ошибка в работе gRPC сервера1")
		}
	}()

	//fmt.Println("awaiting signal")
	//<-done
	//fmt.Println("done")
	fmt.Println("awaiting signal1")
	wg.Wait()

	//defer func() {
	//	dashCon.log.Info("Останавливаем GRPC сервер")
	//	s.GracefulStop()
	//}()
	//<-stop

	//defer func() {
	//	s.GracefulStop()
	//}()

	// block until either OS signal, or server fatal error
	//select {
	//case sign := <-stop:
	//	dashCon.log.Info("Останавливаем GRPC сервер", "Sign", sign)
	//}
}

func (dashCon *DashboardAutoscalerServer) ChangeMode(ctx context.Context, request *model.ChangeModeRequest) (*model.ChangeModeResponse, error) {
	fmt.Println("Поступил запрос ChangeMode из dashboard")
	response := model.ChangeModeResponse{
		Status: model.DashAutoscalerResponseStatus_Success,
	}
	fmt.Println("Успешно отработал сценарий ChangeMode")
	return &response, nil
}

func main() {
	fmt.Println("Starting the GRPC server")
	// Start grpc server
	server := NewDashboardAutoscalerServer()
	go server.InitServer()
	server.InitServerNew()
}
