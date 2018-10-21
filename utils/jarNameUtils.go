package utils

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func NormalizeJarName(jarName string) string {
	const internalImage = "internal_image"
	normalizeString := strings.TrimSuffix(strings.ToLower(jarName), filepath.Ext(jarName))

	return strings.Join([]string{internalImage, normalizeString}, "/")
}

func GetDockerFilePath() string{
	_, rootPath, _, _ := runtime.Caller(1)
	return filepath.Join(path.Dir(rootPath), "../Dockerfile.tar.gz")
}