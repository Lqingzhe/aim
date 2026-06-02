// ============================================
// 会话列表管理
// ============================================
let sessionHasNewMessage = {};

async function loadSessions() {
    const result = await apiCall('POST', '/group/get-group-and-session-id', {});

    if (result && result.code === 0 && result.data) {
        const list = document.getElementById('session-list');
        list.innerHTML = '';
        sessions = [];

        // 加载群聊
        if (result.data.group_info?.group_id_list && result.data.group_info.group_id_list.length > 0) {
            for (const id of result.data.group_info.group_id_list) {
                const groupInfo = await apiCall('POST', '/group/get-group-info', { group_id: String(id) });
                if (groupInfo && groupInfo.code === 0) {
                    addSessionToList({
                        id: String(id),
                        type: 'group',
                        name: groupInfo.data?.group_info?.group_name || `群聊 ${id.substring(0, 8)}`,
                        preview: ''
                    });
                }
            }
        }

        // 加载私聊会话
        if (result.data.session_info?.session_id_list && result.data.session_info.session_id_list.length > 0) {
            const sessionIds = result.data.session_info.session_id_list;
            for (let i = 0; i < sessionIds.length; i++) {
                const sessionId = sessionIds[i];
                addSessionToList({
                    id: String(sessionId),
                    type: 'session',
                    name: `私聊会话`,
                    preview: ''
                });
            }
        }

        if (sessions.length === 0) {
            list.innerHTML = '<div class="empty-tip">📭 暂无会话</div>';
        }
    }
}

function markSessionHasNewMessage(sessionId) {
    sessionHasNewMessage[String(sessionId)] = true;
    const sessionItem = document.querySelector(`.session-item[data-id="${String(sessionId)}"]`);
    if (sessionItem) {
        let badge = sessionItem.querySelector('.new-message-badge');
        if (!badge) {
            badge = document.createElement('span');
            badge.className = 'new-message-badge';
            badge.textContent = '●';
            badge.style.color = '#ef4444';
            badge.style.fontSize = '14px';
            badge.style.marginLeft = '8px';
            sessionItem.querySelector('.session-name').appendChild(badge);
        }
    }

    // 更新页面标题显示（可选）
    updatePageTitle(true);
}

function clearSessionNewMessageMark(sessionId) {
    sessionHasNewMessage[String(sessionId)] = false;
    const sessionItem = document.querySelector(`.session-item[data-id="${String(sessionId)}"]`);
    if (sessionItem) {
        const badge = sessionItem.querySelector('.new-message-badge');
        if (badge) badge.remove();
    }

    // 检查是否还有任何未读消息
    let hasAnyUnread = false;
    for (const key in sessionHasNewMessage) {
        if (sessionHasNewMessage[key]) {
            hasAnyUnread = true;
            break;
        }
    }
    updatePageTitle(hasAnyUnread);
}

// 更新页面标题提示
let originalTitle = document.title;
function updatePageTitle(hasUnread) {
    if (hasUnread) {
        document.title = '🔴 有新消息 - ' + originalTitle;
    } else {
        document.title = originalTitle;
    }
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

    const icon = session.type === 'group' ? '👥' : '💬';
    const typeName = session.type === 'group' ? '群聊' : '私聊';
    const shortId = String(session.id).length > 12 ? String(session.id).substring(0, 12) + '...' : String(session.id);

    // 私聊会话添加删除按钮
    let deleteBtnHtml = '';
    if (session.type === 'session') {
        deleteBtnHtml = `<button class="delete-session-btn" data-id="${session.id}" title="删除好友" style="background:none;border:none;color:#ef4444;cursor:pointer;margin-left:8px;">🗑️</button>`;
    }

    div.innerHTML = `
        <div style="display:flex;justify-content:space-between;align-items:center;">
            <div class="session-name" style="flex:1;">
                ${icon} ${escapeHtml(session.name)}
            </div>
            ${deleteBtnHtml}
        </div>
        <div class="session-id">${typeName} ID: ${shortId}</div>
    `;

    // 删除好友按钮事件
    if (session.type === 'session') {
        const deleteBtn = div.querySelector('.delete-session-btn');
        deleteBtn.onclick = async (e) => {
            e.stopPropagation();
            if (confirm(`确定要删除好友 ${session.name} 吗？`)) {
                const result = await apiCall('POST', '/group/delete-session', { session_id: session.id });
                if (result && result.code === 0) {
                    alert('删除成功');
                    loadSessions();
                    if (currentSession && currentSession.id === session.id) {
                        document.getElementById('messages-container').innerHTML = '';
                        document.getElementById('chat-title').innerHTML = '请选择一个会话';
                        currentSession = null;
                    }
                } else {
                    alert('删除失败: ' + (result?.message || '未知错误'));
                }
            }
        };
    }

    div.onclick = () => selectSession(session);
    list.appendChild(div);
    sessions.push(session);
}

async function selectSession(session) {
    currentSession = session;

    clearSessionNewMessageMark(session.id);

    // 更新高亮
    document.querySelectorAll('.session-item').forEach(item => {
        item.classList.remove('active');
        if (item.getAttribute('data-id') === String(session.id)) {
            item.classList.add('active');
        }
    });

    const icon = session.type === 'group' ? '👥' : '💬';
    const typeName = session.type === 'group' ? '群聊' : '私聊';

    document.getElementById('chat-title').innerHTML = `${icon} ${escapeHtml(session.name)}`;
    document.getElementById('chat-id').innerHTML = `${typeName} ID: ${session.id}`;

    // 清空消息区域
    document.getElementById('messages-container').innerHTML = '';
    addSystemMessage('加载消息中...');

    // 重置分页
    currentGroupId = null;
    hasMoreMessages = true;
    isLoadingMore = false;

    // 加载历史消息
    await loadHistoryMessages(session.id, false);
}
