#!/bin/bash
#
# Universal RPM/DEB packaging script for Go projects
#
# Usage:
#   ./scripts/package.sh --name myapp --version 1.0.0
#   ./scripts/package.sh --format deb
#   ./scripts/package.sh --format rpm
#   ./scripts/package.sh --format all
#   ./scripts/package.sh --arch amd64,arm64
#
# No external dependencies required (no fpm, no goreleaser).
# Builds .deb and .rpm from pure shell.
#

set -euo pipefail

# ======================= Color Output =======================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()    { echo -e "${BLUE}[INFO]${NC} $*"; }
success() { echo -e "${GREEN}[OK]${NC} $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $*"; }
error()   { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

# ======================= Default Config =======================
PKG_NAME=""
PKG_VERSION=""
PKG_DESCRIPTION=""
PKG_MAINTAINER=""
PKG_VENDOR=""
PKG_HOMEPAGE=""
PKG_LICENSE="MIT"
PKG_SECTION="net"
PKG_PRIORITY="optional"
PKG_BINDIR="/usr/bin"
PKG_ETCDIR="/etc/${PKG_NAME}"
PKG_SYSCONFDIR="/etc/systemd/system"
PKG_DATADIR="/var/lib/${PKG_NAME}"
PKG_LOGDIR="/var/log/${PKG_NAME}"
PKG_USER="${PKG_NAME}"
PKG_GROUP="${PKG_NAME}"

FORMAT="all"
ARCH="amd64"
OUTPUT_DIR="dist/packages"
BUILD_DIR="dist/build"
GOOS="linux"
SYSTEMD_SERVICE=1
CONFIG_FILE=""
EXTRA_FILES=""
MAIN_PATH="."
PREINST_SCRIPT=""
POSTINST_SCRIPT=""
PRERM_SCRIPT=""
POSTRM_SCRIPT=""
STRIP=1
VERBOSE=0

# ======================= Parse Arguments =======================
usage() {
    cat <<EOF
Universal RPM/DEB Packaging Script for Go Projects

Usage: $0 [options]

Options:
  --name NAME           Package name (default: derived from go.mod or directory name)
  --version VERSION     Package version (default: git tag or "dev")
  --description DESC    Package description
  --maintainer EMAIL    Maintainer info (e.g. "Name <email>")
  --vendor VENDOR       Vendor name
  --homepage URL        Project homepage URL
  --license LICENSE     License name (default: MIT)
  --format FORMAT       Package format: deb, rpm, all (default: all)
  --arch ARCHS          Target architectures: amd64, arm64, amd64,arm64 (default: amd64)
  --output DIR          Output directory (default: dist/packages)
  --main PATH           Go main package path (default: .)
  --config FILE         Config file to include (e.g. upftp.example.yaml)
  --service FILE        Systemd service file (default: auto-generate)
  --no-service          Skip systemd service file
  --preinst SCRIPT      Pre-install script
  --postinst SCRIPT     Post-install script
  --prerm SCRIPT        Pre-remove script
  --postrm SCRIPT       Post-remove script
  --extra "SRC:DST"     Extra files to include (repeatable)
  --no-strip            Don't strip debug symbols
  --verbose             Verbose output
  --help                Show this help

Examples:
  # Build .deb and .rpm for amd64
  $0 --name upftp --version 1.0.0

  # Build only .deb for arm64
  $0 --format deb --arch arm64

  # Build for multiple architectures
  $0 --arch amd64,arm64

  # Custom config and service files
  $0 --config upftp.example.yaml --service packaging/systemd/upftp.service

EOF
    exit 0
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --name)         PKG_NAME="$2"; shift 2 ;;
        --version)      PKG_VERSION="$2"; shift 2 ;;
        --description)  PKG_DESCRIPTION="$2"; shift 2 ;;
        --maintainer)   PKG_MAINTAINER="$2"; shift 2 ;;
        --vendor)       PKG_VENDOR="$2"; shift 2 ;;
        --homepage)     PKG_HOMEPAGE="$2"; shift 2 ;;
        --license)      PKG_LICENSE="$2"; shift 2 ;;
        --format)       FORMAT="$2"; shift 2 ;;
        --arch)         ARCH="$2"; shift 2 ;;
        --output)       OUTPUT_DIR="$2"; shift 2 ;;
        --main)         MAIN_PATH="$2"; shift 2 ;;
        --config)       CONFIG_FILE="$2"; shift 2 ;;
        --service)      SYSTEMD_SERVICE_FILE="$2"; shift 2 ;;
        --no-service)   SYSTEMD_SERVICE=0; shift ;;
        --preinst)      PREINST_SCRIPT="$2"; shift 2 ;;
        --postinst)     POSTINST_SCRIPT="$2"; shift 2 ;;
        --prerm)        PRERM_SCRIPT="$2"; shift 2 ;;
        --postrm)       POSTRM_SCRIPT="$2"; shift 2 ;;
        --extra)        EXTRA_FILES="${EXTRA_FILES} $2"; shift 2 ;;
        --no-strip)     STRIP=0; shift ;;
        --verbose)      VERBOSE=1; shift ;;
        --help|-h)      usage ;;
        *) warn "Unknown option: $1"; shift ;;
    esac
done

# ======================= Auto-detect Config =======================
detect_name() {
    if [[ -n "${PKG_NAME}" ]]; then
        return
    fi
    if [[ -f "go.mod" ]]; then
        PKG_NAME=$(grep '^module' go.mod | head -1 | awk '{print $2}' | xargs basename 2>/dev/null || true)
    fi
    if [[ -z "${PKG_NAME}" ]]; then
        PKG_NAME=$(basename "$(pwd)")
    fi
    info "Auto-detected package name: ${PKG_NAME}"
}

detect_version() {
    if [[ -n "${PKG_VERSION}" ]]; then
        return
    fi
    PKG_VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    PKG_VERSION="${PKG_VERSION#v}"
    info "Auto-detected version: ${PKG_VERSION}"
}

detect_config() {
    if [[ -n "${CONFIG_FILE}" ]]; then
        return
    fi
    for f in "${PKG_NAME}.yaml" "${PKG_NAME}.yml" "${PKG_NAME}.conf" "${PKG_NAME}.example.yaml" "${PKG_NAME}.example.yml" "config.yaml" "config.yml"; do
        if [[ -f "$f" ]]; then
            CONFIG_FILE="$f"
            info "Auto-detected config file: ${CONFIG_FILE}"
            return
        fi
    done
}

detect_systemd_service() {
    if [[ "${SYSTEMD_SERVICE}" -eq 0 ]]; then
        return
    fi
    if [[ -n "${SYSTEMD_SERVICE_FILE:-}" ]]; then
        return
    fi
    for f in "packaging/systemd/${PKG_NAME}.service" "scripts/${PKG_NAME}.service" "${PKG_NAME}.service"; do
        if [[ -f "$f" ]]; then
            SYSTEMD_SERVICE_FILE="$f"
            info "Auto-detected systemd service: ${SYSTEMD_SERVICE_FILE}"
            return
        fi
    done
}

detect_maintainer() {
    if [[ -n "${PKG_MAINTAINER}" ]]; then
        return
    fi
    local git_user git_email
    git_user=$(git config user.name 2>/dev/null || true)
    git_email=$(git config user.email 2>/dev/null || true)
    if [[ -n "${git_user}" && -n "${git_email}" ]]; then
        PKG_MAINTAINER="${git_user} <${git_email}>"
        PKG_VENDOR="${git_user}"
        info "Auto-detected maintainer: ${PKG_MAINTAINER}"
    else
        PKG_MAINTAINER="Unknown <unknown@example.com>"
        PKG_VENDOR="Unknown"
    fi
}

detect_homepage() {
    if [[ -n "${PKG_HOMEPAGE}" ]]; then
        return
    fi
    local remote_url
    remote_url=$(git config --get remote.origin.url 2>/dev/null || true)
    if [[ -n "${remote_url}" ]]; then
        remote_url="${remote_url%.git}"
        if [[ "${remote_url}" == git@github.com:* ]]; then
            remote_url="https://github.com/${remote_url#git@github.com:}"
        fi
        PKG_HOMEPAGE="${remote_url}"
        info "Auto-detected homepage: ${PKG_HOMEPAGE}"
    else
        PKG_HOMEPAGE="https://example.com"
    fi
}

detect_scripts() {
    local dirs=("packaging/debian" "packaging/scripts" "scripts")
    for dir in "${dirs[@]}"; do
        if [[ -z "${PREINST_SCRIPT}" && -f "${dir}/preinst" ]]; then
            PREINST_SCRIPT="${dir}/preinst"
        fi
        if [[ -z "${POSTINST_SCRIPT}" && -f "${dir}/postinst" ]]; then
            POSTINST_SCRIPT="${dir}/postinst"
        fi
        if [[ -z "${PRERM_SCRIPT}" && -f "${dir}/prerm" ]]; then
            PRERM_SCRIPT="${dir}/prerm"
        fi
        if [[ -z "${POSTRM_SCRIPT}" && -f "${dir}/postrm" ]]; then
            POSTRM_SCRIPT="${dir}/postrm"
        fi
    done
}

detect_description() {
    if [[ -n "${PKG_DESCRIPTION}" ]]; then
        return
    fi
    PKG_DESCRIPTION="${PKG_NAME} - built from source"
    info "Auto-detected description: ${PKG_DESCRIPTION}"
}

# ======================= Go Build =======================
build_binary() {
    local goarch="$1"
    local output="$2"

    local commit build_date
    commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    build_date=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    local ldflags="-X main.Version=${PKG_VERSION} -X main.LastCommit=${commit} -X main.BuildDate=${build_date}"
    if [[ "${STRIP}" -eq 1 ]]; then
        ldflags="${ldflags} -s -w"
    fi

    info "Building ${PKG_NAME} for ${GOOS}/${goarch}..."
    CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${goarch}" go build \
        -trimpath \
        -ldflags "${ldflags}" \
        -o "${output}" \
        "${MAIN_PATH}"

    if [[ ! -f "${output}" ]]; then
        error "Build failed: ${output} not found"
    fi

    local size
    size=$(du -h "${output}" | awk '{print $1}')
    success "Built ${output} (${size})"
}

# ======================= Generate Systemd Service =======================
generate_systemd_service() {
    local output="$1"

    cat > "${output}" <<EOF
[Unit]
Description=${PKG_NAME}
Documentation=${PKG_HOMEPAGE}
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${PKG_USER}
Group=${PKG_GROUP}
ExecStart=${PKG_BINDIR}/${PKG_NAME}
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${PKG_DATADIR}
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF
}

# ======================= Build DEB Package =======================
build_deb() {
    local goarch="$1"
    local deb_arch="${goarch}"

    info "Building .deb package for ${deb_arch}..."

    local staging_dir="${BUILD_DIR}/deb/${deb_arch}"
    rm -rf "${staging_dir}"
    mkdir -p "${staging_dir}/DEBIAN"
    mkdir -p "${staging_dir}${PKG_BINDIR}"
    mkdir -p "${staging_dir}${PKG_ETCDIR}"
    mkdir -p "${staging_dir}${PKG_SYSCONFDIR}"
    mkdir -p "${staging_dir}${PKG_DATADIR}"

    # Copy binary
    local binary="${BUILD_DIR}/${PKG_NAME}_${GOOS}_${goarch}"
    cp "${binary}" "${staging_dir}${PKG_BINDIR}/${PKG_NAME}"
    chmod 755 "${staging_dir}${PKG_BINDIR}/${PKG_NAME}"

    # Copy config
    if [[ -n "${CONFIG_FILE}" && -f "${CONFIG_FILE}" ]]; then
        cp "${CONFIG_FILE}" "${staging_dir}${PKG_ETCDIR}/$(basename "${CONFIG_FILE}")"
        chmod 644 "${staging_dir}${PKG_ETCDIR}/$(basename "${CONFIG_FILE}")"
    fi

    # Copy systemd service
    if [[ "${SYSTEMD_SERVICE}" -eq 1 ]]; then
        local svc_src
        if [[ -n "${SYSTEMD_SERVICE_FILE:-}" && -f "${SYSTEMD_SERVICE_FILE}" ]]; then
            svc_src="${SYSTEMD_SERVICE_FILE}"
        else
            svc_src="${BUILD_DIR}/${PKG_NAME}.service"
            generate_systemd_service "${svc_src}"
        fi
        cp "${svc_src}" "${staging_dir}${PKG_SYSCONFDIR}/${PKG_NAME}.service"
        chmod 644 "${staging_dir}${PKG_SYSCONFDIR}/${PKG_NAME}.service"
    fi

    # Extra files
    for item in ${EXTRA_FILES}; do
        local src="${item%%:*}"
        local dst="${item##*:}"
        if [[ "${src}" == "${dst}" ]]; then
            dst="/${src}"
        fi
        mkdir -p "${staging_dir}/$(dirname "${dst}")"
        cp "${src}" "${staging_dir}${dst}"
    done

    # ---- debian/control ----
    local installed_size
    installed_size=$(du -sk "${staging_dir}" | awk '{print $1}')

    cat > "${staging_dir}/DEBIAN/control" <<EOF
Package: ${PKG_NAME}
Version: ${PKG_VERSION}
Section: ${PKG_SECTION}
Priority: ${PKG_PRIORITY}
Architecture: ${deb_arch}
Maintainer: ${PKG_MAINTAINER}
Vendor: ${PKG_VENDOR}
Homepage: ${PKG_HOMEPAGE}
License: ${PKG_LICENSE}
Installed-Size: ${installed_size}
Depends: libc6 (>= 2.17)
Recommends: systemd
Description: ${PKG_DESCRIPTION}
EOF

    # ---- Maintainer scripts ----
    # conffiles (only when config file is actually included in the package)
    if [[ -n "${CONFIG_FILE}" && -f "${CONFIG_FILE}" ]]; then
        echo "${PKG_ETCDIR}/$(basename "${CONFIG_FILE}")" > "${staging_dir}/DEBIAN/conffiles"
    fi

    # preinst
    cat > "${staging_dir}/DEBIAN/preinst" <<'PREINST_EOF'
#!/bin/bash
set -e
if ! getent group PKG_GROUP >/dev/null; then
    groupadd --system PKG_GROUP
fi
if ! getent passwd PKG_USER >/dev/null; then
    useradd --system --gid PKG_GROUP --shell /bin/false \
        --home-dir PKG_DATADIR --create-home PKG_USER
fi
mkdir -p PKG_DATADIR
chown PKG_USER:PKG_GROUP PKG_DATADIR
chmod 755 PKG_DATADIR
PREINST_EOF
    sed -i.bak "s|PKG_USER|${PKG_USER}|g; s|PKG_GROUP|${PKG_GROUP}|g; s|PKG_DATADIR|${PKG_DATADIR}|g" \
        "${staging_dir}/DEBIAN/preinst"
    rm -f "${staging_dir}/DEBIAN/preinst.bak"

    if [[ -n "${PREINST_SCRIPT}" && -f "${PREINST_SCRIPT}" ]]; then
        cat "${PREINST_SCRIPT}" >> "${staging_dir}/DEBIAN/preinst"
    fi
    chmod 755 "${staging_dir}/DEBIAN/preinst"

    # postinst
    cat > "${staging_dir}/DEBIAN/postinst" <<'POSTINST_EOF'
#!/bin/bash
set -e
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload
fi
echo "PKG_NAME has been installed successfully!"
echo "Start with: sudo systemctl start PKG_NAME"
echo "Enable on boot: sudo systemctl enable PKG_NAME"
POSTINST_EOF
    sed -i.bak "s|PKG_NAME|${PKG_NAME}|g" "${staging_dir}/DEBIAN/postinst"
    rm -f "${staging_dir}/DEBIAN/postinst.bak"

    if [[ -n "${POSTINST_SCRIPT}" && -f "${POSTINST_SCRIPT}" ]]; then
        cat "${POSTINST_SCRIPT}" >> "${staging_dir}/DEBIAN/postinst"
    fi
    chmod 755 "${staging_dir}/DEBIAN/postinst"

    # prerm
    cat > "${staging_dir}/DEBIAN/prerm" <<'PRERM_EOF'
#!/bin/bash
set -e
if command -v systemctl >/dev/null 2>&1; then
    if systemctl is-active --quiet PKG_NAME; then
        systemctl stop PKG_NAME
    fi
    if systemctl is-enabled --quiet PKG_NAME 2>/dev/null; then
        systemctl disable PKG_NAME
    fi
fi
PRERM_EOF
    sed -i.bak "s|PKG_NAME|${PKG_NAME}|g" "${staging_dir}/DEBIAN/prerm"
    rm -f "${staging_dir}/DEBIAN/prerm.bak"

    if [[ -n "${PRERM_SCRIPT}" && -f "${PRERM_SCRIPT}" ]]; then
        cat "${PRERM_SCRIPT}" >> "${staging_dir}/DEBIAN/prerm"
    fi
    chmod 755 "${staging_dir}/DEBIAN/prerm"

    # postrm
    cat > "${staging_dir}/DEBIAN/postrm" <<'POSTRM_EOF'
#!/bin/bash
set -e
case "$1" in
    purge)
        rm -f /etc/systemd/system/PKG_NAME.service
        if command -v systemctl >/dev/null 2>&1; then
            systemctl daemon-reload
        fi
        if getent passwd PKG_USER >/dev/null; then
            userdel PKG_USER
        fi
        if getent group PKG_GROUP >/dev/null; then
            groupdel PKG_GROUP
        fi
        echo "Note: Data directory PKG_DATADIR was not removed."
        echo "Remove manually: sudo rm -rf PKG_DATADIR"
        ;;
esac
POSTRM_EOF
    sed -i.bak "s|PKG_NAME|${PKG_NAME}|g; s|PKG_USER|${PKG_USER}|g; s|PKG_GROUP|${PKG_GROUP}|g; s|PKG_DATADIR|${PKG_DATADIR}|g" \
        "${staging_dir}/DEBIAN/postrm"
    rm -f "${staging_dir}/DEBIAN/postrm.bak"

    if [[ -n "${POSTRM_SCRIPT}" && -f "${POSTRM_SCRIPT}" ]]; then
        cat "${POSTRM_SCRIPT}" >> "${staging_dir}/DEBIAN/postrm"
    fi
    chmod 755 "${staging_dir}/DEBIAN/postrm"

    # Build .deb using ar (available on macOS/Linux)
    local deb_name="${PKG_NAME}_${PKG_VERSION}_${deb_arch}.deb"
    local deb_path="${OUTPUT_DIR}/${deb_name}"
    local work_dir="${BUILD_DIR}/deb_work/${deb_arch}"
    rm -rf "${work_dir}"
    mkdir -p "${work_dir}"

    info "Assembling ${deb_name}..."

    # debian-binary
    echo "2.0" > "${work_dir}/debian-binary"

    # control.tar.gz — force uid/gid 0 so installed files are root-owned
    tar -czf "${work_dir}/control.tar.gz" -C "${staging_dir}/DEBIAN" \
        --owner=0 --group=0 .

    # data.tar.gz — force uid/gid 0 so installed files are root-owned
    tar -czf "${work_dir}/data.tar.gz" -C "${staging_dir}" \
        --owner=0 --group=0 \
        --exclude="./DEBIAN" \
        .

    # Assemble .deb with ar
    local abs_deb_path
    abs_deb_path=$(cd "$(dirname "${deb_path}")" && pwd)/$(basename "${deb_path}")
    rm -f "${abs_deb_path}"

    # macOS ar is broken for non-Mach-O files, use Python assembler instead
    local os_name
    os_name=$(uname -s)

    if [[ "${os_name}" == "Darwin" ]]; then
        info "Using Python-based deb assembler (macOS)..."
        python3 -c "
import os, sys, time
deb_path = '${abs_deb_path}'
work_dir = '${work_dir}'
files = ['debian-binary', 'control.tar.gz', 'data.tar.gz']
with open(deb_path, 'wb') as f:
    f.write(b'!<arch>\n')
    for fn in files:
        path = os.path.join(work_dir, fn)
        data = open(path, 'rb').read()
        ar_name = fn + '/'
        header = '{:<16}{:<12}{:<6}{:<6}{:<8}{:<10}\x60\n'.format(
            ar_name, int(time.time()), 0, 0, '100644', len(data)
        ).encode()
        f.write(header)
        f.write(data)
        if len(data) % 2 != 0:
            f.write(b'\n')
"
    else
        ar -qc "${abs_deb_path}" \
            "${work_dir}/debian-binary" \
            "${work_dir}/control.tar.gz" \
            "${work_dir}/data.tar.gz"
    fi

    if [[ -f "${deb_path}" ]]; then
        local size
        size=$(du -h "${deb_path}" | awk '{print $1}')
        success "Created ${deb_path} (${size})"
    else
        error "Failed to create .deb package"
    fi
}

# ======================= Build RPM Package =======================
build_rpm() {
    local goarch="$1"

    # Map Go arch to RPM arch
    local rpm_arch
    case "${goarch}" in
        amd64)  rpm_arch="x86_64"  ;;
        arm64)  rpm_arch="aarch64" ;;
        386)    rpm_arch="i686"    ;;
        arm)    rpm_arch="armhfp"  ;;
        *)      rpm_arch="${goarch}" ;;
    esac

    # RPM Version must not contain '-'; normalize git describe output
    local rpm_version
    rpm_version="${PKG_VERSION//-/_}"
    if [[ "${rpm_version}" != "${PKG_VERSION}" ]]; then
        warn "Version '${PKG_VERSION}' contains '-'; normalizing to '${rpm_version}' for RPM."
    fi

    info "Building .rpm package for ${rpm_arch}..."

    local staging_dir="${BUILD_DIR}/rpm/${rpm_arch}"
    local rpmbuild_dir="${BUILD_DIR}/rpmbuild/${rpm_arch}"
    rm -rf "${staging_dir}" "${rpmbuild_dir}"
    mkdir -p "${staging_dir}"
    mkdir -p "${rpmbuild_dir}"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

    # Copy binary
    local binary="${BUILD_DIR}/${PKG_NAME}_${GOOS}_${goarch}"
    mkdir -p "${staging_dir}${PKG_BINDIR}"
    cp "${binary}" "${staging_dir}${PKG_BINDIR}/${PKG_NAME}"
    chmod 755 "${staging_dir}${PKG_BINDIR}/${PKG_NAME}"

    # Copy config
    if [[ -n "${CONFIG_FILE}" && -f "${CONFIG_FILE}" ]]; then
        mkdir -p "${staging_dir}${PKG_ETCDIR}"
        cp "${CONFIG_FILE}" "${staging_dir}${PKG_ETCDIR}/$(basename "${CONFIG_FILE}")"
        chmod 644 "${staging_dir}${PKG_ETCDIR}/$(basename "${CONFIG_FILE}")"
    fi

    # Copy systemd service
    if [[ "${SYSTEMD_SERVICE}" -eq 1 ]]; then
        local svc_src
        if [[ -n "${SYSTEMD_SERVICE_FILE:-}" && -f "${SYSTEMD_SERVICE_FILE}" ]]; then
            svc_src="${SYSTEMD_SERVICE_FILE}"
        else
            svc_src="${BUILD_DIR}/${PKG_NAME}.service"
            generate_systemd_service "${svc_src}"
        fi
        mkdir -p "${staging_dir}${PKG_SYSCONFDIR}"
        cp "${svc_src}" "${staging_dir}${PKG_SYSCONFDIR}/${PKG_NAME}.service"
        chmod 644 "${staging_dir}${PKG_SYSCONFDIR}/${PKG_NAME}.service"
    fi

    # Create tarball for rpmbuild
    local tar_name="${PKG_NAME}-${rpm_version}"
    tar -czf "${rpmbuild_dir}/SOURCES/${tar_name}.tar.gz" -C "${staging_dir}" .

    # ---- .spec file ----
    local postinst_content=""
    if [[ -n "${POSTINST_SCRIPT}" && -f "${POSTINST_SCRIPT}" ]]; then
        postinst_content=$(cat "${POSTINST_SCRIPT}")
    fi

    local preinst_content=""
    if [[ -n "${PREINST_SCRIPT}" && -f "${PREINST_SCRIPT}" ]]; then
        preinst_content=$(cat "${PREINST_SCRIPT}")
    fi

    local prerm_content=""
    if [[ -n "${PRERM_SCRIPT}" && -f "${PRERM_SCRIPT}" ]]; then
        prerm_content=$(cat "${PRERM_SCRIPT}")
    fi

    local postrm_content=""
    if [[ -n "${POSTRM_SCRIPT}" && -f "${POSTRM_SCRIPT}" ]]; then
        postrm_content=$(cat "${POSTRM_SCRIPT}")
    fi

    cat > "${rpmbuild_dir}/SPECS/${PKG_NAME}.spec" <<EOF
Name:           ${PKG_NAME}
Version:        ${rpm_version}
Release:        1%{?dist}
Summary:        ${PKG_DESCRIPTION}

License:        ${PKG_LICENSE}
URL:            ${PKG_HOMEPAGE}
Source0:        %{name}-%{version}.tar.gz

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root

Requires:       glibc >= 2.17
Recommends:     systemd

%description
${PKG_DESCRIPTION}

%prep
%setup -q -c

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}

# Copy all files from staging
cp -a * %{buildroot}/ || true

%pre
getent group ${PKG_GROUP} >/dev/null || groupadd -r ${PKG_GROUP}
getent passwd ${PKG_USER} >/dev/null || useradd -r -g ${PKG_GROUP} -s /bin/false -d ${PKG_DATADIR} ${PKG_USER}
mkdir -p ${PKG_DATADIR}
chown ${PKG_USER}:${PKG_GROUP} ${PKG_DATADIR}
${preinst_content}

%post
%systemd_post ${PKG_NAME}.service
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload
fi
${postinst_content}

%preun
%systemd_preun ${PKG_NAME}.service
${prerm_content}

%postun
%systemd_postun ${PKG_NAME}.service
if [ \$1 -eq 0 ]; then
    rm -f ${PKG_SYSCONFDIR}/${PKG_NAME}.service
    if command -v systemctl >/dev/null 2>&1; then
        systemctl daemon-reload
    fi
    getent passwd ${PKG_USER} >/dev/null && userdel ${PKG_USER}
    getent group ${PKG_GROUP} >/dev/null && groupdel ${PKG_GROUP}
fi
${postrm_content}

%files
%defattr(-,root,root,-)
${PKG_BINDIR}/${PKG_NAME}
%if %{?_systemd_unitdir:1}0
${PKG_SYSCONFDIR}/${PKG_NAME}.service
%endif
$(if [[ -n "${CONFIG_FILE}" && -f "${CONFIG_FILE}" ]]; then echo "${PKG_ETCDIR}/$(basename "${CONFIG_FILE}")"; fi)

%changelog
* $(date +'%a %b %d %Y') ${PKG_MAINTAINER} - ${rpm_version}-1
- Initial package
EOF

    # Build with rpmbuild
    if command -v rpmbuild &>/dev/null; then
        info "Building RPM with rpmbuild..."
        rpmbuild -bb \
            --define "_topdir ${rpmbuild_dir}" \
            --target "${rpm_arch}" \
            "${rpmbuild_dir}/SPECS/${PKG_NAME}.spec"

        # Find the built RPM and copy to output
        local built_rpm
        built_rpm=$(find "${rpmbuild_dir}/RPMS" -name "*.rpm" -type f | head -1)
        if [[ -n "${built_rpm}" && -f "${built_rpm}" ]]; then
            cp "${built_rpm}" "${OUTPUT_DIR}/"
            success "Created $(basename "${built_rpm}")"
        fi
    else
        error "rpmbuild is not available. Install it with: yum install rpm-build  OR  apt install rpm"
        return 1
    fi
}

# ======================= Main =======================
main() {
    info "==========================================="
    info "  Universal RPM/DEB Packaging Script"
    info "==========================================="

    # Check Go
    if ! command -v go &>/dev/null; then
        error "Go is not installed"
    fi

    # Auto-detect
    detect_name
    detect_version
    detect_config
    detect_systemd_service
    detect_maintainer
    detect_homepage
    detect_scripts
    detect_description

    # Update derived vars
    PKG_ETCDIR="/etc/${PKG_NAME}"
    PKG_DATADIR="/var/lib/${PKG_NAME}"
    PKG_USER="${PKG_NAME}"
    PKG_GROUP="${PKG_NAME}"

    # Parse arch list
    IFS=',' read -ra ARCH_LIST <<< "${ARCH}"

    info "Package:    ${PKG_NAME}"
    info "Version:    ${PKG_VERSION}"
    info "Format:     ${FORMAT}"
    info "Arch:       ${ARCH}"
    info "Output:     ${OUTPUT_DIR}"

    # Prepare dirs
    mkdir -p "${OUTPUT_DIR}" "${BUILD_DIR}"

    # Build binaries for each arch
    for goarch in "${ARCH_LIST[@]}"; do
        build_binary "${goarch}" "${BUILD_DIR}/${PKG_NAME}_${GOOS}_${goarch}"
    done

    # Build packages
    for goarch in "${ARCH_LIST[@]}"; do
        case "${FORMAT}" in
            deb)
                build_deb "${goarch}"
                ;;
            rpm)
                build_rpm "${goarch}"
                ;;
            all)
                build_deb "${goarch}"
                build_rpm "${goarch}"
                ;;
            *)
                error "Unknown format: ${FORMAT}. Use: deb, rpm, all"
                ;;
        esac
    done

    echo ""
    success "==========================================="
    success "  All packages built successfully!"
    success "==========================================="
    info "Output directory: ${OUTPUT_DIR}/"
    ls -lh "${OUTPUT_DIR}/" 2>/dev/null || true
}

main
