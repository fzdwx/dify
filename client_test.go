package dify

import "testing"

func TestDatasets(t *testing.T) {
	c := NewClient("http://192.168.50.21:88/v1", "dataset-22i7BMiiZaobMzRTLffm3mDX")
	resp, err := c.CreateEmptyDataset(&CreateEmptyDatasetRequest{
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
