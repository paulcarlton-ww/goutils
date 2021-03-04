module github.com/paulcarlton-ww/goutils/pkg/kubectl

go 1.15

require (
	github.com/go-logr/logr v0.4.0
	github.com/golang/mock v1.5.0
	github.com/paulcarlton-ww/goutils/pkg/logging v0.0.3
	github.com/paulcarlton-ww/goutils/pkg/mocks/kubectl v0.0.3
	github.com/paulcarlton-ww/goutils/pkg/mocks/logr v0.0.3
	github.com/pkg/errors v0.9.1
)

replace (
	github.com/paulcarlton-ww/goutils/pkg/mocks/kubectl => ../mocks/kubectl
	github.com/paulcarlton-ww/goutils/pkg/mocks/logr => ../mocks/logr
)
