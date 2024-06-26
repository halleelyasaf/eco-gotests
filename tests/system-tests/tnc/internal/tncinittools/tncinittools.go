package tncinittools

import (
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-gotests/tests/internal/inittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/tnc/internal/tncconfig"
)

var (
	// APIClient provides API access to cluster.
	APIClient *clients.Settings
	// TncConfig provides access to Tnc system tests configuration parameters.
	TncConfig *tncconfig.TncConfig
)

// init loads all variables automatically when this package is imported. Once package is imported a user has full
// access to all vars within init function. It is recommended to import this package using dot import.
func init() {
	TncConfig = tncconfig.NewTncConfig()
	APIClient = inittools.APIClient
}
