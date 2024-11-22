// Services/APIService.swift
enum APIError: Error {
    case invalidURL
    case noData
    case decodingError
    case unauthorized
    case serverError(String)
}

class APIService {
    static let shared = APIService()
    private let baseURL = "https://your-api-url.com/api"
    private let session = URLSession.shared
    
    func request<T: Codable>(
        endpoint: String,
        method: String = "GET",
        body: Data? = nil
    ) async throws -> T {
        guard let url = URL(string: baseURL + endpoint) else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if let token = KeychainService.shared.getToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        if let body = body {
            request.httpBody = body
        }
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.serverError("Invalid response")
        }
        
        switch httpResponse.statusCode {
        case 200...299:
            do {
                return try JSONDecoder().decode(T.self, from: data)
            } catch {
                throw APIError.decodingError
            }
        case 401:
            throw APIError.unauthorized
        default:
            throw APIError.serverError("Server error: \(httpResponse.statusCode)")
        }
    }
}