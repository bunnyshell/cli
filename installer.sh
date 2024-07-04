#!/bin/sh
# Github release installer
#
# Heavily inspired by:
# - https://github.com/goreleaser/godownloader
#
# Shoutout to:
# - https://www.codetinkerer.com/is-shell-portable/
# - https://google.github.io/styleguide/shellguide.html although the guide is for BASH not SH

# run in subshell
(
    # config
    BINARY=bns
    BREW_TAP=bunnyshell/tap/bunnyshell-cli
    REPOSITORY=bunnyshell/cli

    # install options
    RELEASE_VERSION=${RELEASE_VERSION:-'latest'}
    PREFER_BREW=${PREFER_BREW:-false}
    INSTALL_PATH=${INSTALL_PATH:-.}
    SUDO_INSTALL=${SUDO_INSTALL:-false}

    if [ "${DEBUG_INSTALLER}" = true ]; then
        set -x
    fi

    # main functions
    main() {
        if has_brew; then
            install_with_brew
        elif is_command "${BINARY}"; then
            if has_correct_binary; then
                upgrade;
            else
                throw "There is another binary called '${BINARY}' in your PATH"
            fi
        else
            install
        fi
    }
    has_brew() {
        if [ "${PREFER_BREW}" != true ]; then
            return 1
        fi

        is_command brew;
    }
    install_with_brew() {
        log_info 'Using detected brew installation'

        set -x

        brew install "${BREW_TAP}"
    }
    upgrade() {
        INSTALLED_VERSION=$(show_binary_version "${BINARY}")

        log_info "Already installed. ${INSTALLED_VERSION}"
        tmp_install=$(fetch_binary);EXIT_CODE=$?
        if [ "${EXIT_CODE}" -ne 0 ]; then
            issue 'Something went wrong with the installation'
        fi

        UPGRADED_VERSION=$(
            show_binary_version "${tmp_install}"
        );EXIT_CODE=$?

        if [ "${EXIT_CODE}" -ne 0 ]; then
            log_error 'Something went wrong with the installation'
            issue "Consider manually checking ${tmp_install}"
        fi

        if [ "${INSTALLED_VERSION}" = "${UPGRADED_VERSION}" ]; then
            log_info 'You are using the latest version'
            rm -r "${tmp_install}"
            exit 0
        fi

        release "${tmp_install}"
    }

    install() {
        log_info 'Fresh install required...'

        tmp_install=$(fetch_binary);EXIT_CODE=$?
        if [ "${EXIT_CODE}" -ne 0 ]; then
            issue 'Something went wrong with the installation'
        fi

        release "${tmp_install}"
    }

    fetch_binary() {
        if ! is_command tar; then
            throw 'You need tar binary to unpack github archive'
        fi

        version=$(github_release_version);EXIT_CODE=$?
        if [ "${EXIT_CODE}" -ne 0 ]; then
            throw "Installer could not determine version for ${RELEASE_VERSION}."
        fi

        os=$(goreleaser_os);EXIT_CODE=$?
        if [ "${EXIT_CODE}" -ne 0 ]; then
            issue "Installer cannot handle OS ${os}"
        fi

        arch=$(goreleaser_arch);EXIT_CODE=$?
        if [ "${EXIT_CODE}" -ne 0 ]; then
            issue "Installer cannot handle OS ${os} and ARCH ${arch}"
        fi

        source="https://github.com/${REPOSITORY}/releases/download/v${version}/${BINARY}_${version}_${os}_${arch}.tar.gz"

        github_download "$source"
    }

    release() {
        tmp_install=$1
        maybe_sudo=''
        if [ "${SUDO_INSTALL}" = true ]; then
            maybe_sudo='sudo'
            log_info 'Detected SUDO_INSTALL. Password might be required...'
        fi

        install_file="${INSTALL_PATH}/${BINARY}"

        RESULT=$(
            ${maybe_sudo} mv "${tmp_install}" "${install_file}" \
                \
                2>&1 \
            ;
        );EXIT_CODE=$?

        if [ "${EXIT_CODE}" -ne 0 ]; then
            throw "${RESULT}"
        fi

        log_info "Installed '${BINARY}' in ${INSTALL_PATH}"
        log_info "$(show_binary_version "${install_file}")"

        if ! is_command "${BINARY}"; then
            absolute_install_path=$(realpath "${INSTALL_PATH}")
            update_path=$(
                printf 'You will need to update your %sPATH variable or run "mv %s %s"' '$' "${absolute_install_path}/${BINARY}" '/usr/local/bin'
            )
            log_info "${update_path}"
        fi
    }

    # binary functions
    has_correct_binary() {
        HEADER=$(
            ${BINARY} --help | head -1
        );EXIT_CODE=$?

        if [ "${EXIT_CODE}" -ne 0 ]; then
            return 1
        fi

        case $HEADER in
            ("Bunnyshell CLI"*) return 0 ;;
            *) return 1 ;;
        esac
    }
    show_binary_version() {
        $1 version --client=true
    }

    # goreleaser binary naming from .goreleaser.yaml
    goreleaser_os() {
        os=$(uname -s | tr '[:upper:]' '[:lower:]')
        case $os in
            darwin) echo "Darwin" ;;
            linux) echo "Linux" ;;
            *) echo "${os}"; return 1;
        esac
    }
    goreleaser_arch() {
        arch=$(uname -m)
        case $arch in
            x86_64) echo 'x86_64' ;;
            i386) echo 'i386' ;;
            aarch64) echo 'arm64' ;;
            arm64) echo 'arm64' ;;
            *) echo "${arch}"; return 1;
        esac
    }

    # fetcher logic
    fetch() {
        source_url=$1
        tmp=$(tmp_file)
        echo "${tmp}"

        download "$source_url" "${tmp}"
    }
    fetch_response() {
        tmp=$(fetch "$1");
        cat "${tmp}"
        rm -rf "${tmp}"
    }
    download() {
        if is_command curl; then
            download_with_curl "$1" "$2"
        elif is_command wget; then
            download_with_wget "$1" "$2"
        else
            throw 'You need either curl or wget to download remote files'
        fi
    }
    download_with_curl() {
        source_url=$1
        local_file=$2

        RESPONSE=$(
            curl "$source_url" --silent --location --fail \
                --output "$local_file"\
                --header 'Accept: application/json' \
            ;
        );EXIT_CODE=$?

        if [ "${EXIT_CODE}" -ne 0 ]; then
            throw "Could not download ${source_url} to ${local_file}: ${RESPONSE}"
        fi

        return 0
    }
    download_with_wget() {
        source_url=$1
        local_file=$2

        RESPONSE=$(
            wget "$source_url" --quiet \
                --output "$local_file" \
                --header 'Accept: application/json' \
            ;
        );EXIT_CODE=$?

        if [ "${EXIT_CODE}" -ne 0 ]; then
            throw "Could not download ${source_url} to ${local_file}: ${RESPONSE}"
        fi

        return 0
    }

    # github operations
    github_extract_binary() {
        archive=$1

        extract_dir=$(tmp_dir)

        result=$(
            tar --extract --gzip \
                --directory "${extract_dir}" \
                --file "${archive}" \
                ${BINARY} \
                \
                2>&1 \
            ;
        );EXIT_CODE=$?

        if [ "${EXIT_CODE}" -ne 0 ]; then
            issue "Could not extract ${BINARY}: ${result}"
        fi

        tmp_binary=$(tmp_file)
        mv "${extract_dir}/${BINARY}" "${tmp_binary}"
        rm -r "${extract_dir}"
        echo "${tmp_binary}"
    }
    github_download() {
        tmp_archive=$(fetch "$1") || throw 'Could not download archive'
        tmp_binary=$(github_extract_binary "${tmp_archive}")
        rm -rf "${tmp_archive}"

        echo "${tmp_binary}"
    }
    github_release_version() {
        json=$(fetch_response "https://github.com/${REPOSITORY}/releases/${RELEASE_VERSION}")

        test -z "${json}" && return 1
        version=$(
            echo "${json}" \
                | tr -s '\n' ' ' \
                | sed 's/.*"tag_name":"v//' \
                | sed 's/".*//' \
            ;
        )

        test -z "${version}" && return 1
        echo "${version}"
    }

    # helpers
    tmp_dir() {
        mktemp --quiet --directory
    }
    tmp_file() {
        mktemp --quiet
    }
    issue() {
        log_error "$@"
        log_error "Let us know: https://github.com/${REPOSITORY}/issues"

        exit 2
    }
    throw() {
        log_error "$@"

        exit 1
    }
    log_error() {
        echo '[ERROR]' "$@" >&2
    }
    log_info() {
        echo '[INFO]' "$@"
    }

    is_command() {
        command -v "$1" > /dev/null
    }

    main "$@"
)
