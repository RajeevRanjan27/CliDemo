package dirs

import (
	"fmt"
	"io/fs"

	// "os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"go/src/github.com/RajeevRanjan27/golangclidemo/common"

	"github.com/gookit/color"
	"github.com/spf13/viper"
)

type Dir struct {
	Path            string
	Depth           int
	BytesSize       int64
	PrettyBytesSize string
}

var dirsVisited []string
var dirs []Dir

func ReadDirDepth(dirPath string) ([]Dir, error) {
	currentDir := strings.ReplaceAll(dirPath, viper.GetString("path"), "")
	currentDepth := len(strings.Split(currentDir, "/"))

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			_, err := ReadDirDepth(filepath.Join(dirPath, d.Name()))
			if err != nil {
				return err
			}
			if viper.GetInt("depth") >= currentDepth {

				if !common.SliceContains(dirsVisited, filepath.Join(dirPath, d.Name())) {
					dirSize, _ := DirSizeBytes(filepath.Join(dirPath, d.Name()))
					if dirSize > int64(viper.GetInt("mindirsize")*1000000) {
						dir := Dir{}
						dir.Path = filepath.Join(dirPath, d.Name())
						dir.Depth = viper.GetInt("depth")
						dir.BytesSize = dirSize
						dir.PrettyBytesSize = common.PrettyBytes(dirSize)

						dirs = append(dirs, dir)
					}

				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].BytesSize > dirs[j].BytesSize
	})

	return dirs, nil
}

// DirSize returns the size of a directory in bytes.
func DirSizeBytes(dirPath string) (int64, error) {
	var size int64
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}

func PrintResults(dirs []Dir) {
	fmt.Println()
	common.PrintColor("forestgreen", "background", fmt.Sprintf("Largest directories in: %s", viper.GetString("path")))
	fmt.Println("----------------------------------------------------------")
	fmt.Println()

	spacing := make(map[string]int)
	highWaterMark := 0

	for _, dir := range dirs {
		if len(dir.Path) > highWaterMark {
			highWaterMark = len(dir.Path)
		}

		spacing[dir.Path] = len(dir.Path)
	}

	for _, dir := range dirs {
		padding := strconv.Itoa(highWaterMark + 2)
		if dir.BytesSize >= int64(viper.GetInt("highlight")*1000000) {
			color.HEXStyle("000", common.AllHex["yellow1"]).Printf("%-"+padding+"s %10s\n", dir.Path, dir.PrettyBytesSize)
		} else {
			color.HEXStyle(common.AllHex["steelblue2"]).Printf("%-"+padding+"s %10s\n", dir.Path, dir.PrettyBytesSize)
		}
	}
}
