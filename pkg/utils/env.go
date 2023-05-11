package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	BundleData      *BundleDataType
	BundleDir       string
	BundleDirs      *BundleDirsType
	CatalogIndexes  map[string]string
	ClientsBaseURL  = "https://mirror.openshift.com/pub/openshift-v4/clients/ocp"
	RhcosBaseURL    = "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos"
	RhcosPreBaseURL = "https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/pre-release"
	Cwd, Err        = os.Getwd() // nolint:revive
	Logger          = logrus.New()
	OCLocalCmdPath  = fmt.Sprintf("oc-%s-%s", runtime.GOOS, runtime.GOARCH)
)
