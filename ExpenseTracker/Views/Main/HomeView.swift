// Views/Main/HomeView.swift
struct HomeView: View {
    @StateObject private var viewModel = TransactionViewModel()
    @State private var showingAddTransaction = false
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Spending Overview Card
                    VStack(alignment: .leading, spacing: 10) {
                        Text("This Month")
                            .font(.headline)
                        
                        HStack {
                            VStack(alignment: .leading) {
                                Text("Spent")
                                    .foregroundColor(.secondary)
                                Text("$\(viewModel.totalSpent, specifier: "%.2f")")
                                    .font(.title)
                                    .fontWeight(.bold)
                            }
                            
                            Spacer()
                            
                            VStack(alignment: .trailing) {
                                Text("Remaining")
                                    .foregroundColor(.secondary)
                                Text("$\(viewModel.remaining, specifier: "%.2f")")
                                    .font(.title2)
                                    .foregroundColor(viewModel.remaining < 0 ? .red : .green)
                            }
                        }
                        
                        ProgressView(value: viewModel.spendingProgress)
                            .tint(viewModel.spendingProgress > 0.9 ? .red : .blue)
                    }
                    .padding()
                    .background(Color(.systemBackground))
                    .cornerRadius(12)
                    .shadow(radius: 2)
                    
                    // Recent Transactions
                    VStack(alignment: .leading, spacing: 10) {
                        Text("Recent Transactions")
                            .font(.headline)
                        
                        ForEach(viewModel.recentTransactions) { transaction in
                            TransactionRowView(transaction: transaction)
                        }
                        
                        Button("See All") {
                            // Navigate to transaction list
                        }
                        .frame(maxWidth: .infinity)
                        .foregroundColor(.blue)
                    }
                    .padding()
                    .background(Color(.systemBackground))
                    .cornerRadius(12)
                    .shadow(radius: 2)
                }
                .padding()
            }
            .navigationTitle("Overview")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { showingAddTransaction = true }) {
                        Image(systemName: "plus")
                    }
                }
            }
            .sheet(isPresented: $showingAddTransaction) {
                AddTransactionView()
            }
        }
    }
}