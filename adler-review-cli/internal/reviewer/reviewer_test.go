package reviewer_test

import (
	"strings"
	"testing"

	"adler-review-cli/internal/reviewer"
)

func TestSystemInstruction(t *testing.T) {
	instruction := reviewer.SystemInstruction
	if len(instruction) == 0 {
		t.Error("Expected SystemInstruction to be non-empty")
	}

	requiredSections := []string{
		"## Abertura",
		"## 1. Sobre o que é o livro?",
		"## 2. O que está sendo dito?",
		"## 3. O livro está certo?",
		"## 4. Qual a importância?",
	}

	for _, section := range requiredSections {
		if !strings.Contains(instruction, section) {
			t.Errorf("Expected SystemInstruction to contain section %q", section)
		}
	}
}

func TestConfirmBookInstruction(t *testing.T) {
	instruction := reviewer.ConfirmBookInstruction
	if len(instruction) == 0 {
		t.Error("Expected ConfirmBookInstruction to be non-empty")
	}

	requiredFields := []string{
		"**Título:**",
		"**Autor:**",
		"**Ano de Publicação:**",
		"**Gênero:**",
		"**Breve Descrição:**",
	}

	for _, field := range requiredFields {
		if !strings.Contains(instruction, field) {
			t.Errorf("Expected ConfirmBookInstruction to contain field %q", field)
		}
	}
}

func TestSlugify(t *testing.T) {
	// Let's write the exact test cases for our Slugify:
	t1 := reviewer.Slugify("Dom Casmurro")
	if t1 != "dom-casmurro" {
		t.Errorf("Expected 'dom-casmurro', got %q", t1)
	}

	t2 := reviewer.Slugify("O Pequeno Principe")
	if t2 != "o-pequeno-principe" {
		t.Errorf("Expected 'o-pequeno-principe', got %q", t2)
	}

	t3 := reviewer.Slugify("A Morte de Ivan Ilitch!")
	if t3 != "a-morte-de-ivan-ilitch" {
		t.Errorf("Expected 'a-morte-de-ivan-ilitch', got %q", t3)
	}
}
