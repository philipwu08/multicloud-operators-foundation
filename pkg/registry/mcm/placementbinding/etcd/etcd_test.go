package etcd

import (
	"testing"

	mcm "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm"
	hcm "github.com/open-cluster-management/multicloud-operators-foundation/pkg/apis/mcm/install"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistrytest "k8s.io/apiserver/pkg/registry/generic/testing"
	etcd3testing "k8s.io/apiserver/pkg/storage/etcd3/testing"
)

func newStorage(t *testing.T) (*REST, *etcd3testing.EtcdTestServer) {
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	hcm.Install(scheme)
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1alpha1"})
	unversioned := schema.GroupVersion{Group: "mcm.ibm.com", Version: "v1alpha1"}
	scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
		&metav1.ExportOptions{},
		&metav1.WatchEvent{},
		&metav1beta1.Table{})
	server, etcdStorage := etcd3testing.NewUnsecuredEtcd3TestClientServer(t)

	etcdStorage.Codec = codecs.LegacyCodec(
		schema.GroupVersion{Group: "mcm.ibm.com", Version: "v1alpha1"})

	restOptions := generic.RESTOptions{
		StorageConfig:           etcdStorage,
		Decorator:               generic.UndecoratedStorage,
		DeleteCollectionWorkers: 1,
		ResourcePrefix:          "placementbindings",
	}
	rest := NewREST(restOptions)
	return rest, server
}

func validNewPlacementBindings() *mcm.PlacementBinding {
	return &mcm.PlacementBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pb1",
		},
		PlacementPolicyRef: mcm.PlacementPolicyRef{
			Name:     "pp",
			Kind:     "placementPolicy",
			APIGroup: "mcm.ibm.com",
		},
	}
}

func TestGet(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)
	test.TestGet(validNewPlacementBindings())

	storage.ShortNames()
}