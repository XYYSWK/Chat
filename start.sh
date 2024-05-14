#!/bin/sh
# 声明脚本使用的解释器是 /bin/sh，即 Bourne shell

set -e # 设置 shell 在执行脚本过程中遇到任何命令非零退出状态时直接退出，确保脚本在遇到错误时立刻停止执行

# 输出信息，提示即将执行数据库迁移操作
echo "run db migrate"

# 输出信息，指示正在执行数据库迁移的数据库名称
echo "postgres_zr"
#  运行数据库迁移
./migrate -database "postgres://postgres:123456@121.36.86.143:5432/chatroom?sslmode=disable" -path "app/migration" up

# 输出信息，指示即将启动应用程序
echo "start the app"

# 执行传递给脚本的所有参数
echo "$@"