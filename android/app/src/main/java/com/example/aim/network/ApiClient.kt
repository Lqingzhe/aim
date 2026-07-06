package com.example.aim.network

import android.util.Log
import com.example.aim.network.models.*
import kotlinx.coroutines.*
import kotlinx.serialization.json.Json
import okhttp3.*
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import java.util.concurrent.TimeUnit
import kotlin.coroutines.resume

object ApiClient {
    private const val BASE_URL = "http://10.32.216.127:8080"
    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    val json = Json {
        ignoreUnknownKeys = true
        coerceInputValues = true
        isLenient = true
    }

    private val client = OkHttpClient.Builder()
        .connectTimeout(15, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()

    private val jsonMediaType = "application/json; charset=utf-8".toMediaType()

    private var accessToken: String? = null
    private var refreshToken: String? = null
    private var deviceId: String = ""

    private var isRefreshing = false
    private val failedQueue = mutableListOf<suspend (Boolean) -> Unit>()

    fun setTokens(access: String, refresh: String) {
        accessToken = access
        refreshToken = refresh
        Log.e("token",accessToken+"||"+refreshToken)
    }

    fun setDeviceId(id: String) {
        deviceId = id
    }

    fun getAccessToken(): String? = accessToken

    private fun processQueue(success: Boolean) {
        val queue = failedQueue.toList()
        failedQueue.clear()
        queue.forEach { resume -> scope.launch { resume(success) } }
    }

    suspend fun refreshAccessToken(): Boolean {
        val rt = refreshToken ?: return false

        if (isRefreshing) {
            return suspendCancellableCoroutine<Boolean> { cont ->
                failedQueue.add { success -> if (cont.isActive) cont.resume(success) }
            }
        }

        isRefreshing = true
        return try {
            val bodyJson = json.encodeToString(RefreshTokenRequest.serializer(), RefreshTokenRequest(rt))
            val request = Request.Builder()
                .url("$BASE_URL/user/refresh-token")
                .post(bodyJson.toRequestBody(jsonMediaType))
                .header("Content-Type", "application/json")
                .header("X-Device-ID", deviceId)
                .build()
            Log.e("deviceId1=",deviceId)

            val response = client.newCall(request).execute()
            val responseBody = response.body?.string() ?: "{}"
            val apiResponse = json.decodeFromString<ApiResponse>(responseBody)

            if (apiResponse.code == 0 && apiResponse.data?.token_info != null) {
                val tokenInfo = apiResponse.data.token_info
                accessToken = tokenInfo.access_token
                refreshToken = tokenInfo.refresh_token
                com.example.aim.data.TokenStore.accessToken = tokenInfo.access_token
                com.example.aim.data.TokenStore.refreshToken = tokenInfo.refresh_token
                processQueue(true)
                true
            } else {
                processQueue(false)
                false
            }
        } catch (_: Exception) {
            processQueue(false)
            false
        } finally {
            isRefreshing = false
        }
    }

    suspend fun post(path: String, bodyJson: String, auth: Boolean = true): ApiResponse =
        withContext(Dispatchers.IO) {
            val requestBuilder = Request.Builder()
                .url("$BASE_URL$path")
                .post(bodyJson.toRequestBody(jsonMediaType))
                .header("Content-Type", "application/json")

            if (auth && accessToken != null) {
                requestBuilder.header("Authorization", "Bearer $accessToken")
                requestBuilder.header("X-Device-ID", deviceId)
                Log.e("deviceId2=",deviceId)
            }

            if (!auth) {
                requestBuilder.header("X-Device-ID", deviceId)
                Log.e("deviceId3=",deviceId)
            }

            val response = client.newCall(requestBuilder.build()).execute()
            val responseBody = response.body?.string() ?: "{}"
            val apiResponse = json.decodeFromString<ApiResponse>(responseBody)

            if (auth && (apiResponse.code == 1101 || apiResponse.code == 1102)) {
                val refreshed = refreshAccessToken()
                if (refreshed) {
                    return@withContext post(path, bodyJson, auth)
                } else {
                    clearAuth()
                }
            }

            if (auth && apiResponse.code == 1103) {
                clearAuth()
            }

            apiResponse
        }

    private fun clearAuth() {
        accessToken = null
        refreshToken = null
        com.example.aim.data.TokenStore.clearTokens()
    }

    suspend fun postLogin(path: String, bodyJson: String): ApiResponse =
        withContext(Dispatchers.IO) {
            val requestBuilder = Request.Builder()
                .url("$BASE_URL$path")
                .post(bodyJson.toRequestBody(jsonMediaType))
                .header("Content-Type", "application/json")
                .header("X-Device-ID", deviceId)

            val response = client.newCall(requestBuilder.build()).execute()
            val responseBody = response.body?.string() ?: "{}"
            json.decodeFromString<ApiResponse>(responseBody)
        }

    suspend fun postMultipart(path: String, parts: Map<String, String>, fileField: String, fileBytes: ByteArray, fileName: String): ApiResponse =
        withContext(Dispatchers.IO) {
            val bodyBuilder = MultipartBody.Builder()
                .setType(MultipartBody.FORM)
            parts.forEach { (key, value) ->
                bodyBuilder.addFormDataPart(key, value)
            }
            bodyBuilder.addFormDataPart(fileField, fileName, fileBytes.toRequestBody("application/octet-stream".toMediaType()))

            val requestBuilder = Request.Builder()
                .url("$BASE_URL$path")
                .post(bodyBuilder.build())

            if (accessToken != null) {
                requestBuilder.header("Authorization", "Bearer $accessToken")
                requestBuilder.header("X-Device-ID", deviceId)
               // Log.e("deviceId5=",deviceId)
            }

            val response = client.newCall(requestBuilder.build()).execute()
            val responseBody = response.body?.string() ?: "{}"
            val apiResponse = json.decodeFromString<ApiResponse>(responseBody)

            if (apiResponse.code == 1101 || apiResponse.code == 1102) {
                val refreshed = refreshAccessToken()
                if (refreshed) {
                    return@withContext postMultipart(path, parts, fileField, fileBytes, fileName)
                }
            }

            apiResponse
        }

    suspend fun downloadFile(path: String, bodyJson: String): ByteArray? =
        withContext(Dispatchers.IO) {
            val requestBuilder = Request.Builder()
                .url("$BASE_URL$path")
                .post(bodyJson.toRequestBody(jsonMediaType))
                .header("Content-Type", "application/json")

            if (accessToken != null) {
                requestBuilder.header("Authorization", "Bearer $accessToken")
                requestBuilder.header("X-Device-ID", deviceId)
                Log.e("deviceId6=",deviceId)
            }

            val response = client.newCall(requestBuilder.build()).execute()
            response.body?.bytes()
        }
}
