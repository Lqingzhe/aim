package com.example.aim.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.example.aim.network.models.UpdateUserInfoRequest
import com.example.aim.viewmodel.ChatViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProfileScreen(
    onBack: () -> Unit,
    chatViewModel: ChatViewModel = viewModel()
) {
    val userInfo by chatViewModel.userInfo.collectAsState()
    var isEditing by remember { mutableStateOf(false) }
    var userName by remember { mutableStateOf("") }
    var introduction by remember { mutableStateOf("") }
    var birthYear by remember { mutableStateOf("") }
    var birthMonth by remember { mutableStateOf("") }
    var birthDay by remember { mutableStateOf("") }
    var message by remember { mutableStateOf("") }

    LaunchedEffect(Unit) {
        chatViewModel.loadUserInfo()
    }

    LaunchedEffect(userInfo) {
        userInfo?.let {
            userName = it.resolvedUserName()
            introduction = it.resolvedIntroduction()
            birthYear = if (it.resolvedBirthdayYear() > 0) it.resolvedBirthdayYear().toString() else ""
            birthMonth = if (it.resolvedBirthdayMonth() > 0) it.resolvedBirthdayMonth().toString() else ""
            birthDay = if (it.resolvedBirthdayDay() > 0) it.resolvedBirthdayDay().toString() else ""
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("个人中心") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "返回")
                    }
                },
                actions = {
                    TextButton(onClick = {
                        if (isEditing) {
                            val request = UpdateUserInfoRequest(
                                user_name = userName,
                                introduction = introduction,
                                birthday_year = birthYear.toIntOrNull() ?: 0,
                                birthday_month = birthMonth.toIntOrNull() ?: 0,
                                birthday_day = birthDay.toIntOrNull() ?: 0
                            )
                            chatViewModel.updateUserInfo(request) { success ->
                                message = if (success) "保存成功" else "保存失败"
                                isEditing = false
                            }
                        } else {
                            isEditing = true
                        }
                    }) {
                        Text(if (isEditing) "保存" else "编辑")
                    }
                }
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp)
        ) {
            OutlinedTextField(
                value = com.example.aim.data.TokenStore.userId,
                onValueChange = {},
                label = { Text("用户ID") },
                modifier = Modifier.fillMaxWidth(),
                readOnly = true,
                enabled = false
            )

            Spacer(modifier = Modifier.height(12.dp))

            OutlinedTextField(
                value = userName,
                onValueChange = { userName = it },
                label = { Text("昵称") },
                modifier = Modifier.fillMaxWidth(),
                readOnly = !isEditing,
                enabled = isEditing
            )

            Spacer(modifier = Modifier.height(12.dp))

            OutlinedTextField(
                value = introduction,
                onValueChange = { introduction = it },
                label = { Text("个性签名") },
                modifier = Modifier.fillMaxWidth(),
                readOnly = !isEditing,
                enabled = isEditing,
                minLines = 2
            )

            Spacer(modifier = Modifier.height(12.dp))

            Text("生日", style = MaterialTheme.typography.labelMedium)
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                OutlinedTextField(
                    value = birthYear,
                    onValueChange = { birthYear = it },
                    label = { Text("年") },
                    modifier = Modifier.weight(1f),
                    readOnly = !isEditing,
                    enabled = isEditing,
                    singleLine = true
                )
                OutlinedTextField(
                    value = birthMonth,
                    onValueChange = { birthMonth = it },
                    label = { Text("月") },
                    modifier = Modifier.weight(1f),
                    readOnly = !isEditing,
                    enabled = isEditing,
                    singleLine = true
                )
                OutlinedTextField(
                    value = birthDay,
                    onValueChange = { birthDay = it },
                    label = { Text("日") },
                    modifier = Modifier.weight(1f),
                    readOnly = !isEditing,
                    enabled = isEditing,
                    singleLine = true
                )
            }

            if (message.isNotEmpty()) {
                Spacer(modifier = Modifier.height(12.dp))
                Text(text = message, color = MaterialTheme.colorScheme.primary)
            }
        }
    }
}
