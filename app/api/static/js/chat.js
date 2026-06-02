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
    loadSessions();
    if (currentSession) {
        loadHistoryMessages(currentSession.id, false);
    }
    addSystemMessage('✅ 已刷新');
}

// 初始化
function init() {
    connectWebSocket();
    loadSessions();

    document.getElementById('send-btn').onclick = sendMessage;
    document.getElementById('btn-logout').onclick = logout;
    document.getElementById('btn-add-friend').onclick = showAddFriendModal;
    document.getElementById('btn-create-group').onclick = showCreateGroupModal;
    document.getElementById('btn-join-group').onclick = showJoinGroupModal;
    document.getElementById('btn-friend-apply-list').onclick = showFriendApplyList;
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
        if (e.target === modal) modal.style.display = 'none';
    };
}

init();