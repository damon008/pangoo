package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ExecutableFilePath(name string) string {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, name)
}

func ExecutableFileArch(path string) string {
	fileCmd := fmt.Sprintf("file %s", path)
	out := BashEval(fileCmd)
	arm64 := []string{"aarch64", "arm64"}
	amd64 := []string{"x86-64", "x86_64"}

	for _, a := range arm64 {
		if strings.Contains(out, a) {
			return "arm64"
		}
	}
	for _, a := range amd64 {
		if strings.Contains(out, a) {
			return "amd64"
		}
	}
	return ""
}

// FetchSealosAbsPath 获取sealos绝对路径
func FetchSealosAbsPath() string {
	ex, _ := os.Executable()
	exPath, _ := filepath.Abs(ex)
	return exPath
}
