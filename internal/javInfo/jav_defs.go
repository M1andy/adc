package javInfo

type JavInfo struct {
	// local info
	Number       string
	SrcFilePath  string
	OutDir       string
	Type         string
	HasChnSub    bool
	IsUncensored bool

	// online info
	Title        string
	ReleaseDate  string
	VideoLength  string
	Manufacturer string
	Studio       string
	Series       string
	Genre        []string
	Actresses    []string
}

var JavExt = []string{
	".mp4",
	".mkv",
	".avi",
	".rmdb",
	".rm",
	".flv",
	".mov",
}
