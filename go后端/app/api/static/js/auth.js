// ============================================
// 认证相关 - Token 刷新机制
// ============================================

function processQueue(err, token = null) {
    console.log('processQueue 执行, 队列长度:', window.failedQueue?.length || 0);
    if (window.failedQueue) {
        window.failedQueue.forEach(promise => {
            if (err) {
                promise.reject(err);
            } else {
                promise.resolve(token);
            }
        });
        window.failedQueue = [];
    }
}

async function refreshAccessToken() {
    const currentRefreshToken = sessionStorage.getItem('refresh_token');

    console.log('refreshAccessToken 被调用, refreshToken:', currentRefreshToken ? '存在' : '不存在');

    if (!currentRefreshToken) {
        console.log('没有 refresh_token，无法刷新');
        return false;
    }

    if (window.isRefreshing) {
        console.log('正在刷新中，加入等待队列');
        return new Promise((resolve, reject) => {
            if (!window.failedQueue) window.failedQueue = [];
            window.failedQueue.push({ resolve, reject });
        }).then(token => {
            console.log('等待队列完成，获得新token');
            return !!token;
        }).catch(() => false);
    }

    window.isRefreshing = true;
    console.log('开始刷新 access_token...');

    try {
        const deviceID = sessionStorage.getItem('device_id') || 'web-001';

        console.log('发送刷新请求，refresh_token:', currentRefreshToken.substring(0, 20) + '...');

        const response = await fetch(`${API_BASE}/user/refresh-token`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Device-ID': deviceID
            },
            body: JSON.stringify({ refresh_token: currentRefreshToken })
        });

        console.log('刷新接口 HTTP 状态:', response.status);

        const data = await response.json();
        console.log('刷新响应完整数据:', JSON.stringify(data, null, 2));

        if (data.code === 0 && data.data?.token_info) {
            const newAccessToken = data.data.token_info.access_token;
            const newRefreshToken = data.data.token_info.refresh_token;

            console.log('Token 刷新成功，新 access_token:', newAccessToken.substring(0, 30) + '...');

            sessionStorage.setItem('access_token', newAccessToken);
            sessionStorage.setItem('refresh_token', newRefreshToken);

            window.accessToken = newAccessToken;
            window.refreshToken = newRefreshToken;

            processQueue(null, newAccessToken);
            return true;
        }

        console.log('Token 刷新失败, code:', data.code, 'message:', data.message);
        processQueue(new Error('Refresh failed'), null);
        return false;
    } catch (error) {
        console.error('Refresh token 网络错误:', error);
        processQueue(error, null);
        return false;
    } finally {
        window.isRefreshing = false;
    }
}

// 退出单个设备（当前设备）
async function logoutDevice() {
    console.log('执行退出本设备...');
    try {
        const result = await apiCall('POST', '/user/logout-a-device', {});
        console.log('退出本设备结果:', result);
    } catch (e) {
        console.log('退出本设备接口调用失败:', e);
    }
    sessionStorage.clear();
    if (window.webSocket && window.webSocket.readyState === WebSocket.OPEN) {
        window.webSocket.close();
    }
    window.location.href = '/login';
}

// 退出所有设备
async function logoutAll() {
    console.log('执行退出所有设备...');
    try {
        const result = await apiCall('POST', '/user/logout-all-device', {});
        console.log('退出所有设备结果:', result);
    } catch (e) {
        console.log('退出所有设备接口调用失败:', e);
    }
    sessionStorage.clear();
    if (window.webSocket && window.webSocket.readyState === WebSocket.OPEN) {
        window.webSocket.close();
    }
    window.location.href = '/login';
}