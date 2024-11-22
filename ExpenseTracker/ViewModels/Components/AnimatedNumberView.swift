// Views/Components/AnimatedNumberView.swift
struct AnimatedNumberView: View {
    let value: Double
    let format: String
    @State private var animatedValue: Double = 0
    
    var body: some View {
        Text(String(format: format, animatedValue))
            .onAppear {
                withAnimation(.easeOut(duration: 1.0)) {
                    animatedValue = value
                }
            }
            .onChange(of: value) { newValue in
                withAnimation(.easeOut(duration: 0.5)) {
                    animatedValue = newValue
                }
            }
    }
}