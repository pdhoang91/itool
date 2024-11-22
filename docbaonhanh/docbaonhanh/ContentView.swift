//
//  ContentView.swift
//  docbaonhanh
//
//  Created by Pham Hoang on 23/11/24.
//

// ContentView.swift
import SwiftUI

struct ContentView: View {
    var body: some View {
        TabView {
            NewsView()
                .tabItem {
                    Image(systemName: "newspaper")
                    Text("Tin tức")
                }
            
            PodcastView()
                .tabItem {
                    Image(systemName: "headphones")
                    Text("Podcast")
                }
            
            SettingsView()
                .tabItem {
                    Image(systemName: "gear")
                    Text("Tùy chọn")
                }
        }
    }
}

#Preview {
    ContentView()
}
