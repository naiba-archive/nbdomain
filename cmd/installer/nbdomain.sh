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
conf_dir="/etc/nbdomain/"
conf_path="${conf_dir}nbdomain.conf"

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

# -1: 未安装, 0: 已运行, 1: 未运行
sprov_ui_status=-1

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
        restart
    else
        show_menu
    fi
}

before_show_menu() {
    echo && echo -n -e "${yellow}按回车返回主菜单: ${plain}" && read temp
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
    echo -e "正在安装 Docker"
    bash <(curl -fsSL https://get.docker.com -o get-docker.sh)
    if [[ $? != 0 ]]; then
        echo -e "${red}下载脚本失败，请检查本机能否连接 get.docker.com${plain}"
        return 0
    fi
    echo -e "${green}Docker${plain} 安装成功"

    echo -e "正在安装 Docker Compose"
    curl -L "https://github.com/docker/compose/releases/download/1.25.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
        chmod +x /usr/local/bin/docker-compose
    if [[ $? != 0 ]]; then
        echo -e "${red}下载脚本失败，请检查本机能否连接 github.com${plain}"
        return 0
    fi
    echo -e "${green}Docker Compose${plain} 安装成功"

    echo -e "构建系统镜像"
    docker build --no-cache 

}

update() {
    confirm "本功能会强制重装当前最新版，数据不会丢失，是否继续?" "n"
    if [[ $? != 0 ]]; then
        echo -e "${red}已取消${plain}"
        if [[ $# == 0 ]]; then
            before_show_menu
        fi
        return 0
    fi
    install_base
    bash <(curl -L -s https://github.com/seedc/nbdomain/blob/install.sh)
    if [[ $? == 0 ]]; then
        if [[ $# == 0 ]]; then
            restart
        else
            restart 0
        fi
    fi
}

uninstall() {
    confirm "确定要卸载米表吗?" "n"
    if [[ $? != 0 ]]; then
        if [[ $# == 0 ]]; then
            show_menu
        fi
        return 0
    fi
    systemctl stop nbdomain
    systemctl disable nbdomain
    rm /etc/systemd/system/nbdomain.service -f
    systemctl daemon-reload
    systemctl reset-failed
    rm /etc/nbdomain/ -rf
    rm /usr/local/nbdomain/ -rf

    echo ""
    echo -e "${gree}卸载成功${plain}，感谢你的使用，如果你有更多的建议或意见，可以在以下地方进行讨论: "
    echo ""
    echo -e "Telegram 群组: ${green}https://t.me/sprov_blog${plain}"
    echo -e "Github issues: ${green}https://github.com/naiba/nbdomain-theme/issues${plain}"
    echo -e "博客: ${green}https://nai.ba/nbdomain${plain}"

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

modify_user() {
    echo && read -p "请输入用户名: " user
    read -p "请输入密码: " pwd
    if [[ -z "${user}" || -z "${pwd}" ]]; then
        echo -e "${red}用户名和密码不能为空${plain}"
        before_show_menu
        return 1
    fi
    sed -i "s/^username=.*/username=${user}/" ${conf_path}
    sed -i "s/^password=.*/password=${pwd}/" ${conf_path}
    confirm_restart
}

modify_port() {
    echo && read -p "输入米表新端口 [建议10000-65535]: " port
    if [[ -z "${port}" ]]; then
        echo -e "${red}未输入端口${plain}"
        before_show_menu
        return 1
    fi
    sed -i "s/^port=.*/port=${port}/" ${conf_path}
    confirm_restart
}

modify_config() {
    install_soft vim
    echo -e "----------------------------------------------------"
    echo -e "                vim 使用方法说明: "
    echo -e "首先按字母 ${red}i${plain} 进入 ${green}[编辑模式]${plain}"
    echo -e "${green}[编辑模式]${plain} 下，方向键移动光标，和平常编辑文本的习惯一样"
    echo -e "编辑完毕后，按 ${red}Esc${plain} 键退出 ${green}[编辑模式]${plain}"
    echo -e "最后按 ${red}:wq${plain} 保存文件并退出 vim ${yellow}(注意有个英文冒号)${plain}"
    echo -e "----------------------------------------------------"
    echo -e -n "${greed}将会使用 vim 进行编辑，按回车继续，或输入 n 返回: ${plain}"
    read temp
    if [[ x"${temp}" == x"n" || x"${temp}" == x"N" ]]; then
        show_menu
        return 0
    fi
    vim ${conf_path}
    confirm_restart
}

start() {
    check_status
    if [[ $? == 0 ]]; then
        echo ""
        echo -e "${green}米表已运行，无需再次启动，如需重启请选择重启${plain}"
    else
        systemctl start nbdomain
        sleep 2
        check_status
        if [[ $? == 0 ]]; then
            echo -e "${green}nbdomain 启动成功${plain}"
        else
            echo -e "${red}米表启动失败，可能是因为启动时间超过了两秒，请稍后查看日志信息${plain}"
        fi
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

stop() {
    check_status
    if [[ $? == 1 ]]; then
        echo ""
        echo -e "${green}米表已停止，无需再次停止${plain}"
    else
        systemctl stop nbdomain
        sleep 2
        check_status
        if [[ $? == 1 ]]; then
            echo -e "${green}nbdomain 停止成功${plain}"
        else
            echo -e "${red}米表停止失败，可能是因为停止时间超过了两秒，请稍后查看日志信息${plain}"
        fi
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

restart() {
    systemctl restart nbdomain
    sleep 2
    check_status
    if [[ $? == 0 ]]; then
        echo -e "${green}nbdomain 重启成功${plain}"
    else
        echo -e "${red}米表重启失败，可能是因为启动时间超过了两秒，请稍后查看日志信息${plain}"
    fi
    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

enable() {
    systemctl enable nbdomain
    if [[ $? == 0 ]]; then
        echo -e "${green}nbdomain 设置开机自启成功${plain}"
    else
        echo -e "${red}nbdomain 设置开机自启失败${plain}"
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

disable() {
    systemctl disable nbdomain
    if [[ $? == 0 ]]; then
        echo -e "${green}nbdomain 取消开机自启成功${plain}"
    else
        echo -e "${red}nbdomain 取消开机自启失败${plain}"
    fi

    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

show_log() {
    systemctl status nbdomain -l
    if [[ $# == 0 ]]; then
        before_show_menu
    fi
}

install_bbr() {
    bash <(curl -L -s https://github.com/seedc/nbdomain/blob/bbr.sh)
    if [[ $? == 0 ]]; then
        echo ""
        echo -e "${green}安装 bbr 成功${plain}"
    else
        echo ""
        echo -e "${red}下载 bbr 安装脚本失败，请检查本机能否连接 Github${plain}"
    fi

    before_show_menu
}

update_shell() {
    wget -O /usr/bin/nbdomain -N --no-check-certificate https://github.com/seedc/nbdomain/blob/nbdomain.sh
    if [[ $? != 0 ]]; then
        echo ""
        echo -e "${red}下载脚本失败，请检查本机能否连接 Github${plain}"
        before_show_menu
    else
        chmod +x /usr/bin/nbdomain
        echo -e "${green}升级脚本成功，请重新运行脚本${plain}" && exit 0
    fi
}

check_install() {
    command -v git >/dev/null 2>&1 && command -v curl >/dev/null 2>&1 && command -v wget >/dev/null 2>&1
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

show_status() {
    check_status
    case $? in
    0)
        echo -e "米表状态: ${green}已运行${plain}"
        show_enable_status
        ;;
    1)
        echo -e "米表状态: ${yellow}未运行${plain}"
        show_enable_status
        ;;
    2)
        echo -e "米表状态: ${red}未安装${plain}"
        ;;
    esac
}

show_enable_status() {
    check_enabled
    if [[ $? == 0 ]]; then
        echo -e "是否开机自启: ${green}是${plain}"
    else
        echo -e "是否开机自启: ${red}否${plain}"
    fi
}

show_usage() {
    echo "奶霸米表 管理脚本使用方法: "
    echo "------------------------------------------"
    echo "./nbdomain.sh              - 显示管理菜单 (功能更多)"
    echo "./nbdomain.sh start        - 启动 奶霸米表"
    echo "./nbdomain.sh stop         - 停止 奶霸米表"
    echo "./nbdomain.sh restart      - 重启 奶霸米表"
    echo "./nbdomain.sh status       - 查看 奶霸米表 状态"
    echo "./nbdomain.sh enable       - 设置 奶霸米表 开机自启"
    echo "./nbdomain.sh disable      - 取消 奶霸米表 开启自启"
    echo "./nbdomain.sh log          - 查看 奶霸米表 日志"
    echo "./nbdomain.sh update       - 更新 奶霸米表"
    echo "./nbdomain.sh install      - 安装 奶霸米表"
    echo "./nbdomain.sh uninstall    - 卸载 奶霸米表"
    echo "------------------------------------------"
}

show_menu() {
    echo -e "
  ${green}奶霸米表管理脚本${plain} ${red}${version}${plain}
--- https://nai.ba ---
  ${green}0.${plain} 退出脚本
————————————————
  ${green}1.${plain} 安装 奶霸米表
  ${green}2.${plain} 更新 奶霸米表
  ${green}3.${plain} 卸载 奶霸米表
————————————————
  ${green}4.${plain} 修改米表账号密码
  ${green}5.${plain} 修改米表监听端口
  ${green}6.${plain} 手动修改配置
————————————————
  ${green}7.${plain} 启动 奶霸米表
  ${green}8.${plain} 停止 奶霸米表
  ${green}9.${plain} 重启 奶霸米表
 ${green}10.${plain} 查看 奶霸米表 日志
————————————————
 ${green}11.${plain} 设置 奶霸米表 开机自启
 ${green}12.${plain} 取消 奶霸米表 开机自启
————————————————
 ${green}13.${plain} 一键安装 bbr (最新内核)
 ${green}14.${plain} 升级此脚本
 "
    show_status
    echo && read -p "请输入选择 [0-14]: " num

    case "${num}" in
    0)
        exit 0
        ;;
    1)
        check_uninstall && install
        ;;
    2)
        check_install && update
        ;;
    3)
        check_install && uninstall
        ;;
    4)
        check_install && modify_user
        ;;
    5)
        check_install && modify_port
        ;;
    6)
        check_install && modify_config
        ;;
    7)
        check_install && start
        ;;
    8)
        check_install && stop
        ;;
    9)
        check_install && restart
        ;;
    10)
        check_install && show_log
        ;;
    11)
        check_install && enable
        ;;
    12)
        check_install && disable
        ;;
    13)
        install_bbr
        ;;
    14)
        update_shell
        ;;
    *)
        echo -e "${red}请输入正确的数字 [0-14]${plain}"
        ;;
    esac
}

if [[ $# > 0 ]]; then
    case $1 in
    "start")
        check_install 0 && start 0
        ;;
    "stop")
        check_install 0 && stop 0
        ;;
    "restart")
        check_install 0 && restart 0
        ;;
    "status")
        check_install 0 && show_status 0
        ;;
    "enable")
        check_install 0 && enable 0
        ;;
    "disable")
        check_install 0 && disable 0
        ;;
    "log")
        check_install 0 && show_log 0
        ;;
    "update")
        check_install 0 && update 0
        ;;
    "install")
        check_uninstall 0 && install 0
        ;;
    "uninstall")
        check_install 0 && uninstall 0
        ;;
    *) show_usage ;;
    esac
else
    show_menu
fi
