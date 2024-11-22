// ViewModels/OfflineTransactionViewModel.swift
@MainActor
class OfflineTransactionViewModel: ObservableObject {
    @Published var transactions: [Transaction] = []
    private let coreDataManager = CoreDataManager.shared
    
    func saveTransaction(_ transaction: Transaction) async throws {
        // Save to local storage
        let context = coreDataManager.viewContext
        let localTransaction = LocalTransaction(context: context)
        localTransaction.update(from: transaction)
        coreDataManager.saveContext()
        
        // Enqueue for sync
        let syncTask = TransactionSyncTask(transaction: transaction)
        SyncManager.shared.enqueueSyncTask(syncTask)
        
        // Update UI
        await MainActor.run {
            transactions.append(transaction)
        }
    }
}