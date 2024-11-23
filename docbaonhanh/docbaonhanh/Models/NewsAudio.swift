// Models/NewsAudio.swift
import Foundation

struct NewsAudio: Identifiable, Codable {
    let id: Int
    let newsId: Int
    let audioUrl: String
    let duration: Double
    let createdAt: Date
    let updatedAt: Date
}

struct AudioProgress: Codable {
    let currentTime: Double
    let duration: Double
    var progress: Double {
        return currentTime / duration
    }
}