// Tests/TransactionViewModelTests.swift
import XCTest
@testable import ExpenseTracker

class TransactionViewModelTests: XCTestCase {
    var sut: TransactionViewModel!
    var mockAPIService: MockAPIService!
    
    override func setUp() {
        super.setUp()
        mockAPIService = MockAPIService()
        sut = TransactionViewModel(apiService: mockAPIService)
    }
    
    func testFetchTransactions() async throws {
        // Given
        let expectedTransactions = [Transaction.mock(), Transaction.mock()]
        mockAPIService.mockResponse = expectedTransactions
        
        // When
        await sut.fetchTransactions()
        
        // Then
        XCTAssertEqual(sut.transactions.count, expectedTransactions.count)
        XCTAssertEqual(sut.transactions.first?.id, expectedTransactions.first?.id)
    }
    
    func testAddTransaction() async throws {
        // Given
        let transaction = Transaction.mock()
        mockAPIService.mockResponse = transaction
        
        // When
        try await sut.addTransaction(transaction)
        
        // Then
        XCTAssertEqual(sut.transactions.count, 1)
        XCTAssertEqual(sut.transactions.first?.id, transaction.id)
    }
}