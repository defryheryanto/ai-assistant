package config

var (
	GoogleCredentialsFilePath string
	GoogleTokenFilePath       string

	OpenAIToken string
	OpenAIModel string

	WhatsmeowSQLPath string
)

func Init() {
	GoogleCredentialsFilePath = getString("GOOGLE_CREDENTIALS_FILEPATH", "")
	GoogleTokenFilePath = getString("GOOGLE_TOKEN_FILEPATH", "")

	OpenAIToken = getString("OPEN_AI_TOKEN", "")
	OpenAIModel = getString("OPEN_AI_MODEL", "")

	WhatsmeowSQLPath = getString("WHATSMEOW_SQL_PATH", "")
}
