// ============================================
// WebSocket 连接和消息处理
// ============================================

// 使用不同的变量名避免冲突
let webSocket = null;
let heartbeatTimer = null;

function connectWebSocket() {
    const token = sessionStorage.getItem('access_token');
    if (!token) {
        console.log('没有 token，跳过 WebSocket 连接');
        return;
    }

    const wsUrl = `${WS_URL}?token=${encodeURIComponent(token)}`;
    console.log('WebSocket 连接 URL:', wsUrl);

    webSocket = new WebSocket(wsUrl);

    webSocket.onopen = () => {
        console.log('WebSocket 已连接');
        if (typeof addSystemMessage === 'function') {
            addSystemMessage('✅ WebSocket 连接成功');
        }
        startHeartbeatTimer();
    };

    webSocket.onmessage = (event) => {
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

            // 其他所有消息，原样输出到聊天框
            if (typeof addSystemMessage === 'function') {
                addSystemMessage(`📨 ${event.data}`);
            }

        } catch (e) {
            // 不是 JSON 格式，直接输出
            console.log('消息不是 JSON 格式，直接输出');
            if (typeof addSystemMessage === 'function') {
                addSystemMessage(`📨 ${event.data}`);
            }
        }
    };

    webSocket.onclose = (event) => {
        console.log('WebSocket 已断开, code:', event.code, 'reason:', event.reason);
        if (typeof addSystemMessage === 'function') {
            addSystemMessage('⚠️ WebSocket 连接已断开，3秒后重连...');
        }
        stopHeartbeatTimer();

        setTimeout(() => {
            if (sessionStorage.getItem('access_token')) {
                console.log('尝试重新连接 WebSocket...');
                connectWebSocket();
            }
        }, 3000);
    };

    webSocket.onerror = (error) => {
        console.error('WebSocket 错误:', error);
        if (typeof addSystemMessage === 'function') {
            addSystemMessage('❌ WebSocket 连接错误');
        }
    };
}

function startHeartbeatTimer() {
    if (heartbeatTimer) {
        clearInterval(heartbeatTimer);
    }

    heartbeatTimer = setInterval(() => {
        if (webSocket && webSocket.readyState === WebSocket.OPEN) {
            webSocket.send(JSON.stringify({ type: 'ping' }));
            console.log('发送心跳 ping');
        }
    }, 30000);
}

function stopHeartbeatTimer() {
    if (heartbeatTimer) {
        clearInterval(heartbeatTimer);
        heartbeatTimer = null;
    }
}

// 页面关闭前清理
window.addEventListener('beforeunload', () => {
    if (webSocket && webSocket.readyState === WebSocket.OPEN) {
        webSocket.close();
    }
    stopHeartbeatTimer();
});