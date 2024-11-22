// Views/Components/TransactionListAnimation.swift
struct TransactionListAnimation: ViewModifier {
    let index: Int
    
    func body(content: Content) -> some View {
        content
            .offset(x: 0, y: 50)
            .opacity(0)
            .animation(
                .spring(
                    response: 0.5,
                    dampingFraction: 0.8,
                    blendDuration: 0
                )
                .delay(Double(index) * 0.05),
                value: true
            )
    }
}