package support

import (
	"github.com/ali-ghn/AlighnSolution_go/attachment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	Id             primitive.ObjectID `bson:"_id"`
	Title          string
	AdminId        primitive.ObjectID
	SenderId       primitive.ObjectID
	Status         string
	TicketContents []TicketContent
}

type TicketContent struct {
	Id          primitive.ObjectID
	Content     string
	SenderId    primitive.ObjectID
	Attachments []attachment.Attachment
}

type CreateTicketRequest struct {
	Title       string
	Content     string
	Attachments []attachment.Attachment
}

type CreateTicketResponse struct {
	Id      string
	Title   string
	Status  string
	Content string
}

type GetTicketRequest struct {
	Id string
}

type GetTicketResponse struct {
	Id             string
	Title          string
	Status         string
	TicketContents []TicketContent
}

type CreateTicketContentRequest struct {
	Content     string
	TicketId    string
	Attachments []attachment.Attachment
}
