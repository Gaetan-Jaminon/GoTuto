# Go's Short Variable Declaration (:=)

## What is `:=`?

The `:=` operator is Go's **short variable declaration**. It declares and initializes a variable in one line, with Go automatically inferring the type.

## Syntax

```go
variableName := value
```

## Examples

### Basic Usage
```go
name := "John"        // string
age := 25             // int
isActive := true      // bool
price := 19.99        // float64
```

### Compare with Long Form
```go
// Short form (preferred)
name := "John"

// Long form (equivalent)
var name string
name = "John"

// Or even longer
var name string = "John"
```

## Multiple Variable Declaration

```go
// Multiple variables at once
name, age := "John", 25

// With function returns
result, err := someFunction()

// Common pattern with errors
file, err := os.Open("data.txt")
if err != nil {
    // handle error
}
```

## Important Rules

### 1. Only for New Variables
```go
name := "John"    // ✅ OK - declaring new variable
name := "Jane"    // ❌ ERROR - name already declared
name = "Jane"     // ✅ OK - reassigning existing variable
```

### 2. At Least One New Variable (Multiple Assignment)
```go
name := "John"
age := 25

// This works - err is new even though name exists
name, err := getUser()  // ✅ OK

// This fails - both variables already exist
name, age := getNameAge()  // ❌ ERROR
```

### 3. Function Scope Only
```go
var globalVar string  // ✅ OK at package level

// globalShort := "test"  // ❌ ERROR - := not allowed at package level

func main() {
    localVar := "test"  // ✅ OK inside function
}
```

## When to Use `:=` vs `var`

### Use `:=` when:
- Inside functions
- You want type inference
- You have an initial value
- Most common case

```go
func main() {
    name := "John"           // Preferred
    users := make([]User, 0) // Preferred
    count := len(items)      // Preferred
}
```

### Use `var` when:
- Package-level variables
- You need to declare without initializing
- You want to be explicit about type
- Zero value initialization

```go
var globalCounter int  // Package level

func example() {
    var result string     // Will be "" (zero value)
    var numbers []int     // Will be nil
    var user User         // Will be zero value of User struct
}
```

## Common Patterns

### Error Handling
```go
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()
```

### Type Assertion
```go
value, ok := interfaceValue.(string)
if !ok {
    // handle type assertion failure
}
```

### Map Lookup
```go
value, exists := myMap["key"]
if !exists {
    // key doesn't exist
}
```

### Channel Operations
```go
data, more := <-channel
if !more {
    // channel is closed
}
```

## `:=` vs `=` - Key Difference

**`:=` (Short Declaration)**
- **Declares** a new variable AND assigns a value
- Can only be used inside functions
- Go infers the type automatically

**`=` (Assignment)**
- **Assigns** to an existing variable
- Can be used anywhere
- Variable must already exist

### Examples
```go
// := declares a NEW variable
name := "John"     // Creates new variable 'name'

// = assigns to EXISTING variable  
name = "Jane"      // Changes existing 'name' variable

// This would ERROR:
name := "Bob"      // ❌ Can't redeclare 'name'
```

## Coming from .NET

### C# Pattern:
```csharp
var name = "John";      // Declare with type inference
name = "Jane";          // Assign to existing variable
string name = "John";   // Explicit type declaration
```

### Go Equivalent:
```go
name := "John"              // Short declaration (similar to C# var)
name = "Jane"               // Assignment (same as C#)
var name string = "John"    // Explicit declaration (similar to C# explicit)
```

### Key Difference from C#:
- **C#**: `var` is just type inference, but still uses `=`
- **Go**: `:=` means "declare AND assign", `=` means "assign only"

## Key Takeaways

1. **`:=` = declare + assign + infer type**
2. **Only inside functions**
3. **At least one variable must be new (in multiple assignment)**
4. **Preferred way for local variables with initial values**
5. **Cannot redeclare the same variable in the same scope**