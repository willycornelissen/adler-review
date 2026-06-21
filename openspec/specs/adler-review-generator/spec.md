## ADDED Requirements

### Requirement: Structured Review Format Adherence
The system SHALL instruct the Gemini Pro model to output the review strictly in Portuguese and in the exact format defined in `openspec/FORMATO.md`. This format consists of:
1. `## Abertura` (Context)
2. `## 1. Sobre o que é o livro?` (Classification, central thesis, and unity)
3. `## 2. O que está sendo dito?` (Structure, propositions, and argument progression)
4. `## 3. O livro está certo?` (Logical coherence, completeness, correctness, and verisimilitude critique)
5. `## 4. Qual a importância?` (Significance, connections, and utility)

#### Scenario: Successful formatting and generation
- **WHEN** the system processes a valid input file and generates the review
- **THEN** the output file is saved as Markdown, written in Portuguese, containing exactly these five sections matching the headings, tone, and depth shown in the `openspec/FORMATO.md` example.

### Requirement: Structured Prompting with Context and Metamethodology
The system SHALL construct a specialized, high-fidelity system instruction or prompt that educates the Gemini Pro model about Mortimer Adler's analytical reading framework (from *How to Read a Book*) and instructs it to analyze the input text according to that precise framework.

#### Scenario: Crafting the system prompt
- **WHEN** the generation request is prepared
- **THEN** the system generates a prompt containing Adler's reading rules, details from `openspec/FORMATO.md`, and instructions to avoid pre-judging in the first two questions.

### Requirement: Full-Context Ingestion
The system SHALL load the entire content of the input text file and transmit it within the prompt to utilize Gemini Pro's large context window (1M+ tokens), enabling a comprehensive non-chunked reading of the entire text.

#### Scenario: Processing a long book file
- **WHEN** the input file contains a full book text of 500,000 words (approximately 650,000 tokens)
- **THEN** the system reads the full file, verifies that it does not exceed the maximum token limit of the selected model (typically 1,000,000 to 2,000,000 tokens), and passes the entire text directly in a single API call.

### Requirement: Content limit checks
The system SHALL check the approximate token count of the input text (using 1 word ≈ 1.3 tokens as a conservative estimate, or the official token counting API) and print a warning or error if the input exceeds 90% of the selected model's context window.

#### Scenario: Input exceeds safe context window
- **WHEN** the input text is calculated to exceed 1.8M tokens for a model with a 2M token limit (or 900k for 1M limit)
- **THEN** the system prints a warning to stderr recommending truncating or chunking, or exits with code 1 if it exceeds the hard limit.
