// ============================================
// 会话列表管理
// ============================================

let sessionHasNewMessage = {};
window.currentRemarkMap = {};

async function loadSessions() {
    console.log('loadSessions 开始...');

    const list = document.getElementById('session-list');
    list.innerHTML = '<div class="empty-tip">加载会话中...</div>';

    const myInfo = await apiCall('POST', '/user/get-user-info', {});
    console.log('get-user-info 返回:', myInfo);

    window.currentRemarkMap = {};
    if (myInfo && myInfo.code === 0 && myInfo.data?.user_info) {
        let userInfoData = myInfo.data.user_info;
        let remarkInfos = userInfoData.RemarkInfos || userInfoData.remark_info || [];
        for (const r of remarkInfos) {
            const goalId = r.GoalUserID || r.goal_user_id;
            const nickName = r.NickName || r.nick_name;
            if (goalId && nickName) {
                window.currentRemarkMap[String(goalId)] = nickName;
            }
        }
    }

    const result = await apiCall('POST', '/group/get-group-and-session-id', {});
    if (!result || result.code !== 0) {
        list.innerHTML = '<div class="empty-tip">❌ 加载失败，请刷新重试</div>';
        return;
    }

    list.innerHTML = '';
    sessions = [];

    const groupIdList = result.data?.group_info?.group_id_list || [];
    for (const id of groupIdList) {
        const groupInfo = await apiCall('POST', '/group/get-group-info', { group_id: String(id) });
        let groupName = `群聊 ${id}`;
        if (groupInfo && groupInfo.code === 0 && groupInfo.data?.group_info) {
            groupName = groupInfo.data.group_info.group_name || `群聊 ${id}`;
        }

        // 获取当前用户在该群中的权限
        const groupUserInfo = await apiCall('POST', '/group/get-group-info-with-user', { group_id: String(id) });
        let userRole = 'Member';
        if (groupUserInfo && groupUserInfo.code === 0 && groupUserInfo.data?.group_info) {
            userRole = groupUserInfo.data.group_info.group_role || 'Member';
        }

        addSessionToList({
            id: String(id),
            type: 'group',
            name: groupName,
            goalUserId: null,
            userRole: userRole
        });
    }

    const sessionIdList = result.data?.session_info?.session_id_list || [];
    const userOfSessionIdList = result.data?.session_info?.user_of_session_id_list || [];

    for (let i = 0; i < sessionIdList.length; i++) {
        const sessionId = sessionIdList[i];
        const goalUserId = userOfSessionIdList[i];
        const goalUserIdStr = String(goalUserId);

        let displayName = `私聊 ${sessionId}`;
        let hasRemark = false;

        if (window.currentRemarkMap[goalUserIdStr]) {
            displayName = window.currentRemarkMap[goalUserIdStr];
            hasRemark = true;
        } else if (goalUserId) {
            const userInfo = await apiCall('POST', '/user/get-other-user-info', { goal_user_id: goalUserIdStr });
            if (userInfo && userInfo.code === 0 && userInfo.data?.user_info) {
                let userData = userInfo.data.user_info;
                let userName = userData.UserName || userData.user_name;
                if (userName) {
                    displayName = userName;
                }
            }
        }

        addSessionToList({
            id: String(sessionId),
            type: 'session',
            name: displayName,
            goalUserId: goalUserIdStr,
            hasRemark: hasRemark
        });
    }

    if (sessions.length === 0) {
        list.innerHTML = '<div class="empty-tip">📭 暂无会话</div>';
    }
}

function addSessionToList(session) {
    const list = document.getElementById('session-list');
    const existing = document.querySelector(`.session-item[data-id="${String(session.id)}"]`);
    if (existing) return;

    const div = document.createElement('div');
    div.className = 'session-item';
    div.setAttribute('data-id', String(session.id));
    div.setAttribute('data-type', session.type);
    if (session.goalUserId) {
        div.setAttribute('data-goal-user-id', session.goalUserId);
    }
    if (session.userRole) {
        div.setAttribute('data-user-role', session.userRole);
    }

    const icon = session.type === 'group' ? '👥' : '💬';
    const typeName = session.type === 'group' ? '群聊' : '私聊';
    const shortId = String(session.id).length > 12 ? String(session.id).substring(0, 12) + '...' : String(session.id);

    let remarkBadge = '';
    if (session.hasRemark) {
        remarkBadge = '<span style="background:#eef2ff;color:#667eea;font-size:10px;padding:2px 6px;border-radius:10px;margin-left:6px;">备注</span>';
    }

    let roleBadge = '';
    if (session.type === 'group' && session.userRole) {
        const roleColors = {
            'Owner': '#f59e0b',
            'Manager': '#10b981',
            'Member': '#6b7280'
        };
        roleBadge = `<span style="background:${roleColors[session.userRole] || '#6b7280'};color:white;font-size:10px;padding:2px 6px;border-radius:10px;margin-left:6px;">${session.userRole}</span>`;
    }

    let deleteBtnHtml = '';
    if (session.type === 'session') {
        deleteBtnHtml = `<button class="delete-session-btn" data-id="${session.id}" title="删除会话" style="background:none;border:none;color:#ef4444;cursor:pointer;margin-left:8px;">🗑️</button>`;
    } else if (session.type === 'group' && session.userRole === 'Owner') {
        deleteBtnHtml = `<button class="delete-group-btn" data-id="${session.id}" title="解散群聊" style="background:none;border:none;color:#ef4444;cursor:pointer;margin-left:8px;">💣</button>`;
    }

    div.innerHTML = `
        <div style="display:flex;justify-content:space-between;align-items:center;">
            <div class="session-name" style="flex:1;">
                ${icon} ${escapeHtml(session.name)} ${remarkBadge} ${roleBadge}
            </div>
            ${deleteBtnHtml}
        </div>
        <div class="session-id">${typeName} ID: ${shortId}</div>
        ${session.goalUserId ? `<div class="session-extra" style="font-size: 10px; color: #9ca3af; margin-top: 2px;">对方ID: ${session.goalUserId}</div>` : ''}
    `;

    if (session.type === 'session') {
        const deleteBtn = div.querySelector('.delete-session-btn');
        deleteBtn.onclick = async (e) => {
            e.stopPropagation();
            if (confirm(`确定要删除这个私聊会话吗？`)) {
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
    } else if (session.type === 'group' && session.userRole === 'Owner') {
        const deleteGroupBtn = div.querySelector('.delete-group-btn');
        deleteGroupBtn.onclick = async (e) => {
            e.stopPropagation();
            if (confirm(`确定要解散群聊 ${session.name} 吗？此操作不可恢复！`)) {
                const result = await apiCall('POST', '/group/delete-group', { group_id: session.id });
                if (result && result.code === 0) {
                    alert('解散群聊成功');
                    loadSessions();
                    if (currentSession && currentSession.id === session.id) {
                        document.getElementById('messages-container').innerHTML = '';
                        document.getElementById('chat-title').innerHTML = '请选择一个会话';
                        currentSession = null;
                    }
                } else {
                    alert('解散失败: ' + (result?.message || '未知错误'));
                }
            }
        };
    }

    div.onclick = () => selectSession(session);
    list.appendChild(div);
    sessions.push(session);
}

// 获取当前群聊中用户的权限
async function getGroupUserRole(groupId, userId) {
    const result = await apiCall('POST', '/group/get-group-info-with-user', { group_id: String(groupId) });
    if (result && result.code === 0 && result.data?.group_info) {
        return result.data.group_info.group_role || 'Member';
    }
    return 'Member';
}

async function selectSession(session) {
    currentSession = session;
    clearSessionNewMessageMark(session.id);

    document.querySelectorAll('.session-item').forEach(item => {
        item.classList.remove('active');
        if (item.getAttribute('data-id') === String(session.id)) {
            item.classList.add('active');
        }
    });

    const icon = session.type === 'group' ? '👥' : '💬';
    const typeName = session.type === 'group' ? '群聊' : '私聊';

    let roleDisplay = '';
    if (session.type === 'group' && session.userRole) {
        const roleNames = { 'Owner': '群主', 'Manager': '管理员', 'Member': '成员' };
        roleDisplay = `<span style="font-size:12px;color:#667eea;margin-left:8px;">(${roleNames[session.userRole] || session.userRole})</span>`;
    }

    document.getElementById('chat-title').innerHTML = `${icon} ${escapeHtml(session.name)} ${roleDisplay}`;
    document.getElementById('chat-id').innerHTML = `${typeName} ID: ${session.id}`;

    document.getElementById('messages-container').innerHTML = '';
    addSystemMessage('加载消息中...');

    currentGroupId = null;
    hasMoreMessages = true;
    isLoadingMore = false;

    await loadHistoryMessages(session.id, false);
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
    updatePageTitle(true);
}

function clearSessionNewMessageMark(sessionId) {
    sessionHasNewMessage[String(sessionId)] = false;
    const sessionItem = document.querySelector(`.session-item[data-id="${String(sessionId)}"]`);
    if (sessionItem) {
        const badge = sessionItem.querySelector('.new-message-badge');
        if (badge) badge.remove();
    }
    let hasAnyUnread = false;
    for (const key in sessionHasNewMessage) {
        if (sessionHasNewMessage[key]) {
            hasAnyUnread = true;
            break;
        }
    }
    updatePageTitle(hasAnyUnread);
}

let originalTitle = document.title;
function updatePageTitle(hasUnread) {
    if (hasUnread) {
        document.title = '🔴 有新消息 - ' + originalTitle;
    } else {
        document.title = originalTitle;
    }
}