package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"ai_marketer/internal/agent"
	"ai_marketer/internal/config"
	"ai_marketer/internal/llm"
	"ai_marketer/internal/orchestrator"
)

func main() {
	var configPath string

	rootCmd := &cobra.Command{
		Use:   "ai-marketer",
		Short: "AI 멀티 에이전트 SEO 마케터",
		Long:  "마케팅 지식 없이도 웹사이트 SEO 최적화와 상위 노출을 달성할 수 있도록 돕는 AI 멀티 에이전트 시스템",
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "설정 파일 경로")

	runCmd := &cobra.Command{
		Use:   "run [요청]",
		Short: "SEO 캠페인 파이프라인 실행",
		Long:  "전체 SEO 캠페인 파이프라인을 실행합니다 (지식 가이더 → 시행자 → 개발자/수행자 → 신뢰도 담당 → 분석가)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userRequest := strings.Join(args, " ")
			return runPipeline(configPath, userRequest)
		},
	}

	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runPipeline(configPath, userRequest string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("설정 로딩 실패: %w", err)
	}

	client := llm.NewClient(cfg.LLM.APIKey, cfg.LLM.Model)

	// Create and register all agents
	registry := agent.NewRegistry()
	registry.Register(agent.NewKnowledgeGuide(client))
	registry.Register(agent.NewPlanner(client))
	registry.Register(agent.NewDeveloper(client))
	registry.Register(agent.NewExecutor(client))
	registry.Register(agent.NewTrustManager(client))
	registry.Register(agent.NewAnalyst(client))

	// Create orchestrator and run the pipeline
	orch := orchestrator.New(registry)

	log.Printf("SEO 캠페인 시작: %q\n", userRequest)

	wctx, err := orch.Run(context.Background(), orchestrator.FullCampaignPipeline, userRequest)
	if err != nil {
		return fmt.Errorf("파이프라인 실행 실패: %w", err)
	}

	// Print results with clear delimiters for Claude Code parsing
	fmt.Println("\n=== PIPELINE_RESULT START ===")

	for _, name := range []string{"knowledge_guide", "planner", "developer", "executor", "trust_manager", "analyst"} {
		if result, ok := wctx.GetResult(name); ok {
			fmt.Printf("\n=== [%s] START ===\n", name)
			if result.Success {
				output, _ := json.MarshalIndent(result.Output, "", "  ")
				fmt.Println(string(output))
			} else {
				errOut, _ := json.Marshal(map[string]string{
					"error": result.ErrorMessage,
				})
				fmt.Println(string(errOut))
			}
			fmt.Printf("=== [%s] END ===\n", name)
		}
	}

	fmt.Println("\n=== PIPELINE_RESULT END ===")

	return nil
}
