package controller

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"todo/internal/model"
	"todo/internal/repository"
)

type TicketController struct {
	wsServer         *socketio.Server
	ticketRepository *repository.TicketRepository
}

func NewTicketController(wsServer *socketio.Server, ticketRepository *repository.TicketRepository) *TicketController {
	return &TicketController{wsServer, ticketRepository}
}

func (c *TicketController) SetupRoutes(api router.Party) {
	api.Get("/tickets", c.GetAll)
	api.Post("/tickets", c.Add)
	api.Put("/tickets/{id}/complete", c.Complete)
	api.Delete("/tickets/{id}/delete", c.Delete)
}

func (c *TicketController) GetAll(ctx iris.Context) {
	tickets, err := c.ticketRepository.GetAll()
	if err != nil {
		ctx.JSON(map[string]string{
			"errors": err.Error(),
		})
	}
	ctx.JSON(tickets)
}

func (c *TicketController) Add(ctx iris.Context) {
	var ticket model.Ticket
	err := ctx.ReadJSON(&ticket)
	ticket.Status = model.TICKET_STATUS_CREATED

	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket creation failure").DetailErr(err))
		return
	}

	if ticket.Title == "" {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket is empty"))
		return
	}

	err = c.ticketRepository.Create(&ticket)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket creation failure").DetailErr(err))
		return
	}

	ticketJson, err := json.Marshal(ticket)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket creation failure").DetailErr(err))
		return
	}
	println("Created Ticket: " + string(ticketJson))

	c.wsServer.BroadcastToNamespace("/", "ticket-created", string(ticketJson))

	ctx.StatusCode(iris.StatusCreated)
}

func (c *TicketController) Complete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket fail").DetailErr(err))
		return
	}
	var ticket *model.Ticket
	ticket, err = c.ticketRepository.GetById(id)
	if ticket.Status == model.TICKET_STATUS_DELITED {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("You cannot complete a deleted ticket"))
		return
	}

	ticket.Status = model.TICKET_STATUS_COMPLITED

	err = c.ticketRepository.Update(ticket)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket complete fail").DetailErr(err))
		return
	}

	ticketJson, err := json.Marshal(ticket)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket complete fail").DetailErr(err))
		return
	}
	println("Complete Ticket: " + string(ticketJson))

	c.wsServer.BroadcastToNamespace("/", "ticket-completed", string(ticketJson))

	ctx.StatusCode(iris.StatusOK)
}

func (c *TicketController) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket fail").DetailErr(err))
		return
	}
	var ticket *model.Ticket
	ticket, err = c.ticketRepository.GetById(id)

	if ticket.Status == model.TICKET_STATUS_COMPLITED {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("You cannot delete a completed ticket"))
		return
	}

	ticket.Status = model.TICKET_STATUS_DELITED

	err = c.ticketRepository.Update(ticket)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket delete fail").DetailErr(err))
		return
	}

	ticketJson, err := json.Marshal(ticket)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Ticket delete fail").DetailErr(err))
		return
	}
	println("Delete Ticket: " + string(ticketJson))

	c.wsServer.BroadcastToNamespace("/", "ticket-deleted", string(ticketJson))

	ctx.StatusCode(iris.StatusOK)
}
