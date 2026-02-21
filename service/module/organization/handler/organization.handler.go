package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/dto"
)

func (h *OrganizationHandler) GetOrganizations(ctx fiber.Ctx) error {
	// Parse pagination dari query
	pagination := helper.ParsePaginationFromQuery(ctx)

	// Implementation for handling the request to get organizations
	organizations, totalRecords, sysError := h.OrganizationUsecase.GetOrganizations(ctx)
	if sysError != nil {
		return helper.SendErrorResponse(ctx, sysError.GetStatusCode(), sysError.GetMessage(), sysError.GetError())
	}

	// Return paginated response
	return helper.SendPaginatedResponse(ctx, fiber.StatusOK, "Organizations retrieved successfully",
		pagination.Page, pagination.PageSize, totalRecords, organizations)
}

// GetOrganizationByID - Get single organization by ID
// @GET /organizations/:id
func (h *OrganizationHandler) GetOrganizationByID(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return helper.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid organization ID", err)
	}

	res, sysErr := h.OrganizationUsecase.GetOrganizationByID(ctx, id)
	if sysErr != nil {
		return helper.SendErrorResponse(ctx, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(ctx, fiber.StatusOK, "Organization retrieved successfully", res)
}

// CreateOrganization - Create new organization
// @POST /organizations
func (h *OrganizationHandler) CreateOrganization(ctx fiber.Ctx) error {
	var req dto.CreateOrganizationRequest
	if validationErrors, err := helper.ValidateRequest(ctx, &req); err != nil {
		return helper.SendResponse(ctx, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	res, sysErr := h.OrganizationUsecase.CreateOrganization(ctx, req)
	if sysErr != nil {
		return helper.SendErrorResponse(ctx, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(ctx, fiber.StatusCreated, "Organization created successfully", res)
}

// UpdateOrganization - Update existing organization
// @PUT /organizations/:id
func (h *OrganizationHandler) UpdateOrganization(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return helper.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid organization ID", err)
	}

	var req dto.UpdateOrganizationRequest

	if validationErrors, err := helper.ValidateRequest(ctx, &req); err != nil {
		return helper.SendResponse(ctx, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	res, sysErr := h.OrganizationUsecase.UpdateOrganization(ctx, id, req)
	if sysErr != nil {
		return helper.SendErrorResponse(ctx, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(ctx, fiber.StatusOK, "Organization updated successfully", res)
}

// DeleteOrganization - Delete organization
// @DELETE /organizations/:id
func (h *OrganizationHandler) DeleteOrganization(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return helper.SendErrorResponse(ctx, fiber.StatusBadRequest, "Invalid organization ID", err)
	}

	sysErr := h.OrganizationUsecase.DeleteOrganization(ctx, id)
	if sysErr != nil {
		return helper.SendErrorResponse(ctx, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(ctx, fiber.StatusOK, "Organization deleted successfully", nil)
}
