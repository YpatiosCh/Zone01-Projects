package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

type Client struct {
	Name     string
	Conn     net.Conn
	Messages chan string
	Group    string
}

type MessageHistory struct {
	Messages []string
	MaxSize  int
}

type Server struct {
	Clients        map[*Client]bool
	ClientsMutex   sync.Mutex
	Groups         map[string]map[*Client]bool
	GroupHistory   map[string]*MessageHistory
	LogFile        *os.File
	HistoryMaxSize int
}

var (
	server        *Server
	logUpdateCh   = make(chan string, 100)
	groupUpdateCh = make(chan struct{}, 10)
	userUpdateCh  = make(chan struct{}, 10)
	icon          = ` 
    Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |     ".      | "" \Zq
_)      \.___.,|     .'
"\____   )MMMMMP|   .'
     "-'       "--'
    
	Please enter your name: 
     `
)

const (
	sidebarWidth  = 30
	inputHeight   = 3
	statusHeight  = 2
	minMainHeight = 10
	minMainWidth  = 50
)

func NewServer(logFilePath string) *Server {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	server := &Server{
		Clients:        make(map[*Client]bool),
		Groups:         make(map[string]map[*Client]bool),
		GroupHistory:   make(map[string]*MessageHistory),
		LogFile:        logFile,
		HistoryMaxSize: 100, // Keep last 100 messages per group
	}

	// Add default groups with message history
	defaultGroups := []string{"general", "tech", "business"}
	for _, group := range defaultGroups {
		server.Groups[group] = make(map[*Client]bool)
		server.GroupHistory[group] = &MessageHistory{
			Messages: make([]string, 0),
			MaxSize:  server.HistoryMaxSize,
		}
	}

	return server
}

func (s *Server) SendGroupHistory(client *Client) {
	s.ClientsMutex.Lock()
	history := s.GroupHistory[client.Group]
	s.ClientsMutex.Unlock()

	if history != nil {

		for _, msg := range history.Messages {
			client.Messages <- msg
		}

	}
}

func (s *Server) AddMessageToHistory(group, message string) {
	s.ClientsMutex.Lock()
	defer s.ClientsMutex.Unlock()

	if s.GroupHistory[group] == nil {
		s.GroupHistory[group] = &MessageHistory{
			Messages: make([]string, 0),
			MaxSize:  s.HistoryMaxSize,
		}
	}

	history := s.GroupHistory[group]
	history.Messages = append(history.Messages, message)

	// Keep only the last MaxSize messages
	if len(history.Messages) > history.MaxSize {
		history.Messages = history.Messages[len(history.Messages)-history.MaxSize:]
	}
}

func main() {
	port := "8989"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	server = NewServer("server.log")
	defer server.LogFile.Close()

	// Start the server in a goroutine
	go StartServer(server, port)

	// Start the UI
	if err := StartUI(); err != nil {
		if err != gocui.ErrQuit {
			log.Fatalf("Failed to start UI: %v", err)
		}
	}
}

func StartServer(server *Server, port string) {

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	LogMessage(server, fmt.Sprintf("Server started on port %s", port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogMessage(server, fmt.Sprintf("Failed to accept connection: %v", err))
			continue
		}

		client := &Client{
			Name:     "",
			Conn:     conn,
			Messages: make(chan string, 10),
		}

		go HandleClient(server, client)
	}
}

func GetName(client *Client) string {
	client.Conn.Write([]byte(icon))
	reader := bufio.NewReader(client.Conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read name: %v", err)
		return ""
	}
	return strings.TrimSpace(name)
}

func GetGroupSelection(server *Server, client *Client) string {
	client.Conn.Write([]byte("Type < /help > when in chat to see your options\n"))
	client.Conn.Write([]byte("Available groups:\n - general\n - tech\n - business\n"))
	client.Conn.Write([]byte("Type the name of an existing group or create a new one: "))

	reader := bufio.NewReader(client.Conn)
	group, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read group: %v", err)
		return ""
	}
	group = strings.TrimSpace(group)

	server.ClientsMutex.Lock()
	defer server.ClientsMutex.Unlock()

	if _, exists := server.Groups[group]; !exists {
		server.Groups[group] = make(map[*Client]bool)
		LogMessage(server, fmt.Sprintf("Group '%s' created by %s", group, client.Name))
	}

	return group
}

func SendMessages(client *Client) {
	for msg := range client.Messages {
		client.Conn.Write([]byte(msg + "\n"))
	}
}

func ReceiveMessages(server *Server, client *Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		command, content := ParseCommand(msg)
		switch command {
		case "/join":
			HandleJoinCommand(server, client, content)
		case "/list":
			HandleListCommand(server, client)
		case "/help":
			HandleHelpCommand(client)
		case "":
			timestamp := time.Now().Format("15:04:05")
			formattedMsg := fmt.Sprintf("[%s][%s][%s]: %s", timestamp, client.Group, client.Name, strings.TrimSpace(msg))
			LogMessage(server, formattedMsg)
			BroadcastToGroup(server, client.Group, formattedMsg)
			server.AddMessageToHistory(client.Group, formattedMsg)
		}
	}
}

func HandleListCommand(server *Server, client *Client) {
	server.ClientsMutex.Lock()
	defer server.ClientsMutex.Unlock()

	response := "Available groups and members:\n"
	for groupName, clients := range server.Groups {
		response += fmt.Sprintf("\n📁 %s (%d members):\n", groupName, len(clients))
		for member := range clients {
			response += fmt.Sprintf("  - %s\n", member.Name)
		}
	}
	client.Messages <- response
}

func HandleHelpCommand(client *Client) {
	helpMsg := `Available commands:
/join <group> - Join or create a new group
/list - List all groups and their members
/help - Show this help message`
	client.Messages <- helpMsg
}

func ParseCommand(msg string) (string, string) {
	msg = strings.TrimSpace(msg)
	if strings.HasPrefix(msg, "/") {
		parts := strings.SplitN(msg, " ", 2)
		if len(parts) == 2 {
			return parts[0], strings.TrimSpace(parts[1])
		}
		return parts[0], ""
	}
	return "", msg
}

func HandleJoinCommand(server *Server, client *Client, groupName string) {
	if groupName == "" {
		client.Messages <- "Usage: /join <group_name>"
		return
	}

	server.ClientsMutex.Lock()
	if server.Groups[groupName] == nil {
		server.Groups[groupName] = make(map[*Client]bool)
		server.GroupHistory[groupName] = &MessageHistory{
			Messages: make([]string, 0),
			MaxSize:  server.HistoryMaxSize,
		}
		LogMessage(server, fmt.Sprintf("Group '%s' created by %s", groupName, client.Name))
	}
	server.ClientsMutex.Unlock()

	oldGroup := client.Group
	leaveMsg := fmt.Sprintf("🔄 %s has left the group to join %s", client.Name, groupName)
	BroadcastToGroup(server, oldGroup, leaveMsg)
	server.AddMessageToHistory(oldGroup, leaveMsg)

	server.ClientsMutex.Lock()
	if server.Groups[oldGroup] != nil {
		delete(server.Groups[oldGroup], client)
	}
	server.Groups[groupName][client] = true
	client.Group = groupName
	server.ClientsMutex.Unlock()

	joinMsg := fmt.Sprintf("🔄 %s has joined the group from %s", client.Name, oldGroup)
	BroadcastToGroup(server, groupName, joinMsg)
	server.AddMessageToHistory(groupName, joinMsg)

	LogMessage(server, fmt.Sprintf("%s switched from group %s to %s", client.Name, oldGroup, groupName))

	select {
	case client.Messages <- fmt.Sprintf("You joined group: %s", groupName):
	default:
		log.Printf("Warning: Could not send join confirmation to %s", client.Name)
	}

	select {
	case userUpdateCh <- struct{}{}:
	default:
	}
	select {
	case groupUpdateCh <- struct{}{}:
	default:
	}
}

func BroadcastToGroup(server *Server, group string, msg string) {
	if group == "" {
		return
	}

	server.ClientsMutex.Lock()
	clients := make([]*Client, 0)
	if groupClients, ok := server.Groups[group]; ok {
		for client := range groupClients {
			clients = append(clients, client)
		}
	}
	server.ClientsMutex.Unlock()

	// Send messages outside the lock
	for _, client := range clients {
		select {
		case client.Messages <- msg:
		default:
			// If the channel is full, log the error but don't block
			log.Printf("Warning: Could not send message to %s", client.Name)
		}
	}
}

func LogMessage(server *Server, msg string) {
	log.Printf("%s\n", msg)
	server.LogFile.WriteString(msg + "\n")

	select {
	case logUpdateCh <- msg:
	default:
	}
}

func StartUI() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.Mouse = true

	g.SetManagerFunc(layout)

	// Key bindings
	if err := SetKeybindings(g); err != nil {
		return err
	}

	// Start UI update goroutines
	go UpdateLogs(g)
	go UpdateGroups(g)
	go UpdateUsers(g)

	return g.MainLoop()
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Create status bar
	if v, err := g.SetView("status", -1, maxY-statusHeight, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.BgColor = gocui.ColorBlue
		v.FgColor = gocui.ColorWhite
		fmt.Fprintf(v, " TCP Chat Server - Port: 8989 - Press Ctrl+C to quit")
	}

	// Create groups list
	if v, err := g.SetView("groups", 0, 0, sidebarWidth, maxY-statusHeight-inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Groups"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	// Create users list
	if v, err := g.SetView("users", maxX-sidebarWidth, 0, maxX-1, maxY-statusHeight-inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Users"
		v.Wrap = true
	}

	// Create main chat/logs area
	if v, err := g.SetView("logs", sidebarWidth+1, 0, maxX-sidebarWidth-1, maxY-statusHeight-inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat Logs"
		v.Wrap = true
		v.Autoscroll = true
	}

	// Create input field
	inputY := maxY - statusHeight - inputHeight
	if v, err := g.SetView("input", sidebarWidth+1, inputY, maxX-sidebarWidth-1, maxY-statusHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input (Enter to send)"
		v.Wrap = true
		v.Editable = true
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	return nil
}

func SetKeybindings(g *gocui.Gui) error {
	// Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}

	// Send message
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, SendMessage); err != nil {
		return err
	}

	// Switch between views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, NextView); err != nil {
		return err
	}

	// Groups navigation
	if err := g.SetKeybinding("groups", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("groups", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}

	return nil
}

func UpdateLogs(g *gocui.Gui) {
	for msg := range logUpdateCh {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("logs")
			if err != nil {
				return err
			}
			fmt.Fprintln(v, msg)
			return nil
		})
	}
}

func SendMessage(g *gocui.Gui, v *gocui.View) error {
	input := strings.TrimSpace(v.Buffer())
	if input != "" {
		timestamp := time.Now().Format("15:04:05")
		serverMsg := fmt.Sprintf("[%s][SERVER]: %s", timestamp, input)

		// Broadcast to all groups
		server.ClientsMutex.Lock()
		for groupName := range server.Groups {
			for client := range server.Groups[groupName] {
				select {
				case client.Messages <- fmt.Sprintf("\x1b[31m%s\x1b[0m", serverMsg):
				default:
					log.Printf("Warning: Could not send server message to %s", client.Name)
				}
			}
		}
		server.ClientsMutex.Unlock()

		// Log the server message
		LogMessage(server, serverMsg)

		// Clear the input
		v.Clear()
		v.SetCursor(0, 0)
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func NextView(g *gocui.Gui, v *gocui.View) error {
	nextViewName := ""
	if v == nil || v.Name() == "groups" {
		nextViewName = "input"
	} else {
		nextViewName = "groups"
	}

	if _, err := g.SetCurrentView(nextViewName); err != nil {
		return err
	}

	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func HandleClient(server *Server, client *Client) {
	defer func() {
		client.Conn.Close()
		close(client.Messages)

		server.ClientsMutex.Lock()
		delete(server.Clients, client)
		if server.Groups[client.Group] != nil {
			delete(server.Groups[client.Group], client)

			if len(server.Groups[client.Group]) == 0 &&
				client.Group != "general" &&
				client.Group != "tech" &&
				client.Group != "business" {
				delete(server.Groups, client.Group)
				delete(server.GroupHistory, client.Group)
			}
		}
		server.ClientsMutex.Unlock()

		leaveMsg := fmt.Sprintf("🔴 %s has left the chat", client.Name)
		BroadcastToGroup(server, client.Group, leaveMsg)
		server.AddMessageToHistory(client.Group, leaveMsg)
		LogMessage(server, fmt.Sprintf("%s left the group %s", client.Name, client.Group))

		select {
		case userUpdateCh <- struct{}{}:
		default:
		}
		select {
		case groupUpdateCh <- struct{}{}:
		default:
		}
	}()

	client.Name = GetName(client)
	if client.Name == "" {
		return
	}

	client.Group = GetGroupSelection(server, client)
	if client.Group == "" {
		return
	}

	server.ClientsMutex.Lock()
	server.Groups[client.Group][client] = true
	server.Clients[client] = true
	server.ClientsMutex.Unlock()

	// Send message history to new client only when they first connect
	server.SendGroupHistory(client)

	select {
	case userUpdateCh <- struct{}{}:
	default:
	}
	select {
	case groupUpdateCh <- struct{}{}:
	default:
	}

	joinMsg := fmt.Sprintf("🟢 %s has joined the chat", client.Name)
	BroadcastToGroup(server, client.Group, joinMsg)
	server.AddMessageToHistory(client.Group, joinMsg)
	LogMessage(server, fmt.Sprintf("%s joined group %s", client.Name, client.Group))

	go SendMessages(client)
	ReceiveMessages(server, client)
}

func UpdateUsers(g *gocui.Gui) {
	for range userUpdateCh {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("users")
			if err != nil {
				return err
			}
			v.Clear()

			server.ClientsMutex.Lock()
			activeClients := make([]*Client, 0, len(server.Clients))
			for client := range server.Clients {
				activeClients = append(activeClients, client)
			}
			server.ClientsMutex.Unlock()

			// Sort clients by name for consistent display
			sort.Slice(activeClients, func(i, j int) bool {
				return activeClients[i].Name < activeClients[j].Name
			})

			for _, client := range activeClients {
				fmt.Fprintf(v, "🟢 %s (%s)\n", client.Name, client.Group)
			}
			return nil
		})
	}
}

func UpdateGroups(g *gocui.Gui) {
	for range groupUpdateCh {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("groups")
			if err != nil {
				return err
			}
			v.Clear()

			server.ClientsMutex.Lock()
			// Create a sorted list of group names for consistent display
			groupNames := make([]string, 0, len(server.Groups))
			for groupName := range server.Groups {
				if len(server.Groups[groupName]) > 0 ||
					groupName == "general" ||
					groupName == "tech" ||
					groupName == "business" {
					groupNames = append(groupNames, groupName)
				}
			}

			// Sort group names
			sort.Strings(groupNames)

			// Display groups
			for _, groupName := range groupNames {
				clients := server.Groups[groupName]
				fmt.Fprintf(v, "📁 %s (%d)\n", groupName, len(clients))
			}
			server.ClientsMutex.Unlock()
			return nil
		})
	}
}
