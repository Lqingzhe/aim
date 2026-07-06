package com.example.aim.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.automirrored.filled.Send
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.example.aim.network.models.MessageData
import com.example.aim.viewmodel.ChatViewModel
import java.text.SimpleDateFormat
import java.util.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatRoomScreen(
    chatId: String,
    chatType: String,
    onBack: () -> Unit,
    chatViewModel: ChatViewModel = viewModel()
) {
    var inputText by remember { mutableStateOf("") }
    val messages by chatViewModel.messages.collectAsState()
    val uiState by chatViewModel.uiState.collectAsState()
    val currentSession by chatViewModel.currentSession.collectAsState()
    val readStatusMap by chatViewModel.readStatusMap.collectAsState()
    val toastMessage by chatViewModel.toastMessage.collectAsState()
    val listState = rememberLazyListState()
    val snackbarHostState = remember { SnackbarHostState() }
    var showMoreMenu by remember { mutableStateOf(false) }
    var showGroupManagement by remember { mutableStateOf(false) }

    LaunchedEffect(toastMessage) {
        toastMessage?.let {
            snackbarHostState.showSnackbar(it)
            chatViewModel.clearToast()
        }
    }

    LaunchedEffect(chatId) {
        chatViewModel.loadMessages(chatId)
    }

    LaunchedEffect(messages.size) {
        if (messages.isNotEmpty()) {
            listState.animateScrollToItem(messages.size - 1)
        }
    }

    val sessionName = currentSession?.name ?: if (chatType == "group") "群聊" else "会话"
    val roleLabel = when (currentSession?.userRole) {
        "Owner" -> " (群主)"
        "Manager" -> " (管理员)"
        else -> ""
    }

    Scaffold(
        snackbarHost = { SnackbarHost(snackbarHostState) },
        topBar = {
            TopAppBar(
                title = {
                    Column {
                        Text("$sessionName$roleLabel", style = MaterialTheme.typography.titleMedium)
                        Text("ID: $chatId", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    }
                },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "返回")
                    }
                },
                actions = {
                    IconButton(onClick = { showMoreMenu = true }) {
                        Icon(Icons.Default.MoreVert, contentDescription = "更多")
                    }
                    DropdownMenu(expanded = showMoreMenu, onDismissRequest = { showMoreMenu = false }) {
                        if (chatType == "group") {
                            DropdownMenuItem(text = { Text("群管理") }, onClick = { showMoreMenu = false; showGroupManagement = true },
                                leadingIcon = { Icon(Icons.Default.AdminPanelSettings, contentDescription = null) })
                            DropdownMenuItem(text = { Text("发送群通知") }, onClick = { showMoreMenu = false; /* TODO */ },
                                leadingIcon = { Icon(Icons.Default.Notifications, contentDescription = null) })
                        }
                        DropdownMenuItem(text = { Text("设置备注") }, onClick = { showMoreMenu = false; /* TODO */ },
                            leadingIcon = { Icon(Icons.Default.Edit, contentDescription = null) })
                        if (chatType == "group" && currentSession?.userRole != "Owner") {
                            DropdownMenuItem(text = { Text("退出群聊") }, onClick = {
                                showMoreMenu = false
                                chatViewModel.leaveGroup(chatId)
                            }, leadingIcon = { Icon(Icons.Default.ExitToApp, contentDescription = null) })
                        }
                    }
                }
            )
        },
        bottomBar = {
            Surface(tonalElevation = 3.dp) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(8.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    OutlinedTextField(
                        value = inputText,
                        onValueChange = { inputText = it },
                        modifier = Modifier.weight(1f),
                        placeholder = { Text("输入消息...") },
                        singleLine = true
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    IconButton(
                        onClick = {
                            if (inputText.isNotBlank()) {
                                chatViewModel.sendMessage(chatId, inputText) { success ->
                                    if (success) inputText = ""
                                }
                            }
                        }
                    ) {
                        Icon(Icons.AutoMirrored.Filled.Send, contentDescription = "发送")
                    }
                }
            }
        }
    ) { padding ->
        if (uiState.isLoading) {
            Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                CircularProgressIndicator()
            }
        } else if (messages.isEmpty()) {
            Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                Text("暂无消息", color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
        } else {
            LazyColumn(
                state = listState,
                modifier = Modifier.fillMaxSize().padding(padding),
                contentPadding = PaddingValues(8.dp),
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                items(messages) { msg ->
                    val isOwn = msg.sender() == com.example.aim.data.TokenStore.userId
                    val readStatus = if (isOwn) readStatusMap[msg.id()] else null
                    MessageItem(
                        message = msg,
                        isOwnMessage = isOwn,
                        readStatus = readStatus,
                        onWithdraw = { chatViewModel.withdrawMessage(chatId, msg.id()) }
                    )
                }
            }
        }
    }

    if (showGroupManagement) {
        GroupManagementDialog(
            session = currentSession,
            onDismiss = { showGroupManagement = false },
            chatViewModel = chatViewModel
        )
    }
}

@Composable
fun MessageItem(
    message: MessageData,
    isOwnMessage: Boolean,
    readStatus: com.example.aim.viewmodel.ChatViewModel.ReadStatus? = null,
    onWithdraw: () -> Unit,
    senderName: String = ""
) {
    var showMenu by remember { mutableStateOf(false) }
    val timeStr = remember(message.time()) {
        SimpleDateFormat("HH:mm", Locale.getDefault()).format(Date(message.time() * 1000))
    }
    val displayName = if (senderName.isNotEmpty()) senderName else if (isOwnMessage) "我" else "用户 ${message.sender()}"

    Column(
        modifier = Modifier.fillMaxWidth(),
        horizontalAlignment = if (isOwnMessage) Alignment.End else Alignment.Start
    ) {
        if (!isOwnMessage) {
            Text(
                text = displayName,
                style = MaterialTheme.typography.labelSmall,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }

        Box {
            Surface(
                shape = MaterialTheme.shapes.medium,
                color = if (isOwnMessage) MaterialTheme.colorScheme.primaryContainer else MaterialTheme.colorScheme.surfaceVariant,
                modifier = Modifier.widthIn(max = 280.dp)
            ) {
                if (message.withdrawn()) {
                    Text(
                        text = "消息已撤回",
                        modifier = Modifier.padding(12.dp),
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        style = MaterialTheme.typography.bodySmall
                    )
                } else {
                    Text(
                        text = message.content(),
                        modifier = Modifier.padding(12.dp),
                        color = if (isOwnMessage) MaterialTheme.colorScheme.onPrimaryContainer else MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }

            if (isOwnMessage && !message.withdrawn()) {
                DropdownMenu(expanded = showMenu, onDismissRequest = { showMenu = false }) {
                    DropdownMenuItem(
                        text = { Text("撤回") },
                        onClick = { onWithdraw(); showMenu = false },
                        leadingIcon = { Icon(Icons.Default.Delete, contentDescription = null) }
                    )
                }
            }
        }

        Text(
            text = timeStr,
            style = MaterialTheme.typography.labelSmall,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.padding(horizontal = 4.dp)
        )

        if (isOwnMessage && readStatus != null) {
            Text(
                text = readStatus.displayText(),
                style = MaterialTheme.typography.labelSmall,
                color = if (readStatus.isRead()) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurfaceVariant,
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }

        if (isOwnMessage && !message.withdrawn()) {
            Box(modifier = Modifier.fillMaxWidth()) {
                IconButton(
                    onClick = { showMenu = true },
                    modifier = Modifier.align(Alignment.CenterEnd).size(20.dp)
                ) {
                    Icon(Icons.Default.MoreVert, contentDescription = "更多", modifier = Modifier.size(14.dp))
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GroupManagementDialog(
    session: com.example.aim.network.models.SessionItem?,
    onDismiss: () -> Unit,
    chatViewModel: ChatViewModel
) {
    var selectedAction by remember { mutableIntStateOf(0) }
    var targetUserId by remember { mutableStateOf("") }
    var muteTime by remember { mutableStateOf("") }
    var muteReason by remember { mutableStateOf("") }
    var noticeContent by remember { mutableStateOf("") }
    var groupApplyList by remember { mutableStateOf<List<String>>(emptyList()) }

    val groupId = session?.id ?: ""

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("群管理") },
        text = {
            Column {
                val actions = listOf("设置管理员", "罢免管理员", "踢出成员", "转让群主", "禁言", "解除禁言", "群通知", "入群申请")
                SingleChoiceSegmentedButtonRow(modifier = Modifier.fillMaxWidth()) {
                    actions.forEachIndexed { index, label ->
                        SegmentedButton(
                            selected = selectedAction == index,
                            onClick = { selectedAction = index },
                            shape = SegmentedButtonDefaults.itemShape(index, actions.size)
                        ) { Text(label, style = MaterialTheme.typography.labelSmall) }
                    }
                }

                Spacer(Modifier.height(12.dp))

                if (selectedAction in 0..5) {
                    OutlinedTextField(value = targetUserId, onValueChange = { targetUserId = it },
                        label = { Text("目标用户ID") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                }

                if (selectedAction == 4) {
                    Spacer(Modifier.height(8.dp))
                    OutlinedTextField(value = muteTime, onValueChange = { muteTime = it },
                        label = { Text("禁言时长(秒)") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                    Spacer(Modifier.height(8.dp))
                    OutlinedTextField(value = muteReason, onValueChange = { muteReason = it },
                        label = { Text("禁言原因") }, modifier = Modifier.fillMaxWidth())
                }

                if (selectedAction == 6) {
                    OutlinedTextField(value = noticeContent, onValueChange = { noticeContent = it },
                        label = { Text("通知内容") }, modifier = Modifier.fillMaxWidth())
                }

                if (selectedAction == 7) {
                    LaunchedEffect(Unit) {
                        chatViewModel.loadGroupApplyList(groupId) { groupApplyList = it }
                    }
                    if (groupApplyList.isEmpty()) {
                        Text("暂无申请", color = MaterialTheme.colorScheme.onSurfaceVariant)
                    } else {
                        groupApplyList.forEach { userId ->
                            ListItem(
                                headlineContent = { Text(userId) },
                                trailingContent = {
                                    Row {
                                        TextButton(onClick = { chatViewModel.agreeGroupApply(groupId, userId); groupApplyList = groupApplyList - userId }) { Text("同意") }
                                        TextButton(onClick = { chatViewModel.refuseGroupApply(groupId, userId); groupApplyList = groupApplyList - userId }) { Text("拒绝") }
                                    }
                                }
                            )
                        }
                    }
                }
            }
        },
        confirmButton = {
            if (selectedAction != 7) {
                TextButton(onClick = {
                    when (selectedAction) {
                        0 -> chatViewModel.setManager(groupId, targetUserId)
                        1 -> chatViewModel.revokeManager(groupId, targetUserId)
                        2 -> chatViewModel.kickOutGroup(groupId, targetUserId)
                        3 -> chatViewModel.transformGroupOwner(groupId, targetUserId)
                        4 -> chatViewModel.setMute(groupId, targetUserId, muteTime.toIntOrNull() ?: 0, muteReason)
                        5 -> chatViewModel.releaseMute(groupId, targetUserId)
                        6 -> chatViewModel.sendGroupNotice(groupId, noticeContent)
                    }
                    onDismiss()
                }) { Text("确定") }
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("关闭") }
        }
    )
}
