package openstack

import (
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	exutil "github.com/openshift/origin/test/extended/util"
)

var _ = g.Describe("[sig-installer][Feature:openstack] OpenStack platform should", func() {
	defer g.GinkgoRecover()

	oc := exutil.NewCLI("openstack")

	g.It("have Control plane nodes in a Server group", func() {
		skipIfNotOpenStack(oc)

		computeClient, err := client(serviceCompute)
		o.Expect(err).NotTo(o.HaveOccurred())

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
