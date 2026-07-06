// ============================================
// 配置和全局变量
// ============================================

const API_BASE = window.location.origin;
const WS_URL = `ws://${window.location.host}/ws`;

// 从 sessionStorage 读取认证信息
let userID = sessionStorage.getItem('user_id');
let deviceID = sessionStorage.getItem('device_id');
let accessToken = sessionStorage.getItem('access_token');
let refreshToken = sessionStorage.getItem('refresh_token');

// 挂载到 window 以便全局访问
window.accessToken = accessToken;
window.refreshToken = refreshToken;
window.userID = userID;
window.deviceID = deviceID;

// Token 刷新相关变量（挂载到 window 以便跨文件访问）
window.isRefreshing = false;
window.failedQueue = [];

// WebSocket 相关变量
window.webSocket = null;
window.heartbeatTimer = null;

// 会话相关变量
let currentSession = null;
let sessions = [];

// 分页相关
let isLoadingMore = false;
let hasMoreMessages = true;
let currentGroupId = null;

// 打印初始化状态
console.log('Config 初始化:', {
    hasUserID: !!userID,
    hasAccessToken: !!accessToken,
    hasRefreshToken: !!refreshToken
});