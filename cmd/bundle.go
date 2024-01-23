package cmd

import (
	"fmt"

	"github.com/shanedell/openshift4-mirror-go/pkg/app"
	"github.com/shanedell/openshift4-mirror-go/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	openshiftVersion            string
	pullSecret                  string
	platform                    string
	redhatOperatorIndexImage    string
	redhatMarketplaceIndexImage string
	certifiedOperatorIndexImage string
	communityOperatorIndexImage string
	catalogVersion              string
	catalogs                    []string
	skipExisting                bool
	skipRelease                 bool
	skipCatalogs                bool
	skipRhcos                   bool
)

var bundleHelp = "bundle the OpenShift content"

func addBundleSkipFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(
		&skipExisting,
		"skip-existing",
		true,
		"skip downloading content that already exists on disk",
	)
	cmd.PersistentFlags().BoolVar(
		&skipRelease,
		"skip-release",
		false,
		"skip downloading of release content",
	)
	cmd.PersistentFlags().BoolVar(
		&skipCatalogs,
		"skip-catalogs",
		false,
		"skip downloading of catalog content",
	)
	cmd.PersistentFlags().BoolVar(
		&skipRhcos,
		"skip-rhcos",
		false,
		"skip downloading of RHCOS image",
	)
}

func addIndexImageFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&redhatOperatorIndexImage,
		"redhat-operator-index-image",
		"",
		"version of image to use for redhat-operator catalogs",
	)
	cmd.PersistentFlags().StringVar(
		&redhatMarketplaceIndexImage,
		"redhat-marketplace-index-image",
		"",
		"version of image to use for redhat-marketplace catalogs",
	)
	cmd.PersistentFlags().StringVar(
		&certifiedOperatorIndexImage,
		"certified-operator-index-image",
		"",
		"version of image to use for certified-operator catalogs",
	)
	cmd.PersistentFlags().StringVar(
		&communityOperatorIndexImage,
		"community-operator-index-image",
		"",
		"version of image to use for community-operator catalogs",
	)
}

func NewBundleCommand() *cobra.Command {
	bundleCommand := &cobra.Command{
		Use:   "bundle",
		Short: bundleHelp,
		Long:  bundleHelp,
		RunE:  bundleMain,
	}

	bundleCommand.PersistentFlags().StringVarP(
		&openshiftVersion,
		"openshift-version",
		"v",
		"",
		"the OpenShift version (e.g. 4.5.11)",
	)
	err := bundleCommand.MarkPersistentFlagRequired("openshift-version")
	if err != nil {
		panic(err)
	}

	bundleCommand.PersistentFlags().StringVar(
		&pullSecret,
		"pull-secret",
		"",
		"the content of your pull secret (can be found at https://cloud.redhat.com/openshift/install/pull-secret)",
	)
	err = bundleCommand.MarkPersistentFlagRequired("pull-secret")
	if err != nil {
		panic(err)
	}

	bundleCommand.PersistentFlags().StringVar(
		&platform,
		"platform",
		"",
		"target platform for install. platforms: [aws, azure, gcp, metal, openstack, vmware]",
	)
	err = bundleCommand.MarkPersistentFlagRequired("platform")
	if err != nil {
		panic(err)
	}

	bundleCommand.PersistentFlags().StringSliceVar(
		&catalogs,
		"catalogs",
		nil,
		"the catalog(s) content to download. catalogs: [redhat-operators, certified-operators, redhat-marketplace, community-operators]. defaults to all",
	)
	bundleCommand.PersistentFlags().StringVar(
		&catalogVersion,
		"catalog-version",
		"",
		"version of images to use for catalogs",
	)

	addBundleSkipFlags(bundleCommand)
	addIndexImageFlags(bundleCommand)

	return bundleCommand
}

// check platform is supported
func checkPlatform() error {
	switch platform {
	case "aws", "azure", "gcp", "metal", "openstack", "vmware":
		return nil
	default:
		return fmt.Errorf(
			"invalid platform. Allowed platforms: [aws, azure, gcp, metal, openstack, vmware]",
		)
	}
}

// check catalog is supported
func checkCatalogs() error {
	if catalogs == nil {
		catalogs = []string{
			"redhat-operators",
			"certified-operators",
			"redhat-marketplace",
			"community-operators",
		}
		return nil
	}

	for _, catalog := range catalogs {
		switch catalog {
		case "redhat-operators", "certified-operators", "redhat-marketplace", "community-operators":
			continue
		default:
			return fmt.Errorf(
				"invalid catalog. Supported catalogs: [redhat-operators, certified-operators, redhat-marketplace, community-operators]",
			)
		}
	}

	return nil
}

func bundleMain(_ *cobra.Command, _ []string) error {
	if err := checkPlatform(); err != nil {
		return err
	}

	if err := checkCatalogs(); err != nil {
		return err
	}

	if catalogVersion == "" {
		catalogVersion = openshiftVersion
	}

	versionMinor := utils.GetVersionMinor(openshiftVersion)

	if redhatOperatorIndexImage == "" {
		redhatOperatorIndexImage = fmt.Sprintf(
			"registry.redhat.io/redhat/redhat-operator-index:v%s", versionMinor,
		)
	}

	if redhatMarketplaceIndexImage == "" {
		redhatMarketplaceIndexImage = fmt.Sprintf(
			"registry.redhat.io/redhat/redhat-marketplace-index:v%s", versionMinor,
		)
	}

	if certifiedOperatorIndexImage == "" {
		certifiedOperatorIndexImage = fmt.Sprintf(
			"registry.redhat.io/redhat/certified-operator-index:v%s", versionMinor,
		)
	}

	if communityOperatorIndexImage == "" {
		communityOperatorIndexImage = fmt.Sprintf(
			"registry.redhat.io/redhat/community-operator-index:v%s", versionMinor,
		)
	}

	bundleData := &utils.BundleDataType{
		BundleDir:                   bundleDir,
		CatalogVersion:              catalogVersion,
		Catalogs:                    catalogs,
		Platform:                    platform,
		PreRelease:                  preRelease,
		PullSecret:                  pullSecret,
		OpenshiftVersion:            openshiftVersion,
		RedhatOperatorIndexImage:    redhatOperatorIndexImage,
		RedhatMarketplaceIndexImage: redhatMarketplaceIndexImage,
		CertifiedOperatorIndexImage: certifiedOperatorIndexImage,
		CommunityOperatorIndexImage: communityOperatorIndexImage,
		SkipExisting:                skipExisting,
		SkipRelease:                 skipRelease,
		SkipCatalogs:                skipCatalogs,
		SkipRhcos:                   skipRhcos,
		TargetRegistry:              targetRegistry,
	}

	return app.Bundle(bundleData)
}
