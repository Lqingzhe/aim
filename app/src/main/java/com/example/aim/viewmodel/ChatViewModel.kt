package com.example.aim.viewmodel

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.aim.network.ApiService
import com.example.aim.network.WebSocketManager
import com.example.aim.network.models.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import kotlinx.serialization.json.Json

class ChatViewModel : ViewModel() {
    private val _uiState = MutableStateFlow(ChatUiState())
    val uiState: StateFlow<ChatUiState> = _uiState

    private val _messages = MutableStateFlow<List<MessageData>>(emptyList())
    val messages: StateFlow<List<MessageData>> = _messages

    private val _sessions = MutableStateFlow<List<SessionItem>>(emptyList())
    val sessions: StateFlow<List<SessionItem>> = _sessions

    private val _friendRequests = MutableStateFlow<List<String>>(emptyList())
    val friendRequests: StateFlow<List<String>> = _friendRequests

    private val _userInfo = MutableStateFlow<UserInfo?>(null)
    val userInfo: StateFlow<UserInfo?> = _userInfo

    private val _currentSession = MutableStateFlow<SessionItem?>(null)
    val currentSession: StateFlow<SessionItem?> = _currentSession

    private val _remarkMap = MutableStateFlow<Map<String, String>>(emptyMap())
    val remarkMap: StateFlow<Map<String, String>> = _remarkMap

    private val _readStatusMap = MutableStateFlow<Map<String, ReadStatus>>(emptyMap())
    val readStatusMap: StateFlow<Map<String, ReadStatus>> = _readStatusMap

    private val _toastMessage = MutableStateFlow<String?>(null)
    val toastMessage: StateFlow<String?> = _toastMessage

    private val json = Json { ignoreUnknownKeys = true; coerceInputValues = true }

    fun showToast(message: String) {
        _toastMessage.value = message
    }

    fun clearToast() {
        _toastMessage.value = null
    }

    fun loadSessionsAndGroups() {
        viewModelScope.launch {
            try {
                val userInfoResponse = ApiService.getUserInfo()
                if (userInfoResponse.code == 0) {
                    val remarkMap = mutableMapOf<String, String>()
                    userInfoResponse.data?.user_info?.RemarkInfos?.forEach { remark ->
                        val goalId = remark.resolvedGoalUserId()
                        val nickName = remark.resolvedNickName()
                        if (goalId.isNotEmpty() && nickName.isNotEmpty()) {
                            remarkMap[goalId] = nickName
                        }
                    }
                    _remarkMap.value = remarkMap
                }
            } catch (e: Exception) {
                e.printStackTrace()
            }

            try {
                val response = ApiService.getGroupAndSessionId()
                android.util.Log.d("ChatVM", "getGroupAndSessionId code=${response.code} data=${response.data}")
                android.util.Log.d("ChatVM", "group_info=${response.data?.group_info}")
                android.util.Log.d("ChatVM", "session_info=${response.data?.session_info}")

                if (response.code == 0) {
                    val sessionItems = mutableListOf<SessionItem>()

                    val groupIdList = response.data?.group_info?.group_id_list ?: emptyList()
                    android.util.Log.d("ChatVM", "groupIdList size=${groupIdList.size}")
                    for (id in groupIdList) {
                        var groupName = "群聊 $id"
                        var userRole = "Member"
                        try {
                            val groupInfo = ApiService.getGroupInfo(id)
                            if (groupInfo.code == 0) {
                                groupName = groupInfo.data?.group_info?.group_name ?: "群聊 $id"
                            }
                        } catch (_: Exception) {}

                        try {
                            val groupUserInfo = ApiService.getGroupInfoWithUser(id)
                            if (groupUserInfo.code == 0) {
                                userRole = groupUserInfo.data?.group_info?.group_role ?: "Member"
                            }
                        } catch (_: Exception) {}

                        sessionItems.add(SessionItem(id, "group", groupName, null, userRole))
                    }

                    val sessionIdList = response.data?.session_info?.session_id_list ?: emptyList()
                    val userIdList = response.data?.session_info?.user_of_session_id_list ?: emptyList()
                    android.util.Log.d("ChatVM", "sessionIdList size=${sessionIdList.size}, userIdList size=${userIdList.size}")

                    for (i in sessionIdList.indices) {
                        val sessionId = sessionIdList[i]
                        val goalUserId = userIdList.getOrElse(i) { "" }

                        var displayName = "私聊 $sessionId"
                        val goalUserIdStr = goalUserId.toString()

                        if (_remarkMap.value[goalUserIdStr] != null) {
                            displayName = _remarkMap.value[goalUserIdStr]!!
                        } else if (goalUserId.isNotEmpty()) {
                            try {
                                val otherInfo = ApiService.getOtherUserInfo(goalUserIdStr)
                                if (otherInfo.code == 0) {
                                    val userName = otherInfo.data?.user_info?.UserInfo?.user_name
                                        ?: otherInfo.data?.user_info?.user_name
                                    if (!userName.isNullOrEmpty()) {
                                        displayName = userName
                                    }
                                }
                            } catch (_: Exception) {}
                        }

                        sessionItems.add(SessionItem(sessionId, "session", displayName, goalUserIdStr))
                    }

                    _sessions.value = sessionItems
                }
            } catch (_: Exception) {}
        }
    }

    fun selectSession(session: SessionItem) {
        _currentSession.value = session
        loadMessages(session.id)
    }

    fun loadMessages(groupId: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                val now = System.currentTimeMillis() / 1000
                val response = ApiService.getMessageList(groupId, 0, now)
                android.util.Log.d("ChatVM", "loadMessages code=${response.code} msg=${response.message}")
                if (response.code == 0) {
                    val rawMessages = response.data?.message_info?.message_list ?: emptyList()
                    android.util.Log.d("ChatVM", "messages count=${rawMessages.size}")
                    if (rawMessages.isNotEmpty()) {
                        android.util.Log.d("ChatVM", "first msg: id=${rawMessages[0].id()} content=${rawMessages[0].content().take(50)}")
                    }
                    _messages.value = rawMessages.reversed()
                    loadReadStatus(groupId)
                }
                ApiService.setLastVisitTime(groupId)
            } catch (e: Exception) {
                android.util.Log.e("ChatVM", "loadMessages error", e)
            } finally {
                _uiState.value = _uiState.value.copy(isLoading = false)
            }
        }
    }

    private fun resolveSenderName(senderId: String): String {
        if (senderId == com.example.aim.data.TokenStore.userId) return "我"
        _remarkMap.value[senderId]?.let { return it }
        return "用户 $senderId"
    }

    private fun loadReadStatus(groupId: String) {
        viewModelScope.launch {
            try {
                val session = _currentSession.value ?: return@launch
                val myId = com.example.aim.data.TokenStore.userId
                val myMessages = _messages.value.filter { it.sender() == myId }

                if (session.type == "group") {
                    val response = ApiService.getLastVisitTime(groupId)
                    if (response.code == 0) {
                        val visitTimes = response.data?.group_info?.last_visit_time ?: emptyMap()
                        val totalOthers = visitTimes.keys.size - 1
                        val statusMap = mutableMapOf<String, ReadStatus>()
                        for (msg in myMessages) {
                            val msgTime = msg.time()
                            var readCount = 0
                            for ((uid, visitTime) in visitTimes) {
                                if (uid != myId && visitTime >= msgTime) readCount++
                            }
                            statusMap[msg.id()] = ReadStatus(readCount, totalOthers)
                        }
                        _readStatusMap.value = statusMap
                    }
                } else if (session.type == "session" && session.goalUserId != null) {
                    val response = ApiService.getFriendLastVisitTime(groupId, session.goalUserId)
                    if (response.code == 0) {
                        val friendVisitTime = response.data?.session_info?.last_visit_time ?: 0
                        val statusMap = mutableMapOf<String, ReadStatus>()
                        for (msg in myMessages) {
                            val isRead = friendVisitTime >= msg.time()
                            statusMap[msg.id()] = ReadStatus(if (isRead) 1 else 0, 1)
                        }
                        _readStatusMap.value = statusMap
                    }
                }
            } catch (_: Exception) {}
        }
    }

    data class ReadStatus(val readCount: Int, val totalOthers: Int) {
        fun displayText(): String = when {
            totalOthers == 0 -> "未读"
            readCount > 0 -> "${readCount}人已读"
            else -> "未读"
        }
        fun isRead(): Boolean = readCount > 0
    }

    fun sendMessage(groupId: String, content: String, onResult: (Boolean) -> Unit) {
        viewModelScope.launch {
            try {
                val response = ApiService.sendMessage(groupId, content)
                if (response.code == 0) {
                    val newMsg = MessageData(
                        message_id = response.data?.message_info?.message_id ?: "",
                        group_id = groupId,
                        sender_id = com.example.aim.data.TokenStore.userId,
                        message_content = content,
                        send_time_second = System.currentTimeMillis() / 1000,
                        user_id = com.example.aim.data.TokenStore.userId
                    )
                    _messages.value = _messages.value + newMsg
                    onResult(true)
                } else {
                    showToast(response.message.ifEmpty { "发送失败" })
                    onResult(false)
                }
            } catch (_: Exception) {
                showToast("网络错误")
                onResult(false)
            }
        }
    }

    fun handleWsMessage(rawMessage: String) {
        try {
            val data = json.decodeFromString<Map<String, kotlinx.serialization.json.JsonElement>>(rawMessage)
            val type = data["type"]?.toString()?.trim('"') ?: ""
            val messageCode = data["message_code"]?.toString()?.trim('"') ?: ""

            if (type == "pong") return
            if (type == "logout") {
                showToast("已被踢出登录")
                return
            }

            if (messageCode == "group_message" || type == "new_message") {
                val senderId = data["user_id"]?.toString()?.trim('"') ?: ""
                val content = data["message_content"]?.toString()?.trim('"') ?: ""
                val sessionId = data["session_id"]?.toString()?.trim('"')
                    ?: data["group_id"]?.toString()?.trim('"') ?: ""
                val msgId = data["message_id"]?.toString()?.trim('"') ?: ""
                val sendTime = data["send_time_second"]?.toString()?.trim('"')?.toLongOrNull()
                    ?: System.currentTimeMillis() / 1000

                val current = _currentSession.value
                if (current != null && sessionId == current.id) {
                    val newMsg = MessageData(
                        message_id = msgId,
                        group_id = sessionId,
                        sender_id = senderId,
                        message_content = content,
                        send_time_second = sendTime,
                        user_id = senderId
                    )
                    _messages.value = _messages.value + newMsg
                } else {
                    showToast("来自 $sessionId 的新消息")
                }
            }
        } catch (_: Exception) {}
    }

    fun withdrawMessage(groupId: String, messageId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.withdrawMessage(groupId, messageId)
                if (response.code == 0) {
                    _messages.value = _messages.value.map {
                        if (it.id() == messageId) it.copy(is_withdrawn = true) else it
                    }
                    showToast("消息已撤回")
                } else {
                    showToast(response.message.ifEmpty { "撤回失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun loadUserInfo() {
        viewModelScope.launch {
            try {
                val response = ApiService.getUserInfo()
                if (response.code == 0) {
                    _userInfo.value = response.data?.user_info?.resolveUserInfo()
                }
            } catch (_: Exception) {}
        }
    }

    fun loadFriendRequests() {
        viewModelScope.launch {
            try {
                val response = ApiService.getFriendApplyList()
                if (response.code == 0) {
                    _friendRequests.value = response.data?.session_info?.apply_user_list ?: emptyList()
                }
            } catch (_: Exception) {}
        }
    }

    fun applyFriend(goalUserId: String, onResult: (Boolean, String) -> Unit) {
        viewModelScope.launch {
            try {
                val response = ApiService.applyForFriend(goalUserId)
                onResult(response.code == 0, response.message.ifEmpty { if (response.code == 0) "已发送" else "发送失败" })
            } catch (e: Exception) {
                onResult(false, "网络错误")
            }
        }
    }

    fun acceptFriend(goalUserId: String) {
        viewModelScope.launch {
            try {
                ApiService.createSession(goalUserId)
                loadSessionsAndGroups()
                loadFriendRequests()
                showToast("已同意好友申请")
            } catch (_: Exception) {}
        }
    }

    fun refuseFriend(goalUserId: String) {
        viewModelScope.launch {
            try {
                ApiService.refuseFriendApply(goalUserId)
                loadFriendRequests()
                showToast("已拒绝")
            } catch (_: Exception) {}
        }
    }

    fun createGroup(groupName: String, onResult: (Boolean, String) -> Unit) {
        viewModelScope.launch {
            try {
                val response = ApiService.createGroup(groupName)
                if (response.code == 0) {
                    val groupId = response.data?.group_info?.group_id ?: ""
                    loadSessionsAndGroups()
                    onResult(true, groupId)
                } else {
                    onResult(false, response.message.ifEmpty { "创建失败" })
                }
            } catch (e: Exception) {
                onResult(false, "网络错误")
            }
        }
    }

    fun joinGroup(groupId: String, onResult: (Boolean, String) -> Unit) {
        viewModelScope.launch {
            try {
                val response = ApiService.setGroupApply(groupId)
                onResult(response.code == 0, response.message.ifEmpty { if (response.code == 0) "已申请" else "申请失败" })
            } catch (e: Exception) {
                onResult(false, "网络错误")
            }
        }
    }

    fun updateUserInfo(request: UpdateUserInfoRequest, onResult: (Boolean) -> Unit) {
        viewModelScope.launch {
            try {
                val response = ApiService.updateUserInfo(request)
                if (response.code == 0) {
                    loadUserInfo()
                    onResult(true)
                } else {
                    onResult(false)
                }
            } catch (_: Exception) {
                onResult(false)
            }
        }
    }

    fun setRemark(goalUserId: String, nickName: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.setRemark(goalUserId, nickName)
                if (response.code == 0) {
                    showToast("备注设置成功")
                    loadSessionsAndGroups()
                } else {
                    showToast(response.message.ifEmpty { "设置失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun deleteSession(sessionId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.deleteSession(sessionId)
                if (response.code == 0) {
                    showToast("已删除")
                    if (_currentSession.value?.id == sessionId) {
                        _currentSession.value = null
                        _messages.value = emptyList()
                    }
                    loadSessionsAndGroups()
                } else {
                    showToast(response.message.ifEmpty { "删除失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun deleteGroup(groupId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.deleteGroup(groupId)
                if (response.code == 0) {
                    showToast("群聊已解散")
                    if (_currentSession.value?.id == groupId) {
                        _currentSession.value = null
                        _messages.value = emptyList()
                    }
                    loadSessionsAndGroups()
                } else {
                    showToast(response.message.ifEmpty { "解散失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun leaveGroup(groupId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.leaveGroup(groupId)
                if (response.code == 0) {
                    showToast("已退出群聊")
                    if (_currentSession.value?.id == groupId) {
                        _currentSession.value = null
                        _messages.value = emptyList()
                    }
                    loadSessionsAndGroups()
                } else {
                    showToast(response.message.ifEmpty { "退出失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun sendGroupNotice(groupId: String, content: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.sendGroupNotice(groupId, content)
                if (response.code == 0) {
                    showToast("通知已发送")
                } else {
                    showToast(response.message.ifEmpty { "发送失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun kickOutGroup(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.kickOutGroup(groupId, goalUserId)
                if (response.code == 0) {
                    showToast("已踢出")
                } else {
                    showToast(response.message.ifEmpty { "操作失败" })
                }
            } catch (_: Exception) {
                showToast("网络错误")
            }
        }
    }

    fun setManager(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.setManager(groupId, goalUserId)
                showToast(if (response.code == 0) "设置成功" else response.message.ifEmpty { "操作失败" })
            } catch (_: Exception) { showToast("网络错误") }
        }
    }

    fun revokeManager(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.revokeManager(groupId, goalUserId)
                showToast(if (response.code == 0) "罢免成功" else response.message.ifEmpty { "操作失败" })
            } catch (_: Exception) { showToast("网络错误") }
        }
    }

    fun transformGroupOwner(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.transformGroupOwner(groupId, goalUserId)
                if (response.code == 0) {
                    showToast("转让成功")
                    loadSessionsAndGroups()
                } else {
                    showToast(response.message.ifEmpty { "操作失败" })
                }
            } catch (_: Exception) { showToast("网络错误") }
        }
    }

    fun setMute(groupId: String, goalUserId: String, timeSeconds: Int, reason: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.setMute(groupId, goalUserId, timeSeconds, reason)
                showToast(if (response.code == 0) "禁言成功" else response.message.ifEmpty { "操作失败" })
            } catch (_: Exception) { showToast("网络错误") }
        }
    }

    fun releaseMute(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.releaseMute(groupId, goalUserId)
                showToast(if (response.code == 0) "解除禁言成功" else response.message.ifEmpty { "操作失败" })
            } catch (_: Exception) { showToast("网络错误") }
        }
    }

    fun logoutCurrentDevice(onDone: () -> Unit) {
        viewModelScope.launch {
            try { ApiService.logoutCurrentDevice() } catch (_: Exception) {}
            onDone()
        }
    }

    fun logoutAllDevices(onDone: () -> Unit) {
        viewModelScope.launch {
            try { ApiService.logoutAllDevices() } catch (_: Exception) {}
            onDone()
        }
    }

    fun loadGroupApplyList(groupId: String, onResult: (List<String>) -> Unit) {
        viewModelScope.launch {
            try {
                val response = ApiService.getGroupApplyList(groupId)
                if (response.code == 0) {
                    onResult(response.data?.group_info?.group_id_list ?: emptyList())
                } else {
                    onResult(emptyList())
                }
            } catch (_: Exception) { onResult(emptyList()) }
        }
    }

    fun agreeGroupApply(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.agreeGroupApply(groupId, goalUserId)
                showToast(if (response.code == 0) "已同意" else response.message.ifEmpty { "操作失败" })
            } catch (_: Exception) { showToast("网络错误") }
        }
    }

    fun refuseGroupApply(groupId: String, goalUserId: String) {
        viewModelScope.launch {
            try {
                val response = ApiService.refuseGroupApply(groupId, goalUserId)
                showToast(if (response.code == 0) "已拒绝" else response.message.ifEmpty { "操作失败" })
            } catch (_: Exception) { showToast("网络错误") }
        }
    }
}

data class ChatUiState(
    val isLoading: Boolean = false
)
