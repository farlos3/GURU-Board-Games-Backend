package boardgame

import (
	"github.com/gofiber/fiber/v2"
	"guru-game/models"
)

// CreateBoardGame handler
func CreateBoardGame(c *fiber.Ctx) error {
	newGame := new(models.BoardGame)
	if err := c.BodyParser(newGame); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	boardgame, err := CreateNewBoardGame(newGame)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Boardgame already exists"})
	}

	return c.JSON(fiber.Map{"message": "Boardgame created", "boardgame": boardgame})
}

// GetBoardGameByName handler
func GetBoardGameByName(c *fiber.Ctx) error {
	name := c.Params("name")

	boardgame, err := FindBoardGameByName(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Boardgame not found"})
	}

	return c.JSON(boardgame)
}

// GetAllBoardGames handler
func GetAllBoardGames(c *fiber.Ctx) error {
	boardgames, err := GetAllBoardgames()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get boardgames"})
	}
	return c.JSON(boardgames)
}

