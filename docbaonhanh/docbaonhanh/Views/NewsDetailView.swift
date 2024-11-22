// Views/NewsDetailView.swift
import SwiftUI

struct NewsDetailView: View {
    let item: NewsItem
    @StateObject private var audioPlayer = AudioPlayerViewModel()
    
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                // Header Image
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
                    // Title and Meta
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
                    
                    Divider()
                    
                    // Audio Player Button
                    Button(action: {
                        audioPlayer.isPlaying ? audioPlayer.pause() : audioPlayer.play(item)
                    }) {
                        HStack {
                            Image(systemName: audioPlayer.isPlaying ? "pause.circle.fill" : "play.circle.fill")
                                .font(.title)
                            Text(audioPlayer.isPlaying ? "Tạm dừng" : "Nghe tin")
                        }
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.blue)
                        .cornerRadius(10)
                    }
                    
                    // Content
                    Text(item.content)
                        .font(.body)
                }
                .padding(.horizontal)
            }
        }
        .navigationBarTitleDisplayMode(.inline)
    }
}

#Preview {
    ContentView()
}

