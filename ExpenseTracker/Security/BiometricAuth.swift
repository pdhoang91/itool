// Security/BiometricAuth.swift
import LocalAuthentication

class BiometricAuth {
    static let shared = BiometricAuth()
    
    enum BiometricType {
        case none
        case touchID
        case faceID
    }
    
    func authenticate() async throws {
        let context = LAContext()
        var error: NSError?
        
        guard context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error) else {
            throw error ?? LAError(.biometryNotAvailable)
        }
        
        return try await withCheckedThrowingContinuation { continuation in
            context.evaluatePolicy(
                .deviceOwnerAuthenticationWithBiometrics,
                localizedReason: "Authenticate to access your expenses"
            ) { success, error in
                if success {
                    continuation.resume()
                } else {
                    continuation.resume(throwing: error ?? LAError(.authenticationFailed))
                }
            }
        }
    }
}
