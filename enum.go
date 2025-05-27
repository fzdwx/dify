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
