//
//  SuccessView.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/14/25.
//

import SwiftUI

struct SuccessView: View {
    @State var userVM: UserViewModel = UserViewModel()
    @State var currentUser: User?
    @State var jwt: String?
    @State var refreshToken: String?
    @State var isLoading: Bool = false
    @State var errorLoading: Bool = false
    @State var logout: Bool = false
    
    func logoutUser() {
        UserDefaults.standard.removeObject(forKey: "currentUser")
        UserDefaults.standard.removeObject(forKey: "authToken")
        UserDefaults.standard.removeObject(forKey: "refreshToken")
        logout = true
    }
    
    var body: some View {
        NavigationStack {
            if isLoading {
                VStack{
                    ProgressView()
                    Button("Logout", action: {
                        logoutUser()
                    })
                    .fontWeight(.ultraLight)
                    .foregroundColor(.black)
                    .padding()
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(Color.white)
                            .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                    )
                    .navigationDestination(isPresented: $logout) {
                        LoginView().navigationBarBackButtonHidden(true)
                    }
                }
            } else {
                VStack {
                    Text("Success!")
                    Text(currentUser?.name ?? "No user")
                    Button("Logout", action: {
                        logoutUser()
                    })
                    .fontWeight(.ultraLight)
                    .foregroundColor(.black)
                    .padding()
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(Color.white)
                            .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                    )
                    .navigationDestination(isPresented: $logout) {
                        LoginView().navigationBarBackButtonHidden(true)
                    }
                    NavigationView {
                                Group {
                                    if userVM.isLoading {
                                        ProgressView("Loading Users...")
                                    } else if userVM.error != "" {
                                        Text("Error: \(userVM.error)")
                                            .foregroundColor(.red)
                                    } else {
                                        List(userVM.users) { user in
                                            VStack(alignment: .leading) {
                                                Text(user.name)
                                                    .font(.headline)
                                                Text(user.email)
                                                    .font(.subheadline)
                                                    .foregroundColor(.secondary)
                                            }
                                        }
                                    }
                                }
                                .navigationTitle("Users")
                    
                    }.onAppear() {
                        isLoading = true
                        do {
                            currentUser = try User.decode(from: UserDefaults.standard.object(forKey: "currentUser") as? Data ?? Data())
                            jwt = UserDefaults.standard.string(forKey: "authToken")
                            refreshToken = UserDefaults.standard.string(forKey: "refresh_token")
                            Task {
                                userVM.isLoading = true
                                await userVM.getAllUsers()
                                userVM.isLoading = false
                            }
                        } catch {
                            errorLoading = true
                        }
                        isLoading = false
                    }
                }
            }
        }
    }
}

#Preview {
    SuccessView()
}
