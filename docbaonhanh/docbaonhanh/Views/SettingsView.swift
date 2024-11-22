// Views/SettingsView.swift
import SwiftUI

struct SettingsView: View {
    @State private var selectedSources: Set<String> = []
    let availableSources = ["VnExpress", "24h", "Người Đưa Tin", "FireAnt"]
    
    var body: some View {
        NavigationView {
            Form {
                Section(header: Text("Nguồn tin tức")) {
                    ForEach(availableSources, id: \.self) { source in
                        Toggle(source, isOn: Binding(
                            get: { selectedSources.contains(source) },
                            set: { isSelected in
                                if isSelected {
                                    selectedSources.insert(source)
                                } else {
                                    selectedSources.remove(source)
                                }
                            }
                        ))
                    }
                }
                
                Section(header: Text("Thông tin ứng dụng")) {
                    HStack {
                        Text("Phiên bản")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.secondary)
                    }
                }
            }
            .navigationTitle("Tùy chọn")
        }
    }
}

#Preview {
    SettingsView()
}
