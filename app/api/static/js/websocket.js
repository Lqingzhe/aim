// ============================================
// WebSocket 连接和消息处理
// ============================================

let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;

function connectWebSocket() {
    const token = sessionStorage.getItem('access_token');
    if (!token) {
        console.log('没有 token，跳过 WebSocket 连接');
        return;
    }

    const wsUrl = `${WS_URL}?token=${encodeURIComponent(token)}`;
    console.log('WebSocket 连接 URL:', wsUrl);

    window.webSocket = new WebSocket(wsUrl);

    window.webSocket.onopen = () => {
        console.log('WebSocket 已连接');
        reconnectAttempts = 0; // 连接成功，重置重连次数
        if (typeof addSystemMessage === 'function') {
            addSystemMessage('✅ WebSocket 连接成功');
        }
        startHeartbeatTimer();
    };

    window.webSocket.onmessage = (event) => {
        console.log('WebSocket 收到消息:', event.data);

        if (!event.data || event.data.trim() === '') {
            return;
        }

        try {
            const data = JSON.parse(event.data);
            console.log('解析后:', data);

            // 心跳响应
            if (data.type === 'pong') {
                console.log('收到心跳响应');
                return;
            }

            // logout 响应
            if (data.type === 'logout') {
                if (typeof addSystemMessage === 'function') {
                    addSystemMessage('🔚 您已退出登录');
                }
                setTimeout(() => {
                    window.location.href = '/login';
                }, 1000);
                return;
            }

            // 处理新消息通知
            if (data.message_code === 'group_message' || data.type === 'new_message') {
                console.log('收到新消息通知，调用 handleNewMessage');
                if (typeof handleNewMessage === 'function') {
                    handleNewMessage(data);
                }

                // 收到新消息后，延迟更新自己消息的已读状态
                if (currentSession && String(data.session_id || data.group_id) === String(currentSession.id)) {
                    setTimeout(() => {
                        if (typeof updateMyMessagesReadStatus === 'function') {
                            console.log('收到新消息，更新自己消息的已读状态');
                            updateMyMessagesReadStatus();
                        }
                    }, 1000);
                }
            }

            // 处理系统消息通知
            if (data.message_code === 'friend_request' || data.message_code === 'group_apply' ||
                data.message_code === 'group_join' || data.message_code === 'group_leave' ||
                data.message_code === 'group_kick' || data.message_code === 'group_disband') {
                console.log('收到系统通知:', data.message_code);
                if (typeof addSystemMessage === 'function') {
                    addSystemMessage(`📢 [系统通知] ${JSON.stringify(data)}`);
                }
                if (typeof loadSessions === 'function') {
                    loadSessions();
                }
            }

            // 其他所有消息，原样输出到聊天框
            if (typeof addSystemMessage === 'function') {
                let displayText = event.data;
                if (typeof data === 'object' && data.message_code !== 'group_message') {
                    displayText = JSON.stringify(data, null, 2);
                }
                addSystemMessage(`📨 ${displayText}`);
            }

        } catch (e) {
            console.log('消息不是 JSON 格式，直接输出');
            if (typeof addSystemMessage === 'function') {
                addSystemMessage(`📨 ${event.data}`);
            }
        }
    };

    window.webSocket.onclose = async (event) => {
        console.log('WebSocket 已断开, code:', event.code, 'reason:', event.reason);

        // 检查是否是 token 过期导致的关闭（后端关闭连接时可能会发送特定信息）
        // 或者直接尝试刷新 token
        if (sessionStorage.getItem('access_token')) {
            console.log('尝试刷新 token...');

            // 尝试刷新 token
            const refreshed = await refreshAccessToken();

            if (refreshed) {
                console.log('Token 刷新成功，使用新 token 重新连接 WebSocket');
                if (typeof addSystemMessage === 'function') {
                    addSystemMessage('🔄 Token 已刷新，正在重新连接...');
                }
                stopHeartbeatTimer();
                // 使用新 token 重新连接
                connectWebSocket();
            } else {
                console.log('Token 刷新失败，跳转到登录页');
                if (typeof addSystemMessage === 'function') {
                    addSystemMessage('❌ Token 刷新失败，请重新登录');
                }
                sessionStorage.clear();
                window.location.href = '/login';
            }
        } else {
            console.log('没有 token，跳转到登录页');
            if (typeof addSystemMessage === 'function') {
                addSystemMessage('⚠️ WebSocket 连接已断开，请重新登录');
            }
            window.location.href = '/login';
        }
    };

    window.webSocket.onerror = (error) => {
        console.error('WebSocket 错误:', error);
        if (typeof addSystemMessage === 'function') {
            addSystemMessage('❌ WebSocket 连接错误');
        }
    };
}

function startHeartbeatTimer() {
    if (window.heartbeatTimer) {
        clearInterval(window.heartbeatTimer);
    }

    window.heartbeatTimer = setInterval(() => {
        if (window.webSocket && window.webSocket.readyState === WebSocket.OPEN) {
            window.webSocket.send(JSON.stringify({ type: 'ping' }));
            console.log('发送心跳 ping');
        }
    }, 30000);
}

function stopHeartbeatTimer() {
    if (window.heartbeatTimer) {
        clearInterval(window.heartbeatTimer);
        window.heartbeatTimer = null;
    }
}

// 页面关闭前清理
window.addEventListener('beforeunload', () => {
    if (window.webSocket && window.webSocket.readyState === WebSocket.OPEN) {
        window.webSocket.close();
    }
    stopHeartbeatTimer();
});