package rest

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	rpcs   *map[string][]string
	lock   *sync.Mutex
}

func NewServer(healthyRPcs *map[string][]string, healthyLock *sync.Mutex) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		rpcs:   healthyRPcs,
		lock:   healthyLock,
	}

	return server
}

func (s *Server) Start() error {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "bearerToken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	s.router.POST("/", s.RPChub())

	return s.router.Run(":8080")
}

func (s *Server) RPChub() gin.HandlerFunc {
	return func(c *gin.Context) {
		chainId := c.Query("chain_id")
		if chainId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "chain_id is required",
			})
			return
		}

		s.lock.Lock()
		rpcs, ok := (*s.rpcs)[chainId]
		s.lock.Unlock()
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "chain_id not supported",
			})
			return
		}

		//write a better load balancer algorithm
		rpc := rpcs[0]

		proxyRequest(c, rpc)

	}
}

func proxyRequest(c *gin.Context, rpc string) {

	context, cancel := context.WithTimeout(c.Request.Context(), time.Second*30)
	defer cancel()

	req, err := http.NewRequestWithContext(context, c.Request.Method, rpc, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	req.Header = c.Request.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Writer.Write(body)
}
