package handlers

import (
	"log"

	"github.com/aglili/auction-app/internal/utils"
	"github.com/aglili/auction-app/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	connManager *websocket.ConnectionManager
}

func NewWebSocketHandler(connManager *websocket.ConnectionManager) *WebSocketHandler {
	return &WebSocketHandler{
		connManager: connManager,
	}
}

func (h *WebSocketHandler) HandleWSConnections(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.RespondWithError(ctx, err, "invalid session")
		return
	}

	conn, err := websocket.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		utils.RespondWithError(ctx, err, "failed to create websocket connection")
		return
	}

	h.connManager.Register(uid, conn)
	defer h.connManager.UnRegister(uid)

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if messageType == ws.PingMessage {
			conn.WriteMessage(ws.PongMessage, nil)
		}
	}
}
