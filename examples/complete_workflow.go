package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fzdwx/dify"
)

func main() {
	// 创建客户端
	client, err := dify.NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}

	ctx := context.Background()

	// 1. 创建数据集
	fmt.Println("=== 创建数据集 ===")
	datasetName := fmt.Sprintf("示例数据集_%d", time.Now().Unix())
	datasetResp, err := client.CreateEmptyDataset(ctx, &dify.CreateEmptyDatasetRequest{
		Name:              datasetName,
		Description:       "用于聊天应用的示例数据集",
		IndexingTechnique: dify.IndexingTechniqueEconomy,
		Permission:        dify.DatasetPermissionAllTeamMembers,
		Provider:          dify.DatasetProviderVendor,
		RetrievalModel: dify.RetrievalModel{
			SearchMethod:    dify.RetrievalModelSearchMethodHybridSearch,
			RerankingEnable: true,
			TopK:            10,
		},
	})

	if err != nil {
		log.Fatal("创建数据集失败:", err)
	}

	if !datasetResp.IsSuccess() {
		log.Fatal("创建数据集失败:", datasetResp.Message)
	}

	fmt.Printf("✅ 成功创建数据集: %s (ID: %s)\n", datasetResp.Result.Name, datasetResp.Result.ID)

	// 2. 创建聊天应用（只需要提供名称）
	fmt.Println("\n=== 创建聊天应用 ===")
	appName := fmt.Sprintf("智能助手_%d", time.Now().Unix())
	appResp, err := client.CreateChatApp(ctx, &dify.CreateChatAppRequest{
		Name: appName,
	})

	if err != nil {
		log.Fatal("创建聊天应用失败:", err)
	}

	if !appResp.IsSuccess() {
		log.Fatal("创建聊天应用失败:", appResp.Message)
	}

	fmt.Printf("✅ 成功创建聊天应用: %s (ID: %s)\n", appResp.Result.Name, appResp.Result.ID)

	// 3. 更新应用配置（设置模型并绑定数据集）
	fmt.Println("\n=== 更新应用配置 ===")
	updateResp, err := client.UpdateAppModelConfig(ctx, &dify.UpdateAppModelConfigRequest{
		AppID: appResp.Result.ID,
		Model: dify.ModelConfig{
			Provider:         "langgenius/deepseek/deepseek",
			Name:             "deepseek-chat",
			Mode:             "chat",
			CompletionParams: map[string]interface{}{},
		},
		DatasetID: datasetResp.Result.ID, // 绑定数据集
	})

	if err != nil {
		log.Fatal("更新应用配置失败:", err)
	}

	if !updateResp.IsSuccess() {
		log.Fatal("更新应用配置失败:", updateResp.Message)
	}

	fmt.Printf("✅ 成功更新应用配置，已绑定数据集\n")

	// 4. 创建应用访问令牌
	fmt.Println("\n=== 创建应用访问令牌 ===")
	tokenResp, err := client.CreateAppAccessToken(ctx, &dify.CreateAppAccessTokenRequest{
		AppID: appResp.Result.ID,
	})

	if err != nil {
		log.Fatal("创建应用访问令牌失败:", err)
	}

	if !tokenResp.IsSuccess() {
		log.Fatal("创建应用访问令牌失败:", tokenResp.Message)
	}

	fmt.Printf("✅ 成功创建应用访问令牌: %s\n", tokenResp.Result.Token)

	// 5. 创建另一个不绑定数据集的应用
	fmt.Println("\n=== 创建简单应用 ===")
	simpleAppName := fmt.Sprintf("简单助手_%d", time.Now().Unix())
	simpleAppResp, err := client.CreateChatApp(ctx, &dify.CreateChatAppRequest{
		Name: simpleAppName,
	})

	if err != nil {
		log.Fatal("创建简单应用失败:", err)
	}

	if !simpleAppResp.IsSuccess() {
		log.Fatal("创建简单应用失败:", simpleAppResp.Message)
	}

	fmt.Printf("✅ 成功创建简单应用: %s (ID: %s)\n", simpleAppResp.Result.Name, simpleAppResp.Result.ID)

	// 6. 更新简单应用配置（只设置模型，不绑定数据集）
	updateSimpleResp, err := client.UpdateAppModelConfig(ctx, &dify.UpdateAppModelConfigRequest{
		AppID: simpleAppResp.Result.ID,
		Model: dify.ModelConfig{
			Provider:         "langgenius/deepseek/deepseek",
			Name:             "deepseek-chat",
			Mode:             "chat",
			CompletionParams: map[string]interface{}{},
		},
		DatasetID: "", // 不绑定数据集
	})

	if err != nil {
		log.Fatal("更新简单应用配置失败:", err)
	}

	if !updateSimpleResp.IsSuccess() {
		log.Fatal("更新简单应用配置失败:", updateSimpleResp.Message)
	}

	fmt.Printf("✅ 成功更新简单应用配置，未绑定数据集\n")

	// 7. 为简单应用创建访问令牌
	simpleTokenResp, err := client.CreateAppAccessToken(ctx, &dify.CreateAppAccessTokenRequest{
		AppID: simpleAppResp.Result.ID,
	})

	if err != nil {
		log.Fatal("创建简单应用访问令牌失败:", err)
	}

	if !simpleTokenResp.IsSuccess() {
		log.Fatal("创建简单应用访问令牌失败:", simpleTokenResp.Message)
	}

	fmt.Printf("✅ 成功创建简单应用访问令牌: %s\n", simpleTokenResp.Result.Token)

	fmt.Println("\n🎉 完整工作流程完成！")
	fmt.Println("\n📋 总结:")
	fmt.Printf("数据集ID: %s\n", datasetResp.Result.ID)
	fmt.Printf("带数据集的应用ID: %s\n", appResp.Result.ID)
	fmt.Printf("带数据集的应用访问令牌: %s\n", tokenResp.Result.Token)
	fmt.Printf("简单应用ID: %s\n", simpleAppResp.Result.ID)
	fmt.Printf("简单应用访问令牌: %s\n", simpleTokenResp.Result.Token)
}
