package commonmodel

type MessageType string

const (
	KafkaMessageType_Notice  MessageType = "notice"
	KafkaMessageType_Message MessageType = "message"
	KafkaMessageType_System  MessageType = "system"
)

type MessageCode string

const (
	MessageCode_FriendRequest         MessageCode = "friend_request"
	MessageCode_FriendRequest_Success MessageCode = "friend_request_success"
	MessageCode_FriendRequest_Refuse  MessageCode = "friend_request_error"
	MessageCode_FriendDelete          MessageCode = "friend_delete"

	MessageCode_GroupApply          MessageCode = "group_apply"
	MessageCode_GroupDisband        MessageCode = "group_disband"
	MessageCode_GroupInfoChange     MessageCode = "group_info_change"
	MessageCode_GroupNotice         MessageCode = "group_notice"
	MessageCode_GroupJoin           MessageCode = "group_join"
	MessageCode_GroupRefuse         MessageCode = "group_refuse"
	MessageCode_GroupLeave          MessageCode = "group_leave"
	MessageCode_GroupKick           MessageCode = "group_kick"
	MessageCode_TransformGroupOwner MessageCode = "group_transform_owner"
	MessageCode_SetGroupManager     MessageCode = "set_group_manager"
	MessageCode_RevokeGroupManager  MessageCode = "revoke_group_manager"
	MessageCode_GroupSetMute        MessageCode = "group_set_mute"
	MessageCode_GroupReleaseMute    MessageCode = "group_release_mute"

	MessageCode_GroupMessage           MessageCode = "group_message"
	MessageCode_GroupWithdrawMessage   MessageCode = "group_withdraw_message"
	MessageCode_GroupEmpphasizeMessage MessageCode = "group_emphasize_message"
)

type GroupRole string

const (
	GroupOwner GroupRole = "Owner"
	Manager    GroupRole = "Manager"
	Member     GroupRole = "Member"
)

type FileType string

const (
	FileType_File    FileType = "file"
	FileType_Voice   FileType = "voice"
	FileType_Picture FileType = "picture"
)
