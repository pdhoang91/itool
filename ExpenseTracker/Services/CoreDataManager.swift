// Services/CoreDataManager.swift
import CoreData

class CoreDataManager {
    static let shared = CoreDataManager()
    
    private init() {}
    
    lazy var persistentContainer: NSPersistentContainer = {
        let container = NSPersistentContainer(name: "ExpenseTracker")
        container.loadPersistentStores { _, error in
            if let error = error {
                fatalError("Failed to load Core Data stack: \(error)")
            }
        }
        return container
    }()
    
    var viewContext: NSManagedObjectContext {
        persistentContainer.viewContext
    }
    
    func saveContext() {
        guard viewContext.hasChanges else { return }
        
        do {
            try viewContext.save()
        } catch {
            print("Error saving context: \(error)")
        }
    }
}
