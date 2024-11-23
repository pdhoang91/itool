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
    //private let baseURL = "https://apiinsight.site/api/api/"
    private let baseURL = "http://localhost:85/api"
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
            
            // Debug print
            if let jsonString = String(data: data, encoding: .utf8) {
                print("Response JSON: \(jsonString)")
            }
            
            let decoder = JSONDecoder()
            
            // Cấu hình custom date formatter
            let dateFormatter = DateFormatter()
            dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSZ"
            dateFormatter.locale = Locale(identifier: "en_US_POSIX")
            dateFormatter.timeZone = TimeZone(secondsFromGMT: 0)
            
            decoder.dateDecodingStrategy = .custom { decoder in
                let container = try decoder.singleValueContainer()
                let dateString = try container.decode(String.self)
                
                // Thử nhiều format date khác nhau
                let formats = [
                    "yyyy-MM-dd'T'HH:mm:ss.SSSZ",
                    "yyyy-MM-dd'T'HH:mm:ssZ",
                    "yyyy-MM-dd'T'HH:mm:ss.SSSSSSSZ"
                ]
                
                for format in formats {
                    dateFormatter.dateFormat = format
                    if let date = dateFormatter.date(from: dateString) {
                        return date
                    }
                }
                
                throw DecodingError.dataCorruptedError(
                    in: container,
                    debugDescription: "Cannot decode date string: \(dateString)"
                )
            }
            
            decoder.keyDecodingStrategy = .convertFromSnakeCase
            
            return try decoder.decode(NewsResponse.self, from: data)
        } catch let error as DecodingError {
            print("Decoding error: \(error)")
            throw APIError.decodingError(error)
        } catch {
            print("Network error: \(error)")
            throw APIError.networkError(error)
        }
    }
}
