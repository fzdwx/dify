package dify

import (
	"context"
	"testing"
)

func TestDatasets(t *testing.T) {
	c := NewClient("http://192.168.50.21:88/v1", "dataset-22i7BMiiZaobMzRTLffm3mDX")
	ctx := context.Background()
	resp, err := c.CreateEmptyDataset(ctx, &CreateEmptyDatasetRequest{
		Name:              "测试111",
		Description:       "123123123123",
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
	t.Log(resp)
}
