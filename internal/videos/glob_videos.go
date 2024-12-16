package videos

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	str "strings"

	"github.com/joomcode/errorx"

	. "adc/internal/javInfo"
	. "adc/internal/logger"
)

var FilesList []*JavInfo
var minimumFileSize int64 = 300 * 1024 * 1024 // 300MB

type avMatchRule struct {
	Rule *regexp.Regexp
	Type string
}

var (
	javMatchReg     = regexp.MustCompile("([a-zA-Z]{2,10})[-_](\\d{2,5})")
	fc2MatchReg     = regexp.MustCompile("FC2+(-PPV|)-[0-9]+")
	heydouMatchReg  = regexp.MustCompile("(HEYDOUGA)[-_]*(\\d{4})[-_]0?(\\d{3,5})")
	getchuMatchReg  = regexp.MustCompile("GETCHU[-_]*(\\d+)")
	gyuttoMatchReg  = regexp.MustCompile("GYUTTO-(\\d+)")
	luxuMatchReg    = regexp.MustCompile("259LUXU-(\\d+)")
	hasChnSubReg    = regexp.MustCompile("(-c|-C|-UC|-U)$")
	isUncensoredReg = regexp.MustCompile("(\\d{6}[-_]\\d{2,3})")
)

var avMatchRegList []*avMatchRule

func init() {
	avMatchRegList = []*avMatchRule{
		{javMatchReg, "jav"},
		{fc2MatchReg, "fc2"},
		{heydouMatchReg, "heydou"},
		{getchuMatchReg, "getchu"},
		{gyuttoMatchReg, "gyutto"},
		{luxuMatchReg, "luxu"},
	}
}

func JavWalk(srcDir string) error {
	// clear FilesList
	FilesList = make([]*JavInfo, 0)
	err := filepath.Walk(srcDir, javGlobFunc)
	if err != nil {
		return errorx.Decorate(err, "Glob Jav error!")
	}
	return nil
}

func javGlobFunc(p string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// ignore directory
	if info.IsDir() {
		return nil
	}

	// check file extension
	fileExt := path.Ext(p)
	if !isValidExt(fileExt) {
		return nil
	}

	// check file size
	fileSize := info.Size()
	if fileSize <= minimumFileSize {
		return nil
	}

	// check file name
	fileName := info.Name()
	fileNameWithoutExt := str.TrimSuffix(fileName, fileExt)
	javType, ok := isValidJavNumber(fileNameWithoutExt)
	if !ok {
		return nil
	}

	// append to fileList
	javInfo, err := newJavInfo(fileNameWithoutExt, javType, p)
	if err != nil {
		Logger.Infoln(err)
		return nil
	}

	FilesList = append(FilesList, javInfo)
	Logger.WithField("number", javInfo.Number).Debugf("Path: %s", filepath.ToSlash(p))

	return nil
}

func isValidExt(ext string) bool {
	for _, s := range JavExt {
		if ext == s {
			return true
		}
	}
	return false
}

func isValidJavNumber(fileName string) (name string, ok bool) {
	for _, reg := range avMatchRegList {
		if reg.Rule.MatchString(fileName) {
			return reg.Type, true
		}
	}
	return "", false
}

func newJavInfo(fileNameWithoutExt, javType, javFilePath string) (*JavInfo, error) {
	info := &JavInfo{}

	var rule *regexp.Regexp
	for _, reg := range avMatchRegList {
		if reg.Type == javType {
			rule = reg.Rule
		}
	}
	javNumbers := rule.FindAllString(fileNameWithoutExt, -1)
	if javNumbers == nil {
		return info, fmt.Errorf("%s is not a valid jav file", fileNameWithoutExt)
	}

	info.Number = str.ToUpper(javNumbers[0])
	info.SrcFilePath = filepath.ToSlash(javFilePath)
	info.Type = javType
	info.HasChnSub = hasChnSubReg.MatchString(fileNameWithoutExt)
	info.IsUncensored = isUncensoredReg.MatchString(javFilePath)

	return info, nil
}
