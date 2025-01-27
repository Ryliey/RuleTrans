package converter

import (
	"github.com/Ryliey/RuleTrans/pkg/fileutil"
)

type Converter interface {
	Convert(path string) error
	GetTargetPath(path string) string
}

type BaseConverter struct {
	SourceDir string
	TargetDir string
	SourceExt string
	TargetExt string
}

func (c *BaseConverter) GetTargetPath(sourcePath string) string {
	// 转换基础路径
	targetPath := fileutil.ConvertPath(sourcePath, c.SourceDir, c.TargetDir)

	// 仅当源路径是文件时修改扩展名
	if !fileutil.IsDir(sourcePath) {
		targetPath = fileutil.ChangeExtension(targetPath, c.TargetExt)
	}

	return targetPath
}
