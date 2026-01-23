package gigachat_adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/google/uuid"
)

const (
	authURL        = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	completionsURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
)

type GigaChatAiAPI struct {
	apiKey     string
	httpClient *http.Client
	prompts    *Prompts
}

type Prompts struct {
	SystemPrompt       string `json:"system_prompt"`
	UserPromptTemplate string `json:"user_prompt_template"`
}

func NewGigaChatAiAPI(apiKey, promptsPath string) (*GigaChatAiAPI, error) {
	prompts, err := loadPrompts(promptsPath)
	if err != nil {
		return nil, err
	}

	return &GigaChatAiAPI{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		prompts: prompts,
	}, nil
}

// GenerateNode генерирует новую ноду, используя ВСЮ цепочку родителей
func (a *GigaChatAiAPI) GenerateNode(
	ctx context.Context,
	parentNode *dilemma_entity.DilemmaNode,
) (*dilemma_entity.DilemmaNode, error) {
	const op = "ai - GigaChatAdapter - GenerateNode"

	accessToken, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, berrors.Wrap(op, "failed to get access token", err)
	}

	userPrompt, err := a.buildUserPrompt(parentNode)
	if err != nil {
		return nil, berrors.Wrap(op, "failed to build user prompt", err)
	}

	reqBody := CompletionRequest{
		Model: "GigaChat:latest",
		Messages: []Message{
			{
				Role:    "system",
				Content: a.prompts.SystemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		Temperature: 0.7,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		completionsURL,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, berrors.InternalFromErr(
			op,
			fmt.Errorf("gigachat API error: status=%d body=%s", resp.StatusCode, string(body)),
		)
	}

	var completionResp CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	if len(completionResp.Choices) == 0 {
		return nil, berrors.InternalFromErr(
			op,
			fmt.Errorf("gigachat API returned no choices"),
		)
	}

	content := completionResp.Choices[0].Message.Content

	// строго ожидаем JSON
	var dto struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal([]byte(content), &dto); err != nil {
		return nil, berrors.InternalFromErr(
			op,
			fmt.Errorf("invalid AI response, expected JSON, got: %s", content),
		)
	}

	if dto.Name == "" || dto.Value == "" {
		return nil, berrors.InternalFromErr(
			op,
			fmt.Errorf("AI returned empty name or value: %s", content),
		)
	}

	return &dilemma_entity.DilemmaNode{
		Name:  dto.Name,
		Value: dto.Value,
	}, nil
}

func (a *GigaChatAiAPI) getAccessToken(ctx context.Context) (string, error) {
	payload := strings.NewReader("scope=GIGACHAT_API_PERS")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, payload)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+a.apiKey)
	req.Header.Set("RqUID", uuid.New().String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"gigachat auth error: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.AccessToken, nil
}

// buildUserPrompt строит prompt на основе ВСЕЙ цепочки родителей
func (a *GigaChatAiAPI) buildUserPrompt(
	node *dilemma_entity.DilemmaNode,
) (string, error) {

	chain := buildContextChain(node)

	tmpl, err := template.New("user_prompt").Parse(a.prompts.UserPromptTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, chain); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type NodeContext struct {
	Name  string
	Value string
}

// от корня → к текущей ноде
func buildContextChain(node *dilemma_entity.DilemmaNode) []NodeContext {
	var chain []NodeContext

	for n := node; n != nil; n = n.Parent {
		chain = append(chain, NodeContext{
			Name:  n.Name,
			Value: n.Value,
		})
	}

	// reverse
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	return chain
}

func loadPrompts(path string) (*Prompts, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var prompts Prompts
	if err := json.NewDecoder(file).Decode(&prompts); err != nil {
		return nil, err
	}

	return &prompts, nil
}

// ===== GigaChat API DTO =====

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}
