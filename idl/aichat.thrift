namespace go kitexaiservice

include"commonmodel/common_config.thrift"

service KitexAiService{
    SendMessageToAiResp SendMessageToAi(1:SendMessageToAiReq req)
    DeleteChatContextResp DeleteChatContext(1:DeleteChatContextReq req)
    GetAiConfigResp GetAiConfig(1:GetAiConfigReq req)
    UpdateAiConfigResp UpdateAiConfig(1:UpdateAiConfigReq req)
    DeleteAiConfigResp DeleteAiConfig(1:DeleteAiConfigReq req)
}

struct SendMessageToAiReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:string message
}
struct SendMessageToAiResp{}

struct DeleteChatContextReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
}
struct DeleteChatContextResp{}

struct GetAiConfigReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
}
struct GetAiConfigResp{
    1:string model_name
    2:string base_url
    3:string api_key
    4:string role
    5:string prompt
}

struct UpdateAiConfigReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:string model_name
    4:string base_url
    5:string api_key
    6:string role
    7:string prompt
}
struct UpdateAiConfigResp{}

struct DeleteAiConfigReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
}
struct DeleteAiConfigResp{}