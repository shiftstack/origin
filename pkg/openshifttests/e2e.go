package openshifttests

import (
	"regexp"
	"strings"
	"time"

	"github.com/openshift/origin/pkg/test/ginkgo"
	exutil "github.com/openshift/origin/test/extended/util"

	// these register framework.NewFrameworkExtensions responsible for
	// executing post-test actions, here debug and metrics gathering
	// see https://github.com/kubernetes/kubernetes/blob/v1.26.0/test/e2e/framework/framework.go#L175-L181
	// for more details
	_ "k8s.io/kubernetes/test/e2e/framework/debug/init"
	_ "k8s.io/kubernetes/test/e2e/framework/metrics/init"
)

func IsDisabled(name string) bool {
	if strings.Contains(name, "[Disabled") {
		return true
	}

	return shouldSkipUntil(name)
}

// shouldSkipUntil allows a test to be skipped with a time limit.
// the test should be annotated with the 'SkippedUntil' tag, as shown below.
//
//	[SkippedUntil:05092022:blocker-bz/123456]
//
// - the specified date should conform to the 'MMDDYYYY' format.
// - a valid blocker BZ must be specified
// if the specified date in the tag has not passed yet, the test
// will be skipped by the runner.
func shouldSkipUntil(name string) bool {
	re, err := regexp.Compile(`\[SkippedUntil:(\d{8}):blocker-bz\/([a-zA-Z0-9]+)\]`)
	if err != nil {
		// it should only happen with a programmer error and unit
		// test will prevent that
		return false
	}
	matches := re.FindStringSubmatch(name)
	if len(matches) != 3 {
		return false
	}

	skipUntil, err := time.Parse("01022006", matches[1])
	if err != nil {
		return false
	}

	if skipUntil.After(time.Now()) {
		return true
	}
	return false
}

type TestSuite struct {
	ginkgo.TestSuite

	PreSuite  func(opt *RunOptions) error
	PostSuite func(opt *RunOptions)

	PreTest func() error
}

type TestSuites []TestSuite

func (s TestSuites) TestSuites() []*ginkgo.TestSuite {
	copied := make([]*ginkgo.TestSuite, 0, len(s))
	for i := range s {
		copied = append(copied, &s[i].TestSuite)
	}
	return copied
}

// IsStandardEarlyTest returns true if a test is considered part of the normal
// pre or post condition tests.
func IsStandardEarlyTest(name string) bool {
	if !strings.Contains(name, "[Early]") {
		return false
	}
	return strings.Contains(name, "[Suite:openshift/conformance/parallel")
}

// IsStandardEarlyOrLateTest returns true if a test is considered part of the normal
// pre or post condition tests.
func IsStandardEarlyOrLateTest(name string) bool {
	if !strings.Contains(name, "[Early]") && !strings.Contains(name, "[Late]") {
		return false
	}
	return strings.Contains(name, "[Suite:openshift/conformance/parallel")
}

// suiteWithInitializedProviderPreSuite loads the provider info, but does not
// exclude any tests specific to that provider.
func SuiteWithInitializedProviderPreSuite(opt *RunOptions) error {
	config, err := decodeProvider(opt.Provider, opt.DryRun, true, nil)
	if err != nil {
		return err
	}
	opt.config = config

	opt.Provider = config.ToJSONString()
	return nil
}

// SuiteWithProviderPreSuite ensures that the suite filters out tests from providers
// that aren't relevant (see exutilcluster.ClusterConfig.MatchFn) by loading the
// provider info from the cluster or flags.
func SuiteWithProviderPreSuite(opt *RunOptions) error {
	if err := SuiteWithInitializedProviderPreSuite(opt); err != nil {
		return err
	}
	opt.MatchFn = opt.config.MatchFn()
	return nil
}

// suiteWithNoProviderPreSuite blocks out provider settings from being passed to
// child tests. Used with suites that should not have cloud specific behavior.
func SuiteWithNoProviderPreSuite(opt *RunOptions) error {
	opt.Provider = `none`
	return SuiteWithProviderPreSuite(opt)
}

// suiteWithKubeTestInitialization invokes the Kube suite in order to populate
// data from the environment for the CSI suite. Other suites should use
// SuiteWithProviderPreSuite.
func SuiteWithKubeTestInitializationPreSuite(opt *RunOptions) error {
	if err := SuiteWithProviderPreSuite(opt); err != nil {
		return err
	}
	return initializeTestFramework(exutil.TestContext, opt.config, opt.DryRun)
}
