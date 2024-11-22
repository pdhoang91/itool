// Views/Transaction/TransactionDetailView.swift
struct TransactionDetailView: View {
    let transaction: Transaction
    @Environment(\.dismiss) private var dismiss
    @StateObject private var viewModel = TransactionDetailViewModel()
    @State private var showingDeleteAlert = false
    
    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                // Header with amount
                HStack {
                    Spacer()
                    Text(transaction.formattedAmount)
                        .font(.system(size: 48, weight: .bold))
                        .foregroundColor(transaction.amount < 0 ? .red : .green)
                        .transition(.scale.combined(with: .opacity))
                    Spacer()
                }
                .padding()
                
                // Category Icon
                CategoryIconView(category: transaction.category)
                    .frame(width: 80, height: 80)
                    .transition(.scale)
                
                // Details
                VStack(alignment: .leading, spacing: 15) {
                    DetailRow(title: "Category", value: transaction.category)
                    DetailRow(title: "Date", value: transaction.date.formatted())
                    if let note = transaction.note {
                        DetailRow(title: "Note", value: note)
                    }
                }
                .padding()
                .background(Color(.systemBackground))
                .cornerRadius(12)
                .shadow(radius: 2)
                
                // Location if available
                if let location = transaction.location {
                    MapView(coordinate: location)
                        .frame(height: 200)
                        .cornerRadius(12)
                        .shadow(radius: 2)
                }
                
                // Attachments
                if !viewModel.attachments.isEmpty {
                    AttachmentsGridView(attachments: viewModel.attachments)
                }
                
                // Delete Button
                Button(action: { showingDeleteAlert = true }) {
                    Text("Delete Transaction")
                        .foregroundColor(.red)
                        .frame(maxWidth: .infinity)
                }
                .padding()
                .background(Color(.systemBackground))
                .cornerRadius(12)
                .shadow(radius: 2)
            }
            .padding()
        }
        .navigationBarTitleDisplayMode(.inline)
        .alert("Delete Transaction", isPresented: $showingDeleteAlert) {
            Button("Delete", role: .destructive) {
                Task {
                    await viewModel.deleteTransaction(transaction)
                    dismiss()
                }
            }
            Button("Cancel", role: .cancel) {}
        }
    }
}