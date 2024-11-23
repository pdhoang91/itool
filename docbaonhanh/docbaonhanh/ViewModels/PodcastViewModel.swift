// ViewModels/PodcastViewModel.swift
import Foundation
import Combine

@MainActor
class PodcastViewModel: ObservableObject {
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
                let response = try await APIService.shared.fetchNews(page: 1)
                newsItems = response.items
                hasMorePages = response.hasNext
                currentPage = 1
            } catch {
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
