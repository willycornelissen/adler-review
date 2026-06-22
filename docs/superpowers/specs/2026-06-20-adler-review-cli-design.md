# Design Spec: Adler Review CLI (Go)

- **Data:** 20 de Junho de 2026
- **Autor:** OpenCode CLI Agent
- **Status:** Aprovado / Pronto para Implementação

---

## 1. Contexto e Objetivo

A geração manual de resenhas analíticas detalhadas baseadas no rigor intelectual de Mortimer Adler (autor de *Como Ler Livros*, 1940) é um processo altamente complexo e demorado.
Este projeto define o design de uma ferramenta de linha de comando (CLI) desenvolvida em **Go (Golang)** que automatiza a geração inicial destas resenhas, estruturando-as rigorosamente conforme o formato oficial descrito em `openspec/FORMATO.md`.
A ferramenta consome a API do Google AI Pro (Gemini) de forma altamente eficiente, aproveitando a janela de contexto ultra-longo do Gemini para processar livros ou manuscritos completos em uma única chamada contextual.

---

## 2. Metas e Não-Metas

### Metas (Goals):
- Binário compilado rápido, portável e independente de runtimes externos.
- Detecção e seleção dinâmica do melhor modelo Pro disponível (Gemini 2.5 Pro > Gemini 1.5 Pro).
- Envio integrado do contexto total do livro ou texto de entrada diretamente no prompt da API.
- Saída formatada em Markdown estritamente aderente ao formato exigido (`openspec/FORMATO.md`).
- Se o arquivo de saída não for especificado, salvar no caminho padrão e também exibir o resultado no terminal (`stdout`).
- Tratamento explícito de rate limits (HTTP 429) usando recuo exponencial.

### Não-Metas (Non-Goals):
- Interface gráfica (GUI) ou painel web.
- Processamento em lote de múltiplos arquivos de entrada simultâneos.
- Conversão ou OCR de arquivos binários complexos (PDFs protegidos, EPUBs, imagens). Foco exclusivo em UTF-8 texto puro (`.txt`, `.md`).

---

## 3. Arquitetura e Estrutura de Pacotes

A aplicação Go será estruturada para separação rígida de responsabilidades de forma idiomática e modular:

```
adler-review-cli/
├── go.mod                  # Inicialização do módulo Go (adler-review-cli)
├── go.sum                  # Verificação de integridade de dependências
├── cmd/
│   └── adler-review/       # Executável CLI principal
│       └── main.go         # Inicialização, parsing de flags, validação de escrita e orquestração do fluxo
└── internal/
    ├── config/             # Gerenciamento de chaves (GEMINI_API_KEY) e envs (.env)
    │   └── config.go
    ├── googleai/           # Wrapper do SDK do Gemini, listagem e seleção de modelo Pro
    │   └── client.go
    └── reviewer/           # Core: leitura do texto, contagem de tokens e prompts Adler
        └── reviewer.go
```

---

## 4. Detalhamento de Componentes

### 4.1 CLI e Orquestração (`cmd/adler-review/main.go`)
- **Flags Suportadas:**
  - `-o`, `--output`: Caminho para salvar a resenha (Default: `<nome-de-entrada>-resenha.md`).
  - `-k`, `--key`: Chave da API do Google AI (Fallback para `GEMINI_API_KEY` do ambiente).
  - `-m`, `--model`: Sobrescrita manual do modelo para evitar seleção automática.
  - `-h`, `--help`: Mensagem informativa de ajuda.
- **Fluxo de Inicialização:**
  1. Carrega variáveis de ambiente via `godotenv`.
  2. Valida se a chave de API está presente.
  3. Valida a acessibilidade do arquivo de entrada e permissão de gravação no arquivo de saída.
  4. Dispara a lógica de leitura e validação de tokens.
  5. Invoca a API do Gemini Pro com timeout estrito de **5 minutos**.
  6. Grava o Markdown resultante no arquivo de saída e o imprime no terminal (`stdout`) se o caminho de saída tiver sido omitido.

### 4.2 Configurações e Ambiente (`internal/config/config.go`)
- Responsável por carregar o arquivo `.env` usando `github.com/joho/godotenv` se ele estiver presente na pasta atual.
- Exporta uma estrutura `Config` contendo o `APIKey` e o `ModelOverride`.
- Consolida a priorização da chave de API: Flag CLI > Variável de Ambiente `GEMINI_API_KEY` > Falha.

### 4.3 Cliente do Google Gemini (`internal/googleai/client.go`)
- Usa o SDK oficial unificado `github.com/google/generative-ai-go/genai` e `google.golang.org/api/option`.
- Implementa a função de seleção dinâmica de modelo:
  - Executa `client.ListModels(ctx)` para buscar todos os modelos da assinatura.
  - Filtra por modelos ativos com `"pro"` no ID e compatibilidade com `generateContent`.
  - Ordena a prioridade: `gemini-2.5-pro` > `gemini-1.5-pro`.
  - Em caso de falha de chamada de listagem, emite um aviso ao terminal e assume `gemini-1.5-pro` como fallback seguro.
- **Tratamento de Rate Limit (HTTP 429):**
  - Loop de retentativa exponencial (max 3 tentativas) para contornar limites transitórios de concorrência de requisição.

### 4.4 Gerador de Resenhas Analíticas (`internal/reviewer/reviewer.go`)
- **Leitura do Livro:** Carrega o conteúdo textual completo para a memória (`os.ReadFile`).
- **Validação de Tamanho:** Utiliza uma estimativa conservadora de tokens baseada em contagem de palavras (`Tokens = Palavras * 1.3`). Se exceder 90% da janela do modelo selecionado (ex: 900K para 1M de tokens), emite um aviso no terminal informando que o texto pode ser truncado.
- **System Prompt (Engenharia de Prompt em Português):**
  Incorpora o guia metodológico de Mortimer Adler baseado no `openspec/FORMATO.md`:
  - **Abertura:** Contexto histórico, autor, momento da carreira.
  - **Pergunta 1:** Classificação do livro, tese em 1-3 frases, arquitetura global.
  - **Pergunta 2:** Estruturação interna, termos-chave, condução do argumento.
  - **Pergunta 3:** Julgamento imparcial (desinformado, mal informado, ilógico, incompleto).
  - **Pergunta 4:** Importância e conexões sindônticas atuais.
  - O prompt exige cabeçalhos exatos em Markdown (`## Abertura`, `## 1. Sobre o que é o livro?`, etc.).

---

## 5. Tratamento de Erros e Saídas (Exit Codes)

| Caso de Erro | Comportamento da CLI | Exit Code |
|---|---|---|
| Chave de API ausente | Solicita a chave de forma interativa no terminal. Se continuar vazia, imprime erro e sai. | 1 |
| Arquivo de entrada inacessível | Imprime erro no `stderr` notificando arquivo inexistente/ilegível. | 1 |
| Pasta de saída sem permissão | Imprime erro no `stderr` antes de chamar a API, evitando custos. | 1 |
| Erro de Rate Limit persistente | Imprime erro após 3 tentativas de recuo exponencial. | 1 |
| Tempo de API Excedido (5m) | Imprime erro de timeout de rede e orienta tentar novamente. | 1 |
| Execução com Sucesso | Grava arquivo Markdown e imprime o texto gerado na stdout (se omitida saída). | 0 |
