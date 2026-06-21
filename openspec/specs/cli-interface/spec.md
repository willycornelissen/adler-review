## ADDED Requirements

### Requirement: Command execution and argument parsing
The CLI tool SHALL accept a positional argument representing the input file path, and optional arguments for the output file path (`--output` or `-o`), Google API Key (`--key` or `-k`), and Model name (`--model` or `-m`).

#### Scenario: Valid execution with only required argument
- **WHEN** the CLI is executed with `adler-review-cli path/to/book.txt`
- **THEN** the system parses the input path as `path/to/book.txt` and uses default values for output path (`path/to/book-resenha.md`), checks environment variable for API key, and automatically selects the best available model.

#### Scenario: Valid execution with all arguments provided
- **WHEN** the CLI is executed with `adler-review-cli path/to/book.txt -o reviews/my-review.md -k MY_SECRET_KEY -m gemini-2.5-pro`
- **THEN** the system uses the specified output path `reviews/my-review.md`, the API key `MY_SECRET_KEY`, and overrides model selection to `gemini-2.5-pro`.

### Requirement: Input file validation
The CLI tool SHALL verify that the input file exists, is readable, and is a valid text file (UTF-8 format). If not, it SHALL print a clear error message and terminate with exit code 1.

#### Scenario: Input file does not exist
- **WHEN** the CLI is executed with a path to a non-existent file
- **THEN** the system prints a message like "Error: Input file does not exist" to stderr and exits with code 1.

### Requirement: Help command
The CLI tool SHALL provide a clear help message detailing its usage, options, and environment variables when `--help` or `-h` is requested.

#### Scenario: Help requested
- **WHEN** the CLI is executed with `-h` or `--help`
- **THEN** the system outputs a usage description detailing all commands, options, and exit codes to stdout and exits with code 0.
