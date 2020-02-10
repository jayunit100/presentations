/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package network

import (
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"

	"fmt"

	"github.com/onsi/ginkgo"
)

/*
The following Network Policy tests verify that policy object definitions
are correctly enforced by a networking plugin. It accomplishes this by launching
a simple netcat server, and two clients with different
attributes. Each test case creates a network policy which should only allow
connections from one of the clients. The test then asserts that the clients
failed or successfully connected as expected.
*/

// The first test can be described using the following truth table entry
var podServerLabelSelector string
var nsANamePtr *string
var nsBNamePtr *string
var nsCNamePtr *string

type TruthTableEntry struct {
	OldDescription string
	Description    string
	Policy         *networkingv1.NetworkPolicy
	Whitelist      map[string]bool
}

// The following set of selectors and rules can be used to make selectors for peers
var emptyLabelSelector = metav1.LabelSelector{}

var _ = SIGDescribe("NetworkPolicy2 [LinuxOnly]", func() {
	var podServer *v1.Pod
	var service *v1.Service
	f := framework.NewDefaultFramework("network-policy")

	// Create the namespaces used by all tests.
	// var nsA = f.Namespace

	nsBName := f.BaseName + "-b"
	nsBNamePtr = &nsBName

	_, err := f.CreateNamespace(f.BaseName, map[string]string{
		"ns-name": *nsBNamePtr,
	})
	if err != nil {
		panic("failed creating ns")
	}
	nsCName := f.BaseName + "-c"
	nsCNamePtr = &nsCName

	_, err = f.CreateNamespace(f.BaseName, map[string]string{
		"ns-name": *nsCNamePtr,
	})
	if err != nil {
		panic("failed creating ns")
	}
	// Create the pods used by all tests

	podServer, service = createServerPodAndService(f, f.Namespace, "server", []int{80, 81})

	//  now initialize
	podServerLabelSelector = podServer.ObjectMeta.Labels["pod-name"]

	ginkgo.BeforeEach(func() {
		// Windows does not support network policies.
		e2eskipper.SkipIfNodeOSDistroIs("windows")
	})

	ginkgo.Context("NetworkPolicy between server and client", func() {

		p80 := 80
		p81 := 81

		ginkgo.BeforeSuite(func() {
			// create nsA, nsB
			// create pods a, b : in each framework, nsA, nsb
			// verify that all pods can talk to all other pods.
		})

		ginkgo.AfterSuite(func() {
			// delete nsA, nsB as they were created outside of the framework context.
		})
		var A_V = "ASDF"
		ginkgo.It("should support a 'default-deny' policy [Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			builder := &NetworkPolicySpecBuilder{}
			builder.SetName("deny-all").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(
				nil, nil, nil, nil, nil, nil, nil, nil)
			// verify all connectivity in whitelist

			
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy to allow traffic from pods within server namespace based on PodSelector [Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			// Create AllowInnerNamespaceSelector.Policy
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, &map[string]string{"pod-name": "client-a"}, nil, nil, nil)
			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy to allow traffic only from a different namespace, based on NamespaceSelector [Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, nil, &map[string]string{"ns-name": *nsBNamePtr}, nil, nil)
			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy based on PodSelector with MatchExpressions[Feature:NetworkPolicy]", func() {

			// Verify all connectivity
			matchpodNameIsClientA := &[]metav1.LabelSelectorRequirement{{
				Key:      "pod-name",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"client-a"},
			}}

			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{}).SetTypeIngress()
			builder = builder.AddIngress(nil, &p80, nil, nil, nil, nil, matchpodNameIsClientA, nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})
		notNsC := &[]metav1.LabelSelectorRequirement{{
			Key:      "ns-name",
			Operator: metav1.LabelSelectorOpNotIn,
			// see the -c above, i think thats what we mean here.
			Values: []string{*nsCNamePtr},
		}}
		ginkgo.It("should enforce policy based on NamespaceSelector with MatchExpressions[Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			notNsC := &[]metav1.LabelSelectorRequirement{{
				Key:      "ns-name",
				Operator: metav1.LabelSelectorOpNotIn,
				// see the -c above, i think thats what we mean here.
				Values: []string{*nsCNamePtr},
			}}

			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{}).SetTypeIngress()
			builder = builder.AddIngress(nil, &p80, nil, nil, nil, nil, notNsC, nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy based on NamespaceSelector with MatchExpressions[Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{}).SetTypeIngress()
			builder = builder.AddIngress(nil, &p80, nil, nil, nil, nil, notNsC, nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy based on PodSelector or NamespaceSelector [Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress()
			builder.AddIngress(nil, &p80, nil, nil,
				&map[string]string{"pod-name": "client-b"}, nil, nil, nil)
			builder.AddIngress(nil, &p80, nil, nil,
				nil, &map[string]string{"ns-name": *nsBNamePtr}, nil, nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]", func() {
			// Verify all connectivity
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-pod-selector").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress()
			builder.AddIngress(nil, &p80, nil, nil,
				&map[string]string{"pod-name": "client-b"},
				&map[string]string{"ns-name": *nsBNamePtr}, nil, nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		// TEST ASDF
		ginkgo.It("should enforce policy to allow traffic only from a pod in a different namespace based on PodSelector and NamespaceSelector [Feature:NetworkPolicy]", func() {

			// todo 2:37 pm
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ns-b-client-a-via-namespace-pod-selector").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress()
			builder.AddIngress(nil, &p80, nil, nil,
				&map[string]string{"pod-name": "client-a"},
				&map[string]string{"ns-name": *nsBNamePtr},
				nil,
				nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should enforce policy based on Ports [Feature:NetworkPolicy]", func() {
			// todo 2:37 pm
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ingress-on-port-81").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress()
			builder.AddIngress(nil, &p81, nil, nil,
				nil,
				nil,
				nil,
				nil)

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		// TODO 12333
		ginkgo.It("should enforce multiple, stacked policies with overlapping podSelectors [Feature:NetworkPolicy]", func() {
			// todo 2:37 pm
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ingress-on-port-81").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress()
			builder.AddIngress(nil, &p80, nil, nil,
				nil,
				nil,
				nil,
				nil)
			// CREATE the first policy...
			builder2 := &NetworkPolicySpecBuilder{}
			builder2 = builder2.SetName("allow-ingress-on-port-81").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder2.SetTypeIngress()
			builder2.AddIngress(nil, &p81, nil, nil,
				nil,
				nil,
				nil,
				nil)
			// CREATE the second policy...

			// union the port 80 + 81 policies
			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		ginkgo.It("should support allow-all policy [Feature:NetworkPolicy]", func() {

			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-all").SetPodSelector(map[string]string{})
			builder.SetTypeIngress()
			builder.AddIngress(nil, nil, nil, nil,
				nil,
				nil,
				nil,
				nil) // make sure AddIngress creates an empty array !

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})
		s80 := "serve-80"

		ginkgo.It("should allow ingress access on one named port [Feature:NetworkPolicy]", func() {
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-all").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder = builder.SetTypeIngress()
			builder.AddIngress(nil, nil, &s80, nil,
				nil,
				nil,
				nil,
				nil) // make sure named ports work

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		// i think we should qualify this as 'exactly'  vs 'one or more'
		ginkgo.It("should allow ingress access from namespace on one named port [Feature:NetworkPolicy]", func() {

			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-client-a-via-named-port-egress-rule").SetPodSelector(map[string]string{"pod-name": "client-a"})
			builder = builder.SetTypeIngress()
			builder.AddIngress(nil, nil, &s80, nil,
				nil,
				&map[string]string{"ns-name": *nsBNamePtr},
				nil,
				nil) // make sure named ports work

			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity
		})

		// we dont specify the namedport/dns dependency explicitly here, i think we should.
		ginkgo.It("should allow egress access on one named port [Feature:NetworkPolicy]", func() {

			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-all").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder = builder.SetTypeIngress()

			// TODO add support for a second Egress port somehow?
			builder.AddEgress(nil, nil, &s80, nil,
				nil,
				nil,
				nil,
				nil) // make sure named ports work

			// this appends the 53 udp port to every single egress rule.
			builder.WithEgressDNS()

			// make sure named ports work
			// verify all connectivity in whitelist
			// complement the whitelist and verify disconnectivity

		})

		ginkgo.It("should enforce updated policy [Feature:NetworkPolicy]", func() {

			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ingress").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, &map[string]string{"pod-name": "client-a"}, nil, nil, nil)
			_, err = f.ClientSet.NetworkingV1().NetworkPolicies(f.Namespace.Name).Create(builder2.Get())

			// test 

			builder2 := &NetworkPolicySpecBuilder{}
			builder2 = builder2.SetName("allow-ingress").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder2.SetTypeIngress().AddIngress(nil, &p81, nil, nil, &map[string]string{"pod-name": "client-b"}, nil, nil, nil)
			_, err = f.ClientSet.NetworkingV1().NetworkPolicies(f.Namespace.Name).Update(builder2.Get())

			// test updates
		})

		ginkgo.It("should allow ingress access from updated namespace [Feature:NetworkPolicy]", func() {
			// rewrite this test
		})

		ginkgo.It("should allow ingress access from updated pod [Feature:NetworkPolicy]", func() {
			// rewrite this test
		})

		// we should say "matching BOTH podSelector AND Namespaceselector"
		ginkgo.It("should enforce egress policy allowing traffic to a server in a different namespace based PodSelector and NamespaceSelector [Feature:NetworkPolicy]", func() {
			// rewrite this test

		})

		ginkgo.It("should enforce multiple ingress policies with ingress allow-all policy taking precedence [Feature:NetworkPolicy]", func() {			
			
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ingress").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, &map[string]string{"pod-name": "client-b"}, nil, nil, nil)


			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-all").SetPodSelector(map[string]string{})
			builder.SetTypeIngress()
			builder.AddIngress(nil, nil, nil, nil,
				nil,
				nil,
				nil,
				nil) // make sure AddIngress creates an empty array !

			// verify all connectivity 

		})

		ginkgo.It("should enforce multiple egress policies with egress allow-all policy taking precedence [Feature:NetworkPolicy]", func() {
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ingress").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeEgress()().AddIngress(nil, &p80, nil, nil, &map[string]string{"pod-name": "client-b"}, nil, nil, nil)


			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-all").SetPodSelector(map[string]string{})
			builder.SetTypeEgress()
			builder.AddEgress(nil, nil, nil, nil,
				nil,
				nil,
				nil,
				nil) // make sure AddIngress creates an empty array !

			// verify all connectivity 

		})

		ginkgo.It("should stop enforcing policies after they are deleted [Feature:NetworkPolicy]", func() {
			builder := &NetworkPolicySpecBuilder{}
			builder.SetName("deny-all").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(
				nil, nil, nil, nil, nil, nil, nil, nil)

			policyDenyAll, err := f.ClientSet.NetworkingV1().NetworkPolicies(f.Namespace.Name).Create(policyDenyAll)


			ginkgo.By("Creating a network policy for the server which allows traffic only from client-a.")
			
			builder := &NetworkPolicySpecBuilder{}
			builder = builder.SetName("allow-ingress").SetPodSelector(map[string]string{"pod-name": podServerLabelSelector})
			builder.SetTypeIngress().AddIngress(nil, &p80, nil, nil, &map[string]string{"pod-name": "client-a"}, nil, nil, nil)
			_, err = f.ClientSet.NetworkingV1().NetworkPolicies(f.Namespace.Name).Create(builder2.Get())

			// delete allow policy, confirm traffic is denies

			// delete deny-all, confirm now all traffic is back online.
		})

		// we should update this description to clarify that we only test one pod with a /32, and btw,
		// is this something we can run on ipv6 clusters ?
		ginkgo.It("should allow egress access to server in CIDR block [Feature:NetworkPolicy]", func() {
			// do this later
		})

		// i think we also need to make sure there is a test for enforcing multiple policies with a NamespaceSelector.

		// i think we mean "should enforce multiple separate policies against the same PodSelector"
		ginkgo.It("should enforce policies to check ingress and egress policies can be controlled independently based on PodSelector [Feature:NetworkPolicy]", func() {
			// do this later
		})
	}
})
