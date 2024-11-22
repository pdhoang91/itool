// ViewModels/NewsViewModel.swift
import Foundation

class NewsViewModel: ObservableObject {
    @Published var newsItems: [NewsItem] = []
    
    func fetchNews() {
        // Mock data
        newsItems = [
            NewsItem(
                title: "Tin tức 1",
                summary: "Tóm tắt tin tức 1",
                content: "Nội dung chi tiết tin tức 1...",
                source: "VnExpress",
                publishedDate: Date(),
                imageUrl: "https://example.com/image1.jpg"
            ),
            // Add more mock items
        ]
    }
}