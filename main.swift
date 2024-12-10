import Foundation

// Command line argument structure
struct CommandLineOptions {
    var n: Int = 0
    var coin: Bool = false
    var `repeat`: UInt = 1
    var repeatSet: Bool = false
    var lines: Bool = false
    var tokens: Bool = false
    var shuffle: Bool = false
    var newLine: Bool = false
    var zero: Bool = false
    var comma: Bool = false
}

let coins = ["HEADS", "tails"]

// Parse command line arguments
func parseArguments() -> CommandLineOptions {
    var options = CommandLineOptions()
    let arguments = CommandLine.arguments

    var index = 1
    while index < arguments.count {
        let arg = arguments[index]
        switch arg {
        case "-n":
            if index + 1 < arguments.count, let value = Int(arguments[index + 1]) {
                options.n = value
                index += 1
            }
        case "-c":
            options.coin = true
        case "-l":
            options.lines = true
        case "-t":
            options.tokens = true
        case "-r":
            if index + 1 < arguments.count, let value = UInt(arguments[index + 1]) {
                options.repeat = value
                options.repeatSet = true
                index += 1
            }
        case "-s", "--shuffle":
            options.shuffle = true
        case "-nl":
            options.newLine = true
        case "-0", "--null":
            options.zero = true
        case "-d":
            options.comma = true
        default:
            break
        }
        index += 1
    }

    return options
}

// Generate random values based on the options
func generateValue(options: CommandLineOptions) -> Any {
    if options.coin {
        return coins[Int.random(in: 0...1)]
    } else if options.n > 0 {
        return Int.random(in: 1...options.n)
    } else if options.lines {
        // Read lines from standard input
        var lines: [String] = []
        while let line = readLine() {
            lines.append(line)
        }
        return lines[Int.random(in: 0..<lines.count)]
    } else if options.tokens {
        // Read and tokenize input
        var tokens: [String] = []
        while let line = readLine() {
            tokens.append(contentsOf: line.split(separator: " ").map(String.init))
        }
        return tokens[Int.random(in: 0..<tokens.count)]
    }
    return ""
}

func main() {
    let options = parseArguments()

    // Validate options
    if options.repeat < 1 && !options.shuffle {
        FileHandle.standardError.write("Error: -r argument must be > 0\n".data(using: .utf8)!)
        exit(1)
    }

    // Generate results
    var results: [Any] = []

    if options.shuffle {
        var uniqueResults = Set<String>()
        let timeout = Date().addingTimeInterval(1)  // 1 second timeout

        while uniqueResults.count < Int(options.repeat) {
            if Date() > timeout {
                FileHandle.standardError.write(
                    "Error: timeout exceeded when generating results.\n".data(using: .utf8)!)
                break
            }

            let value = String(describing: generateValue(options: options))
            uniqueResults.insert(value)
        }

        results = Array(uniqueResults)
        results.shuffle()
    } else {
        for _ in 0..<options.repeat {
            results.append(generateValue(options: options))
        }
    }

    // Output results
    if options.newLine {
        results.forEach { print($0) }
    } else if options.zero {
        results.forEach { print($0, terminator: "\0") }
    } else if options.comma {
        results.forEach { print($0, terminator: ",") }
    } else {
        print(results.map { String(describing: $0) }.joined(separator: " "))
    }
}

// Run the program
main()
