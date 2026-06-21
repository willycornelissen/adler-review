## 1. Project Scaffolding & Dependencies

- [x] 1.1 Initialize the Go module with `go mod init adler-review-cli`.
- [x] 1.2 Scaffold the directory layout: create `cmd/adler-review/`, `internal/config/`, `internal/googleai/`, and `internal/reviewer/`.
- [x] 1.3 Add dependencies to `go.mod`: fetch `github.com/google/generative-ai-go/genai`, `google.golang.org/api/option`, and `github.com/joho/godotenv`.

## 2. CLI Interface & Configuration

- [x] 2.1 Implement `internal/config/config.go` to handle reading the API Key (with `--key` and `GEMINI_API_KEY` fallback) and loading `.env` files via `godotenv`.
- [x] 2.2 Implement argument parsing in `cmd/adler-review/main.go` using the standard `flag` library (or `cobra`) supporting `-o/--output`, `-k/--key`, `-m/--model`, and `-h/--help`.
- [x] 2.3 Implement file-existence and format validation inside `internal/reviewer/reviewer.go` using `os.Stat` and returning Go errors if unreadable.
- [x] 2.4 Add dynamic output name mapping inside `cmd/adler-review/main.go` so that if no output path is provided, it defaults to `<input-base>-resenha.md`.

## 3. Google Gen AI SDK Integration

- [x] 3.1 Implement client initialization in `internal/googleai/client.go` using `genai.NewClient` and `option.WithAPIKey`.
- [x] 3.2 Implement the `ListAndSelectBestProModel` function to iterate through the results of `client.ListModels(ctx)` in Go, filter for names containing "pro", and select the highest ranking model (Gemini 2.5 Pro > 1.5 Pro) with fallback.
- [x] 3.3 Add automatic exponential backoff retry logic inside `client.go` using standard Go sleep loops when encountering a rate limit (HTTP 429).

## 4. Mortimer Adler Review Generation

- [x] 4.1 Define the system instruction prompt in Portuguese inside `internal/reviewer/reviewer.go` incorporating Adler's 5-part analytical reading framework from `openspec/FORMATO.md`.
- [x] 4.2 Implement file reading logic in `internal/reviewer/reviewer.go` to load the full text of the input file into memory.
- [x] 4.3 Add a token verification check utilizing `client.CountTokens` (or a fallback word-count logic) to check if the text fits within the model's safe context limit before sending.
- [x] 4.4 Implement the generation call linking `googleai` and `reviewer` with a 5-minute context timeout (`context.WithTimeout`).

## 5. End-to-End Integration & Verification

- [x] 5.1 Implement output writing in `cmd/adler-review/main.go` using `os.WriteFile` to write the review in Markdown.
- [x] 5.2 Validate successful compilation with `go build -o adler-review cmd/adler-review/main.go`.
- [x] 5.3 Test the application end-to-end with a sample text file and verify that the generated review structure strictly conforms to `openspec/FORMATO.md`.
