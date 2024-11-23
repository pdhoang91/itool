// Models/NewsItem.swift
import Foundation

struct NewsItem: Identifiable, Codable {
    let id: Int
    let title: String
    let summary: String
    let content: String
    let sourceId: Int
    let categoryId: Int
    let imageUrl: String
    let audioUrl: String?
    let isBookmarked: Bool
    let publishedDate: Date
    let createdAt: Date
    let updatedAt: Date
    
    // Computed property để hiển thị source
    var source: String {
        // Tạm thời return string cứng, sau này có thể map với sourceId
        return "News Source"
    }
}
