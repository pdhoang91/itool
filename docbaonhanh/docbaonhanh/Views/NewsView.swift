// Views/NewsView.swift
import SwiftUI

struct NewsView: View {
    @StateObject private var viewModel = NewsViewModel()
    
    var body: some View {
        NavigationView {
            List {
                ForEach(viewModel.newsItems) { item in
                    NavigationLink(destination: NewsDetailView(item: item)) {
                        NewsItemRow(item: item)
                    }
                }
            }
            .navigationTitle("Tin tức")
            .onAppear {
                viewModel.fetchNews()
            }
        }
    }
}

#Preview {
    NewsView()
}
