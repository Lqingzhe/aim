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

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function showUserInfoModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');

    modalBody.innerHTML = '<div style="text-align:center;padding:20px;">加载中...</div>';
    document.getElementById('modal-title').textContent = '个人信息';
    modal.style.display = 'flex';

    try {
        const result = await apiCall('POST', '/user/get-user-info', {});
        console.log('get-user-info 返回:', result);

        if (result && result.code === 0 && result.data?.user_info) {
            let userInfoData = result.data.user_info;
            let userInfo = userInfoData.UserInfo || userInfoData;
            let remarkInfos = userInfoData.RemarkInfos || userInfoData.remark_info || [];

            console.log('用户信息:', userInfo);
            console.log('备注信息:', remarkInfos);

            modalBody.innerHTML = `
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">用户ID</label>
                    <input type="text" value="${userInfo.UserID || userInfo.user_id || ''}" disabled style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;background:#f5f5f5;">
                </div>
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">用户名</label>
                    <input type="text" id="edit-username" value="${escapeHtml(userInfo.UserName || userInfo.user_name || '')}" placeholder="请输入用户名" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                </div>
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">个性签名</label>
                    <textarea id="edit-introduction" rows="3" placeholder="请输入个性签名" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">${escapeHtml(userInfo.Introduction || userInfo.introduction || '')}</textarea>
                </div>
                <div style="margin-bottom:15px;">
                    <label style="display:block;margin-bottom:5px;font-weight:bold;">生日</label>
                    <div style="display:flex;gap:8px;">
                        <input type="number" id="edit-birthday-year" placeholder="年" value="${userInfo.BirthdayYear || userInfo.birthday_year || ''}" style="width:33%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                        <input type="number" id="edit-birthday-month" placeholder="月" value="${userInfo.BirthdayMonth || userInfo.birthday_month || ''}" style="width:33%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                        <input type="number" id="edit-birthday-day" placeholder="日" value="${userInfo.BirthdayDay || userInfo.birthday_day || ''}" style="width:33%;padding:8px;border:1px solid #ddd;border-radius:4px;">
                    </div>
                </div>
                <button id="save-user-info" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">保存修改</button>
            `;

            document.getElementById('save-user-info').onclick = async () => {
                const updateData = {
                    user_name: document.getElementById('edit-username')?.value || '',
                    introduction: document.getElementById('edit-introduction')?.value || '',
                    birthday_year: parseInt(document.getElementById('edit-birthday-year')?.value) || 0,
                    birthday_month: parseInt(document.getElementById('edit-birthday-month')?.value) || 0,
                    birthday_day: parseInt(document.getElementById('edit-birthday-day')?.value) || 0
                };

                const saveResult = await apiCall('POST', '/user/update-user-info', updateData);
                if (saveResult && saveResult.code === 0) {
                    alert('保存成功');
                    modal.style.display = 'none';
                    const newInfo = await apiCall('POST', '/user/get-user-info', {});
                    if (newInfo && newInfo.code === 0) {
                        let userInfoData = newInfo.data.user_info;
                        let userName = userInfoData?.UserInfo?.UserName || userInfoData?.user_name || '用户';
                        document.getElementById('user-name').innerHTML = `${escapeHtml(userName)} <span class="status"></span>`;
                    }
                } else {
                    alert('保存失败: ' + (saveResult?.message || '未知错误'));
                }
            };
        } else {
            modalBody.innerHTML = `<div style="color:red;">获取用户信息失败</div>
                <button id="close-modal" style="margin-top:10px;padding:8px 16px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">关闭</button>`;
            document.getElementById('close-modal').onclick = () => modal.style.display = 'none';
        }
    } catch (error) {
        console.error('获取用户信息异常:', error);
        modalBody.innerHTML = `<div style="color:red;">网络错误</div>
            <button id="close-modal" style="margin-top:10px;padding:8px 16px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">关闭</button>`;
        document.getElementById('close-modal').onclick = () => modal.style.display = 'none';
    }
}