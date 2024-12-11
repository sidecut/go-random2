import ArgumentParser
import Foundation

@main
struct Random: ParsableCommand {
    static var configuration: CommandConfiguration {
        CommandConfiguration(
            commandName: "random",
            abstract: "Generate random numbers, flip coins, or select random items"
        )
    }

    var tokensEmptyError: Error {
        ValidationError("No tokens were provided")
    }

    var linesEmptyError: Error {
        ValidationError("No lines were provided")
    }

    @Option(name: .shortAndLong, help: "Generate a random number between 1 and N")
    var n: Int?

    @Flag(name: .shortAndLong, help: "Flip a coin (HEADS or tails)")
    var coin = false

    @Option(name: .shortAndLong, help: "Repeat the operation N times")
    var repeatCount: UInt = 1

    @Flag(name: .shortAndLong, help: "Read lines from stdin and select one randomly")
    var lines = false

    @Flag(
        name: .shortAndLong, help: "Read space-separated tokens from stdin and select one randomly")
    var tokens = false

    @Flag(name: .shortAndLong, help: "Ensure unique results when using -r")
    var shuffle = false

    @Flag(name: [.customLong("nl"), .long], help: "Print each result on a new line")
    var newLine = false

    @Flag(name: [.customShort("0"), .long], help: "Print results with null terminator")
    var zero = false

    @Flag(name: .shortAndLong, help: "Print results with comma separator")
    var delimiter = false

    var coins = ["HEADS", "tails"]

    var invalidArgsError: Error {
        ValidationError("Invalid arguments provided")
    }

    func generateValue() -> Any {
        if coin {
            return coins[Int.random(in: 0...1)]
        } else if let maxN = n {
            return Int.random(in: 1...maxN)
        } else if lines {
            // Read lines from standard input
            var lines: [String] = []
            while let line = readLine() {
                lines.append(line)
            }
            guard !lines.isEmpty else { Random.exit(withError: linesEmptyError) }
            return lines[Int.random(in: 0..<lines.count)]
        } else if tokens {
            // Read and tokenize input
            var tokens: [String] = []
            while let line = readLine() {
                tokens.append(contentsOf: line.split(separator: " ").map(String.init))
            }
            guard !tokens.isEmpty else { Random.exit(withError: tokensEmptyError) }
            return tokens[Int.random(in: 0..<tokens.count)]
        }
        Random.exit(withError: invalidArgsError)
    }

    mutating func run() throws {
        // Validate options
        if repeatCount < 1 && !shuffle {
            throw ValidationError("Error: -r argument must be > 0")
        }

        // Generate results
        var results: [Any] = []

        if shuffle {
            var uniqueResults = Set<String>()
            let timeout = Date().addingTimeInterval(1)  // 1 second timeout

            while uniqueResults.count < Int(repeatCount) {
                if Date() > timeout {
                    let warningMessage =
                        "Warning: timeout exceeded when generating results. The repeat count may be too high.\n"
                    FileHandle.standardError.write(Data(warningMessage.utf8))
                    break
                }

                let value = String(describing: generateValue())
                uniqueResults.insert(value)
            }

            results = Array(uniqueResults)
            results.shuffle()
        } else {
            for _ in 0..<repeatCount {
                results.append(generateValue())
            }
        }

        // Output results
        if newLine {
            results.forEach { print($0) }
        } else if zero {
            results.forEach { print($0, terminator: "\0") }
        } else if delimiter {
            results.forEach { print($0, terminator: ",") }
        } else {
            print(results.map { String(describing: $0) }.joined(separator: " "))
        }
    }
}
