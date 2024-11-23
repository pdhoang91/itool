// Views/PodcastView.swift
import SwiftUI
import AVFoundation

struct PodcastView: View {
    @StateObject private var viewModel = PodcastViewModel()
    @ObservedObject private var audioService = AudioPlayerService.shared
    @State private var showingFullPlayer = false
    
    var body: some View {
        NavigationView {
            ZStack {
                List {
                    ForEach(viewModel.newsItems) { item in
                        PodcastItemRow(
                            item: item,
                            isPlaying: audioService.isPlaying &&
                                     audioService.currentIndex == viewModel.newsItems.firstIndex(where: { $0.id == item.id }) ?? -1,
                            onTap: {
                                if let index = viewModel.newsItems.firstIndex(where: { $0.id == item.id }) {
                                    audioService.setPlaylist(viewModel.newsItems)
                                    audioService.currentIndex = index
                                    if let audioUrl = URL(string: item.audioUrl ?? "") {
                                        audioService.play(url: audioUrl)
                                        showingFullPlayer = true
                                    }
                                }
                            }
                        )
                        .onAppear {
                            viewModel.loadMoreIfNeeded(currentItem: item)
                        }
                    }
                }
                .refreshable {
                    viewModel.fetchInitialNews()
                }
                
                if viewModel.isLoading && viewModel.newsItems.isEmpty {
                    ProgressView()
                }
                
                if let error = viewModel.error {
                    VStack {
                        Text("Có lỗi xảy ra")
                            .font(.headline)
                        Text(error.localizedDescription)
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                        Button("Thử lại") {
                            viewModel.fetchInitialNews()
                        }
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    .shadow(radius: 4)
                }
            }
            .navigationTitle("Podcast")
            .sheet(isPresented: $showingFullPlayer) {
                if let currentItem = viewModel.newsItems[safe: audioService.currentIndex] {
                    AudioPlayerView(audioService: audioService, newsTitle: currentItem.title)
                        .presentationDetents([.height(200)])
                }
            }
            .overlay(
                VStack {
                    Spacer()
                    if !viewModel.newsItems.isEmpty &&
                       audioService.isPlaying,
                       let currentItem = viewModel.newsItems[safe: audioService.currentIndex] {
                        MiniPlayerView(audioService: audioService, newsTitle: currentItem.title)
                            .onTapGesture {
                                showingFullPlayer = true
                            }
                    }
                }
            )
        }
        .onAppear {
            if viewModel.newsItems.isEmpty {
                viewModel.fetchInitialNews()
            }
        }
    }
}

// Helper extension để an toàn khi truy cập array
extension Array {
    subscript(safe index: Int) -> Element? {
        return indices.contains(index) ? self[index] : nil
    }
}

struct PodcastItemRow: View {
    let item: NewsItem
    let isPlaying: Bool
    let onTap: () -> Void
    
    var body: some View {
        HStack(spacing: 12) {
            Button(action: onTap) {
                Image(systemName: isPlaying ? "pause.circle.fill" : "play.circle.fill")
                    .font(.title)
                    .foregroundColor(.blue)
            }
            
            VStack(alignment: .leading, spacing: 4) {
                Text(item.title)
                    .font(.headline)
                    .lineLimit(2)
                
                Text(item.source)
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            
            Spacer()
            
            Text(item.publishedDate, style: .date)
                .font(.caption)
                .foregroundColor(.gray)
        }
        .padding(.vertical, 8)
    }
}

#Preview {
    NewsView()
}

