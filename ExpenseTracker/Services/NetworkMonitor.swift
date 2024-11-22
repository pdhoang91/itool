// Services/NetworkMonitor.swift
import Network

class NetworkMonitor: ObservableObject {
    static let shared = NetworkMonitor()
    private let monitor = NWPathMonitor()
    @Published private(set) var isConnected = true
    
    private init() {
        startMonitoring()
    }
    
    func startMonitoring() {
        monitor.pathUpdateHandler = { [weak self] path in
            DispatchQueue.main.async {
                self?.isConnected = path.status == .satisfied
            }
        }
        monitor.start(queue: DispatchQueue.global())
    }
}