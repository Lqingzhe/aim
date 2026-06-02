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
        // 自己的消息直接显示
        addMessageToChat({
            message_id: result.data?.message_info?.message_id,
            user_id: userID,
            message_content: content,
            send_time_second: Math.floor(Date.now() / 1000),
            group_id: currentSession.id
        });
    } else {
        addSystemMessage('发送失败: ' + (result?.message || '未知错误'));
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
                addMessageToChat({
                    message_id: msg.message_id,
                    user_id: String(msg.user_id),
                    message_content: msg.message_content,
                    send_time_second: msg.send_time_second,
                    group_id: groupId
                });
            }
            addSystemMessage(`共加载 ${messages.length} 条消息`);
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