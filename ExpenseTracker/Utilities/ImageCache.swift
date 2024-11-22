// Utilities/ImageCache.swift
actor ImageCache {
    static let shared = ImageCache()
    
    private var cache: NSCache<NSString, UIImage> = {
        let cache = NSCache<NSString, UIImage>()
        cache.countLimit = 100
        cache.totalCostLimit = 50 * 1024 * 1024 // 50MB
        return cache
    }()
    
    func image(for key: String) -> UIImage? {
        cache.object(forKey: key as NSString)
    }
    
    func insertImage(_ image: UIImage, for key: String) {
        cache.setObject(image, forKey: key as NSString)
    }
}   