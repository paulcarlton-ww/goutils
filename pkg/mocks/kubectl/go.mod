module github.com/paulcarlton-ww/goutils/pkg/mocks/kubectl

go 1.15

require (
	github.com/go-logr/logr v0.4.0
	github.com/golang/mock v1.5.0
)

replace (
	github.com/paulcarlton-ww/goutils/pkg/kubectl => ../../kubectl
)
