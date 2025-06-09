package dify

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"resty.dev/v3"
	"strings"
)

// CreateChatApp 创建聊天应用
func (c *client) CreateChatApp(ctx context.Context, req *CreateChatAppRequest) (*Response[CreateChatAppResponse], error) {
	var resultErr error
	var resp = &CreateChatAppResponse{}
	var finalResponse *resty.Response

	_, err := c.executeConsoleWithRetry(func() (*resty.Response, error) {
		response, err := c.console().
			WithContext(ctx).
			SetContentType("application/json").
			SetBody(&CreateChatAppInternalRequest{
				Name:           req.Name,
				IconType:       "emoji",
				Icon:           "🤖",
				IconBackground: "#FFEAD5",
				Mode:           "chat",
				Description:    "Created by Dify Go Client",
			}).
			SetResult(&resp).
			Post("/console/api/apps")

		finalResponse = response

		if err != nil {
			resultErr = fmt.Errorf("failed to create chat app: %w", err)
			return response, err
		}

		if response.IsError() {
			resultErr = fmt.Errorf("failed to create chat app with status %d: %s", response.StatusCode(), response.String())
			return response, nil // Don't return error here, let executeConsoleWithRetry handle 401
		}

		return response, nil
	})

	if err != nil {
		return nil, err
	}

	if resultErr != nil {
		return nil, resultErr
	}

	return buildResponse[CreateChatAppResponse](finalResponse, resp), nil
}

// UpdateAppModelConfig 更新应用模型配置
func (c *client) UpdateAppModelConfig(ctx context.Context, req *UpdateAppModelConfigRequest) (*Response[UpdateAppModelConfigResponse], error) {
	var resultErr error
	var resp = &UpdateAppModelConfigResponse{}
	var finalResponse *resty.Response

	// 构建数据集配置
	var datasets []DatasetConfig
	if req.DatasetID != "" {
		datasets = append(datasets, DatasetConfig{
			Dataset: DatasetInfo{
				Enabled: true,
				ID:      req.DatasetID,
			},
		})
	}

	_, err := c.executeConsoleWithRetry(func() (*resty.Response, error) {
		response, err := c.console().
			WithContext(ctx).
			SetContentType("application/json").
			SetBody(&UpdateAppModelConfigInternalRequest{
				PrePrompt:                     "根据用户提问进行回答",
				PromptType:                    "simple",
				ChatPromptConfig:              map[string]interface{}{},
				CompletionPromptConfig:        map[string]interface{}{},
				UserInputForm:                 []interface{}{},
				DatasetQueryVariable:          "",
				MoreLikeThis:                  MoreLikeThisConfig{Enabled: false},
				OpeningStatement:              "",
				SuggestedQuestions:            []interface{}{},
				SensitiveWordAvoidance:        SensitiveWordAvoidanceConfig{Enabled: false, Type: "", Configs: []interface{}{}},
				SpeechToText:                  SpeechToTextConfig{Enabled: false},
				TextToSpeech:                  TextToSpeechConfig{Enabled: false},
				FileUpload:                    getDefaultFileUploadConfig(),
				SuggestedQuestionsAfterAnswer: SuggestedQuestionsAfterAnswerConfig{Enabled: false},
				RetrieverResource:             RetrieverResourceConfig{Enabled: true},
				AgentMode:                     getDefaultAgentModeConfig(),
				Model:                         req.Model,
				DatasetConfigs: DatasetConfigs{
					RetrievalModel:  "multiple",
					TopK:            4,
					RerankingMode:   "reranking_model",
					RerankingModel:  RerankingModelConfig{RerankingProviderName: "", RerankingModelName: ""},
					RerankingEnable: false,
					Datasets:        DatasetsWrapper{Datasets: datasets},
				},
			}).
			SetResult(&resp).
			Post(fmt.Sprintf("/console/api/apps/%s/model-config", req.AppID))

		finalResponse = response

		if err != nil {
			resultErr = fmt.Errorf("failed to update app model config: %w", err)
			return response, err
		}

		if response.IsError() {
			resultErr = fmt.Errorf("failed to update app model config with status %d: %s", response.StatusCode(), response.String())
			return response, nil // Don't return error here, let executeConsoleWithRetry handle 401
		}

		return response, nil
	})

	if err != nil {
		return nil, err
	}

	if resultErr != nil {
		return nil, resultErr
	}

	return buildResponse[UpdateAppModelConfigResponse](finalResponse, resp), nil
}

// CreateAppAccessToken 创建应用访问令牌
func (c *client) CreateAppAccessToken(ctx context.Context, req *CreateAppAccessTokenRequest) (*Response[CreateAppAccessTokenResponse], error) {
	var resultErr error
	var resp = &CreateAppAccessTokenResponse{}
	var finalResponse *resty.Response

	_, err := c.executeConsoleWithRetry(func() (*resty.Response, error) {
		response, err := c.console().
			WithContext(ctx).
			SetContentType("application/json").
			SetResult(&resp).
			Post(fmt.Sprintf("/console/api/apps/%s/api-keys", req.AppID))

		finalResponse = response

		if err != nil {
			resultErr = fmt.Errorf("failed to create app access token: %w", err)
			return response, err
		}

		if response.IsError() {
			resultErr = fmt.Errorf("failed to create app access token with status %d: %s", response.StatusCode(), response.String())
			return response, nil // Don't return error here, let executeConsoleWithRetry handle 401
		}

		return response, nil
	})

	if err != nil {
		return nil, err
	}

	if resultErr != nil {
		return nil, resultErr
	}

	return buildResponse[CreateAppAccessTokenResponse](finalResponse, resp), nil
}

// 辅助函数
func getDefaultFileUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		Image: ImageUploadConfig{
			Detail:          "high",
			Enabled:         false,
			NumberLimits:    3,
			TransferMethods: []string{"remote_url", "local_file"},
		},
		Enabled:                  false,
		AllowedFileTypes:         []string{},
		AllowedFileExtensions:    []string{".JPG", ".JPEG", ".PNG", ".GIF", ".WEBP", ".SVG", ".MP4", ".MOV", ".MPEG", ".WEBM"},
		AllowedFileUploadMethods: []string{"remote_url", "local_file"},
		NumberLimits:             3,
	}
}

func (c *client) CallWorkflowAppBlocking(ctx context.Context, req *CallWorkflowRequest) (*Response[CallWorkflowCompletionResponse], error) {
	var resultErr error
	var resp = &CallWorkflowCompletionResponse{}
	var finalResponse *resty.Response

	req.ResponseMode = ResponseModeBlocking
	_, err := c.executeConsoleWithRetry(func() (*resty.Response, error) {
		response, err := c.console().
			WithContext(ctx).
			SetContentType("application/json").
			SetBody(req).
			SetHeader("Authorization", "Bearer "+req.Token).
			SetResult(&CallWorkflowCompletionResponse{}).
			Post(fmt.Sprintf("/v1/workflows/run"))

		finalResponse = response

		if err != nil {
			return response, fmt.Errorf("failed to call workflow app: %w", err)
		}

		if response.IsError() {
			resultErr = fmt.Errorf("failed to call workflow app with status %d: %s", response.StatusCode(), response.String())
			return response, nil // Don't return error here, let executeConsoleWithRetry handle 401
		}

		resp = response.Result().(*CallWorkflowCompletionResponse)
		return response, nil
	})
	if err != nil {
		return nil, err
	}
	if resultErr != nil {
		return nil, resultErr
	}
	return buildResponse[CallWorkflowCompletionResponse](finalResponse, resp), nil
}

func (c *client) CallWorkflowAppStreaming(ctx context.Context, req *CallWorkflowRequest) (chan *CallWorkflowChunkCompletionResponse, error) {
	var resultErr error
	req.ResponseMode = ResponseModeStreaming
	resp, err := c.executeConsoleWithRetry(func() (*resty.Response, error) {
		response, err := c.console().
			SetDoNotParseResponse(true).
			WithContext(ctx).
			SetContentType("application/json").
			SetBody(req).
			SetHeader("Authorization", "Bearer "+req.Token).
			Post(fmt.Sprintf("/v1/workflows/run"))

		if err != nil {
			return response, fmt.Errorf("failed to call workflow app: %w", err)
		}
		if response.IsError() {
			resultErr = fmt.Errorf("failed to call workflow app with status %d: %s", response.StatusCode(), response.String())
			return response, nil // Don't return error here, let executeConsoleWithRetry handle 401
		}
		return response, nil
	})
	if err != nil {
		return nil, err
	}
	if resultErr != nil {
		return nil, resultErr
	}

	// 处理 SSE 流
	chunks := make(chan *CallWorkflowChunkCompletionResponse)
	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if matched, _ := regexp.MatchString(`^(data|event|id|retry):\s?.*`, line); matched {
				sseLine, err := parseSSELine(line)
				if err != nil {
					continue
				}
				chunks <- sseLine
			}
		}
	}()
	return chunks, nil
}

func parseSSELine(line string) (*CallWorkflowChunkCompletionResponse, error) {
	if strings.HasPrefix(line, "data: ") {
		rawJSON := strings.TrimPrefix(line, "data: ")
		var msg CallWorkflowChunkCompletionResponse
		err := json.Unmarshal([]byte(rawJSON), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	}
	return nil, fmt.Errorf("invalid line: does not start with 'data: '")
}

func getDefaultAgentModeConfig() AgentModeConfig {
	return AgentModeConfig{
		Enabled:      false,
		MaxIteration: 5,
		Strategy:     "function_call",
		Tools:        []interface{}{},
	}
}
