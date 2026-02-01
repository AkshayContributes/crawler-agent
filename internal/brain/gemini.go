package brain

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiStrategy struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

func NewGeminiStrategy(apiKey string, modelName string) (*GeminiStrategy, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		fmt.Printf("Error initializing Gemini strategy: %v\n", err)
		return nil, err
	}

	model := client.GenerativeModel(modelName)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	fmt.Println("Initialized Gemini strategy successfully")

	return &GeminiStrategy{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}

func (g *GeminiStrategy) Summarize(text string) (string, error) {
	prompt := fmt.Sprintf(`
	You are a Staff Engineer's assistant.
	Summarize the following technical article into 3 high-signal bullet points.
	Focus on architecture, trade-offs, and numbers. Ignore marketing fluff.
	
	Article Content:
	%s
	`, text)

	resp, err := g.model.GenerateContent(g.ctx, genai.Text(prompt))

	if err != nil {
		fmt.Printf("Error generating summary with Gemini: %v\n", err)
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", err
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			sb.WriteString(string(t))
		}
	}
	return sb.String(), nil
}

func (g *GeminiStrategy) Close() error {
	return g.client.Close()
}
