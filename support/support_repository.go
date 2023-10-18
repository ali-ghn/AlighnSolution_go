package support

import (
	"context"

	"github.com/ali-ghn/AlighnSolution_go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "Tickets"
)

type ISupportRepository interface {
	CreateTicket(ticket *Ticket) (*Ticket, error)
	GetTicket(id string) (*Ticket, error)
	GetTickets(filter interface{}, options *options.FindOptions) (*[]Ticket, error)
	GetTicketCount(filter interface{}) (int64, error)
	UpdateTicket(ticket *Ticket) (*Ticket, error)
}

type SupportRepository struct {
	Client *mongo.Client
}

func NewSupportRepository(client *mongo.Client) SupportRepository {
	return SupportRepository{
		Client: client,
	}
}

func (sr SupportRepository) CreateTicket(ticket *Ticket) (*Ticket, error) {
	ticket.Id = primitive.NewObjectID()
	res, err := sr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).InsertOne(context.TODO(), ticket)
	if err != nil {
		return nil, err
	}
	ticket.Id = res.InsertedID.(primitive.ObjectID)
	return ticket, nil
}

func (sr SupportRepository) GetTicket(id string) (*Ticket, error) {
	bId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: bId}}
	var ticket Ticket
	err = sr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (sr SupportRepository) GetTickets(filter interface{}, options *options.FindOptions) (*[]Ticket, error) {
	var tickets []Ticket
	cur, err := sr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &tickets)
	if err != nil {
		return nil, err
	}
	return &tickets, nil
}

func (sr SupportRepository) GetTicketCount(filter interface{}) (int64, error) {
	count, err := sr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).CountDocuments(context.TODO(), filter)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (sr SupportRepository) UpdateTicket(ticket *Ticket) (*Ticket, error) {
	filter := bson.D{{Key: "_id", Value: ticket.Id}}
	err := sr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOneAndReplace(context.TODO(), filter, ticket).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
