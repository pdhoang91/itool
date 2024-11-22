// ViewModels/TransactionViewModel.swift
@MainActor
class TransactionViewModel: ObservableObject {
    @Published var transactions: [Transaction] = []
    @Published var recentTransactions: [Transaction] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var totalSpent: Double = 0
    @Published var spendingLimit: Double = 1000 // Default value
    
    var remaining: Double {
        spendingLimit - totalSpent
    }
    
    var spendingProgress: Double {
        totalSpent / spendingLimit
    }
    
    func fetchTransactions() async {
        isLoading = true
        
        do {
            let response: [Transaction] = try await APIService.shared.request(endpoint: "/transactions")
            transactions = response
            recentTransactions = Array(response.prefix(5))
            calculateTotalSpent()
        } catch {
            errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
    
    func addTransaction(_ transaction: Transaction) async throws {
        let jsonData = try JSONEncoder().encode(transaction)
        
        do {
            let response: Transaction = try await APIService.shared.request(
                endpoint: "/transactions",
                method: "POST",
                body: jsonData
            )
            transactions.append(response)
            calculateTotalSpent()
        } catch {
            errorMessage = error.localizedDescription
            throw error
        }
    }
    
    private func calculateTotalSpent() {
        totalSpent = transactions.reduce(0) { $0 + $1.amount }
    }
}