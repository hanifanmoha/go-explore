package behaviouralpatterns

import "fmt"

// https://refactoring.guru/design-patterns/command/go/example

// Client -> Invoker -> Command -> Receiver

// ====== COMMAND INTERFACE =====

type Command interface {
	Execute() error
}

// ====== COMMAND IMPLEMENTATION =====

type CreateOrderCommand struct {
	Receiver *ReceiverDB
	OrderID  string
	Quantity int
}

func (c *CreateOrderCommand) Execute() error {
	c.Receiver.Save(c.OrderID, c.Quantity)
	return nil
}

type DeleteOrderCommand struct {
	Receiver *ReceiverDB
	OrderID  string
}

func (c *DeleteOrderCommand) Execute() error {
	c.Receiver.Delete(c.OrderID)
	return nil
}

type SendEmailCommand struct {
	Receiver *EmailService
	To       string
	Subject  string
	Body     string
}

func (c *SendEmailCommand) Execute() error {
	c.Receiver.Send(c.To, c.Subject, c.Body)
	return nil
}

// ====== INVOKER =====

type Invoker struct {
	commands []Command
}

func (i *Invoker) AddCommand(cmd Command) {
	i.commands = append(i.commands, cmd)
}

func (i *Invoker) ExecuteCommands() error {
	for _, cmd := range i.commands {
		if err := cmd.Execute(); err != nil {
			return err
		}
	}
	return nil
}

// ====== RECEIVER =====

type EmailService struct{}

func (e *EmailService) Send(to, subject, body string) {
	fmt.Printf("Sending email \nTo: %s\nSubject: %s\nBody: %s\n", to, subject, body)
}

type ReceiverDB struct {
	orders map[string]int
}

func (r *ReceiverDB) Save(orderID string, quantity int) {
	r.orders[orderID] = quantity
}

func (r *ReceiverDB) Delete(orderID string) {
	delete(r.orders, orderID)
}

// ====== CLIENT =====

type Client struct {
	DB           *ReceiverDB
	EmailService *EmailService
	// CreateOrderInvoker *Invoker
	// DeleteOrderInvoker *Invoker
}

func NewClient() *Client {
	return &Client{
		DB:           &ReceiverDB{orders: make(map[string]int)},
		EmailService: &EmailService{},
	}
}

func (c *Client) CreateOrder(orderID string, quantity int) {
	invoker := &Invoker{}
	invoker.AddCommand(&CreateOrderCommand{
		Receiver: c.DB,
		OrderID:  orderID,
		Quantity: quantity,
	})
	invoker.AddCommand(&SendEmailCommand{
		Receiver: c.EmailService,
		To:       "customer@example.com",
		Subject:  "Order Confirmation",
		Body:     fmt.Sprintf("Your order %s has been created.", orderID),
	})
	invoker.ExecuteCommands()
}

// questions: is invoker cannot receive params? how to make invoker reusable for different commands with different params?
