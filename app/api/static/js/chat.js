// ============================================
// 配置
// ============================================
const API_BASE = window.location.origin;
const WS_URL = `ws://${window.location.host}/ws`;

let userID = localStorage.getItem('user_id');
let deviceID = localStorage.getItem('device_id');
let accessToken = localStorage.getItem('access_token');
let refreshToken = localStorage.getItem('refresh_token');

let ws = null;
let currentSession = null;
let sessions = [];

// 创建 axios 实例
const api = axios.create({
    baseURL: API_BASE,
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
});

// 请求拦截器：自动添加 Token 和设备ID
api.interceptors.request.use(
    config => {
        if (accessToken) {
            config.headers['Authorization'] = `Bearer ${accessToken}`;
        }
        if (deviceID) {
            config.headers['X-Device-ID'] = deviceID;
        }
        return config;
    },
    error => Promise.reject(error)
);

// 响应拦截器：处理 Token 过期
api.interceptors.response.use(
    response => response,
    async error => {
        const originalRequest = error.config;

        // Token 过期（根据你的错误码调整）
        if (error.response?.data?.code === 1101 && !originalRequest._retry) {
            originalRequest._retry = true;
            const refreshed = await refreshAccessToken();
            if (refreshed) {
                return api(originalRequest);
            } else {
                localStorage.clear();
                window.location.href = '/login';
                return Promise.reject(error);
            }
        }
        return Promise.reject(error);
    }
);

// 检查登录
if (!accessToken || !userID) {
    window.location.href = '/login';
}

// 显示用户信息
document.getElementById('user-name').innerHTML = `用户 ${userID.substring(0, 10)}...`;
document.getElementById('user-id').textContent = `ID: ${userID}`;

// ============================================
// 刷新 Token
// ============================================
async function refreshAccessToken() {
    if (!refreshToken) return false;
    try {
        const response = await api.post('/user/refresh-token', {
            refresh_token: refreshToken
        });
        const data = response.data;
        if (data.code === 0 && data.data?.token_info) {
            accessToken = data.data.token_info.access_token;
            refreshToken = data.data.token_info.refresh_token;
            localStorage.setItem('access_token', accessToken);
            localStorage.setItem('refresh_token', refreshToken);
            return true;
        }
        return false;
    } catch (error) {
        console.error('Refresh token failed:', error);
        return false;
    }
}

// ============================================
// API 调用封装
// ============================================
async function apiCall(method, path, data = null) {
    try {
        const config = { method, url: path };
        if (method === 'GET') {
            config.params = data;
        } else {
            // 将 ID 字段转为字符串
            if (data) {
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
        }
        const response = await api(config);
        const result = response.data;
        if (result.code !== 0 && result.code !== 1301) {
            console.error(`API Error: ${path}`, result);
        }
        return result;
    } catch (error) {
        console.error(`API Exception: ${path}`, error);
        return { code: -1, message: error.message };
    }
}

// ============================================
// WebSocket 连接
// ============================================
function connectWebSocket() {
    ws = new WebSocket(`${WS_URL}?token=${encodeURIComponent(accessToken)}`);
    ws.onopen = () => {
        console.log('WebSocket 已连接');
        addSystemMessage('✅ WebSocket 已连接');
    };
    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.Code === 0 && data.Data) {
                try {
                    const msgData = JSON.parse(data.Data);
                    addMessageToChat({
                        user_id: String(msgData.user_id),
                        message_content: msgData.message_content || msgData.Content,
                        send_time_second: msgData.send_time_second
                    });
                } catch (e) {
                    addSystemMessage(`📨 收到: ${data.Data.substring(0, 100)}`);
                }
            }
        } catch (e) {
            console.log('WebSocket 消息:', event.data);
        }
    };
    ws.onclose = () => {
        console.log('WebSocket 已断开');
        addSystemMessage('⚠️ WebSocket 已断开，3秒后重连...');
        setTimeout(() => {
            if (accessToken) connectWebSocket();
        }, 3000);
    };
}

// ============================================
// UI 辅助函数
// ============================================
function addSystemMessage(text) {
    const container = document.getElementById('messages-container');
    const div = document.createElement('div');
    div.style.textAlign = 'center';
    div.style.color = '#9ca3af';
    div.style.fontSize = '12px';
    div.style.padding = '8px';
    div.textContent = text;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function addMessageToChat(message) {
    const container = document.getElementById('messages-container');
    const isSelf = String(message.user_id) === String(userID);
    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${isSelf ? 'self' : 'other'}`;
    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    bubble.textContent = message.message_content || '';
    const meta = document.createElement('div');
    meta.className = 'message-meta';
    const time = message.send_time_second
        ? new Date(message.send_time_second * 1000).toLocaleTimeString()
        : new Date().toLocaleTimeString();
    meta.textContent = `用户 ${message.user_id} ${time}`;
    msgDiv.appendChild(bubble);
    msgDiv.appendChild(meta);
    container.appendChild(msgDiv);
    container.scrollTop = container.scrollHeight;
}

function addSessionToList(session) {
    const list = document.getElementById('session-list');
    const existing = document.querySelector(`.session-item[data-id="${String(session.id)}"]`);
    if (existing) return;

    const emptyTip = list.querySelector('.empty-tip');
    if (emptyTip) list.innerHTML = '';

    const div = document.createElement('div');
    div.className = 'session-item';
    div.setAttribute('data-id', String(session.id));
    div.setAttribute('data-type', session.type);
    div.setAttribute('data-name', session.name);
    div.innerHTML = `<div class="session-name">${escapeHtml(session.name)}</div><div class="session-preview">${escapeHtml(session.preview || '点击开始聊天')}</div>`;
    div.onclick = () => selectSession(session);
    list.appendChild(div);

    sessions.push(session);
    if (!currentSession) selectSession(session);
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function selectSession(session) {
    currentSession = session;
    document.querySelectorAll('.session-item').forEach(item => {
        item.classList.remove('active');
        if (item.getAttribute('data-id') === String(session.id)) item.classList.add('active');
    });
    document.getElementById('chat-title').textContent = session.name;
    document.getElementById('messages-container').innerHTML = '';
    await loadMessages(session.id);
}

// ============================================
// 加载会话列表
// ============================================
async function loadSessions() {
    const result = await apiCall('POST', '/group/get-group-and-session-id', {});
    if (result && result.code === 0 && result.data) {
        const list = document.getElementById('session-list');
        list.innerHTML = '';
        sessions = [];

        if (result.data.group_info?.group_id_list) {
            for (const id of result.data.group_info.group_id_list) {
                const groupInfo = await apiCall('POST', '/group/get-group-info', { group_id: String(id) });
                if (groupInfo && groupInfo.code === 0) {
                    addSessionToList({
                        id: String(id),
                        type: 'group',
                        name: groupInfo.data?.group_info?.group_name || `群组 ${id.substring(0, 8)}`,
                        preview: ''
                    });
                }
            }
        }

        if (result.data.session_info?.session_id_list) {
            for (const id of result.data.session_info.session_id_list) {
                addSessionToList({
                    id: String(id),
                    type: 'session',
                    name: `会话 ${id.substring(0, 8)}`,
                    preview: ''
                });
            }
        }

        if (sessions.length === 0) {
            list.innerHTML = '<div class="empty-tip">暂无会话，请添加好友或创建群聊</div>';
        }
    }
}

// ============================================
// 加载消息历史
// ============================================
async function loadMessages(groupId) {
    const result = await apiCall('POST', '/message/get-message-list', {
        group_id: String(groupId),
        start_time_second: 0,
        end_time_second: 0
    });

    if (result && result.code === 0 && result.data?.message_info?.message_list) {
        const messages = result.data.message_info.message_list;
        for (let i = messages.length - 1; i >= 0; i--) {
            const msg = messages[i];
            addMessageToChat({
                user_id: String(msg.user_id),
                message_content: msg.message_content,
                send_time_second: msg.send_time_second
            });
        }
    }
}

// ============================================
// 发送消息
// ============================================
async function sendMessage() {
    const input = document.getElementById('message-input');
    const content = input.value.trim();
    if (!content || !currentSession) return;

    const result = await apiCall('POST', '/message/send-message', {
        group_id: String(currentSession.id),
        message_content: content
    });

    if (result && result.code === 0) {
        input.value = '';
        addMessageToChat({
            user_id: userID,
            message_content: content,
            send_time_second: Math.floor(Date.now() / 1000)
        });
    } else {
        addSystemMessage('发送失败: ' + (result?.message || '未知错误'));
    }
}

// ============================================
// 好友申请
// ============================================
function showAddFriendModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <input type="text" id="friend-user-id" placeholder="对方用户ID" style="width:100%;padding:10px;margin-bottom:15px;border:1px solid #ddd;border-radius:6px;">
        <button id="submit-add-friend" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">发送好友申请</button>
    `;
    document.getElementById('modal-title').textContent = '添加好友';
    modal.style.display = 'flex';

    document.getElementById('submit-add-friend').onclick = async () => {
        const goalUserId = document.getElementById('friend-user-id').value.trim();
        if (!goalUserId) {
            alert('请输入对方用户ID');
            return;
        }

        const result = await apiCall('POST', '/group/apply-for-friend', { goal_user_id: goalUserId });
        if (result && result.code === 0) {
            alert('好友申请已发送');
            modal.style.display = 'none';
        } else {
            alert('发送失败: ' + (result?.message || '未知错误'));
        }
    };
}

// ============================================
// 创建群聊
// ============================================
function showCreateGroupModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <input type="text" id="group-name" placeholder="群名称" style="width:100%;padding:10px;margin-bottom:15px;border:1px solid #ddd;border-radius:6px;">
        <button id="submit-create-group" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">创建群聊</button>
    `;
    document.getElementById('modal-title').textContent = '创建群聊';
    modal.style.display = 'flex';

    document.getElementById('submit-create-group').onclick = async () => {
        const groupName = document.getElementById('group-name').value.trim();
        if (!groupName) {
            alert('请输入群名称');
            return;
        }

        const result = await apiCall('POST', '/group/create-group', { group_name: groupName });
        if (result && result.code === 0 && result.data?.group_info?.group_id) {
            alert(`群聊创建成功！群ID: ${result.data.group_info.group_id}`);
            modal.style.display = 'none';
            loadSessions();
        } else {
            alert('创建失败: ' + (result?.message || '未知错误'));
        }
    };
}

// ============================================
// 加入群聊
// ============================================
function showJoinGroupModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <input type="text" id="join-group-id" placeholder="群ID" style="width:100%;padding:10px;margin-bottom:15px;border:1px solid #ddd;border-radius:6px;">
        <button id="submit-join-group" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">申请加入</button>
    `;
    document.getElementById('modal-title').textContent = '加入群聊';
    modal.style.display = 'flex';

    document.getElementById('submit-join-group').onclick = async () => {
        const groupId = document.getElementById('join-group-id').value.trim();
        if (!groupId) {
            alert('请输入群ID');
            return;
        }

        const result = await apiCall('POST', '/group/set-group-apply', { group_id: groupId });
        if (result && result.code === 0) {
            alert('入群申请已发送，等待群主/管理员审核');
            modal.style.display = 'none';
        } else {
            alert('申请失败: ' + (result?.message || '未知错误'));
        }
    };
}

// ============================================
// 获取好友申请列表
// ============================================
async function showFriendApplyList() {
    const result = await apiCall('POST', '/group/get-friend-apply-list', {});
    if (result && result.code === 0 && result.data?.session_info?.apply_user_list) {
        const applyList = result.data.session_info.apply_user_list;
        if (applyList.length === 0) {
            alert('暂无好友申请');
            return;
        }

        const modal = document.getElementById('modal');
        const modalBody = document.getElementById('modal-body');
        let html = '<div style="max-height:300px;overflow-y:auto;">';
        for (const userId of applyList) {
            html += `<div class="list-item" style="display:flex;justify-content:space-between;align-items:center;padding:10px;border-bottom:1px solid #eee;">
                        <span>用户 ${userId}</span>
                        <button class="agree-apply" data-user-id="${userId}" style="padding:4px 12px;background:#667eea;color:white;border:none;border-radius:4px;cursor:pointer;">同意</button>
                    </div>`;
        }
        html += '</div>';
        modalBody.innerHTML = html;
        document.getElementById('modal-title').textContent = '好友申请列表';
        modal.style.display = 'flex';

        document.querySelectorAll('.agree-apply').forEach(btn => {
            btn.onclick = async () => {
                const goalUserId = btn.dataset.userId;
                const createResult = await apiCall('POST', '/group/creat-session', { goal_user_id: goalUserId });
                if (createResult && createResult.code === 0) {
                    alert(`已同意并创建会话，会话ID: ${createResult.data?.session_info?.session_id}`);
                    modal.style.display = 'none';
                    loadSessions();
                } else {
                    alert('操作失败: ' + (createResult?.message || '未知错误'));
                }
            };
        });
    } else {
        alert('获取申请列表失败: ' + (result?.message || '未知错误'));
    }
}

// ============================================
// 退出登录
// ============================================
async function logout() {
    await apiCall('POST', '/user/logout-a-device', {});
    localStorage.clear();
    if (ws) ws.close();
    window.location.href = '/login';
}

// ============================================
// 初始化
// ============================================
function init() {
    connectWebSocket();
    loadSessions();

    document.getElementById('send-btn').onclick = sendMessage;
    document.getElementById('btn-logout').onclick = logout;
    document.getElementById('btn-add-friend').onclick = showAddFriendModal;
    document.getElementById('btn-create-group').onclick = showCreateGroupModal;
    document.getElementById('btn-join-group').onclick = showJoinGroupModal;
    document.getElementById('btn-friend-apply-list').onclick = showFriendApplyList;
    document.getElementById('btn-refresh').onclick = loadSessions;

    document.getElementById('message-input').addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    document.getElementById('modal-close').onclick = () => {
        document.getElementById('modal').style.display = 'none';
    };
    window.onclick = (e) => {
        const modal = document.getElementById('modal');
        if (e.target === modal) modal.style.display = 'none';
    };
}

init();