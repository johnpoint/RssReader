# web

## 依赖安装 Project setup

```bash
yarn install
```

### 编译和热重载的开发模式 Compiles and hot-reloads for development
```
yarn serve
```

### 编译生成生产文件 Compiles and minifies for production
```
yarn build
```

### Lints and fixes files
```
yarn lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

### 本地缓存数据结构 Local cache data structure

### basic

```
jwt: String
login: Bool (true/false)
i18n: String (zh/en) 
```

#### feeds
```json
[
    {
        "id":40,
        "title":"johnpoint's blog",
        "link":"https://blog.lvcshu.com/atom.xml",
        "unread":0
    }
]
```

#### posts
```json
[
    {
        "ID":4624,
        "Title":"派早报：华为将发布 Mate 40、一加或推出中低端机型、特斯拉 Model 3 大更新等",
        "Source":"少数派",
        "date":"2020-10-12 08:59:26",
        "link":"https://sspai.com/post/63092",
        "read":true,
        "readAfter":false
    }
]
```
#### post[postID]
```json
{
    "ID":4624,
    "FID":57,
    "Title":"派早报：华为将发布 Mate 40、一加或推出中低端机型、特斯拉 Model 3 大更新等",
    "Content":"",
    "Url":"https://sspai.com/post/63092",
    "Description":"...",
    "Published":"1602464366",
    "Source":"少数派",
    "read":true,
    "readafter":false
}
```