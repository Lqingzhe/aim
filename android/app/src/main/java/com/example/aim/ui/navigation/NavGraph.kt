package com.example.aim.ui.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.navArgument
import com.example.aim.ui.screens.*
import com.example.aim.viewmodel.ChatViewModel

sealed class Screen(val route: String) {
    data object Login : Screen("login")
    data object Register : Screen("register")
    data object Home : Screen("home")
    data object Profile : Screen("profile")
    data object ChatRoom : Screen("chat/{chatId}/{chatType}") {
        fun createRoute(chatId: String, chatType: String) = "chat/$chatId/$chatType"
    }
}

@Composable
fun AppNavigation(
    navController: NavHostController,
    startDestination: String = Screen.Login.route,
    chatViewModel: ChatViewModel,
    onLogout: () -> Unit
) {
    NavHost(navController = navController, startDestination = startDestination) {
        composable(Screen.Login.route) {
            LoginScreen(
                onLoginSuccess = {
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.Login.route) { inclusive = true }
                    }
                },
                onNavigateToRegister = {
                    navController.navigate(Screen.Register.route)
                }
            )
        }

        composable(Screen.Register.route) {
            RegisterScreen(
                onRegisterSuccess = {
                    navController.popBackStack()
                },
                onBack = {
                    navController.popBackStack()
                }
            )
        }

        composable(Screen.Home.route) {
            HomeScreen(
                onOpenChat = { chatId, chatType ->
                    navController.navigate(Screen.ChatRoom.createRoute(chatId, chatType))
                },
                onOpenProfile = {
                    navController.navigate(Screen.Profile.route)
                },
                onLogout = onLogout,
                chatViewModel = chatViewModel
            )
        }

        composable(Screen.Profile.route) {
            ProfileScreen(
                onBack = { navController.popBackStack() },
                chatViewModel = chatViewModel
            )
        }

        composable(
            route = Screen.ChatRoom.route,
            arguments = listOf(
                navArgument("chatId") { type = NavType.StringType },
                navArgument("chatType") { type = NavType.StringType }
            )
        ) { backStackEntry ->
            val chatId = backStackEntry.arguments?.getString("chatId") ?: ""
            val chatType = backStackEntry.arguments?.getString("chatType") ?: "group"
            ChatRoomScreen(
                chatId = chatId,
                chatType = chatType,
                onBack = { navController.popBackStack() },
                chatViewModel = chatViewModel
            )
        }
    }
}
