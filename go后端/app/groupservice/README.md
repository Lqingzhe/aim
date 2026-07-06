# GroupService
## 职责
- 群聊管理，包括：群聊的创建、解散，群聊申请的发送、审核，群基本信息的更改，群公告的发送
- 群成员管理：禁言，管理员的设立、罢免，群主的转移
- 私聊管理，包括：好友的申请的发起、同意和拒绝
## 数据结构
### Mysql
| 表名                   | 字段                                                 | 说明                      |
|----------------------|----------------------------------------------------|-------------------------|
| group_info           | group_id,group_name                                | 群聊的基本新信息，还可以添加申请条件、群头像等 |
| group_with_user_info | group_id, user_id (复合唯一), group_remark_name, role, | 用户在群中的信息，还可以添加头衔、称号等    |
| group_apply_info     | goal_id, apply_user_id                             | 群申请记录                   |
| group_mute_info      | group_id, user_id, mute_end_time, mute_reason,     | 禁言记录                    |
| session_info         | session_id, user_id, goal_user_id                  | 私聊会话                    |
### Redis
| Key格式                             | 数据结构                                  | 用途                               |
|-----------------------------------|---------------------------------------|----------------------------------|
| group_member:role:{groupID}       | ZSet (member=userID, score=role)      | 群成员角色缓存，快速获取用户的群权限以及是否加入会话（包括私聊） |
| group_member:visit_time:{groupID} | ZSet (member=userID, score=timestamp) | 群成员最后访问时间，用来判断消息的已读状态            |
| group_mute_info:{groupID}{userID} | Hash                                  | 禁言信息，Redis中失效就表示禁言结束             |
## 其他
### 一点点设计思路
在我的看法中，私聊可以看成特殊的群聊，再加上groupID和sessionID都是由雪花算法生成的，极小概率出现冲突的情况，所以对于 `判断用户是否在会话中以及用户是否有权限进行群操作` 这种高频情况，我将私聊和群聊统一放在Redis的 group_member:role:{groupID}中，既实现了redis缓存，也不用按照私聊和群聊分别维护两个redis的Key。