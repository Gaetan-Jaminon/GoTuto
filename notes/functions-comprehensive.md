# Go Functions - Comprehensive Guide

## Basic Function Declaration

### Syntax
```go
func functionName(parameters) returnType {
    // function body
    return value
}
```

### Simple Examples
```go
// No parameters, no return
func sayHello() {
    fmt.Println("Hello, World!")
}

// With parameters
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// With return value
func add(a, b int) int {
    return a + b
}

// Multiple parameters of same type (shorthand)
func multiply(x, y, z int) int {
    return x * y * z
}
```

## Function Parameters

### Different Parameter Types
```go
// Mixed types
func createUser(name string, age int, isActive bool) {
    fmt.Printf("User: %s, Age: %d, Active: %t\n", name, age, isActive)
}

// Variadic parameters (accepts multiple arguments)
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

// Calling variadic function
result := sum(1, 2, 3, 4, 5)        // Pass individual values
slice := []int{1, 2, 3, 4, 5}
result2 := sum(slice...)            // Spread slice with ...
```

## Return Values

### Single Return
```go
func getFullName(first, last string) string {
    return first + " " + last
}
```

### Multiple Returns (Very Common in Go)
```go
// Multiple unnamed returns
func divmod(a, b int) (int, int) {
    return a / b, a % b
}

// Usage
quotient, remainder := divmod(10, 3)
fmt.Printf("10 / 3 = %d remainder %d\n", quotient, remainder)
```

### Named Return Values
```go
func calculate(a, b int) (sum, product int) {
    sum = a + b
    product = a * b
    return // naked return - returns named values
}

// Or you can still return explicitly
func calculateExplicit(a, b int) (sum, product int) {
    sum = a + b
    product = a * b
    return sum, product
}
```

## Error Handling in Go

### Standard Error Pattern
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Usage with error handling
func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Result: %.2f\n", result)
}
```

### Custom Error Types
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Message: "cannot be negative",
        }
    }
    if age > 150 {
        return ValidationError{
            Field:   "age", 
            Message: "unrealistic value",
        }
    }
    return nil
}
```

## If/Else Statements

### Basic If/Else
```go
func checkAge(age int) string {
    if age < 18 {
        return "Minor"
    } else if age < 65 {
        return "Adult"
    } else {
        return "Senior"
    }
}
```

### If with Short Statement (Go idiom)
```go
func processFile(filename string) error {
    // Declare and check in same line
    if file, err := os.Open(filename); err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    } else {
        defer file.Close() // file is available in this scope
        // process file...
        return nil
    }
}

// More common pattern - early return
func processFilePreferred(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    // process file...
    return nil
}
```

### Error Handling with If
```go
func readUserData(id int) (*User, error) {
    user, err := database.GetUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    
    if user.IsDeleted {
        return nil, errors.New("user has been deleted")
    }
    
    return user, nil
}
```

## Switch Statements

### Basic Switch
```go
func getDay(day int) string {
    switch day {
    case 1:
        return "Monday"
    case 2:
        return "Tuesday"
    case 3:
        return "Wednesday"
    case 4:
        return "Thursday"
    case 5:
        return "Friday"
    case 6, 7: // Multiple values
        return "Weekend"
    default:
        return "Invalid day"
    }
}
```

### Switch with Expressions
```go
func categorizeScore(score int) string {
    switch {
    case score >= 90:
        return "Excellent"
    case score >= 80:
        return "Good"
    case score >= 70:
        return "Average"
    case score >= 60:
        return "Below Average"
    default:
        return "Poor"
    }
}
```

### Switch on Types (Type Switch)
```go
func handleValue(value interface{}) string {
    switch v := value.(type) {
    case string:
        return fmt.Sprintf("String: %s", v)
    case int:
        return fmt.Sprintf("Integer: %d", v)
    case bool:
        return fmt.Sprintf("Boolean: %t", v)
    case nil:
        return "Nil value"
    default:
        return fmt.Sprintf("Unknown type: %T", v)
    }
}
```

## Advanced Function Features

### Functions as First-Class Citizens
```go
// Function type
type Operation func(int, int) int

// Functions as variables
var add Operation = func(a, b int) int {
    return a + b
}

var multiply Operation = func(a, b int) int {
    return a * b
}

// Function as parameter
func calculate(a, b int, op Operation) int {
    return op(a, b)
}

// Usage
result1 := calculate(5, 3, add)      // 8
result2 := calculate(5, 3, multiply) // 15
```

### Anonymous Functions and Closures
```go
func createCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// Usage
counter1 := createCounter()
counter2 := createCounter()

fmt.Println(counter1()) // 1
fmt.Println(counter1()) // 2
fmt.Println(counter2()) // 1 (separate closure)
```

### Defer Statement
```go
func processData() error {
    file, err := os.Open("data.txt")
    if err != nil {
        return err
    }
    defer file.Close() // Executes when function returns
    
    mutex.Lock()
    defer mutex.Unlock() // Multiple defers executed in LIFO order
    
    // Process data...
    // file.Close() and mutex.Unlock() called automatically
    return nil
}
```

### Method vs Function
```go
// Function
func calculateArea(radius float64) float64 {
    return math.Pi * radius * radius
}

// Method (function with receiver)
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Usage
area1 := calculateArea(5.0)     // Function call
circle := Circle{Radius: 5.0}
area2 := circle.Area()          // Method call
```

## Function Visibility and Scope

### Exported vs Unexported
```go
// Exported (public) - starts with capital letter
func PublicFunction() {
    // Can be called from other packages
}

// Unexported (private) - starts with lowercase
func privateFunction() {
    // Can only be called within same package
}
```

## Common Patterns and Best Practices

### Early Return Pattern
```go
func validateUser(user *User) error {
    if user == nil {
        return errors.New("user cannot be nil")
    }
    
    if user.Name == "" {
        return errors.New("name is required")
    }
    
    if user.Email == "" {
        return errors.New("email is required")
    }
    
    // All validations passed
    return nil
}
```

### Builder Pattern with Functions
```go
type Config struct {
    Host    string
    Port    int
    Timeout time.Duration
}

type ConfigOption func(*Config)

func WithHost(host string) ConfigOption {
    return func(c *Config) {
        c.Host = host
    }
}

func WithPort(port int) ConfigOption {
    return func(c *Config) {
        c.Port = port
    }
}

func NewConfig(options ...ConfigOption) *Config {
    config := &Config{
        Host:    "localhost",
        Port:    8080,
        Timeout: 30 * time.Second,
    }
    
    for _, option := range options {
        option(config)
    }
    
    return config
}

// Usage
config := NewConfig(
    WithHost("example.com"),
    WithPort(9000),
)
```

## Coming from .NET

### Function Declaration Comparison
```csharp
// C#
public int Add(int a, int b)
{
    return a + b;
}

public (int sum, int product) Calculate(int a, int b)
{
    return (a + b, a * b);
}
```

```go
// Go
func Add(a, b int) int {
    return a + b
}

func Calculate(a, b int) (int, int) {
    return a + b, a * b
}
```

### Error Handling Comparison
```csharp
// C# - Exceptions
try 
{
    var result = Divide(10, 0);
    Console.WriteLine(result);
}
catch (DivideByZeroException ex)
{
    Console.WriteLine($"Error: {ex.Message}");
}
```

```go
// Go - Error values
result, err := Divide(10, 0)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
fmt.Println(result)
```

## Key Takeaways

1. **Multiple returns are idiomatic** - especially for error handling
2. **Error handling is explicit** - no exceptions, use error values
3. **Functions are first-class citizens** - can be passed around
4. **Defer is powerful** - for cleanup and resource management
5. **Early returns improve readability** - avoid deep nesting
6. **Capital letters = exported** - visibility is determined by naming
7. **Switch doesn't fall through** - no break needed (use fallthrough if needed)
8. **If statements can include initialization** - `if err := doSomething(); err != nil`