// Views/PodcastView.swift
import SwiftUI

struct PodcastView: View {
    @StateObject private var viewModel = PodcastViewModel()
    
    var body: some View {
        NavigationView {
            List {
                if viewModel.isLoading {
                    ProgressView()
                } else {
                    ForEach(viewModel.newsItems) { item in
                        PodcastItemRow(
                            item: item,
                            isPlaying: viewModel.currentPlayingItem?.id == item.id && viewModel.isPlaying,
                            onTap: {
                                if viewModel.currentPlayingItem?.id == item.id {
                                    viewModel.isPlaying ? viewModel.pausePodcast() : viewModel.playPodcast(item)
                                } else {
                                    viewModel.playPodcast(item)
                                }
                            }
                        )
                    }
                }
            }
            .navigationTitle("Podcast")
            .overlay(
                Group {
                    if let currentItem = viewModel.currentPlayingItem {
                        VStack {
                            Spacer()
                            MiniPlayer(
                                item: currentItem,
                                isPlaying: viewModel.isPlaying,
                                onPlayPause: {
                                    viewModel.isPlaying ? viewModel.pausePodcast() : viewModel.playPodcast(currentItem)
                                }
                            )
                        }
                    }
                }
            )
        }
        .onAppear {
            viewModel.fetchPodcasts()
        }
    }
}

struct PodcastItemRow: View {
    let item: NewsItem
    let isPlaying: Bool
    let onTap: () -> Void
    
    var body: some View {
        HStack(spacing: 12) {
            // Left side - Info
            VStack(alignment: .leading, spacing: 8) {
                Text(item.title)
                    .font(.headline)
                    .lineLimit(2)
                
                Text(item.source)
                    .font(.caption)
                    .foregroundColor(.blue)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            
            // Right side - Image and Play Button
            ZStack {
                AsyncImage(url: URL(string: item.imageUrl)) { image in
                    image
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                } placeholder: {
                    Color.gray
                }
                .frame(width: UIScreen.main.bounds.width * 0.25)
                .clipped()
                .cornerRadius(8)
                
                Button(action: onTap) {
                    Image(systemName: isPlaying ? "pause.circle.fill" : "play.circle.fill")
                        .font(.title)
                        .foregroundColor(.white)
                        .shadow(radius: 4)
                }
            }
        }
        .frame(height: 80)
        .padding(.vertical, 4)
    }
}

struct MiniPlayer: View {
    let item: NewsItem
    let isPlaying: Bool
    let onPlayPause: () -> Void
    
    var body: some View {
        HStack(spacing: 16) {
            AsyncImage(url: URL(string: item.imageUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Color.gray
            }
            .frame(width: 40, height: 40)
            .cornerRadius(4)
            
            Text(item.title)
                .lineLimit(1)
            
            Spacer()
            
            Button(action: onPlayPause) {
                Image(systemName: isPlaying ? "pause.circle.fill" : "play.circle.fill")
                    .font(.title)
            }
        }
        .padding()
        .background(Color(UIColor.systemBackground))
        .shadow(radius: 4)
    }
}

struct PodcastView_Previews: PreviewProvider {
    static var previews: some View {
        PodcastView()
    }
}

#Preview {
    NewsView()
}

