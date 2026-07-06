namespace go kitexfileservice

include"commonmodel/common_config.thrift"

service KitexFileService{
    CreateFileResp CreateFile(1:CreateFileReq req)
    DeleteFileResp DeleteFile(1:DeleteFileReq req)
    GetFileResp GetFile(1:GetFileReq req)
}

struct CreateFileReq{
    1:common_config.CommonInfo common_info
    2:string file_name
    3:binary data_stream
    4:string file_type
    5:string content_type
    6:i64 voice_duration_time_second

}
struct CreateFileResp{
    1:i64 file_id
}

struct DeleteFileReq{
    1:common_config.CommonInfo common_info
    2:i64 file_id

}
struct DeleteFileResp{}

struct GetFileReq{
    1:common_config.CommonInfo common_info
    2:i64 file_id
}
struct GetFileResp{
    1:binary data_stream
    2:string content_type
}