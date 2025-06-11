package services

import (
	"container/list"
	"errors"
)

type ListVacancyQueue struct {
	queue *list.List
}

func NewListVacancyQueue() *ListVacancyQueue {
	return &ListVacancyQueue{queue: list.New()}
}

func (q *ListVacancyQueue) Enqueue(vacancyID string) error {
	q.queue.PushBack(vacancyID)
	return nil
}

func (q *ListVacancyQueue) Dequeue() (string, error) {
	front := q.queue.Front()
	if front == nil {
		return "", errors.New("queue is empty")
	}
	vacancyID := front.Value.(string)
	q.queue.Remove(front)
	return vacancyID, nil
}

func (q *ListVacancyQueue) Len() int {
	return q.queue.Len()
}
