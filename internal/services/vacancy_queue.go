package services

type VacancyQueue interface {
	Enqueue(vacancyID string) error
	Dequeue() (string, error)
	Len() int
}
