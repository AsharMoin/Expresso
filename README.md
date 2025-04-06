# 🚀 Expresso

<p align="center">
  <span style="font-size: 1.5em; font-weight: bold; color: #0066cc;">Express-o your commands effortlessly</span>
</p>
<p align="center">
  <em>Transform natural language into powerful shell commands with AI</em>
</p>
<p align="center">

<p align="center">
  <a href="https://raw.githubusercontent.com/AsharMoin/Expresso/main/assets/demo.gif">
    <img src="https://raw.githubusercontent.com/AsharMoin/Expresso/main/assets/demo.gif" alt="Expresso Demo" width="800" />
  </a>
</p>

## 🌟 Overview

Expresso is a command-line utility that transforms your natural language requests into precise, executable shell commands. No more searching for the exact syntax or command flags - just tell Expresso what you want to accomplish in plain English, and it will generate the appropriate command for your shell.

## ✨ Features

- Convert natural language to shell commands
- Works across multiple shells (Bash, Zsh, PowerShell, cmd)
- Powered by OpenAI's GPT models
- Provides command descriptions
- Command confirmation before execution
- Secure API key management

## 📋 Requirements

- An OpenAI API key
- Go 1.18+ (for building from source)

## ⚙️ Installation

### Option 1: Quick Install (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/AsharMoin/Expresso/main/install.sh | bash
```

### Option 2: Download Binary

Visit the [releases page](https://github.com/AsharMoin/Expresso/releases) to download the latest binary for your platform.

## 🚀 Getting Started

1. Run Expresso for the first time:
   ```bash
   expresso
   ```

2. When prompted, enter your OpenAI API key (required only on first run)

## 💡 Usage Examples

| Natural Language Request | Generated Command |
|--------------------------|-------------------|
| `expresso list all nodejs processes` | `ps aux \| grep node` |
| `expresso create a backup of my project folder` | `tar -czvf project_backup.tar.gz ./project` |
| `expresso find large files over 100MB` | `find . -type f -size +100M -exec ls -lh {} \;` |
| `expresso show me memory usage` | `free -h` |
| `expresso download video from youtube` | `youtube-dl "https://www.youtube.com/watch?v=VIDEO_ID"` |

## 🛠️ Configuration

Expresso stores your configuration in `~/.config/expresso/config.yaml`. The configuration includes your OpenAI API key.

To change your API key, you can edit this file directly or delete it to be prompted again on next run.

## 🧹 Uninstallation

To remove Expresso from your system:

```bash
curl -fsSL https://raw.githubusercontent.com/AsharMoin/Expresso/main/uninstall.sh | bash
```

---

<p align="center">Made with ☕ and 💻 by <a href="https://github.com/AsharMoin">Ashar Moin</a></p>
