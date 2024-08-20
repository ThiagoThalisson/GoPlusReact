package api

import (
  "net/http"
  "log/slog"
  "errors"

  "github.com/thiagothalisson/goplusreact/internal/store/pgstore"
  "github.com/go-chi/cors"
  "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
  "github.com/google/uuid"

  "github.com/jackc/pgx/v5"
  
  "github.com/gorilla/websocket"
)


type apiHandler struct {
  q *pgstore.Queries
  r *chi.Mux
  upgrader websocket.Upgrader
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  h.r.ServeHTTP(w, r)

}

func NewHandler(q *pgstore.Queries) http.Handler {
  a := apiHandler {
    q: q,
    upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
  }
  
  r := chi.NewRouter()
  r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)
  
  r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"https://*", "http://*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: false,
    MaxAge:           300,
  }))

  r.Get("/subscribe/{room_id}", a.handleSubscribe)
  
  r.Route("/api", func(r chi.Router) {
    r.Route("/rooms", func(r chi.Router) {
      r.Post("/", a.handleCreateRoom)
      r.Get("/", a.handleGetRooms)
      
      r.Route("/{room_id}/messages", func(r chi.Router){
        r.Post("/", a.handleCreateRoomMessage)
        r.Get("/", a.handleGetRoomMessages)
        
        r.Route("/{message_id}", func(r chi.Router){
          r.Get("/", a.handleGetRoomMessage)
          r.Patch("/react", a.handleReactToMessage)
          r.Delete("/react", a.handleRemoveReactFromMessage)
          r.Patch("/answer", a.handleMarkMessageAsAnswered)
        })
      }) 
      
    })
  })

  a.r = r
  return a
}

func (h apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
  rawRoomID := chi.URLParam(r, "room_id")
  roomID, err := uuid.Parse(rawRoomID)
  
  if err != nil {
    http.Error(w, "Invalid room id", http.StatusBadGateway)
    return
  }
  
  _, err = h.q.GetRoom(r.Context(), roomID)
  if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
      http.Error(w, "room not found", http.StatusBadRequest)
      return
    }
    
    http.Error(w, "something went wrong", http.StatusInternalServerError)
    return
  }
  
  c, err := h.upgrader.Upgrade(w, r, nil)
  if err != nil {
    slog.Warn("failed to upgrade connection", "error", err)
    http.Error(w, "failed to upgrade to ws connection", http.StatusBadRequest)
    return
  }
  
  defer c.Close()
}

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {

}

func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleMarkMessageAsAnswered(w http.ResponseWriter, r *http.Request) {}


