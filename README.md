# Adler Review CLI 📚✨

**Adler Review CLI** é uma ferramenta de linha de comando escrita em Go que automatiza a criação de resenhas literárias analíticas e profundas baseadas estritamente na renomada metodologia de Mortimer Adler exposta em seu clássico livro *"Como Ler Livros"* (1940). 

Utilizando a API do **Google GenAI (Gemini)**, o utilitário identifica com precisão a obra desejada, confirma os metadados bibliográficos interativamente com o usuário e gera um documento estruturado em Markdown com alto rigor analítico e vocabulário erudito.

---

## 🧠 A Metodologia Mortimer Adler

A resenha gerada segue rigorosamente a estrutura de leitura analítica proposta por Adler, dividida nas quatro perguntas fundamentais que todo leitor ativo deve fazer ao ler um livro:

1. **Abertura:** Contextualização histórica e biográfica sobre o autor e a obra, ancorando o leitor no mundo real sem antecipar as respostas das seções seguintes.
2. **Sobre o que é o livro? (Classificação e Visão Panorâmica):** Define o gênero e a categoria da obra (Teórico ou Prático), a tese central (1-3 frases) e a sua unidade estrutural.
3. **O que está sendo dito? (Estrutura e Proposições):** Mapeamento neutro da organização interna da obra (capítulos, partes, progressão narrativa) e suas proposições/termos-chave.
4. **O livro está certo? (Julgamento Crítico Fundamentado):** Uma avaliação intelectualmente honesta que julga a obra com base em quatro critérios objetivos de Adler:
   - Onde o autor é desinformado (falta de dados).
   - Onde o autor está mal informado (dados incorretos).
   - Onde o autor é ilógico (contradições ou raciocínio falho).
   - Onde a análise do autor é incompleta.
5. **Qual a importância? (Significado e Leitura Sintópica):** Conecta a obra com o resto do mundo, abordando sua relevância, utilidade prática ou estética hoje e conexões com outros livros do mesmo tema.

---

## 🚀 Funcionalidades

- **Fluxo Interativo de Confirmação:** Busca e valida detalhes bibliográficos (título oficial, autor, ano de publicação original, gênero e breve sinopse) com o Google AI antes de iniciar a escrita da resenha.
- **Estruturação Markdown Impecável:** Saída estruturada de forma limpa, ideal para publicação em blogs, repositórios de notas pessoais (Obsidian, Notion) ou publicação acadêmica.
- **Seleção Inteligente de Modelo:** Descobre e seleciona automaticamente os melhores modelos Pro disponíveis do Gemini (como o `gemini-1.5-pro` ou `gemini-2.5-pro`) para garantir uma análise literária densa e criativa.
- **Prevenção de Erros de Escrita:** Realiza testes de gravação no arquivo de saída antes de invocar APIs externas de longa duração, evitando a perda de créditos ou tempo em caso de caminhos inválidos ou falhas de permissão.

---

## 🛠️ Instalação e Compilação

### Pré-requisitos
- **Go 1.25** ou superior instalado em seu sistema.
- Uma chave de API da **Google AI (Gemini)**. Você pode obter uma gratuitamente no [Google AI Studio](https://aistudio.google.com/).

### Compilação do Código Fonte

1. Clone o repositório ou navegue até o diretório do projeto:
   ```bash
   git clone git@github.com:willycornelissen/adler-review.git
   cd adler-review
   ```

2. Compile o binário do CLI:
   ```bash
   cd adler-review-cli
   go build -o adler-review ./cmd/adler-review
   ```

3. Mova o executável gerado para o seu PATH global (opcional) ou execute-o localmente:
   ```bash
   ./adler-review --help
   ```

---

## 📖 Como Usar

### Configuração da Chave da API

A forma recomendada de fornecer a chave do Gemini é por meio da variável de ambiente `GEMINI_API_KEY`:

```bash
export GEMINI_API_KEY="sua-chave-api-aqui"
```

Alternativamente, você pode fornecer a chave diretamente via flag ao executar o utilitário.

---

### Executando a Ferramenta

Para iniciar o fluxo interativo de criação de uma resenha:

```bash
./adler-review
```

#### Opções de Linha de Comando (Flags)

| Flag (Longa) | Flag (Curta) | Tipo | Descrição |
|---|---|---|---|
| `--output` | `-o` | string | Caminho personalizado para salvar o arquivo final (Default: `<slug-do-titulo>.md`) |
| `--key` | `-k` | string | Sobrescreve/fornece a Chave de API do Google AI diretamente por comando |
| `--model` | `-m` | string | Sobrescrita manual para usar um modelo Gemini específico (ex: `gemini-1.5-pro`) |
| `--help` | `-h` | booleano | Exibe o guia de ajuda e lista de comandos |

#### Exemplo de Uso com Flags:

```bash
./adler-review -o resenha_cristianismo_puro_e_simples.md -m gemini-1.5-pro
```

---

### Exemplo de Fluxo Interativo

```text
--- Definição do Livro para Resenha ---
Digite o título do livro: Cristianismo Puro e Simples
Digite o autor do livro (opcional, pressione Enter para pular): C.S. Lewis

Identificando o livro com o Google AI...

--- Livro Identificado ---
- **Título:** Cristianismo Puro e Simples
- **Autor:** C.S. Lewis
- **Ano de Publicação:** 1952
- **Gênero:** Apologética Cristã / Teologia / Filosofia da Religião
- **Breve Descrição:** Uma defesa e explicação das crenças cristãs fundamentais baseada em palestras de rádio ministradas pelo autor durante a Segunda Guerra Mundial.
--------------------------

Confirma este livro para a geração da resenha? (S/N): s

Model selected: gemini-1.5-pro
Processing review... (this may take up to 5 minutes)

Success! Review generated and saved to: cristianismo-puro-e-simples.md
```

---

## 📁 Estrutura do Repositório

O projeto está organizado da seguinte forma:

```text
├── adler-review-cli/            # Código-fonte principal do CLI em Go
│   ├── cmd/
│   │   └── adler-review/        # Ponto de entrada do programa (main.go)
│   ├── internal/
│   │   ├── config/              # Carregamento e validação de configurações e flags
│   │   ├── googleai/            # Wrapper de cliente para a API do Google AI (Gemini)
│   │   └── reviewer/            # Motores de prompt, formatação de resenhas e lógica de Adler
│   ├── go.mod                   # Arquivo de módulo Go
│   └── go.sum                   # Somas de verificação de dependências Go
├── docs/                        # Documentações adicionais, planos e especificações de design
├── openspec/                    # Especificações técnicas e formatos do ecossistema OpenSpec
└── README.md                    # Este arquivo guia do projeto
```

---

## 📝 Licença

Este projeto está licenciado sob a Licença MIT. Consulte o arquivo `LICENSE` para obter mais informações (se disponível).

---
*Desenvolvido com carinho para elevar o nível das discussões e leituras intelectuais.* 🧠✍️
