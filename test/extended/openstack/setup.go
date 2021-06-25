package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/openstack/clientconfig"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	exutil "github.com/openshift/origin/test/extended/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
)

const (
	serviceCompute = "compute"
)

func skipIfNotOpenStack(oc *exutil.CLI) {
	g.By("checking platform type")

	infra, err := oc.AdminConfigClient().ConfigV1().Infrastructures().Get(context.Background(), "cluster", metav1.GetOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())

	if infra.Status.Platform != configv1.OpenStackPlatformType {
		e2eskipper.Skipf("No OpenStack platform detected")
	}
}

// client generates a Gophercloud client for the given service. Available
// services are listed above as constants.
func client(service string) (*gophercloud.ServiceClient, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: "moc",
	}
	return clientconfig.NewServiceClient(service, opts)
}
