# Adler Review CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a robust, portably compiled Go command-line tool that takes a text file as input and generates an analytical review in Portuguese matching the `openspec/FORMATO.md` format using Google AI Pro's dynamic model selection.

**Architecture:** Modularized Go command-line layout with `cmd/adler-review/main.go` driving overall execution, `internal/config` handling chaves/variables, `internal/googleai` wrapper for client & dynamic model discovery, and `internal/reviewer` for prompt design and reading texts.

**Tech Stack:** Go (1.21+), `github.com/google/generative-ai-go/genai`, `google.golang.org/api/option`, `github.com/joho/godotenv`.

---

## File Structure

```
.
├── go.mod
├── go.sum
├── cmd/
│   └── adler-review/
│       └── main.go
└── internal/
    ├── config/
    │   └── config.go
    ├── googleai/
    │   └── client.go
    └── reviewer/
        └── reviewer.go
```

## Tasks

### Task 1: Go Module Initialization & Scaffolding

**Files:**
- Create: `go.mod`
- Create: `cmd/adler-review/main.go` (placeholder)

- [ ] **Step 1: Initialize the Go module**
Run: `go mod init adler-review-cli`
Expected: `go.mod` file is created with module name `adler-review-cli`.

- [ ] **Step 2: Install required packages**
Run:
```bash
go get github.com/google/generative-ai-go/genai
go get google.golang.org/api/option
go get github.com/joho/godotenv
```
Expected: `go.mod` and `go.sum` are updated with the installed dependencies.

- [ ] **Step 3: Create directory structure**
Run:
```bash
mkdir -p cmd/adler-review internal/config internal/googleai internal/reviewer
```
Expected: All directories are created.

- [ ] **Step 4: Create a minimal hello main.go**
Write to `cmd/adler-review/main.go`:
```go
package main

import "fmt"

func main() {
	fmt.Println("Adler Review CLI")
}
```

- [ ] **Step 5: Compile hello world**
Run: `go build -o adler-review cmd/adler-review/main.go`
Expected: Binário `adler-review` compilado e executa imprimindo "Adler Review CLI".

- [ ] **Step 6: Commit**
Run:
```bash
git add go.mod go.sum cmd/adler-review/main.go
git commit -m "feat: scaffold go project and module"
```

### Task 2: Config Loading & Key Verification

**Files:**
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`

- [ ] **Step 1: Write config_test.go**
Write to `internal/config/config_test.go`:
```go
package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("GEMINI_API_KEY", "test-key-123")
	defer os.Unsetenv("GEMINI_API_KEY")

	cfg, err := LoadConfig("flag-key-override")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey != "flag-key-override" {
		t.Errorf("expected flag key override to take precedence, got %s", cfg.APIKey)
	}

	cfg2, err := LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg2.APIKey != "test-key-123" {
		t.Errorf("expected env key fallback, got %s", cfg2.APIKey)
	}
}
```

- [ ] **Step 2: Run test and verify it fails**
Run: `go test ./internal/config/...`
Expected: FAIL due to "config" or "LoadConfig" undefined.

- [ ] **Step 3: Implement LoadConfig**
Write to `internal/config/config.go`:
```go
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey string
}

func LoadConfig(cliKey string) (*Config, error) {
	// Tenta carregar .env de forma silenciosa
	_ = godotenv.Load()

	apiKey := cliKey
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	if apiKey == "" {
		return nil, errors.New("chave de API do Google não encontrada. Defina a variável de ambiente GEMINI_API_KEY ou use o parâmetro -k")
	}

	return &Config{APIKey: apiKey}, nil
}
```

- [ ] **Step 4: Run test to verify it passes**
Run: `go test -v ./internal/config/...`
Expected: PASS

- [ ] **Step 5: Commit**
Run:
```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: implement config loading and key validation"
```

### Task 3: CLI Parsing & Validation

**Files:**
- Modify: `cmd/adler-review/main.go`

- [ ] **Step 1: Write parsing logic with flags in main.go**
Replace `cmd/adler-review/main.go` with:
```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"adler-review-cli/internal/config"
)

func main() {
	outputFlag := flag.String("output", "", "Caminho do arquivo de saída Markdown")
	flag.StringVar(outputFlag, "o", "", "Caminho do arquivo de saída Markdown (abreviação)")

	keyFlag := flag.String("key", "", "Chave de API do Google AI Pro")
	flag.StringVar(keyFlag, "k", "", "Chave de API do Google AI Pro (abreviação)")

	modelFlag := flag.String("model", "", "Sobrescrita manual do modelo")
	flag.StringVar(modelFlag, "m", "", "Sobrescrita manual do modelo (abreviação)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: adler-review-cli <caminho-do-arquivo-de-entrada> [flags]\n\nOpções:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Erro: arquivo de entrada não especificado.")
		flag.Usage()
		os.Exit(1)
	}

	inputFile := args[0]

	// Validação básica de arquivo
	info, err := os.Stat(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: O arquivo de entrada '%s' não existe ou não pôde ser acessado.\n", inputFile)
		os.Exit(1)
	}
	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "Erro: O caminho '%s' é um diretório, não um arquivo regular.\n", inputFile)
		os.Exit(1)
	}

	// Carrega chaves
	cfg, err := config.LoadConfig(*keyFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	// Determinar arquivo de saída padrão se omitido
	outputFile := *outputFlag
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		base := strings.TrimSuffix(filepath.Base(inputFile), ext)
		outputFile = base + "-resenha.md"
	}

	// Validação de permissão de escrita de saída prévia
	testFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: Sem permissão de escrita no caminho de saída '%s'.\n", outputFile)
		os.Exit(1)
	}
	testFile.Close()

	fmt.Printf("Configurado: Entrada=%s, Saída=%s, ChaveCarregada=%t, ModeloSobrescrita=%s\n", inputFile, outputFile, cfg.APIKey != "", *modelFlag)
}
```

- [ ] **Step 2: Test CLI locally for error handling**
Run: `go build -o adler-review cmd/adler-review/main.go`
Expected: Binário compilado sem erros.

- [ ] **Step 3: Run with missing positional argument**
Run: `./adler-review`
Expected: Saída: "Erro: arquivo de entrada não especificado." e encerra com exit code 1.

- [ ] **Step 4: Run with non-existent input file**
Run: `./adler-review non_existent_file.txt`
Expected: Saída: "Erro: O arquivo de entrada 'non_existent_file.txt' não existe..." e encerra com exit code 1.

- [ ] **Step 5: Run with valid input file but missing API Key**
Run:
```bash
touch temp_test.txt
unset GEMINI_API_KEY
./adler-review temp_test.txt
rm temp_test.txt
```
Expected: Saída: "Erro: chave de API do Google não encontrada..." e encerra com exit code 1.

- [ ] **Step 6: Commit**
Run:
```bash
git add cmd/adler-review/main.go
git commit -m "feat: implement CLI parsing, key validation and write check"
```

### Task 4: File Loading, Token Counter & Context Limitation

**Files:**
- Create: `internal/reviewer/reviewer.go`
- Test: `internal/reviewer/reviewer_test.go`

- [ ] **Step 1: Write reviewer_test.go**
Write to `internal/reviewer/reviewer_test.go`:
```go
package reviewer

import (
	"os"
	"testing"
)

func TestReadInputFile(t *testing.T) {
	content := "Este é um arquivo de teste com dez palavras para validação."
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, _ = tmpFile.WriteString(content)
	tmpFile.Close()

	text, words, tokens, err := ReadInputFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error reading file: %v", err)
	}
	if text != content {
		t.Errorf("expected text %q, got %q", content, text)
	}
	if words != 10 {
		t.Errorf("expected 10 words, got %d", words)
	}
	expectedTokens := int(float64(words) * 1.3)
	if tokens != expectedTokens {
		t.Errorf("expected %d tokens, got %d", expectedTokens, tokens)
	}
}
```

- [ ] **Step 2: Run test and verify it fails**
Run: `go test ./internal/reviewer/...`
Expected: FAIL due to `ReadInputFile` undefined.

- [ ] **Step 3: Implement ReadInputFile**
Write to `internal/reviewer/reviewer.go`:
```go
package reviewer

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func ReadInputFile(path string) (text string, words int, tokens int, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", 0, 0, err
	}

	if !utf8.Valid(b) {
		return "", 0, 0, fmt.Errorf("arquivo não contém texto UTF-8 válido")
	}

	text = string(b)
	fields := strings.Fields(text)
	words = len(fields)
	// Estimativa conservadora de tokens baseada em contagem de palavras (1 palavra ≈ 1.3 tokens)
	tokens = int(float64(words) * 1.3)

	return text, words, tokens, nil
}
```

- [ ] **Step 4: Run test to verify it passes**
Run: `go test -v ./internal/reviewer/...`
Expected: PASS

- [ ] **Step 5: Commit**
Run:
```bash
git add internal/reviewer/reviewer.go internal/reviewer/reviewer_test.go
git commit -m "feat: implement file reading and conservative token estimation"
```

### Task 5: Google Gemini Integration & Dynamic Model Selection

**Files:**
- Create: `internal/googleai/client.go`

- [ ] **Step 1: Write client.go code with ListModels selection and Exponential Backoff**
Write to `internal/googleai/client.go`:
```go
package googleai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type ClientWrapper struct {
	Client *genai.Client
}

func NewClient(ctx context.Context, apiKey string) (*ClientWrapper, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &ClientWrapper{Client: client}, nil
}

func (cw *ClientWrapper) SelectBestProModel(ctx context.Context) string {
	iter := cw.Client.ListModels(ctx)
	var availablePro []string

	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("[Aviso] Erro ao listar modelos do Google AI: %v. Usando fallback de preferência.", err)
			break
		}

		// Filtrar por modelos que contenham "pro" e suportem generateContent
		isPro := strings.Contains(strings.ToLower(m.Name), "pro") || strings.Contains(strings.ToLower(m.DisplayName), "pro")
		supportsGeneration := false
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				supportsGeneration = true
				break
			}
		}

		if isPro && supportsGeneration {
			availablePro = append(availablePro, m.Name)
		}
	}

	// Ordena por preferência estática se encontrar múltiplos modelos Pro
	for _, mName := range availablePro {
		// Modelos retornados normalmente começam com "models/"
		if strings.Contains(mName, "gemini-2.5-pro") {
			return mName
		}
	}
	for _, mName := range availablePro {
		if strings.Contains(mName, "gemini-1.5-pro") {
			return mName
		}
	}

	if len(availablePro) > 0 {
		return availablePro[0] // Retorna o primeiro Pro disponível
	}

	return "models/gemini-1.5-pro" // Fallback seguro
}

func (cw *ClientWrapper) GenerateWithRetry(ctx context.Context, modelName string, systemInstruction string, prompt string) (string, error) {
	model := cw.Client.GenerativeModel(modelName)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	backoffs := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second}
	var lastErr error

	for i := 0; i <= len(backoffs); i++ {
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err == nil {
			if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
				return "", fmt.Errorf("a API do Gemini retornou uma resposta sem conteúdo")
			}
			var builder strings.Builder
			for _, part := range resp.Candidates[0].Content.Parts {
				if textPart, ok := part.(genai.Text); ok {
					builder.WriteString(string(textPart))
				}
			}
			return builder.String(), nil
		}

		lastErr = err
		// Verifica se o erro é rate limit (HTTP 429 ou RESOURCE_EXHAUSTED)
		isRateLimit := strings.Contains(strings.ToLower(err.Error()), "429") || strings.Contains(strings.ToLower(err.Error()), "resource_exhausted")
		if !isRateLimit || i == len(backoffs) {
			break
		}

		log.Printf("[Rate Limit] API esgotada (429). Tentando novamente em %v (Tentativa %d de %d)...", backoffs[i], i+1, len(backoffs))
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(backoffs[i]):
		}
	}

	return "", fmt.Errorf("falha ao gerar conteúdo após retentativas: %w", lastErr)
}
```

- [ ] **Step 2: Test compilation of ClientWrapper**
Run: `go build -o /dev/null cmd/adler-review/main.go`
Expected: Build succeeds with zero errors.

- [ ] **Step 3: Commit**
Run:
```bash
git add internal/googleai/client.go
git commit -m "feat: implement client wrapper with dynamic model selection and exponential retry"
```

### Task 6: Prompt Engineering & Adler Generation Engine

**Files:**
- Modify: `internal/reviewer/reviewer.go`

- [ ] **Step 1: Implement prompt design and generation runner in reviewer.go**
Add the following functions to `internal/reviewer/reviewer.go`:
```go
package reviewer

import (
	"context"
	"fmt"
)

// SystemInstruction retorna as orientações estritas para o modelo baseadas em FORMATO.md
func SystemInstruction() string {
	return `Você é um Resenhista Analítico especialista, altamente treinado na metodologia de leitura analítica proposta por Mortimer Adler em "Como Ler Livros".
Sua tarefa é analisar rigorosamente o texto de entrada fornecido (que pode ser uma obra inteira, capítulos ou resumos detalhados) e gerar uma Resenha Analítica estruturada estritamente em português brasileiro.

A saída gerada deve conter EXATAMENTE as seguintes seções em Markdown, sem nenhuma alteração nos cabeçalhos:

## Abertura
[Escreva de 1 a 2 parágrafos contextualizando o autor do livro, o momento histórico da publicação e o lugar da obra na sua carreira ou no debate contemporâneo. NÃO antecipe as respostas das quatro perguntas fundamentais de Adler aqui]

## 1. Sobre o que é o livro?
[Defina de forma exata a classificação de gênero/categoria do livro (Teórico ou Prático; se Ficção ou Não-Ficção). Declare em 1 a 3 frases a tese central ou mensagem dominante. Descreva de forma concisa a unidade ou arquitetura geral da obra (início, meio e fim)]

## 2. O que está sendo dito?
[Descreva com profundidade a estrutura interna do livro, seus capítulos ou atos principais. Liste e analise as proposições ou termos-chave defendidos pelo autor e como os argumentos são progressivamente encadeados]

## 3. O livro está certo?
[Realize o julgamento crítico fundamentado. Lembre-se da regra de ouro de Adler: você só pode criticar após demonstrar que compreendeu perfeitamente. Avalie a coerência lógica e a verdade das teses. Aponte especificamente se o autor é:
1. Desinformado (faltam dados fundamentais)
2. Mal informado (dados errados)
3. Ilógico (contradições ou raciocínios falhos)
4. Incompleto em sua análise.
Se for literatura imaginativa, avalie a verossimilhança existencial, coerência interna do universo e profundidade estética]

## 4. Qual a importância?
[Descreva qual o impacto, relevância atual e conexões que esta obra estabelece hoje. Qual o significado prático para a vida do leitor ou para o campo conceitual da área]

REGRAS CRÍTICAS DE ESTILO:
- Escreva de forma profunda, intelectualizada, objetiva e clara. Evite resumos banais ou superficiais.
- Use os cabeçalhos em Markdown exatamente como descritos acima.
- Nunca adicione preâmbulos ou introduções fora dos cabeçalhos. Comece diretamente com "## Abertura".`
}

type AIClientInterface interface {
	SelectBestProModel(ctx context.Context) string
	GenerateWithRetry(ctx context.Context, modelName string, systemInstruction string, prompt string) (string, error)
}

func GenerateAdlerReview(ctx context.Context, client AIClientInterface, text string, modelOverride string) (string, string, error) {
	model := modelOverride
	if model == "" {
		model = client.SelectBestProModel(ctx)
	}

	prompt := fmt.Sprintf("Analise detalhadamente o seguinte texto do livro para gerar a Resenha Analítica:\n\n---\n%s\n---", text)

	result, err := client.GenerateWithRetry(ctx, model, SystemInstruction(), prompt)
	if err != nil {
		return "", "", err
	}

	return result, model, nil
}
```

- [ ] **Step 2: Verify compilation**
Run: `go build -o /dev/null cmd/adler-review/main.go`
Expected: Build succeeds with zero errors.

- [ ] **Step 3: Commit**
Run:
```bash
git add internal/reviewer/reviewer.go
git commit -m "feat: implement adler system instruction and reviewer execution engine"
```

### Task 7: Full End-to-End Integration

**Files:**
- Modify: `cmd/adler-review/main.go`

- [ ] **Step 1: Link config, reviewer, and googleai together in main.go**
Replace `cmd/adler-review/main.go` with the complete production code:
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"adler-review-cli/internal/config"
	"adler-review-cli/internal/googleai"
	"adler-review-cli/internal/reviewer"
)

func main() {
	outputFlag := flag.String("output", "", "Caminho do arquivo de saída Markdown")
	flag.StringVar(outputFlag, "o", "", "Caminho do arquivo de saída Markdown (abreviação)")

	keyFlag := flag.String("key", "", "Chave de API do Google AI Pro")
	flag.StringVar(keyFlag, "k", "", "Chave de API do Google AI Pro (abreviação)")

	modelFlag := flag.String("model", "", "Sobrescrita manual do modelo")
	flag.StringVar(modelFlag, "m", "", "Sobrescrita manual do modelo (abreviação)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: adler-review-cli <caminho-do-arquivo-de-entrada> [flags]\n\nOpções:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Erro: arquivo de entrada não especificado.")
		flag.Usage()
		os.Exit(1)
	}

	inputFile := args[0]

	// 1. Validar e ler arquivo de entrada
	info, err := os.Stat(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: O arquivo de entrada '%s' não existe ou não pôde ser acessado.\n", inputFile)
		os.Exit(1)
	}
	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "Erro: O caminho '%s' é um diretório, não um arquivo regular.\n", inputFile)
		os.Exit(1)
	}

	text, words, estimatedTokens, err := reviewer.ReadInputFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler o arquivo de entrada: %v\n", err)
		os.Exit(1)
	}

	// 2. Determinar caminho de saída e testar permissão de gravação
	outputFile := *outputFlag
	hasCustomOutput := outputFile != ""
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		base := strings.TrimSuffix(filepath.Base(inputFile), ext)
		outputFile = base + "-resenha.md"
	}

	testFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: Sem permissão de escrita no caminho de saída '%s'.\n", outputFile)
		os.Exit(1)
	}
	testFile.Close()

	// 3. Carregar Chave de API e Inicializar Cliente
	cfg, err := config.LoadConfig(*keyFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	// Configurar contexto de execução com Timeout Robusto de 5 minutos
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := googleai.NewClient(ctx, cfg.APIKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao conectar-se à API do Google AI: %v\n", err)
		os.Exit(1)
	}
	defer client.Client.Close()

	// 4. Seleção de modelo e alerta de contexto
	selectedModel := *modelFlag
	if selectedModel == "" {
		selectedModel = client.SelectBestProModel(ctx)
	}

	fmt.Printf("[Info] Processando '%s' (%d palavras, aprox. %d tokens estimado).\n", inputFile, words, estimatedTokens)
	if estimatedTokens > 900000 {
		fmt.Println("[Aviso] O arquivo de entrada é extremamente longo e pode estar próximo ou exceder a janela ideal do modelo.")
	}

	fmt.Printf("[Info] Usando o modelo Pro: %s\n", selectedModel)
	fmt.Println("[Info] Enviando solicitação para o Google AI Pro...")

	// 5. Geração da Resenha Analítica
	reviewContent, finalModel, err := reviewer.GenerateAdlerReview(ctx, client, text, *modelFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro na geração da resenha: %v\n", err)
		os.Exit(1)
	}

	// 6. Gravar Arquivo de Saída
	err = os.WriteFile(outputFile, []byte(reviewContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar arquivo de saída '%s': %v\n", outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("[Sucesso] Resenha analítica gerada com sucesso e salva em: %s (Modelo: %s)\n", outputFile, finalModel)

	// Se não foi fornecido arquivo de saída customizado, imprime o resultado no stdout
	if !hasCustomOutput {
		fmt.Println("\n--- RESENHA GERADA (STDOUT) ---")
		fmt.Println(reviewContent)
		fmt.Println("-------------------------------")
	}
}
```

- [ ] **Step 2: Build the final production application**
Run: `go build -o adler-review cmd/adler-review/main.go`
Expected: Binário `adler-review` compilado de forma limpa e sem erros.

- [ ] **Step 3: Run all unit tests**
Run: `go test -v ./...`
Expected: All tests PASS.

- [ ] **Step 4: Commit**
Run:
```bash
git add cmd/adler-review/main.go
git commit -m "feat: complete full CLI end-to-end integration"
```
