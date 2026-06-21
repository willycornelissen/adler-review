## Context

A geração de resenhas analíticas de alta qualidade requer a leitura e análise aprofundada de textos que podem variar de pequenos artigos a livros inteiros. A metodologia de Mortimer Adler (definida em `openspec/FORMATO.md`) exige um processo estruturado em 5 fases (Abertura + 4 Perguntas).
Para atender aos requisitos de processamento de textos longos e de seleção inteligente de modelos da assinatura do Google AI Pro, este projeto propõe uma ferramenta CLI moderna baseada em Go (Golang) que consome a API do Gemini Pro, realizando a leitura integral de obras de forma nativa devido à ampla janela de contexto da plataforma.

## Goals / Non-Goals

**Goals:**
- Prover uma ferramenta CLI robusta e fácil de usar que receba um arquivo de texto e gere a resenha analítica formatada em Markdown.
- Integrar com o SDK oficial da API do Google Gen AI para Go (`github.com/google/generative-ai-go/genai`).
- Descobrir de forma autônoma e dinâmica o melhor modelo de nível Pro disponível na assinatura do usuário (ex: `gemini-2.5-pro`, `gemini-1.5-pro`).
- Formatar o arquivo de saída rigorosamente de acordo com as especificações e o exemplo prático descritos em `openspec/FORMATO.md`.
- Lidar com arquivos de texto de tamanhos variados usando tratamento de erros apropriado.

**Non-Goals:**
- Conversão de arquivos de entrada complexos (como PDFs protegidos por DRM, EPUBs com criptografia ou imagens/scans via OCR). O escopo limita-se a arquivos de texto puro (`.txt` ou `.md`).
- Interface gráfica de usuário (GUI) ou aplicação web (focado apenas em CLI).
- Geração de resenhas em lote (processamento concorrente de múltiplos livros em uma única execução).

## Architecture & Directory Structure

A aplicação seguirá a estrutura de pacotes recomendada para projetos em Go:

```
adler-review-cli/
├── go.mod                  # Definição do módulo Go e dependências
├── go.sum                  # Verificação de integridade das dependências
├── cmd/
│   └── adler-review/       # Ponto de entrada do executável da CLI
│       └── main.go         # Inicialização, parsing de flags e tratamento global de erros
└── internal/
    ├── config/             # Gerenciamento de configurações e variáveis de ambiente
    │   └── config.go
    ├── googleai/           # Wrapper do cliente do Google Gen AI e seleção de modelos
    │   └── client.go
    └── reviewer/           # Core de leitura, cálculo de tokens e prompt de resenha de Adler
        └── reviewer.go
```

## Decisions

### 1. Linguagem de Programação e Runtime: Go (Golang)
- **Escolha:** Go (v1.21+).
- **Alternativa Considerada:** Python, Node.js.
- **Razão:** Go é uma linguagem de programação compilada, altamente performática e estaticamente tipada, perfeita para a criação de binários CLI rápidos, portáteis e sem dependências externas de runtime (como Node.js ou interpretador Python).

### 2. SDK de IA: `github.com/google/generative-ai-go/genai` (SDK oficial do Google para Go)
- **Escolha:** O pacote `github.com/google/generative-ai-go/genai`.
- **Alternativa Considerada:** Chamadas HTTP diretas via REST.
- **Razão:** O SDK oficial do Google encapsula a lógica de conexão com o Gemini de forma Go-idiomática e simplifica tarefas como listagem de modelos, controle de parâmetros e geração de conteúdo.

### 3. Seleção Dinâmica do Melhor Modelo Pro
- **Escolha:** A aplicação criará o cliente `genai.NewClient` e executará uma iteração sobre `client.ListModels(ctx)`. Filtrará os modelos cujo campo `Name` ou `DisplayName` contenham "pro" e que suportem a ação `generateContent`. Os modelos Pro serão ordenados de forma descendente por capacidade (preferindo `gemini-2.5-pro`, seguido de `gemini-1.5-pro`). Caso a consulta falhe ou retorne lista vazia, usará o fallback estático `gemini-1.5-pro` para máxima resiliência.
- **Alternativa Considerada:** Hardcode fixo do modelo.
- **Razão:** Garante que o aplicativo use sempre a melhor inteligência disponível na assinatura do Google AI Pro do usuário, mesmo quando novos modelos Pro forem lançados.

### 4. Controle de Conexão e Timeout via Context
- **Escolha:** Toda requisição à API do Google AI Pro usará um `context.WithTimeout` com limite de 5 minutos, repassando o contexto para o SDK.
- **Alternativa Considerada:** Timeout indefinido.
- **Razão:** O processamento de livros ou textos muito extensos pode demorar devido ao tempo de geração (que pode produzir milhares de palavras). Um timeout de 5 minutos previne que o CLI trave indefinidamente e gerencie a liberação de recursos de rede de forma eficiente.

### 5. Tratamento de Grandes Contextos
- **Escolha:** Passar o conteúdo total do texto diretamente no prompt da API, uma vez que a família Gemini Pro possui janelas de contexto massivas (1 milhão a 2 milhões de tokens). Se o arquivo exceder um limite de segurança prática (ex: 1,5M de tokens, aproximadamente 1,1M de palavras), um aviso será exibido ao usuário, e o aplicativo tentará truncar ou falhar graciosamente.
- **Alternativa Considerada:** Chunking de texto e Map-Reduce.
- **Razão:** O processamento em contexto único (single-context prompt) permite ao modelo compreender melhor o livro como um todo e relacionar diferentes capítulos de maneira muito mais profunda do que em chunks isolados.

## Risks / Trade-offs

- **[Risco]** Tamanho do Texto Excede a Janela de Contexto ou Limites da API → **[Mitigação]** O CLI calculará o número estimado de tokens do arquivo de entrada e, se ultrapassar o limite suportado pelo modelo, alertará o usuário antes de iniciar a requisição custosa.
- **[Risco]** API Rate Limits (Limites de Taxa por Minuto/Dia) → **[Mitigação]** Tratamento de erros de limite de taxa (`429 Too Many Requests`) com retentativa exponencial (exponential backoff).
- **[Risco]** Falta de Chave de API no Ambiente → **[Mitigação]** Validação imediata na inicialização do CLI, exibindo uma mensagem amigável instruindo o usuário a definir a variável `GEMINI_API_KEY`.
