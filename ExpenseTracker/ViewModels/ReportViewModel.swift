// ViewModels/ReportViewModel.swift
@MainActor
class ReportViewModel: ObservableObject {
    @Published var report: Report?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    func fetchReport(timeRange: String) async {
        isLoading = true
        
        do {
            report = try await APIService.shared.request(endpoint: "/reports?range=\(timeRange)")
        } catch {
            errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
}