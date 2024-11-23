// Views/NewsDetailView.swift
// Views/NewsDetailView.swift
import SwiftUI
import AVFoundation

struct NewsDetailView: View {
    let item: NewsItem
    @State private var showAudioPlayer = false
    let audioService = AudioPlayerService.shared
    
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                AsyncImage(url: URL(string: item.imageUrl)) { image in
                    image
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                } placeholder: {
                    Color.gray
                }
                .frame(height: 200)
                .clipped()
                
                VStack(alignment: .leading, spacing: 12) {
                    Text(item.title)
                        .font(.title)
                        .fontWeight(.bold)
                    
                    HStack {
                        Text(item.source)
                            .foregroundColor(.blue)
                        Spacer()
                        Text(item.publishedDate, style: .date)
                            .foregroundColor(.gray)
                    }
                    .font(.subheadline)
                    
                    if let audioUrl = URL(string: item.audioUrl ?? "") {
                        Button(action: {
                            audioService.play(url: audioUrl)
                            showAudioPlayer = true
                        }) {
                            HStack {
                                Image(systemName: "play.circle.fill")
                                    .font(.title2)
                                Text("Nghe tin")
                                    .font(.headline)
                            }
                            .foregroundColor(.blue)
                            .padding(.vertical, 8)
                        }
                    }
                    
                    Divider()
                    
                    Text(item.content)
                        .font(.body)
                }
                .padding(.horizontal)
            }
        }
        .navigationBarTitleDisplayMode(.inline)
        .sheet(isPresented: $showAudioPlayer) {
            AudioPlayerView(audioService: audioService, newsTitle: item.title)
                .presentationDetents([.height(200)])
        }
        .overlay(
            VStack {
                Spacer()
                if showAudioPlayer {
                    MiniPlayerView(audioService: audioService, newsTitle: item.title)
                }
            }
        )
    }
}

#Preview {
    ContentView()
}

