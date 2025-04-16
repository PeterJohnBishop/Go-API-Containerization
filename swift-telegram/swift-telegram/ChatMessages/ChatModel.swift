//
//  ChatModel.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/15/25.
//

import Foundation
import Observation

struct Chat: Codable, Identifiable, Equatable {
    var id: String
    var users: [String]
    var messages: [String]?
    var active: Double
    
    func encode() throws -> Data {
            let encoder = JSONEncoder()
            return try encoder.encode(self)
        }

    static func decode(from data: Data) throws -> Chat {
            let decoder = JSONDecoder()
            return try decoder.decode(Chat.self, from: data)
        }
}
