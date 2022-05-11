module github.com/comfysweet/grpc-go-basic

go 1.14

require (
	github.com/go-logr/logr v0.4.0
	github.com/golang/protobuf v1.4.3
	github.com/pkg/errors v0.9.1
	go.mongodb.org/mongo-driver v1.4.0
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
	sber-tech.com/synapse/synai-gen-go/dashboard/autoscaler v1.0.0-beta.20
	sigs.k8s.io/controller-runtime v0.8.3
)

replace sber-tech.com/synapse/synai-gen-go/dashboard/autoscaler => sber-tech.com/synapse/synai-gen-go.git/dashboard/autoscaler v1.0.0-beta.20
