# GONEList

<img align="right" width="240" src="https://gonelist-doc.cugxuan.cn/img/logo/logo.png">

[![Build Status](https://github.com/gonelist/gonelist/actions/workflows/multi-arch.yml/badge.svg)](https://github.com/gonelist/gonelist/actionst)
[![Latest Release](https://img.shields.io/github/release/gonelist/gonelist.svg)](../../releases)
[![All Releases Download](https://img.shields.io/github/downloads/cugxuan/gonelist/total.svg)](../../releases)

这是一款使用 `golang + vue` 编写的基于 onedrive 的**在线共享网盘**
效果展示：https://gonelist.cugxuan.cn  
后端项目地址：https://github.com/gonelist/gonelist  
前端项目地址：https://github.com/gonelist/gonelist-web  
详细文档地址：https://gonelist-doc.cugxuan.cn  
有问题请提 issue，也可以进入 QQ 群交流，群号：1083165608

# 功能特性

- 「跨平台」，兼容 Linux/Windows/MacOS 等多个平台
- 「极速部署」，下载对应程序，修改配置即可前后端一键部署
- 「直链下载」，文件直链下载，下载不消耗服务器流量
- 「自动刷新」，自动刷新缓存，可自定义时间
- 「自定义目录」，支持将 onedrive 的某个目录作为根目录
- 「在线播放」，支持在线播放音频和视频，在线浏览图片
- 「海量文件」，支持单目录下上千个文件，正常使用
- 「多平台」，支持个人版、教育账号、世纪互联等
- 「README」，支持页面添加 README
- 「加密目录」，支持给目录加密
- 「登陆缓存」，登陆 onedrive 之后会有缓存，下次直接启动无需登录
- ...

注：支持绝大部分教育账号，部分 **教育账号** 因为需要管理员同意无法使用

# 安装配置

# 实体服务安装教程

如果您的**整个微软账号和 onedrive 网盘**内，**没有隐私内容**，可以按照下面的流程快速配置体验效果，完整的下载安装流程请看 [安装文档](https://gonelist-doc.cugxuan.cn)

## 快速配置体验

下载 [Github Release](https://github.com/cugxuan/gonelist/releases)
或者 [gonelist-release](https://gonelist.cugxuan.cn/#/gonelist-release) 中对应的包，Linux 系统下载
gonelist_linux_amd64.tar.gz，直接运行即可启动，以 Linux 系统本地启动为例

```
// 下载对应的安装包，也可下载 gonelist-release 中的包，下面命令不一定是最新版本
$ wget https://github.com/cugxuan/gonelist/releases/download/v0.5.3/gonelist_linux_amd64.tar.gz
$ tar -zxf gonelist_linux_amd64.tar.gz && cd gonelist_linux_amd64
$ ./gonelist_linux_amd64
```

打开 http://localhost:8000 按照提示登录后即可。如果是在本地部署，登陆成功会跳转到首页，此时已经完成部署。  
如果是在服务器部署，登陆成功会跳转到 http://localhost:8000/auth?code=xxx ，将当前网址改成 http://yoursite:8000/auth?code=xxx
再回车等待文件加载后，会自动跳转你的网站 http://yoursite:8000 。如果登陆后一直没有反应，可能是因为文件夹数量过多导致，建议设置「子文件夹」选项  
默认情况下读取当前路径的 `config.yml` 作为配置文件，或加参数 `--conf=dir1/file.yml` 指定配置文件路径

## 实体systemd服务安装

视频教程(包含了Azure应用程序的配置) https://www.bilibili.com/video/BV1PA411t7Jw/

## docker运行

视频教程 https://www.bilibili.com/video/BV1Vz4y1R7EK/

直接使用项目的`docker-compose.yml`去`docker-compose up -d`即可，建议把配置文件放在一个文件夹里，把文件夹挂载进去，否则直挂文件 docker 挂载的是 inode。
如果是群晖的 docker 上运行的话会不支持 docker 的 command 似乎，可以把配置文件的目录挂载到容器里，例如`/etc/config`，创建容器的时候加上环境变量`CONF_PATH=/etc/config/config.yml`。
token_path 写`/etc/config/`，然后创建容器的最后地方的`Entrypoint`和`命令`空着

```
.
├── config
│   └── config.yml
└── docker-compose.yml
```

## config.yml

如果需要修改配置，在配置文件中填对应内容即可

```
# gonelist 配置文件，注意配置字段和信息中间有个空格

# name 表示你的站点的名字，会显示在每个页面的左上角
name: GONEList

# Remote name，可选 onedrive, chinacloud
remote: onedrive

# onedrive 的获取层级，默认获取两层
level: 2

# 提供 onedrive 的应用配置，建议自己创建应用
client_id: 16e320f7-e427-4612-88da-f3d03e944d40
client_secret: lURpL3U@bBlmJ0:_dnU.LeLOGNGdVT30
# 提供 chinacloud 的应用配置，建议自己创建应用
#client_id: 2b54b127-b403-42a3-8b55-d25f3119aa13
#client_secret: a0CGqBT3f_8U5gztxKjxR-LNW-ZnTe.m

# 不建议修改，需要和应用中心设置的 redirect_url 一致
redirect_url: http://localhost:8000/auth

# 随意设置一个你喜欢的字符串，在 onedrive 认证时会使用
state: 23333

# 默认在当前文件夹下，不建议修改
# token 实现了下次启动时不需要重新登陆验证的功能
token_path:

# 可以配置 CDN 加速重定向，url 前缀
download_redirect_prefix:

# gonelist 服务设置，不建议修改
server:
  ReadTimeout: 0
  WriteTimeout: 0
  bind_global: true
  dist_path: ./dist/
  # 子文件夹设置，比如你只想挂载你盘根目录下的 public 文件夹，就使用 /public
  folder_sub: /
  gzip: true # 是否开启 gzip 加速, 默认开启
  port: 8000
  # 自动刷新时间
  refresh_time: 10
  site_url: http://localhost:8000 # 不建议修改，在启动后会自适应调整
  enable_upload: false # 是否允许文件上传

# 这是一个数组，可以对不同文件夹设置密码
# 但还是推荐在需要设置密码的文件夹下创建 .password 文件
# 在该文件中填写密码即可，文件格式最好是 UTF-8
pass_list:
  - pass:
    path:
#  - pass:
#    path:
```

# Contributors

- 开发：<a href="https://github.com/cugxuan"><img src="https://avatars1.githubusercontent.com/u/23120372?s=400&v=4" width="30"></a>
<a href="https://github.com/Sillywa/"><img src="https://avatars0.githubusercontent.com/u/22909601?s=400&v=4" width="30"></a>
<a href="https://github.com/zhangguanzhang"><img src="https://avatars3.githubusercontent.com/u/18641678?s=400&v=4" width="30"></a>
<a href="https://github.com/StringKe"><img src="https://avatars.githubusercontent.com/u/31089228?s=400&v=4" width="30"></a>
- logo
  设计：<a href="http://lambertchan.me/"><img src="https://avatars0.githubusercontent.com/u/39192150?s=400&v=4" width="30"></a>