// Models/Report.swift
struct Report: Codable {
    let totalSpent: Double
    let categoryBreakdown: [CategoryTotal]
    let timeRange: String
}

struct CategoryTotal: Identifiable, Codable {
    let id: String
    let category: String
    let total: Double
    let percentage: Double
}