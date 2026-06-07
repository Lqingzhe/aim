# GateWay
## 职责
- 负责连接外网和内部微服务实例，将HTTP协议转换为RPC协议
- 负责链路追踪id的生成
- 全局的ctx超时的设置
- 用户单设备级别的redis分布式限流
- AccessToken和RefreshToken的生成和管理
- 消费Kafka中的NewMessage,SystemMessage和GroupNotice,通过websocket推送至用户
## 数据结构
### Redis
| Key格式                        | 数据结构 | 用途              | 示例                      |
|------------------------------|------|-----------------|-------------------------|
| limiter:{userID}{deviceID}   | Hash | API 限流令牌桶       | limiter:123web-001      |
| refresh_token:{refreshToken} | Hash | RefreshToken 存储 | refresh_token:abc123... |
## 其他
### token
#### AccessToken
AccessToken采用jwt无状态会，包含了用户的UserID, DeviceID, 解析的时候会检验携带的设备信息和当前设备发送的是否相同，减小会话劫持的风险。一切需要登录的路由的UserID均由AccessToken携带
#### RefreshToken
RefreshToken采用cookie的模式，在redis中储存了用户的相关信息，在生成的时候采用`salt+userID(十六进制)+deviceID(用户的)`，用户可以持有有效的RefreshToken获取有效的AccessToken。  
考虑到RefreshToken有效时间较长，一旦被劫持，就会长时间被非法持有，所以我才用刷新机制： *每一次刷新AccessToken都会从redis中删除就RefreshToken并生成一个新的RefreshToken进行储存*  
这样不仅很大程度上减小了会话劫持的影响，同时对于经常使用的用户来说可以不用每次RefreshToken失效就登录，因为每一次刷新都会重置redis的储存时间。    
这套方案的一个小问题是无法做到立即登出，但是考虑到AIM系统对于登出延迟的影响较小，添加黑名单进制的付出小于收益，便采用的删除redis中的RefreshToken进行登出，AccessToken自然过期
### limiter
结合redis实现的令牌桶限流，不过由于是项目初期写的，没有使用Lua脚本，在数据一致性上没有做到原子性，所以后期会重构为Lua脚本的原子性实现。