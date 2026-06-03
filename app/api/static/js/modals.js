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

// ============================================
// 群管理功能
// ============================================

// 1. 获取群申请列表
function showGroupApplyListModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">群ID</label>
            <input type="text" id="group-apply-group-id" placeholder="请输入群ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <button id="fetch-group-apply-list" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">获取申请列表</button>
        <div id="group-apply-list-container" style="margin-top:15px;"></div>
    `;
    document.getElementById('modal-title').textContent = '群申请列表';
    modal.style.display = 'flex';

    document.getElementById('fetch-group-apply-list').onclick = async () => {
        const groupId = document.getElementById('group-apply-group-id').value.trim();
        if (!groupId) {
            alert('请输入群ID');
            return;
        }

        const result = await apiCall('POST', '/group/get-group-apply-list', { group_id: groupId });
        console.log('群申请列表:', result);

        const container = document.getElementById('group-apply-list-container');
        if (result && result.code === 0 && result.data?.group_info?.group_id_list) {
            const applyList = result.data.group_info.group_id_list;
            if (applyList.length === 0) {
                container.innerHTML = '<div style="text-align:center;padding:20px;color:#9ca3af;">暂无申请</div>';
                return;
            }

            let html = '<div style="max-height:300px;overflow-y:auto;">';
            for (const userId of applyList) {
                html += `<div style="display:flex;justify-content:space-between;align-items:center;padding:10px;border-bottom:1px solid #eee;">
                            <span>申请用户ID: ${userId}</span>
                            <div>
                                <button class="agree-group-apply" data-group-id="${groupId}" data-user-id="${userId}" style="padding:4px 12px;background:#10b981;color:white;border:none;border-radius:4px;cursor:pointer;margin-right:8px;">同意</button>
                                <button class="refuse-group-apply" data-group-id="${groupId}" data-user-id="${userId}" style="padding:4px 12px;background:#ef4444;color:white;border:none;border-radius:4px;cursor:pointer;">拒绝</button>
                            </div>
                        </div>`;
            }
            html += '</div>';
            container.innerHTML = html;

            document.querySelectorAll('.agree-group-apply').forEach(btn => {
                btn.onclick = async () => {
                    const groupId = btn.dataset.groupId;
                    const goalUserId = btn.dataset.userId;
                    const result = await apiCall('POST', '/group/agree-group-apply', { group_id: groupId, goal_user_id: goalUserId });
                    if (result && result.code === 0) {
                        alert('已同意入群申请');
                        modal.style.display = 'none';
                        loadSessions();
                    } else {
                        alert('操作失败: ' + (result?.message || '未知错误'));
                    }
                };
            });

            document.querySelectorAll('.refuse-group-apply').forEach(btn => {
                btn.onclick = async () => {
                    const groupId = btn.dataset.groupId;
                    const goalUserId = btn.dataset.userId;
                    const result = await apiCall('POST', '/group/refuse-group-apply', { group_id: groupId, goal_user_id: goalUserId });
                    if (result && result.code === 0) {
                        alert('已拒绝入群申请');
                        modal.style.display = 'none';
                    } else {
                        alert('操作失败: ' + (result?.message || '未知错误'));
                    }
                };
            });
        } else {
            container.innerHTML = '<div style="text-align:center;padding:20px;color:#ef4444;">获取失败: ' + (result?.message || '未知错误') + '</div>';
        }
    };
}

// 2. 设置/罢免管理员
function showManageAdminModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">群ID</label>
            <input type="text" id="admin-group-id" placeholder="请输入群ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">目标用户ID</label>
            <input type="text" id="admin-goal-user-id" placeholder="请输入目标用户ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="display:flex;gap:10px;">
            <button id="set-manager-btn" style="flex:1;padding:10px;background:#10b981;color:white;border:none;border-radius:6px;cursor:pointer;">设为管理员</button>
            <button id="revoke-manager-btn" style="flex:1;padding:10px;background:#ef4444;color:white;border:none;border-radius:6px;cursor:pointer;">罢免管理员</button>
        </div>
    `;
    document.getElementById('modal-title').textContent = '管理员管理';
    modal.style.display = 'flex';

    document.getElementById('set-manager-btn').onclick = async () => {
        const groupId = document.getElementById('admin-group-id').value.trim();
        const goalUserId = document.getElementById('admin-goal-user-id').value.trim();
        if (!groupId || !goalUserId) {
            alert('请填写群ID和目标用户ID');
            return;
        }
        const result = await apiCall('POST', '/group/set-manager', { group_id: groupId, goal_user_id: goalUserId });
        if (result && result.code === 0) {
            alert('设置管理员成功');
            modal.style.display = 'none';
        } else {
            alert('操作失败: ' + (result?.message || '未知错误'));
        }
    };

    document.getElementById('revoke-manager-btn').onclick = async () => {
        const groupId = document.getElementById('admin-group-id').value.trim();
        const goalUserId = document.getElementById('admin-goal-user-id').value.trim();
        if (!groupId || !goalUserId) {
            alert('请填写群ID和目标用户ID');
            return;
        }
        const result = await apiCall('POST', '/group/revoke-manager', { group_id: groupId, goal_user_id: goalUserId });
        if (result && result.code === 0) {
            alert('罢免管理员成功');
            modal.style.display = 'none';
        } else {
            alert('操作失败: ' + (result?.message || '未知错误'));
        }
    };
}

// 3. 踢出群聊
function showKickOutModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">群ID</label>
            <input type="text" id="kick-group-id" placeholder="请输入群ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">目标用户ID</label>
            <input type="text" id="kick-goal-user-id" placeholder="请输入要踢出的用户ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <button id="submit-kick-out" style="width:100%;padding:10px;background:#ef4444;color:white;border:none;border-radius:6px;cursor:pointer;">确认踢出</button>
    `;
    document.getElementById('modal-title').textContent = '踢出群聊';
    modal.style.display = 'flex';

    document.getElementById('submit-kick-out').onclick = async () => {
        const groupId = document.getElementById('kick-group-id').value.trim();
        const goalUserId = document.getElementById('kick-goal-user-id').value.trim();
        if (!groupId || !goalUserId) {
            alert('请填写群ID和目标用户ID');
            return;
        }
        if (confirm(`确定要将用户 ${goalUserId} 踢出群聊 ${groupId} 吗？`)) {
            const result = await apiCall('POST', '/group/kick-out-group', { group_id: groupId, goal_user_id: goalUserId });
            if (result && result.code === 0) {
                alert('踢出成功');
                modal.style.display = 'none';
                loadSessions();
            } else {
                alert('操作失败: ' + (result?.message || '未知错误'));
            }
        }
    };
}

// 4. 转让群主
function showTransferOwnerModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">群ID</label>
            <input type="text" id="transfer-group-id" placeholder="请输入群ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">新群主用户ID</label>
            <input type="text" id="transfer-goal-user-id" placeholder="请输入新群主用户ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <button id="submit-transfer" style="width:100%;padding:10px;background:#f59e0b;color:white;border:none;border-radius:6px;cursor:pointer;">确认转让</button>
    `;
    document.getElementById('modal-title').textContent = '转让群主';
    modal.style.display = 'flex';

    document.getElementById('submit-transfer').onclick = async () => {
        const groupId = document.getElementById('transfer-group-id').value.trim();
        const goalUserId = document.getElementById('transfer-goal-user-id').value.trim();
        if (!groupId || !goalUserId) {
            alert('请填写群ID和新群主用户ID');
            return;
        }
        if (confirm(`确定要将群 ${groupId} 转让给用户 ${goalUserId} 吗？`)) {
            const result = await apiCall('POST', '/group/transform-group-owner', { group_id: groupId, goal_user_id: goalUserId });
            if (result && result.code === 0) {
                alert('转让成功');
                modal.style.display = 'none';
                loadSessions();
            } else {
                alert('操作失败: ' + (result?.message || '未知错误'));
            }
        }
    };
}

// 5. 禁言
function showSetMuteModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">群ID</label>
            <input type="text" id="mute-group-id" placeholder="请输入群ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">目标用户ID</label>
            <input type="text" id="mute-goal-user-id" placeholder="请输入要禁言的用户ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">禁言时长(秒)</label>
            <input type="number" id="mute-time-seconds" placeholder="请输入禁言时长，单位秒" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">禁言原因</label>
            <textarea id="mute-reason" rows="2" placeholder="请输入禁言原因" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;"></textarea>
        </div>
        <button id="submit-set-mute" style="width:100%;padding:10px;background:#f59e0b;color:white;border:none;border-radius:6px;cursor:pointer;">确认禁言</button>
    `;
    document.getElementById('modal-title').textContent = '禁言';
    modal.style.display = 'flex';

    document.getElementById('submit-set-mute').onclick = async () => {
        const groupId = document.getElementById('mute-group-id').value.trim();
        const goalUserId = document.getElementById('mute-goal-user-id').value.trim();
        const muteTimeSeconds = parseInt(document.getElementById('mute-time-seconds').value);
        const muteReason = document.getElementById('mute-reason').value.trim();

        if (!groupId || !goalUserId) {
            alert('请填写群ID和目标用户ID');
            return;
        }
        if (!muteTimeSeconds || muteTimeSeconds <= 0) {
            alert('请输入有效的禁言时长');
            return;
        }

        const result = await apiCall('POST', '/group/set-mute', {
            group_id: groupId,
            goal_user_id: goalUserId,
            mute_time_seconds: muteTimeSeconds,
            mute_reason: muteReason
        });

        if (result && result.code === 0) {
            alert('禁言设置成功');
            modal.style.display = 'none';
        } else {
            alert('操作失败: ' + (result?.message || '未知错误'));
        }
    };
}

// 6. 解除禁言
function showReleaseMuteModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">群ID</label>
            <input type="text" id="release-group-id" placeholder="请输入群ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">目标用户ID</label>
            <input type="text" id="release-goal-user-id" placeholder="请输入要解除禁言的用户ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <button id="submit-release-mute" style="width:100%;padding:10px;background:#10b981;color:white;border:none;border-radius:6px;cursor:pointer;">确认解除禁言</button>
    `;
    document.getElementById('modal-title').textContent = '解除禁言';
    modal.style.display = 'flex';

    document.getElementById('submit-release-mute').onclick = async () => {
        const groupId = document.getElementById('release-group-id').value.trim();
        const goalUserId = document.getElementById('release-goal-user-id').value.trim();
        if (!groupId || !goalUserId) {
            alert('请填写群ID和目标用户ID');
            return;
        }
        const result = await apiCall('POST', '/group/release-mute', { group_id: groupId, goal_user_id: goalUserId });
        if (result && result.code === 0) {
            alert('解除禁言成功');
            modal.style.display = 'none';
        } else {
            alert('操作失败: ' + (result?.message || '未知错误'));
        }
    };
}

// 备注功能
function showRemarkModal() {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = `
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">对方用户ID</label>
            <input type="text" id="remark-goal-user-id" placeholder="请输入要设置备注的用户ID" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
        </div>
        <div style="margin-bottom:15px;">
            <label style="display:block;margin-bottom:5px;font-weight:bold;">备注名</label>
            <input type="text" id="remark-nickname" placeholder="输入备注名（留空表示删除备注）" style="width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;">
            <div style="font-size:12px;color:#9ca3af;margin-top:5px;">💡 提示：备注名留空表示删除备注</div>
        </div>
        <button id="submit-remark" style="width:100%;padding:10px;background:#667eea;color:white;border:none;border-radius:6px;cursor:pointer;">保存备注</button>
    `;
    document.getElementById('modal-title').textContent = '设置备注';
    modal.style.display = 'flex';

    document.getElementById('submit-remark').onclick = async () => {
        const goalUserId = document.getElementById('remark-goal-user-id').value.trim();
        const nickName = document.getElementById('remark-nickname').value.trim();

        if (!goalUserId) {
            alert('请输入对方用户ID');
            return;
        }

        const result = await apiCall('POST', '/user/remark', {
            goal_user_id: goalUserId,
            nick_name: nickName
        });

        if (result && result.code === 0) {
            if (nickName) {
                alert(`备注设置成功！`);
            } else {
                alert(`备注已取消！`);
            }
            modal.style.display = 'none';
            loadSessions();
        } else {
            alert('备注设置失败: ' + (result?.message || '未知错误'));
        }
    };
}