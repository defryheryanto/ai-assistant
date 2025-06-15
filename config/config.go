package config

var (
	AppName                         string
	DatabaseConnectionString        string
	IsUserWhitelistEnabled          bool
	IsWhatsAppGroupWhitelistEnabled bool

	GoogleCredentialsFilePath string
	GoogleTokenFilePath       string

	OpenAIToken string
	OpenAIModel string

	WhatsmeowSQLPath string
)

func Init() {
	AppName = getString("APP_NAME", "ai-assistant")
	DatabaseConnectionString = getString("DATABASE_CONNECTION_STRING", "")
	IsUserWhitelistEnabled = getBool("IS_USER_WHITELIST_ENABLED", false)
	IsWhatsAppGroupWhitelistEnabled = getBool("IS_WHATSAPP_GROUP_WHITELIST_ENABLED", false)

	GoogleCredentialsFilePath = getString("GOOGLE_CREDENTIALS_FILEPATH", "")
	GoogleTokenFilePath = getString("GOOGLE_TOKEN_FILEPATH", "")

	OpenAIToken = getString("OPEN_AI_TOKEN", "")
	OpenAIModel = getString("OPEN_AI_MODEL", "")

	WhatsmeowSQLPath = getString("WHATSMEOW_SQL_PATH", "")
}
