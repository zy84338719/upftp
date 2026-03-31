# 📸 UPFTP Showcase

## 🖥️ Web Interface

```
┌──────────────────────────────────────────────────────────────────┐
│ UPFTP - AI-First File Server                    v0.1.3  [🌐 EN]  │
├──────────────────────────────────────────────────────────────────┤
│ [📁 Files] [📤 Upload] [⚙️ Settings] [ℹ️ About]                  │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  📁 documents/                                📅 2026-03-31      │
│     📄 sample.txt                             337 bytes          │
│     📄 data.json                              215 bytes          │
│                                                                   │
│  📁 images/                                   📅 2026-03-31      │
│     🖼️ logo.png                              2.3 MB             │
│                                                                   │
│  📁 projects/                                 📅 2026-03-31      │
│     📦 source.tar.gz                          1.5 MB             │
│                                                                   │
│  📄 README.md                                 337 bytes          │
│                                                                   │
├──────────────────────────────────────────────────────────────────┤
│ 📊 Total: 5 items (4.2 MB)    [📤 Upload Files] [📁 New Folder] │
└──────────────────────────────────────────────────────────────────┘
```

## 💻 Command Line Interface

```bash
$ upftp -auto

Available network interfaces:
[0] 10.10.30.99
[1] 192.168.31.185

Starting UPFTP v0.1.3...
✓ HTTP server: http://10.10.30.99:10002
✓ FTP server: ftp://10.10.30.99:2123
✓ WebDAV server: http://10.10.30.99:8082
✓ NFS server: nfs://10.10.30.99:2051
✓ File upload: enabled (max size: 100 MB)

Server is running. Press Ctrl+C to stop.
```

## 📊 Server Status Dashboard

```
╔══════════════════════════════════════════════════════════════╗
║                    Server Status                              ║
╠══════════════════════════════════════════════════════════════╣
║ ✅ HTTP Server      http://10.10.30.99:10002                  ║
║ ✅ FTP Server       ftp://10.10.30.99:2123                    ║
║ ✅ WebDAV Server    http://10.10.30.99:8082                   ║
║ ✅ NFS Server       nfs://10.10.30.99:2051                    ║
║                                                              ║
║ 📁 Shared Directory: /tmp/upftp-share                        ║
║ 📤 File Upload: enabled (max size: 100 MB)                   ║
║ 👤 FTP Credentials: zhangyi / ****                            ║
╚══════════════════════════════════════════════════════════════╝
```

---

## 📷 Real Screenshots

For real screenshots of the web interface, please visit our [GitHub repository](https://github.com/zy84338719/upftp).

To create your own screenshots:

1. **Start the server**:
   ```bash
   upftp -auto -d /path/to/share
   ```

2. **Open in browser**:
   ```bash
   open http://localhost:10002
   ```

3. **Take screenshot**:
   - **macOS**: `Cmd + Shift + 4`, then `Space`, click browser window
   - **Windows**: `Win + Shift + S`
   - **Linux**: Use `gnome-screenshot` or `flameshot`
