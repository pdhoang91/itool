// Views/MiniPlayerView.swift
import SwiftUI
import AVFoundation

struct MiniPlayerView: View {
    @ObservedObject var audioService: AudioPlayerService
    let newsTitle: String
    
    var body: some View {
        VStack(spacing: 0) {
            Divider()
            
            HStack(spacing: 12) {
                // Title
                Text(newsTitle)
                    .lineLimit(1)
                    .font(.subheadline)
                
                Spacer()
                
                // Controls
                Button(action: {
                    audioService.togglePlayPause()
                }) {
                    Image(systemName: audioService.isPlaying ? "pause.fill" : "play.fill")
                        .font(.title3)
                }
                .foregroundColor(.blue)
            }
            .padding(.horizontal)
            .padding(.vertical, 8)
            
            // Progress bar
            GeometryReader { geometry in
                Rectangle()
                    .fill(Color.blue)
                    .frame(width: geometry.size.width * audioService.progress, height: 2)
            }
            .frame(height: 2)
        }
        .background(Color(UIColor.systemBackground))
    }
}
