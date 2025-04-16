//
//  ChatViewModel.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/15/25.
//

import Foundation
import Observation

@Observable class ChatMessageViewModel: ObservableObject {
    var chat: Chat = Chat(id: "", users: [], messages: [], active: 0)
    var chats: [Chat] = []
    var message: Message = Message(id: "",  sender: "", text: "", media: [], date: 0)
    var messages: [Message] = []
    var error: String = ""
    var isLoading: Bool = false
    
    func createNewChat() async -> Bool {
        
        struct ChatResponse: Codable {
            var message: String
            var chatid: String
            
            private enum CodingKeys: String, CodingKey {
                case message = "message"
                case chatid = "chat.id"
            }
        }
        
        guard let url = URL(string: "\(Global.baseURL)/chats/new") else { return false }
        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            return false
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let body: [String: Any] = [
            "users": chat.users
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else { return false }
        request.httpBody = jsonData

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 201 {
                let decoder = JSONDecoder()
                let chatResponse = try decoder.decode(ChatResponse.self, from: data)
                print("Save ChatId: \(chatResponse.chatid)")
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
    
    func getAllChats() async -> Bool {
        
        struct ChatsResponse: Codable {
            let message: String
            let chats: [Chat]
        }
        
        guard let url = URL(string: "\(Global.baseURL)/chats/all") else { return false }
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
            
            if let jsonString = String(data: data, encoding: .utf8) {
                print("Raw response: \(jsonString)")
            }

            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                let decoded = try JSONDecoder().decode(ChatsResponse.self, from: data)
                self.chats = decoded.chats
                print("found \(decoded.chats)")
                print("set \(self.chats)")
                return true
            } else {
                self.error = "Error: \(response)"
                print(self.error)
                return false
            }
        } catch {
            self.error = "Error: \(error.localizedDescription)"
            print(self.error)
            return false
        }
    }
}
