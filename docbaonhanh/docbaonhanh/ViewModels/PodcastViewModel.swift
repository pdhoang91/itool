// ViewModels/PodcastViewModel.swift
import Foundation

class PodcastViewModel: ObservableObject {
    @Published var newsItems: [NewsItem] = []
    
    init() {
        // In a real app, you would fetch news items that have audio available
        fetchNewsItems()
    }
    
    private func fetchNewsItems() {
        // Mock data
        newsItems = [
            NewsItem(
                title: "Podcast 1",
                summary: "Tóm tắt podcast 1",
                content: "Nội dung chi tiết podcast 1...",
                source: "VnExpress",
                publishedDate: Date(),
                imageUrl: "https://example.com/image1.jpg",
                audioUrl: "https://example.com/audio1.mp3"
            ),
            // Add more mock items
        ]
    }
}