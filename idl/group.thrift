namespace go kitexgroupservice

include"commonmodel/common_config.thrift"

service KitexGroupService{
    GetGroupInfoResp GetGroupInfo(1:GetGroupInfoReq req)
    ChangeGroupInfoResp ChangeGroupInfo(1:ChangeGroupInfoReq req)
    SearchGroupResp SearchGroup(1:SearchGroupReq req)
    CreateGroupResp CreateGroup(1:CreateGroupReq req)
    DeleteGroupResp DeleteGroup(1:DeleteGroupReq req)
    LeaveGroupResp LeaveGroup(1:LeaveGroupReq req)
    SetGroupApplyResp SetGroupApply(1:SetGroupApplyReq req)
    GetGroupApplyListResp GetGroupApplyList(1:GetGroupApplyListReq req)
    GetLastVisitTimeResp GetLastVisitTime(1:GetLastVisitTimeReq req)
    AgreeGroupApplyResp AgreeGroupApply(1:AgreeGroupApplyReq req)
    RefuseGroupApplyResp RefuseGroupApply(1:RefuseGroupApplyReq req)
    TransformGroupOwnerResp TransformGroupOwner(1:TransformGroupOwnerReq req)
    KickOutGroupResp KickOutGroup(1:KickOutGroupReq req)
    SetManagerResp SetManager(1:SetManagerReq req)
    RevokeManagerResp RevokeManager(1:RevokeManagerReq req)
    GetGroupInfoWithUserResp GetGroupInfoWithUser(1:GetGroupInfoWithUserReq req)
    UpdateGroupInfoWithUserResp UpdateGroupInfoWithUser(1:UpdateGroupInfoWithUserReq req)

    CreatSessionResp CreatSession(1:CreatSessionReq req)
    DeleteSessionResp DeleteSession(1:DeleteSessionReq req)
    GetFriendLastVisitTimeResp GetFriendLastVisitTime(1:GetFriendLastVisitTimeReq req)
    ApplyForFriendResp ApplyForFriend(1:ApplyForFriendReq req)
    GetFriendApplyListResp GetFriendApplyList(1:GetFriendApplyListReq req)
    RefuseFriendApplyResp RefuseFriendApply(1:RefuseFriendApplyReq req)

    GetGroupAndSessionIDResp GetGroupAndSessionID(1:GetGroupAndSessionIDReq req)
    GetGroupOrSessionRoleAndExistResp GetGroupOrSessionRoleAndExist(1:GetGroupOrSessionRoleAndExistReq req)
    SetLastVisitTimeResp SetLastVisitTime(1:SetLastVisitTimeReq req)
    GetGroupUserIDResp GetGroupUserID(1:GetGroupUserIDReq req)

    SetMuteResp SetMute(1:SetMuteReq req)
    ReleaseMuteResp ReleaseMute(1:ReleaseMuteReq req)
    GetMuteStatusResp GetMuteStatus(1:GetMuteStatusReq req)
}

//group service
struct GetGroupInfoReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
}
struct GetGroupInfoResp{
    1:i64 group_id
    2:string group_name
}

struct ChangeGroupInfoReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
    4:string group_name
}
struct ChangeGroupInfoResp{}

struct SearchGroupReq{
    1:common_config.CommonInfo common_info
    2:string group_name
}
struct SearchGroupResp{
    1:list<i64> group_id_list
}

struct CreateGroupReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:string group_name
}
struct CreateGroupResp{
    1:i64 group_id
}

struct DeleteGroupReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct DeleteGroupResp{}

struct LeaveGroupReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct LeaveGroupResp{}

struct SetGroupApplyReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct SetGroupApplyResp{}

struct GetGroupApplyListReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct GetGroupApplyListResp{
    1:list<i64> apply_user_id_list
}

struct GetLastVisitTimeReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct GetLastVisitTimeResp{
    1:list<i64> user_id_list
    2:list<i64> last_visit_time_list
}

struct AgreeGroupApplyReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct AgreeGroupApplyResp{}

struct RefuseGroupApplyReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct RefuseGroupApplyResp{}

struct TransformGroupOwnerReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct TransformGroupOwnerResp{}

struct KickOutGroupReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct KickOutGroupResp{}

struct SetManagerReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct SetManagerResp{}

struct RevokeManagerReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct RevokeManagerResp{}

struct GetGroupInfoWithUserReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct GetGroupInfoWithUserResp{
    1:i64 group_id
    2:string group_remark_name
    3:string role

}

struct UpdateGroupInfoWithUserReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:string group_remark_name
}
struct UpdateGroupInfoWithUserResp{}

struct GetGroupUserIDReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
}
struct GetGroupUserIDResp{
    1:list<i64> user_id_list
}

//session service
struct CreatSessionReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 goal_user_id
}
struct CreatSessionResp{
    1:i64 session_id
}

struct DeleteSessionReq{
    1:common_config.CommonInfo common_info
    2:i64 session_id
    3:i64 user_id
}
struct DeleteSessionResp{}

struct GetFriendLastVisitTimeReq{
    1:common_config.CommonInfo common_info
    2:i64 session_id
    3:i64 goal_user_id
}
struct GetFriendLastVisitTimeResp{
    1:string last_visit_time
}

struct ApplyForFriendReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 goal_user_id
}
struct ApplyForFriendResp{}

struct GetFriendApplyListReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
}
struct GetFriendApplyListResp{
    1:list<i64> apply_user_id_list
}

struct RefuseFriendApplyReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 goal_user_id
}
struct RefuseFriendApplyResp{}

//
struct GetGroupAndSessionIDReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
}
struct GetGroupAndSessionIDResp{
    1:list<i64> group_id_list
    2:list<i64> session_id_list
    3:list<i64> user_of_session_id_list
}

struct GetGroupOrSessionRoleAndExistReq{
    1:common_config.CommonInfo common_info
    2:i64 group_id
    3:i64 user_id
}
struct GetGroupOrSessionRoleAndExistResp{
    1:string role
    2:bool exist
}

struct SetLastVisitTimeReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct SetLastVisitTimeResp{}
//mute
struct SetMuteReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
    5:i64 mute_time_second
    6:string mute_reason
}
struct SetMuteResp{}

struct ReleaseMuteReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
    4:i64 goal_user_id
}
struct ReleaseMuteResp{}

struct GetMuteStatusReq{
    1:common_config.CommonInfo common_info
    2:i64 user_id
    3:i64 group_id
}
struct GetMuteStatusResp{
    1:string mute_reason
    2:string mute_end_time
    3:bool is_mute
}