# Amprobe

![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
![Go Reference](https://pkg.go.dev/badge/github.com/shirou/gopsutil/v3.svg)
[![GitHub stars](https://img.shields.io/github/stars/amuluze/amprobe)](https://github.com/amuluze/amprobe/stargazers)

中文 | [English](./README.en.md)

## 简介

`Amprobe` 是一款轻量级主机及 `Docker` 容器监控工具，它可以轻松的帮助我们完成以下几方面的工作：

### 容器管理

- 查看 `Docker` 版本信息
- 容器的创建、启动、停止、重启、删除、查看容器日志
- 镜像的导入、导出、删除，虚悬镜像清理
- 网络的创建、删除、状态查看

### 主机管理

- 查看主机名称、启动时间、发行版本、内核版本、系统类型
- 查看主机 `CPU`、内存、磁盘 `IO`、网络 `IO` 监控

### 审计

- 查看用户登录、登出、操作记录

官网地址：[官网 | Amprobe (official.amprobe.amuluze.com)](https://official.amprobe.amuluze.com/)

> **docker 版本要求：>= 20.10.9**

## 技术栈

`Amprobe` 采用前后端分离的技术架构。

**前端技术栈:**

- `Vue3`
- `TypeScript`
- `Element+`
- `Vue-router`
- `Pinia`

**后端技术栈:**

- `Golang`
- `Fiber`
- `Sqlite`

## 下载与安装

### 方式一：直接下载二进制文件

从 [Releases](https://github.com/amuluze/amprobe/releases) 页面下载对应您系统的最新版本二进制文件。

### 方式二：从源码编译

1. 克隆仓库
```bash
git clone https://github.com/amuluze/amprobe.git
cd amprobe
```

2. 编译前端
```bash
cd web
make install  # 安装依赖
make build    # 构建前端
cd ..
```

3. 编译后端
```bash
make assets   # 打包前端资源到后端
make wire     # 生成依赖注入代码
make bin      # 编译二进制文件
```

## 使用方法

安装使用参见：[使用手册](https://official.amprobe.amuluze.com/document)

## 请作者喝杯咖啡

非常感谢大家使用 `Amprobe`, 目前该项目由个人用业余时间在维护，如果本项目有帮助到你的话，可以考虑请作者喝杯咖啡~

<img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/amprobe/202403171446310.jpg" alt="阿慕"  width="300" height="300" />

## 其他

作者阿慕，作为一名 35 岁的临“退”程序员，目前正在尝试去探索一条能够延长职业生涯的可行方案，欢迎添加作者微信或关注作者公众号进行交流呀！

<img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/amprobe/202403171449114.jpg" alt="阿慕微信"  width="300" height="350" /> <img src="https://cdn.jsdelivr.net/gh/amuluze/picgo@main/amprobe/202403171450306.png" alt="公众号"  width="300" />

开源不易，最后别忘了给本项目点个 **Star** ~

## License

MIT
