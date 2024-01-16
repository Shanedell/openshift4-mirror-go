package utils

type BundleDataType struct {
	BundleDir                string
	CatalogVersion           string
	Catalogs                 []string
	OpenshiftVersion         string
	Platform                 string
	PreRelease               bool
	PullSecret               string
	RedhatOperatorIndexImage string
	SkipExisting             bool
	SkipRelease              bool
	SkipCatalogs             bool
	SkipRhcos                bool
	TargetRegistry           string
}

type BundleDirsType struct {
	Bin      string
	Release  string
	Rhcos    string
	Catalogs string
	Clients  string
}

type ContainerDataType struct {
	OpenshiftVersion string
	Runtime          string
	Image            string
}

type SaveFileToFrom struct {
	To   string
	From string
}

type PruneDataType struct {
	ImageToPrune string
	Operators    []string
	OpmVersion   string
	PruneType    string
	TargetImage  string
	FolderName   string
}
