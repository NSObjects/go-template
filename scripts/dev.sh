#!/bin/bash

# Echo Admin 快速开发脚本
# 使用方法: ./scripts/dev.sh [command] [options]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Echo Admin 快速开发脚本"
    echo ""
    echo "用法: $0 [command] [options]"
    echo ""
    echo "命令:"
    echo "  setup          - 初始化开发环境"
    echo "  new <module>   - 创建新的API模块"
    echo "  build          - 构建项目"
    echo "  run            - 运行项目"
    echo "  test           - 运行测试"
    echo "  lint           - 代码检查"
    echo "  clean          - 清理生成的文件"
    echo "  db-setup       - 设置数据库"
    echo "  db-gen         - 生成数据库代码"
    echo "  help           - 显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 setup"
    echo "  $0 new user"
    echo "  $0 build"
    echo "  $0 run"
}

# 检查依赖
check_dependencies() {
    print_info "检查依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go 1.19+"
        exit 1
    fi
    
    # 检查gentool
    if ! command -v gentool &> /dev/null; then
        print_warning "gentool 未安装，正在安装..."
        go install gorm.io/gen/tools/gentool@latest
    fi
    
    print_success "依赖检查完成"
}

# 初始化开发环境
setup() {
    print_info "初始化开发环境..."
    
    # 检查依赖
    check_dependencies
    
    # 下载依赖
    print_info "下载Go模块依赖..."
    go mod download
    
    # 生成错误码
    print_info "生成错误码..."
    make gen-code
    
    # 创建必要的目录
    print_info "创建必要的目录..."
    mkdir -p logs
    mkdir -p tmp
    
    print_success "开发环境初始化完成"
}

# 创建新模块
new_module() {
    if [ -z "$1" ]; then
        print_error "请指定模块名称"
        echo "用法: $0 new <module_name>"
        exit 1
    fi
    
    module_name="$1"
    print_info "创建新模块: $module_name"
    
    # 使用改进的modgen工具生成模块
    go run tools/modgen/main.go --name="$module_name" --force
    
    print_success "模块 $module_name 创建完成"
}

# 构建项目
build() {
    print_info "构建项目..."
    
    # 格式化代码
    make fmt
    
    # 代码检查
    make vet
    
    # 构建
    make build
    
    print_success "项目构建完成"
}

# 运行项目
run() {
    print_info "启动项目..."
    
    # 检查配置文件
    if [ ! -f "configs/config.toml" ]; then
        print_error "配置文件不存在: configs/config.toml"
        exit 1
    fi
    
    # 运行项目
    make run
}

# 运行测试
test() {
    print_info "运行测试..."
    
    # 设置测试环境
    export RUN_ENVIRONMENT=test
    
    # 运行测试
    make test
    
    print_success "测试完成"
}

# 代码检查
lint() {
    print_info "代码检查..."
    
    # 格式化
    make fmt
    
    # 静态检查
    make vet
    
    # Linter
    make lint
    
    print_success "代码检查完成"
}

# 清理
clean() {
    print_info "清理生成的文件..."
    
    make clean
    
    # 清理临时文件
    rm -rf tmp/*
    rm -rf logs/*
    
    print_success "清理完成"
}

# 数据库设置
db_setup() {
    print_info "设置数据库..."
    
    # 检查数据库连接
    print_info "检查数据库连接..."
    
    # 生成数据库代码
    print_info "生成数据库代码..."
    make db-gen
    
    print_success "数据库设置完成"
}

# 生成数据库代码
db_gen() {
    print_info "生成数据库代码..."
    
    # 生成模型和查询
    make db-gen
    
    # 生成Dynamic SQL
    make db-gen-dynamic
    
    print_success "数据库代码生成完成"
}

# 主函数
main() {
    case "${1:-help}" in
        setup)
            setup
            ;;
        new)
            new_module "$2"
            ;;
        build)
            build
            ;;
        run)
            run
            ;;
        test)
            test
            ;;
        lint)
            lint
            ;;
        clean)
            clean
            ;;
        db-setup)
            db_setup
            ;;
        db-gen)
            db_gen
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
