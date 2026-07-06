package com.example.aim.network.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class ApiResponse(
    val code: Int = -1,
    val message: String = "",
    val data: ResponseData? = null
)

@Serializable
data class ResponseData(
    val token_info: TokenInfo? = null,
    val user_info: UserInfoData? = null,
    val session_info: SessionInfoData? = null,
    val group_info: GroupInfoData? = null,
    val message_info: MessageInfoData? = null,
    val message_list: List<MessageData>? = null,
    val ai_config: AiConfig? = null
)

@Serializable
data class TokenInfo(
    val access_token: String = "",
    val refresh_token: String = ""
)

@Serializable
data class UserInfoData(
    val UserInfo: UserInfo? = null,
    val RemarkInfos: List<RemarkInfo>? = null,
    val user_id: String? = null,
    val user_name: String = "",
    val introduction: String = "",
    val birthday_year: Int = 0,
    val birthday_month: Int = 0,
    val birthday_day: Int = 0,
    val is_connect: Boolean = false
) {
    fun resolveUserInfo(): UserInfo {
        val ui = UserInfo ?: return UserInfo(
            user_id = user_id ?: "",
            user_name = user_name,
            introduction = introduction,
            birthday_year = birthday_year,
            birthday_month = birthday_month,
            birthday_day = birthday_day,
            is_connect = is_connect
        )
        return ui
    }
}

@Serializable
data class UserInfo(
    @SerialName("user_id") val user_id: String = "",
    @SerialName("Userid") val Userid: String = "",
    @SerialName("user_name") val user_name: String = "",
    @SerialName("UserName") val UserName: String = "",
    @SerialName("introduction") val introduction: String = "",
    @SerialName("Introduction") val `Introduction_`: String = "",
    @SerialName("birthday_year") val birthday_year: Int = 0,
    @SerialName("BirthdayYear") val BirthdayYear: Int = 0,
    @SerialName("birthday_month") val birthday_month: Int = 0,
    @SerialName("BirthdayMonth") val BirthdayMonth: Int = 0,
    @SerialName("birthday_day") val birthday_day: Int = 0,
    @SerialName("BirthdayDay") val BirthdayDay: Int = 0,
    @SerialName("is_connect") val is_connect: Boolean = false,
    @SerialName("IsConnect") val IsConnect: Boolean = false
) {
    fun resolvedUserId(): String = user_id.ifEmpty { Userid }
    fun resolvedUserName(): String = user_name.ifEmpty { UserName }
    fun resolvedIntroduction(): String = introduction.ifEmpty { `Introduction_` }
    fun resolvedBirthdayYear(): Int = if (birthday_year > 0) birthday_year else BirthdayYear
    fun resolvedBirthdayMonth(): Int = if (birthday_month > 0) birthday_month else BirthdayMonth
    fun resolvedBirthdayDay(): Int = if (birthday_day > 0) birthday_day else BirthdayDay
    fun resolvedIsConnect(): Boolean = is_connect || IsConnect
}

@Serializable
data class RemarkInfo(
    val goal_user_id: String = "",
    val nick_name: String = "",
    val GoalUserID: String = "",
    val NickName: String = ""
) {
    fun resolvedGoalUserId(): String = goal_user_id.ifEmpty { GoalUserID }
    fun resolvedNickName(): String = nick_name.ifEmpty { NickName }
}

@Serializable
data class SessionInfoData(
    val session_id_list: List<String> = emptyList(),
    val user_of_session_id_list: List<String> = emptyList(),
    val session_id: String? = null,
    val apply_user_list: List<String> = emptyList(),
    val goal_user_id: String? = null,
    val last_visit_time: Long? = null
)

@Serializable
data class GroupInfoData(
    val group_id_list: List<String> = emptyList(),
    val group_id: String? = null,
    val group_name: String? = null,
    val group_remark_name: String? = null,
    val group_role: String? = null,
    val last_visit_time: Map<String, Long>? = null
)

@Serializable
data class MessageInfoData(
    val message_id: String = "",
    val message_list: List<MessageData> = emptyList()
)

@Serializable
data class MessageData(
    val message_id: String = "",
    val group_id: String = "",
    val sender_id: String = "",
    val message_content: String = "",
    val message_type: String = "",
    val send_time: Long = 0,
    val send_time_second: Long = 0,
    val is_withdrawn: Boolean = false,
    val user_id: String = "",
    val user_name: String = "",
    val session_id: String = ""
) {
    fun id(): String = message_id
    fun sender(): String = sender_id.ifEmpty { user_id }
    fun content(): String = message_content
    fun time(): Long = when {
        send_time_second > 0 -> send_time_second
        send_time > 0 -> send_time
        else -> 0
    }
    fun gid(): String = group_id.ifEmpty { session_id }
    fun type(): String = message_type
    fun name(): String = user_name
    fun withdrawn(): Boolean = is_withdrawn
}

@Serializable
data class AiConfig(
    val model_name: String = "",
    val base_url: String = "",
    val api_key: String = "",
    val role: String = "",
    val prompt: String = ""
)

// Auth
@Serializable
data class RegisterRequest(val password: String)

@Serializable
data class LoginRequest(val user_id: String, val password: String)

@Serializable
data class RefreshTokenRequest(val refresh_token: String)

// User
@Serializable
data class UpdateUserInfoRequest(
    val user_name: String = "",
    val introduction: String = "",
    val birthday_year: Int = 0,
    val birthday_month: Int = 0,
    val birthday_day: Int = 0
)

@Serializable
data class RemarkRequest(val goal_user_id: String, val nick_name: String)

// Group
@Serializable
data class CreateGroupRequest(val group_name: String)

@Serializable
data class GroupIdRequest(val group_id: String)

@Serializable
data class ChangeGroupInfoRequest(val group_id: String, val group_name: String)

@Serializable
data class SearchGroupRequest(val group_name: String)

@Serializable
data class GroupUserRequest(val group_id: String, val goal_user_id: String)

@Serializable
data class SetGroupApplyRequest(val group_id: String)

@Serializable
data class AgreeGroupApplyRequest(val group_id: String, val goal_user_id: String)

@Serializable
data class MuteRequest(
    val group_id: String,
    val goal_user_id: String,
    val mute_time_seconds: Int = 0,
    val mute_reason: String = ""
)

@Serializable
data class UpdateGroupRemarkRequest(val group_id: String, val group_remark_name: String)

// Friend
@Serializable
data class ApplyFriendRequest(val goal_user_id: String)

@Serializable
data class CreateSessionRequest(val goal_user_id: String)

@Serializable
data class DeleteSessionRequest(val session_id: String)

@Serializable
data class FriendLastVisitRequest(val session_id: String, val goal_user_id: String)

// Message
@Serializable
data class SendMessageRequest(val group_id: String, val message_content: String)

@Serializable
data class WithdrawMessageRequest(val group_id: String, val message_id: String)

@Serializable
data class GetMessageListRequest(
    val group_id: String,
    val start_time_second: Long = 0,
    val end_time_second: Long = 0
)

@Serializable
data class GetNewMessageRequest(val group_id: String)

@Serializable
data class GetFileContentRequest(val group_id: String, val message_id: String)

@Serializable
data class SendGroupNoticeRequest(val group_id: String, val message_content: String)

// AI
@Serializable
data class AiConfigRequest(
    val model_name: String = "",
    val base_url: String = "",
    val api_key: String = "",
    val role: String = "",
    val prompt: String = ""
)

// Session display model
data class SessionItem(
    val id: String,
    val type: String,
    val name: String,
    val goalUserId: String? = null,
    val userRole: String? = null,
    val hasNewMessage: Boolean = false
)
