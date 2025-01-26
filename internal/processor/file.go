package processor

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ryliey/RuleTrans/internal/converter"
	"github.com/Ryliey/RuleTrans/internal/doc"
	"github.com/Ryliey/RuleTrans/internal/git"
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

func (p *FileProcessor) Process(changes []git.FileChange) {
	processedDirs := make(map[string]bool)

	for _, fc := range changes {
		switch {
		case strings.HasSuffix(fc.Path, ".json"):
			p.handleSingBoxFile(fc, processedDirs)
		case strings.HasSuffix(fc.Path, ".yaml"):
			p.handleClashFile(fc, processedDirs)
		}
	}
}

func (p *FileProcessor) handleClashFile(fc git.FileChange, processedDirs map[string]bool) {
	switch fc.Status {
	case "A", "M":
		if err := p.ClashConverter.Convert(fc.Path); err != nil {
			log.Printf("Error converting file %s: %v", fc.Path, err)
		} else {
			// 更新 Clash 目录的 README
			updateReadme(fc.Path)

			// 获取转换后的目标路径（Sing-Box 目录）
			targetPath := p.ClashConverter.GetTargetPath(fc.Path)
			// 更新 Sing-Box 目录的 README
			updateReadme(targetPath)
		}
	case "D":
		p.cleanTargetResources(fc.Path, p.ClashConverter, processedDirs)
	}
}

func (p *FileProcessor) handleSingBoxFile(fc git.FileChange, processedDirs map[string]bool) {
	switch fc.Status {
	case "A", "M":
		if err := p.SingBoxConverter.Convert(fc.Path); err != nil {
			log.Printf("Error converting file %s: %v", fc.Path, err)
		} else {
			// 更新 Sing-Box 目录的 README
			updateReadme(fc.Path)

			// 获取转换后的目标路径（Clash 目录）
			targetPath := p.SingBoxConverter.GetTargetPath(fc.Path)
			// 更新 Clash 目录的 README
			updateReadme(targetPath)
		}
	case "D":
		p.cleanTargetResources(fc.Path, p.SingBoxConverter, processedDirs)
	}
}

func updateReadme(path string) {
	dir := filepath.Dir(path)
	readmePath := filepath.Join(dir, "README.md")
	log.Printf("Generating README for directory: %s", dir)
	if err := doc.GenerateReadme(readmePath); err != nil {
		log.Printf("Error updating README: %v", err)
	} else {
		log.Printf("Successfully generated README: %s", readmePath)
	}
}

func (p *FileProcessor) cleanTargetResources(path string, conv converter.Converter, processedDirs map[string]bool) {
	sourceDir := filepath.Dir(path)

	targetDir := conv.GetTargetPath(sourceDir)

	// 手动移除目标目录的扩展名
	targetDirWithoutExt := strings.TrimSuffix(targetDir, filepath.Ext(targetDir))

	log.Printf("Directory mapping: [%s] => [%s]", sourceDir, targetDirWithoutExt)

	// 需要删除的目录列表
	dirsToDelete := []string{
		sourceDir,           // 原始目录
		targetDirWithoutExt, // 目标目录（移除扩展名）
	}

	for _, dir := range dirsToDelete {
		if processedDirs[dir] {
			continue
		}
		processedDirs[dir] = true

		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			log.Printf("Delete error: %v", err)
		} else {
			log.Printf("Successfully deleted: %s", dir)
		}
	}
}
