package dto

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status 		int;
	Message 	string;
	Data		interface{};
}

func (r *Response) JSON() fiber.Map {
	return fiber.Map{
		"status": r.Status,
		"message": r.Message,
		"data": r.Data,
	}
}