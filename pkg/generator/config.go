package generator

// TemplateConfig 模板仓库配置
type TemplateConfig struct {
	RPCTemplateRepo string `yaml:"rpc_repo"`        // RPC 模板仓库地址
	APITemplateRepo string `yaml:"api_repo"`        // API 模板仓库地址
	DefaultVersion  string `yaml:"default_version"` // 默认版本
}

// NewTemplateConfig 默认配置（用户可通过配置文件覆盖）
func NewTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		RPCTemplateRepo: "github.com/chinese-db/service",
		APITemplateRepo: "github.com/chinese-db/gateway",
		DefaultVersion:  "main", // 默认使用 main 分支
	}
}
