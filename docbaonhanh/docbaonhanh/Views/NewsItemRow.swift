// Views/NewsItemRow.swift
import SwiftUI

struct NewsItemRow: View {
    let item: NewsItem
    
    var body: some View {
        HStack(spacing: 12) {
            // Left side - Info
            VStack(alignment: .leading, spacing: 8) {
                Text(item.title)
                    .font(.headline)
                    .lineLimit(2)
                
                Text(item.summary)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                    .lineLimit(2)
                
                HStack {
                    Text(item.source)
                        .font(.caption)
                        .foregroundColor(.blue)
                    
                    Spacer()
                    
                    Text(item.publishedDate, style: .date)
                        .font(.caption)
                        .foregroundColor(.gray)
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .leading)
            
            // Right side - Image (1/3 width)
            AsyncImage(url: URL(string: item.imageUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Color.gray
            }
            .frame(width: UIScreen.main.bounds.width * 0.25)
            .clipped()
            .cornerRadius(8)
        }
        .frame(height: 100)
        .padding(.vertical, 4)
    }
}