package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/up9inc/mizu/shared"
	"mizuserver/pkg/utils"
)

var TapStatus shared.TapStatus

func GetTappingStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(TapStatus)
}


func AnalyzeInformation(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&shared.AnalyzeStatus{
		IsAnalyzing: IsAnalyzing,
		RemoteUrl: utils.GetRemoteUrl(AnalyzeDestination, AnalyzeToken),
		IsRemoteReady: utils.CheckIfModelReady(AnalyzeDestination, AnalyzedModel, AnalyzeToken),
	})
}