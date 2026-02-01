Create a file named `README.md` in your root folder and paste this in.

---

# 🤖 Staff-Bot: AI-Powered Engineering News Radar

> **A personal automated research assistant that monitors engineering blogs, summarizes them using GenAI, and delivers high-signal insights directly to your phone.**

## 🏗 Architecture

This is not a simple script; it is a **stateful, resilient daemon** designed for low-resource environments (e.g., local machines, Raspberry Pi, M1 Air).

* **Discovery Engine:** Polling RSS feeds (Netflix, Uber, AWS) using `gofeed`.
* **State Management:** JSON-based persistence (`history.json`) prevents re-processing articles across restarts.
* **The Brain (Strategy Pattern):** modular AI interface supporting **Google Gemini** (Cloud) or **Ollama** (Local).
* **Resilience:** Configurable rate limiting (RPM) and exponential backoff to handle strict API quotas.
* **Notification:** Push delivery via Telegram (with automatic message truncation/safe-guards).

## 🚀 Features

* **Smart Deduplication:** Never reads the same article twice.
* **Hardware Optimized:** Single-threaded design with `time.Sleep` pacing to respect API limits (e.g., 4 RPM for Gemini Free Tier).
* **Pluggable AI:** Easily swap between `gemini-3-flash`, `gemini-2.5-flash`, or local `llama3`.
* **Fault Tolerant:** Continues running even if one feed fails or the API rate limits.

## 🛠 Prerequisites

* **Go 1.21+** installed.
* **Google AI Studio Key** (Free tier is sufficient).
* **Telegram Bot Token** & Chat ID.

## ⚙️ Configuration & Setup

### 1. Clone & Install Dependencies

```bash
git clone https://github.com/yourname/crawler-agent.git
cd crawler-agent
go mod tidy

```

### 2. Environment Variables

Create a `.env` file or export these variables in your terminal:

```bash
# Required: Your Google AI Studio Key
export GEMINI_API_KEY="AIzaSy..."

# Optional: Select your model (defaults to gemini-2.5-flash)
# Use "gemini-3-flash" for speed, "gemini-1.5-pro" for depth.
export GEMINI_MODEL="gemini-3-flash"

# Required: Telegram Credentials
export TELEGRAM_TOKEN="123456:ABC-..."
export TELEGRAM_CHAT_ID="987654321"

```

### 3. Tuning Rate Limits

Modify `internal/config/config.go` to match your API tier:

```go
type AgentConfig struct {
    // ...
    RequestsPerMinute: 4, // 4 RPM = 1 request every 15s (Safe for Free Tier)
}

```

## 🏃‍♂️ Usage

Run the agent as a foreground process:

```bash
go run cmd/agent/main.go

```

**What happens next:**

1. The agent loads `history.json` (or creates it).
2. It scans the configured RSS feeds for new links.
3. It filters out seen URLs.
4. It processes new articles **one by one** (respecting the RPM limit).
5. It sends a summarized digest to your Telegram.
6. It sleeps for the configured interval (e.g., 6 hours) before the next scan.

## 📂 Project Structure

```text
.
├── cmd/
│   └── agent/         # Main entry point (The Orchestrator)
├── internal/
│   ├── brain/         # AI Logic (Gemini/Ollama Strategy Pattern)
│   ├── config/        # Centralized configuration & Rate limits
│   ├── crawler/       # RSS fetching & HTML parsing
│   ├── notifier/      # Telegram/Discord clients
│   └── state/         # JSON file persistence logic
├── history.json       # (Ignored) Local database of seen URLs
├── go.mod             # Dependency definitions
└── README.md          # This file

```

## 🔮 Future Roadmap

* [ ] **Discord Webhooks:** Rich embed support for better readability.
* [ ] **Static Site Gen:** Auto-publish daily digests to a Hugo blog.
* [ ] **Topic Filtering:** Only notify about "Database" or "System Design" articles.

## 🤝 Contributing

1. Fork the repo.
2. Create a feature branch (`git checkout -b feature/discord-integration`).
3. Commit your changes.
4. Push to the branch.
5. Open a Pull Request.

---

*Built with ❤️ by [Your Name]*
