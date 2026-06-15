package com.example.aim.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.example.aim.viewmodel.AuthViewModel
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

@Composable
fun RegisterScreen(
    onRegisterSuccess: () -> Unit,
    onBack: () -> Unit,
    authViewModel: AuthViewModel = viewModel()
) {
    var password by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var message by remember { mutableStateOf("") }
    var isSuccess by remember { mutableStateOf(false) }
    val uiState by authViewModel.uiState.collectAsState()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(32.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "注册",
            style = MaterialTheme.typography.headlineLarge,
            modifier = Modifier.padding(bottom = 32.dp)
        )

        OutlinedTextField(
            value = password,
            onValueChange = { password = it },
            label = { Text("密码") },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            visualTransformation = PasswordVisualTransformation()
        )

        Spacer(modifier = Modifier.height(12.dp))

        OutlinedTextField(
            value = confirmPassword,
            onValueChange = { confirmPassword = it },
            label = { Text("确认密码") },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            visualTransformation = PasswordVisualTransformation()
        )

        if (message.isNotEmpty()) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = message,
                color = if (isSuccess) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.error
            )
        }

        Spacer(modifier = Modifier.height(24.dp))

        Button(
            onClick = {
                if (password.isBlank()) {
                    message = "请输入密码"
                    isSuccess = false
                    return@Button
                }
                if (password != confirmPassword) {
                    message = "两次密码不一致"
                    isSuccess = false
                    return@Button
                }
                authViewModel.register(password) { success, msg ->
                    message = msg
                    isSuccess = success
                    if (success) {
                        CoroutineScope(Dispatchers.Main).launch {
                            delay(1500)
                            onRegisterSuccess()
                        }
                    }
                }
            },
            modifier = Modifier.fillMaxWidth(),
            enabled = !uiState.isLoading
        ) {
            if (uiState.isLoading) {
                CircularProgressIndicator(modifier = Modifier.size(20.dp))
            } else {
                Text("注册")
            }
        }

        Spacer(modifier = Modifier.height(12.dp))

        TextButton(onClick = onBack) {
            Text("返回登录")
        }
    }
}
