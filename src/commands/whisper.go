package commands

// func sendWhisper() {
// 	// Check if it's a whisper command
// 	if strings.HasPrefix(v, "/whisper") {
// 		whisper := strings.TrimPrefix(v, "/whisper ")
// 		parts := strings.SplitN(whisper, " ", 2)
//
// 		if len(parts) < 2 {
// 			break
// 		}
//
// 		// Set the target user and content for whisper
// 		userMessage.To = parts[0]
// 		userMessage.Content = parts[1]
//
// 		newMessage := fmt.Sprintf("%s (Whisper to %s): %s ", userMessage.Timestamp, userMessage.To, parts[1])
//
// 		m.MessageList.Messages = append(m.MessageList.Messages, newMessage)
// 		m.MessageList.Count++
//
// 		// if there are more messages than limit pop the oldest from array
// 		if m.MessageList.Count > messageLimit {
// 			m.MessageList.Messages = m.MessageList.Messages[1:]
// 			m.MessageList.Count--
// 		}
//
// 		m.Viewport.SetContent(strings.Join(m.MessageList.Messages, "\n"))
// 		m.Viewport.GotoBottom()
//
// 		//	m.placeMessage(newMessage)
// 	} else {
// 		// Normal message
// 		userMessage.To = "all"
// 		userMessage.Content = v
// 	}
// }
