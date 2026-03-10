# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 🎨 **Enhanced UI/UX Improvements**
  - Beautiful ASCII art logo banner in CLI
  - Improved command menu with emoji icons
  - New server status and about information screens
  - Better formatted output with box drawings
  - Version info command (`v` or `version`)
  
- 📝 **Configuration Enhancements**
  - Improved configuration file with detailed comments
  - Added usage examples in config file
  - Better section organization with visual separators
  
- 🔧 **Extended File Type Support**
  - Added support for 30+ new file extensions:
    - Modern image formats (HEIC, HEIF, AVIF, TIFF)
    - Additional video formats (MPEG, 3GP, OGV)
    - More audio formats (OPUS, AIFF, APE)
    - Extended programming languages (Vue, Svelte, Dart, Kotlin, etc.)
    - Container formats (Dockerfile, Makefile)
    
- 📊 **Better Logging**
  - Improved log format with visual indicators
  - Color-coded log levels with bullet points
  - More informative startup messages

### Enhanced
- 🚀 **CLI Experience**
  - Removed empty directory check at startup
  - Better error messages with emoji indicators
  - Improved file listing with pagination (max 20 items)
  - Enhanced download examples with MCP integration info
  
- 🌐 **Web Interface**
  - Updated page title to "AI-First File Server"
  - Added version badge in header
  - New feature showcase bar
  - Improved upload section UI
  - Better responsive design for mobile
  - Added QR code button for mobile access
  
- 🔒 **Security & Configuration**
  - Configuration path tracking and display
  - Better validation for upload paths
  - Enhanced error handling

### Technical
- Added `GetConfigPath()` function to config package
- Improved logger output formatting
- Extended file type detection system
- Better startup logging with detailed server info

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
