// Security/EncryptionManager.swift
import CryptoKit

class EncryptionManager {
    static let shared = EncryptionManager()
    
    private let keychain = KeychainManager.shared
    
    func encryptData(_ data: Data) throws -> Data {
        let key = try getOrCreateKey()
        let sealedBox = try AES.GCM.seal(data, using: key)
        return sealedBox.combined!
    }
    
    func decryptData(_ data: Data) throws -> Data {
        let key = try getOrCreateKey()
        let sealedBox = try AES.GCM.SealedBox(combined: data)
        return try AES.GCM.open(sealedBox, using: key)
    }
    
    private func getOrCreateKey() throws -> SymmetricKey {
        if let existingKey = try? keychain.getKey() {
            return existingKey
        }
        
        let key = SymmetricKey(size: .bits256)
        try keychain.saveKey(key)
        return key
    }
}