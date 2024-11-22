// Security/CertificatePinning.swift
class CertificatePinning: NSObject, URLSessionDelegate {
    static let shared = CertificatePinning()
    
    private let pinnedCertificates: [Data] = {
        let certificateNames = ["certificate1", "certificate2"]
        return certificateNames.compactMap { name in
            guard let certificatePath = Bundle.main.path(forResource: name, ofType: "cer"),
                  let certificateData = try? Data(contentsOf: URL(fileURLWithPath: certificatePath)) else {
                return nil
            }
            return certificateData
        }
    }()
    
    func urlSession(
        _ session: URLSession,
        didReceive challenge: URLAuthenticationChallenge,
        completionHandler: @escaping (URLSession.AuthChallengeDisposition, URLCredential?) -> Void
    ) {
        guard let serverTrust = challenge.protectionSpace.serverTrust,
              let certificate = SecTrustGetCertificateAtIndex(serverTrust, 0) else {
            completionHandler(.cancelAuthenticationChallenge, nil)
            return
        }
        
        let serverCertificateData = SecCertificateCopyData(certificate) as Data
        
        if pinnedCertificates.contains(serverCertificateData) {
            completionHandler(.useCredential, URLCredential(trust: serverTrust))
        } else {
            completionHandler(.cancelAuthenticationChallenge, nil)
        }
    }
}