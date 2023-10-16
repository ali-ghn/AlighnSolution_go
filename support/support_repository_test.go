package support

import (
	"context"
	"fmt"
	"testing"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	sr     SupportRepository
	client *mongo.Client
)

func init() {
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	sr = NewSupportRepository(client)
}

func TestCreateTicket(t *testing.T) {
	bId, err := primitive.ObjectIDFromHex("63397973d4f9905057ebfe95")
	if err != nil {
		t.Error(err)
	}
	ticket := Ticket{
		Title:    "Test Ticket",
		AdminId:  bId,
		SenderId: primitive.NewObjectID(),
		Status:   shared.TicketStatusNew,
		TicketContents: []TicketContent{
			{
				Content:  "Hello this is content for test",
				Id:       primitive.NewObjectID(),
				SenderId: primitive.NewObjectID(),
			},
		},
	}
	res, err := sr.CreateTicket(&ticket)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res.Id)
}

func TestGetTicket(t *testing.T) {
	ticketId := "63397973d4f9905057ebfe99"
	ticket, err := sr.GetTicket(ticketId)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ticket.Id)
}

func TestGetTickets(t *testing.T) {
	filter := bson.D{{Key: "title", Value: "Test Ticket"}}
	tickets, err := sr.GetTickets(filter, nil)
	if err != nil {
		t.Error(err)
	}
	for _, v := range *tickets {
		fmt.Println(v.Id)
	}
}

func TestGetTicketsByStatus(t *testing.T) {

	bId, err := primitive.ObjectIDFromHex("63397973d4f9905057ebfe95")
	if err != nil {
		t.Error(err)
	}
	filter := bson.D{{"$and",
		bson.A{
			bson.D{{
				Key: "adminid", Value: bId,
			}},
			bson.D{{
				"$or", bson.A{
					bson.D{{"status", shared.TicketStatusNew}},
					bson.D{{"status", shared.TicketStatusInProgress}},
					bson.D{{"status", shared.TicketStatusPending}},
					bson.D{{"status", shared.TicketStatusSupportResponse}},
					bson.D{{"status", shared.TicketStatusUserResponse}},
				},
			}},
		},
	}}
	// options := options.Find().SetLimit(10)
	tickets, err := sr.GetTickets(filter, nil)
	if err != nil {
		t.Error(err)
	}
	for _, v := range *tickets {
		fmt.Println(v.Title)
	}
}

func TestTicketCount(t *testing.T) {

	bId, err := primitive.ObjectIDFromHex("63397973d4f9905057ebfe95")
	if err != nil {
		t.Error(err)
	}
	filter := bson.D{{"$and",
		bson.A{
			bson.D{{
				Key: "adminid", Value: bId,
			}},
			bson.D{{
				"$or", bson.A{
					bson.D{{"status", shared.TicketStatusNew}},
					bson.D{{"status", shared.TicketStatusInProgress}},
					bson.D{{"status", shared.TicketStatusPending}},
					bson.D{{"status", shared.TicketStatusSupportResponse}},
					bson.D{{"status", shared.TicketStatusUserResponse}},
				},
			}},
		},
	}}
	// options := options.Find().SetLimit(10)
	count, err := sr.GetTicketCount(filter)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}
