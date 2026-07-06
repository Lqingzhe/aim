# UserService
## 职责
- 用户注册和登录的验证
- 用户基本信息管理
- 用户之间的备注功能
## 数据结构
### Mysql
| 表名              | 字段                                                                            | 说明     |
|-----------------|-------------------------------------------------------------------------------|--------|
| user_info       | user_id, user_name, introduction, birthday_year, birthday_month, birthday_day | 用户基本信息 |
| user_login_info | user_id, password, salt                                                       | 用户登录信息 |
| remark_info     | user_id, goal_user_id, nick_name                                              | 用户备注信息 |
