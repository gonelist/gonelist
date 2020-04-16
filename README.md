# GONEList

这是一款使用 `golang + vue` 编写的基于 onedrive 的**在线共享网盘**  
效果展示：https://gonelist.cugxuan.cn  
前端项目地址：https://github.com/Sillywa/gonelist-web
有问题请提 issue，也可以进入 QQ 群交流，群号：1083165608

# 功能特性

- 「跨平台」，兼容 Linux/Windows/MacOS 等多个平台
- 「极速部署」，下载对应程序，修改配置即可前后端一键部署
- 「直链下载」，文件直链下载，下载不消耗服务器流量
- 「自动刷新」，自动刷新缓存，可自定义时间
- 「自定义目录」，支持将 onedrive 的某个目录作为根目录
- ...

# 安装配置

如果您的 onedrive 网盘内，**没有隐私内容**，可以按照下面的流程快速配置体验效果，完整的下载安装流程请看 [安装文档](https://github.com/cugxuan/gonelist/wiki/Install)  


下载 [Github Release](https://github.com/cugxuan/GOIndex/releases) 或者 [gonelist-release](https://gonelist.cugxuan.cn/#/gonelist-release) 中对应的包，直接运行即可启动，以 Linux 系统本地启动为例
```
// 下载对应的安装包，也可以下载 github 的 release 链接
$ wget https://gonelist.cugxuan.cn/d/gonelist-release/gonelist_linux_amd64.tar.gz
$ tar -zxf gonelist_linux_amd64.tar.gz && cd gonelist_linux_amd64
$ ./gonelist_linux_amd64
```
打开 http://localhost:8000 按照提示登录后即可

默认情况下读取当前路径的 `config.json` 作为配置文件，或加参数 `--conf=dir1/file.json` 指定配置文件路径

## config.json

如果需要修改配置，在配置文件中填对对应的内容即可
```
{
  //------建议填入自己的 id 和 secret --------
  "client_id": "88966400-cb81-49cb-89c2-6d09f0a3d9e2",
  "redirect_url": "http://localhost:8000/auth",
  "client_secret": "/FKad]FPtKNk-=j11aPwEOBSxYUYUU54",
  // 设置一个自己喜欢的字符串
  "state": "23333",
  "server": {
    // 监听的端口
    "port": 8000,
    // 自动刷新的时间单位是分钟，默认 10 分钟，不要超过 1 小时
    "refresh_time": 10,
    // 登陆成功后，跳转的 URL，可不设置，新版已自动跳转
    "site_url": "http://localhost:8000",
    // 自定义 onedrive 的子文件夹
    "folder_sub": "/",
    //静态页面的目录，默认当前路径下的dist目录
    "dist_path": "./dist/",
    // 是否绑定到 0.0.0.0
    "bind_global": true
  }
}
```

# 参考项目

前端页面 UI 参考：
https://moeclub.org/onedrive/