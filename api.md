## sql_images
### 表 images
    1. fileid 唯一
    2. id 
    3. sha256 文件的shahash
    4. name 文件的名称
    5. path文件的存储地址
    6. removed 是否删除 0未 1已
### 表 prepare 上传
    1.fileid
    2.name
    3.upload 
    4.timestamp
### 方法 插入
    先检查removed 有就更新 没有就直接插入
    插入 fileid sha256 name path removed
    先保存文件 在插入
### 方法 删除
    标记 文件的removed为0
### 方法 查询
    筛选 removed 不为0
### 方法 查询通过sha256
    筛选 sha256 相同的
## sql_token
### user
    id 主键
    username 唯一
    password 不唯一 hash一下
    userid
### token
    token 
    refresh 
    userid 外键
### 生成userid
    返回 userid
### 检查 userid 和username 
    返回 bool 和错误
### 查询token
    传入 userid
    返回TOKEN 和error
### 更新token
    传入 TOKEN
    返回一个 TOKEN，error
### 查询user
    传入username
    返回一个 USER 和 error
### 插入 user token
    传入 USER TOKEN
    返回 bool，error
### 加密password
    加密password
    返回 字符串,error
### 验证password
    传入 服务器所存密码，和password
    返回bool
### 生成token
    
## 文件
### 上传
### fileid
### 删除


## 前端 (数据在data中)
### 准备
    /api/prepare
    post
    1.sha256
    2.filename

    返回值：code：0 1 2 
    特 2 有newFileid

### 上传
    /api/upload
    post
    1.file
    2.fileid
    
    返回 code 0 ，1

### 删除
    /api/remove
    post
    1.fileid
     
    返回 code 0,1

### 响应文件
    /file/*filepath
    get
    fileid 

    返回 html 文件

### 响应全部文件
    /api/allfile
    get
    返回 0 1

## 登录
### 响应login.html
### 响应Register.html
### api
    1.登录
        /api/login
        post username password
        返回 0 1
    2.注册
        /api/register
        post username password
        返回 0 1
        