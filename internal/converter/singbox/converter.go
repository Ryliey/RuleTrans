package singbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Ryliey/RuleTrans/internal/converter"
	"github.com/Ryliey/RuleTrans/pkg/fileutil"
	"github.com/Ryliey/RuleTrans/pkg/types"
	"gopkg.in/yaml.v3"
)

type SingBoxConverter struct {
	converter.BaseConverter
}

func NewConverter() *SingBoxConverter {
	return &SingBoxConverter{
		BaseConverter: converter.BaseConverter{
			SourceDir: "Sing-Box",
			TargetDir: "Clash",
			SourceExt: ".json",
			TargetExt: ".yaml",
		},
	}
}

func (c *SingBoxConverter) GetTargetPath(sourcePath string) string {
	targetPath := fileutil.ConvertPath(sourcePath, c.SourceDir, c.TargetDir)
	return fileutil.ChangeExtension(targetPath, c.TargetExt)
}

func (c *SingBoxConverter) Convert(path string) error {
	// 读取 JSON 文件
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	// 解析 JSON
	var singboxConfig types.SingBoxConfig
	if err := json.Unmarshal(jsonData, &singboxConfig); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 处理规则
	reverseMapping := types.GetReverseMapping()
	var clashRules []string

	for _, rule := range singboxConfig.Rules {
		ruleValue := reflect.ValueOf(rule)
		ruleType := ruleValue.Type()

		for i := 0; i < ruleValue.NumField(); i++ {
			field := ruleValue.Field(i)
			jsonTag := ruleType.Field(i).Tag.Get("json")
			if jsonTag == "" {
				continue
			}
			jsonTag = jsonTag[:strings.Index(jsonTag, ",")]

			if field.Kind() == reflect.Slice && !field.IsNil() && field.Len() > 0 {
				if clashType, exists := reverseMapping[jsonTag]; exists {
					for j := 0; j < field.Len(); j++ {
						value := field.Index(j).String()
						clashRules = append(clashRules, fmt.Sprintf("%s,%s", clashType, value))
					}
				}
			}
		}
	}

	// 创建并输出 Clash 配置
	yamlData, err := yaml.Marshal(types.ClashConfig{Payload: clashRules})
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	// 构造输出路径
	outputPath := strings.Replace(path, "Sing-Box", "Clash", 1)
	outputPath = strings.Replace(outputPath, ".json", ".yaml", 1)

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 写入 YAML 文件
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	fmt.Printf("Successfully converted and saved to: %s\n", outputPath)
	return nil
}
