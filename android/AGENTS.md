# AGENTS.md

## Project Overview

Android chat app (Jetpack Compose) client for a Go microservice IM backend.
- **Package**: `com.example.aim`
- **Min SDK**: 26, Target SDK: 36, Compile SDK: 36
- **Language**: Kotlin 2.2.10, AGP 9.2.1, Compose BOM 2026.02.01

## Build & Run

```bash
./gradlew assembleDebug          # Build debug APK
./gradlew installDebug           # Build + install on connected device
./gradlew test                   # Run unit tests
./gradlew connectedAndroidTest   # Run instrumented tests
```

**JAVA_HOME required**: Set to `D:\JetBrains\Toolbox\PyCharm\jbr` (or any JDK 25+ installation).

Gradle daemon is enabled with config cache. JVM args: `-Xmx2048m`.

## Backend API

The backend is documented in `README (1).md`. Key points:
- **Base URL**: `http://10.0.2.2:8080` (localhost for Android emulator)
- **Auth**: JWT `access_token` + `refresh_token` dual-token mechanism
- **Headers for protected routes**: `Authorization: Bearer {access_token}`, `X-Device-ID: {device_id}`
- **WebSocket**: `ws://10.0.2.2:8080/ws?token={access_token}` — real-time push, heartbeat ping/pong
- **Response format**: JSON `{"code": 0, "message": "success", "data": {...}}`

### API Modules

| Module | Prefix | Purpose |
|--------|--------|---------|
| User | `/user/*` | Register, login, profile, logout |
| Group | `/group/*` | Group CRUD, membership, friend/session mgmt |
| Message | `/message/*` | Send text/file/voice/image, history, withdraw |
| AI | `/ai/*` | AI config, chat context |
| WebSocket | `/ws` | Real-time message push |

## Project Structure

```
app/src/main/java/com/example/aim/
├── MainActivity.kt                    # Entry point, initializes TokenStore & navigation
├── data/
│   └── TokenStore.kt                  # SharedPreferences for tokens & user ID
├── network/
│   ├── ApiClient.kt                   # OkHttp REST client (JSON serialization)
│   ├── ApiService.kt                  # High-level API functions per module
│   ├── WebSocketManager.kt            # WebSocket connection with ping/reconnect
│   └── models/
│       └── Models.kt                  # All request/response data classes
├── viewmodel/
│   ├── AuthViewModel.kt               # Login/register logic
│   └── ChatViewModel.kt               # Chat, contacts, groups, profile logic
└── ui/
    ├── navigation/
    │   └── NavGraph.kt                # Compose Navigation routes
    ├── screens/
    │   ├── LoginScreen.kt             # Login form
    │   ├── RegisterScreen.kt          # Register form
    │   ├── HomeScreen.kt              # Tab bar: Chats / Contacts / Settings
    │   ├── ChatRoomScreen.kt          # Message list + input
    │   └── ProfileScreen.kt           # User info edit
    └── theme/                         # Material3 theme (dynamic color)
```

## Conventions

- Use Jetpack Compose (Material3) for all UI — no XML layouts
- Follow existing package structure under `com.example.aim`
- Compose BOM manages all Compose dependency versions
- Network calls run in `viewModelScope` or `Dispatchers.IO`
- WebSocket reconnects automatically with 3s delay
- Token store uses `SharedPreferences` (not DataStore)

## Key Dependencies

| Library | Purpose |
|---------|---------|
| OkHttp | HTTP client + multipart uploads |
| kotlinx.serialization | JSON parsing |
| Navigation Compose | Screen routing |
| Coil | Image loading (future use) |
| Material Icons Extended | Full icon set |

## Testing

- Unit tests: `app/src/test/java/com/example/aim/`
- Instrumented tests: `app/src/androidTest/java/com/example/aim/`
- Framework: JUnit4 + Espresso + Compose testing
