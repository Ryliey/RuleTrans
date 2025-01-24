package processor

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ryliey/RuleTrans/internal/converter"
	"github.com/Ryliey/RuleTrans/internal/doc"
	"github.com/Ryliey/RuleTrans/internal/git" // 导入 git 包以使用 FileChange 类型
)

type FileProcessor struct {
	ClashConverter   converter.Converter
	SingBoxConverter converter.Converter
}

func NewProcessor(clash, singbox converter.Converter) *FileProcessor {
	return &FileProcessor{
		ClashConverter:   clash,
		SingBoxConverter: singbox,
	}
}

// 使用 git.FileChange 类型
func (p *FileProcessor) Process(changes []git.FileChange) {
	var processedPaths []string

	for _, fc := range changes {
		switch {
		case strings.HasSuffix(fc.Path, ".json"):
			p.handleSingBoxFile(fc, processedPaths)
		case strings.HasSuffix(fc.Path, ".yaml"):
			p.handleClashFile(fc, processedPaths)
		}
	}
}

// 使用 git.FileChange 类型
func (p *FileProcessor) handleSingBoxFile(fc git.FileChange, processedPaths []string) {
	switch fc.Status {
	case "A", "M":
		if err := p.SingBoxConverter.Convert(fc.Path); err != nil {
			log.Printf("Error converting file %s: %v", fc.Path, err)
		}
		updateReadme(fc.Path)
	case "D":
		cleanDeletedFile(fc.Path, p.SingBoxConverter)
	}
}

// 使用 git.FileChange 类型
func (p *FileProcessor) handleClashFile(fc git.FileChange, processedPaths []string) {
	switch fc.Status {
	case "A", "M":
		if err := p.ClashConverter.Convert(fc.Path); err != nil {
			log.Printf("Error converting file %s: %v", fc.Path, err)
		}
		updateReadme(fc.Path)
	case "D":
		cleanDeletedFile(fc.Path, p.ClashConverter)
	}
}

func updateReadme(path string) {
	dir := filepath.Dir(path)
	if err := doc.GenerateReadme(filepath.Join(dir, "README.md")); err != nil {
		log.Printf("Error updating README: %v", err)
	}
}

func cleanDeletedFile(path string, converter converter.Converter) {
	target := converter.GetTargetPath(path)
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		log.Printf("Error removing file: %v", err)
	}
}
