# ai-assistant

A modular AI agent framework built in Go, designed for easy integration with various clients (such as WhatsApp) and providing a flexible, importable tooling system for function-calling and tool execution.

---

## Features

- **AI Agent Core**  
  Easily integrate with messaging clients (currently supports WhatsApp) to provide conversational AI features.

- **Importable Tooling Registry**  
  Provides a registry for function-calling tools, making it simple to register, manage, and execute new tools in any Go project.

- **Reusable Tool Implementations**  
  Includes a growing set of tools that can be used directly or as references for building your own.

---

## Setting Up
### WhatsApp
1. Clone the repository<br> `git clone git@github.com:defryheryanto/ai-assistant.git`
2. Set Environment Variables<br>
See [.env.example](https://github.com/defryheryanto/ai-assistant/blob/main/.env.example) file to see the environment variables required to run the project
3. Running the application<br>
Run the application using this command<br>
`go run ./cmd/whatsapp/...`<br>
or build the application binary<br>
`go build ./cmd/whatsapp -o {binary_filename}`
4. Copy the QRCode text from the terminal `(Only need to setup once)`
5. Open https://www.the-qrcode-generator.com/ and navigate to 'Free Text' tab `(Only need to setup once)`
6. Paste the QR Code text `(Only need to setup once)`
7. Scan the generated QR from your WhatsApp `(Only need to setup once)`
8. Done! Now chat with your WhatsApp account via Personal Chat to interact

---

## Usage as a Library
### Tooling System
- Registry<br>
Centralized struct to register and execute tools, supporting integration with any LLM that implements `llms.Model`.
- Tool Interface<br>
Implement `Tool` [interface](https://github.com/defryheryanto/ai-assistant/blob/cffc53f22279208a31233f1bf896621cb018960c/pkg/tools/registry.go#L10) in `pkg/tools/registry.go` to build your own tools.

### Custom Tool Integration
Visit [weather-forecast](https://github.com/defryheryanto/ai-assistant/tree/main/example/weather-forecast) example to see the implementation for custom tool

---
