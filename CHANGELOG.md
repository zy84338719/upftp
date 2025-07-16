# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 🌍 **Multi-language Support**: Complete Chinese/English interface
  - Auto-detect browser language preference
  - Manual language switching with real-time effect
  - Language preference persistence using localStorage
  - Complete translation of all UI elements including:
    - Header titles and descriptions
    - Server information cards
    - Navigation breadcrumbs
    - Search placeholders
    - File type labels
    - Action buttons (Preview, Download)
    - Modal dialogs and error messages
    - Empty states and loading messages

### Enhanced
- 🎨 **Improved User Experience**
  - Added language selector buttons in header (🇺🇸 English / 🇨🇳 中文)
  - Language selection buttons with active state indicators
  - Seamless language switching without page refresh
  - Better accessibility with proper language attributes

### Technical
- JavaScript-based dynamic language switching
- Data attributes for storing translations (`data-en`, `data-zh`)
- Browser language detection using `navigator.language`
- Local storage integration for preference persistence

## [1.0.0] - 2025-07-16

### Added
- 🚀 **Cross-platform File Sharing Server**
  - Modern HTTP web interface with responsive design
  - FTP server support with authentication
  - File preview capabilities for images, videos, audio, text, and code
  - Real-time file search functionality
  - Interactive command-line interface

### Features
- **File Preview Support**:
  - Images: JPG, PNG, GIF, SVG, WebP, AVIF, BMP, TIFF
  - Videos: MP4, AVI, MOV, WebM, MKV, FLV, 3GP
  - Audio: MP3, WAV, FLAC, AAC, OGG, M4A, WMA
  - Text/Code: TXT, MD, JSON, XML, YAML, and various programming languages
  - Archives: ZIP, RAR, 7Z, TAR, GZ
  - Documents: PDF, DOC, DOCX, PPT, PPTX, XLS, XLSX

- **Web Interface**:
  - Beautiful modern UI with gradient backgrounds
  - File type icons and metadata display
  - Download buttons and preview modals
  - Mobile-responsive design
  - Search and filter capabilities

- **Server Features**:
  - HTTP server with embed file system
  - Optional FTP server with user authentication
  - Cross-platform support (Linux, macOS, Windows)
  - Multiple architecture support (amd64, arm64, 386)
  - Auto network interface selection
  - Configurable ports and directories

### Technical
- Built with Go 1.21+
- Embed file system for templates
- Modular architecture with separate packages:
  - `config/`: Configuration management
  - `network/`: Network interface detection
  - `ftp/`: FTP server implementation
  - `filehandler/`: File type detection and handling
  - `logic/`: HTTP server and UI logic

### Build System
- Comprehensive Makefile with multiple targets
- Cross-platform build support
- Release packaging automation
- Go module with China proxy support
- GitHub Actions integration ready
