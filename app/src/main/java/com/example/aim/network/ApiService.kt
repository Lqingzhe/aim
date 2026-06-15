package com.example.aim.network

import com.example.aim.network.models.*

object ApiService {

    suspend fun register(password: String): ApiResponse {
        return ApiClient.postLogin("/user/register", ApiClient.json.encodeToString(RegisterRequest.serializer(), RegisterRequest(password)))
    }

    suspend fun login(userId: String, password: String): ApiResponse {
        return ApiClient.postLogin("/user/login", ApiClient.json.encodeToString(LoginRequest.serializer(), LoginRequest(userId, password)))
    }

    suspend fun refreshToken(): Boolean {
        return ApiClient.refreshAccessToken()
    }

    suspend fun logoutAllDevices(): ApiResponse {
        return ApiClient.post("/user/logout-all-device", "{}")
    }

    suspend fun logoutCurrentDevice(): ApiResponse {
        return ApiClient.post("/user/logout-a-device", "{}")
    }

    suspend fun getUserInfo(): ApiResponse {
        return ApiClient.post("/user/get-user-info", "{}")
    }

    suspend fun getOtherUserInfo(goalUserId: String): ApiResponse {
        return ApiClient.post("/user/get-other-user-info", """{"goal_user_id":"$goalUserId"}""")
    }

    suspend fun updateUserInfo(request: UpdateUserInfoRequest): ApiResponse {
        return ApiClient.post("/user/update-user-info", ApiClient.json.encodeToString(UpdateUserInfoRequest.serializer(), request))
    }

    suspend fun setRemark(goalUserId: String, nickName: String): ApiResponse {
        return ApiClient.post("/user/remark", ApiClient.json.encodeToString(RemarkRequest.serializer(), RemarkRequest(goalUserId, nickName)))
    }

    // Group
    suspend fun createGroup(groupName: String): ApiResponse {
        return ApiClient.post("/group/create-group", ApiClient.json.encodeToString(CreateGroupRequest.serializer(), CreateGroupRequest(groupName)))
    }

    suspend fun deleteGroup(groupId: String): ApiResponse {
        return ApiClient.post("/group/delete-group", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    suspend fun leaveGroup(groupId: String): ApiResponse {
        return ApiClient.post("/group/leave-group", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    suspend fun getGroupInfo(groupId: String): ApiResponse {
        return ApiClient.post("/group/get-group-info", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    suspend fun changeGroupInfo(groupId: String, groupName: String): ApiResponse {
        return ApiClient.post("/group/change-group-info", ApiClient.json.encodeToString(ChangeGroupInfoRequest.serializer(), ChangeGroupInfoRequest(groupId, groupName)))
    }

    suspend fun searchGroup(groupName: String): ApiResponse {
        return ApiClient.post("/group/search-group", ApiClient.json.encodeToString(SearchGroupRequest.serializer(), SearchGroupRequest(groupName)))
    }

    suspend fun getGroupInfoWithUser(groupId: String): ApiResponse {
        return ApiClient.post("/group/get-group-info-with-user", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    suspend fun updateGroupRemark(groupId: String, remarkName: String): ApiResponse {
        return ApiClient.post("/group/update-group-info-with-user", ApiClient.json.encodeToString(UpdateGroupRemarkRequest.serializer(), UpdateGroupRemarkRequest(groupId, remarkName)))
    }

    suspend fun getGroupAndSessionId(): ApiResponse {
        return ApiClient.post("/group/get-group-and-session-id", "{}")
    }

    // Group Apply
    suspend fun setGroupApply(groupId: String): ApiResponse {
        return ApiClient.post("/group/set-group-apply", ApiClient.json.encodeToString(SetGroupApplyRequest.serializer(), SetGroupApplyRequest(groupId)))
    }

    suspend fun getGroupApplyList(groupId: String): ApiResponse {
        return ApiClient.post("/group/get-group-apply-list", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    suspend fun agreeGroupApply(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/agree-group-apply", ApiClient.json.encodeToString(AgreeGroupApplyRequest.serializer(), AgreeGroupApplyRequest(groupId, goalUserId)))
    }

    suspend fun refuseGroupApply(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/refuse-group-apply", ApiClient.json.encodeToString(AgreeGroupApplyRequest.serializer(), AgreeGroupApplyRequest(groupId, goalUserId)))
    }

    // Group Member
    suspend fun transformGroupOwner(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/transform-group-owner", ApiClient.json.encodeToString(GroupUserRequest.serializer(), GroupUserRequest(groupId, goalUserId)))
    }

    suspend fun setManager(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/set-manager", ApiClient.json.encodeToString(GroupUserRequest.serializer(), GroupUserRequest(groupId, goalUserId)))
    }

    suspend fun revokeManager(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/revoke-manager", ApiClient.json.encodeToString(GroupUserRequest.serializer(), GroupUserRequest(groupId, goalUserId)))
    }

    suspend fun kickOutGroup(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/kick-out-group", ApiClient.json.encodeToString(GroupUserRequest.serializer(), GroupUserRequest(groupId, goalUserId)))
    }

    // Mute
    suspend fun setMute(groupId: String, goalUserId: String, timeSeconds: Int, reason: String): ApiResponse {
        return ApiClient.post("/group/set-mute", ApiClient.json.encodeToString(MuteRequest.serializer(), MuteRequest(groupId, goalUserId, timeSeconds, reason)))
    }

    suspend fun releaseMute(groupId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/release-mute", ApiClient.json.encodeToString(GroupUserRequest.serializer(), GroupUserRequest(groupId, goalUserId)))
    }

    // Last Visit
    suspend fun getLastVisitTime(groupId: String): ApiResponse {
        return ApiClient.post("/group/get-last-visit-time", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    suspend fun setLastVisitTime(groupId: String): ApiResponse {
        return ApiClient.post("/group/set-last-visit-time", ApiClient.json.encodeToString(GroupIdRequest.serializer(), GroupIdRequest(groupId)))
    }

    // Friend
    suspend fun applyForFriend(goalUserId: String): ApiResponse {
        return ApiClient.post("/group/apply-for-friend", ApiClient.json.encodeToString(ApplyFriendRequest.serializer(), ApplyFriendRequest(goalUserId)))
    }

    suspend fun getFriendApplyList(): ApiResponse {
        return ApiClient.post("/group/get-friend-apply-list", "{}")
    }

    suspend fun refuseFriendApply(goalUserId: String): ApiResponse {
        return ApiClient.post("/group/refuse-friend-apply", ApiClient.json.encodeToString(ApplyFriendRequest.serializer(), ApplyFriendRequest(goalUserId)))
    }

    suspend fun createSession(goalUserId: String): ApiResponse {
        return ApiClient.post("/group/creat-session", ApiClient.json.encodeToString(CreateSessionRequest.serializer(), CreateSessionRequest(goalUserId)))
    }

    suspend fun deleteSession(sessionId: String): ApiResponse {
        return ApiClient.post("/group/delete-session", ApiClient.json.encodeToString(DeleteSessionRequest.serializer(), DeleteSessionRequest(sessionId)))
    }

    suspend fun getFriendLastVisitTime(sessionId: String, goalUserId: String): ApiResponse {
        return ApiClient.post("/group/get-friend-last-visit-time", ApiClient.json.encodeToString(FriendLastVisitRequest.serializer(), FriendLastVisitRequest(sessionId, goalUserId)))
    }

    // Message
    suspend fun sendMessage(groupId: String, content: String): ApiResponse {
        return ApiClient.post("/message/send-message", ApiClient.json.encodeToString(SendMessageRequest.serializer(), SendMessageRequest(groupId, content)))
    }

    suspend fun sendFile(groupId: String, fileName: String, fileBytes: ByteArray): ApiResponse {
        return ApiClient.postMultipart(
            "/message/send-file",
            mapOf("group_id" to groupId, "file_name" to fileName),
            "file", fileBytes, fileName
        )
    }

    suspend fun sendVoice(groupId: String, voiceTime: String, fileBytes: ByteArray): ApiResponse {
        return ApiClient.postMultipart(
            "/message/send-voice",
            mapOf("group_id" to groupId, "voice_time" to voiceTime),
            "file", fileBytes, "voice.aac"
        )
    }

    suspend fun sendPicture(groupId: String, fileBytes: ByteArray): ApiResponse {
        return ApiClient.postMultipart(
            "/message/send-picture",
            mapOf("group_id" to groupId),
            "picture", fileBytes, "picture.jpg"
        )
    }

    suspend fun withdrawMessage(groupId: String, messageId: String): ApiResponse {
        return ApiClient.post("/message/withdraw-message", ApiClient.json.encodeToString(WithdrawMessageRequest.serializer(), WithdrawMessageRequest(groupId, messageId)))
    }

    suspend fun getMessageList(groupId: String, startTime: Long, endTime: Long): ApiResponse {
        return ApiClient.post("/message/get-message-list", ApiClient.json.encodeToString(GetMessageListRequest.serializer(), GetMessageListRequest(groupId, startTime, endTime)))
    }

    suspend fun getNewMessage(groupId: String): ApiResponse {
        return ApiClient.post("/message/get-new-message", ApiClient.json.encodeToString(GetNewMessageRequest.serializer(), GetNewMessageRequest(groupId)))
    }

    suspend fun getFileContent(groupId: String, messageId: String): ByteArray? {
        return ApiClient.downloadFile("/message/get-file-content", ApiClient.json.encodeToString(GetFileContentRequest.serializer(), GetFileContentRequest(groupId, messageId)))
    }

    suspend fun sendGroupNotice(groupId: String, content: String): ApiResponse {
        return ApiClient.post("/message/send-group-notice", ApiClient.json.encodeToString(SendGroupNoticeRequest.serializer(), SendGroupNoticeRequest(groupId, content)))
    }

    // AI
    suspend fun deleteChatContext(): ApiResponse {
        return ApiClient.post("/ai/delete-chat-context", "{}")
    }

    suspend fun getAiConfig(): ApiResponse {
        return ApiClient.post("/ai/get-ai-config", "{}")
    }

    suspend fun updateAiConfig(config: AiConfigRequest): ApiResponse {
        return ApiClient.post("/ai/update-ai-config", ApiClient.json.encodeToString(AiConfigRequest.serializer(), config))
    }

    suspend fun deleteAiConfig(): ApiResponse {
        return ApiClient.post("/ai/delete-ai-config", "{}")
    }
}
