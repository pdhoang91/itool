// Services/AudioPlayerService.swift
import AVFoundation
import Combine

class AudioPlayerService: ObservableObject {
    static let shared = AudioPlayerService()
    
    @Published var isPlaying = false
    @Published var currentTime: Double = 0
    @Published var duration: Double = 0
    @Published var progress: Double = 0
    @Published var currentIndex: Int = 0
    @Published var playlist: [NewsItem] = []
    
    private var player: AVPlayer?
    private var timeObserver: Any?
    private var playerItemObserver: AnyCancellable?
    
    private init() {
        setupAudioSession()
    }
    
    private func setupAudioSession() {
        do {
            try AVAudioSession.sharedInstance().setCategory(.playback, mode: .default)
            try AVAudioSession.sharedInstance().setActive(true)
        } catch {
            print("Failed to setup audio session: \(error)")
        }
    }
    
    func setPlaylist(_ items: [NewsItem]) {
        playlist = items
        if let firstItem = items.first, let audioUrl = URL(string: firstItem.audioUrl ?? "") {
            currentIndex = 0
            play(url: audioUrl)
        }
    }
    
    func play(url: URL) {
        // Remove existing player and observer
        stop()
        
        let playerItem = AVPlayerItem(url: url)
        player = AVPlayer(playerItem: playerItem)
        
        // Add time observer
        timeObserver = player?.addPeriodicTimeObserver(forInterval: CMTime(seconds: 0.5, preferredTimescale: 600), queue: .main) { [weak self] time in
            guard let self = self else { return }
            self.currentTime = time.seconds
            if let duration = self.player?.currentItem?.duration.seconds, duration > 0 {
                self.duration = duration
                self.progress = self.currentTime / duration
            }
        }
        
        // Observe when item finishes playing
        playerItemObserver = NotificationCenter.default.publisher(for: .AVPlayerItemDidPlayToEndTime)
            .sink { [weak self] _ in
                self?.playNextItem()
            }
        
        player?.play()
        isPlaying = true
    }
    
    func playNextItem() {
        guard !playlist.isEmpty else { return }
        
        currentIndex = (currentIndex + 1) % playlist.count
        if let audioUrl = URL(string: playlist[currentIndex].audioUrl ?? "") {
            play(url: audioUrl)
        }
    }
    
    func playPreviousItem() {
        guard !playlist.isEmpty else { return }
        
        currentIndex = currentIndex == 0 ? playlist.count - 1 : currentIndex - 1
        if let audioUrl = URL(string: playlist[currentIndex].audioUrl ?? "") {
            play(url: audioUrl)
        }
    }
    
    func togglePlayPause() {
        if isPlaying {
            player?.pause()
        } else {
            player?.play()
        }
        isPlaying.toggle()
    }
    
    func seek(to time: Double) {
        let cmTime = CMTime(seconds: time, preferredTimescale: 600)
        player?.seek(to: cmTime)
    }
    
    func stop() {
        player?.pause()
        if let timeObserver = timeObserver {
            player?.removeTimeObserver(timeObserver)
        }
        playerItemObserver?.cancel()
        player = nil
        timeObserver = nil
        isPlaying = false
        currentTime = 0
        duration = 0
        progress = 0
    }
}
