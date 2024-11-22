// Models/SpendingLimit.swift
struct SpendingLimit: Codable {
    let amount: Double
    let period: String // "weekly" or "monthly"
}