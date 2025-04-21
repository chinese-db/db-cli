package cmd

import "github.com/spf13/cobra"

// rootCmd 定义根命令
var rootCmd = &cobra.Command{
	Use:   "your-cli",
	Short: "微服务框架脚手架工具 - 生成 RPC 服务和 API 网关(董博)",
}

// Execute 启动 CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
