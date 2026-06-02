// ============================================
// WebSocket 连接和消息处理
// ============================================

function connectWebSocket() {
    const token = sessionStorage.getItem('access_token');
    if (!token) {
        console.log('没有 token，跳过 WebSocket 连接');
        return;
    }

    ws = new WebSocket(`${WS_URL}?token=${encodeURIComponent(token)}`);

    ws.onopen = () => {
        console.log('WebSocket 已连接');
        addSystemMessage('✅ WebSocket 已连接');
    };

    ws.onmessage = (event) => {
        if (!event.data || event.data.trim() === '') return;

        // 直接输出原始消息，不做任何处理
        addSystemMessage(`📨 ${event.data}`);
    };

    ws.onclose = () => {
        console.log('WebSocket 已断开');
        addSystemMessage('⚠️ WebSocket 已断开，3秒后重连...');
        setTimeout(() => {
            if (sessionStorage.getItem('access_token')) connectWebSocket();
        }, 3000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket 错误:', error);
        addSystemMessage('❌ WebSocket 错误');
    };
}