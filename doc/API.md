# API
记录接口和对应的功能

# 返回格式
暂时对于正常的返回使用的 http code 都是 200，如果出现 503 等错误是服务器内部错误
```json
{
    "code": 200,
    "msg": "ok",
    "data": []
}
```
具体错误根据返回内容中的 code 判断，对应的 code 和错误对照表


## 登陆部分

- GET /login
登陆接口，访问时判断是否已经登录  
    - 如果已经登录，跳转到 /onedrive/getallfiles
    - 如果没有登陆，跳转到 /loginmg

- GET /loginmg
跳转到微软登陆验证页面

- GET /auth
用来接收 code，接收后会请求 AccessToken，初始化自动刷新

## 获取部分

- GET /onedrive/getallfiles
获取所有文件的树结构，返回 data 字段结构如下
```json
{
    "name": "macOS Catalina",
    "path": "/macOS Catalina",
    "is_folder": true,
    "last_modify_time": "2020-03-04T06:51:34Z",
    "children": [
        {
            "name": "Dark Mode JPEG_Original.jpeg",
            "path": "/macOS Catalina/Dark Mode JPEG_Original.jpeg",
            "is_folder": false,
            "last_modify_time": "2020-03-04T06:52:27Z",
            "children": null
        }
    ]
}
```

- GET /onedrive/getpath?path=xxx
通过相对路径获取文件列表，和 `/onedrive/getallfiles` 返回一样