#!/bin/bash


# Original source is located at github.com/cugxuan/gonelist/scripts/install-release.sh

# If not specify, default meaning of return value:
# 0: Success
# 1: System error
# 2: Application error
# 3: Network error

# CLI arguments
PROXY=''
HELP=''
FORCE=''
CHECK=''
REMOVE=''
VERSION=''
VSRC_ROOT='/tmp/gonelist'
EXTRACT_ONLY=''
LOCAL=''
LOCAL_INSTALL=''
DIST_SRC='github'
ERROR_IF_UPTODATE=''

CUR_VER=""
NEW_VER=""
VDIS=''
TARFILE="/tmp/gonelist/gonelist.tar.gz"
GONELIST_RUNNING=0

CMD_INSTALL=""
CMD_UPDATE=""
SOFTWARE_UPDATED=0

SYSTEMCTL_CMD=$(command -v systemctl 2>/dev/null)
SERVICE_CMD=$(command -v service 2>/dev/null)

#######color code########
RED="31m"      # Error message
GREEN="32m"    # Success message
YELLOW="33m"   # Warning message
BLUE="36m"     # Info message


#########################
while [[ $# > 0 ]]; do
    case "$1" in
        -p|--proxy)
        PROXY="-x ${2}"
        shift # past argument
        ;;
        -h|--help)
        HELP="1"
        ;;
        -f|--force)
        FORCE="1"
        ;;
        -c|--check)
        CHECK="1"
        ;;
        --remove)
        REMOVE="1"
        ;;
        --version)
        VERSION="$2"
        shift
        ;;
        --extract)
        VSRC_ROOT="$2"
        shift
        ;;
        --extractonly)
        EXTRACT_ONLY="1"
        ;;
        -l|--local)
        LOCAL="$2"
        LOCAL_INSTALL="1"
        shift
        ;;
        --source)
        DIST_SRC="$2"
        shift
        ;;
        --errifuptodate)
        ERROR_IF_UPTODATE="1"
        ;;
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done

###############################
colorEcho(){
    echo -e "\033[${1}${@:2}\033[0m" 1>& 2
}

archAffix(){
    case "${1:-"$(uname -m)"}" in
        i686|i386)
            echo '386'
        ;;
        x86_64|amd64)
            echo 'amd64'
        ;;
        *armv7*|armv6l)
            echo 'arm'
        ;;
        *armv8*|aarch64)
            echo 'arm64'
        ;;
        *mips64le*)
            echo 'mips64le'
        ;;
        *mips64*)
            echo 'mips64'
        ;;
        *mipsle*)
            echo 'mipsle'
        ;;
        *mips*)
            echo 'mips'
        ;;
        *s390x*)
            echo 's390x'
        ;;
        ppc64le)
            echo 'ppc64le'
        ;;
        ppc64)
            echo 'ppc64'
        ;;
        *)
            return 1
        ;;
    esac

	return 0
}


# ex: VISD=amd64
download_gonelist(){
    rm -rf /tmp/gonelist
    mkdir -p /tmp/gonelist
    if [[ "${DIST_SRC}" == "GONELIST" ]]; then
        DOWNLOAD_LINK="https://gonelist.cugxuan.cn/d/gonelist-release/gonelist_linux_${VDIS}.tar.gz"
    else
        DOWNLOAD_LINK="https://github.com/cugxuan/gonelist/releases/download/${NEW_VER}/gonelist_linux_${VDIS}.tar.gz"
    fi
    colorEcho ${BLUE} "Downloading Gonelist: ${DOWNLOAD_LINK}"
    curl ${PROXY} -L -H "Cache-Control: no-cache" -o ${TARFILE} ${DOWNLOAD_LINK}
    if [ $? != 0 ];then
        colorEcho ${RED} "Failed to download! Please check your network or try again."
        return 3
    fi
    return 0
}

installSoftware(){
    COMPONENT=$1
    if [[ -n `command -v $COMPONENT` ]]; then
        return 0
    fi

    getPMT
    if [[ $? -eq 1 ]]; then
        colorEcho ${RED} "The system package manager tool isn't APT or YUM, please install ${COMPONENT} manually."
        return 1
    fi
    if [[ $SOFTWARE_UPDATED -eq 0 ]]; then
        colorEcho ${BLUE} "Updating software repo"
        $CMD_UPDATE
        SOFTWARE_UPDATED=1
    fi

    colorEcho ${BLUE} "Installing ${COMPONENT}"
    $CMD_INSTALL $COMPONENT
    if [[ $? -ne 0 ]]; then
        colorEcho ${RED} "Failed to install ${COMPONENT}. Please install it manually."
        return 1
    fi
    return 0
}

# return 1: not apt, yum, or zypper
getPMT(){
    if [[ -n `command -v apt-get` ]];then
        CMD_INSTALL="apt-get -y -qq install"
        CMD_UPDATE="apt-get -qq update"
    elif [[ -n `command -v yum` ]]; then
        CMD_INSTALL="yum -y -q install"
        CMD_UPDATE="yum -q makecache"
    elif [[ -n `command -v zypper` ]]; then
        CMD_INSTALL="zypper -y install"
        CMD_UPDATE="zypper ref"
    else
        return 1
    fi
    return 0
}

normalizeVersion() {
    if [ -n "$1" ]; then
        case "$1" in
            v*)
                echo "$1"
            ;;
            *)
                echo "v$1"
            ;;
        esac
    else
        echo ""
    fi
}

# 1: new gonelist. 0: no. 2: not installed. 3: check failed. 4: don't check.
getVersion(){
    if [[ -n "$VERSION" ]]; then
        NEW_VER="$(normalizeVersion "$VERSION")"
        return 4
    else
        #get the latest release
        TAG_URL="https://api.github.com/repos/cugxuan/gonelist/releases/latest"
        NEW_VER="$(normalizeVersion "$(curl ${PROXY} --retry 6 \
            -H "Accept: application/json" \
            -H "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:74.0) Gecko/20100101 Firefox/74.0" \
            -s "${TAG_URL}" --connect-timeout 20 | awk -F'[ "]+' '$0~"tag_name"{print $4;exit}' )")"

        [ ! -f /usr/bin/gonelist/gonelist ] && return 2
        VER="$(/usr/bin/gonelist/gonelist --version | awk '/Version/{print $NF;exit}')"
        RETVAL=$?
        CUR_VER="$(normalizeVersion "$VER")"

        if [[ $? -ne 0 ]] || [[ $NEW_VER == "" ]]; then
            colorEcho ${RED} "Failed to fetch release information. Please check your network or try again."
            return 3
        elif [[ $RETVAL -ne 0 ]];then
            return 2
        elif [[ $NEW_VER != $CUR_VER ]];then
            return 1
        fi
        return 0
    fi
}

stop_gonelist(){
    colorEcho ${BLUE} "Shutting down Gonelist service."
    if [[ -n "${SYSTEMCTL_CMD}" ]] || [[ -f "/lib/systemd/system/gonelist.service" ]] || [[ -f "/etc/systemd/system/gonelist.service" ]]; then
        ${SYSTEMCTL_CMD} stop gonelist
    elif [[ -n "${SERVICE_CMD}" ]] || [[ -f "/etc/init.d/gonelist" ]]; then
        ${SERVICE_CMD} gonelist stop
    fi
    if [[ $? -ne 0 ]]; then
        colorEcho ${YELLOW} "Failed to shutdown Gonelist service."
        return 2
    fi
    return 0
}

start_gonelist(){
    if [ -n "${SYSTEMCTL_CMD}" ] && [[ -f "/lib/systemd/system/gonelist.service" || -f "/etc/systemd/system/gonelist.service" ]]; then
        ${SYSTEMCTL_CMD} start gonelist
    elif [ -n "${SERVICE_CMD}" ] && [ -f "/etc/init.d/gonelist" ]; then
        ${SERVICE_CMD} gonelist start
    fi
    if [[ $? -ne 0 ]]; then
        colorEcho ${YELLOW} "Failed to start Gonelist service."
        return 2
    fi
    return 0
}


install_gonelist(){
    # Install gonelist binary and dist dir to /usr/local/gonelist/
    local file=$1
    local arch=$2
    mkdir -p '/etc/gonelist/' '/usr/local/gonelist/' && \
    [ -d /usr/local/gonelist/dist ] && rm -rf /usr/local/gonelist/dist
    tar zxf ${file} -C '/usr/local/gonelist/' --strip-components=1 && \
    mv /usr/local/gonelist/gonelist_linux_${arch} /usr/local/gonelist/gonelist && \
    chmod +x '/usr/local/gonelist/gonelist' || {
        colorEcho ${RED} "Failed to copy gonelist binary and resources."
        return 1
    }

    # Install gonelist server config to /etc/gonelist
    if [ ! -f '/etc/gonelist/config.yml' ]; then
        cp /usr/local/gonelist/config.yml /etc/gonelist/config.yml
        sed -ri '/dist_path/s#: "[^"]+#: "/usr/local/gonelist/dist/#' /etc/gonelist/config.yml
    fi
    
}


installInitScript(){
    if [[ -n "${SYSTEMCTL_CMD}" ]]; then
        if [[ ! -f "/etc/systemd/system/gonelist.service" && ! -f "/lib/systemd/system/gonelist.service" ]]; then
            cat>/etc/systemd/system/gonelist.service<<'EOF' 
[Unit]
Description=gonelist - Golang Onedrive List
Documentation=https://github.com/cugxuan/gonelist
After=network.target
Wants=network-online.target

[Service]
# If the version of systemd is 240 or above, then uncommenting Type=exec and commenting out Type=simple
#Type=exec
Type=simple
#User=root
NoNewPrivileges=yes
ExecStart=/usr/local/gonelist/gonelist --conf /etc/gonelist/config.yml
Restart=on-failure
RestartSec=4s
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target    
EOF
    systemctl enable gonelist.service
        fi
    elif [[ -n "${SERVICE_CMD}" ]] && [[ ! -f "/etc/init.d/gonelist" ]]; then
        installSoftware 'daemon' && \
        echo todo rc.d && \
        chmod +x '/etc/init.d/gonelist' && \
        update-rc.d gonelist defaults
    fi
}

Help(){
  cat - 1>& 2 << EOF
./install-release.sh [-h] [-c] [--remove] [-p proxy] [-f] [--version vx.y.z] [-l file]
  -h, --help            Show help
  -p, --proxy           To download through a proxy server, use -p socks5://127.0.0.1:1080 or -p http://127.0.0.1:3128 etc
  -f, --force           Force install
  -l, --local           Install from a local file
      --remove          Remove installed gonelist
  -c, --check           Check for update
EOF
}

remove(){
    if [[ -n "${SYSTEMCTL_CMD}" ]] && [[ -f "/etc/systemd/system/gonelist.service" ]];then
        if pgrep "gonelist" > /dev/null ; then
            stop_gonelist
        fi
        systemctl disable gonelist.service
        rm -rf "/usr/local/gonelist" "/etc/systemd/system/gonelist.service"
        if [[ $? -ne 0 ]]; then
            colorEcho ${RED} "Failed to remove gonelist."
            return 0
        else
            colorEcho ${GREEN} "Removed gonelist successfully."
            colorEcho ${BLUE} "If necessary, please remove configuration file and log file manually."
            return 0
        fi
    elif [[ -n "${SYSTEMCTL_CMD}" ]] && [[ -f "/lib/systemd/system/gonelist.service" ]];then
        if pgrep "gonelist" > /dev/null ; then
            stop_gonelist
        fi
        systemctl disable gonelist.service
        rm -rf "/usr/local/gonelist" "/lib/systemd/system/gonelist.service"
        if [[ $? -ne 0 ]]; then
            colorEcho ${RED} "Failed to remove gonelist."
            return 0
        else
            colorEcho ${GREEN} "Removed gonelist successfully."
            colorEcho ${BLUE} "If necessary, please remove configuration file and log file manually."
            return 0
        fi
    elif [[ -n "${SERVICE_CMD}" ]] && [[ -f "/etc/init.d/gonelist" ]]; then
        if pgrep "gonelist" > /dev/null ; then
            stop_gonelist
        fi
        rm -rf "/usr/local/gonelist" "/etc/init.d/gonelist"
        if [[ $? -ne 0 ]]; then
            colorEcho ${RED} "Failed to remove gonelist."
            return 0
        else
            colorEcho ${GREEN} "Removed gonelist successfully."
            colorEcho ${BLUE} "If necessary, please remove configuration file and log file manually."
            return 0
        fi
    else
        colorEcho ${YELLOW} "gonelist not found."
        return 0
    fi
}

checkUpdate(){
    echo "Checking for update."
    VERSION=""
    getVersion
    RETVAL="$?"
    if [[ $RETVAL -eq 1 ]]; then
        colorEcho ${BLUE} "Found new version ${NEW_VER} for gonelist.(Current version:$CUR_VER)"
    elif [[ $RETVAL -eq 0 ]]; then
        colorEcho ${BLUE} "No new version. Current version is ${NEW_VER}."
    elif [[ $RETVAL -eq 2 ]]; then
        colorEcho ${YELLOW} "No gonelist installed."
        colorEcho ${BLUE} "The newest version for gonelist is ${NEW_VER}."
    fi
    return 0
}

main(){
    #helping information
    [[ "$HELP" == "1" ]] && Help && return
    [[ "$CHECK" == "1" ]] && checkUpdate && return
    [[ "$REMOVE" == "1" ]] && remove && return

    local ARCH=$(uname -m)
    VDIS="$(archAffix)"

    # extract local file
    if [[ $LOCAL_INSTALL -eq 1 ]]; then
        colorEcho ${YELLOW} "Installing Gonelist via local file. Please make sure the file is a valid Gonelist package, as we are not able to determine that."
        NEW_VER=local
        rm -rf /tmp/gonelist
        TARFILE="$LOCAL"
        #FILEVDIS=`ls /tmp/gonelist |grep gonelist-v |cut -d "-" -f4`
        #SYSTEM=`ls /tmp/gonelist |grep gonelist-v |cut -d "-" -f3`
        #if [[ ${SYSTEM} != "linux" ]]; then
        #    colorEcho ${RED} "The local gonelist can not be installed in linux."
        #    return 1
        #elif [[ ${FILEVDIS} != ${VDIS} ]]; then
        #    colorEcho ${RED} "The local gonelist can not be installed in ${ARCH} system."
        #    return 1
        #else
        #    NEW_VER=`ls /tmp/gonelist |grep gonelist-v |cut -d "-" -f2`
        #fi
    else
        # download via network and extract
        installSoftware "curl" || return $?
        getVersion
        RETVAL="$?"
        if [[ $RETVAL == 0 ]] && [[ "$FORCE" != "1" ]]; then
            colorEcho ${BLUE} "Latest version ${CUR_VER} is already installed."
            if [ -n "${ERROR_IF_UPTODATE}" ]; then
              return 10
            fi
            return
        elif [[ $RETVAL == 3 ]]; then
            return 3
        else
            colorEcho ${BLUE} "Installing Gonelist ${NEW_VER} on ${ARCH}"
            download_gonelist || return $?
        fi
    fi
    
    if [ -n "${EXTRACT_ONLY}" ]; then
        colorEcho ${BLUE} "Extracting gonelist package to ${VSRC_ROOT}."
        if tar -zxf "${TARFILE}" --strip-components=1 -C ${VSRC_ROOT}; then
            colorEcho ${GREEN} "gonelist extracted to ${VSRC_ROOT}, and exiting..."
            return 0
        else
            colorEcho ${RED} "Failed to extract gonelist."
            return 2
        fi
    fi

    if pgrep "gonelist" > /dev/null ; then
        GONELIST_RUNNING=1
        stop_gonelist
    fi
    install_gonelist "${TARFILE}" "${VDIS}" || return $?
    installInitScript || return $?
    if [[ ${GONELIST_RUNNING} -eq 1 ]];then
        colorEcho ${BLUE} "Restarting gonelist service."
        stop_gonelist
        start_gonelist
    fi
    colorEcho ${GREEN} "gonelist ${NEW_VER} is installed."
    rm -rf /tmp/gonelist
    return 0
}

main