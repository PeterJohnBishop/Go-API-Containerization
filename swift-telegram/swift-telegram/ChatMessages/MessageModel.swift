//
//  MessageModel.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/15/25.
//

import Foundation

struct Message: Codable, Identifiable, Equatable {
    var id: String
    var sender: String
    var text: String
    var media: [String]
    var date: Double
    
    func encode() throws -> Data {
            let encoder = JSONEncoder()
            return try encoder.encode(self)
        }

    static func decode(from data: Data) throws -> Message {
            let decoder = JSONDecoder()
        return try decoder.decode(Message.self, from: data)
        }
}
