// ============================================
// 消息相关 - 添加已读状态
// ============================================

// 全局消息处理函数，由 WebSocket 调用
async function handleNewMessage(data) {
    console.log('handleNewMessage 处理新消息:', data);

    let messageId = data.message_id || data.MessageId || Date.now();
    let userId = data.user_id || data.UserId;
    let messageContent = data.message_content || data.MessageContent || data.data?.message_content || '';
    let sessionId = data.session_id || data.SessionId || data.group_id || data.GroupId;
    let sendTime = data.send_time_second || data.SendTimeSecond || Math.floor(Date.now() / 1000);

    if (!messageContent && data.data) {
        if (typeof data.data === 'string') {
            messageContent = data.data;
        } else if (data.data.message_content) {
            messageContent = data.data.message_content;
        } else {
            messageContent = JSON.stringify(data.data);
        }
    }

    console.log('解析后的消息:', { messageId, userId, messageContent, sessionId, sendTime });

    if (currentSession && String(sessionId) === String(currentSession.id)) {
        console.log('消息属于当前会话，直接显示');

        let senderDisplayName = '用户';
        if (String(userId) === String(userID)) {
            senderDisplayName = '我';
        } else if (window.currentRemarkMap && window.currentRemarkMap[String(userId)]) {
            senderDisplayName = window.currentRemarkMap[String(userId)];
        } else {
            senderDisplayName = `用户 ${userId}`;
        }

        addMessageToChat({
            message_id: messageId,
            user_id: userId,
            message_content: messageContent,
            send_time_second: sendTime,
            group_id: sessionId,
            user_name: senderDisplayName
        });

        clearSessionNewMessageMark(currentSession.id);

        // 收到新消息后，延迟更新已读状态
        setTimeout(() => {
            updateMessagesReadStatus();
        }, 500);
    } else if (sessionId) {
        console.log('消息属于其他会话，标记未读:', sessionId);
        markSessionHasNewMessage(sessionId);
        if (typeof addSystemMessage === 'function') {
            addSystemMessage(`📬 来自会话 ${sessionId} 的新消息: ${messageContent.substring(0, 50)}...`);
        }
    }

    if (typeof loadSessions === 'function') {
        loadSessions();
    }
}

async function sendMessage() {
    const input = document.getElementById('message-input');
    const content = input.value.trim();
    if (!content || !currentSession) {
        if (!currentSession) addSystemMessage('请先选择一个会话');
        return;
    }

    console.log('发送消息:', content, 'groupId:', currentSession.id);

    const result = await apiCall('POST', '/message/send-message', {
        group_id: String(currentSession.id),
        message_content: content
    });

    console.log('发送消息结果:', result);

    if (result && result.code === 0) {
        input.value = '';
        const messageId = result.data?.message_info?.message_id || Date.now();
        const sendTime = Math.floor(Date.now() / 1000);

        addMessageToChat({
            message_id: messageId,
            user_id: userID,
            message_content: content,
            send_time_second: sendTime,
            group_id: currentSession.id,
            user_name: '我'
        });
    } else {
        addSystemMessage('发送失败: ' + (result?.message || '未知错误'));
    }
}

// ============================================
// 消息相关 - 添加最后访问时间更新
// ============================================

// 更新用户在会话中的最后访问时间（群聊和私聊都调用）
async function updateLastVisitTime(sessionId, sessionType) {
    if (!sessionId) {
        console.log('updateLastVisitTime: sessionId 为空，跳过');
        return;
    }

    console.log('========== 调用 updateLastVisitTime ==========');
    console.log('sessionId:', sessionId, 'sessionType:', sessionType);

    try {
        // 群聊和私聊都使用 set-last-visit-time 接口
        const result = await apiCall('POST', '/group/set-last-visit-time', {
            group_id: String(sessionId)
        });
        console.log('updateLastVisitTime 结果:', JSON.stringify(result, null, 2));
    } catch (error) {
        console.error('updateLastVisitTime 失败:', error);
    }
}

async function loadHistoryMessages(groupId, loadMore = false) {
    if (currentGroupId !== groupId) {
        currentGroupId = groupId;
        hasMoreMessages = true;
    }

    if (isLoadingMore) return;
    isLoadingMore = true;

    let endTime = Math.floor(Date.now() / 1000);
    if (loadMore && hasMoreMessages) {
        const firstMessage = document.querySelector('.message:first-child');
        if (firstMessage && firstMessage.dataset.time) {
            endTime = parseInt(firstMessage.dataset.time);
        }
    }

    const result = await apiCall('POST', '/message/get-message-list', {
        group_id: String(groupId),
        start_time_second: 0,
        end_time_second: endTime
    });

    isLoadingMore = false;

    if (result && result.code === 0 && result.data?.message_info?.message_list) {
        const messages = result.data.message_info.message_list;

        if (messages.length === 0) {
            if (loadMore) addSystemMessage('没有更早的消息了');
            hasMoreMessages = false;
            if (!loadMore) addSystemMessage('暂无消息');
            return;
        }

        hasMoreMessages = messages.length >= 50;

        if (loadMore) {
            for (let i = messages.length - 1; i >= 0; i--) {
                const msg = messages[i];
                addMessageToChatTop(msg, groupId);
            }
            addSystemMessage(`已加载 ${messages.length} 条历史消息`);
        } else {
            const container = document.getElementById('messages-container');
            container.innerHTML = '';
            for (let i = messages.length - 1; i >= 0; i--) {
                const msg = messages[i];
                let senderName = `用户 ${msg.user_id}`;
                if (String(msg.user_id) === String(userID)) {
                    senderName = '我';
                } else if (window.currentRemarkMap && window.currentRemarkMap[String(msg.user_id)]) {
                    senderName = window.currentRemarkMap[String(msg.user_id)];
                } else if (msg.user_name) {
                    senderName = msg.user_name;
                }
                addMessageToChat({
                    message_id: msg.message_id,
                    user_id: msg.user_id,
                    message_content: msg.message_content,
                    send_time_second: msg.send_time_second,
                    group_id: groupId,
                    user_name: senderName
                });
            }
            addSystemMessage(`共加载 ${messages.length} 条消息`);

            // 拉取消息成功后，更新最后访问时间（群聊和私聊都调用）
            console.log('准备调用 updateLastVisitTime, sessionId:', groupId, 'currentSession.type:', currentSession?.type);
            await updateLastVisitTime(groupId, currentSession?.type);

            // 消息加载完成后，更新已读状态
            setTimeout(() => {
                if (typeof updateMessagesReadStatus === 'function') {
                    updateMessagesReadStatus();
                }
            }, 500);
        }

        const container = document.getElementById('messages-container');
        container.onscroll = () => {
            if (container.scrollTop === 0 && !isLoadingMore && hasMoreMessages) {
                loadHistoryMessages(groupId, true);
            }
        };
    } else if (!loadMore) {
        addSystemMessage('暂无消息');
    }
}

// 获取群聊中各成员的最后访问时间
async function getGroupReadStatus(groupId) {
    console.log('调用 getGroupReadStatus, groupId:', groupId);
    try {
        const result = await apiCall('POST', '/group/get-last-visit-time', { group_id: String(groupId) });
        console.log('群聊已读状态返回:', result);

        if (result && result.code === 0 && result.data?.group_info?.last_visit_time) {
            return result.data.group_info.last_visit_time;
        }
        return {};
    } catch (error) {
        console.error('获取群聊已读状态失败:', error);
        return {};
    }
}

// 获取私聊的已读状态
async function getPrivateReadStatus(sessionId) {
    console.log('调用 getPrivateReadStatus, sessionId:', sessionId);
    try {
        const session = sessions.find(s => s.id === String(sessionId));
        if (!session || !session.goalUserId) {
            console.log('找不到对方用户ID');
            return 0;
        }

        console.log('对方用户ID:', session.goalUserId);

        const result = await apiCall('POST', '/group/get-friend-last-visit-time', {
            session_id: String(sessionId),
            goal_user_id: session.goalUserId
        });
        console.log('私聊已读状态返回:', result);

        if (result && result.code === 0 && result.data?.session_info?.last_visit_time) {
            return result.data.session_info.last_visit_time;
        }
        return 0;
    } catch (error) {
        console.error('获取私聊已读状态失败:', error);
        return 0;
    }
}

// 判断消息是否已读
function checkIfMessageRead(message, readStatusMap, session) {
    const sendTime = message.send_time_second;

    if (session.type === 'group') {
        const currentUserId = String(userID);
        let readCount = 0;
        let totalOthers = 0;

        for (const [userId, lastVisitTime] of Object.entries(readStatusMap)) {
            if (userId !== currentUserId) {
                totalOthers++;
                if (lastVisitTime >= sendTime) {
                    readCount++;
                }
            }
        }

        if (totalOthers === 0) {
            return 0;
        }
        if (readCount === totalOthers) {
            return true;
        }
        return readCount;
    } else {
        const friendLastVisit = readStatusMap;
        return friendLastVisit >= sendTime;
    }
}

// 更新群聊消息的已读状态
async function updateGroupReadStatus(groupId) {
    console.log('updateGroupReadStatus 被调用, groupId:', groupId);
    if (!groupId || currentSession?.type !== 'group') {
        return;
    }

    const readStatusMap = await getGroupReadStatus(groupId);

    if (Object.keys(readStatusMap).length === 0) {
        return;
    }

    const messages = document.querySelectorAll('.message');

    for (const msgElement of messages) {
        const sendTime = parseInt(msgElement.dataset.time);
        const userId = msgElement.dataset.userId;

        if (userId !== String(userID)) {
            const readResult = checkIfMessageRead({ send_time_second: sendTime }, readStatusMap, currentSession);
            const readSpan = msgElement.querySelector('.message-read-status');

            if (readSpan) {
                if (readResult === true) {
                    readSpan.textContent = '✓ 全部已读';
                    readSpan.style.color = '#10b981';
                } else if (typeof readResult === 'number') {
                    readSpan.textContent = `${readResult} 人已读`;
                    readSpan.style.color = readResult > 0 ? '#10b981' : '#9ca3af';
                }
            }
        }
    }
}

// 更新私聊消息的已读状态
async function updatePrivateReadStatus(sessionId) {
    console.log('updatePrivateReadStatus 被调用, sessionId:', sessionId);
    if (!sessionId || currentSession?.type !== 'session') {
        return;
    }

    const friendLastVisit = await getPrivateReadStatus(sessionId);

    if (friendLastVisit === 0) {
        return;
    }

    const messages = document.querySelectorAll('.message.self');

    for (const msgElement of messages) {
        const sendTime = parseInt(msgElement.dataset.time);
        const readSpan = msgElement.querySelector('.message-read-status');

        if (readSpan && friendLastVisit >= sendTime) {
            readSpan.textContent = '✓ 已读';
            readSpan.style.color = '#10b981';
        }
    }
}

// 更新所有消息的已读状态（统一入口）
async function updateMessagesReadStatus() {
    console.log('updateMessagesReadStatus 被调用, currentSession:', currentSession);
    if (!currentSession) return;

    if (currentSession.type === 'group') {
        await updateGroupReadStatus(currentSession.id);
    } else if (currentSession.type === 'session') {
        await updatePrivateReadStatus(currentSession.id);
    }
}

async function withdrawMessage(messageId, groupId) {
    if (!messageId || !groupId) {
        addSystemMessage('撤回失败：缺少消息ID或群ID');
        return;
    }

    const result = await apiCall('POST', '/message/withdraw-message', {
        group_id: String(groupId),
        message_id: String(messageId)
    });

    if (result && result.code === 0) {
        const msgElement = document.querySelector(`.message[data-message-id="${messageId}"]`);
        if (msgElement) {
            const bubble = msgElement.querySelector('.message-bubble');
            bubble.textContent = '⚠️ 该消息已撤回';
            bubble.style.opacity = '0.6';
            bubble.style.fontStyle = 'italic';
            const withdrawBtn = msgElement.querySelector('.withdraw-btn');
            if (withdrawBtn) withdrawBtn.remove();
            const readSpan = msgElement.querySelector('.message-read-status');
            if (readSpan) readSpan.remove();
        }
        addSystemMessage('消息已撤回');
    } else {
        addSystemMessage('撤回失败: ' + (result?.message || '未知错误'));
    }
}

async function withdrawLatestMessage() {
    if (!currentSession) {
        addSystemMessage('请先选择一个会话');
        return;
    }

    const container = document.getElementById('messages-container');
    const messages = container.querySelectorAll(`.message[data-user-id="${userID}"]`);
    if (messages.length === 0) {
        addSystemMessage('没有可以撤回的消息');
        return;
    }

    const latestSelfMsg = messages[messages.length - 1];
    const messageId = latestSelfMsg.dataset.messageId;

    if (!messageId) {
        addSystemMessage('无法获取消息ID');
        return;
    }

    await withdrawMessage(messageId, currentSession.id);
}

function addMessageToChat(message) {
    const container = document.getElementById('messages-container');
    const isSelf = String(message.user_id) === String(userID);
    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${isSelf ? 'self' : 'other'}`;
    msgDiv.dataset.messageId = message.message_id;
    msgDiv.dataset.userId = message.user_id;
    msgDiv.dataset.time = message.send_time_second;

    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    bubble.textContent = message.message_content || '';

    const meta = document.createElement('div');
    meta.className = 'message-meta';
    const time = message.send_time_second
        ? new Date(message.send_time_second * 1000).toLocaleTimeString()
        : new Date().toLocaleTimeString();

    let senderName;
    if (isSelf) {
        senderName = '我';
    } else if (message.user_name && message.user_name !== `用户 ${message.user_id}`) {
        senderName = message.user_name;
    } else {
        senderName = `用户 ${message.user_id}`;
    }

    // 添加已读状态显示（初始显示加载中，后续更新）
    let readStatusHtml = '';
    if (isSelf) {
        readStatusHtml = '<span class="message-read-status" style="margin-left:8px;font-size:10px;color:#9ca3af;">未读</span>';
    } else if (message.isRead !== undefined) {
        if (message.isRead === true) {
            readStatusHtml = '<span class="message-read-status" style="margin-left:8px;font-size:10px;color:#10b981;">✓ 全部已读</span>';
        } else if (typeof message.isRead === 'number') {
            readStatusHtml = `<span class="message-read-status" style="margin-left:8px;font-size:10px;color:${message.isRead > 0 ? '#10b981' : '#9ca3af'};">${message.isRead} 人已读</span>`;
        }
    }

    meta.innerHTML = `${senderName} ${time}${readStatusHtml}`;

    msgDiv.appendChild(bubble);
    msgDiv.appendChild(meta);
    container.appendChild(msgDiv);
    container.scrollTop = container.scrollHeight;
}

function addMessageToChatTop(message, groupId) {
    const container = document.getElementById('messages-container');
    const isSelf = String(message.user_id) === String(userID);
    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${isSelf ? 'self' : 'other'}`;
    msgDiv.dataset.messageId = message.message_id;
    msgDiv.dataset.userId = message.user_id;
    msgDiv.dataset.time = message.send_time_second;

    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    bubble.textContent = message.message_content || '';

    const meta = document.createElement('div');
    meta.className = 'message-meta';
    const time = message.send_time_second
        ? new Date(message.send_time_second * 1000).toLocaleString()
        : new Date().toLocaleString();

    let senderName;
    if (isSelf) {
        senderName = '我';
    } else if (message.user_name) {
        senderName = message.user_name;
    } else {
        senderName = `用户 ${message.user_id}`;
    }

    let readStatusHtml = '';
    if (isSelf) {
        readStatusHtml = '<span class="message-read-status" style="margin-left:8px;font-size:10px;color:#9ca3af;">未读</span>';
    }

    meta.innerHTML = `${senderName} ${time}${readStatusHtml}`;

    msgDiv.appendChild(bubble);
    msgDiv.appendChild(meta);
    container.insertBefore(msgDiv, container.firstChild);
}