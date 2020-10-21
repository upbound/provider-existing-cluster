module github.com/turkenh/provider-existing-cluster

go 1.13

require (
	github.com/crossplaneio/crossplane v0.8.0
	github.com/crossplaneio/crossplane-runtime v0.5.0
	github.com/crossplaneio/crossplane-tools v0.0.0-20200214190114-c7c4365eeb95
	github.com/pkg/errors v0.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/controller-tools v0.2.4
	sigs.k8s.io/yaml v1.1.0
)
