module github.com/dweomer/linode-controller-manager/tool/kubeconform

// go get go@1.24.0 toolchain@go1.25

go 1.24.0

toolchain go1.25.10

// go mod edit -tool github.com/yannh/kubeconform/cmd/kubeconform -require github.com/yannh/kubeconform@v0.7.0

tool github.com/yannh/kubeconform/cmd/kubeconform

require (
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.1 // indirect
	github.com/yannh/kubeconform v0.7.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
