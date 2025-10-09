package version

var (
	version   = "1.0.0"
	builddate = "YYYY-mm-dd"
)

func GetVersion() string {
	return version
}

func GetBuildDate() string {
	return builddate
}
