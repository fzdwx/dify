package dify

// IndexingTechnique 索引模式
type IndexingTechnique string

const (
	IndexingTechniqueHighQuality IndexingTechnique = "high_quality" // 高质量
	IndexingTechniqueEconomy     IndexingTechnique = "economy"      // 经济
)

// DatasetPermission 知识库权限
type DatasetPermission string

const (
	DatasetPermissionOnlyMe         DatasetPermission = "only_me"          // 仅自己
	DatasetPermissionAllTeamMembers DatasetPermission = "all_team_members" // 所有团队成员
	DatasetPermissionPartialMembers DatasetPermission = "partial_members"  // 部分成员
)

type DatasetProvider string

const (
	DatasetProviderVendor   DatasetProvider = "vendor"   // 上传文件
	DatasetProviderExternal DatasetProvider = "external" //  外部知识库
)

type RetrievalModelSearchMethod string

const (
	RetrievalModelSearchMethodHybridSearch   RetrievalModelSearchMethod = "hybrid_search"    // 混合检索
	RetrievalModelSearchMethodSemanticSearch RetrievalModelSearchMethod = "semantic_search"  // 语义检索
	RetrievalModelSearchMethodFullTextSearch RetrievalModelSearchMethod = "full_text_search" // 全文检索
)

type DocForm string

const (
	DocFormTextModel         DocForm = "text_model"         // ext 文档直接 embedding，经济模式默认为该模式
	DocFormHierarchicalModel DocForm = "hierarchical_model" // parent-child 模式
	DocFormQAModel           DocForm = "qa_model"           //  Q&A 模式：为分片文档生成 Q&A 对，然后对问题进行 embedding
)

type ProcessMode string

const (
	ProcessModeAutomatic ProcessMode = "automatic" // 自动处理
	ProcessModeManual    ProcessMode = "custom"    // 自定义
)

type PreProcessingRulesID string

const (
	PreProcessingRulesIDRemoveExtraSpaces PreProcessingRulesID = "remove_extra_spaces" // 替换连续空格、换行符、制表符
	PreProcessingRulesIDRemoveUrlsEmails  PreProcessingRulesID = "remove_urls_emails"  // 删除 URL、电子邮件地址
)
