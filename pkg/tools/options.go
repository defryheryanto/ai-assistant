package tools

type Option func(*registry)

func WithLoggerOption() Option {
	return func(r *registry) {
		r.enableLog = true
	}
}

func WithSystemPromptOption(systemPrompt string) Option {
	return func(r *registry) {
		r.systemPrompt = systemPrompt
	}
}

func WithContextWindowManager(cw ContextWindowManager) Option {
	return func(r *registry) {
		r.contextWindowManager = cw
	}
}
