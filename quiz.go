package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "log"
    "math/rand"
    "os"
    "strings"
    "time"
)

type Quiz struct {
    question string
    answer   string
}

func main() {
    // Command-line flags for file name, time limit, and shuffling
    filename := flag.String("csv", "problems.csv", "a CSV file in 'question,answer' format")
    timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
    shuffle := flag.Bool("shuffle", false, "shuffle the quiz questions")
    flag.Parse()

    // Read and parse the CSV file
    questions, err := parseCSV(*filename)
    if err != nil {
        log.Fatalf("Failed to parse the provided CSV file: %v\n", err)
    }

    // Shuffle questions if the shuffle flag is set
    if *shuffle {
        rand.Seed(time.Now().UnixNano())
        rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
    }

    // Run the quiz
    correct := runQuiz(questions, *timeLimit)
    fmt.Printf("\nYou scored %d out of %d.\n", correct, len(questions))
}

// parseCSV reads the CSV file and returns a slice of Quiz structs.
func parseCSV(filename string) ([]Quiz, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    lines, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    var questions []Quiz
    for _, line := range lines {
        questions = append(questions, Quiz{
            question: line[0],
            answer:   strings.TrimSpace(line[1]),
        })
    }
    return questions, nil
}

// runQuiz administers the quiz to the user with a time limit.
func runQuiz(questions []Quiz, timeLimit int) int {
    timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
    defer timer.Stop()

    fmt.Println("Press Enter to start the quiz.")
    fmt.Scanln() // Wait for user input to start

    correct := 0
    answerCh := make(chan string)

    for i, q := range questions {
        fmt.Printf("Problem #%d: %s = ", i+1, q.question)

        go func() {
            var answer string
            fmt.Scanln(&answer)
            answerCh <- answer
        }()

        select {
        case <-timer.C:
            fmt.Println("\nTime's up!")
            return correct
        case answer := <-answerCh:
            // Check the answer, ignoring case and trimming spaces
            if strings.TrimSpace(strings.ToLower(answer)) == strings.ToLower(q.answer) {
                correct++
            }
        }
    }
    return correct
}
											 