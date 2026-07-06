namespace go kitexmessageservice

include"commonmodel/common_config.thrift"
include"commonmodel/common_message.thrift"

service KitexMessageService{
    SendMessageResp SendMessage(1:SendMessageReq req)
    SendFileResp SendFile(1:SendFileReq req)
    SendVoiceResp SendVoice(1:SendVoiceReq req)
    SendPictureResp SendPicture(1:SendPictureReq req)
    WithdrawMessageResp WithdrawMessage(1:WithdrawMessageReq req)
    DeleteMessageAllGroupResp DeleteMessageAllGroup(1:DeleteMessageAllGroupReq req)
    GetMessageListResp GetMessageList(1:GetMessageListReq req)
    GetNewMessageResp GetNewMessage(1:GetNewMessageReq req)
    GetFileContentResp GetFileContent(1:GetFileContentReq req)
    SendGroupNoticeResp SendGroupNotice(1:SendGroupNoticeReq req)
    SetOfflineMessageResp SetOfflineMessage(1:SetOfflineMessageReq req)
    GetOfflineMessageListResp GetOfflineMessageList(1:GetOfflineMessageListReq req)
}

struct SendMessageReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:string message_content
    5:bool is_ai
}
struct SendMessageResp{
    1:i64 message_id
}

struct SendFileReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:string file_name
    5:string content_type
    6:binary data_stream
}
struct SendFileResp{
    1:i64 message_id
}

struct SendVoiceReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:string content_type
    5:i64 voice_time_second
    6:binary data_stream
}
struct SendVoiceResp{
    1:i64 message_id
}

struct SendPictureReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:string content_type
    5:binary data_stream
}
struct SendPictureResp{
    1:i64 message_id
}

struct WithdrawMessageReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:i64 message_id
}
struct WithdrawMessageResp{}

struct DeleteMessageAllGroupReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
}
struct DeleteMessageAllGroupResp{}

struct GetMessageListReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:i64 start_time_second
    5:i64 end_time_second
}
struct  GetMessageListResp{
    1:list<common_message.KitexMessageInfo> message_info
}

struct GetNewMessageReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
}
struct GetNewMessageResp{
    1:list<i64> message_id
    2:list<i64> send_time_second
    3:list<string> message_type
    4:list<string> message_content
}

struct GetFileContentReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:i64 message_id
}
struct GetFileContentResp{
    1:binary data_stream
    2:string content_type
}

struct SendGroupNoticeReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:string message_content
}
struct SendGroupNoticeResp{}

struct SetOfflineMessageReq{
    1:common_config.CommonInfo common_info
    2:list<string> goal_user_and_device_id
    3:binary json_data
}
struct SetOfflineMessageResp{}

struct GetOfflineMessageListReq{
    1:common_config.CommonInfo common_info
    2:string user_and_device_id
}
struct GetOfflineMessageListResp{
    1:list<binary> json_data
    2:bool exist
}
