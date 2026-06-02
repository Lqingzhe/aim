// ============================================
// UI 辅助函数
// ============================================
function addSystemMessage(text) {
    const container = document.getElementById('messages-container');
    const div = document.createElement('div');
    div.className = 'system-message';
    div.textContent = text;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function addMessageToChat(message) {
    const container = document.getElementById('messages-container');
    const isSelf = String(message.user_id) === String(userID);
    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${isSelf ? 'self' : 'other'}`;

    msgDiv.dataset.messageId = message.message_id;
    msgDiv.dataset.userId = message.user_id;
    msgDiv.dataset.time = message.send_time_second;

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

function addMessageToChatTop(message, groupId) {
    const container = document.getElementById('messages-container');
    const isSelf = String(message.user_id) === String(userID);
    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${isSelf ? 'self' : 'other'}`;
    msgDiv.dataset.messageId = message.message_id;
    msgDiv.dataset.userId = message.user_id;
    msgDiv.dataset.time = message.send_time_second;

    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    bubble.textContent = message.message_content || '';

    const meta = document.createElement('div');
    meta.className = 'message-meta';
    const time = message.send_time_second
        ? new Date(message.send_time_second * 1000).toLocaleString()
        : new Date().toLocaleString();
    meta.textContent = `用户 ${message.user_id} ${time}`;

    msgDiv.appendChild(bubble);
    msgDiv.appendChild(meta);
    container.insertBefore(msgDiv, container.firstChild);
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
function showUserInfoModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');

    apiCall('POST', '/user/get-user-info', {}).then(result => {
        if (result && result.code === 0 && result.data?.user_info) {
            const userInfo = result.data.user_info;
            modalBody.innerHTML = `
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">用户ID</label>
                    <input type="text" value="${userInfo.user_id}" disabled style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;background:#f5f5f5;">
                </div>
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">用户名</label>
                    <input type="text" id="edit-username" value="${escapeHtml(userInfo.user_name || '')}" placeholder="请输入用户名" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                </div>
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">个性签名</label>
                    <textarea id="edit-introduction" rows="3" placeholder="请输入个性签名" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">${escapeHtml(userInfo.introduction || '')}</textarea>
                </div>
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">生日</label>
                    <div style="display:flex;gap:8px;">
                        <input type="number" id="edit-birthday-year" placeholder="年" value="${userInfo.birthday_year || ''}" style="width:33%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                        <input type="number" id="edit-birthday-month" placeholder="月" value="${userInfo.birthday_month || ''}" style="width:33%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                        <input type="number" id="edit-birthday-day" placeholder="日" value="${userInfo.birthday_day || ''}" style="width:33%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                    </div>
                </div>
                <button id="save-user-info" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">保存修改</button>
            `;
        } else {
            modalBody.innerHTML = '<div>获取用户信息失败</div><button id="close-modal" style="margin-top:10px;padding:8px 16px;">关闭</button>';
        }

        document.getElementById('save-user-info')?.addEventListener('click', async () => {
            const userName = document.getElementById('edit-username')?.value;
            const introduction = document.getElementById('edit-introduction')?.value;
            const birthdayYear = parseInt(document.getElementById('edit-birthday-year')?.value) || 0;
            const birthdayMonth = parseInt(document.getElementById('edit-birthday-month')?.value) || 0;
            const birthdayDay = parseInt(document.getElementById('edit-birthday-day')?.value) || 0;

            const result = await apiCall('POST', '/user/update-user-info', {
                user_name: userName,
                introduction: introduction,
                birthday_year: birthdayYear,
                birthday_month: birthdayMonth,
                birthday_day: birthdayDay
            });

            if (result && result.code === 0) {
                alert('保存成功');
                modal.style.display = 'none';
                // 刷新显示
                const newInfo = await apiCall('POST', '/user/get-user-info', {});
                if (newInfo && newInfo.code === 0) {
                    document.getElementById('user-name').innerHTML = `${escapeHtml(newInfo.data.user_info.user_name || '用户')} <span class="status"></span>`;
                }
            } else {
                alert('保存失败: ' + (result?.message || '未知错误'));
            }
        });
    });

    document.getElementById('modal-title').textContent = '个人信息';
    modal.style.display = 'flex';
}