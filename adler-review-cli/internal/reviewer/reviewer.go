package reviewer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"adler-review-cli/internal/googleai"
)

// SystemInstruction contains Adler's 5-part analytical reading framework in Portuguese.
const SystemInstruction = `Você é um resenhista literário e analista intelectual altamente rigoroso, especializado na metodologia de leitura analítica desenvolvida por Mortimer Adler em seu livro "Como Ler Livros" (1940).
Sua tarefa é analisar o livro especificado (título, autor e detalhes identificados) e produzir uma resenha analítica profunda, precisa e estruturada em português brasileiro.

A estrutura do seu documento de saída deve seguir RIGOROSAMENTE o formato definido abaixo, usando cabeçalhos exatos em Markdown. Não adicione textos explicativos ou preâmbulos no início ou fim. Comece diretamente com o cabeçalho "## Abertura" e termine com o fim da seção "## 4. Qual a importância?".

### Estrutura Exigida:

## Abertura
Apresente um parágrafo contextualizador sobre a obra e o autor. O objetivo é ancorar o leitor no mundo real, oferecendo contexto externo.
- Inclua: Quem é o autor, o momento histórico da publicação, o lugar do livro em sua carreira ou o que motiva esta resenha hoje.
- Regra de ouro: Não antecipe as respostas das quatro perguntas fundamentais de Adler nesta seção.

## 1. Sobre o que é o livro?
Esta seção foca na classificação e na visão panorâmica do livro. Deve responder neutra e satisfatoriamente a:
- Classificação: Defina o gênero e a categoria da obra (Distinguir entre Teórico — ciência, filosofia, história, etc. — ou Prático — como-fazer, guias, ética/política aplicada, ou Literatura Imaginativa — romance, novela, conto, poesia).
- Tese central: Declare em 1-3 frases a mensagem principal ou o tema dominante.
- Unidade: Descreva resumidamente a arquitetura geral da obra (o início, o meio e o fim ou a premissa fundamental).

## 2. O que está sendo dito?
Nesta seção, demonstre que compreendeu profundamente a estrutura do livro de forma neutra. Deve responder a:
- Estrutura interna: Como o livro é organizado (partes, capítulos, atos ou saltos cronológicos).
- Proposições ou termos-chave: Quais as principais ideias que o autor defende e como as sustenta.
- Condução: O encadeamento do argumento ou a progressão da trama/narrativa.

## 3. O livro está certo?
Esta seção apresenta o julgamento crítico fundamentado. Adler adverte: não julgue sem antes compreender.
- Regra de ouro: Você só pode criticar o livro de forma fundamentada e imparcial. Nunca faça julgamentos puramente subjetivos.
- Formas de crítica:
  1. Mostre onde o autor é desinformado (faltam dados importantes).
  2. Mostre onde o autor está mal informado (dados errados).
  3. Mostre onde o autor é ilógico (contradição ou raciocínio falho).
  4. Mostre onde a análise do autor é incompleta.
  Se o livro estiver totalmente correto, fundamente por que os argumentos e dados apresentados sustentam a tese de forma irrefutável.
- Para literatura imaginativa: Avalie a coerência interna do universo, a verossimilhança dos personagens e a profundidade de sua verdade existencial/estética.

## 4. Qual a importância?
Conecte a obra com o restante do mundo.
- Fatores de valor: O significado do livro para o seu gênero, as conexões com outras obras (leitura sintópica) e sua utilidade teórica, prática ou estética hoje.

Adote um tom sóbrio, formal, intelectualmente denso e de alto nível de vocabulário, similar ao exemplo prático fornecido nas instruções de Adler.`

// ConfirmBookInstruction guides Gemini in identifying and formatting the book details.
const ConfirmBookInstruction = `Você é um assistente bibliográfico preciso. Dado o título de um livro e opcionalmente o autor fornecidos pelo usuário, sua tarefa é identificar com precisão o livro correspondente e retornar suas informações estruturadas em português brasileiro.
Se as informações fornecidas forem ambíguas ou insuficientes, tente identificar a obra mais famosa que corresponda aos termos.

Retorne a resposta EXATAMENTE no seguinte formato (substitua os valores entre colchetes):

- **Título:** [Título Oficial do Livro]
- **Autor:** [Nome Completo do Autor]
- **Ano de Publicação:** [Ano da primeira publicação]
- **Gênero:** [Gênero literário]
- **Breve Descrição:** [Uma frase resumindo a premissa do livro]`

// ConfirmBookDetails identifies and formats book details using Gemini.
func ConfirmBookDetails(ctx context.Context, client *googleai.Client, modelName string, title, author string) (string, error) {
	genCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	prompt := fmt.Sprintf("Identifique o livro:\nTítulo: %s\nAutor: %s", title, author)
	details, err := client.GenerateContentWithRetry(genCtx, modelName, ConfirmBookInstruction, prompt)
	if err != nil {
		return "", fmt.Errorf("falha ao identificar detalhes do livro: %w", err)
	}

	return details, nil
}

// GenerateReviewForBook manages context timeouts and triggers content generation for a specific book.
func GenerateReviewForBook(ctx context.Context, client *googleai.Client, modelName string, bookDetails string) (string, error) {
	genCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	prompt := fmt.Sprintf("Gere uma resenha analítica completa baseada na metodologia de Mortimer Adler para o seguinte livro identificado:\n\n%s", bookDetails)
	review, err := client.GenerateContentWithRetry(genCtx, modelName, SystemInstruction, prompt)
	if err != nil {
		return "", err
	}

	return review, nil
}

// Helper function to slugify book title for default output filename.
func Slugify(s string) string {
	s = strings.ToLower(s)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}
	res := result.String()
	// Replace multiple consecutive hyphens
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	return strings.Trim(res, "-")
}
