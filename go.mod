module github.com/turkenh/stack-existing-cluster

go 1.12

require (
	github.com/crossplaneio/crossplane v0.6.0-rc.0.20191226165033-a452562456e0
	github.com/crossplaneio/crossplane-runtime v0.3.1-0.20200115232149-cd8c52b483c3
	github.com/crossplaneio/crossplane-tools v0.0.0-20191220202319-9033bd8a02ce
	github.com/pkg/errors v0.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/controller-tools v0.2.4
	sigs.k8s.io/yaml v1.1.0
)
