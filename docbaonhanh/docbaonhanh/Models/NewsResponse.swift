// Models/NewsResponse.swift
import Foundation

struct NewsResponse: Codable {
    let items: [NewsItem]
    let pagination: PaginationInfo
}

struct PaginationInfo: Codable {
    let currentPage: Int
    let totalPages: Int
    let totalItems: Int
    let hasNext: Bool
}
