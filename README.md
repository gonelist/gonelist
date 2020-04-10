# GOIndex

这是一款使用 `golang + vue` 编写的基于 onedrive 的在线网盘  
效果展示 https://goindex.cugxuan.cn

# 功能特性

- 文件直链下载
- 自动刷新缓存

# 安装配置

默认情况下读取当前路径的 `conf/config.json` 作为配置文件，或者启动加参数`--conf=dir1/file.json`指定配置文件路径

在配置文件中填对对应的内容即可，至少需要修改 `sub_path, site_url`
```
{
  //---不建议改动--------
  "client_id": "88966400-cb81-49cb-89c2-6d09f0a3d9e2",
  "redirect_url": "http://localhost:8000/auth",
  "client_secret": "/FKad]FPtKNk-=j11aPwEOBSxYUYUU54",
  // 设置一个自己喜欢的字符串
  "state": "23333",
  "server": {
    "run_mode": "run_mode",
    // 可以修改监听的端口
    "http_port": 8000,
    // 自动刷新的时间单位是分钟，默认 10 分钟
    "refresh_time": 10,
    "bind_global": false,
    // 设置路径前缀，如 https://yoursite.com/goindex/
    "sub_path": "goindex",
    "site_url": "https://goindex.cugxuan.cn",
  }
}
```