package files

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"go/src/github.com/RajeevRanjan27/golangclidemo/common"

	"github.com/gookit/color"
	"github.com/spf13/viper"
)

type File struct {
	Path            string
	BytesSize       int64
	PrettyBytesSize string
}

func ReadDirRecursively(dirPath string) ([]File, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var results []File
	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())

		if file.IsDir() {
			subFiles, err := ReadDirRecursively(filePath)
			if err != nil {
				return nil, err
			}
			results = append(results, subFiles...)
		} else {
			if file.Size() >= int64(viper.GetInt("minfilesize")*1000000) {
				foundFile := File{
					Path:            filePath,
					BytesSize:       file.Size(),
					PrettyBytesSize: common.PrettyBytes(file.Size()),
				}
				results = append(results, foundFile)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].BytesSize > results[j].BytesSize
	})

	return results, nil
}

func PrintResults(files []File) {
	fmt.Println()
	common.PrintColor("forestgreen", "background", fmt.Sprintf("Largest files in: %s", viper.GetString("path")))
	fmt.Println("----------------------------------------------------------")
	fmt.Println()

	spacing := make(map[string]int)
	highWaterMark := 0

	for _, file := range files {
		if len(file.Path) > highWaterMark {
			highWaterMark = len(file.Path)
		}

		spacing[file.Path] = len(file.Path)
	}

	for _, file := range files {
		padding := strconv.Itoa(highWaterMark + 2)

		if file.BytesSize >= int64(viper.GetInt("highlight")*1000000) {
			color.HEXStyle("000", common.AllHex["yellow1"]).Printf("%-"+padding+"s %10s\n", file.Path, file.PrettyBytesSize)
		} else {
			color.HEXStyle(common.AllHex["steelblue2"]).Printf("%-"+padding+"s %10s\n", file.Path, file.PrettyBytesSize)
		}
	}
}
