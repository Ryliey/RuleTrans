package types

// Rule 定义了规则的结构
type Rule struct {
	Domain           []string `json:"domain,omitempty" yaml:"domain,omitempty"`
	DomainSuffix     []string `json:"domain_suffix,omitempty" yaml:"domain_suffix,omitempty"`
	DomainKeyword    []string `json:"domain_keyword,omitempty" yaml:"domain_keyword,omitempty"`
	DomainRegex      []string `json:"domain_regex,omitempty" yaml:"domain_regex,omitempty"`
	IpCidr           []string `json:"ip_cidr,omitempty" yaml:"ip_cidr,omitempty"`
	IpCidr6          []string `json:"ip_cidr6,omitempty" yaml:"ip_cidr6,omitempty"`
	SourceIpCidr     []string `json:"source_ip_cidr,omitempty" yaml:"source_ip_cidr,omitempty"`
	Port             []string `json:"port,omitempty" yaml:"port,omitempty"`
	SourcePort       []string `json:"source_port,omitempty" yaml:"source_port,omitempty"`
	ProcessName      []string `json:"process_name,omitempty" yaml:"process_name,omitempty"`
	ProcessPath      []string `json:"process_path,omitempty" yaml:"process_path,omitempty"`
	ProcessPathRegex []string `json:"process_path_regex,omitempty" yaml:"process_path_regex,omitempty"`
}

// SingBoxConfig 定义了 sing-box 配置的结构
type SingBoxConfig struct {
	Version int    `json:"version"`
	Rules   []Rule `json:"rules"`
}

// ClashConfig 定义了 Clash 配置的结构
type ClashConfig struct {
	Payload []string `yaml:"payload"`
}

// RuleTypeMapping 定义了规则类型的双向映射
var RuleTypeMapping = map[string]string{
	// Clash to Sing-box
	"DOMAIN":             "domain",
	"DOMAIN-SUFFIX":      "domain_suffix",
	"DOMAIN-KEYWORD":     "domain_keyword",
	"DOMAIN-REGEX":       "domain_regex",
	"IP-CIDR":            "ip_cidr",
	"IP-CIDR6":           "ip_cidr6",
	"SRC-IP-CIDR":        "source_ip_cidr",
	"DST-PORT":           "port",
	"SRC-PORT":           "source_port",
	"PROCESS-NAME":       "process_name",
	"PROCESS-PATH":       "process_path",
	"PROCESS-PATH-REGEX": "process_path_regex",
}

// GetReverseMapping 返回反向映射
func GetReverseMapping() map[string]string {
	reverse := make(map[string]string, len(RuleTypeMapping))
	for k, v := range RuleTypeMapping {
		reverse[v] = k
	}
	return reverse
}
