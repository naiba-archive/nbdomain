#!/bin/bash

#======================================================
#   System Required: CentOS 7+ / Debian 8+ / Ubuntu 16+
#   Description: 奶霸米表管理脚本
#   version: v1.0.0
#   Author: 奶霸
#   Blog: https://nai.ba
#   Github - nbdomain: https://github.com/naiba/nbdomain-theme
#======================================================

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

version="v1.0.0"

# check root
[[ $EUID -ne 0 ]] && echo -e "${red}错误: ${plain} 必须使用root用户运行此脚本！\n" && exit 1

# check os
if [[ -f /etc/redhat-release ]]; then
    release="centos"
elif cat /etc/issue | grep -Eqi "debian"; then
    release="debian"
elif cat /etc/issue | grep -Eqi "ubuntu"; then
    release="ubuntu"
elif cat /etc/issue | grep -Eqi "centos|red hat|redhat"; then
    release="centos"
elif cat /proc/version | grep -Eqi "debian"; then
    release="debian"
elif cat /proc/version | grep -Eqi "ubuntu"; then
    release="ubuntu"
elif cat /proc/version | grep -Eqi "centos|red hat|redhat"; then
    release="centos"
else
    echo -e "${red}未检测到系统版本，请联系脚本作者！${plain}\n" && exit 1
fi

os_version=""

# os version
if [[ -f /etc/os-release ]]; then
    os_version=$(awk -F'[= ."]' '/VERSION_ID/{print $3}' /etc/os-release)
fi
if [[ -z "$os_version" && -f /etc/lsb-release ]]; then
    os_version=$(awk -F'[= ."]+' '/DISTRIB_RELEASE/{print $2}' /etc/lsb-release)
fi

if [[ x"${release}" == x"centos" ]]; then
    if [[ ${os_version} -le 6 ]]; then
        echo -e "${red}请使用 CentOS 7 或更高版本的系统！${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"ubuntu" ]]; then
    if [[ ${os_version} -lt 16 ]]; then
        echo -e "${red}请使用 Ubuntu 16 或更高版本的系统！${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"debian" ]]; then
    if [[ ${os_version} -lt 8 ]]; then
        echo -e "${red}请使用 Debian 8 或更高版本的系统！${plain}\n" && exit 1
    fi
fi

confirm() {
    if [[ $# > 1 ]]; then
        echo && read -p "$1 [默认$2]: " temp
        if [[ x"${temp}" == x"" ]]; then
            temp=$2
        fi
    else
        read -p "$1 [y/n]: " temp
    fi
    if [[ x"${temp}" == x"y" || x"${temp}" == x"Y" ]]; then
        return 0
    else
        return 1
    fi
}

confirm_restart() {
    confirm "是否重启米表" "y"
    if [[ $? == 0 ]]; then
        docker-compose down
        docker-compose up -d
    else
        show_menu
    fi
}

before_show_menu() {
    echo && echo -n -e "${yellow}* 按回车返回主菜单 *${plain}" && read temp
    show_menu
}

install_base() {
    (command -v git >/dev/null 2>&1 && command -v curl >/dev/null 2>&1 && command -v wget >/dev/null 2>&1) ||
        (command -v yum >/dev/null 2>&1 && yum install curl wget git -y) ||
        (command -v apt >/dev/null 2>&1 && apt install curl wget git -y) ||
        (command -v apt-get >/dev/null 2>&1 && apt-get install curl wget git -y)
}

install_soft() {
    (command -v $1 >/dev/null 2>&1) ||
        (command -v yum >/dev/null 2>&1 && yum install $1 -y) ||
        (command -v apt >/dev/null 2>&1 && apt install $1 -y) ||
        (command -v apt-get >/dev/null 2>&1 && apt-get install $1 -y)
}

install() {
    install_base

    # 上传文件夹
    mkdir -p data/nbdomain/upload/logo
    chmod 777 -R data/nbdomain/upload

    command -v docker >/dev/null 2>&1
    if [[ $? != 0 ]]; then
        echo -e "正在安装 Docker"
        bash <(curl -sL https://get.docker.com) >/dev/null 2>&1
        if [[ $? != 0 ]]; then
            echo -e "${red}下载脚本失败，请检查本机能否连接 get.docker.com${plain}"
            return 0
        fi
        systemctl enable docker.service
        systemctl start docker.service
        echo -e "${green}Docker${plain} 安装成功"
    fi

    command -v docker-compose >/dev/null 2>&1
    if [[ $? != 0 ]]; then
        echo -e "正在安装 Docker Compose"
        curl -L "https://github.com/docker/compose/releases/download/1.25.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose >/dev/null 2>&1 &&
            chmod +x /usr/local/bin/docker-compose
        if [[ $? != 0 ]]; then
            echo -e "${red}下载脚本失败，请检查本机能否连接 github.com${plain}"
            return 0
        fi
        echo -e "${green}Docker Compose${plain} 安装成功"
    fi

    modify_config 0
    rebuild 0
    start 0

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

rebuild() {
    echo -e "> 构建系统镜像"
    docker-compose build --no-cache nbdomain >/dev/null 2>&1
    if [[ $? == 0 ]]; then
        echo -e "系统镜像 ${green}构建成功${plain}"
    else
        echo -e "系统镜像 ${red}构建失败${plain}"
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

modify_config() {
    echo -e "> 修改系统配置"
    mkdir -p data/nbdomain/
    mkdir -p data/caddy/
    cp -rf ./config.yaml data/nbdomain/config.yaml
    cp -rf ./Caddyfile data/caddy/Caddyfile
    read -p "请输入管理后台域名: " domain &&
        read -p "请输入管理员邮箱: " admin_email &&
        read -p "请输入 reCAPTCHA Key: " recaptcha_key
    read -p "请输入 reCAPTCHA Secret: " recaptcha_secret
    if [[ -z "${domain}" || -z "${recaptcha_key}" || -z "${recaptcha_secret}" || -z "${admin_email}" ]]; then
        echo -e "${red}所有选项都不能为空${plain}"
        before_show_menu
        return 1
    fi

    sed -i "s/^example.com/${domain}/" data/caddy/Caddyfile
    sed -i "s/master@example.com/${admin_email}/" data/caddy/Caddyfile
    sed -i "s/example.com/${domain}/" data/nbdomain/config.yaml
    sed -i "s/recaptcha_site_key/${recaptcha_key}/" data/nbdomain/config.yaml
    sed -i "s/recaptcha_site_secret/${recaptcha_secret}/" data/nbdomain/config.yaml
    echo -e "系统配置 ${green}修改成功，请重新编译镜像并启动${plain}"

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

start() {
    docker-compose up -d
    if [[ $? == 0 ]]; then
        echo -e "${green}奶霸米表 启动成功${plain}"
    else
        echo -e "${red}米表启动失败，请稍后查看日志信息${plain}"
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

stop() {
    docker-compose down
    if [[ $? == 0 ]]; then
        echo -e "${green}奶霸米表 停止成功${plain}"
    else
        echo -e "${red}米表停止失败，请稍后查看日志信息${plain}"
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

restart() {
    docker-compose restart
    if [[ $? == 0 ]]; then
        echo -e "${green}奶霸米表 重启成功${plain}"
    else
        echo -e "${red}米表重启失败，可能是因为启动时间超过了两秒，请稍后查看日志信息${plain}"
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

show_log() {
    docker-compose logs -f

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

check_install() {
    command -v docker >/dev/null 2>&1 && command -v docker-compose >/dev/null 2>&1 && command -v git >/dev/null 2>&1
    if [[ $? != 0 ]]; then
        echo ""
        echo -e "${red}请先安装米表${plain}"
        if [[ $# == 0 ]]; then
            before_show_menu
        fi
        return 1
    else
        return 0
    fi
}

show_usage() {
    echo "奶霸米表 管理脚本使用方法: "
    echo "------------------------------------------"
    echo "./nbdomain.sh              - 显示管理菜单"
    echo "./nbdomain.sh install      - 安装"
    echo "./nbdomain.sh rebuild      - 重建镜像"
    echo "./nbdomain.sh config       - 修改配置"
    echo "./nbdomain.sh start        - 启动"
    echo "./nbdomain.sh stop         - 停止"
    echo "./nbdomain.sh restart      - 重启"
    echo "./nbdomain.sh log          - 查看日志"
    echo "------------------------------------------"
}

show_menu() {
    echo -e "
    ${green}奶霸米表管理脚本${plain} ${red}${version}${plain}
    --- https://nai.ba ---
    ${green}0.${plain} 退出脚本
    ————————————————
    ${green}1.${plain} 安装
    ${green}2.${plain} 重建镜像
    ————————————————
    ${green}3.${plain} 修改配置
    ————————————————
    ${green}4.${plain} 启动
    ${green}5.${plain} 停止
    ${green}6.${plain} 重启
    ${green}7.${plain} 查看日志
    "
    echo && read -p "请输入选择 [0-14]: " num

    case "${num}" in
    0)
        exit 0
        ;;
    1)
        install
        ;;
    2)
        check_install && rebuild
        ;;
    3)
        check_install && modify_config
        ;;
    4)
        check_install && start
        ;;
    5)
        check_install && stop
        ;;
    6)
        check_install && restart
        ;;
    7)
        check_install && show_log
        ;;
    *)
        echo -e "${red}请输入正确的数字 [0-7]${plain}"
        ;;
    esac
}

if [[ $# > 0 ]]; then
    case $1 in
    "install")
        install 0
        ;;
    "rebuild")
        check_install 0 && rebuild 0
        ;;
    "config")
        check_install 0 && modify_config 0
        ;;
    "start")
        check_install 0 && start 0
        ;;
    "stop")
        check_install 0 && stop 0
        ;;
    "restart")
        check_install 0 && restart 0
        ;;
    "log")
        check_install 0 && show_log 0
        ;;
    *) show_usage ;;
    esac
else
    show_menu
fi
