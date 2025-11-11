# Docker 镜像拉取工具（Linux）

一个基于渡渡鸟镜像源（https://docker.aityp.com）的轻量级 Docker 镜像拉取工具，提供更便捷的镜像搜索和拉取体验。

## 1.1 功能特性

- 🔍 **智能搜索**: 从渡渡鸟镜像源搜索Docker镜像
- 🎯 **平台支持**: 支持指定平台架构（amd64、arm64等）
- 🏷️ **标签过滤**: 按标签筛选镜像版本
- 💬 **交互式选择**: 当有多个镜像可选时，提供交互式选择界面
- 📊 **详细信息**: 显示镜像大小、创建时间、平台等详细信息
- 🚀 **一键拉取**: 自动调用docker pull命令拉取镜像

![{BB2913E3-3F22-4B89-A827-A1A8045603B8}](./md_images/%7BBB2913E3-3F22-4B89-A827-A1A8045603B8%7D.png)

![{8316DE4C-30F6-4997-A035-B5D4D74826D0}](./md_images/%7B8316DE4C-30F6-4997-A035-B5D4D74826D0%7D.png)

## 1.2 安装

### 1.2.1 直接安装

> 前提条件：
>
> - 已安装 Docker
> - 只支持 Linux 环境

```shell
# 克隆仓库
git clone https://github.com/Stars1eap/Docker-Images-Pull-From-Dodo-Bird.git

# 进入目录
cd Docker-Images-Pull-From-Dodo-Bird

# 使用
./dimages pull 镜像名

# 将可执行文件移动到系统PATH（可选）
mv dimages /usr/local/bin/
```

### 1.2.2 编译安装

> 前提条件：
>
> - 已安装 Go 1.16 或更高版本
> - 已安装 Docker

```shell
# 克隆仓库
git clone https://github.com/Stars1eap/Docker-Images-Pull-From-Dodo-Bird.git

# 进入目录
cd Docker-Images-Pull-From-Dodo-Bird

# 编译项目
go build -o dimages

# 将可执行文件移动到系统PATH（可选）
mv dimages /usr/local/bin/
```

## 1.3 使用方法

### 1.3.1 基本使用

```shell
# 拉取镜像（默认选择第一个匹配的镜像）
./dimages pull nginx

# 显示帮助信息
./dimages help
```

### 1.3.2 高级选项

```shell
# 指定平台架构
./dimages pull -platform linux/arm64 nginx

# 指定标签版本
./dimages pull -tag latest nginx

# 交互式选择镜像
./dimages pull -i nginx

# 组合使用
./dimages pull -platform linux/amd64 -tag alpine -i nginx
```

### 1.3.3 交互式选择示例

当使用 `-i` 参数或找到多个匹配镜像时，会显示交互式选择界面：

```shell
找到 3 个镜像:
[1] nginx:latest
    镜像: registry.docker.aityp.com/nginx:latest
    平台: linux/amd64, 大小: 142MB, 创建时间: 2024-01-15

[2] nginx:alpine
    镜像: registry.docker.aityp.com/nginx:alpine
    平台: linux/amd64, 大小: 41MB, 创建时间: 2024-01-14

[3] nginx:1.23
    镜像: registry.docker.aityp.com/nginx:1.23
    平台: linux/arm64, 大小: 138MB, 创建时间: 2024-01-13

请选择要拉取的镜像 (输入编号): 2
```

## 1.4 命令行选项

| 选项        | 说明           | 示例                         |
| :---------- | :------------- | :--------------------------- |
| `-platform` | 指定平台架构   | `linux/amd64`, `linux/arm64` |
| `-tag`      | 指定镜像标签   | `latest`, `alpine`, `1.23`   |
| `-i`        | 交互式选择模式 | 显示可选镜像列表             |
| `help`      | 显示帮助信息   | `./dimages help`             |

## 1.5 项目结构

```
Docker-Images-Pull-From-Dodo-Bird/
├── main.go          # 主程序文件
├── README.md        # 项目说明文档
└── dimages         # Linux 环境编译的二进制程序
```

## 1.6 技术栈

- **语言**: Go
- **镜像源**: 渡渡鸟Docker镜像仓库（https://docker.aityp.com/）
- **依赖**: 标准库（net/http, encoding/json, flag等）

### 1.6.1 API接口

工具使用渡渡鸟镜像源的公开API：

```shell
GET https://docker.aityp.com/api/v1/image?search=<镜像名>
```

## 贡献

欢迎提交Issue和Pull Request来改进这个项目！

## 致谢

- 感谢 [渡渡鸟镜像源](https://docker.aityp.com/) 提供的镜像服务

------

**作者**: Starsleap
**仓库**: https://github.com/Stars1eap/Docker-Images-Pull-From-Dodo-Bird