// Models/NewsItem.swift
import Foundation

struct NewsItem: Identifiable {
    let id: UUID = UUID()
    let title: String
    let summary: String
    let content: String
    let source: String
    let publishedDate: Date
    let imageUrl: String
    var isBookmarked: Bool = false
    var audioUrl: String? // URL for cached audio file
}