//
//  UserModel.swift
//  swift-telegram
//
//  Created by Peter Bishop on 4/14/25.
//

import Foundation
import Observation

struct User: Codable, Equatable {
    var id: String
    var name: String
    var email: String
    var password: String
    
    func encode() throws -> Data {
            let encoder = JSONEncoder()
            return try encoder.encode(self)
        }

    static func decode(from data: Data) throws -> User {
            let decoder = JSONDecoder()
            return try decoder.decode(User.self, from: data)
        }
}
