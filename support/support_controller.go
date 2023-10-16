package support

import (
	"fmt"
	"net/http"

	"github.com/ali-ghn/Coinopay_Go/attachment"
	"github.com/ali-ghn/Coinopay_Go/auth"
	"github.com/ali-ghn/Coinopay_Go/shared"
	"github.com/ali-ghn/Coinopay_Go/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupportController struct {
	sr   ISupportRepository
	auth auth.IAuth
	ar   attachment.IAttachmentRepository
	ur   user.IUserRepository
}

func NewSupportController(sr SupportRepository, auth auth.Auth,
	ar attachment.AttachmentRepository, ur user.UserRepository) SupportController {
	return SupportController{
		sr:   sr,
		auth: auth,
		ar:   ar,
		ur:   ur,
	}
}

func (sc SupportController) CreateTicket(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	var cTicketReq CreateTicketRequest
	err := c.Bind(&cTicketReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := sc.ur.GetUserByEmail(claims.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}

	filter := bson.D{{Key: "roles", Value: shared.SUPPORT_ROLE}}
	supports, err := sc.ur.GetUsers(filter, nil)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	freestUser := user.User{}
	supportCount := 0
	for _, v := range *supports {
		filter := bson.D{{Key: "$and",
			Value: bson.A{
				bson.D{{
					Key: "adminid", Value: v.Id,
				}},
				bson.D{{
					Key: "$or", Value: bson.A{
						bson.D{{Key: "status", Value: shared.TicketStatusNew}},
						bson.D{{Key: "status", Value: shared.TicketStatusInProgress}},
						bson.D{{Key: "status", Value: shared.TicketStatusPending}},
						bson.D{{Key: "status", Value: shared.TicketStatusSupportResponse}},
						bson.D{{Key: "status", Value: shared.TicketStatusUserResponse}},
					},
				}},
			},
		}}
		ticketCount, err := sc.sr.GetTicketCount(filter)
		if err != nil {
			c.String(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}
		if supportCount == 0 {
			freestUser = v
		} else if ticketCount < int64(supportCount) {
			freestUser = v
		}
	}

	ticket := Ticket{
		Title:    cTicketReq.Title,
		SenderId: cUser.Id,
		AdminId:  freestUser.Id,
		Status:   shared.TicketStatusNew,
		TicketContents: []TicketContent{
			{
				Id:          primitive.NewObjectID(),
				Content:     cTicketReq.Content,
				SenderId:    cUser.Id,
				Attachments: cTicketReq.Attachments,
			},
		},
	}

	res, err := sc.sr.CreateTicket(&ticket)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	cResponse := CreateTicketResponse{
		Id:      res.Id.Hex(),
		Title:   res.Title,
		Status:  res.Status,
		Content: cTicketReq.Content,
	}
	c.JSON(http.StatusCreated, cResponse)
}

func (sc SupportController) GetTicket(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := sc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	var getTicketRequest GetTicketRequest
	err = c.BindJSON(&getTicketRequest)

	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}

	ticket, err := sc.sr.GetTicket(getTicketRequest.Id)

	if err != nil {
		c.String(http.StatusBadRequest, "Ticket doesn't exist")
		return
	}
	if ticket.AdminId == cUser.Id || ticket.SenderId == cUser.Id {
		getTicketResponse := GetTicketResponse{
			Id:             ticket.Id.Hex(),
			Title:          ticket.Title,
			Status:         ticket.Status,
			TicketContents: ticket.TicketContents,
		}
		c.JSON(http.StatusOK, getTicketResponse)
		return
	}

	c.String(http.StatusForbidden, "You don't have access to this resource")
}

func (sc SupportController) GetTickets(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := sc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	if err != nil {
		return
	}
	filter := bson.D{{Key: "$or",
		Value: bson.A{
			bson.D{{Key: "adminid", Value: cUser.Id}},
			bson.D{{Key: "senderid", Value: cUser.Id}},
		},
	}}
	tickets, err := sc.sr.GetTickets(filter, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	resTickets := []GetTicketResponse{}
	for _, v := range *tickets {
		resTickets = append(resTickets, GetTicketResponse{
			Id:             v.Id.Hex(),
			Title:          v.Title,
			Status:         v.Status,
			TicketContents: v.TicketContents,
		})
	}
	c.JSON(http.StatusOK, resTickets)
}

func (sc SupportController) CreateTicketContent(c *gin.Context) {
	var ticketContentReq CreateTicketContentRequest
	err := c.Bind(&ticketContentReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := sc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	ticket, err := sc.sr.GetTicket(ticketContentReq.TicketId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	if ticket.SenderId == cUser.Id || ticket.AdminId == cUser.Id {
		ticket.TicketContents = append(ticket.TicketContents, TicketContent{
			Id:          primitive.NewObjectID(),
			Content:     ticketContentReq.Content,
			SenderId:    cUser.Id,
			Attachments: ticketContentReq.Attachments,
		})
		_, err := sc.sr.UpdateTicket(ticket)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}
		c.String(http.StatusOK, "Ticket content has created")
		return
	} else {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
}
