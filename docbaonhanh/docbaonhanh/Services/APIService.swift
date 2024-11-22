// Services/APIService.swift
import Foundation

enum APIError: Error {
    case invalidURL
    case invalidResponse
    case httpError(Int)
    case decodingError(Error)
    case networkError(Error)
}

class APIService {
    static let shared = APIService()
    private let baseURL = "https://your-api-base-url.com/api"
    private let pageSize = 20
    
    private init() {}
    
    func fetchNews(page: Int) async throws -> NewsResponse {
        guard let url = URL(string: "\(baseURL)/news?page=\(page)&size=\(pageSize)") else {
            throw APIError.invalidURL
        }
        
        do {
            let (data, response) = try await URLSession.shared.data(from: url)
            
            guard let httpResponse = response as? HTTPURLResponse else {
                throw APIError.invalidResponse
            }
            
            guard (200...299).contains(httpResponse.statusCode) else {
                throw APIError.httpError(httpResponse.statusCode)
            }
            
            let decoder = JSONDecoder()
            decoder.dateDecodingStrategy = .iso8601
            
            return try decoder.decode(NewsResponse.self, from: data)
        } catch let error as DecodingError {
            throw APIError.decodingError(error)
        } catch {
            throw APIError.networkError(error)
        }
    }
}
