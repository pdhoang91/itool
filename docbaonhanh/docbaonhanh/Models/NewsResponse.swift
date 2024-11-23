// Models/NewsResponse.swift
import Foundation

// Models/NewsResponse.swift
// Models/NewsResponse.swift
struct NewsResponse: Codable {
    let items: [NewsItem]
    let currentPage: Int
    let totalPages: Int
    let totalItems: Int
    let hasNext: Bool
}

struct Pagination: Codable {
    let currentPage: Int
    let totalPages: Int
    let totalItems: Int
    let hasNext: Bool
}
