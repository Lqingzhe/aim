# FileService
## 职责
- 用于储存文件、图片和音频，同时生成并返回file_id
## 数据结构
### Mysql
| 表明          | 字段                                                                                                                             | 说明    |
|-------------|--------------------------------------------------------------------------------------------------------------------------------|-------|
| file_models | file_id(message_id),file_name,file_type, content_type, voice_duration_second, storage_path, created_at, updated_at, deleted_at | 文件元数据 |
### 本地文件系统
| 路径                                                                        | 说明        |
|---------------------------------------------------------------------------|-----------|
| {file_storage_path}/data/aim-files/{year}/{month}/{day}/{file_id}{suffix} | 文件实际的储存路径 |
|                                                                           |           |