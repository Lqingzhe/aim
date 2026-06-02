// ============================================
// 认证相关
// ============================================
function processQueue(err, token = null) {
    failedQueue.forEach(promise => {
        if (err) {
            promise.reject(err);
        } else {
            promise.resolve(token);
        }
    });
    failedQueue = [];
}

async function refreshAccessToken() {
    if (!refreshToken) return false;

    if (isRefreshing) {
        return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject });
        }).then(token => !!token).catch(() => false);
    }

    isRefreshing = true;

    try {
        const response = await fetch(`${API_BASE}/user/refresh-token`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshToken })
        });
        const data = await response.json();

        if (data.code === 0 && data.data?.token_info) {
            accessToken = data.data.token_info.access_token;
            refreshToken = data.data.token_info.refresh_token;
            sessionStorage.setItem('access_token', accessToken);
            sessionStorage.setItem('refresh_token', refreshToken);
            processQueue(null, accessToken);
            return true;
        }
        processQueue(new Error('Refresh failed'), null);
        return false;
    } catch (error) {
        console.error('Refresh token failed:', error);
        processQueue(error, null);
        return false;
    } finally {
        isRefreshing = false;
    }
}

async function logout() {
    await apiCall('POST', '/user/logout-a-device', {});
    sessionStorage.clear();
    if (ws) ws.close();
    window.location.href = '/login';
}