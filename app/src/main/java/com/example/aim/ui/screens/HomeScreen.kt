package com.example.aim.ui.screens

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.example.aim.network.models.SessionItem
import com.example.aim.viewmodel.ChatViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(
    onOpenChat: (String, String) -> Unit,
    onOpenProfile: () -> Unit,
    onLogout: () -> Unit,
    chatViewModel: ChatViewModel = viewModel()
) {
    var selectedTab by remember { mutableIntStateOf(0) }
    val sessions by chatViewModel.sessions.collectAsState()
    val friendRequests by chatViewModel.friendRequests.collectAsState()
    val toastMessage by chatViewModel.toastMessage.collectAsState()
    val snackbarHostState = remember { SnackbarHostState() }

    var showAddFriend by remember { mutableStateOf(false) }
    var showCreateGroup by remember { mutableStateOf(false) }
    var showJoinGroup by remember { mutableStateOf(false) }

    LaunchedEffect(Unit) {
        chatViewModel.loadSessionsAndGroups()
        chatViewModel.loadFriendRequests()
    }

    LaunchedEffect(toastMessage) {
        toastMessage?.let {
            snackbarHostState.showSnackbar(it)
            chatViewModel.clearToast()
        }
    }

    Scaffold(
        snackbarHost = { SnackbarHost(snackbarHostState) },
        topBar = {
            TopAppBar(
                title = { Text("AIM") },
                actions = {
                    when (selectedTab) {
                        0 -> {
                            IconButton(onClick = { showCreateGroup = true }) {
                                Icon(Icons.Default.GroupAdd, contentDescription = "创建群聊")
                            }
                        }
                        1 -> {
                            IconButton(onClick = { showJoinGroup = true }) {
                                Icon(Icons.Default.PersonAdd, contentDescription = "添加好友")
                            }
                        }
                    }
                    IconButton(onClick = onOpenProfile) {
                        Icon(Icons.Default.Person, contentDescription = "个人中心")
                    }
                }
            )
        },
        bottomBar = {
            NavigationBar {
                NavigationBarItem(
                    icon = { Icon(Icons.Default.Email, contentDescription = null) },
                    label = { Text("消息") },
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 }
                )
                NavigationBarItem(
                    icon = { Icon(Icons.Default.Contacts, contentDescription = null) },
                    label = { Text("联系人") },
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 }
                )
                NavigationBarItem(
                    icon = { Icon(Icons.Default.Settings, contentDescription = null) },
                    label = { Text("设置") },
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 }
                )
            }
        }
    ) { padding ->
        when (selectedTab) {
            0 -> ChatListTab(
                sessions = sessions,
                onRefresh = { chatViewModel.loadSessionsAndGroups() },
                onOpenChat = { session ->
                    chatViewModel.selectSession(session)
                    onOpenChat(session.id, session.type)
                },
                onDeleteSession = { chatViewModel.deleteSession(it) },
                onDeleteGroup = { chatViewModel.deleteGroup(it) },
                modifier = Modifier.padding(padding)
            )
            1 -> ContactsTab(
                sessions = sessions,
                friendRequests = friendRequests,
                onOpenChat = { session ->
                    chatViewModel.selectSession(session)
                    onOpenChat(session.id, session.type)
                },
                onAcceptFriend = { chatViewModel.acceptFriend(it) },
                onRefuseFriend = { chatViewModel.refuseFriend(it) },
                onAddFriend = { showAddFriend = true },
                modifier = Modifier.padding(padding)
            )
            2 -> SettingsTab(
                onLogout = { all ->
                    if (all) chatViewModel.logoutAllDevices { onLogout() }
                    else chatViewModel.logoutCurrentDevice { onLogout() }
                },
                modifier = Modifier.padding(padding)
            )
        }
    }

    if (showAddFriend) {
        AddFriendDialog(
            onDismiss = { showAddFriend = false },
            onConfirm = { userId ->
                chatViewModel.applyFriend(userId) { success, msg ->
                    showAddFriend = false
                    chatViewModel.showToast(msg)
                }
            }
        )
    }

    if (showCreateGroup) {
        CreateGroupDialog(
            onDismiss = { showCreateGroup = false },
            onCreate = { name ->
                chatViewModel.createGroup(name) { success, _ ->
                    showCreateGroup = false
                    if (success) chatViewModel.loadSessionsAndGroups()
                }
            }
        )
    }

    if (showJoinGroup) {
        JoinGroupDialog(
            onDismiss = { showJoinGroup = false },
            onJoin = { groupId ->
                chatViewModel.joinGroup(groupId) { success, msg ->
                    showJoinGroup = false
                    chatViewModel.showToast(msg)
                }
            }
        )
    }
}

@Composable
fun ChatListTab(
    sessions: List<SessionItem>,
    onRefresh: () -> Unit,
    onOpenChat: (SessionItem) -> Unit,
    onDeleteSession: (String) -> Unit,
    onDeleteGroup: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val listState = rememberLazyListState()
    var isRefreshing by remember { mutableStateOf(false) }

    LaunchedEffect(listState.firstVisibleItemScrollOffset) {
        if (listState.firstVisibleItemIndex == 0 && listState.firstVisibleItemScrollOffset == 0 && isRefreshing) {
            isRefreshing = false
        }
    }

    Box(modifier = modifier.fillMaxSize()) {
        if (sessions.isEmpty() && !isRefreshing) {
            Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                Text("下拉刷新或暂无会话", color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
        }

        LazyColumn(
            state = listState,
            modifier = Modifier.fillMaxSize()
        ) {
            items(sessions) { session ->
                SessionListItem(
                    session = session,
                    onClick = { onOpenChat(session) },
                    onDelete = {
                        if (session.type == "session") onDeleteSession(session.id)
                        else if (session.userRole == "Owner") onDeleteGroup(session.id)
                    }
                )
            }
        }

        if (isRefreshing) {
            CircularProgressIndicator(
                modifier = Modifier.align(Alignment.TopCenter).padding(16.dp)
            )
        }

        IconButton(
            onClick = {
                isRefreshing = true
                onRefresh()
                isRefreshing = false
            },
            modifier = Modifier.align(Alignment.BottomEnd).padding(16.dp)
        ) {
            Icon(Icons.Default.Refresh, contentDescription = "刷新")
        }
    }
}

@Composable
fun ContactsTab(
    sessions: List<SessionItem>,
    friendRequests: List<String>,
    onOpenChat: (SessionItem) -> Unit,
    onAcceptFriend: (String) -> Unit,
    onRefuseFriend: (String) -> Unit,
    onAddFriend: () -> Unit,
    modifier: Modifier = Modifier
) {
    LazyColumn(modifier = modifier) {
        item {
            Row(
                modifier = Modifier.fillMaxWidth().padding(16.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text("好友申请", style = MaterialTheme.typography.titleMedium)
                if (friendRequests.isNotEmpty()) {
                    Badge { Text("${friendRequests.size}") }
                }
            }
        }

        if (friendRequests.isNotEmpty()) {
            items(friendRequests) { userId ->
                ListItem(
                    headlineContent = { Text(userId) },
                    trailingContent = {
                        Row {
                            TextButton(onClick = { onAcceptFriend(userId) }) { Text("同意", color = MaterialTheme.colorScheme.primary) }
                            TextButton(onClick = { onRefuseFriend(userId) }) { Text("拒绝", color = MaterialTheme.colorScheme.error) }
                        }
                    }
                )
                HorizontalDivider()
            }
        }

        item {
            TextButton(onClick = onAddFriend, modifier = Modifier.fillMaxWidth().padding(8.dp)) {
                Icon(Icons.Default.PersonAdd, contentDescription = null)
                Spacer(Modifier.width(8.dp))
                Text("添加好友")
            }
        }

        val groupSessions = sessions.filter { it.type == "group" }
        if (groupSessions.isNotEmpty()) {
            item {
                Text("群聊", style = MaterialTheme.typography.titleMedium, modifier = Modifier.padding(16.dp, 12.dp))
            }
            items(groupSessions) { session ->
                SessionListItem(session = session, onClick = { onOpenChat(session) }, onDelete = {})
            }
        }

        val privateSessions = sessions.filter { it.type == "session" }
        if (privateSessions.isNotEmpty()) {
            item {
                Text("好友", style = MaterialTheme.typography.titleMedium, modifier = Modifier.padding(16.dp, 12.dp))
            }
            items(privateSessions) { session ->
                SessionListItem(session = session, onClick = { onOpenChat(session) }, onDelete = {})
            }
        }
    }
}

@Composable
fun SessionListItem(
    session: SessionItem,
    onClick: () -> Unit,
    onDelete: () -> Unit
) {
    val icon = if (session.type == "group") Icons.Default.Group else Icons.Default.ChatBubble
    val roleColors = mapOf("Owner" to MaterialTheme.colorScheme.tertiary, "Manager" to MaterialTheme.colorScheme.primary)
    val roleLabels = mapOf("Owner" to "群主", "Manager" to "管理员")

    ListItem(
        headlineContent = {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(session.name, fontWeight = if (session.hasNewMessage) FontWeight.Bold else FontWeight.Normal)
                if (session.userRole == "Owner" || session.userRole == "Manager") {
                    Spacer(Modifier.width(6.dp))
                    SuggestionChip(
                        onClick = {},
                        label = { Text(roleLabels[session.userRole] ?: session.userRole!!, style = MaterialTheme.typography.labelSmall) }
                    )
                }
            }
        },
        supportingContent = {
            Text(
                if (session.type == "group") "群聊 ID: ${session.id.take(12)}" else "私聊 · ${session.goalUserId ?: session.id.take(12)}",
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
        },
        leadingContent = { Icon(icon, contentDescription = null) },
        trailingContent = {
            if ((session.type == "session") || (session.type == "group" && session.userRole == "Owner")) {
                IconButton(onClick = onDelete) {
                    Icon(Icons.Default.Delete, contentDescription = "删除", tint = MaterialTheme.colorScheme.error)
                }
            }
        },
        modifier = Modifier.clickable(onClick = onClick)
    )
    HorizontalDivider()
}

@Composable
fun SettingsTab(onLogout: (Boolean) -> Unit, modifier: Modifier = Modifier) {
    var showLogoutDialog by remember { mutableStateOf(false) }
    var logoutAll by remember { mutableStateOf(false) }

    Column(modifier = modifier) {
        ListItem(headlineContent = { Text("账号设置") }, leadingContent = { Icon(Icons.Default.AccountCircle, contentDescription = null) })
        HorizontalDivider()
        ListItem(headlineContent = { Text("退出当前设备") }, leadingContent = { Icon(Icons.Default.Logout, contentDescription = null) },
            modifier = Modifier.clickable { logoutAll = false; showLogoutDialog = true })
        HorizontalDivider()
        ListItem(headlineContent = { Text("退出所有设备") }, leadingContent = { Icon(Icons.Default.ExitToApp, contentDescription = null) },
            modifier = Modifier.clickable { logoutAll = true; showLogoutDialog = true })
    }

    if (showLogoutDialog) {
        AlertDialog(
            onDismissRequest = { showLogoutDialog = false },
            title = { Text("确认退出") },
            text = { Text(if (logoutAll) "确定要退出所有设备吗？" else "确定要退出当前设备吗？") },
            confirmButton = {
                TextButton(onClick = { showLogoutDialog = false; onLogout(logoutAll) }) { Text("确定") }
            },
            dismissButton = {
                TextButton(onClick = { showLogoutDialog = false }) { Text("取消") }
            }
        )
    }
}

@Composable
fun AddFriendDialog(onDismiss: () -> Unit, onConfirm: (String) -> Unit) {
    var userId by remember { mutableStateOf("") }
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("添加好友") },
        text = {
            OutlinedTextField(
                value = userId,
                onValueChange = { userId = it },
                label = { Text("对方用户ID") },
                singleLine = true
            )
        },
        confirmButton = {
            TextButton(onClick = { if (userId.isNotBlank()) onConfirm(userId) }) { Text("发送申请") }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("取消") }
        }
    )
}

@Composable
fun CreateGroupDialog(onDismiss: () -> Unit, onCreate: (String) -> Unit) {
    var groupName by remember { mutableStateOf("") }
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("创建群聊") },
        text = {
            OutlinedTextField(
                value = groupName,
                onValueChange = { groupName = it },
                label = { Text("群名称") },
                singleLine = true
            )
        },
        confirmButton = {
            TextButton(onClick = { if (groupName.isNotBlank()) onCreate(groupName) }) { Text("创建") }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("取消") }
        }
    )
}

@Composable
fun JoinGroupDialog(onDismiss: () -> Unit, onJoin: (String) -> Unit) {
    var groupId by remember { mutableStateOf("") }
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("加入群聊") },
        text = {
            OutlinedTextField(
                value = groupId,
                onValueChange = { groupId = it },
                label = { Text("群ID") },
                singleLine = true
            )
        },
        confirmButton = {
            TextButton(onClick = { if (groupId.isNotBlank()) onJoin(groupId) }) { Text("申请加入") }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("取消") }
        }
    )
}
