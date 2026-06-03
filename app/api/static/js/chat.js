// ============================================
// 主入口 - 整合所有模块
// ============================================

if (!accessToken || !userID) {
    window.location.href = '/login';
}

document.getElementById('user-name').innerHTML = `用户 ${userID.substring(0, 10)}... <span class="status"></span>`;
document.getElementById('user-id').textContent = `ID: ${userID}`;

function refreshAll() {
    if (typeof loadSessions === 'function') {
        loadSessions();
    }
    if (currentSession && typeof loadHistoryMessages === 'function') {
        loadHistoryMessages(currentSession.id, false);
    }
    if (typeof addSystemMessage === 'function') {
        addSystemMessage('✅ 已刷新');
    }
}

function init() {
    console.log('初始化开始...');

    if (typeof connectWebSocket === 'function') {
        connectWebSocket();
    }

    if (typeof loadSessions === 'function') {
        loadSessions();
    }

    // 绑定按钮事件
    document.getElementById('send-btn').onclick = sendMessage;
    document.getElementById('btn-logout').onclick = logout;
    document.getElementById('btn-add-friend').onclick = showAddFriendModal;
    document.getElementById('btn-create-group').onclick = showCreateGroupModal;
    document.getElementById('btn-join-group').onclick = showJoinGroupModal;
    document.getElementById('btn-friend-apply-list').onclick = showFriendApplyList;
    document.getElementById('btn-group-apply-list').onclick = showGroupApplyListModal;
    document.getElementById('btn-remark').onclick = showRemarkModal;
    document.getElementById('btn-manage-admin').onclick = showManageAdminModal;
    document.getElementById('btn-kick-out').onclick = showKickOutModal;
    document.getElementById('btn-transfer-owner').onclick = showTransferOwnerModal;
    document.getElementById('btn-set-mute').onclick = showSetMuteModal;
    document.getElementById('btn-release-mute').onclick = showReleaseMuteModal;
    document.getElementById('btn-refresh').onclick = refreshAll;
    document.getElementById('btn-withdraw-latest').onclick = withdrawLatestMessage;
    document.getElementById('btn-user-info').onclick = showUserInfoModal;

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
        if (modal && e.target === modal) modal.style.display = 'none';
    };

    console.log('初始化完成');
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}