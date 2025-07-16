#!/bin/bash

# Script to publish packages to private APT repository
# This script should be run after GoReleaser creates the packages

set -e

REPO_PATH="/var/www/apt"
DIST="stable"
COMPONENT="main"
ARCHITECTURES="amd64 arm64"

echo "Publishing packages to APT repository..."

# Create repository structure
for arch in $ARCHITECTURES; do
    mkdir -p "${REPO_PATH}/dists/${DIST}/${COMPONENT}/binary-${arch}"
done

# Copy packages
find dist/ -name "*.deb" -exec cp {} "${REPO_PATH}/pool/${COMPONENT}/" \;

# Generate Packages files
for arch in $ARCHITECTURES; do
    echo "Generating Packages file for ${arch}..."
    cd "${REPO_PATH}"
    dpkg-scanpackages "pool/${COMPONENT}" /dev/null | gzip -9c > "dists/${DIST}/${COMPONENT}/binary-${arch}/Packages.gz"
    dpkg-scanpackages "pool/${COMPONENT}" /dev/null > "dists/${DIST}/${COMPONENT}/binary-${arch}/Packages"
done

# Generate Release file
cd "${REPO_PATH}/dists/${DIST}"
cat > Release << EOF
Origin: UPFTP Repository
Label: UPFTP
Suite: ${DIST}
Codename: ${DIST}
Date: $(date -Ru)
Architectures: ${ARCHITECTURES}
Components: ${COMPONENT}
Description: UPFTP Private Repository
EOF

# Generate checksums
echo "MD5Sum:" >> Release
find . -name "Packages*" -exec md5sum {} \; | sed 's/\.\///g' | awk '{print " " $1 " " $2 " " $3}' >> Release

echo "SHA1:" >> Release  
find . -name "Packages*" -exec sha1sum {} \; | sed 's/\.\///g' | awk '{print " " $1 " " $2 " " $3}' >> Release

echo "SHA256:" >> Release
find . -name "Packages*" -exec sha256sum {} \; | sed 's/\.\///g' | awk '{print " " $1 " " $2 " " $3}' >> Release

# Sign the release (optional, requires GPG key)
if command -v gpg &> /dev/null && [ -n "$GPG_KEY_ID" ]; then
    echo "Signing release with GPG key: $GPG_KEY_ID"
    gpg --default-key "$GPG_KEY_ID" --armor --detach-sign --sign --output Release.gpg Release
    gpg --default-key "$GPG_KEY_ID" --clear-sign --output InRelease Release
fi

echo "APT repository updated successfully!"
echo "Repository URL: https://your-domain.com/apt"
echo ""
echo "To add this repository:"
echo "echo 'deb https://your-domain.com/apt ${DIST} ${COMPONENT}' | sudo tee /etc/apt/sources.list.d/upftp.list"
echo "sudo apt update"
echo "sudo apt install upftp"
