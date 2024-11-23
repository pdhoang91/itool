// ViewModels/NewsViewModel.swift
// import Foundation

// class NewsViewModel: ObservableObject {
//     @Published var newsItems: [NewsItem] = []
    
//     func fetchNews() {
//         // Mock data
//         newsItems = [
//             NewsItem(
//                 title: "Tin tức 1",
//                 summary: "Tóm tắt tin tức 1",
//                 content: "Nội dung chi tiết tin tức 1...",
//                 source: "VnExpress",
//                 publishedDate: Date(),
//                 imageUrl: "https://example.com/image1.jpg"
//             ),
//             // Add more mock items
//         ]
//     }
// }

// ViewModels/NewsViewModel.swift
import Foundation
import SwiftUI

@MainActor
class NewsViewModel: ObservableObject {
    @Published var newsItems: [NewsItem] = []
    @Published var isLoading = false
    @Published var error: Error?
    @Published var hasMorePages = true
    
    private var currentPage = 1
    private var isFetching = false
    
func fetchInitialNews() {
    Task {
        do {
            isLoading = true
            error = nil
            print("Fetching initial news...")
            let response = try await APIService.shared.fetchNews(page: 1)
            print("Received \(response.items.count) items")
            newsItems = response.items
            hasMorePages = response.hasNext
            currentPage = 1
        } catch {
            print("Error fetching news: \(error)")
            self.error = error
        }
        isLoading = false
    }
}
    
    func loadMoreIfNeeded(currentItem item: NewsItem?) {
        guard let item = item else { return }
        
        let thresholdIndex = newsItems.index(newsItems.endIndex, offsetBy: -5)
        if newsItems.firstIndex(where: { $0.id == item.id }) == thresholdIndex {
            loadNextPage()
        }
    }
    
    private func loadNextPage() {
        guard !isFetching && hasMorePages && !isLoading else { return }
        
        Task {
            do {
                isFetching = true
                let nextPage = currentPage + 1
                let response = try await APIService.shared.fetchNews(page: nextPage)
                
                newsItems.append(contentsOf: response.items)
                hasMorePages = response.hasNext
                currentPage = nextPage
                isFetching = false
            } catch {
                self.error = error
                isFetching = false
            }
        }
    }
}
