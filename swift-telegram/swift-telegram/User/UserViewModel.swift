//
//  UserViewModel.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/14/25.
//

import Foundation
import Observation

@Observable class UserViewModel: ObservableObject {
    var user: User = User(id: "", name: "", email: "", password: "")
    var users: [User] = []
    var error: String = ""
    var isLoading: Bool = false

    func createNewUser() async -> Bool {
        guard let url = URL(string: "\(Global.baseURL)/users/new") else { return false }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "name": user.name,
            "email": user.email,
            "password": user.password!
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else { return false }
        request.httpBody = jsonData

        do {
            let (_, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 201 {
                return true
            } else {
                self.error = "Error: \(response)"
                return false
            }
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }

    func Login() async -> Bool {
        
        struct LoginResponse: Codable {
            let message: String
            let refreshToken: String
            let token: String
            let user: User

            private enum CodingKeys: String, CodingKey {
                case message
                case refreshToken = "refresh_token"
                case token
                case user
            }
        }
        
        guard let url = URL(string: "\(Global.baseURL)/users/login") else { return false }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "email": user.email,
            "password": user.password ?? ""
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else { return false }
        request.httpBody = jsonData

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                let decoder = JSONDecoder()
                decoder.dateDecodingStrategy = .iso8601
                let loginResponse = try decoder.decode(LoginResponse.self, from: data)

                UserDefaults.standard.setValue(loginResponse.token, forKey: "authToken")
                UserDefaults.standard.setValue(loginResponse.refreshToken, forKey: "refresh_token")

                if let encodedUser = try? JSONEncoder().encode(loginResponse.user) {
                    UserDefaults.standard.setValue(encodedUser, forKey: "currentUser")
                }
                return true
            } else {
                self.error = "Error: \(response)"
                return false
            }
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }
    
    func getAllUsers() async -> Bool {
        
        struct UsersResponse: Codable {
            let message: String
            let users: [User]
        }
        
        guard let url = URL(string: "\(Global.baseURL)/users/all") else { return false }
        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                let decoded = try JSONDecoder().decode(UsersResponse.self, from: data)
                self.users = decoded.users
                return true
            } else {
                self.error = "Error: \(response)"
                return false
            }
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }

    func getUserById(id: String) async -> Bool {
        
        guard let url = URL(string: "\(Global.baseURL)/users/id/\(id)") else { return false }
        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            return false
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let fetchedUser = try JSONDecoder().decode(User.self, from: data)
            self.user = fetchedUser
            return true
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }

    func updateUser(id: String) async -> Bool {
        
        guard let url = URL(string: "\(Global.baseURL)/users/update") else { return false }
        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            print(self.error)
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let body: [String: Any] = [
            "id": user.id,
            "name": user.name,
            "email": user.email,
            "password": user.password!
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else { return false }
        request.httpBody = jsonData

        do {
            let (_, response) = try await URLSession.shared.data(for: request)
            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                return true
            } else {
                self.error = "Error: \(response)"
                return false
            }
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }

    func deleteUser(id: String) async -> Bool {
        
        guard let url = URL(string: "\(Global.baseURL)/users/\(id)") else { return false }
        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            print(self.error)
            return false
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        do {
            let (_, response) = try await URLSession.shared.data(for: request)
            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                return true
            } else {
                self.error = "Error: \(response)"
                return false
            }
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }
}
