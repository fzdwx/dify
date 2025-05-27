package dify

func (c *client) CreateEmptyDataset(req *CreateEmptyDatasetRequest) (*Response[CreateEmptyDatasetResponse], error) {
	var resp = &CreateEmptyDatasetResponse{}
	response, err := c.r().
		SetContentType("application/json").
		SetBody(&req).
		SetResult(&resp).
		Post("/datasets")
	if err != nil {
		return nil, err
	}
	return buildResponse[CreateEmptyDatasetResponse](response, resp), nil
}

type CreateEmptyDatasetRequest struct {
	Name                   string            `json:"name"`
	Description            string            `json:"description"`
	IndexingTechnique      IndexingTechnique `json:"indexing_technique"`
	Permission             DatasetPermission `json:"permission"`
	Provider               DatasetProvider   `json:"provider"`
	ExternalKnowledgeApiID string            `json:"external_knowledge_api_id"`
	ExternalKnowledgeId    string            `json:"external_knowledge_id"`
	EmbeddingModel         string            `json:"embedding_model"`
	EmbeddingModelProvider string            `json:"embedding_model_provider"`
	RetrievalModel         RetrievalModel    `json:"retrieval_model"`
}

type RetrievalModel struct {
	SearchMethod    RetrievalModelSearchMethod `json:"search_method"`
	RerankingEnable bool                       `json:"reranking_enable"`
	// Rerank 模型配置
	RerankingModel RerankingModel `json:"reranking_model"`
	TopK           int64          `json:"top_k"`
	// 是否开启召回分数限制
	ScoreThresholdEnabled bool `json:"score_threshold_enabled"`
	// 召回分数限制
	ScoreThreshold float64 `json:"score_threshold"`
}

type RerankingModel struct {
	RerankingProviderName string `json:"reranking_provider_name"`
	RerankingModelName    string `json:"reranking_model_name"`
}

type CreateEmptyDatasetResponse struct {
	Id                     string      `json:"id"`
	Name                   string      `json:"name"`
	Description            interface{} `json:"description"`
	Provider               string      `json:"provider"`
	Permission             string      `json:"permission"`
	DataSourceType         interface{} `json:"data_source_type"`
	IndexingTechnique      interface{} `json:"indexing_technique"`
	AppCount               int         `json:"app_count"`
	DocumentCount          int         `json:"document_count"`
	WordCount              int         `json:"word_count"`
	CreatedBy              string      `json:"created_by"`
	CreatedAt              int         `json:"created_at"`
	UpdatedBy              string      `json:"updated_by"`
	UpdatedAt              int         `json:"updated_at"`
	EmbeddingModel         interface{} `json:"embedding_model"`
	EmbeddingModelProvider interface{} `json:"embedding_model_provider"`
	EmbeddingAvailable     interface{} `json:"embedding_available"`
}
