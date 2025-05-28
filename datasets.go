package dify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func (c *client) CreateEmptyDataset(ctx context.Context, req *CreateEmptyDatasetRequest) (*Response[CreateEmptyDatasetResponse], error) {
	var resp = &CreateEmptyDatasetResponse{}
	response, err := c.r().
		WithContext(ctx).
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
	ID                     string            `json:"id"`
	Name                   string            `json:"name"`
	Description            string            `json:"description"`
	Provider               string            `json:"provider"`
	Permission             string            `json:"permission"`
	DataSourceType         string            `json:"data_source_type"`
	IndexingTechnique      IndexingTechnique `json:"indexing_technique"`
	AppCount               int               `json:"app_count"`
	DocumentCount          int               `json:"document_count"`
	WordCount              int               `json:"word_count"`
	CreatedBy              string            `json:"created_by"`
	CreatedAt              int               `json:"created_at"`
	UpdatedBy              string            `json:"updated_by"`
	UpdatedAt              int               `json:"updated_at"`
	EmbeddingModel         string            `json:"embedding_model"`
	EmbeddingModelProvider string            `json:"embedding_model_provider"`
	EmbeddingAvailable     bool              `json:"embedding_available"`
}

func (c *client) CreateByFile(ctx context.Context, req *CreateByFileRequest) (*Response[CreateByFileResponse], error) {
	bytes, jsonErr := json.Marshal(req)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to marshal CreateByFileRequest: %w", jsonErr)
	}
	var resp = &CreateByFileResponse{}
	response, err := c.r().
		WithContext(ctx).
		SetFileReader("file", req.Filename, req.FileBody).
		SetFormData(map[string]string{
			"data": string(bytes),
		}).
		SetResult(&resp).
		Post(fmt.Sprintf("/datasets/%s/document/create-by-file", req.DatasetsID))
	if err != nil {
		return nil, err
	}
	return buildResponse[CreateByFileResponse](response, resp), nil
}

type CreateByFileRequest struct {
	DatasetsID        string
	Filename          string
	FileBody          io.Reader
	IndexingTechnique IndexingTechnique `json:"indexing_technique"` // 索引方式
	DocForm           DocForm           `json:"doc_form"`           // 索引内容的形式
	DocLanguage       string            `json:"doc_language"`       //  在 Q&A 模式下，指定文档的语言，例如：English、Chinese
	ProcessRule       ProcessRule       `json:"process_rule"`       // 文档处理规则
}

type ProcessRule struct {
	Mode  ProcessMode  `json:"mode"`  // 清洗、分段模式 ，automatic 自动 / custom 自定义
	Rules ProcessRules `json:"rules"` // 自定义规则（自动模式下，该字段为空）
}

type PreProcessingRules struct {
	ID      PreProcessingRulesID `json:"id"`
	Enabled bool                 `json:"enabled"`
}

type ProcessRules struct {
	PreProcessingRules []PreProcessingRules `json:"pre_processing_rules"` // 文档预处理规则
	Segmentation       Segmentation         `json:"segmentation"`         // 分段规则
}

type Segmentation struct {
	Separator    string `json:"separator"`     //  自定义分段标识符，目前仅允许设置一个分隔符。默认为 \n
	MaxTokens    int64  `json:"max_tokens"`    // 最大长度（token）默认为 1000
	ChunkOverlap int64  `json:"chunk_overlap"` // 分段重叠长度（token），默认为 0
}

type CreateByFileResponse struct {
	Document struct {
		Id             string `json:"id"`
		Position       int    `json:"position"`
		DataSourceType string `json:"data_source_type"`
		DataSourceInfo struct {
			UploadFileId string `json:"upload_file_id"`
		} `json:"data_source_info"`
		DatasetProcessRuleId string      `json:"dataset_process_rule_id"`
		Name                 string      `json:"name"`
		CreatedFrom          string      `json:"created_from"`
		CreatedBy            string      `json:"created_by"`
		CreatedAt            int64       `json:"created_at"`
		Tokens               int64       `json:"tokens"`
		IndexingStatus       string      `json:"indexing_status"`
		Error                interface{} `json:"error"`
		Enabled              bool        `json:"enabled"`
		DisabledAt           int64       `json:"disabled_at"`
		DisabledBy           string      `json:"disabled_by"`
		Archived             bool        `json:"archived"`
		DisplayStatus        string      `json:"display_status"`
		WordCount            int64       `json:"word_count"`
		HitCount             int64       `json:"hit_count"`
		DocForm              string      `json:"doc_form"`
	} `json:"document"`
	Batch string `json:"batch"`
}
