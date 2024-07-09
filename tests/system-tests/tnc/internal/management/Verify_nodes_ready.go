package management

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openshift-kni/eco-goinfra/pkg/clusteroperator"
	"github.com/openshift-kni/eco-goinfra/pkg/clusterversion"
	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/nodes"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/tnc/internal/tncinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/tnc/internal/tncparams"
)

// SetSystemReservedMemoryForMasterNodes assert system reserved memory for masters succeeded.
func SetSystemReservedMemoryForMasterNodes(ctx SpecContext) {
	glog.V(tncparams.TncLogLevel).Infof("Verify system reserved memory config for masters succeeded")

	systemReservedBuilder := mco.NewKubeletConfigBuilder(APIClient, tncparams.KubeletConfigName).
		WithMCPoolSelector("pools.operator.machineconfiguration.openshift.io/master", "").
		WithSystemReserved(tncparams.SystemReservedCPU, tncparams.SystemReservedMemory)

	if !systemReservedBuilder.Exists() {
		glog.V(tncparams.TncLogLevel).Infof("Create system-reserved configuration")

		systemReserved, err := systemReservedBuilder.Create()
		Expect(err).ToNot(HaveOccurred(), "Failed to create %s kubeletConfig objects "+
			"with system-reserved definition", tncparams.KubeletConfigName)

		_, err = nodes.WaitForAllNodesToReboot(
			APIClient,
			40*time.Minute,
			TncConfig.ControlPlaneLabelListOption)
		Expect(err).ToNot(HaveOccurred(), "Nodes failed to reboot after applying %s config; %s",
			tncparams.KubeletConfigName, err)

		Expect(systemReserved.Exists()).To(Equal(true),
			"Failed to setup master system reserved memory, %s kubeletConfig not found; %s",
			tncparams.KubeletConfigName, err)

		glog.V(tncparams.TncLogLevel).Infof("Checking all master nodes are Ready")

		var isReady bool
		isReady, err = nodes.WaitForAllNodesAreReady(
			APIClient,
			30*time.Second,
			TncConfig.ControlPlaneLabelListOption)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Some nodes are not ready due to: %v", err))
		Expect(isReady).To(Equal(true),
			fmt.Sprintf("Failed master nodes status, not all Master node are Ready; %v", isReady))

		glog.V(tncparams.TncLogLevel).Infof("Checking that the clusterversion is available")

		_, err = clusterversion.Pull(APIClient)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Error accessing csv: %v", err))

		var coBuilder []*clusteroperator.Builder
		coBuilder, err = clusteroperator.List(APIClient)
		Expect(err).To(BeNil(), fmt.Sprintf("ClusterOperator List not found: %v", err))
		Expect(len(coBuilder)).ToNot(Equal(0), "Empty clusterOperators list received")

		_, err = clusteroperator.WaitForAllClusteroperatorsAvailable(APIClient, 60*time.Second)
		Expect(err).ToNot(HaveOccurred(),
			fmt.Sprintf("Error waiting for all available clusteroperators: %v", err))
	}
}

// VerifyPostDeploymentConfig container that contains tests for basic post-deployment config verification.
func VerifyPostDeploymentConfig() {
	Describe(
		"Post-deployment config validation",
		Label(tncparams.LabelTncDeployment), func() {
			BeforeAll(func() {
				By(fmt.Sprintf("Asserting %s folder exists", tncparams.ConfigurationFolderName))

				homeDir, err := os.UserHomeDir()
				Expect(err).To(BeNil(), fmt.Sprint(err))

				tncConfigsFolder := filepath.Join(homeDir, tncparams.ConfigurationFolderName)

				glog.V(tncparams.TncLogLevel).Infof("tncConfigsFolder: %s", tncConfigsFolder)

				if err := os.Mkdir(tncConfigsFolder, 0755); os.IsExist(err) {
					glog.V(tncparams.TncLogLevel).Infof("%s folder already exists", tncConfigsFolder)
				}
			})

			It("Verifies system reserved memory for masters succeeded",
				Label("system-reserved"), reportxml.ID("60045"), SetSystemReservedMemoryForMasterNodes)
		})
}
