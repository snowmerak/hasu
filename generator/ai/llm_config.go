package ai

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
	"gopkg.in/yaml.v3"
)

const (
	configFileName = ".ako/llm.config.yaml"
)

type Config struct {
	Ollama struct {
		Enable bool   `yaml:"enable"`
		Host   string `yaml:"host"`
		Model  string `yaml:"model"`
	} `yaml:"ollama"`
	Gemini struct {
		Enable bool   `yaml:"enable"`
		Model  string `yaml:"model"`
		APIKey string `yaml:"api_key"`
	} `yaml:"gemini"`
	Vertex struct {
		Enable   bool   `yaml:"enable"`
		Model    string `yaml:"model"`
		APIKey   string `yaml:"api_key"`
		Location string `yaml:"location"`
		Project  string `yaml:"project"`
	}
	Anthropic struct {
		Enable bool   `yaml:"enable"`
		Model  string `yaml:"model"`
		APIKey string `yaml:"api_key"`
	} `yaml:"anthropic"`
	OpenAI struct {
		Enable bool   `yaml:"enable"`
		Model  string `yaml:"model"`
		APIKey string `yaml:"api_key"`
	} `yaml:"openai"`
}

var GlobalConfig Config

var globalConfigNotExists bool

func init() {
	if err := ReadConfig(); err != nil {
		globalConfigNotExists = true
	}
}

func ReadConfig() error {
	configFile, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&GlobalConfig); err != nil {
		return err
	}

	return nil
}

func SaveConfig() error {
	if err := os.MkdirAll(".ako", os.ModePerm); err != nil {
		return err
	}

	configFile, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := yaml.NewEncoder(configFile)
	if err := encoder.Encode(GlobalConfig); err != nil {
		return err
	}

	return nil
}

func InitConfig() error {
	GlobalConfig = Config{
		Ollama: struct {
			Enable bool   `yaml:"enable"`
			Host   string `yaml:"host"`
			Model  string `yaml:"model"`
		}{
			Enable: true,
			Host:   "http://localhost:11434",
			Model:  "gemma3:1b",
		},
	}

	if err := SaveConfig(); err != nil {
		return err
	}

	return nil
}

type LLMClient interface {
	GenerateCommitMessage(ctx context.Context, gitDiff string) (<-chan string, error)
}

func NewLLMClient(ctx context.Context) (LLMClient, error) {
	if globalConfigNotExists {
		return nil, fmt.Errorf("config file not exists")
	}

	switch {
	case GlobalConfig.Ollama.Enable:
		client, err := NewOllamaClient(GlobalConfig.Ollama.Host, GlobalConfig.Ollama.Model)
		if err != nil {
			return nil, err
		}
		return client, nil
	case GlobalConfig.Gemini.Enable:
		client, err := NewGeminiClient(genai.BackendGeminiAPI, GlobalConfig.Gemini.APIKey, GlobalConfig.Gemini.Model, "", "")
		if err != nil {
			return nil, err
		}
		return client, nil
	case GlobalConfig.Vertex.Enable:
		client, err := NewGeminiClient(genai.BackendVertexAI, GlobalConfig.Vertex.APIKey, GlobalConfig.Vertex.Model, GlobalConfig.Vertex.Location, GlobalConfig.Vertex.Project)
		if err != nil {
			return nil, err
		}
		return client, nil
	case GlobalConfig.Anthropic.Enable:
		client, err := NewAnthropicClient(GlobalConfig.Anthropic.APIKey, GlobalConfig.Anthropic.Model)
		if err != nil {
			return nil, err
		}
		return client, nil
	case GlobalConfig.OpenAI.Enable:
		client, err := NewOpenAIClient(GlobalConfig.OpenAI.APIKey, GlobalConfig.OpenAI.Model)
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	return nil, fmt.Errorf("no LLM client enabled")
}

func GenerateCommitMessage(ctx context.Context, gitDiff string) (<-chan string, error) {
	client, err := NewLLMClient(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := client.GenerateCommitMessage(ctx, gitDiff)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
