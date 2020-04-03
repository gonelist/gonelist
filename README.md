# GOIndex

默认情况下读取当前路径的 `conf/config.json` 作为配置文件，或者启动加参数`--conf=dir1/file.json`指定配置文件路径

在配置文件中填对对应的内容即可
```
{
  //---不建议改动--------
  "client_id": "88966400-cb81-49cb-89c2-6d09f0a3d9e2",
  "redirect_url": "http://localhost:8000/auth",
  "client_secret": "/FKad]FPtKNk-=j11aPwEOBSxYUYUU54",
  // 设置一个自己喜欢的字符串
  "state": "23333",
  // 设置路径前缀，如 http://yoursite.com/goindex/
  "sub_path": "goindex",
  "server": {
    "run_mode": "run_mode",
    "http_port": 8000,
    "refresh_time": 10
  }
}
```
- client_id 客户端 id 
- redirect_url 填 http://localhost:8000/auth // 不建议改动
- client_secret 客户端密码 // 不建议改动
- state 填一个随机字符串  // 可随意改动
