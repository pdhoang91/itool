// Views/Components/LazyLoadingImage.swift
struct LazyLoadingImage: View {
    let url: URL
    @State private var image: UIImage?
    @State private var isLoading = true
    
    var body: some View {
        Group {
            if let image = image {
                Image(uiImage: image)
                    .resizable()
                    .aspectRatio(contentMode: .fit)
            } else if isLoading {
                ProgressView()
            }
        }
        .task {
            await loadImage()
        }
    }
    
    private func loadImage() async {
        if let cachedImage = await ImageCache.shared.image(for: url.absoluteString) {
            image = cachedImage
            isLoading = false
            return
        }
        
        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            guard let downloadedImage = UIImage(data: data) else { return }
            await ImageCache.shared.insertImage(downloadedImage, for: url.absoluteString)
            image = downloadedImage
        } catch {
            print("Failed to load image: \(error)")
        }
        
        isLoading = false
    }
}