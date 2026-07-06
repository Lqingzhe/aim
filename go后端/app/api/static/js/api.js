// ============================================
// axios 配置和 API 调用
// ============================================

const api = axios.create({
    baseURL: API_BASE,
    timeout: 30000
});

// 请求拦截器
api.interceptors.request.use(
    config => {
        const token = sessionStorage.getItem('access_token');
        const device = sessionStorage.getItem('device_id');

        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        if (device) {
            config.headers['X-Device-ID'] = device;
        }
        config.headers['Content-Type'] = 'application/json';

        return config;
    },
    error => {
        console.error('请求拦截器错误:', error);
        return Promise.reject(error);
    }
);

// 响应拦截器
api.interceptors.response.use(
    response => {
        // 检查业务状态码是否为 token 相关错误
        const responseData = response.data;
        // 1101: CodeAccessTokenInvalid (无效)
        // 1102: CodeAccessTokenExpired (过期)
        // 这两个都需要刷新 token
        if (responseData && (responseData.code === 1101 || responseData.code === 1102)) {
            console.log('检测到 token 问题 (code:', responseData.code, ')，准备刷新');
            return Promise.reject({
                response: response,
                config: response.config,
                isTokenExpired: true,
                errorCode: responseData.code
            });
        }
        return response;
    },
    async error => {
        const originalRequest = error.config;

        // 检查是否是 token 过期错误（支持多种错误码）
        const isTokenExpired = error.isTokenExpired ||
            error.response?.data?.code === 1101 ||
            error.response?.data?.code === 1102 ||
            error.response?.status === 401;

        console.log('响应拦截器捕获错误:', {
            url: originalRequest?.url,
            status: error.response?.status,
            code: error.response?.data?.code,
            isTokenExpired: isTokenExpired,
            retry: originalRequest?._retry
        });

        // 如果是 token 相关问题且没有重试过
        if (isTokenExpired && !originalRequest._retry) {
            originalRequest._retry = true;
            console.log('Token 问题，尝试刷新...');

            const refreshed = await refreshAccessToken();

            if (refreshed) {
                console.log('Token 刷新成功，重试原请求');
                const newToken = sessionStorage.getItem('access_token');
                originalRequest.headers['Authorization'] = `Bearer ${newToken}`;
                return api(originalRequest);
            } else {
                console.log('Token 刷新失败，跳转到登录页');
                sessionStorage.clear();
                if (window.webSocket && window.webSocket.readyState === WebSocket.OPEN) {
                    window.webSocket.close();
                }
                window.location.href = '/login';
                return Promise.reject(error);
            }
        }

        return Promise.reject(error);
    }
);

async function apiCall(method, path, data = null) {
    try {
        const config = { method, url: path };
        if (method === 'GET' && data) {
            config.params = data;
        } else if (data) {
            const stringData = {};
            const idFields = ['group_id', 'user_id', 'goal_user_id', 'session_id', 'message_id', 'file_id'];
            for (const [key, value] of Object.entries(data)) {
                if (idFields.includes(key) && value !== undefined && value !== null) {
                    stringData[key] = String(value);
                } else {
                    stringData[key] = value;
                }
            }
            config.data = stringData;
        }
        const response = await api(config);
        return response.data;
    } catch (error) {
        // 如果是 token 相关问题，等待刷新后重试
        const errorCode = error.response?.data?.code;
        if (errorCode === 1101 || errorCode === 1102) {
            console.log('apiCall 中检测到 token 问题 (code:', errorCode, ')，等待刷新后重试');
            const refreshed = await refreshAccessToken();
            if (refreshed) {
                return apiCall(method, path, data);
            }
        }
        console.error(`API Exception: ${path}`, error);
        return { code: -1, message: error.message || '请求失败' };
    }
}