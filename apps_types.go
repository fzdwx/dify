package dify

// CreateChatAppRequest 创建聊天应用请求（简化版，只需要名称）
type CreateChatAppRequest struct {
	Name string `json:"name"` // 应用名称
}

// CreateChatAppInternalRequest 内部使用的完整请求结构
type CreateChatAppInternalRequest struct {
	Name           string `json:"name"`
	IconType       string `json:"icon_type"`
	Icon           string `json:"icon"`
	IconBackground string `json:"icon_background"`
	Mode           string `json:"mode"`
	Description    string `json:"description"`
}

// CreateChatAppResponse 创建聊天应用响应
type CreateChatAppResponse struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Mode                string      `json:"mode"`
	Icon                string      `json:"icon"`
	IconBackground      string      `json:"icon_background"`
	EnableSite          bool        `json:"enable_site"`
	EnableAPI           bool        `json:"enable_api"`
	ModelConfig         interface{} `json:"model_config"`
	Workflow            interface{} `json:"workflow"`
	Tracing             interface{} `json:"tracing"`
	UseIconAsAnswerIcon bool        `json:"use_icon_as_answer_icon"`
	CreatedBy           string      `json:"created_by"`
	CreatedAt           int64       `json:"created_at"`
	UpdatedBy           string      `json:"updated_by"`
	UpdatedAt           int64       `json:"updated_at"`
}

// UpdateAppModelConfigRequest 更新应用模型配置请求（简化版）
type UpdateAppModelConfigRequest struct {
	AppID     string      `json:"-"`     // 应用ID，不包含在JSON中
	Model     ModelConfig `json:"model"` // 模型配置
	DatasetID string      `json:"-"`     // 数据集ID，不包含在JSON中
}

// ModelConfig 模型配置
type ModelConfig struct {
	Provider         string                 `json:"provider"`          // 模型提供商
	Name             string                 `json:"name"`              // 模型名称
	Mode             string                 `json:"mode"`              // 模式，通常是 "chat"
	CompletionParams map[string]interface{} `json:"completion_params"` // 完成参数
}

// UpdateAppModelConfigInternalRequest 内部使用的完整更新请求结构
type UpdateAppModelConfigInternalRequest struct {
	PrePrompt                     string                              `json:"pre_prompt"`
	PromptType                    string                              `json:"prompt_type"`
	ChatPromptConfig              map[string]interface{}              `json:"chat_prompt_config"`
	CompletionPromptConfig        map[string]interface{}              `json:"completion_prompt_config"`
	UserInputForm                 []interface{}                       `json:"user_input_form"`
	DatasetQueryVariable          string                              `json:"dataset_query_variable"`
	MoreLikeThis                  MoreLikeThisConfig                  `json:"more_like_this"`
	OpeningStatement              string                              `json:"opening_statement"`
	SuggestedQuestions            []interface{}                       `json:"suggested_questions"`
	SensitiveWordAvoidance        SensitiveWordAvoidanceConfig        `json:"sensitive_word_avoidance"`
	SpeechToText                  SpeechToTextConfig                  `json:"speech_to_text"`
	TextToSpeech                  TextToSpeechConfig                  `json:"text_to_speech"`
	FileUpload                    FileUploadConfig                    `json:"file_upload"`
	SuggestedQuestionsAfterAnswer SuggestedQuestionsAfterAnswerConfig `json:"suggested_questions_after_answer"`
	RetrieverResource             RetrieverResourceConfig             `json:"retriever_resource"`
	AgentMode                     AgentModeConfig                     `json:"agent_mode"`
	Model                         ModelConfig                         `json:"model"`
	DatasetConfigs                DatasetConfigs                      `json:"dataset_configs"`
}

// UpdateAppModelConfigResponse 更新应用模型配置响应
type UpdateAppModelConfigResponse struct {
	Result string `json:"result,omitempty"`
}

// 配置相关结构体
type MoreLikeThisConfig struct {
	Enabled bool `json:"enabled"`
}

type SensitiveWordAvoidanceConfig struct {
	Enabled bool          `json:"enabled"`
	Type    string        `json:"type"`
	Configs []interface{} `json:"configs"`
}

type SpeechToTextConfig struct {
	Enabled bool `json:"enabled"`
}

type TextToSpeechConfig struct {
	Enabled bool `json:"enabled"`
}

type FileUploadConfig struct {
	Image                    ImageUploadConfig `json:"image"`
	Enabled                  bool              `json:"enabled"`
	AllowedFileTypes         []string          `json:"allowed_file_types"`
	AllowedFileExtensions    []string          `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []string          `json:"allowed_file_upload_methods"`
	NumberLimits             int               `json:"number_limits"`
}

type ImageUploadConfig struct {
	Detail          string   `json:"detail"`
	Enabled         bool     `json:"enabled"`
	NumberLimits    int      `json:"number_limits"`
	TransferMethods []string `json:"transfer_methods"`
}

type SuggestedQuestionsAfterAnswerConfig struct {
	Enabled bool `json:"enabled"`
}

type RetrieverResourceConfig struct {
	Enabled bool `json:"enabled"`
}

type AgentModeConfig struct {
	Enabled      bool          `json:"enabled"`
	MaxIteration int           `json:"max_iteration"`
	Strategy     string        `json:"strategy"`
	Tools        []interface{} `json:"tools"`
}

type DatasetConfigs struct {
	RetrievalModel  string               `json:"retrieval_model"`
	TopK            int                  `json:"top_k"`
	RerankingMode   string               `json:"reranking_mode"`
	RerankingModel  RerankingModelConfig `json:"reranking_model"`
	RerankingEnable bool                 `json:"reranking_enable"`
	Datasets        DatasetsWrapper      `json:"datasets"`
}

type RerankingModelConfig struct {
	RerankingProviderName string `json:"reranking_provider_name"`
	RerankingModelName    string `json:"reranking_model_name"`
}

type DatasetsWrapper struct {
	Datasets []DatasetConfig `json:"datasets"`
}

type DatasetConfig struct {
	Dataset DatasetInfo `json:"dataset"`
}

type DatasetInfo struct {
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
}
