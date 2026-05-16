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

    enum GeneratedValue: CustomStringConvertible {
        case string(String)
        case int(Int)

        var description: String {
            switch self {
            case let .string(value):
                return value
            case let .int(value):
                return String(value)
            }
        }
    }

    func generateValue() throws -> GeneratedValue {
        if coin {
            return .string(coins[Int.random(in: 0...1)])
        } else if let maxN = n {
            return .int(Int.random(in: 1...maxN))
        } else if lines {
            var lines: [String] = []
            while let line = readLine() {
                lines.append(line)
            }
            guard !lines.isEmpty else { throw linesEmptyError }
            return .string(lines[Int.random(in: 0..<lines.count)])
        } else if tokens {
            var tokens: [String] = []
            while let line = readLine() {
                tokens.append(contentsOf: line.split(separator: " ").map(String.init))
            }
            guard !tokens.isEmpty else { throw tokensEmptyError }
            return .string(tokens[Int.random(in: 0..<tokens.count)])
        }
        throw ValidationError("Invalid arguments provided")
    }

    mutating func run() throws {
        // Validate options
        if repeatCount == 0 && !shuffle {
            throw ValidationError("Error: repeatCount must be > 0")
        }
        if [newLine, zero, delimiter].filter({ $0 }).count > 1 {
            throw ValidationError("Error: -nl, -0, and delimiter output are mutually exclusive")
        }

        // Generate results
        var results: [GeneratedValue] = []

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

                let value = try generateValue()
                uniqueResults.insert(value.description)
            }

            results = uniqueResults.map { .string($0) }
            results.shuffle()
        } else {
            for _ in 0..<repeatCount {
                results.append(try generateValue())
            }
        }

        // Output results
        if newLine {
            results.forEach { print($0) }
        } else if zero {
            results.forEach { print($0, terminator: "\0") }
        } else if delimiter {
            print(results.map(\.description).joined(separator: ","))
        } else {
            print(results.map(\.description).joined(separator: " "))
        }
    }
}
