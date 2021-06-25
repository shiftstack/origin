package openstack

import (
	"fmt"
	"strings"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	exutil "github.com/openshift/origin/test/extended/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type FieldGetterFunc func(obj map[string]interface{}, fields ...string) (interface{}, bool, error)

func baremetalClient(dc dynamic.Interface) dynamic.ResourceInterface {
	baremetalClient := dc.Resource(schema.GroupVersionResource{Group: "metal3.io", Resource: "baremetalhosts", Version: "v1alpha1"})
	return baremetalClient.Namespace("openshift-machine-api")
}

func expectField(host unstructured.Unstructured, nestedField string, fieldGetter FieldGetterFunc) o.Assertion {
	fields := strings.Split(nestedField, ".")

	value, found, err := fieldGetter(host.Object, fields...)
	o.Expect(err).NotTo(o.HaveOccurred())
	o.Expect(found).To(o.BeTrue(), fmt.Sprintf("baremetalhost field `%s` not found", nestedField))
	return o.Expect(value)
}

func expectStringField(host unstructured.Unstructured, nestedField string) o.Assertion {
	return expectField(host, nestedField, func(obj map[string]interface{}, fields ...string) (interface{}, bool, error) {
		return unstructured.NestedString(host.Object, fields...)
	})
}

func expectBoolField(host unstructured.Unstructured, nestedField string) o.Assertion {
	return expectField(host, nestedField, func(obj map[string]interface{}, fields ...string) (interface{}, bool, error) {
		return unstructured.NestedBool(host.Object, fields...)
	})
}

var _ = g.Describe("[sig-installer][Feature:openstack] OpenStack platform should", func() {
	defer g.GinkgoRecover()

	oc := exutil.NewCLI("openstack")

	g.It("have Control plane nodes in a Server group", func() {
		skipIfNotOpenStack(oc)

		// e2e.LoadClientset()
		// o.Expect(configFile).To(o.Equal("ahah"))

		// c, err := e2e.LoadClientset()
		// o.Expect(err).ToNot(o.HaveOccurred())

		// metal3, err := c.AppsV1().Deployments("openshift-machine-api").Get(context.Background(), "metal3", metav1.GetOptions{})
		// o.Expect(err).NotTo(o.HaveOccurred())
		// o.Expect(metal3.Status.AvailableReplicas).To(o.BeEquivalentTo(1))

		// o.Expect(metal3.Annotations).Should(o.HaveKey("baremetal.openshift.io/owned"))
		// o.Expect(metal3.Labels).Should(o.HaveKeyWithValue("baremetal.openshift.io/cluster-baremetal-operator", "metal3-state"))
	})

	// g.It("have baremetalhost resources", func() {
	// 	skipIfNotOpenStack(oc)

	// 	dc := oc.AdminDynamicClient()
	// 	bmc := baremetalClient(dc)

	// 	hosts, err := bmc.List(context.Background(), v1.ListOptions{})
	// 	o.Expect(err).NotTo(o.HaveOccurred())
	// 	o.Expect(hosts.Items).ToNot(o.BeEmpty())

	// 	for _, h := range hosts.Items {
	// 		expectStringField(h, "status.operationalStatus").To(o.BeEquivalentTo("OK"))
	// 		expectStringField(h, "status.provisioning.state").To(o.Or(o.BeEquivalentTo("provisioned"), o.BeEquivalentTo("externally provisioned")))
	// 		expectBoolField(h, "spec.online").To(o.BeTrue())
	// 	}
	// })
})
