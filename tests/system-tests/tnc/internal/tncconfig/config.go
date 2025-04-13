package tncconfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kelseyhightower/envconfig"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/systemtestsconfig"
	"gopkg.in/yaml.v2"
)

const (
	// PathToDefaultTncParamsFile path to config file with default tnc parameters.
	PathToDefaultTncParamsFile = "./default.yaml"
)

// TncConfig type keeps tnc configuration.
type TncConfig struct {
	*systemtestsconfig.SystemTestsConfig
	Namespace                   string `yaml:"tnc_default_ns" envconfig:"ECO_SYSTNC_NS"`
	OdfMCPName                  string `yaml:"odf_mcp" envconfig:"ECO_SYSTEM_TNC_ODF_MCP"`
	TncPpMCPName                string `yaml:"tnc_pp_mcp" envconfig:"ECO_SYSTEM_TNC_PP_MCP"`
	TncCpMCPName                string `yaml:"tnc_cp_mcp" envconfig:"ECO_SYSTEM_TNC_CP_MCP"`
	Host                        string `yaml:"host" envconfig:"ECO_SYSTEM_TNC_HOST"`
	User                        string `yaml:"user" envconfig:"ECO_SYSTEM_TNC_USER"`
	Pass                        string `yaml:"pass" envconfig:"ECO_SYSTEM_TNC_PASS"`
	MirrorRegistryUser          string `yaml:"mirror_registry_user" envconfig:"ECO_SYSTEM_TNC_MIRROR_REGISTRY_USER"`
	MirrorRegistryPass          string `yaml:"mirror_registry_pass" envconfig:"ECO_SYSTEM_TNC_MIRROR_REGISTRY_PASSWORD"`
	CombinedPullSecretFile      string `yaml:"combined_pull_secret" envconfig:"ECO_SYSTEM_TNC_COMBINED_PULL_SECRET"`
	PrivateKey                  string `yaml:"private_key" envconfig:"ECO_SYSTEM_TNC_PRIVATE_KEY"`
	RegistryRepository          string `yaml:"registry_repository" envconfig:"ECO_SYSTEM_TNC_REGISTRY_REPOSITORY"`
	CPUIsolated                 string `yaml:"cpu_isolated" envconfig:"ECO_SYSTEM_TNC_CPU_ISOLATED"`
	CPUReserved                 string `yaml:"cpu_reserved" envconfig:"ECO_SYSTEM_TNC_CPU_RESERVED"`
	OcpUpgradeUpstreamURL       string `yaml:"ocpUpgradeUpstreamUrl" envconfig:"ECO_CNF_TNC_OCP_UPGRADE_UPSTREAM_URL"`
	OdfLabel                    string
	TncPpLabel                  string
	TncCpLabel                  string
	ControlPlaneLabelListOption metav1.ListOptions
	OdfLabelListOption          metav1.ListOptions
	TncPpLabelListOption        metav1.ListOptions
	TncCpLabelListOption        metav1.ListOptions
	OdfLabelMap                 map[string]string
	TncPpLabelMap               map[string]string
	TncCpLabelMap               map[string]string
}

// NewTncConfig returns instance of TncConfig config type.
func NewTncConfig() *TncConfig {
	log.Print("Creating new TncConfig struct")

	var tncConf TncConfig
	tncConf.SystemTestsConfig = systemtestsconfig.NewSystemTestsConfig()

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	confFile := filepath.Join(baseDir, PathToDefaultTncParamsFile)
	err := readFile(&tncConf, confFile)

	if err != nil {
		log.Printf("Error to read config file %s", confFile)

		return nil
	}

	err = readEnv(&tncConf)

	if err != nil {
		log.Print("Error to read environment variables")

		return nil
	}

	return &tncConf
}

func readFile(tncConfig *TncConfig, cfgFile string) error {
	openedCfgFile, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	defer func() {
		_ = openedCfgFile.Close()
	}()

	decoder := yaml.NewDecoder(openedCfgFile)
	err = decoder.Decode(&tncConfig)

	if err != nil {
		return err
	}

	return nil
}

func readEnv(tncConfig *TncConfig) error {
	err := envconfig.Process("", tncConfig)
	if err != nil {
		return err
	}

	tncConfig.OdfLabel = fmt.Sprintf("%s/%s", tncConfig.KubernetesRolePrefix, tncConfig.OdfMCPName)
	tncConfig.TncPpLabel = fmt.Sprintf("%s/%s", tncConfig.KubernetesRolePrefix, tncConfig.TncPpMCPName)
	tncConfig.TncCpLabel = fmt.Sprintf("%s/%s", tncConfig.KubernetesRolePrefix, tncConfig.TncCpMCPName)
	tncConfig.ControlPlaneLabelListOption = metav1.ListOptions{LabelSelector: tncConfig.ControlPlaneLabel}
	tncConfig.OdfLabelListOption = metav1.ListOptions{LabelSelector: tncConfig.OdfLabel}
	tncConfig.TncPpLabelListOption = metav1.ListOptions{LabelSelector: tncConfig.TncPpLabel}
	tncConfig.TncCpLabelListOption = metav1.ListOptions{LabelSelector: tncConfig.TncCpLabel}
	tncConfig.OdfLabelMap = map[string]string{tncConfig.OdfLabel: ""}
	tncConfig.TncPpLabelMap = map[string]string{tncConfig.TncPpLabel: ""}
	tncConfig.TncCpLabelMap = map[string]string{tncConfig.TncCpLabel: ""}

	return nil
}
