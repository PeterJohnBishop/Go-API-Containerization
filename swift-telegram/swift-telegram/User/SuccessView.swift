//
//  SuccessView.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/14/25.
//

import SwiftUI

struct SuccessView: View {
    @StateObject private var userVM = UserViewModel()
    @StateObject private var chatVM = ChatMessageViewModel()

    @State private var currentUser: User?
    @State private var jwt: String?
    @State private var refreshToken: String?
    @State private var errorLoading: Bool = false
    @State private var logout: Bool = false

    private func logoutUser() {
        UserDefaults.standard.removeObject(forKey: "currentUser")
        UserDefaults.standard.removeObject(forKey: "authToken")
        UserDefaults.standard.removeObject(forKey: "refreshToken")
        logout = true
    }

    private var logoutButton: some View {
        Button("Logout", action: logoutUser)
            .fontWeight(.ultraLight)
            .foregroundColor(.black)
            .padding()
            .background(
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color.white)
                    .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
            )
    }
    
    private func startChat(userId: String) async -> String {
        print(chatVM.chats)
        guard let currentUser = currentUser else {
            print("No current user found.")
            return ""
        }

        let targetUsersSet: Set<String> = [currentUser.id, userId]

        if let index = chatVM.chats.firstIndex(where: { chat in
            Set(chat.users) == targetUsersSet
        }) {
            print("found existing chat \(chatVM.chats[index].id)")
            return chatVM.chats[index].id
        } else {
            chatVM.chat.users = [userId, currentUser.id]
            await chatVM.createNewChat()
            print("created new chat \(chatVM.chat.id)")
            return chatVM.chat.id
        }
    }

    var body: some View {
        NavigationStack {
            VStack {
                if errorLoading {
                    Text("Failed to load user info.")
                        .foregroundColor(.red)
                        .padding()
                }

                if userVM.isLoading {
                    ProgressView("Loading...")
                        .padding()
                    logoutButton
                } else {
                    HStack {
                        Text(currentUser?.name ?? "No user")
                        Spacer()
                        logoutButton
                    }
                    .padding()

                    NavigationView {
                        Group {
                            if userVM.isLoading {
                                ProgressView("Loading Users...")
                            } else if !userVM.error.isEmpty {
                                Text("Error: \(userVM.error)")
                                    .foregroundColor(.red)
                            } else {
                                List(userVM.users) { user in
                                    if user.id != currentUser?.id {
                                        HStack {
                                            VStack(alignment: .leading) {
                                                Text(user.name)
                                                    .font(.headline)
                                                Text(user.email)
                                                    .font(.subheadline)
                                                    .foregroundColor(.secondary)
                                            }
                                            Spacer()
                                            Button {
                                                Task{
                                                    await startChat(userId: user.id)
                                                    await chatVM.getAllChats()
                                                }
                                            } label: {
                                                Image(systemName: "message")
                                            }
                                        }
                                    }
                                }
                            }
                        }
                        .navigationTitle("Users")
                    }
                }
            }
            .navigationDestination(isPresented: $logout) {
                LoginView().navigationBarBackButtonHidden(true)
            }
            .onAppear {
                do {
                    currentUser = try User.decode(from: UserDefaults.standard.object(forKey: "currentUser") as? Data ?? Data())
                    jwt = UserDefaults.standard.string(forKey: "authToken")
                    refreshToken = UserDefaults.standard.string(forKey: "refresh_token")
                    Task {
                        await userVM.getAllUsers()
                        await chatVM.getAllChats()
                    }
                } catch {
                    errorLoading = true
                }
            }
        }
    }
}

#Preview {
    SuccessView()
}
