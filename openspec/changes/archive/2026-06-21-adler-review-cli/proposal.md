## Why

Manualmente produzir resenhas analíticas de alta qualidade seguindo a metodologia rigorosa de Mortimer Adler (como detalhado em `openspec/FORMATO.md`) exige um esforço cognitivo e de tempo substanciais. Com o avanço de modelos de linguagem de contexto ultra-longo do Google AI Pro (como a família Gemini 1.5/2.5 Pro), tornou-se viável automatizar a geração inicial dessas resenhas com profundidade, precisão e aderência estrita à estrutura exigida, escolhendo dinamicamente o melhor modelo disponível.

## What Changes

Criação de um utilitário de linha de comando (CLI) que:
- Recebe como entrada um arquivo de texto contendo a obra, capítulos ou notas detalhadas.
- Consome a API do Google AI Pro usando a chave fornecida por variáveis de ambiente ou arquivo de configuração.
- Seleciona automaticamente o melhor modelo de nível Pro disponível (como `gemini-1.5-pro` ou superior) para processar o texto, aproveitando a janela de contexto massiva de mais de 1 milhão de tokens para ler obras extensas de forma nativa e íntegra.
- Executa a análise analítica estruturada em 5 partes conforme o `openspec/FORMATO.md` (Abertura, P1: Sobre o que é o livro, P2: O que está sendo dito, P3: O livro está certo, P4: Qual a importância).
- Salva o resultado final em um arquivo Markdown estruturado e formatado de forma impecável em português.

## Capabilities

### New Capabilities
- `cli-interface`: Interface de linha de comando para receber caminhos de arquivos de entrada/saída, parâmetros de modelo e tratamento de erros de uso.
- `google-ai-client`: Integração com a API do Google Gen AI, autenticação via chave de API (`GEMINI_API_KEY`), listagem e seleção dinâmica do melhor modelo Pro disponível.
- `adler-review-generator`: Orquestração do processamento do texto, montagem de prompts otimizados e formatação estruturada do resultado final em estrito alinhamento com `openspec/FORMATO.md`.

### Modified Capabilities
Nenhuma.

## Impact

- **Novas Dependências:** Adição do SDK oficial do Google Gen AI para Go (`github.com/google/generative-ai-go/genai`).
- **Estrutura de Arquivos:** Novo diretório de código-fonte para o aplicativo CLI.
- **Configuração:** Necessidade de configuração da chave de API `GEMINI_API_KEY`.
