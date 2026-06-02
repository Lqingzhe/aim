// ============================================
// 弹窗功能
// ============================================
function showAddFriendModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <input type="text" id="friend-user-id" placeholder="对方用户ID">
        <button id="submit-add-friend">发送好友申请</button>
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

function showCreateGroupModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <input type="text" id="group-name" placeholder="群名称">
        <button id="submit-create-group">创建群聊</button>
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

function showJoinGroupModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <input type="text" id="join-group-id" placeholder="群ID">
        <button id="submit-join-group">申请加入</button>
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
            html += `<div style="display:flex;justify-content:space-between;align-items:center;padding:10px;border-bottom:1px solid #eee;">
                        <span>用户 ${userId}</span>
                        <div>
                            <button class="agree-apply" data-user-id="${userId}" style="padding:4px 12px;background:#10b981;color:white;border:none;border-radius:4px;cursor:pointer;margin-right:8px;">同意</button>
                            <button class="refuse-apply" data-user-id="${userId}" style="padding:4px 12px;background:#ef4444;color:white;border:none;border-radius:4px;cursor:pointer;">拒绝</button>
                        </div>
                    </div>`;
        }
        html += '</div>';
        modalBody.innerHTML = html;
        document.getElementById('modal-title').textContent = '好友申请列表';
        modal.style.display = 'flex';

        // 同意按钮
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

        // 拒绝按钮
        document.querySelectorAll('.refuse-apply').forEach(btn => {
            btn.onclick = async () => {
                const goalUserId = btn.dataset.userId;
                const refuseResult = await apiCall('POST', '/group/refuse-friend-apply', { goal_user_id: goalUserId });
                if (refuseResult && refuseResult.code === 0) {
                    alert('已拒绝好友申请');
                    modal.style.display = 'none';
                    loadSessions();
                } else {
                    alert('操作失败: ' + (refuseResult?.message || '未知错误'));
                }
            };
        });
    } else {
        alert('获取申请列表失败: ' + (result?.message || '未知错误'));
    }
}