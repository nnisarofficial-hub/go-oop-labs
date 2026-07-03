# Go OOP Labs — 20 Real-World Challenges
### Structs, Interfaces, Embedding, and Dependency Injection

---

## Before You Start

Go does not have classes, inheritance, or the `extends` keyword. It achieves everything OOP can do through:

- **Structs** — data bundled together (like a class without methods baked in)
- **Methods** — functions attached to a struct
- **Interfaces** — contracts any type can satisfy, just by having the right methods
- **Embedding** — composing structs inside other structs (Go's answer to inheritance)

This is not a limitation — it's a deliberate design. Go favors **composition over inheritance**, which produces code that is easier to test, easier to change, and easier to understand.

These 20 labs build a mental model of how Go thinks about OOP. Work through them in order.

---

## Key Concepts Cheat Sheet

```go
// Struct — a named collection of fields
type Person struct {
    Name string
    Age  int
}

// Method — a function with a receiver
func (p Person) Greet() string {           // value receiver — p is a copy
    return "Hello, I'm " + p.Name
}

func (p *Person) Birthday() {              // pointer receiver — p is the real thing
    p.Age++                                // this change persists
}

// Interface — a set of method signatures
type Speaker interface {
    Greet() string
}

// Any type that has Greet() satisfies Speaker automatically — no "implements" keyword
var s Speaker = Person{Name: "Ali"}

// Embedding — include one struct inside another
type Employee struct {
    Person                                 // embedded — Employee now has Name, Age, Greet()
    Department string
}

// Constructor function — Go's pattern for building structs safely
func NewPerson(name string, age int) *Person {
    return &Person{Name: name, Age: age}
}
```

---

## Group 1 — Structs, Methods, and Constructors

---

### Lab 1 — Notification Service

**The scenario:** You're building a notification system. A `Notification` has a title, message, type (info/warning/error), and whether it's been read.

**You will practice:** Structs, value and pointer receivers, constructor functions, encapsulation with unexported fields

**Build this:**

```go
// Define this struct
type Notification struct {
    id      int    // unexported — can only be set by the constructor
    Title   string
    Message string
    Type    string // "info", "warning", "error"
    read    bool   // unexported — changed only through MarkAsRead()
}

// Constructor — the only way to create a Notification
func NewNotification(id int, title, message, notifType string) *Notification

// Methods to implement:
func (n *Notification) MarkAsRead()
func (n Notification) IsRead() bool
func (n Notification) Summary() string // returns: "[INFO] Title: Your message here"
func (n Notification) ID() int         // getter for unexported field
```

**Why unexported fields?** `id` and `read` should never be set directly from outside the package. `id` is assigned on creation, `read` only changes through `MarkAsRead()`. This is Go's encapsulation — not `private/public` keywords, but exported (uppercase) vs unexported (lowercase).

**Expected output:**
```
ID: 1
Summary: [WARNING] Low disk space: Your disk is 90% full
Is read: false

After marking as read:
Is read: true
```

**Hint:** The constructor function returns a `*Notification` (pointer), not a value. This is standard Go — when a struct has pointer receiver methods, always work with pointers.

---

### Lab 2 — Shape Area Calculator

**The scenario:** A geometry utility that can calculate areas and perimeters for any shape.

**You will practice:** Interfaces, implicit satisfaction, polymorphic functions

**Build this:**

```go
// Define this interface
type Shape interface {
    Area() float64
    Perimeter() float64
    Name() string
}

// Implement these three shapes:
type Circle struct {
    Radius float64
}

type Rectangle struct {
    Width, Height float64
}

type Triangle struct {
    A, B, C float64 // three sides
}
// Triangle area: Heron's formula — s = (a+b+c)/2, area = √(s(s-a)(s-b)(s-c))

// Write this function — it works with ANY shape
func PrintShapeInfo(s Shape) {
    // print name, area (2 decimal places), perimeter (2 decimal places)
}

// Write this function — finds the largest shape by area
func LargestShape(shapes []Shape) Shape {
    // TODO
}
```

**Expected output:**
```
Circle
  Area:      78.54
  Perimeter: 31.42

Rectangle
  Area:      24.00
  Perimeter: 20.00

Triangle
  Area:      6.00
  Perimeter: 12.00

Largest shape by area: Circle (78.54)
```

**Key insight:** `PrintShapeInfo` and `LargestShape` don't know what shape they're working with. They only know it satisfies the `Shape` interface. This is polymorphism in Go.

---

### Lab 3 — Employee and Manager with Embedding

**The scenario:** An HR system with two roles. A Manager is an Employee who also manages a team.

**You will practice:** Struct embedding, overriding promoted methods, accessing embedded fields

**Build this:**

```go
type Employee struct {
    ID         int
    Name       string
    Department string
    Salary     float64
}

func (e Employee) Describe() string {
    return fmt.Sprintf("Employee #%d: %s (%s) — PKR %.0f/month", e.ID, e.Name, e.Department, e.Salary)
}

func (e Employee) AnnualSalary() float64 {
    return e.Salary * 12
}

// Manager embeds Employee and adds team management
type Manager struct {
    Employee                // embed — Manager gets all Employee fields and methods
    Reports  []*Employee    // the team this manager is responsible for
}

// Override Describe() for Manager — should mention number of reports
func (m Manager) Describe() string {
    // TODO — call e.Employee.Describe() inside here, then add team size info
}

func (m *Manager) AddReport(e *Employee) {
    // TODO
}

func (m Manager) TeamAnnualCost() float64 {
    // TODO — sum of manager's salary + all reports' salaries (all × 12)
}
```

**Expected output:**
```
Employee #1: Sara Ahmed (Engineering) — PKR 150000/month
Annual salary: PKR 1800000

Manager #2: Ali Khan (Engineering) — PKR 220000/month — manages 3 reports
Team annual cost: PKR 5880000
```

**Key insight:** When you embed `Employee` in `Manager`, the manager automatically gets `Name`, `Department`, `Salary`, and `AnnualSalary()` for free — you access them directly: `m.Name`, `m.AnnualSalary()`. But you can override `Describe()` by defining it on `Manager` too.

---

### Lab 4 — Vehicle Fleet with Embedding

**The scenario:** A vehicle fleet management system. All vehicles share common properties but each type has unique behavior.

**You will practice:** Multi-level embedding, method promotion, constructor with validation

**Build this:**

```go
type Vehicle struct {
    Make         string
    Model        string
    Year         int
    Mileage      float64 // km
    FuelCapacity float64 // liters
}

func (v Vehicle) Info() string {
    return fmt.Sprintf("%d %s %s (%.0f km)", v.Year, v.Make, v.Model, v.Mileage)
}

func (v *Vehicle) Drive(km float64) {
    v.Mileage += km
}

// Truck extends Vehicle
type Truck struct {
    Vehicle
    PayloadCapacity float64 // tonnes
    CurrentLoad     float64
}

func NewTruck(make, model string, year int, payload float64) *Truck

func (t *Truck) LoadCargo(tonnes float64) error {
    // return error if load exceeds PayloadCapacity
}

func (t *Truck) Info() string {
    // override — include payload info: "2020 Toyota Hilux (0 km) | Load: 2.0/5.0t"
}

// ElectricCar extends Vehicle
type ElectricCar struct {
    Vehicle
    BatteryCapacity float64 // kWh
    ChargeLevel     float64 // 0.0 to 1.0 (percentage as decimal)
}

func NewElectricCar(make, model string, year int, battery float64) *ElectricCar

func (e *ElectricCar) Charge(toLevel float64) error {
    // error if toLevel > 1.0 or < current level
}

func (e *ElectricCar) RangeRemaining() float64 {
    // assume 6 km per kWh
    return e.BatteryCapacity * e.ChargeLevel * 6
}

func (e *ElectricCar) Info() string {
    // override — include battery info: "2023 Tesla Model 3 (0 km) | Battery: 80% | Range: ~288 km"
}
```

**Expected output:**
```
2020 Toyota Hilux (0 km) | Load: 0.0/5.0t
After loading 3 tonnes:
2020 Toyota Hilux (500 km) | Load: 3.0/5.0t

2023 Tesla Model 3 (0 km) | Battery: 20% | Range: ~72 km
After charging to 90%:
2023 Tesla Model 3 (250 km) | Battery: 90% | Range: ~324 km

Error: cannot load 4.0t, only 2.0t capacity remaining
```

---

## Group 2 — Interfaces

---

### Lab 5 — The Stringer Interface

**The scenario:** Go's `fmt.Println` calls `.String()` automatically if a type implements the `fmt.Stringer` interface. This makes your types print beautifully without any extra work from the caller.

**You will practice:** `fmt.Stringer`, interface satisfaction, formatting output

```go
// The fmt.Stringer interface (built into Go's fmt package):
// type Stringer interface {
//     String() string
// }
```

**Build this:**

Define these three types and implement `String()` on each so `fmt.Println(x)` prints nicely:

```go
type Color struct {
    R, G, B uint8 // 0–255
}
// String() → "rgb(255, 128, 0)"
// Also implement: Hex() string → "#FF8000"

type Temperature struct {
    Celsius float64
}
// String() → "37.0°C (98.6°F)"
// Also implement: IsFever() bool — true if > 37.5°C

type Money struct {
    Amount   float64
    Currency string
}
// String() → "PKR 1,250.00" or "USD 9.99"
// Also implement: Add(other Money) (Money, error) — error if currencies differ
//                 IsZero() bool
```

**Expected output:**
```
Color: rgb(255, 128, 0)
Hex:   #FF8000

Temperature: 38.5°C (101.3°F)
Is fever: true

Price: PKR 1,250.00
Total: PKR 2,500.00

Error: cannot add PKR and USD
```

**Key insight:** Any type with a `String() string` method automatically satisfies `fmt.Stringer`. You never write `fmt.Println(t.String())` — just `fmt.Println(t)` and Go calls `.String()` for you.

---

### Lab 6 — Custom Errors

**The scenario:** A user registration system that returns specific, informative errors rather than generic strings.

**You will practice:** The `error` interface, custom error types, type assertion on errors, `errors.As`

```go
// Go's built-in error interface:
// type error interface {
//     Error() string
// }
// Any struct with Error() string satisfies it.
```

**Build this:**

```go
// Define these custom error types:

type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string
// → "validation error on field 'email': invalid format"

type DuplicateError struct {
    Resource string
    Value    string
}
func (e *DuplicateError) Error() string
// → "duplicate 'username': 'ali123' already exists"

type NotFoundError struct {
    Resource string
    ID       string
}
func (e *NotFoundError) Error() string
// → "user with id '999' not found"

// User and registration logic
type User struct {
    ID       int
    Username string
    Email    string
    Age      int
}

// In-memory store — a map[int]*User
type UserStore struct {
    users  map[int]*User
    nextID int
}

func NewUserStore() *UserStore

func (s *UserStore) Register(username, email string, age int) (*User, error) {
    // Validate: username must be 3+ chars → ValidationError
    // Validate: email must contain "@" → ValidationError
    // Validate: age must be 18+ → ValidationError
    // Check: username must not already exist → DuplicateError
    // If all good: create and store the user
}

func (s *UserStore) FindByID(id int) (*User, error) {
    // Return NotFoundError if not found
}
```

In `main`, call `Register` with bad data and handle each error type differently:

```go
_, err := store.Register("al", "not-an-email", 15)
// use errors.As to check which type of error it is:
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("Validation failed:", valErr.Field, "—", valErr.Message)
}
```

**Expected output:**
```
Register("al", "ali@email.com", 25):
  Error: validation error on field 'username': must be at least 3 characters

Register("ali", "not-an-email", 25):
  Error: validation error on field 'email': invalid format

Register("ali", "ali@email.com", 15):
  Error: validation error on field 'age': must be 18 or older

Register("ali", "ali@email.com", 25): ✓ User created (ID: 1)

Register("ali", "ali2@email.com", 30):
  Error: duplicate 'username': 'ali' already exists

FindByID(999):
  Error: user with id '999' not found
```

---

### Lab 7 — Type Assertions and Type Switches

**The scenario:** A data pipeline that receives mixed event types and processes each differently.

**You will practice:** `interface{}` / `any`, type assertion `x.(Type)`, type switch `switch v := x.(type)`

**Build this:**

```go
// Events in the system
type LoginEvent struct {
    UserID    int
    IPAddress string
    Success   bool
}

type PurchaseEvent struct {
    UserID    int
    ProductID string
    Amount    float64
}

type ErrorEvent struct {
    Code    int
    Message string
    Fatal   bool
}

// Process handles any event type
func Process(event any) string {
    // Use a type switch to handle each type differently
    // LoginEvent:   "LOGIN user=42 ip=192.168.1.1 success=true"
    // PurchaseEvent: "PURCHASE user=42 product=SKU-001 amount=PKR 1500.00"
    // ErrorEvent:   "ERROR [404] page not found (fatal=false)"
    // unknown:      "UNKNOWN event type: <type>"
}

// ProcessBatch processes a slice of mixed events
func ProcessBatch(events []any) {
    for i, e := range events {
        fmt.Printf("[%d] %s\n", i+1, Process(e))
    }
}
```

In `main`, create a mixed slice of events and pass to `ProcessBatch`:

```go
events := []any{
    LoginEvent{UserID: 1, IPAddress: "192.168.1.1", Success: true},
    PurchaseEvent{UserID: 1, ProductID: "SKU-001", Amount: 1500},
    LoginEvent{UserID: 2, IPAddress: "10.0.0.5", Success: false},
    ErrorEvent{Code: 500, Message: "database timeout", Fatal: true},
    "unexpected string",  // unknown type
}
```

**Expected output:**
```
[1] LOGIN user=1 ip=192.168.1.1 success=true
[2] PURCHASE user=1 product=SKU-001 amount=PKR 1500.00
[3] LOGIN user=2 ip=10.0.0.5 success=false
[4] ERROR [500] database timeout (fatal=true)
[5] UNKNOWN event type: string
```

---

### Lab 8 — Interface Segregation

**The scenario:** A document management system. Not every document handler needs to do everything — some can only read, some can read and write, some can also delete.

**You will practice:** Small focused interfaces, composing interfaces, accepting the smallest interface you need

**Build this:**

```go
// Define three separate interfaces — small and focused
type Reader interface {
    Read(id string) (string, error)
}

type Writer interface {
    Write(id, content string) error
}

type Deleter interface {
    Delete(id string) error
}

// Compose them
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteDeleter interface {
    Reader
    Writer
    Deleter
}

// Implement an in-memory document store
type DocumentStore struct {
    docs map[string]string
}

func NewDocumentStore() *DocumentStore
func (d *DocumentStore) Read(id string) (string, error)
func (d *DocumentStore) Write(id, content string) error
func (d *DocumentStore) Delete(id string) error

// These functions accept the SMALLEST interface they need
// — they don't demand full access when they only need part of it

func DisplayDocument(r Reader, id string) {
    // reads and prints the document
}

func BackupDocument(rw ReadWriter, fromID, toID string) {
    // reads from fromID, writes a copy to toID
}

func ArchiveDocument(rwd ReadWriteDeleter, id, archiveID string) {
    // reads, writes to archive location, deletes original
}
```

**Expected output:**
```
Document 'doc-1': Hello, this is document one.

Backed up 'doc-1' → 'doc-1-backup'
doc-1-backup: Hello, this is document one.

Archived 'doc-2': moved to 'archive-doc-2', original deleted
Read deleted doc: document with id 'doc-2' not found
```

**Key insight:** `DisplayDocument` only needs a `Reader`. It doesn't know if the store can also write or delete. This is interface segregation — functions accept the minimum interface they need, which makes them easier to test and reuse.

---

### Lab 9 — Implementing sort.Interface

**The scenario:** A leaderboard system that needs to sort players by score, and by name as a tiebreaker.

**You will practice:** Implementing an external interface (`sort.Interface`), `sort.Sort`, `sort.Reverse`

```go
// Go's sort.Interface (from the standard library):
// type Interface interface {
//     Len() int
//     Less(i, j int) bool
//     Swap(i, j int)
// }
```

**Build this:**

```go
type Player struct {
    Name   string
    Score  int
    Rank   int
}

// ByScore sorts players by score descending, then by name ascending as tiebreaker
type ByScore []Player

func (b ByScore) Len() int
func (b ByScore) Less(i, j int) bool  // higher score = "less" (comes first)
func (b ByScore) Swap(i, j int)

// ByName sorts players alphabetically
type ByName []Player

func (b ByName) Len() int
func (b ByName) Less(i, j int) bool
func (b ByName) Swap(i, j int)

// After sorting, assign ranks 1, 2, 3...
func AssignRanks(players []Player) {
    // rank 1 = highest score
}

func PrintLeaderboard(players []Player) {
    // print formatted: "Rank 1 | Ali Khan      | 9850 pts"
}
```

In `main`:
```go
players := []Player{
    {Name: "Sara", Score: 9200},
    {Name: "Ali", Score: 9850},
    {Name: "Umar", Score: 9200},   // tied with Sara — Sara should come first (alphabetical)
    {Name: "Fatima", Score: 10100},
}

sort.Sort(ByScore(players))
AssignRanks(players)
PrintLeaderboard(players)
```

**Expected output:**
```
Rank 1 | Fatima        | 10100 pts
Rank 2 | Ali           |  9850 pts
Rank 3 | Sara          |  9200 pts
Rank 4 | Umar          |  9200 pts
```

---

## Group 3 — Dependency Injection via Interfaces

---

### Lab 10 — Logger Interface

**The scenario:** An application that uses different logging backends depending on the environment — console output in development, file output in production.

**You will practice:** Defining an interface for infrastructure, injecting it into a service, swapping implementations without changing service code

**Build this:**

```go
// The interface your service depends on
type Logger interface {
    Info(msg string)
    Warn(msg string)
    Error(msg string)
}

// Implementation 1: prints to console with color and timestamp
type ConsoleLogger struct {
    prefix string
}
func NewConsoleLogger(prefix string) *ConsoleLogger
func (l *ConsoleLogger) Info(msg string)  // "[INFO]  2024/01/15 10:30:00 [prefix] msg"
func (l *ConsoleLogger) Warn(msg string)  // "[WARN]  ..."
func (l *ConsoleLogger) Error(msg string) // "[ERROR] ..."

// Implementation 2: writes to a file
type FileLogger struct {
    filename string
    file     *os.File
}
func NewFileLogger(filename string) (*FileLogger, error)
func (l *FileLogger) Info(msg string)
func (l *FileLogger) Warn(msg string)
func (l *FileLogger) Error(msg string)
func (l *FileLogger) Close()

// Implementation 3: silent — discards all logs (useful in tests)
type NoopLogger struct{}
func (l *NoopLogger) Info(msg string)  {}
func (l *NoopLogger) Warn(msg string)  {}
func (l *NoopLogger) Error(msg string) {}

// A service that needs logging — it never imports ConsoleLogger or FileLogger directly
type OrderService struct {
    logger Logger // injected
    orders map[int]string
}

func NewOrderService(logger Logger) *OrderService

func (s *OrderService) PlaceOrder(id int, item string) {
    s.logger.Info(fmt.Sprintf("placing order %d for item: %s", id, item))
    s.orders[id] = item
    s.logger.Info(fmt.Sprintf("order %d placed successfully", id))
}

func (s *OrderService) CancelOrder(id int) error {
    if _, exists := s.orders[id]; !exists {
        s.logger.Warn(fmt.Sprintf("attempted to cancel non-existent order %d", id))
        return fmt.Errorf("order %d not found", id)
    }
    delete(s.orders, id)
    s.logger.Info(fmt.Sprintf("order %d cancelled", id))
    return nil
}
```

In `main`, run the same operations with different loggers:

```go
// Development
svc := NewOrderService(NewConsoleLogger("orders"))
svc.PlaceOrder(1, "Barbari Buck — Premium Grade")
svc.CancelOrder(999)

// Testing — no output at all
testSvc := NewOrderService(&NoopLogger{})
testSvc.PlaceOrder(2, "test item") // silent
```

**Expected output:**
```
[INFO]  2024/01/15 10:30:00 [orders] placing order 1 for item: Barbari Buck — Premium Grade
[INFO]  2024/01/15 10:30:00 [orders] order 1 placed successfully
[WARN]  2024/01/15 10:30:00 [orders] attempted to cancel non-existent order 999
```

**Key insight:** `OrderService` only knows about `Logger`. It doesn't import `ConsoleLogger` or `FileLogger`. You can swap logging behavior without touching the service — and in tests you pass `&NoopLogger{}` to silence all output.

---

### Lab 11 — Payment Processor

**The scenario:** An e-commerce checkout that supports multiple payment methods. Adding a new payment method should never require changing the checkout logic.

**You will practice:** DI with interfaces, error handling through interface, strategy pattern

**Build this:**

```go
type PaymentResult struct {
    TransactionID string
    Amount        float64
    Method        string
    Success       bool
}

type PaymentProcessor interface {
    Name() string
    Charge(amount float64, reference string) (*PaymentResult, error)
    Refund(transactionID string) error
}

// Implement three payment processors:

type CreditCard struct {
    CardNumber string // last 4 digits only for display
    HolderName string
}
func NewCreditCard(number, holder string) *CreditCard
func (c *CreditCard) Name() string // "CreditCard (**** 4242)"
func (c *CreditCard) Charge(amount float64, reference string) (*PaymentResult, error)
// Simulate: reject if amount > 100000 (card limit)

type JazzCash struct {
    PhoneNumber string
}
func NewJazzCash(phone string) *JazzCash
func (j *JazzCash) Name() string
func (j *JazzCash) Charge(amount float64, reference string) (*PaymentResult, error)
// Simulate: reject if amount > 25000 (wallet limit)

type BankTransfer struct {
    AccountNumber string
    BankName      string
}
func NewBankTransfer(account, bank string) *BankTransfer
func (b *BankTransfer) Name() string
func (b *BankTransfer) Charge(amount float64, reference string) (*PaymentResult, error)
// No limit — always succeeds (bank transfer)

// CheckoutService depends only on the interface
type CheckoutService struct {
    processor PaymentProcessor
    logger    Logger // reuse your Logger from Lab 10
}

func NewCheckoutService(processor PaymentProcessor, logger Logger) *CheckoutService

func (c *CheckoutService) Checkout(orderID string, amount float64) (*PaymentResult, error) {
    c.logger.Info(fmt.Sprintf("charging PKR %.0f via %s for order %s", amount, c.processor.Name(), orderID))
    result, err := c.processor.Charge(amount, orderID)
    if err != nil {
        c.logger.Error("payment failed: " + err.Error())
        return nil, err
    }
    c.logger.Info(fmt.Sprintf("payment successful: txn %s", result.TransactionID))
    return result, nil
}
```

**Expected output:**
```
Charging PKR 5000 via CreditCard (**** 4242) for order ORD-001
Payment successful: txn CC-A3F2B1

Charging PKR 30000 via JazzCash (0300-1234567) for order ORD-002
Payment failed: JazzCash: amount 30000 exceeds wallet limit of 25000

Charging PKR 30000 via BankTransfer (HBL) for order ORD-003
Payment successful: txn BT-C7D9E2
```

---

### Lab 12 — Multi-Channel Notification System

**The scenario:** An alerting system that can send the same notification through multiple channels simultaneously.

**You will practice:** Slice of interfaces, fan-out pattern, collecting errors from multiple implementations

**Build this:**

```go
type Notifier interface {
    Send(to, subject, body string) error
    Name() string
}

// Implement three notifiers:

type EmailNotifier struct {
    SMTPServer string
}
func (e *EmailNotifier) Name() string { return "Email" }
func (e *EmailNotifier) Send(to, subject, body string) error {
    // Simulate sending — just print what would be sent
    // Simulate failure: if to contains "invalid" return an error
}

type SMSNotifier struct {
    APIKey string
}
func (s *SMSNotifier) Name() string { return "SMS" }
func (s *SMSNotifier) Send(to, subject, body string) error {
    // SMS only sends the body (no subject) — truncate to 160 chars
    // Simulate failure: if to doesn't start with "03" return error
}

type SlackNotifier struct {
    Channel string
}
func (s *SlackNotifier) Name() string { return "Slack" }
func (s *SlackNotifier) Send(to, subject, body string) error {
    // Simulate sending to a Slack channel — always succeeds
}

// NotificationService sends to ALL notifiers and collects errors
type NotificationService struct {
    notifiers []Notifier
}

func NewNotificationService(notifiers ...Notifier) *NotificationService

// Broadcast sends to all notifiers and returns a map of notifier name → error (nil if success)
func (ns *NotificationService) Broadcast(to, subject, body string) map[string]error {
    // TODO: send to each notifier, collect results
    // don't stop if one fails — try all of them
}

func (ns *NotificationService) PrintReport(results map[string]error) {
    // print: "Email: ✓ sent" or "SMS: ✗ failed — error message"
}
```

**Expected output:**
```
Sending alert to ali@example.com / 0300-1234567...
Email: ✓ sent
SMS:   ✓ sent
Slack: ✓ sent

Sending alert to invalid@example.com / 0300-1234567...
Email: ✗ failed — invalid recipient address
SMS:   ✓ sent
Slack: ✓ sent
2 of 3 channels successful
```

---

### Lab 13 — Repository Pattern

**The scenario:** A product catalog service. The service doesn't know or care if products are stored in memory, a database, or a file. It just uses a `ProductRepository` interface.

**You will practice:** Repository pattern, DI, separating business logic from storage, writing testable code

**Build this:**

```go
type Product struct {
    ID       string
    Name     string
    Price    float64
    Stock    int
    Category string
}

// The repository interface — storage is hidden behind this contract
type ProductRepository interface {
    FindByID(id string) (*Product, error)
    FindAll() ([]*Product, error)
    FindByCategory(category string) ([]*Product, error)
    Save(p *Product) error
    Delete(id string) error
}

// In-memory implementation
type InMemoryProductRepo struct {
    products map[string]*Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo
func (r *InMemoryProductRepo) FindByID(id string) (*Product, error)
func (r *InMemoryProductRepo) FindAll() ([]*Product, error)
func (r *InMemoryProductRepo) FindByCategory(category string) ([]*Product, error)
func (r *InMemoryProductRepo) Save(p *Product) error
func (r *InMemoryProductRepo) Delete(id string) error

// Business logic — depends ONLY on the interface
type ProductService struct {
    repo   ProductRepository
    logger Logger
}

func NewProductService(repo ProductRepository, logger Logger) *ProductService

func (s *ProductService) AddProduct(id, name, category string, price float64, stock int) error {
    // validate: name not empty, price > 0, stock >= 0
    // log the action
    // save via repo
}

func (s *ProductService) GetCatalog() ([]*Product, error) {
    // return all products sorted by category then name
}

func (s *ProductService) UpdateStock(id string, delta int) error {
    // find product, update stock (reject if it would go below 0)
    // log the action
}

func (s *ProductService) GetLowStockAlerts(threshold int) ([]*Product, error) {
    // return products where stock <= threshold
}
```

**Expected output:**
```
Catalog (4 products):
  [electronics] Laptop Pro 15  — PKR 250,000  | Stock: 5
  [electronics] Wireless Mouse — PKR   3,500  | Stock: 50
  [feed]        Barseem Hay    — PKR     800  | Stock: 3
  [medicine]    Ivermectin     — PKR   1,200  | Stock: 2

Low stock (threshold: 5):
  Barseem Hay (3 remaining)
  Ivermectin (2 remaining)

After selling 4 units of Barseem Hay:
Error: insufficient stock (have 3, requested 4)
```

---

### Lab 14 — Cache Interface

**The scenario:** An API that fetches weather data from an external service. Fetching is slow and expensive — add a caching layer that can be swapped between an in-memory cache and a no-cache option.

**You will practice:** Interface wrapping, the decorator pattern, `time.Duration` for TTL

**Build this:**

```go
type Cache interface {
    Get(key string) (string, bool)        // bool = found
    Set(key string, value string, ttl time.Duration)
    Delete(key string)
    Clear()
}

// In-memory cache with TTL (time to live)
type entry struct {
    value     string
    expiresAt time.Time
}

type InMemoryCache struct {
    data map[string]entry
}

func NewInMemoryCache() *InMemoryCache
func (c *InMemoryCache) Get(key string) (string, bool)   // return false if expired
func (c *InMemoryCache) Set(key string, value string, ttl time.Duration)
func (c *InMemoryCache) Delete(key string)
func (c *InMemoryCache) Clear()

// NoCache — always misses (useful to disable caching in development/testing)
type NoCache struct{}
func (n *NoCache) Get(key string) (string, bool)                          { return "", false }
func (n *NoCache) Set(key string, value string, ttl time.Duration)        {}
func (n *NoCache) Delete(key string)                                      {}
func (n *NoCache) Clear()                                                 {}

// WeatherService uses the cache interface
type WeatherService struct {
    cache      Cache
    callCount  int // tracks how many real API calls were made
}

func NewWeatherService(cache Cache) *WeatherService

func (w *WeatherService) GetWeather(city string) string {
    key := "weather:" + city
    // 1. Check cache — if found, return cached value
    // 2. If not found: "call the API" (simulate with a formatted string + increment callCount)
    // 3. Store in cache with 5 minute TTL
    // 4. Return the result
}

func (w *WeatherService) Stats() string {
    return fmt.Sprintf("API calls made: %d", w.callCount)
}
```

**Expected output:**
```
With cache:
Karachi: Sunny, 35°C  (fetched from API)
Lahore:  Cloudy, 28°C (fetched from API)
Karachi: Sunny, 35°C  (served from cache)
Karachi: Sunny, 35°C  (served from cache)
API calls made: 2

With NoCache:
Karachi: Sunny, 35°C  (fetched from API)
Karachi: Sunny, 35°C  (fetched from API)
Karachi: Sunny, 35°C  (fetched from API)
API calls made: 3
```

---

## Group 4 — Design Patterns in Go

---

### Lab 15 — Observer Pattern (Event System)

**The scenario:** A farm monitoring system. When an animal's status changes, multiple parts of the system need to be notified — the dashboard, the health tracker, and the alert system.

**You will practice:** Observer pattern, slice of interfaces, event-driven thinking

**Build this:**

```go
type Event struct {
    Type    string // "status_changed", "weight_logged", "health_event"
    Payload map[string]string
}

// Observer — anything that wants to be notified
type Observer interface {
    OnEvent(event Event)
    Name() string
}

// EventEmitter — anything that can emit events and manage observers
type EventEmitter struct {
    observers []Observer
}

func (e *EventEmitter) Subscribe(o Observer) {
    e.observers = append(e.observers, o)
}

func (e *EventEmitter) Unsubscribe(name string) {
    // remove the observer with the given name
}

func (e *EventEmitter) Emit(event Event) {
    // call OnEvent on every observer
}

// Concrete observers
type DashboardObserver struct{}
func (d *DashboardObserver) Name() string { return "Dashboard" }
func (d *DashboardObserver) OnEvent(event Event) {
    // print: "[Dashboard] event_type: key=value key=value"
}

type HealthTrackerObserver struct{}
func (h *HealthTrackerObserver) Name() string { return "HealthTracker" }
func (h *HealthTrackerObserver) OnEvent(event Event) {
    // only care about "health_event" type — ignore others
}

type AlertObserver struct {
    threshold string // only alert on "critical" events
}
func (a *AlertObserver) Name() string { return "AlertSystem" }
func (a *AlertObserver) OnEvent(event Event) {
    // only alert if payload["severity"] == "critical"
}

// AnimalMonitor emits events when things happen
type AnimalMonitor struct {
    EventEmitter // embed to get Subscribe/Unsubscribe/Emit for free
    AnimalID string
}

func (m *AnimalMonitor) LogHealthEvent(condition, severity string) {
    m.Emit(Event{
        Type: "health_event",
        Payload: map[string]string{
            "animal_id": m.AnimalID,
            "condition": condition,
            "severity":  severity,
        },
    })
}

func (m *AnimalMonitor) LogWeight(weightKg string) {
    m.Emit(Event{
        Type:    "weight_logged",
        Payload: map[string]string{"animal_id": m.AnimalID, "weight_kg": weightKg},
    })
}
```

**Expected output:**
```
Subscribing: Dashboard, HealthTracker, AlertSystem

Event: weight_logged
  [Dashboard] weight_logged: animal_id=D-001 weight_kg=42.5
  [HealthTracker] ignoring non-health event

Event: health_event (severity: warning)
  [Dashboard] health_event: animal_id=D-001 condition=limping severity=warning
  [HealthTracker] recording health event for D-001: limping (warning)
  [AlertSystem] no alert — severity is warning, not critical

Event: health_event (severity: critical)
  [Dashboard] health_event: animal_id=D-001 condition=PPR severity=critical
  [HealthTracker] recording health event for D-001: PPR (critical)
  [AlertSystem] 🚨 CRITICAL ALERT: animal D-001 has PPR
```

---

### Lab 16 — Strategy Pattern (Pricing Engine)

**The scenario:** A sales system where pricing rules change based on season, customer type, and promotions — without changing the core order code.

**You will practice:** Strategy pattern, swapping behavior at runtime, functions as strategies

**Build this:**

```go
type PricingStrategy interface {
    Apply(basePrice float64) float64
    Description() string
}

// Implement these strategies:

type NoDiscount struct{}
func (n NoDiscount) Apply(price float64) float64   { return price }
func (n NoDiscount) Description() string           { return "No discount" }

type PercentageDiscount struct {
    Percent float64 // e.g., 10.0 for 10%
}
func (p PercentageDiscount) Apply(price float64) float64
func (p PercentageDiscount) Description() string // "10% discount"

type QurbaniSeasonDiscount struct{} // flat 20% off + free delivery (subtract 500)
func (q QurbaniSeasonDiscount) Apply(price float64) float64
func (q QurbaniSeasonDiscount) Description() string

type BulkDiscount struct {
    Quantity  int
    Threshold int     // quantity must exceed this for discount to apply
    Percent   float64
}
func (b BulkDiscount) Apply(price float64) float64
func (b BulkDiscount) Description() string // "15% bulk discount (10+ units)"

// Order uses whatever strategy it's given
type Order struct {
    ID       string
    Item     string
    Quantity int
    UnitPrice float64
    strategy  PricingStrategy
}

func NewOrder(id, item string, qty int, unitPrice float64, strategy PricingStrategy) *Order

func (o *Order) Total() float64 {
    return o.strategy.Apply(o.UnitPrice * float64(o.Quantity))
}

func (o *Order) PrintReceipt() {
    // print item, quantity, unit price, strategy description, and final total
}
```

**Expected output:**
```
Order: ORD-001
  Item:     Barbari Buck
  Qty:      1 × PKR 50,000
  Pricing:  No discount
  Total:    PKR 50,000.00

Order: ORD-002
  Item:     Barbari Buck
  Qty:      1 × PKR 50,000
  Pricing:  Qurbani Season (20% off + free delivery)
  Total:    PKR 39,500.00

Order: ORD-003
  Item:     Barseem Hay (per bag)
  Qty:      15 × PKR 800
  Pricing:  15% bulk discount (10+ units)
  Total:    PKR 10,200.00
```

---

### Lab 17 — Decorator Pattern (Middleware Chain)

**The scenario:** An HTTP-like request handler where behavior (logging, authentication, timing) is added by wrapping handlers — not by modifying them.

**You will practice:** Decorator pattern, function types as interfaces, composing behavior through wrapping

**Build this:**

```go
type Request struct {
    Method  string
    Path    string
    Headers map[string]string
    Body    string
}

type Response struct {
    StatusCode int
    Body       string
}

// Handler is a function that processes a request
type Handler func(req Request) Response

// Decorators — each wraps a Handler and returns a new Handler

// WithLogging logs the request and response
func WithLogging(logger Logger, next Handler) Handler {
    return func(req Request) Response {
        logger.Info(fmt.Sprintf("→ %s %s", req.Method, req.Path))
        resp := next(req)
        logger.Info(fmt.Sprintf("← %d %s", resp.StatusCode, req.Path))
        return resp
    }
}

// WithAuth checks for an Authorization header; returns 401 if missing
func WithAuth(validToken string, next Handler) Handler {
    return func(req Request) Response {
        token := req.Headers["Authorization"]
        if token != "Bearer "+validToken {
            return Response{StatusCode: 401, Body: "unauthorized"}
        }
        return next(req)
    }
}

// WithTiming prints how long the handler took
func WithTiming(next Handler) Handler {
    return func(req Request) Response {
        start := time.Now()
        resp := next(req)
        fmt.Printf("Request took: %v\n", time.Since(start))
        return resp
    }
}

// Your actual handler — knows nothing about logging or auth
func ProductHandler(req Request) Response {
    if req.Method == "GET" && req.Path == "/products" {
        return Response{StatusCode: 200, Body: `[{"id":"1","name":"Barbari Buck"}]`}
    }
    return Response{StatusCode: 404, Body: "not found"}
}
```

In `main`, compose handlers:
```go
handler := WithTiming(
    WithLogging(logger,
        WithAuth("secret-token",
            ProductHandler,
        ),
    ),
)

// Valid request
resp := handler(Request{
    Method:  "GET",
    Path:    "/products",
    Headers: map[string]string{"Authorization": "Bearer secret-token"},
})

// No auth token
resp2 := handler(Request{
    Method: "GET",
    Path:   "/products",
})
```

**Expected output:**
```
[INFO] → GET /products
[INFO] ← 200 /products
Request took: 41µs
Status: 200 | Body: [{"id":"1","name":"Barbari Buck"}]

[INFO] → GET /products
[INFO] ← 401 /products
Request took: 12µs
Status: 401 | Body: unauthorized
```

---

### Lab 18 — Implementing io.Reader

**The scenario:** A data transform pipeline. Build a custom `io.Reader` that uppercases all text as it's read — without loading the entire input into memory first.

**You will practice:** `io.Reader` interface, wrapping an existing reader, streaming data

```go
// Go's io.Reader interface:
// type Reader interface {
//     Read(p []byte) returns (n int, err error)
// }
// Read fills p with up to len(p) bytes and returns how many were read.
// Return io.EOF when there is no more data.
```

**Build this:**

```go
// UpperCaseReader wraps another reader and uppercases all bytes on the fly
type UpperCaseReader struct {
    source io.Reader
}

func NewUpperCaseReader(r io.Reader) *UpperCaseReader

func (u *UpperCaseReader) Read(p []byte) (int, error) {
    // 1. Read from u.source into p
    // 2. Uppercase every byte that is a lowercase letter (a-z → A-Z)
    // 3. Return n and err unchanged
    // Hint: for i := 0; i < n; i++ { if p[i] >= 'a' && p[i] <= 'z' { p[i] -= 32 } }
}

// LineCountReader wraps a reader and counts how many newlines it has seen
type LineCountReader struct {
    source io.Reader
    Lines  int
}

func NewLineCountReader(r io.Reader) *LineCountReader
func (l *LineCountReader) Read(p []byte) (int, error)
// count '\n' in each chunk you read
```

In `main`:
```go
text := "hello world\ngo is great\ninterfaces are powerful\n"

// Chain: text → uppercase → count lines → print
source := strings.NewReader(text)
lineCounter := NewLineCountReader(source)
upper := NewUpperCaseReader(lineCounter)

// Read all and print
result, _ := io.ReadAll(upper)
fmt.Println(string(result))
fmt.Printf("Line count: %d\n", lineCounter.Lines)
```

**Expected output:**
```
HELLO WORLD
GO IS GREAT
INTERFACES ARE POWERFUL

Line count: 3
```

---

### Lab 19 — Plugin Registry System

**The scenario:** A report generation system where report types are registered at startup and can be looked up and run by name — new report types can be added without touching the registry or runner code.

**You will practice:** Map of interfaces, registry pattern, factory functions

**Build this:**

```go
type ReportData struct {
    Title     string
    From      time.Time
    To        time.Time
    Records   []map[string]string
}

type ReportGenerator interface {
    Name() string
    Description() string
    Generate(data ReportData) (string, error)
}

// Registry — stores and retrieves report generators by name
type ReportRegistry struct {
    generators map[string]ReportGenerator
}

func NewReportRegistry() *ReportRegistry
func (r *ReportRegistry) Register(g ReportGenerator) error       // error if name already registered
func (r *ReportRegistry) Get(name string) (ReportGenerator, error)
func (r *ReportRegistry) List() []string                         // sorted list of registered names

// Implement three report generators:

type AnimalCountReport struct{}
func (a *AnimalCountReport) Name() string { return "animal-count" }
func (a *AnimalCountReport) Description() string { return "Count of animals by type" }
func (a *AnimalCountReport) Generate(data ReportData) (string, error) {
    // count records by data["type"], format as a table
}

type HealthSummaryReport struct{}
func (h *HealthSummaryReport) Name() string { return "health-summary" }
func (h *HealthSummaryReport) Description() string { return "Health events summary" }
func (h *HealthSummaryReport) Generate(data ReportData) (string, error) {
    // group records by data["event_type"], count each
}

type FinancialReport struct{}
func (f *FinancialReport) Name() string { return "financial" }
func (f *FinancialReport) Description() string { return "Income and expense summary" }
func (f *FinancialReport) Generate(data ReportData) (string, error) {
    // sum records where data["type"]=="income" vs "expense", compute net
}
```

In `main`:
```go
registry := NewReportRegistry()
registry.Register(&AnimalCountReport{})
registry.Register(&HealthSummaryReport{})
registry.Register(&FinancialReport{})

fmt.Println("Available reports:", registry.List())

gen, _ := registry.Get("financial")
output, _ := gen.Generate(data)
fmt.Println(output)

// Try an unknown report
_, err := registry.Get("nonexistent")
fmt.Println(err)
```

**Expected output:**
```
Available reports: [animal-count financial health-summary]

Financial Report: Jan 01 – Dec 31
  Total income:   PKR 450,000
  Total expenses: PKR 185,000
  Net profit:     PKR 265,000

Error: report generator 'nonexistent' not found
```

---

### Lab 20 — Full System: Animal Health Tracker

**The scenario:** Combine everything from the previous labs into a mini system. This is your capstone for this lab set.

**You will practice:** All previous concepts together — embedding, interfaces, DI, repository pattern, custom errors, observer pattern, strategy

**Build this:**

The system tracks animal health records on the Markhaur farm.

```go
// Domain types
type AnimalType string
const (
    Doe  AnimalType = "doe"
    Buck AnimalType = "buck"
    Kid  AnimalType = "kid"
)

type Animal struct {
    TagID  string
    Name   string
    Type   AnimalType
    WeightKg float64
}

type HealthRecord struct {
    ID        string
    AnimalTag string
    EventType string // "vaccination", "illness", "checkup"
    Notes     string
    Date      time.Time
    Cost      float64
}

// Repository interface
type HealthRepository interface {
    SaveRecord(record HealthRecord) error
    FindByAnimal(tagID string) ([]HealthRecord, error)
    FindAll() ([]HealthRecord, error)
}

// In-memory implementation
type InMemoryHealthRepo struct {
    records []HealthRecord
}
// implement all 3 methods

// Cost strategy — different events have different cost rules
type CostStrategy interface {
    CalculateCost(event string, baseAmount float64) float64
}

type StandardCost struct{}
func (s StandardCost) CalculateCost(event string, baseAmount float64) float64 {
    // vaccination: baseAmount * 1.0 (no markup)
    // illness:     baseAmount * 1.5 (50% surcharge for treatment)
    // checkup:     500.0 fixed cost regardless of baseAmount
}

// Observer for alerts
type HealthAlertObserver struct {
    logger Logger
}
func (h *HealthAlertObserver) Name() string { return "HealthAlert" }
func (h *HealthAlertObserver) OnEvent(e Event) {
    // alert on illness events
}

// The main service — wires everything together
type HealthService struct {
    EventEmitter                    // embed observer support
    repo     HealthRepository       // injected
    logger   Logger                 // injected
    strategy CostStrategy           // injected
}

func NewHealthService(repo HealthRepository, logger Logger, strategy CostStrategy) *HealthService

func (s *HealthService) RecordEvent(animal Animal, eventType, notes string, baseCost float64) error {
    // 1. Calculate cost using strategy
    // 2. Save record via repo
    // 3. Emit event via EventEmitter
    // 4. Log the action
}

func (s *HealthService) GetAnimalHistory(tagID string) ([]HealthRecord, error)

func (s *HealthService) GetTotalCost(from, to time.Time) (float64, error) {
    // sum all record costs within the date range
}
```

In `main`, wire the full system and simulate a week of farm operations:

```go
repo     := &InMemoryHealthRepo{}
logger   := NewConsoleLogger("health")
strategy := StandardCost{}

svc := NewHealthService(repo, logger, strategy)
svc.Subscribe(&HealthAlertObserver{logger: logger})

doe := Animal{TagID: "D-001", Name: "Rani", Type: Doe, WeightKg: 38}

svc.RecordEvent(doe, "vaccination", "PPR vaccine — annual", 1500)
svc.RecordEvent(doe, "checkup", "routine monthly checkup", 0)
svc.RecordEvent(doe, "illness", "limping — possible injury", 2000)

history, _ := svc.GetAnimalHistory("D-001")
fmt.Printf("\nHealth history for %s (%s):\n", doe.Name, doe.TagID)
for _, r := range history {
    fmt.Printf("  [%s] %s — PKR %.0f\n", r.EventType, r.Notes, r.Cost)
}

total, _ := svc.GetTotalCost(time.Now().Add(-7*24*time.Hour), time.Now())
fmt.Printf("\nTotal health cost this week: PKR %.0f\n", total)
```

**Expected output:**
```
[INFO] [health] recording vaccination for D-001 (Rani) — cost: PKR 1500
[INFO] [health] recording checkup for D-001 (Rani) — cost: PKR 500
[INFO] [health] recording illness for D-001 (Rani) — cost: PKR 3000
[HealthAlert] 🚨 illness event for D-001: limping — possible injury

Health history for Rani (D-001):
  [vaccination] PPR vaccine — annual — PKR 1500
  [checkup]     routine monthly checkup — PKR 500
  [illness]     limping — possible injury — PKR 3000

Total health cost this week: PKR 5000
```

---

## What You've Learned

After completing these 20 labs, you understand:

| Concept | Lab(s) |
|---|---|
| Structs, methods, constructors | 1, 2 |
| Pointer vs value receivers | 1, 3 |
| Struct embedding | 3, 4, 15, 20 |
| Interface definition and satisfaction | 2, 5 |
| `fmt.Stringer` | 5 |
| Custom error types | 6 |
| Type assertions and type switches | 7 |
| Interface segregation | 8 |
| Implementing standard library interfaces | 9, 18 |
| Dependency injection via interfaces | 10, 11, 12, 13, 14 |
| Observer pattern | 15, 20 |
| Strategy pattern | 16, 20 |
| Decorator / middleware pattern | 17 |
| `io.Reader` / streaming | 18 |
| Registry pattern | 19 |
| Combining all patterns | 20 |

---

## Go vs OOP Languages — Key Differences

| OOP (Java/C++) | Go equivalent |
|---|---|
| `class` | `struct` + methods |
| `extends` | struct embedding |
| `implements` | implicit — just have the methods |
| `abstract class` | interface |
| `private` | lowercase field/method name |
| `public` | uppercase field/method name |
| `constructor` | `NewXxx()` function by convention |
| `override` | define the same method on the outer struct |
| dependency injection framework | pass interfaces to constructors |
