# upftp 一行安装脚本 —— Windows (PowerShell)
#
# 用法(在 PowerShell 中执行):
#   irm https://github.com/zy84338719/upftp/raw/main/scripts/install.ps1 | iex
#
# 或(指定版本):
#   & ([scriptblock]::Create((irm https://github.com/zy84338719/upftp/raw/main/scripts/install.ps1))) 'v1.0.0'
#
# 行为:自动检测架构,从 GitHub Releases 下载 zip,解压到 %USERPROFILE%\.upftp\bin,
#       并把该目录加入用户 PATH(永久)。无需管理员权限。

$ErrorActionPreference = "Stop"

# --- 配置 ---
$Owner = "zy84338719"
$Repo = "upftp"
$InstallDir = if ($env:UPFTP_INSTALL_DIR) { $env:UPFTP_INSTALL_DIR } else { Join-Path $env:USERPROFILE ".upftp\bin" }

function Write-Info($m) { Write-Host "▸ $m" -ForegroundColor Cyan }
function Write-Ok($m)   { Write-Host "✓ $m" -ForegroundColor Green }
function Write-Warn($m) { Write-Host "⚠ $m" -ForegroundColor Yellow }
function Write-Err($m)  { Write-Host "✗ $m" -ForegroundColor Red }

# --- 检测架构 ---
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    "386"
}
$OS = "windows"

# --- 确定版本(支持通过参数指定)---
$Tag = $args[0]
if (-not $Tag) {
    Write-Info "查询最新版本..."
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Owner/$Repo/releases/latest" -Headers @{ "User-Agent" = "upftp-installer" }
        $Tag = $release.tag_name
    } catch {
        Write-Err "无法查询最新版本: $_"
        Write-Err "可手动指定: irm .../install.ps1 | iex 'v1.0.0'"
        exit 1
    }
}
$Version = $Tag.TrimStart("v")
Write-Ok "目标版本: $Tag ($OS/$Arch)"

# --- 拼接下载地址 ---
$Asset = "upftp_${Version}_${OS}_${Arch}.zip"
$DownloadUrl = "https://github.com/$Owner/$Repo/releases/download/$Tag/$Asset"

# --- 临时目录 ---
$TmpDir = Join-Path $env:TEMP "upftp-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null
try {
    # --- 下载 ---
    Write-Info "下载 $Asset ..."
    $ZipPath = Join-Path $TmpDir $Asset
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath -UseBasicParsing
    } catch {
        Write-Err "下载失败: $DownloadUrl"
        Write-Err "可能该架构未发布二进制。检查: https://github.com/$Owner/$Repo/releases/tag/$Tag"
        exit 1
    }

    # --- 校验 SHA256 ---
    Write-Info "校验完整性..."
    $ChecksumUrl = "https://github.com/$Owner/$Repo/releases/download/$Tag/checksums.txt"
    try {
        $Checksums = Invoke-RestMethod -Uri $ChecksumUrl -UseBasicParsing
        $ExpectedLine = ($Checksums -split "`n" | Where-Object { $_ -match $Asset })
        if ($ExpectedLine) {
            $Expected = ($ExpectedLine -split "\s+")[0]
            $Actual = (Get-FileHash $ZipPath -Algorithm SHA256).Hash.ToLower()
            if ($Expected.ToLower() -ne $Actual) {
                Write-Err "SHA256 校验失败! 期望 $Expected, 实际 $Actual"
                exit 1
            }
            Write-Ok "SHA256 校验通过"
        } else {
            Write-Warn "checksums.txt 中未找到 $Asset,跳过校验"
        }
    } catch {
        Write-Warn "无法下载 checksums.txt,跳过校验"
    }

    # --- 解压 ---
    Write-Info "解压..."
    Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force
    $Binary = Get-ChildItem -Path $TmpDir -Filter "upftp.exe" -Recurse | Select-Object -First 1
    if (-not $Binary) {
        Write-Err "解压后未找到 upftp.exe"
        exit 1
    }

    # --- 安装 ---
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    $Dest = Join-Path $InstallDir "upftp.exe"
    Move-Item -Path $Binary.FullName -Destination $Dest -Force
    Write-Ok "已安装到 $Dest"

    # --- 加入用户 PATH ---
    $UserPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [System.Environment]::SetEnvironmentVariable("Path", "$InstallDir;$UserPath", "User")
        $env:Path = "$InstallDir;$env:Path"
        Write-Ok "已将 $InstallDir 加入用户 PATH"
        Write-Warn "请重新打开终端使 PATH 生效"
    }

    # --- 验证 ---
    $VersionInstalled = & $Dest -version
    Write-Host ""
    Write-Host "  ✓ upftp 安装成功! " -ForegroundColor Green -NoNewline
    Write-Host $VersionInstalled
    Write-Host ""
    Write-Host "  快速开始:"
    Write-Host "    upftp                  # 分享当前目录"
    Write-Host "    upftp -d D:\share      # 分享指定目录"
    Write-Host ""
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}
