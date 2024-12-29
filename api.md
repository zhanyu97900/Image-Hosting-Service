## sql
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

