package model

import "time"

const TICKET_STATUS_CREATED = 1
const TICKET_STATUS_COMPLITED = 2
const TICKET_STATUS_DELITED = 3

type Ticket struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
