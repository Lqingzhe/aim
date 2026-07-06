namespace go kitexuserservice

include"commonmodel/common_user.thrift"
include"commonmodel/common_config.thrift"

service KitexUserService{
    GetUserInfoResp GetUserInfo(1:GetUserInfoReq req)
    GetOtherUserInfoResp GetOtherUserInfo(1:GetOtherUserInfoReq req)
    UpdateUserInfoResp UpdateUserInfo(1:UpdateUserInfoReq req)
    RemarkResp Remark(1:RemarkReq req)
    RegisterResp Register(1:RegisterReq req)
    LoginResp Login(1:LoginReq req)

}

struct GetUserInfoReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id

}
struct GetUserInfoResp{
    1:common_user.UserInfo user_info
    2:list<common_user.RemarkInfo> remark_info
}
struct GetOtherUserInfoReq{
    1:common_config.CommonInfo common_info
    2:i64 goal_user_id
}
struct GetOtherUserInfoResp{
    1:common_user.UserInfo user_info
}

struct UpdateUserInfoReq{
    1:common_config.CommonInfo common_info
    2:common_user.UserInfo user_info
}
struct UpdateUserInfoResp{}

struct RemarkReq{
    1:common_config.CommonInfo common_info
    2:common_user.RemarkInfo remark_info
}
struct RemarkResp{}

struct RegisterReq{
    1:common_config.CommonInfo common_info
    2:string password
}
struct RegisterResp{
    1:i64 user_id
}

struct LoginReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:string password
}
struct LoginResp{}

