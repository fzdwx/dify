package dify

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDatasets(t *testing.T) {
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// Use timestamp to create unique dataset name
	uniqueName := fmt.Sprintf("测试数据集_%d", time.Now().Unix())
	resp, err := c.CreateEmptyDataset(ctx, &CreateEmptyDatasetRequest{
		Name:              uniqueName,
		Description:       "测试数据集描述",
		IndexingTechnique: IndexingTechniqueEconomy,
		Permission:        DatasetPermissionAllTeamMembers,
		Provider:          DatasetProviderVendor,
		RetrievalModel: RetrievalModel{
			SearchMethod:    RetrievalModelSearchMethodHybridSearch,
			RerankingEnable: true,
			TopK:            10,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsSuccess() {
		t.Fatal(resp.Message)
	}
	t.Log(resp)

	file, err := os.OpenFile("./testdata/aaaa.docx", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	createFileResp, err := c.CreateByFile(ctx, &CreateByFileRequest{
		DatasetsID:        resp.Result.ID,
		Filename:          "aaaa.docx",
		FileBody:          file,
		IndexingTechnique: IndexingTechniqueEconomy,
		DocForm:           DocFormTextModel,
		ProcessRule: ProcessRule{
			Mode: ProcessModeManual,
			Rules: ProcessRules{
				PreProcessingRules: []PreProcessingRules{
					{
						ID:      PreProcessingRulesIDRemoveExtraSpaces,
						Enabled: true,
					},
					{
						ID:      PreProcessingRulesIDRemoveUrlsEmails,
						Enabled: false,
					},
				},
				Segmentation: Segmentation{
					Separator:    "\n\n",
					MaxTokens:    1024,
					ChunkOverlap: 50,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(createFileResp)
}
