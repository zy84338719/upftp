class Upftp < Formula
  desc "A lightweight file sharing server"
  homepage "https://github.com/zy84338719/upftp"
  version "v1.0.0" # 替换为您的版本

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_darwin_amd64.tar.gz"
    sha256 "xxx" # 替换为实际的 SHA256
  elsif OS.mac? && Hardware::CPU.arm?
    url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_darwin_arm64.tar.gz"
    sha256 "xxx" # 替换为实际的 SHA256
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_linux_amd64.tar.gz"
    sha256 "xxx" # 替换为实际的 SHA256
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/zy84338719/upftp/releases/download/v1.0.0/upftp_linux_arm64.tar.gz"
    sha256 "xxx" # 替换为实际的 SHA256
  end

  def install
    bin.install "upftp"
  end

  test do
    system "#{bin}/upftp", "-h"
  end
end 