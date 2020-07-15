package platform

import (
	"os"
	"time"

	"github.com/jenkins-x/bdd-jx/test/helpers"
	"github.com/jenkins-x/bdd-jx/test/utils"
	"github.com/jenkins-x/bdd-jx/test/utils/runner"
	cmd "github.com/jenkins-x/jx/v2/pkg/cmd/clients"
	"github.com/jenkins-x/jx/v2/pkg/jenkins"
	"github.com/jenkins-x/jx/v2/pkg/kube"
	"github.com/jenkins-x/jx/v2/pkg/kube/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
)

type testCaseUpgradePlatform struct {
	*runner.JxRunner
	version   string
	client    kubernetes.Interface
	namespace string
}

func newTestCaseUpgradePlatfrom(cwd string, version string, factory cmd.Factory) (*testCaseUpgradePlatform, error) {
	client, ns, err := factory.CreateKubeClient()
	if err != nil {
		return nil, err
	}

	return &testCaseUpgradePlatform{
		JxRunner:  runner.New(cwd, nil, 0),
		version:   version,
		client:    client,
		namespace: ns,
	}, nil
}

func (t *testCaseUpgradePlatform) Upgrade(args ...string) {
	allargs := []string{"upgrade", "platform",
		"--version=" + t.version, "-b"}
	allargs = append(allargs, args...)
	t.Run(allargs...)
}

func (t *testCaseUpgradePlatform) CheckJenkins() {
	url, err := services.FindServiceURL(t.client, t.namespace, kube.ServiceJenkins)
	Expect(err).NotTo(HaveOccurred())
	utils.LogInfof("Checking health of Jenkins service: %q\n", url)
	err = jenkins.CheckHealth(url, time.Minute*5)
	Expect(err).NotTo(HaveOccurred())
}

var _ = Describe("upgrade platform", func() {
	var test *testCaseUpgradePlatform
	skipJenkinsCheck := false

	BeforeEach(func() {
		version := os.Getenv("PLATFORM_VERSION")
		_, skipJenkinsCheck = os.LookupEnv("SKIP_JENKINS_CHECK")

		utils.LogInfof("Using platform version: %q\n", version)
		var err error
		test, err = newTestCaseUpgradePlatfrom(helpers.WorkDir, version, cmd.NewFactory())
		Expect(err).NotTo(HaveOccurred())
		Expect(test).NotTo(BeNil())
	})

	Describe("Given valid parameters", func() {
		Context("when running upgrade platform", func() {
			It("updates the platform to the given version", func() {
				test.Upgrade()
				if !skipJenkinsCheck {
					test.CheckJenkins()
				}
			})
		})
	})

	Describe("Given valid parameters", func() {
		Context("when running upgrade platform in force mode", func() {
			It("updates always the platform to the given version", func() {
				test.Upgrade("--always-upgrade=true")
				if !skipJenkinsCheck {
					test.CheckJenkins()
				}
			})
		})
	})
})
