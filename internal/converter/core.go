package converter

import (
	"fmt"

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
	targetPath := fileutil.ConvertPath(sourcePath, c.SourceDir, c.TargetDir)
	return fileutil.ChangeExtension(targetPath, c.TargetExt)
}

func ProcessFile(conv Converter, path string) error {
	targetPath := conv.GetTargetPath(path)

	if err := fileutil.EnsureDirectory(targetPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return conv.Convert(path)
}
