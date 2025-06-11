package services

import "errors"

type SliceVacancyQueue struct {
	queue []string
}

func NewSliceVacancyQueue() *SliceVacancyQueue {
	return &SliceVacancyQueue{queue: make([]string, 0)}
}

func (q *SliceVacancyQueue) Enqueue(vacancyID string) error {
	q.queue = append(q.queue, vacancyID)
	return nil
}

func (q *SliceVacancyQueue) Dequeue() (string, error) {
	if len(q.queue) == 0 {
		return "", errors.New("queue is empty")
	}
	vacancyID := q.queue[0]
	q.queue[0] = "" // зануляем для GC
	q.queue = q.queue[1:]
	return vacancyID, nil
}

func (q *SliceVacancyQueue) Len() int {
	return len(q.queue)
}
