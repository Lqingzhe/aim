# MessageService
## 职责
- 负责消息的发送，拉取
- 负责离线消息的储存
## 数据结构
### MongoDB
| 集合      | 字段                                                                                                                                                | 说明     |
|---------|---------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| message | _id (message_id), group_id, user_id, message_content, message_type, file_storage_id, content_type, voice_duration_second, is_ai, send_time_second | 聊天消息存储 |
### Mysql
| 表名                   | 字段                                      | 说明        |
|----------------------|-----------------------------------------|-----------|
| offline_message_info | message_id, json_data, send_time_second | 离线消息的详细数据 |
### Redis
| Key格式                             | 数据结构 | 用途             |
|-----------------------------------|------|----------------|
| offset_message:{userID}{deviceID} | Set  | 离线消息ID集合（7天过期） |
## 其他
### “@"功能
可以`@UserID1\t@UserID2\t`的形式call多个用户  
在没有`@UserID `的情况下可以`@bot `调用Ai助手