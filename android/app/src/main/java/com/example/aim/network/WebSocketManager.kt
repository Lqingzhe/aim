package com.example.aim.network

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.SharedFlow
import okhttp3.*
import java.util.concurrent.TimeUnit

class WebSocketManager {
    private val client = OkHttpClient.Builder()
        .readTimeout(0, TimeUnit.MILLISECONDS)
        .build()

    private var webSocket: WebSocket? = null
    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())
    private var pingJob: Job? = null
    private var isConnected = false
    private var shouldReconnect = true
    private var reconnectAttempts = 0
    private val maxReconnectAttempts = 10

    private val _messages = MutableSharedFlow<String>(extraBufferCapacity = 64)
    val messages: SharedFlow<String> = _messages

    private val _connectionState = MutableSharedFlow<Boolean>(extraBufferCapacity = 8)
    val connectionState: SharedFlow<Boolean> = _connectionState

    fun connect(token: String) {
        if (isConnected) return
        shouldReconnect = true

        val request = Request.Builder()
            .url("ws://10.17.153.4:8080/ws?token=$token")
            .build()

        webSocket = client.newWebSocket(request, object : WebSocketListener() {
            override fun onOpen(webSocket: WebSocket, response: Response) {
                isConnected = true
                reconnectAttempts = 0
                scope.launch { _connectionState.emit(true) }
                startPing()
            }

            override fun onMessage(webSocket: WebSocket, text: String) {
                scope.launch {
                    _messages.emit(text)
                }
            }

            override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
                webSocket.close(1000, null)
                handleDisconnect()
            }

            override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
                handleDisconnect()
            }
        })
    }

    private fun handleDisconnect() {
        isConnected = false
        scope.launch { _connectionState.emit(false) }
        stopPing()
        if (shouldReconnect) scheduleReconnect()
    }

    private fun startPing() {
        pingJob?.cancel()
        pingJob = scope.launch {
            while (isActive && isConnected) {
                delay(30_000)
                webSocket?.send("""{"type":"ping"}""")
            }
        }
    }

    private fun stopPing() {
        pingJob?.cancel()
        pingJob = null
    }

    private fun scheduleReconnect() {
        if (reconnectAttempts >= maxReconnectAttempts) return
        reconnectAttempts++

        scope.launch {
            val delayMs = minOf(3000L * reconnectAttempts, 30_000L)
            delay(delayMs)

            val refreshed = ApiClient.refreshAccessToken()
            if (refreshed) {
                val token = ApiClient.getAccessToken() ?: return@launch
                connect(token)
            }
        }
    }

    fun send(text: String) {
        webSocket?.send(text)
    }

    fun disconnect() {
        shouldReconnect = false
        reconnectAttempts = maxReconnectAttempts
        stopPing()
        webSocket?.close(1000, "Client disconnect")
        isConnected = false
    }

    fun destroy() {
        disconnect()
        scope.cancel()
    }
}
