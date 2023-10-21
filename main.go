package main

import (
	"github.com/ali-ghn/AlighnSolution_go/attachment"
	"github.com/ali-ghn/AlighnSolution_go/auth"
	"github.com/ali-ghn/AlighnSolution_go/blog"
	"github.com/ali-ghn/AlighnSolution_go/cors"
	"github.com/ali-ghn/AlighnSolution_go/cryptography"
	"github.com/ali-ghn/AlighnSolution_go/email"
	"github.com/ali-ghn/AlighnSolution_go/invoice"
	"github.com/ali-ghn/AlighnSolution_go/paymentProcessor"
	"github.com/ali-ghn/AlighnSolution_go/services"
	"github.com/ali-ghn/AlighnSolution_go/settings"
	"github.com/ali-ghn/AlighnSolution_go/store"
	"github.com/ali-ghn/AlighnSolution_go/support"
	"github.com/ali-ghn/AlighnSolution_go/user"
	"github.com/ali-ghn/AlighnSolution_go/wallet"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var routerFactoryState []*gin.Engine

func init() {

}

func main() {
	configHelper, err := services.NewConfigHelper("appSettings.json")
	if err != nil {
		panic(err)
	}
	// Site settings
	key, err := configHelper.GetSection("Key")
	if err != nil {
		panic(err)
	}
	// Smtp Settings
	smtpHost, err := configHelper.GetSection("SmtpHost")
	if err != nil {
		panic(err)
	}
	smtpPassword, err := configHelper.GetSection("SmtpPassword")
	if err != nil {
		panic(err)
	}
	// Payment Processor Settings
	ppToken, err := configHelper.GetSection("PaymentProcessorToken")
	if err != nil {
		panic(err)
	}
	ppHost, err := configHelper.GetSection("PaymentProcessorHost")
	if err != nil {
		panic(err)
	}
	ppStoreId, err := configHelper.GetSection("PaymentProcessorStoreId")
	if err != nil {
		panic(err)
	}

	uc := user.NewUserController(
		user.NewUserRepository(client),
		cryptography.NewEncryptionHelper(key.(string)),
		auth.NewAuth([]byte(key.(string))),
		email.NewEmailSender(smtpHost.(string), smtpPassword.(string)))

	sc := store.NewStoreController(auth.NewAuth([]byte(key.(string))),
		store.NewStoreRepository(client), user.NewUserRepository(client))

	ic := invoice.NewInvoiceController(paymentProcessor.NewPaymentProcessor(ppToken.(string),
		ppHost.(string), ppStoreId.(string), *resty.New()),
		invoice.NewInvoiceRepository(client), store.NewStoreRepository(client),
		auth.NewAuth([]byte(key.(string))))

	suc := support.NewSupportController(support.NewSupportRepository(client),
		auth.NewAuth([]byte(key.(string))),
		attachment.NewAttachmentRepository(client),
		user.NewUserRepository(client))

	stc := settings.NewSettingsController(settings.NewSettingsRepository(client),
		auth.NewAuth([]byte(key.(string))), user.NewUserRepository(client))

	bc := blog.NewBlogController(blog.NewBlogRepository(client),
		auth.NewAuth([]byte(key.(string))),
		user.NewUserRepository(client))

	wc := wallet.NewWalletController(wallet.NewWalletHelper(ppToken.(string), ppHost.(string), ppStoreId.(string), *resty.New()),
		wallet.NewWalletRepository(client), user.NewUserRepository(client), auth.NewAuth([]byte(key.(string))))

	// fmt.Println(tel.Bot.Client)
	r := gin.Default()
	// CORS middleware
	r.Use(cors.CORSMiddleware())
	// User controller
	r.POST("/signUp", uc.SignUp)
	r.POST("/signIn", uc.SignIn)
	r.POST("/verifyEmail", uc.VerifyEmail)
	r.POST("/forgetPassword", uc.ForgetPassword)
	r.POST("/verifyForgetPassword", uc.VerifyForgetPassword)
	r.POST("/user", uc.InsertUser)
	r.GET("/user", uc.GetUser)
	r.GET("/users", uc.GetUsers)
	r.GET("/users/role", uc.GetUserByRole)
	r.PUT("/user", uc.UpdateUser)
	// Store Controller
	r.POST("/store", sc.CreateStore)
	r.GET("/store", sc.GetStore)
	r.GET("/stores", sc.GetStores)
	r.GET("/stores/user", sc.GetStoresByUser)
	// Invoice Controller
	r.POST("/invoice", ic.CreateInvoice)
	r.GET("/invoice", ic.GetInvoice)
	r.GET("/invoices", ic.GetInvoices)
	// Support Controller
	r.POST("/ticket", suc.CreateTicket)
	r.GET("/ticket", suc.GetTicket)
	r.GET("/tickets", suc.GetTickets)
	r.POST("/ticketContent", suc.CreateTicketContent)
	// Settings Controller
	r.POST("/siteSettings", stc.CreateSiteSettings)
	r.GET("/siteSettings", stc.GetSiteSettings)
	// Blog Controller
	r.POST("/blogPost", bc.CreateBlogPost)
	r.GET("/blogPost", bc.GetBlogPost)
	r.GET("/blogPosts/admin", bc.GetBlogPostsByAdmin)
	r.GET("/blogPosts", bc.GetPublishedBlogPosts)
	// Wallet Controller
	r.GET("/wallet", wc.GetWalletOverview)
	r.GET("/wallet/address", wc.GetWalletAddress)
	r.GET("/wallets", wc.GetWallets)
	r.Run(":8082")
}
