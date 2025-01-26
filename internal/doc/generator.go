// doc/doc.go
package doc

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const readmeTmpl = `# {{.Name}}

## 链接

**Github**
` + "```" + `
{{.GithubURL}}
` + "```" + `

**CDN**
` + "```" + `
{{.CDNURL}}
` + "```" + `

**GitHub Proxy**
` + "```" + `
{{.ProxyURL}}
` + "```" + `
`

type ReadmeData struct {
	Name      string
	GithubURL string
	CDNURL    string
	ProxyURL  string
}

// GenerateReadme 生成或更新 README.md 文件
func GenerateReadme(readmePath string) error {
	// 从仓库根目录获取相对路径（统一转换为正斜杠）
	relPath := strings.ReplaceAll(filepath.Dir(readmePath), string(filepath.Separator), "/")
	relPath = strings.TrimPrefix(relPath, ".")
	relPath = strings.TrimPrefix(relPath, "/")

	// 从目录名获取规则名称
	ruleName := filepath.Base(filepath.Dir(readmePath))

	// 生成 URL（强制使用正斜杠）
	baseURL := "https://raw.githubusercontent.com/Ryliey/Rules/main/"
	data := ReadmeData{
		Name:      ruleName,
		GithubURL: baseURL + relPath + "/" + ruleName + getExtension(relPath),
		CDNURL:    "https://cdn.jsdelivr.net/gh/ryliey/Rules@main/" + relPath + "/" + ruleName + getExtension(relPath),
		ProxyURL:  "https://ghgo.xyz/" + baseURL + relPath + "/" + ruleName + getExtension(relPath),
	}

	// 解析并执行模板
	tmpl, err := template.New("readme").Parse(readmeTmpl)
	if err != nil {
		return err
	}

	// 创建或打开 README.md 文件
	f, err := os.Create(readmePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 将模板数据写入文件
	return tmpl.Execute(f, data)
}

// getExtension 根据路径判断文件扩展名
func getExtension(path string) string {
	if strings.Contains(path, "Clash") {
		return ".yaml"
	} else if strings.Contains(path, "Sing-Box") {
		return ".json"
	}
	return ""
}
