package usecases

import (
	"errors"
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/middlewares"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"sort"
	"strings"
)

type UserUsecase interface {
	UserLogin(input dtos.UserLoginInput) (dtos.UserInformationResponse, error)
	UserRegister(input dtos.UserRegisterInput) (dtos.UserInformationResponse, error)
	AdminRegister(input dtos.UserRegisterInput) (dtos.UserInformationResponse, error)
	UserUpdatePassword(userId uint, input dtos.UserUpdatePasswordInput) (dtos.UserInformationResponse, error)
	UserUpdateProfile(userId uint, input dtos.UserUpdateProfileInput) (dtos.UserInformationResponse, error)
	UserCredential(userId uint) (dtos.UserInformationResponse, error)
	UserUpdatePhotoProfile(userId uint, input dtos.UserUpdatePhotoProfileInput) (dtos.UserInformationResponse, error)
	UserDeletePhotoProfile(userId uint) (dtos.UserInformationResponse, error)
	UserGetAll(page, limit int, search, sortBy, filter string) ([]dtos.UserInformationResponse, int, error)
	UserGetDetail(id int, isDeleted bool) (dtos.UserInformationResponse, error)
	UserAdminRegister(input dtos.UserRegisterInput) (dtos.UserInformationResponse, error)
	UserAdminUpdate(id uint, input dtos.UserRegisterInputUpdateByAdmin) (dtos.UserInformationResponse, error)
}

type userUsecase struct {
	userRepo         repositories.UserRepository
	notificationRepo repositories.NotificationRepository
}

func NewUserUsecase(userRepo repositories.UserRepository, notificationRepo repositories.NotificationRepository) UserUsecase {
	return &userUsecase{userRepo, notificationRepo}
}

// UserLogin godoc
// @Summary      Login
// @Description  Login an account
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserLoginInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /login [post]
func (u *userUsecase) UserLogin(input dtos.UserLoginInput) (dtos.UserInformationResponse, error) {
	var (
		userResponse dtos.UserInformationResponse
		accessToken  string
	)

	input.Email = strings.ToLower(input.Email)

	user, err := u.userRepo.UserGetByEmail(input.Email)
	if err != nil {
		return userResponse, errors.New("Email or password is wrong")
	}

	valid := helpers.ComparePassword(input.Password, user.Password)
	if !valid {
		return userResponse, errors.New("Email or password is wrong")
	}

	accessToken, err = middlewares.CreateToken(user.ID, user.Role)
	if err != nil {
		return userResponse, err
	}


	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.Token = &accessToken
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, nil
}

// UserRegister godoc
// @Summary      Register
// @Description  Register an account
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserRegisterInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.UserCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /register [post]
func (u *userUsecase) UserRegister(input dtos.UserRegisterInput) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	input.Email = strings.ToLower(input.Email)

	user, err := u.userRepo.UserGetByEmail(input.Email)
	if user.ID > 0 {
		return userResponse, errors.New("Email already used")
	}

	if input.Password != input.ConfirmPassword {
		return userResponse, errors.New("Password does not match")
	}

	password, err := helpers.HashPassword(input.Password)
	if err != nil {
		return userResponse, err
	}

	if input.Email == "" || input.FullName == "" || input.Password == "" || input.ConfirmPassword == "" || input.Role == "admin" || input.NIMNIP == "" {
		return userResponse, errors.New("Failed to create user")
	}

	user.FullName = input.FullName
	user.Email = input.Email
	user.Password = password
	user.NIMNIP = input.NIMNIP
	user.ProfilePicture = "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"
	user.Role = "user"

	user, err = u.userRepo.UserCreate(user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

func (u *userUsecase) AdminRegister(input dtos.UserRegisterInput) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	input.Email = strings.ToLower(input.Email)

	user, err := u.userRepo.UserGetByEmail(input.Email)
	if user.ID > 0 {
		return userResponse, errors.New("Email already used")
	}

	if input.Password != input.ConfirmPassword {
		return userResponse, errors.New("Password does not match")
	}

	password, err := helpers.HashPassword(input.Password)
	if err != nil {
		return userResponse, err
	}

	if input.Email == "" || input.FullName == "" || input.Password == "" || input.ConfirmPassword == "" || input.Role == "admin" || input.NIMNIP == "" {
		return userResponse, errors.New("Failed to create user")
	}

	user.FullName = input.FullName
	user.Email = input.Email
	user.Password = password
	user.NIMNIP = input.NIMNIP
	user.ProfilePicture = "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"
	user.Role = "admin"

	user, err = u.userRepo.UserCreate(user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

// UserUpdatePassword godoc
// @Summary      Update Password
// @Description  User update an password
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserUpdatePasswordInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/update-password [put]
// @Security BearerAuth
func (u *userUsecase) UserUpdatePassword(userId uint, input dtos.UserUpdatePasswordInput) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	user, err := u.userRepo.UserGetById(userId)
	if err != nil {
		return userResponse, errors.New("User not found")
	}

	if input.OldPassword == "" || input.NewPassword == "" || input.ConfirmPassword == "" {
		return userResponse, errors.New("Failed to update user password")
	}

	password := helpers.ComparePassword(input.OldPassword, user.Password)
	if !password {
		return userResponse, errors.New("Old password is incorrect")
	}

	if input.OldPassword == input.NewPassword {
		return userResponse, errors.New("New password must be different")
	}

	if input.NewPassword != input.ConfirmPassword {
		return userResponse, errors.New("Confirm password does not match")
	}

	newPassword, err := helpers.HashPassword(input.ConfirmPassword)
	if err != nil {
		return userResponse, err
	}

	user.Password = newPassword

	user, err = u.userRepo.UserUpdate(user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

// UserUpdateProfile godoc
// @Summary      Update Profile
// @Description  User update an profile
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserUpdateProfileInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/update-profile [put]
// @Security BearerAuth
func (u *userUsecase) UserUpdateProfile(userId uint, input dtos.UserUpdateProfileInput) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	user, err := u.userRepo.UserGetById(userId)
	if err != nil {
		return userResponse, errors.New("User not found")
	}

	user.FullName = input.FullName
	user.Email = input.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture

	user, err = u.userRepo.UserUpdate(user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

// UserCredential godoc
// @Summary      Get Credentials
// @Description  User get credentials
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user [get]
// @Security BearerAuth
func (u *userUsecase) UserCredential(userId uint) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	user, err := u.userRepo.UserGetById(userId)
	if err != nil {
		return userResponse, errors.New("User not found")
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

// UserUpdatePhotoProfile godoc
// @Summary      Update Photo Profile
// @Description  User update an photo profile
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        file formData file false "Photo file"
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/update-photo-profile [put]
// @Security BearerAuth
func (u *userUsecase) UserUpdatePhotoProfile(userId uint, input dtos.UserUpdatePhotoProfileInput) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	user, err := u.userRepo.UserGetById(userId)
	if err != nil {
		return userResponse, errors.New("User not found")
	}

	user.ProfilePicture = input.ProfilePicture

	user, err = u.userRepo.UserUpdate(user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

// UserDeletePhotoProfile godoc
// @Summary      Update Information
// @Description  User update an information
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/delete-photo-profile [delete]
// @Security BearerAuth
func (u *userUsecase) UserDeletePhotoProfile(userId uint) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	user, err := u.userRepo.UserGetById(userId)
	if err != nil {
		return userResponse, errors.New("User not found")
	}

	user.ProfilePicture = "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"

	user, err = u.userRepo.UserUpdate(user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	return userResponse, err
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get all users
// @Tags         Admin - User
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search data"
// @Param sort_by query string false "Sort by name" Enums(asc, desc)
// @Param filter query string false "Filter data" Enums(active, inactive)
// @Success      200 {object} dtos.GetAllUserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/user [get]
// @Security BearerAuth
func (u *userUsecase) UserGetAll(page, limit int, search, sortBy, filter string) ([]dtos.UserInformationResponse, int, error) {
	var userResponses []dtos.UserInformationResponse
	users, count, err := u.userRepo.UserGetAll(page, limit, search)
	if err != nil {
		return userResponses, count, err
	}
	filter = strings.ToLower(filter)

	for _, user := range users {
		deletedUser := ""

		if filter == "inactive" && user.DeletedAt.Time.IsZero() {
			continue
		} else if filter == "active" && !user.DeletedAt.Time.IsZero() {
			continue
		}

		if !user.DeletedAt.Time.IsZero() {
			deletedUser = user.DeletedAt.Time.Format("2006-01-02T15:04:05.000-07:00")
		}

		userResponse := dtos.UserInformationResponse{
			ID:             user.ID,
			FullName:       strings.ToUpper(user.FullName),
			Email:          user.Email,
			NIMNIP:         user.NIMNIP,
			ProfilePicture: user.ProfilePicture,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			DeletedAt:      &deletedUser,
		}
		userResponses = append(userResponses, userResponse)
	}

	// Sort trainResponses based on price
	if strings.ToLower(sortBy) == "asc" {
		sort.SliceStable(userResponses, func(i, j int) bool {
			return userResponses[i].FullName < userResponses[j].FullName
		})
	} else if strings.ToLower(sortBy) == "desc" {
		sort.SliceStable(userResponses, func(i, j int) bool {
			return userResponses[i].FullName > userResponses[j].FullName
		})
	}
	return userResponses, count, nil
}

// UserGetDetail godoc
// @Summary      Get detail user
// @Description  Get detail user
// @Tags         Admin - User
// @Accept       json
// @Produce      json
// @Param id query int false "User ID"
// @Param isDeleted query bool false "Use this params if user want to be delete"
// @Success      200 {object} dtos.GetAllUserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/user/detail [get]
// @Security BearerAuth
func (u *userUsecase) UserGetDetail(id int, isDeleted bool) (dtos.UserInformationResponse, error) {
	var userResponses dtos.UserInformationResponse
	user, err := u.userRepo.UserGetById2(uint(id))
	if err != nil {
		return userResponses, err
	}
	user, err = u.userRepo.UserGetDetail(uint(id), isDeleted)
	if err != nil {
		return userResponses, err
	}

	deletedUser := ""
	if !user.DeletedAt.Time.IsZero() {
		deletedUser = user.DeletedAt.Time.Format("2006-01-02T15:04:05.000-07:00")
	}

	userResponses = dtos.UserInformationResponse{
		ID:             user.ID,
		FullName:       user.FullName,
		Email:          user.Email,
		NIMNIP:         user.NIMNIP,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		DeletedAt:      &deletedUser,
	}
	return userResponses, nil
}

// UserAdminRegister godoc
// @Summary      Register user
// @Description  Register an account
// @Tags         Admin - User
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserRegisterInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.UserCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/user/register [post]
// @Security BearerAuth
func (u *userUsecase) UserAdminRegister(input dtos.UserRegisterInput) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	input.Email = strings.ToLower(input.Email)

	user, err := u.userRepo.UserGetByEmail(input.Email)
	if user.ID > 0 {
		return userResponse, errors.New("Email already used")
	}

	if input.Password != input.ConfirmPassword {
		return userResponse, errors.New("Password does not match")
	}

	password, err := helpers.HashPassword(input.Password)
	if err != nil {
		return userResponse, err
	}

	if input.Email == "" || input.FullName == "" || input.Password == "" || input.ConfirmPassword == "" || input.Role == "admin" || input.NIMNIP == "" {
		return userResponse, errors.New("Failed to create user")
	}

	user.FullName = input.FullName
	user.Email = input.Email
	user.Password = password
	user.NIMNIP = input.NIMNIP
	user.ProfilePicture = "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"
	user.Role = "user"

	isActive := false // Default value if the pointer is nil

	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	user, err = u.userRepo.UserCreate2(user, isActive)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	deletedUser := ""
	if !user.DeletedAt.Time.IsZero() {
		deletedUser = user.DeletedAt.Time.Format("2006-01-02T15:04:05.000-07:00")
	}

	userResponse.DeletedAt = &deletedUser

	return userResponse, err
}

// UserAdminRegister godoc
// @Summary      Update user
// @Description  Register an account
// @Tags         Admin - User
// @Accept       json
// @Produce      json
// @Param id path integer true "ID user"
// @Param        request body dtos.UserRegisterInputUpdateByAdmin true "Payload Body [RAW]"
// @Success      201 {object} dtos.UserCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/user/update/{id} [put]
// @Security BearerAuth
func (u *userUsecase) UserAdminUpdate(id uint, input dtos.UserRegisterInputUpdateByAdmin) (dtos.UserInformationResponse, error) {
	var (
		user         models.User
		userResponse dtos.UserInformationResponse
	)

	user, err := u.userRepo.UserGetById2(id)
	if err != nil {
		return userResponse, errors.New("User not found")
	}

	userCheck, err := u.userRepo.UserGetByEmail2(id, input.Email)
	if userCheck.ID != 0 {
		return userResponse, errors.New("Email already used")
	}

	if input.Email == "" || input.FullName == "" || input.Role == "admin" || input.NIMNIP == "" {
		return userResponse, errors.New("Failed to create user")
	}

	user.FullName = input.FullName
	user.Email = input.Email
	user.NIMNIP = input.NIMNIP
	user.ProfilePicture = "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg"
	user.Role = "user"

	isActive := false // Default value if the pointer is nil

	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	user, err = u.userRepo.UserUpdate2(user, isActive)
	if err != nil {
		return userResponse, err
	}

	userResponse.ID = user.ID
	userResponse.FullName = user.FullName
	userResponse.Email = user.Email
	userResponse.NIMNIP = user.NIMNIP
	userResponse.ProfilePicture = user.ProfilePicture
	userResponse.Role = &user.Role
	userResponse.CreatedAt = user.CreatedAt
	userResponse.UpdatedAt = user.UpdatedAt

	deletedUser := ""
	if !user.DeletedAt.Time.IsZero() {
		deletedUser = user.DeletedAt.Time.Format("2006-01-02T15:04:05.000-07:00")
	}

	userResponse.DeletedAt = &deletedUser

	return userResponse, err
}
