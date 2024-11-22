// Views/Authentication/LoginView.swift
struct LoginView: View {
    @StateObject private var viewModel = AuthViewModel()
    @State private var email = ""
    @State private var password = ""
    @State private var showingAlert = false
    
    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                Text("Expense Tracker")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                
                TextField("Email", text: $email)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                    .keyboardType(.emailAddress)
                    .autocapitalization(.none)
                
                SecureField("Password", text: $password)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                
                Button(action: login) {
                    if viewModel.isLoading {
                        ProgressView()
                    } else {
                        Text("Login")
                            .frame(maxWidth: .infinity)
                    }
                }
                .buttonStyle(.borderedProminent)
                .disabled(viewModel.isLoading)
                
                NavigationLink("Create Account", destination: RegisterView())
                    .foregroundColor(.secondary)
            }
            .padding()
            .alert("Error", isPresented: $showingAlert) {
                Button("OK", role: .cancel) {}
            } message: {
                Text(viewModel.errorMessage ?? "An error occurred")
            }
        }
    }
    
    private func login() {
        Task {
            do {
                try await viewModel.login(email: email, password: password)
            } catch {
                showingAlert = true
            }
        }
    }
}