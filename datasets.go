package dify

func (c *client) CreateEmptyDataset() {
	//TODO implement me
	panic("implement me")
}

type CreateEmptyDatasetRequest struct {
	Name              string
	Description       string
	IndexingTechnique IndexingTechnique
}

type IndexingTechnique string

const (
	HighQuality IndexingTechnique = "high_quality"
	Economy     IndexingTechnique = "economy"
)
