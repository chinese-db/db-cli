package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/chinese-db/db-cli/pkg/generator"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "生成微服务项目（RPC/API）",
	Run: func(cmd *cobra.Command, args []string) {
		// 交互式问答
		answers, err := promptUser()
		if err != nil {
			fmt.Println("❌ 错误:", err)
			return
		}

		// 生成服务
		if err := generator.GenerateService(
			answers.ServiceType,
			answers.ServiceName,
			answers.Version,
			answers.Port,
		); err != nil {
			fmt.Println("❌ 生成失败:", err)
		}
	},
}

func promptUser() (*generator.ServiceConfig, error) {
	answers := &generator.ServiceConfig{}
	qs := []*survey.Question{
		{
			Name: "ServiceType",
			Prompt: &survey.Select{
				Message: "请选择服务类型:",
				Options: []string{"rpc", "api"},
				Default: "rpc",
			},
		},
		{
			Name:     "ServiceName",
			Prompt:   &survey.Input{Message: "输入服务名称（英文,推荐使用_）:"},
			Validate: survey.Required,
		},
		{
			Name:   "Version",
			Prompt: &survey.Input{Message: "输入模板版本（默认 main）:", Default: "main"},
		},
		{
			Name:   "Port",
			Prompt: &survey.Input{Message: "输入服务端口:", Default: "8080"},
		},
	}
	return answers, survey.Ask(qs, answers)
}

func init() {
	rootCmd.AddCommand(createCmd)
}
