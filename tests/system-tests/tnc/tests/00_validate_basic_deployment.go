package tnc_system_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift-kni/eco-gotests/tests/system-tests/tnc/internal/management"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/tnc/internal/tncparams"
)

var _ = Describe(
	"TNC managemant Basic Deployment Suite",
	Ordered,
	ContinueOnFailure,
	Label(tncparams.Label), func() {

		management.VerifyPostDeploymentConfig()
	})
