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
    error => Promise.reject(error)
);

// 响应拦截器
api.interceptors.response.use(
    response => response,
    async error => {
        const originalRequest = error.config;

        if (error.response?.data?.code === 1101 && !originalRequest._retry) {
            originalRequest._retry = true;
            const refreshed = await refreshAccessToken();
            if (refreshed) {
                return api(originalRequest);
            } else {
                sessionStorage.clear();
                if (ws) ws.close();
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
        console.error(`API Exception: ${path}`, error);
        return { code: -1, message: error.message };
    }
}