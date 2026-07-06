namespace go kitexcommonmodel

struct KitexMessageInfo{
    1:string group_id
    2:string message_id
    3:string user_id
    4:string message_content
    5:string content_type
    6:i64 voice_duration_second
    7:bool is_ai
    8:string message_type
    9:i64 send_time_second
}