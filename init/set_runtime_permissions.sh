#!/bin/sh
# Run on the Linux host: sudo sh /path/to/chatwiki/init/set_runtime_permissions.sh
set -eu

if ! command -v setfacl >/dev/null 2>&1; then
    echo '缺少 setfacl。请先安装 acl 软件包，再重新运行脚本。' >&2
    echo 'Debian/Ubuntu: sudo apt-get install acl' >&2
    echo 'RHEL/Rocky/AlmaLinux: sudo dnf install acl' >&2
    exit 1
fi

script=$(realpath -- "$0")
project_dir=$(dirname -- "$(dirname -- "$script")")
set -- "$project_dir/clawbot/working_dir" "$project_dir/internal/app/code_run/logs"

# Validate both paths before modifying either directory.
for dir do
    if [ "$(realpath -m -- "$dir")" != "$dir" ] || { [ -e "$dir" ] && [ ! -d "$dir" ]; }; then
        printf '目标路径必须是项目内的真实目录：%s\n' "$dir" >&2
        exit 1
    fi
done

for dir do
    mkdir -p -- "$dir"
    (
        set -x
        setfacl -R -m u:10001:rwx "$dir"
        setfacl -R -m d:u:10001:rwx "$dir"
    )
    printf '已设置 UID 10001 的访问和继承权限：%s\n' "$dir"
done
