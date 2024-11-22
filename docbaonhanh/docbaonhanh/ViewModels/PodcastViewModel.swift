// ViewModels/PodcastViewModel.swift
import Foundation
import SwiftUI

@MainActor
class PodcastViewModel: ObservableObject {
    @Published var newsItems: [NewsItem] = []
    @Published var isLoading = false
    @Published var error: Error?
    @Published var currentPlayingItem: NewsItem?
    @Published var isPlaying = false
    
    private let apiService = APIService.shared
    
    func fetchPodcasts() {
        Task {
            isLoading = true
            do {
                let response = try await apiService.fetchNews(page: 1)
                newsItems = response.items
            } catch {
                self.error = error
            }
            isLoading = false
        }
    }
    
    func playPodcast(_ item: NewsItem) {
        currentPlayingItem = item
        isPlaying = true
        // Implement actual audio playback here
    }
    
    func pausePodcast() {
        isPlaying = false
    }
}
