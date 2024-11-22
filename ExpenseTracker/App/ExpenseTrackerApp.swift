
import SwiftUI

struct ContentView: View {
    var body: some View {
        VStack {
            Text("Chào mừng bạn đến với ứng dụng của tôi!")
                .font(.largeTitle)
                .padding()
            Button(action: {
                print("Button tapped!")
            }) {
                Text("Nhấn vào đây")
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)
                    .cornerRadius(8)
            }
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
