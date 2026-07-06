package com.example.aim

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.compose.rememberNavController
import com.example.aim.data.TokenStore
import com.example.aim.network.ApiClient
import com.example.aim.network.ApiService
import com.example.aim.network.WebSocketManager
import com.example.aim.ui.navigation.AppNavigation
import com.example.aim.ui.navigation.Screen
import com.example.aim.ui.theme.AimTheme
import com.example.aim.viewmodel.ChatViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.MainScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : ComponentActivity() {
    private val webSocketManager = WebSocketManager()
    private val scope = MainScope()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        TokenStore.init(this)
        ApiClient.setDeviceId(TokenStore.deviceId)

        setContent {
            AimTheme {
                val chatViewModel: ChatViewModel = ViewModelProvider(this@MainActivity)[ChatViewModel::class.java]
                var isLoggedIn by remember { mutableStateOf(TokenStore.isLoggedIn) }
                var isChecking by remember { mutableStateOf(TokenStore.isLoggedIn) }

                if (TokenStore.isLoggedIn) {
                    LaunchedEffect(Unit) {
                        val valid = checkAndRefreshToken()
                        isChecking = false
                        if (valid) {
                            isLoggedIn = true
                            connectWebSocket(chatViewModel)
                        } else {
                            TokenStore.clearTokens()
                            isLoggedIn = false
                        }
                    }
                }

                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    if (isChecking) return@Surface

                    val navController = rememberNavController()
                    val startDestination = if (isLoggedIn) Screen.Home.route else Screen.Login.route

                    AppNavigation(
                        navController = navController,
                        startDestination = startDestination,
                        chatViewModel = chatViewModel,
                        onLogout = {
                            webSocketManager.disconnect()
                            scope.launch { ApiService.logoutCurrentDevice() }
                            TokenStore.clearTokens()
                            isLoggedIn = false
                            navController.navigate(Screen.Login.route) {
                                popUpTo(0) { inclusive = true }
                            }
                        }
                    )
                }
            }
        }
    }

    private suspend fun checkAndRefreshToken(): Boolean = withContext(Dispatchers.IO) {
        val accessToken = TokenStore.accessToken
        val refreshToken = TokenStore.refreshToken

        if (accessToken.isNullOrEmpty()) return@withContext false
        ApiClient.setTokens(accessToken, refreshToken ?: "")

        try {
            val response = ApiService.getUserInfo()
            when (response.code) {
                0 -> true
                1101, 1102 -> {
                    val refreshed = ApiClient.refreshAccessToken()
                    if (refreshed) {
                        val newAccess = ApiClient.getAccessToken()
                        if (newAccess != null) {
                            TokenStore.accessToken = newAccess
                            TokenStore.isLoggedIn = true
                            true
                        } else false
                    } else false
                }
                else -> false
            }
        } catch (_: Exception) {
            false
        }
    }

    private fun connectWebSocket(chatViewModel: ChatViewModel) {
        val token = ApiClient.getAccessToken() ?: return
        webSocketManager.connect(token)

        scope.launch {
            webSocketManager.messages.collect { message ->
                chatViewModel.handleWsMessage(message)
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        webSocketManager.destroy()
    }
}
