// Views/AudioPlayerView.swift
import SwiftUI
import AVFoundation

struct AudioPlayerView: View {
    @ObservedObject var audioService: AudioPlayerService
    let newsTitle: String
    
    private func formatTime(_ time: Double) -> String {
        let minutes = Int(time) / 60
        let seconds = Int(time) % 60
        return String(format: "%d:%02d", minutes, seconds)
    }
    
    var body: some View {
        VStack(spacing: 12) {
            // Title
            Text(newsTitle)
                .font(.headline)
                .lineLimit(2)
                .multilineTextAlignment(.center)
                .padding(.horizontal)
            
            // Progress Slider
            HStack {
                Text(formatTime(audioService.currentTime))
                    .font(.caption)
                    .monospacedDigit()
                
                Slider(value: Binding(
                    get: { audioService.progress },
                    set: { newValue in
                        audioService.seek(to: newValue * audioService.duration)
                    }
                ))
                
                Text(formatTime(audioService.duration))
                    .font(.caption)
                    .monospacedDigit()
            }
            .padding(.horizontal)
            
            // Controls
            HStack(spacing: 40) {
                Button(action: {
                    audioService.seek(to: max(0, audioService.currentTime - 15))
                }) {
                    Image(systemName: "gobackward.15")
                        .font(.title2)
                }
                
                Button(action: {
                    audioService.togglePlayPause()
                }) {
                    Image(systemName: audioService.isPlaying ? "pause.circle.fill" : "play.circle.fill")
                        .font(.title)
                }
                
                Button(action: {
                    audioService.seek(to: min(audioService.duration, audioService.currentTime + 15))
                }) {
                    Image(systemName: "goforward.15")
                        .font(.title2)
                }
            }
            .foregroundColor(.blue)
        }
        .padding(.vertical)
        .background(Color(UIColor.systemBackground))
        .cornerRadius(12)
        .shadow(radius: 4)
    }
}
