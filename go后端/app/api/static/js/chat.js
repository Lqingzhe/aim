// ============================================
// 主入口 - 整合所有模块
// ============================================

// 检查登录
if (!accessToken || !userID) {
    window.location.href = '/login';
}

// 显示用户信息
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

// 使用事件委托绑定所有按钮（更可靠）
function bindAllButtons() {
    console.log('开始绑定按钮事件...');

    document.getElementById('send-btn').onclick = sendMessage;
    document.getElementById('btn-logout-device').onclick = logoutDevice;
    document.getElementById('btn-logout-all').onclick = logoutAll;
    document.getElementById('btn-add-friend').onclick = showAddFriendModal;
    document.getElementById('btn-create-group').onclick = showCreateGroupModal;
    document.getElementById('btn-join-group').onclick = showJoinGroupModal;
    document.getElementById('btn-friend-apply-list').onclick = showFriendApplyList;
    document.getElementById('btn-group-apply-list').onclick = showGroupApplyListModal;
    document.getElementById('btn-send-group-notice').onclick = showSendGroupNoticeModal;
    document.getElementById('btn-remark').onclick = showRemarkModal;
    document.getElementById('btn-manage-admin').onclick = showManageAdminModal;
    document.getElementById('btn-kick-out').onclick = showKickOutModal;
    document.getElementById('btn-transfer-owner').onclick = showTransferOwnerModal;
    document.getElementById('btn-set-mute').onclick = showSetMuteModal;
    document.getElementById('btn-release-mute').onclick = showReleaseMuteModal;
    document.getElementById('btn-ai-config').onclick = showGetAIConfigModal;
    document.getElementById('btn-clear-chat-context').onclick = showDeleteChatContextModal;
    document.getElementById('btn-refresh').onclick = refreshAll;
    document.getElementById('btn-withdraw-latest').onclick = withdrawLatestMessage;
    document.getElementById('btn-user-info').onclick = showUserInfoModal;

    console.log('所有按钮事件绑定完成');
}
function init() {
    console.log('初始化开始...');

    if (typeof connectWebSocket === 'function') {
        connectWebSocket();
    } else {
        console.error('connectWebSocket 函数未定义');
    }

    if (typeof loadSessions === 'function') {
        loadSessions();
    } else {
        console.error('loadSessions 函数未定义');
    }

    // 绑定所有按钮事件
    bindAllButtons();

    // 消息输入框回车发送
    const messageInput = document.getElementById('message-input');
    if (messageInput) {
        messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                sendMessage();
            }
        });
    }

    // 弹窗关闭
    const modalClose = document.getElementById('modal-close');
    if (modalClose) {
        modalClose.onclick = () => {
            const modal = document.getElementById('modal');
            if (modal) modal.style.display = 'none';
        };
    }

    window.onclick = (e) => {
        const modal = document.getElementById('modal');
        if (modal && e.target === modal) modal.style.display = 'none';
    };

    console.log('初始化完成');
}

// 确保 DOM 加载完成后执行
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM 已加载，立即执行
    init();
}