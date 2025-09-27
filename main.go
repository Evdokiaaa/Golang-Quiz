package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)




type problem struct {
	question string
	answer   string
}
type wrongProblem struct {
	problem
	userAnswer string
}

const (
	timeLimitSec = 20
)


func ErrorWrapper(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func getAndReadFile() (records [][]string, err error ) {
	defer func () { err = ErrorWrapper("Error reading file", err) }()

	//Человек вводит go run . -cvs problems.csv и должен получить по файл нейму
	fileName := flag.String("csv","problems.csv","Имя файла с данными")
	flag.Parse();
	if filepath.Ext(*fileName) != ".csv"{
		fmt.Println("Файл должен быть в формате csv")
		os.Exit(1)
	}
	//Открытие и чтение файла 
	file, err := os.Open(*fileName) //Открыли файла
	err = ErrorWrapper("Error opening file", err)
	r := csv.NewReader(file) //Создали считыватель

	lines, err := r.ReadAll() //Читали файл. Массив строк
	err = ErrorWrapper("Error reading file", err)
	
	return lines, err
}



func main(){
	
	timeLimit := flag.Int("time",timeLimitSec,"Время на тест в секундах")
	records, err := getAndReadFile()
	if err != nil {
		fmt.Println(err)
	}
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second) //создаем таймер и превращаем в секунды
	showResults(records,timer)

	
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem,len(lines)) // [{2+2,4}, {1+1,2}, {8+3,11}]
	for i, line := range lines {
		problems[i] = problem {
			question:line[0],
			answer:strings.TrimSpace(line[1]),
		}
	}
	return problems
}

//Counter - количество правильных ответов. Идет в возращаемых значениях
func showResults(records [][]string, timer *time.Timer)  {
	problems := parseLines(records)
	wrongProblems := make([]wrongProblem,len(problems))
	counter := 0

	problemLoop:
	for id, problem := range problems {
		fmt.Printf("Problem #%d; Problem question:  %s - " , id+1, problem.question)
		ansCh := make(chan string)

		go func(){
			var answer string;
			fmt.Scanf("%s\n",&answer)
			ansCh <- answer
		}()
		select {
			//В канале будет какое что значение, когда ввели. Присвоим его перемееной
			case ans := <- ansCh:
				if ans == problem.answer {
					counter++
				} else {
					wrongProblems = append(wrongProblems, wrongProblem{problem, ans})
					
				}
			case <-timer.C:
				fmt.Println()
				fmt.Println("\nВремя истекло")
				break problemLoop
		}
	}
	
	fmt.Printf("Всего правильных ответов: %d из %d \n",counter,len(problems))

	fmt.Println("------------------------------")
	

	id := 1
	for _, w := range wrongProblems {
		//Если пользователь дал неправильный ответ, то только выводим его
		if w.userAnswer != "" {
			fmt.Printf("Вопрос %d: %s\n", id ,w.question)
			fmt.Printf("Правильный ответ: %s\n",w.answer)
			fmt.Println("Ваш ответ:" + w.userAnswer)

			fmt.Println("--------------------")
			id++
		}
			
	}
	
}



