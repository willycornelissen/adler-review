## ADDED Requirements

### Requirement: API Authentication Validation
The system SHALL validate that a Google AI API key is available via either the CLI `--key` argument or the `GEMINI_API_KEY` environment variable. If neither is available, it SHALL print an error and exit with code 1.

#### Scenario: Missing API Key
- **WHEN** the CLI is executed with no `--key` option and `GEMINI_API_KEY` is not set in the environment
- **THEN** the system prints "Error: Google API key is missing. Please set the GEMINI_API_KEY environment variable or provide it via the --key parameter." to stderr and exits with code 1.

### Requirement: Dynamic Pro Model Selection
The system SHALL automatically discover and select the best available Pro-tier Gemini model from the user's account by querying the list of available models, prioritizing newer or higher-capability Pro models, and falling back to a hardcoded default Pro-tier model if model listing fails or returns no active models.

#### Scenario: Successful automatic model discovery
- **WHEN** the CLI is executed with automatic model selection, and the API returns `gemini-2.5-pro` and `gemini-1.5-pro` in the active models list
- **THEN** the system automatically selects `gemini-2.5-pro` as the highest capability Pro model for generation.

#### Scenario: Fallback when listing fails or returns empty
- **WHEN** the model listing API call fails or returns empty, and no model override is provided
- **THEN** the system logs a warning and falls back to using `gemini-1.5-pro` as the default model.

### Requirement: Rate Limit Handling
The system SHALL handle API rate limits (HTTP 429 / RESOURCE_EXHAUSTED) gracefully by employing an exponential backoff retry mechanism before giving up and throwing an error.

#### Scenario: Rate limit encountered with eventual success
- **WHEN** a rate limit error is received from the Google Gen AI API
- **THEN** the system waits for an increasing backoff delay, retries the request up to 3 times, and successfully completes when a retry succeeds.
