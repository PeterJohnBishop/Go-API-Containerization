//
//  SuccessView.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/14/25.
//

import SwiftUI

struct SuccessView: View {
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
                    Text(jwt ?? "No JWT")
                    Text(refreshToken ?? "No refresh token")
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
                }.onAppear() {
                    isLoading = true
                    do {
                        currentUser = try User.decode(from: UserDefaults.standard.object(forKey: "currentUser") as? Data ?? Data())
                        jwt = UserDefaults.standard.string(forKey: "authToken")
                        refreshToken = UserDefaults.standard.string(forKey: "refresh_token")
                    } catch {
                        errorLoading = true
                    }
                    isLoading = false
                }
            }
        }
    }
}

#Preview {
    SuccessView()
}
