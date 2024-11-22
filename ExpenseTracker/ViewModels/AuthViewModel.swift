// ViewModels/AuthViewModel.swift
@MainActor
class AuthViewModel: ObservableObject {
    @Published var isAuthenticated = false
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    func login(email: String, password: String) async throws {
        isLoading = true
        errorMessage = nil
        
        do {
            let credentials = ["email": email, "password": password]
            let jsonData = try JSONEncoder().encode(credentials)
            
            let response: AuthResponse = try await APIService.shared.request(
                endpoint: "/auth/login",
                method: "POST",
                body: jsonData
            )
            
            KeychainService.shared.saveToken(response.token)
            isAuthenticated = true
        } catch {
            errorMessage = error.localizedDescription
            throw error
        } finally {
            isLoading = false
        }
    }
}