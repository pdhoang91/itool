// Views/NewsDetailView.swift
// Views/NewsDetailView.swift
import SwiftUI

struct NewsDetailView: View {
    let item: NewsItem
    
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 16) {
                AsyncImage(url: URL(string: item.imageUrl)) { image in
                    image
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                } placeholder: {
                    Color.gray
                }
                .frame(height: 200)
                .clipped()
                
                VStack(alignment: .leading, spacing: 12) {
                    Text(item.title)
                        .font(.title)
                        .fontWeight(.bold)
                    
                    HStack {
                        Text(item.source)
                            .foregroundColor(.blue)
                        Spacer()
                        Text(item.publishedDate, style: .date)
                            .foregroundColor(.gray)
                    }
                    .font(.subheadline)
                    
                    Divider()
                    
                    Text(item.content)
                        .font(.body)
                }
                .padding(.horizontal)
            }
        }
        .navigationBarTitleDisplayMode(.inline)
    }
}

#Preview {
    ContentView()
}

