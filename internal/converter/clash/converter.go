package clash

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Ryliey/RuleTrans/internal/converter"
	"github.com/Ryliey/RuleTrans/pkg/types"
	"gopkg.in/yaml.v3"
)

type ClashConverter struct {
	converter.BaseConverter
}

func NewConverter() *ClashConverter {
	return &ClashConverter{
		BaseConverter: converter.BaseConverter{
			SourceDir: "Clash",
			TargetDir: "Sing-Box",
			SourceExt: ".yaml",
			TargetExt: ".json",
		},
	}
}

func (c *ClashConverter) Convert(path string) error {
	// 读取 YAML 文件
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// 解析 YAML
	var clashConfig map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &clashConfig); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// 处理 payload
	payload, ok := clashConfig["payload"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid YAML format: 'payload' not found or of incorrect type")
	}

	// 创建规则
	rules := types.Rule{}
	ruleValue := reflect.ValueOf(&rules).Elem()
	fieldMap := make(map[string]reflect.Value)

	// 构建字段映射
	for i := 0; i < ruleValue.NumField(); i++ {
		field := ruleValue.Field(i)
		jsonTag := strings.Split(ruleValue.Type().Field(i).Tag.Get("json"), ",")[0]
		fieldMap[jsonTag] = field
	}

	// 处理规则
	for _, item := range payload {
		parts := strings.Split(fmt.Sprint(item), ",")
		if len(parts) < 2 {
			continue
		}

		clashType := parts[0]
		if singboxType, exists := types.RuleTypeMapping[clashType]; exists {
			if field, ok := fieldMap[singboxType]; ok {
				value := parts[1]
				field.Set(reflect.Append(field, reflect.ValueOf(value)))
			}
		}
	}

	// 创建并输出 sing-box 配置
	output := types.SingBoxConfig{
		Version: 3,
		Rules:   []types.Rule{rules},
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate JSON: %w", err)
	}

	// 构造输出路径
	outputPath := strings.Replace(path, "Clash", "Sing-Box", 1)
	outputPath = strings.Replace(outputPath, ".yaml", ".json", 1)

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 写入 JSON 文件
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("Successfully converted and saved to: %s\n", outputPath)
	return nil
}
