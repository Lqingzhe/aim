namespace go kitexcommonmodel

struct KitexMessageInfo{
    1:i64 group_id
    2:i64 message_id
    3:i64 user_id
    4:string message_content
    5:string content_type
    6:i64 voice_duration_second
    7:bool is_ai
    8:string message_type
    9:i64 send_time_second
}