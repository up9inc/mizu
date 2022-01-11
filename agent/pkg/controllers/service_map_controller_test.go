package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/up9inc/mizu/shared"
)

type ServiceMapControllerSuite struct {
	suite.Suite

	c *ServiceMapController
	w *httptest.ResponseRecorder
	g *gin.Context
}

func (s *ServiceMapControllerSuite) SetupTest() {
	s.c = NewServiceMapController()
	s.c.service.SetConfig(&shared.MizuAgentConfig{
		ServiceMap: true,
	})
	s.c.service.AddEdge("a", "b", "p")

	s.w = httptest.NewRecorder()
	s.g, _ = gin.CreateTestContext(s.w)
}

func (s *ServiceMapControllerSuite) TestGetStatus() {
	assert := s.Assert()

	s.c.Status(s.g)
	assert.Equal(http.StatusOK, s.w.Code)

	var status shared.ServiceMapStatus
	err := json.Unmarshal(s.w.Body.Bytes(), &status)
	assert.NoError(err)
	assert.Equal("enabled", status.Status)
	assert.Equal(1, status.EntriesProcessedCount)
	assert.Equal(2, status.NodeCount)
	assert.Equal(1, status.EdgeCount)
}

func (s *ServiceMapControllerSuite) TestGet() {
	assert := s.Assert()

	s.c.Get(s.g)
	assert.Equal(http.StatusOK, s.w.Code)

	var response shared.ServiceMapResponse
	err := json.Unmarshal(s.w.Body.Bytes(), &response)
	assert.NoError(err)

	// response status
	assert.Equal("enabled", response.Status.Status)
	assert.Equal(1, response.Status.EntriesProcessedCount)
	assert.Equal(2, response.Status.NodeCount)
	assert.Equal(1, response.Status.EdgeCount)

	// response nodes
	aNode := shared.ServiceMapNode{
		Name:     "a",
		Id:       1,
		Protocol: "p",
		Count:    1,
	}
	bNode := shared.ServiceMapNode{
		Name:     "b",
		Id:       2,
		Protocol: "p",
		Count:    1,
	}
	assert.Equal([]shared.ServiceMapNode{
		aNode,
		bNode,
	}, response.Nodes)

	// response edges
	assert.Equal([]shared.ServiceMapEdge{
		{
			Source:      aNode,
			Destination: bNode,
			Count:       1,
		},
	}, response.Edges)
}

func (s *ServiceMapControllerSuite) TestGetReset() {
	assert := s.Assert()

	s.c.Reset(s.g)
	assert.Equal(http.StatusOK, s.w.Code)

	var status shared.ServiceMapStatus
	err := json.Unmarshal(s.w.Body.Bytes(), &status)
	assert.NoError(err)
	assert.Equal("enabled", status.Status)
	assert.Equal(0, status.EntriesProcessedCount)
	assert.Equal(0, status.NodeCount)
	assert.Equal(0, status.EdgeCount)
}

func TestServiceMapControllerSuite(t *testing.T) {
	suite.Run(t, new(ServiceMapControllerSuite))
}
