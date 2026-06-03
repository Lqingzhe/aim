// ============================================
// 配置和全局变量
// ============================================
const API_BASE = window.location.origin;
const WS_URL = `ws://${window.location.host}/ws`;

let userID = sessionStorage.getItem('user_id');
let deviceID = sessionStorage.getItem('device_id');
let accessToken = sessionStorage.getItem('access_token');
let refreshToken = sessionStorage.getItem('refresh_token');

// 注意：这里不声明任何与 WebSocket 相关的变量

let currentSession = null;
let sessions = [];
let isRefreshing = false;
let failedQueue = [];

// 分页相关
let isLoadingMore = false;
let hasMoreMessages = true;
let currentGroupId = null;