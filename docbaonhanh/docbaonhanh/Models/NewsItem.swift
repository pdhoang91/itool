// Models/NewsItem.swift
// import Foundation

// struct NewsItem: Identifiable {
//     let id: UUID = UUID()
//     let title: String
//     let summary: String
//     let content: String
//     let source: String
//     let publishedDate: Date
//     let imageUrl: String
//     var isBookmarked: Bool = false
//     var audioUrl: String? // URL for cached audio file
// }

// Models/NewsItem.swift
import Foundation
struct NewsItem: Identifiable, Codable {
    let id: UUID
    let title: String
    let summary: String
    let content: String
    let source: String
    let publishedDate: Date
    let imageUrl: String
    var isBookmarked: Bool = false
    var audioUrl: String?
    
    enum CodingKeys: String, CodingKey {
        case id, title, summary, content, source
        case publishedDate = "published_date"
        case imageUrl = "image_url"
        case isBookmarked = "is_bookmarked"
        case audioUrl = "audio_url"
    }
}
