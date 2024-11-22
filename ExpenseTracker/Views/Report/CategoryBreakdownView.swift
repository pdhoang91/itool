// Views/Report/CategoryBreakdownView.swift
struct CategoryBreakdownView: View {
    let categories: [CategoryTotal]
    @State private var selectedCategory: CategoryTotal?
    @State private var animateChart = false
    
    var body: some View {
        VStack {
            // Pie Chart
            PieChartView(
                categories: categories,
                selectedCategory: $selectedCategory
            )
            .frame(height: 250)
            .padding()
            .animation(.easeInOut(duration: 0.5), value: animateChart)
            
            // Category List
            List(categories) { category in
                CategoryRowView(category: category)
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(category.id == selectedCategory?.id ? 
                                  Color.blue.opacity(0.1) : Color.clear)
                    )
                    .onTapGesture {
                        withAnimation {
                            selectedCategory = category
                        }
                    }
            }
        }
        .onAppear {
            animateChart = true
        }
    }
}