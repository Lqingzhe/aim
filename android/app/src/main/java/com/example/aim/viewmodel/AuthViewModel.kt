package com.example.aim.viewmodel

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.aim.network.ApiClient
import com.example.aim.network.ApiService
import com.example.aim.network.WebSocketManager
import com.example.aim.network.models.*
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

class AuthViewModel : ViewModel() {
    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState

    fun login(userId: String, password: String, onResult: (Boolean, String) -> Unit) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                val response = ApiService.login(userId, password)
                if (response.code == 0) {
                    val tokenInfo = response.data?.token_info
                    if (tokenInfo != null) {
                        ApiClient.setTokens(tokenInfo.access_token, tokenInfo.refresh_token)
                        ApiClient.setDeviceId(com.example.aim.data.TokenStore.deviceId)
                        Log.i("id_login",com.example.aim.data.TokenStore.deviceId)
                        com.example.aim.data.TokenStore.saveTokens(tokenInfo.access_token, tokenInfo.refresh_token)
                        com.example.aim.data.TokenStore.userId = userId
                        onResult(true, "登录成功")
                    } else {
                        onResult(false, "登录失败: 无效响应")
                    }
                } else {
                    onResult(false, response.message.ifEmpty { "登录失败" })
                }
            } catch (e: Exception) {
                onResult(false, "网络错误: ${e.message}")
            } finally {
                _uiState.value = _uiState.value.copy(isLoading = false)
            }
        }
    }

    fun register(password: String, onResult: (Boolean, String) -> Unit) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                val response = ApiService.register(password)
                if (response.code == 0) {
                    val userId = response.data?.user_info?.user_id ?: ""
                    onResult(true, "注册成功，你的ID: $userId")
                } else {
                    onResult(false, response.message.ifEmpty { "注册失败" })
                }
            } catch (e: Exception) {
                onResult(false, "网络错误: ${e.message}")
            } finally {
                _uiState.value = _uiState.value.copy(isLoading = false)
            }
        }
    }
}

data class AuthUiState(
    val isLoading: Boolean = false
)
