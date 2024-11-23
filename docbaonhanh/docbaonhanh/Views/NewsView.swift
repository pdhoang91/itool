
// Views/NewsView.swift
import SwiftUI
struct NewsView: View {
    @StateObject private var viewModel = NewsViewModel()
    
    var body: some View {
        NavigationView {
            ZStack {
                List {
                    ForEach(viewModel.newsItems) { item in
                        NavigationLink(destination: NewsDetailView(item: item)) {
                            NewsItemRow(item: item)
                                .onAppear {
                                    viewModel.loadMoreIfNeeded(currentItem: item)
                                }
                        }
                    }
                    
                    if viewModel.hasMorePages {
                        ProgressView()
                            .frame(maxWidth: .infinity)
                            .padding()
                    }
                }
                .refreshable {
                    viewModel.fetchInitialNews()
                }
                
                if viewModel.isLoading && viewModel.newsItems.isEmpty {
                    ProgressView()
                }
                
                if let error = viewModel.error {
                    VStack {
                        Text("Có lỗi xảy ra")
                            .font(.headline)
                        Text(error.localizedDescription)
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                        Button("Thử lại") {
                            viewModel.fetchInitialNews()
                        }
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    .shadow(radius: 4)
                }
            }
            .navigationTitle("Tin tức")
        }
        .onAppear {
            if viewModel.newsItems.isEmpty {
                viewModel.fetchInitialNews()
            }
        }
    }
}

#Preview {
    NewsView()
}
