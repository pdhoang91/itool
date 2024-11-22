// Views/PodcastView.swift
// Views/PodcastView.swift
import SwiftUI

struct PodcastView: View {
    @StateObject private var viewModel = PodcastViewModel()
    @StateObject private var audioPlayer = AudioPlayerViewModel()
    
    var body: some View {
        NavigationView {
            List {
                ForEach(viewModel.newsItems) { item in
                    PodcastItemRow(item: item, audioPlayer: audioPlayer)
                }
            }
            .navigationTitle("Podcast")
            .overlay(
                // Mini Player if something is playing
                audioPlayer.currentItem.map { item in
                    VStack {
                        Spacer()
                        MiniPlayer(item: item, audioPlayer: audioPlayer)
                    }
                }
            )
        }
    }
}

struct PodcastItemRow: View {
    let item: NewsItem
    @ObservedObject var audioPlayer: AudioPlayerViewModel
    
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
                
                // Play/Pause Button overlay
                Button(action: {
                    if audioPlayer.currentItem?.id == item.id {
                        audioPlayer.isPlaying ? audioPlayer.pause() : audioPlayer.play(item)
                    } else {
                        audioPlayer.play(item)
                    }
                }) {
                    Image(systemName: audioPlayer.currentItem?.id == item.id && audioPlayer.isPlaying ? 
                          "pause.circle.fill" : "play.circle.fill")
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
    @ObservedObject var audioPlayer: AudioPlayerViewModel
    
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
            
            Button(action: {
                audioPlayer.isPlaying ? audioPlayer.pause() : audioPlayer.play(item)
            }) {
                Image(systemName: audioPlayer.isPlaying ? "pause.circle.fill" : "play.circle.fill")
                    .font(.title)
            }
        }
        .padding()
        .background(Color(UIColor.systemBackground))
        .shadow(radius: 4)
    }
}

#Preview {
    PodcastView()
}
