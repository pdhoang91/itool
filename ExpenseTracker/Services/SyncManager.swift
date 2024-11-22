// Services/SyncManager.swift
class SyncManager {
    static let shared = SyncManager()
    private let queue = OperationQueue()
    private var syncTasks: [SyncTask] = []
    
    func enqueueSyncTask(_ task: SyncTask) {
        syncTasks.append(task)
        attemptSync()
    }
    
    private func attemptSync() {
        guard NetworkMonitor.shared.isConnected else { return }
        
        Task {
            for task in syncTasks {
                do {
                    try await task.execute()
                    syncTasks.removeAll { $0.id == task.id }
                } catch {
                    print("Sync failed: \(error)")
                }
            }
        }
    }
}