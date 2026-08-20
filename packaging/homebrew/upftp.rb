# Homebrew formula 参考模板。
#
# 实际使用时,goreleaser 会根据 .goreleaser.yml 的 brews 配置,
# 自动生成此文件并推送到 github.com/zy84338719/homebrew-tap 仓库。
# 用户无需手动维护,这里仅作参考与文档。
#
# 用户安装:
#   brew tap zy84338719/tap
#   brew install upftp

class Upftp < Formula
  desc "快速、即开即用的 FTP 文件分享工具"
  homepage "https://github.com/zy84338719/upftp"
  url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_1.0.0_darwin_arm64.tar.gz"
  sha256 "REPLACED_BY_GORELEASER"
  version "1.0.0"
  license "MIT"

  # goreleaser 用 on_system 按架构选择不同下载包,
  # 这里展示 arm64 分支;amd64 与 linux 分支由 goreleaser 自动补全。
  on_macos do
    on_arm do
      url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_1.0.0_darwin_arm64.tar.gz"
      sha256 "REPLACED_BY_GORELEASER"
    end
    on_intel do
      url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_1.0.0_darwin_amd64.tar.gz"
      sha256 "REPLACED_BY_GORELEASER"
    end
  end

  def install
    bin.install "upftp"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/upftp -version")
  end
end
