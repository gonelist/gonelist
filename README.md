# GONEList

这是一款使用 `golang + vue` 编写的基于 onedrive 的在线网盘  
效果展示：https://gonelist.cugxuan.cn  
前端项目地址：https://github.com/Sillywa/gonelist-web

# 功能特性

- 「跨平台」，兼容 Linux/Windows/MacOS 等多个平台
- 「极速部署」，下载对应程序，修改配置即可前后端一键部署
- 「直链下载」，文件直链下载，下载不消耗服务器流量
- 「自动刷新」，自动刷新缓存，可自定义时间
- 「自定义目录」，支持将 onedrive 的某个目录作为根目录
- ...

# 安装配置

完整的下载安装流程请看 [安装文档](https://github.com/cugxuan/gonelist/wiki/Install)
下载对应的 [Github Release](https://github.com/cugxuan/GOIndex/releases)，直接运行即可启动
默认情况下读取当前路径的 `config.json` 作为配置文件，或者启动加参数 `--conf=dir1/file.json` 指定配置文件路径

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
    // 登陆成功后，跳转的 URL，可不设置
    "site_url": "http://localhost:8000",
    // 自定义 onedrive 的子文件夹
    "folder_sub": "/"
    // 是否绑定到 0.0.0.0
    "bind_global": true,
  }
}
```

# 参考项目

前端页面 UI 参考：
https://moeclub.org/onedrive/