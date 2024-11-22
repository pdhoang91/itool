// // Views/NewsView.swift
// import SwiftUI

// struct NewsView: View {
//     @StateObject private var viewModel = NewsViewModel()
    
//     var body: some View {
//         NavigationView {
//             List {
//                 ForEach(viewModel.newsItems) { item in
//                     NavigationLink(destination: NewsDetailView(item: item)) {
//                         NewsItemRow(item: item)
//                     }
//                 }
//             }
//             .navigationTitle("Tin tức")
//             .onAppear {
//                 viewModel.fetchNews()
//             }
//         }
//     }
// }

// #Preview {
//     NewsView()
// }

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

//// Loading Indicator Component
//struct LoadingRow: View {
//    var body: some View {
//        HStack {
//            Spacer()
//            ProgressView()
//            Spacer()
//        }
//        .padding()
//    }
//}
//
//// Error View Component
//struct ErrorView: View {
//    let error: Error
//    let retryAction: () -> Void
//    
//    var body: some View {
//        VStack(spacing: 16) {
//            Text("Có lỗi xảy ra")
//                .font(.headline)
//            Text(error.localizedDescription)
//                .font(.subheadline)
//                .foregroundColor(.secondary)
//                .multilineTextAlignment(.center)
//            Button(action: retryAction) {
//                Text("Thử lại")
//                    .padding(.horizontal, 24)
//                    .padding(.vertical, 12)
//                    .background(Color.blue)
//                    .foregroundColor(.white)
//                    .cornerRadius(8)
//            }
//        }
//        .padding()
//        .background(Color(UIColor.systemBackground))
//        .cornerRadius(12)
//        .shadow(radius: 4)
//    }
//}

#Preview {
    NewsView()
}
