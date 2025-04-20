package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateService 动态生成服务
func GenerateService(serviceType, serviceName, version, port string) error {
	config := NewTemplateConfig()

	// 根据类型选择模板仓库
	var repo string
	switch serviceType {
	case "rpc":
		repo = config.RPCTemplateRepo
	case "api":
		repo = config.APITemplateRepo
	default:
		return fmt.Errorf("错误: 不支持的服务类型 '%s'", serviceType)
	}

	// 1. 拉取模板仓库
	tmpDir, err := cloneTemplateRepo(repo, version)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// 2. 准备模板数据
	data := map[string]interface{}{
		"ServiceName": serviceName,
		"Port":        port,
		// 可扩展其他变量，如数据库配置
	}

	// 3. 渲染模板
	destDir := filepath.Join(".", serviceName)
	if err := renderAllTemplates(tmpDir, destDir, data); err != nil {
		return fmt.Errorf("模板渲染失败: %v", err)
	}

	fmt.Printf("✅ 服务 '%s' 生成成功！目录: %s\n", serviceName, destDir)
	return nil
}

// cloneTemplateRepo 克隆模板仓库到临时目录
func cloneTemplateRepo(repo, version string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "your-cli-*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 支持 Git URL 格式转换（SSH/HTTPS）
	repoURL := convertRepoURL(repo)
	cmd := exec.Command("git", "clone", "--depth=1", "--branch", version, repoURL, tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("克隆仓库失败\n仓库: %s\n版本: %s\n错误: %s", repo, version, output)
	}
	return tmpDir, nil
}

// convertRepoURL 统一仓库地址格式
func convertRepoURL(repo string) string {
	if strings.HasPrefix(repo, "github.com") {
		return "https://" + repo + ".git"
	}
	return repo // 假设其他仓库已为完整 URL
}

// renderAllTemplates 渲染所有模板文件
func renderAllTemplates(srcDir, destDir string, data map[string]interface{}) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略 .git 目录
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// 计算相对路径
		relPath, _ := filepath.Rel(srcDir, path)
		destPath := filepath.Join(destDir, renderPath(relPath, data))

		// 处理目录
		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// 处理文件
		return renderFile(path, destPath, data)
	})
}

// renderPath 替换路径中的变量
func renderPath(path string, data map[string]interface{}) string {
	tpl := template.Must(template.New("path").Parse(path))
	var buf strings.Builder
	if err := tpl.Execute(&buf, data); err != nil {
		return path // 回退原始路径
	}
	return buf.String()
}

// renderFile 渲染单个文件
func renderFile(src, dest string, data map[string]interface{}) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取文件失败: %s", src)
	}

	tpl, err := template.New("file").Parse(string(content))
	if err != nil {
		return fmt.Errorf("解析模板失败: %s", src)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建文件失败: %s", dest)
	}
	defer file.Close()

	return tpl.Execute(file, data)
}
