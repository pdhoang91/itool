// Models/Transaction.swift
struct Transaction: Identifiable, Codable {
    let id: Int
    let amount: Double
    let category: String
    let note: String?
    let date: Date
    
    var formattedAmount: String {
        return String(format: "%.2f", amount)
    }
}