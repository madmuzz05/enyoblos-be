package usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/service/module/user/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/user/entity"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecase) CreateUser(ctx *fiber.Ctx, req dto.CreateUserRequest) (res dto.GetUserResponse, sysError syserror.SysError) {

	// Start transaction
	tx := u.mainDB.TxCreate(ctx)
	defer func() {
		u.mainDB.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	// Validate email tidak sudah terdaftar
	issetEmailUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err.GetStatusCode() != fiber.StatusNotFound {
		sysError = err
		return
	}

	if issetEmailUser != (entity.User{}) {
		sysError = syserror.CreateError(fiber.ErrBadRequest, fiber.StatusBadRequest, "Email sudah terdaftar")
		return
	}

	// Hash password
	hashedPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if bcryptErr != nil {
		sysError = syserror.CreateError(bcryptErr, fiber.StatusInternalServerError, bcryptErr.Error())
		return
	}

	// Create user
	model, repoErr := u.userRepo.CreateUser(ctx, entity.User{
		Name:           req.Name,
		ShortName:      req.ShortName,
		Email:          req.Email,
		Age:            req.Age,
		Password:       string(hashedPassword),
		OrganizationID: req.OrganizationID,
	})
	if repoErr != nil {
		sysError = repoErr
		return
	}

	// Retrieve full user data to include in response
	res = dto.GetUserResponse{
		ID:             model.ID,
		Name:           model.Name,
		ShortName:      model.ShortName,
		Email:          model.Email,
		Age:            model.Age,
		OrganizationID: model.OrganizationID,
	}

	// Fetch organization if exists
	organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, res.ID)
	if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
		sysError = orgErr
		return
	}
	res.Organization = &organization

	return
}

func (u *UserUsecase) GetUserByID(ctx *fiber.Ctx, Id string) (res dto.GetUserResponse, sysError syserror.SysError) {
	user, repoErr := u.userRepo.GetUserByID(ctx, Id)
	if repoErr != nil {
		sysError = repoErr
		return
	}

	res = dto.GetUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		ShortName:      user.ShortName,
		Email:          user.Email,
		Age:            user.Age,
		OrganizationID: user.OrganizationID,
	}

	// Fetch organization if exists
	organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, res.ID)
	if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
		sysError = orgErr
		return
	}
	res.Organization = &organization

	return
}
func (u *UserUsecase) GetUsers(ctx *fiber.Ctx) (res []dto.GetUserResponse, totalRecords int64, sysError syserror.SysError) {
	users, total, repoErr := u.userRepo.GetUsers(ctx)
	if repoErr != nil {
		sysError = repoErr
		return
	}

	for _, user := range users {
		userDto := dto.GetUserResponse{
			ID:             user.ID,
			Name:           user.Name,
			ShortName:      user.ShortName,
			Email:          user.Email,
			Age:            user.Age,
			OrganizationID: user.OrganizationID,
		}

		// Fetch organization if exists
		organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, user.ID)
		if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
			sysError = orgErr
			return
		}
		userDto.Organization = &organization

		res = append(res, userDto)
	}

	totalRecords = total
	return
}

func (u *UserUsecase) DeleteUser(ctx *fiber.Ctx, Id string) (sysError syserror.SysError) {
	tx := u.mainDB.TxCreate(ctx)
	defer func() {
		u.mainDB.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	sysError = u.userRepo.DeleteUser(ctx, Id)
	return
}
func (u *UserUsecase) UpdateUser(ctx *fiber.Ctx, Id string, req dto.CreateUserRequest) (res dto.GetUserResponse, sysError syserror.SysError) {
	tx := u.mainDB.TxCreate(ctx)
	defer func() {
		u.mainDB.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	// Update user
	updatedUser, repoErr := u.userRepo.UpdateUser(ctx, Id, entity.User{
		Name:           req.Name,
		ShortName:      req.ShortName,
		Email:          req.Email,
		Age:            req.Age,
		Password:       req.Password,
		OrganizationID: req.OrganizationID,
	})
	if repoErr != nil {
		sysError = repoErr
		return
	}

	res = dto.GetUserResponse{
		ID:             updatedUser.ID,
		Name:           updatedUser.Name,
		ShortName:      updatedUser.ShortName,
		Email:          updatedUser.Email,
		Age:            updatedUser.Age,
		OrganizationID: updatedUser.OrganizationID,
	}
	// Fetch organization if exists
	organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, res.ID)
	if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
		sysError = orgErr
		return
	}
	res.Organization = &organization

	return
}

func (u *UserUsecase) UpdateProfileUser(ctx *fiber.Ctx, Id string, req dto.CreateUserRequest) (res dto.GetUserResponse, sysError syserror.SysError) {
	tx := u.mainDB.TxCreate(ctx)
	defer func() {
		u.mainDB.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	// Update user profile
	updatedUser, repoErr := u.userRepo.UpdateProfileUser(ctx, Id, entity.User{
		Name:      req.Name,
		ShortName: req.ShortName,
		Email:     req.Email,
		Age:       req.Age,
	})
	if repoErr != nil {
		sysError = repoErr
		return
	}

	res = dto.GetUserResponse{
		ID:             updatedUser.ID,
		Name:           updatedUser.Name,
		ShortName:      updatedUser.ShortName,
		Email:          updatedUser.Email,
		Age:            updatedUser.Age,
		OrganizationID: updatedUser.OrganizationID,
	}
	// Fetch organization if exists
	organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, res.ID)
	if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
		sysError = orgErr
		return
	}
	res.Organization = &organization

	return
}

func (u *UserUsecase) GetUsersByOrganizationID(ctx *fiber.Ctx, organizationID string) (res []dto.GetUserResponse, totalRecords int64, sysError syserror.SysError) {
	users, total, repoErr := u.userRepo.GetUsersByOrganizationID(ctx, organizationID)
	if repoErr != nil {
		sysError = repoErr
		return
	}

	for _, user := range users {
		userDto := dto.GetUserResponse{
			ID:             user.ID,
			Name:           user.Name,
			ShortName:      user.ShortName,
			Email:          user.Email,
			Age:            user.Age,
			OrganizationID: user.OrganizationID,
		}

		// Fetch organization if exists
		organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, user.ID)
		if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
			sysError = orgErr
			return
		}
		userDto.Organization = &organization

		res = append(res, userDto)
	}

	totalRecords = total
	return
}

func (u *UserUsecase) GetUserByEmail(ctx *fiber.Ctx, email string) (res dto.GetUserResponse, sysError syserror.SysError) {
	// fetch redis first
	cacheKey := fmt.Sprintf("getUserByEmail:%s", email)

	val, err := u.redisDb.Client.Get(u.redisDb.Ctx, cacheKey).Result()
	if err == nil && val != "" {
		// Unmarshal JSON -> struct
		var cached dto.GetUserResponse
		if unmarshalErr := json.Unmarshal([]byte(val), &cached); unmarshalErr == nil {
			return cached, nil
		}
	}

	// Fetch user by email
	user, repoErr := u.userRepo.GetUserByEmail(ctx, email)
	if repoErr != nil {
		sysError = repoErr
		return
	}

	res = dto.GetUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		ShortName:      user.ShortName,
		Email:          user.Email,
		Age:            user.Age,
		OrganizationID: user.OrganizationID,
	}

	// Fetch organization if exists
	organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, res.ID)
	if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
		sysError = orgErr
		return
	}
	res.Organization = &organization

	// Cache the result
	jsonData, err := json.Marshal(res)
	if err == nil {
		u.redisDb.Client.Set(u.redisDb.Ctx, cacheKey, jsonData, 60*time.Minute)
	}

	return
}

func (u *UserUsecase) GetUserByEmailAndId(ctx *fiber.Ctx, email string, Id string) (res dto.GetUserResponse, sysError syserror.SysError) {
	user, sysError := u.userRepo.GetUserByEmailAndId(ctx, email, Id)
	if sysError != nil {
		return
	}

	res = dto.GetUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		ShortName:      user.ShortName,
		Email:          user.Email,
		Age:            user.Age,
		OrganizationID: user.OrganizationID,
	}

	// Fetch organization if exists
	organization, orgErr := u.organizationUse.GetOrganizationByID(ctx, res.ID)
	if orgErr != nil && orgErr.GetStatusCode() != fiber.StatusNotFound {
		sysError = orgErr
		return
	}
	res.Organization = &organization

	return
}
func (u *UserUsecase) GetPasswordById(ctx *fiber.Ctx, Id int) (password string, sysError syserror.SysError) {
	// fetch redis first
	cacheKey := fmt.Sprintf("getPasswordById:%d", Id)

	val, err := u.redisDb.Client.Get(u.redisDb.Ctx, cacheKey).Result()
	if err == nil && val != "" {
		return val, nil
	}
	// Fetch password by ID
	password, sysError = u.userRepo.GetPasswordById(ctx, Id)

	// cache the result
	u.redisDb.Client.Set(u.redisDb.Ctx, cacheKey, password, 10*time.Minute)
	return

}
