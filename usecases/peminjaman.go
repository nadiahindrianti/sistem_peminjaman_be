package usecases

import (
	"errors"
	"sistem_peminjaman_be/dtos"
	"sistem_peminjaman_be/helpers"
	"sistem_peminjaman_be/models"
	"sistem_peminjaman_be/repositories"
	"sort"
	"strings"
	"time"
)

type PeminjamanUsecase interface {
	GetPeminjamans(page, limit int, userID uint, search, nameHotel, orderDateHotel, sort, status string) ([]dtos.PeminjamanResponse, int, error)
	GetPeminjamansByAdmin(page, limit, ratingClass int, search, dateStart, dateEnd, orderBy, status string) ([]dtos.PeminjamanResponse, int, error)
	GetPeminjamansDetailByAdmin(peminjamanId uint) (dtos.PeminjamanResponse, error)
	GetPeminjamanByID(userID, hotelOrderId uint, isCheckIn, isCheckOut bool) (dtos.PeminjamanResponse, error)
	CreatePeminjaman(userID uint, peminjamanInput dtos.PeminjamanInput) (dtos.PeminjamanResponse, error)
	UpdatePeminjaman(userID, peminjamanID uint, status string) (dtos.PeminjamanResponse, error)
}

type peminjamanUsecase struct {
	peminjamanRepo            repositories.PeminjamanRepository
	labRepo                   repositories.LabRepository
	labImageRepo              repositories.LabImageRepository
	suratRekomendasiImageRepo repositories.SuratRekomendasiImageRepository
	userRepo                  repositories.UserRepository
	notificationRepo          repositories.NotificationRepository
}

func NewPeminjamanUsecase(peminjamanRepo repositories.PeminjamanRepository, labRepo repositories.LabRepository, labImageRepo repositories.LabImageRepository, suratRekomendasiImageRepo repositories.SuratRekomendasiImageRepository, userRepo repositories.UserRepository, notificationRepo repositories.NotificationRepository) PeminjamanUsecase {
	return &peminjamanUsecase{peminjamanRepo, labRepo, labImageRepo, suratRekomendasiImageRepo, userRepo, notificationRepo}
}

// GetHotelOrders godoc
// @Summary      Get Hotel Order User
// @Description  Get Hotel Order User
// @Tags         User - Order
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search order"
// @Param name query string false "Filter by name hotel"
// @Param address query string false "Filter by address hotel"
// @Param order_date query string false "Filter by order date hotel"
// @Param order_by query string false "Filter order by" Enums(latest, oldest, highest_price, lowest_price)
// @Param status query string false "Filter by status order" Enums(unpaid, paid, done, canceled, refund)
// @Success      200 {object} dtos.GetAllHotelOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/order/hotel [get]
// @Security BearerAuth
func (u *peminjamanUsecase) GetPeminjamans(page, limit int, userID uint, search, nameHotel, addressHotel, orderDateHotel, orderBy, status string) ([]dtos.HotelOrderResponse, int, error) {
	var hotelOrderResponses []dtos.HotelOrderResponse

	hotelOrders, _, err := u.hotelOrderRepo.GetHotelOrders(page, limit, userID, status)
	if err != nil {
		return hotelOrderResponses, 0, err
	}

	for _, hotelOrder := range hotelOrders {
		getHotel, err := u.hotelRepo.GetHotelByID2(hotelOrder.HotelID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}

		// Check if the search query matches the hotel name, address
		if search != "" &&
			!strings.Contains(strings.ToLower(hotelOrder.HotelOrderCode), strings.ToLower(search)) &&
			!strings.Contains(strings.ToLower(getHotel.Name), strings.ToLower(search)) &&
			!strings.Contains(strings.ToLower(getHotel.Address), strings.ToLower(search)) {
			continue // Skip hotel order if it doesn't match the search query
		}

		// Apply filters based on nameHotel, addressHotel, orderDateHotel
		if nameHotel != "" && !strings.Contains(strings.ToLower(getHotel.Name), strings.ToLower(nameHotel)) {
			continue
		}
		if addressHotel != "" && !strings.Contains(strings.ToLower(getHotel.Address), strings.ToLower(addressHotel)) {
			continue
		}
		if orderDateHotel != "" && helpers.FormatDateToYMD(&hotelOrder.DateStart) != orderDateHotel || orderDateHotel != "" && helpers.FormatDateToYMD(&hotelOrder.DateEnd) != orderDateHotel {
			continue
		}

		getHotelImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotelOrder.HotelID)
		if err != nil {
			continue
		}
		var hotelImageResponses []dtos.HotelImageResponse
		for _, hotelImage := range getHotelImage {
			hotelImageResponse := dtos.HotelImageResponse{
				HotelID:  hotelImage.HotelID,
				ImageUrl: hotelImage.ImageUrl,
			}
			hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
		}
		getHotelFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotelOrder.HotelID)
		if err != nil {
			continue
		}
		getHotelPolicies, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotelOrder.HotelID)
		if err != nil {
			continue
		}
		var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
		for _, hotelFacilities := range getHotelFacilities {
			hotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
				HotelID: hotelFacilities.HotelID,
				Name:    hotelFacilities.Name,
			}
			hotelFacilitiesResponses = append(hotelFacilitiesResponses, hotelFacilitiesResponse)
		}
		getHotelRoom, err := u.hotelRoomRepo.GetHotelRoomByID2(hotelOrder.HotelRoomID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}
		getHotelRoomImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(getHotelRoom.ID)
		if err != nil {
			continue
		}
		var hotelRoomImageResponses []dtos.HotelRoomImageResponse
		for _, hotelRoomImage := range getHotelRoomImage {
			hotelRoomImageResponse := dtos.HotelRoomImageResponse{
				HotelID:     hotelRoomImage.HotelID,
				HotelRoomID: hotelRoomImage.ID,
				ImageUrl:    hotelRoomImage.ImageUrl,
			}
			hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
		}
		getHotelRoomFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByHotelRoomID(getHotelRoom.ID)
		if err != nil {
			continue
		}
		var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
		for _, hotelRoomFacilities := range getHotelRoomFacilities {
			hotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
				HotelID:     hotelRoomFacilities.HotelID,
				HotelRoomID: hotelRoomFacilities.ID,
				Name:        hotelRoomFacilities.Name,
			}
			hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, hotelRoomFacilitiesResponse)
		}
		getPayment, err := u.paymentRepo.GetPaymentByID(uint(hotelOrder.PaymentID))
		if err != nil {
			continue
		}
		getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByHotelOrderID(hotelOrder.ID)
		if err != nil {
			continue
		}

		var travelerDetailResponses []dtos.TravelerDetailResponse

		for _, travelerDetail := range getTravelerDetail {
			travelerDetailResponse := dtos.TravelerDetailResponse{
				ID:           int(travelerDetail.ID),
				Title:        travelerDetail.Title,
				FullName:     travelerDetail.FullName,
				IDCardNumber: *travelerDetail.IDCardNumber,
			}
			travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
		}

		hotelOrderResponse := dtos.HotelOrderResponse{
			HotelOrderID:     int(hotelOrder.ID),
			QuantityAdult:    hotelOrder.QuantityAdult,
			QuantityInfant:   hotelOrder.QuantityInfant,
			NumberOfNight:    hotelOrder.NumberOfNight,
			DateStart:        helpers.FormatDateToYMD(&hotelOrder.DateStart),
			DateEnd:          helpers.FormatDateToYMD(&hotelOrder.DateEnd),
			Price:            hotelOrder.Price,
			TotalAmount:      hotelOrder.TotalAmount,
			NameOrder:        hotelOrder.NameOrder,
			EmailOrder:       hotelOrder.EmailOrder,
			PhoneNumberOrder: hotelOrder.PhoneNumberOrder,
			SpecialRequest:   hotelOrder.SpecialRequest,
			HotelOrderCode:   hotelOrder.HotelOrderCode,
			Status:           hotelOrder.Status,
			Hotel: dtos.HotelByIDResponses{
				HotelID:         getHotel.ID,
				Name:            getHotel.Name,
				Class:           getHotel.Class,
				Description:     getHotel.Description,
				PhoneNumber:     getHotel.PhoneNumber,
				Email:           getHotel.Email,
				Address:         getHotel.Address,
				HotelImage:      hotelImageResponses,
				HotelFacilities: hotelFacilitiesResponses,
				HotelPolicy: dtos.HotelPoliciesResponse{
					HotelID:            getHotelPolicies.HotelID,
					IsCheckInCheckOut:  getHotelPolicies.IsCheckInCheckOut,
					TimeCheckIn:        getHotelPolicies.TimeCheckIn,
					TimeCheckOut:       getHotelPolicies.TimeCheckOut,
					IsPolicyCanceled:   getHotelPolicies.IsPolicyCanceled,
					PolicyMinimumAge:   getHotelPolicies.PolicyMinimumAge,
					IsPolicyMinimumAge: getHotelPolicies.IsPolicyMinimumAge,
					IsCheckInEarly:     getHotelPolicies.IsCheckInEarly,
					IsCheckOutOverdue:  getHotelPolicies.IsCheckOutOverdue,
					IsBreakfast:        getHotelPolicies.IsBreakfast,
					TimeBreakfastStart: getHotelPolicies.TimeBreakfastStart,
					TimeBreakfastEnd:   getHotelPolicies.TimeBreakfastEnd,
					IsSmoking:          getHotelPolicies.IsSmoking,
					IsPet:              getHotelPolicies.IsPet,
				},
				HotelRoom: dtos.HotelRoomHotelIDResponse{
					HotelRoomID:       getHotelRoom.ID,
					HotelID:           getHotelRoom.HotelID,
					Name:              getHotelRoom.Name,
					SizeOfRoom:        getHotelRoom.SizeOfRoom,
					QuantityOfRoom:    getHotelRoom.QuantityOfRoom,
					Description:       getHotelRoom.Description,
					NormalPrice:       getHotelRoom.NormalPrice,
					Discount:          getHotelRoom.Discount,
					DiscountPrice:     getHotelRoom.DiscountPrice,
					NumberOfGuest:     getHotelRoom.NumberOfGuest,
					MattressSize:      getHotelRoom.MattressSize,
					NumberOfMattress:  getHotelRoom.NumberOfMattress,
					HotelRoomImage:    hotelRoomImageResponses,
					HotelRoomFacility: hotelRoomFacilitiesResponses,
				},
			},
			Payment: &dtos.PaymentResponses{
				ID:            int(getPayment.ID),
				Type:          getPayment.Type,
				ImageUrl:      getPayment.ImageUrl,
				Name:          getPayment.Name,
				AccountName:   getPayment.AccountName,
				AccountNumber: getPayment.AccountNumber,
			},
			TravelerDetail: travelerDetailResponses,
			CreatedAt:      hotelOrder.CreatedAt,
			UpdatedAt:      hotelOrder.UpdatedAt,
		}
		hotelOrderResponses = append(hotelOrderResponses, hotelOrderResponse)
	}

	// Sort hotelOrderResponses based on the orderBy parameter
	switch orderBy {
	case "latest":
		// Sort hotelOrderResponses by descending order of CreatedAt
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].CreatedAt.After(hotelOrderResponses[j].CreatedAt)
		})
	case "oldest":
		// Sort hotelOrderResponses by ascending order of CreatedAt
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].CreatedAt.Before(hotelOrderResponses[j].CreatedAt)
		})
	case "highest_price":
		// Sort hotelOrderResponses by descending order of Price
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].TotalAmount > hotelOrderResponses[j].TotalAmount
		})
	case "lowest_price":
		// Sort hotelOrderResponses by ascending order of TotalAmount
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].TotalAmount < hotelOrderResponses[j].TotalAmount
		})
	}
	return hotelOrderResponses, len(hotelOrderResponses), nil
}

// GetHotelOrdersByAdmin godoc
// @Summary      Get Hotel Order User
// @Description  Get Hotel Order User
// @Tags         Admin - Order
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "search hotel name"
// @Param rating_class query int false "Hotel rating class" Enums(1,2,3,4,5)
// @Param date_start query string false "Date start"
// @Param date_end query string false "Date end"
// @Param order_by query string false "Filter order by" Enums(latest, oldest, highest_price, lowest_price)
// @Param status query string false "Filter by status order" Enums(unpaid, paid, done, canceled, refund)
// @Success      200 {object} dtos.GetAllHotelOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/order/hotel [get]
// @Security BearerAuth
func (u *hotelOrderUsecase) GetHotelOrdersByAdmin(page, limit, ratingClass int, search, dateStart, dateEnd, orderBy, status string) ([]dtos.HotelOrderResponse, int, error) {
	var hotelOrderResponses []dtos.HotelOrderResponse

	hotelOrders, _, err := u.hotelOrderRepo.GetHotelOrders(page, limit, 1, status)
	if err != nil {
		return hotelOrderResponses, 0, err
	}

	// Parse dateStart and dateEnd strings into time.Time objects
	var startDate, endDate time.Time
	if dateStart != "" {
		startDate, err = time.Parse("2006-01-02", dateStart)
		if err != nil {
			return hotelOrderResponses, 0, errors.New("invalid dateStart format")
		}
		startDate = startDate.AddDate(0, 0, -1) // Subtract 1 day from startDate
	}
	if dateEnd != "" {
		endDate, err = time.Parse("2006-01-02", dateEnd)
		if err != nil {
			return hotelOrderResponses, 0, errors.New("invalid dateEnd format")
		}
	}

	for _, hotelOrder := range hotelOrders {
		// Filter hotel orders based on dateStart and dateEnd
		if !startDate.IsZero() && hotelOrder.DateStart.Before(startDate) {
			continue // Skip hotel order if its dateStart is before the specified startDate
		}
		if !endDate.IsZero() && hotelOrder.DateStart.After(endDate) {
			continue // Skip hotel order if its dateEnd is after the specified endDate
		}

		getHotel, err := u.hotelRepo.GetHotelByID2(hotelOrder.HotelID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}

		// Check if the search query matches the hotel name, address, or traveler detail name
		if search != "" &&
			!strings.Contains(strings.ToLower(getHotel.Name), strings.ToLower(search)) &&
			!strings.Contains(strings.ToLower(getHotel.Address), strings.ToLower(search)) &&
			!hasMatchingTravelerDetail(hotelOrder.ID, search, u.travelerDetailRepo) {
			continue // Skip hotel order if it doesn't match the search query
		}

		getUser, err := u.userRepo.UserGetById2(hotelOrder.UserID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}
		getHotelImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotelOrder.HotelID)
		if err != nil {
			continue
		}
		var hotelImageResponses []dtos.HotelImageResponse
		for _, hotelImage := range getHotelImage {
			hotelImageResponse := dtos.HotelImageResponse{
				HotelID:  hotelImage.HotelID,
				ImageUrl: hotelImage.ImageUrl,
			}
			hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
		}
		getHotelFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotelOrder.HotelID)
		if err != nil {
			continue
		}
		getHotelPolicies, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotelOrder.HotelID)
		if err != nil {
			continue
		}
		var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
		for _, hotelFacilities := range getHotelFacilities {
			hotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
				HotelID: hotelFacilities.HotelID,
				Name:    hotelFacilities.Name,
			}
			hotelFacilitiesResponses = append(hotelFacilitiesResponses, hotelFacilitiesResponse)
		}
		getHotelRoom, err := u.hotelRoomRepo.GetHotelRoomByID2(hotelOrder.HotelRoomID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}
		getHotelRoomImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(getHotelRoom.ID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}
		var hotelRoomImageResponses []dtos.HotelRoomImageResponse
		for _, hotelRoomImage := range getHotelRoomImage {
			hotelRoomImageResponse := dtos.HotelRoomImageResponse{
				HotelID:     hotelRoomImage.HotelID,
				HotelRoomID: hotelRoomImage.ID,
				ImageUrl:    hotelRoomImage.ImageUrl,
			}
			hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
		}
		getHotelRoomFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByHotelRoomID(getHotelRoom.ID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}
		var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
		for _, hotelRoomFacilities := range getHotelRoomFacilities {
			hotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
				HotelID:     hotelRoomFacilities.HotelID,
				HotelRoomID: hotelRoomFacilities.ID,
				Name:        hotelRoomFacilities.Name,
			}
			hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, hotelRoomFacilitiesResponse)
		}
		getPayment, err := u.paymentRepo.GetPaymentByID(uint(hotelOrder.PaymentID))
		if err != nil {
			return hotelOrderResponses, 0, err
		}
		getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByHotelOrderID(hotelOrder.ID)
		if err != nil {
			return hotelOrderResponses, 0, err
		}

		var travelerDetailResponses []dtos.TravelerDetailResponse

		for _, travelerDetail := range getTravelerDetail {
			travelerDetailResponse := dtos.TravelerDetailResponse{
				ID:           int(travelerDetail.ID),
				Title:        travelerDetail.Title,
				FullName:     travelerDetail.FullName,
				IDCardNumber: *travelerDetail.IDCardNumber,
			}
			travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
		}

		hotelOrderResponse := dtos.HotelOrderResponse{
			HotelOrderID:     int(hotelOrder.ID),
			QuantityAdult:    hotelOrder.QuantityAdult,
			QuantityInfant:   hotelOrder.QuantityInfant,
			NumberOfNight:    hotelOrder.NumberOfNight,
			DateStart:        helpers.FormatDateToYMD(&hotelOrder.DateStart),
			DateEnd:          helpers.FormatDateToYMD(&hotelOrder.DateEnd),
			Price:            hotelOrder.Price,
			TotalAmount:      hotelOrder.TotalAmount,
			NameOrder:        hotelOrder.NameOrder,
			EmailOrder:       hotelOrder.EmailOrder,
			PhoneNumberOrder: hotelOrder.PhoneNumberOrder,
			SpecialRequest:   hotelOrder.SpecialRequest,
			HotelOrderCode:   hotelOrder.HotelOrderCode,
			Status:           hotelOrder.Status,
			Hotel: dtos.HotelByIDResponses{
				HotelID:         getHotel.ID,
				Name:            getHotel.Name,
				Class:           getHotel.Class,
				Description:     getHotel.Description,
				PhoneNumber:     getHotel.PhoneNumber,
				Email:           getHotel.Email,
				Address:         getHotel.Address,
				HotelImage:      hotelImageResponses,
				HotelFacilities: hotelFacilitiesResponses,
				HotelPolicy: dtos.HotelPoliciesResponse{
					HotelID:            getHotelPolicies.HotelID,
					IsCheckInCheckOut:  getHotelPolicies.IsCheckInCheckOut,
					TimeCheckIn:        getHotelPolicies.TimeCheckIn,
					TimeCheckOut:       getHotelPolicies.TimeCheckOut,
					IsPolicyCanceled:   getHotelPolicies.IsPolicyCanceled,
					PolicyMinimumAge:   getHotelPolicies.PolicyMinimumAge,
					IsPolicyMinimumAge: getHotelPolicies.IsPolicyMinimumAge,
					IsCheckInEarly:     getHotelPolicies.IsCheckInEarly,
					IsCheckOutOverdue:  getHotelPolicies.IsCheckOutOverdue,
					IsBreakfast:        getHotelPolicies.IsBreakfast,
					TimeBreakfastStart: getHotelPolicies.TimeBreakfastStart,
					TimeBreakfastEnd:   getHotelPolicies.TimeBreakfastEnd,
					IsSmoking:          getHotelPolicies.IsSmoking,
					IsPet:              getHotelPolicies.IsPet,
				},
				HotelRoom: dtos.HotelRoomHotelIDResponse{
					HotelRoomID:       getHotelRoom.ID,
					HotelID:           getHotelRoom.HotelID,
					Name:              getHotelRoom.Name,
					SizeOfRoom:        getHotelRoom.SizeOfRoom,
					QuantityOfRoom:    getHotelRoom.QuantityOfRoom,
					Description:       getHotelRoom.Description,
					NormalPrice:       getHotelRoom.NormalPrice,
					Discount:          getHotelRoom.Discount,
					DiscountPrice:     getHotelRoom.DiscountPrice,
					NumberOfGuest:     getHotelRoom.NumberOfGuest,
					MattressSize:      getHotelRoom.MattressSize,
					NumberOfMattress:  getHotelRoom.NumberOfMattress,
					HotelRoomImage:    hotelRoomImageResponses,
					HotelRoomFacility: hotelRoomFacilitiesResponses,
				},
			},
			Payment: &dtos.PaymentResponses{
				ID:            int(getPayment.ID),
				Type:          getPayment.Type,
				ImageUrl:      getPayment.ImageUrl,
				Name:          getPayment.Name,
				AccountName:   getPayment.AccountName,
				AccountNumber: getPayment.AccountNumber,
			},
			TravelerDetail: travelerDetailResponses,
			User: &dtos.UserInformationResponses{
				ID:             getUser.ID,
				FullName:       getUser.FullName,
				Email:          getUser.Email,
				PhoneNumber:    getUser.PhoneNumber,
				BirthDate:      helpers.FormatDateToYMD(getUser.BirthDate),
				ProfilePicture: getUser.ProfilePicture,
				Citizen:        getUser.Citizen,
			},
			CreatedAt: hotelOrder.CreatedAt,
			UpdatedAt: hotelOrder.UpdatedAt,
		}

		if ratingClass != 0 && getHotel.Class != ratingClass {
			continue // Skip the hotel if its rating class is below the specified ratingClass
		}
		hotelOrderResponses = append(hotelOrderResponses, hotelOrderResponse)
	}

	// Sort hotelOrderResponses based on the orderBy parameter
	switch orderBy {
	case "latest":
		// Sort hotelOrderResponses by descending order of CreatedAt
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].CreatedAt.After(hotelOrderResponses[j].CreatedAt)
		})
	case "oldest":
		// Sort hotelOrderResponses by ascending order of CreatedAt
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].CreatedAt.Before(hotelOrderResponses[j].CreatedAt)
		})
	case "highest_price":
		// Sort hotelOrderResponses by descending order of Price
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].Price > hotelOrderResponses[j].Price
		})
	case "lowest_price":
		// Sort hotelOrderResponses by ascending order of Price
		sort.SliceStable(hotelOrderResponses, func(i, j int) bool {
			return hotelOrderResponses[i].Price < hotelOrderResponses[j].Price
		})
	}
	return hotelOrderResponses, len(hotelOrderResponses), nil
}

// GetHotelOrdersByAdmin godoc
// @Summary      Get Hotel Order User
// @Description  Get Hotel Order User
// @Tags         Admin - Order
// @Accept       json
// @Produce      json
// @Param hotel_order_id query int true "Hotel Order ID"
// @Success      200 {object} dtos.HotelOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/order/hotel/detail [get]
// @Security BearerAuth
func (u *hotelOrderUsecase) GetHotelOrdersDetailByAdmin(hotelOrderId uint) (dtos.HotelOrderResponse, error) {
	var hotelOrderResponses dtos.HotelOrderResponse

	hotelOrder, err := u.hotelOrderRepo.GetHotelOrderByID(hotelOrderId, 1)
	if err != nil {
		return hotelOrderResponses, err
	}

	getHotel, err := u.hotelRepo.GetHotelByID2(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelImageResponses []dtos.HotelImageResponse
	for _, hotelImage := range getHotelImage {
		hotelImageResponse := dtos.HotelImageResponse{
			HotelID:  hotelImage.HotelID,
			ImageUrl: hotelImage.ImageUrl,
		}
		hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
	}
	getHotelFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelPolicies, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
	for _, hotelFacilities := range getHotelFacilities {
		hotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
			HotelID: hotelFacilities.HotelID,
			Name:    hotelFacilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, hotelFacilitiesResponse)
	}
	getHotelRoom, err := u.hotelRoomRepo.GetHotelRoomByID2(hotelOrder.HotelRoomID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelRoomImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(getHotelRoom.ID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelRoomImageResponses []dtos.HotelRoomImageResponse
	for _, hotelRoomImage := range getHotelRoomImage {
		hotelRoomImageResponse := dtos.HotelRoomImageResponse{
			HotelID:     hotelRoomImage.HotelID,
			HotelRoomID: hotelRoomImage.ID,
			ImageUrl:    hotelRoomImage.ImageUrl,
		}
		hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
	}
	getHotelRoomFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByHotelRoomID(getHotelRoom.ID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
	for _, hotelRoomFacilities := range getHotelRoomFacilities {
		hotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
			HotelID:     hotelRoomFacilities.HotelID,
			HotelRoomID: hotelRoomFacilities.ID,
			Name:        hotelRoomFacilities.Name,
		}
		hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, hotelRoomFacilitiesResponse)
	}
	getPayment, err := u.paymentRepo.GetPaymentByID(uint(hotelOrder.PaymentID))
	if err != nil {
		return hotelOrderResponses, err
	}
	getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByHotelOrderID(hotelOrder.ID)
	if err != nil {
		return hotelOrderResponses, err
	}

	var travelerDetailResponses []dtos.TravelerDetailResponse

	for _, travelerDetail := range getTravelerDetail {
		travelerDetailResponse := dtos.TravelerDetailResponse{
			ID:           int(travelerDetail.ID),
			Title:        travelerDetail.Title,
			FullName:     travelerDetail.FullName,
			IDCardNumber: *travelerDetail.IDCardNumber,
		}
		travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
	}

	hotelOrderResponses = dtos.HotelOrderResponse{
		HotelOrderID:     int(hotelOrder.ID),
		QuantityAdult:    hotelOrder.QuantityAdult,
		QuantityInfant:   hotelOrder.QuantityInfant,
		NumberOfNight:    hotelOrder.NumberOfNight,
		DateStart:        helpers.FormatDateToYMD(&hotelOrder.DateStart),
		DateEnd:          helpers.FormatDateToYMD(&hotelOrder.DateEnd),
		Price:            hotelOrder.Price,
		TotalAmount:      hotelOrder.TotalAmount,
		NameOrder:        hotelOrder.NameOrder,
		EmailOrder:       hotelOrder.EmailOrder,
		PhoneNumberOrder: hotelOrder.PhoneNumberOrder,
		SpecialRequest:   hotelOrder.SpecialRequest,
		HotelOrderCode:   hotelOrder.HotelOrderCode,
		Status:           hotelOrder.Status,
		Hotel: dtos.HotelByIDResponses{
			HotelID:         getHotel.ID,
			Name:            getHotel.Name,
			Class:           getHotel.Class,
			Description:     getHotel.Description,
			PhoneNumber:     getHotel.PhoneNumber,
			Email:           getHotel.Email,
			Address:         getHotel.Address,
			HotelImage:      hotelImageResponses,
			HotelFacilities: hotelFacilitiesResponses,
			HotelPolicy: dtos.HotelPoliciesResponse{
				HotelID:            getHotelPolicies.HotelID,
				IsCheckInCheckOut:  getHotelPolicies.IsCheckInCheckOut,
				TimeCheckIn:        getHotelPolicies.TimeCheckIn,
				TimeCheckOut:       getHotelPolicies.TimeCheckOut,
				IsPolicyCanceled:   getHotelPolicies.IsPolicyCanceled,
				PolicyMinimumAge:   getHotelPolicies.PolicyMinimumAge,
				IsPolicyMinimumAge: getHotelPolicies.IsPolicyMinimumAge,
				IsCheckInEarly:     getHotelPolicies.IsCheckInEarly,
				IsCheckOutOverdue:  getHotelPolicies.IsCheckOutOverdue,
				IsBreakfast:        getHotelPolicies.IsBreakfast,
				TimeBreakfastStart: getHotelPolicies.TimeBreakfastStart,
				TimeBreakfastEnd:   getHotelPolicies.TimeBreakfastEnd,
				IsSmoking:          getHotelPolicies.IsSmoking,
				IsPet:              getHotelPolicies.IsPet,
			},
			HotelRoom: dtos.HotelRoomHotelIDResponse{
				HotelRoomID:       getHotelRoom.ID,
				HotelID:           getHotelRoom.HotelID,
				Name:              getHotelRoom.Name,
				SizeOfRoom:        getHotelRoom.SizeOfRoom,
				QuantityOfRoom:    getHotelRoom.QuantityOfRoom,
				Description:       getHotelRoom.Description,
				NormalPrice:       getHotelRoom.NormalPrice,
				Discount:          getHotelRoom.Discount,
				DiscountPrice:     getHotelRoom.DiscountPrice,
				NumberOfGuest:     getHotelRoom.NumberOfGuest,
				MattressSize:      getHotelRoom.MattressSize,
				NumberOfMattress:  getHotelRoom.NumberOfMattress,
				HotelRoomImage:    hotelRoomImageResponses,
				HotelRoomFacility: hotelRoomFacilitiesResponses,
			},
		},
		Payment: &dtos.PaymentResponses{
			ID:            int(getPayment.ID),
			Type:          getPayment.Type,
			ImageUrl:      getPayment.ImageUrl,
			Name:          getPayment.Name,
			AccountName:   getPayment.AccountName,
			AccountNumber: getPayment.AccountNumber,
		},
		TravelerDetail: travelerDetailResponses,
		CreatedAt:      hotelOrder.CreatedAt,
		UpdatedAt:      hotelOrder.UpdatedAt,
	}
	return hotelOrderResponses, nil
}

// GetHotelOrderByID godoc
// @Summary      Get Hotel Order User by ID
// @Description  Get Hotel Order User by ID
// @Tags         User - Order
// @Accept       json
// @Produce      json
// @Param hotel_order_id query int true "Hotel Order ID"
// @Param update_check_in query bool false "Use this params if want update status order check in"
// @Param update_check_out query bool false "Use this params if want update status order check out"
// @Success      200 {object} dtos.HotelOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/order/hotel/detail [get]
// @Security BearerAuth
func (u *hotelOrderUsecase) GetHotelOrderByID(userID, hotelOrderId uint, isCheckIn, isCheckOut bool) (dtos.HotelOrderResponse, error) {
	var hotelOrderResponses dtos.HotelOrderResponse

	hotelOrder, err := u.hotelOrderRepo.GetHotelOrderByID(hotelOrderId, userID)
	if err != nil {
		return hotelOrderResponses, err
	}

	if hotelOrder.IsCheckIn == false && hotelOrder.IsCheckOut == false && isCheckIn == true {
		hotelOrder.IsCheckIn = true
		_, _ = u.hotelOrderRepo.UpdateHotelOrder(hotelOrder)
	}

	if hotelOrder.IsCheckIn == true && hotelOrder.IsCheckOut == false && isCheckOut == true {
		hotelOrder.IsCheckOut = true
		hotelOrder.Status = "done"
		_, _ = u.hotelOrderRepo.UpdateHotelOrder(hotelOrder)
	}

	getHotel, err := u.hotelRepo.GetHotelByID2(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelImageResponses []dtos.HotelImageResponse
	for _, hotelImage := range getHotelImage {
		hotelImageResponse := dtos.HotelImageResponse{
			HotelID:  hotelImage.HotelID,
			ImageUrl: hotelImage.ImageUrl,
		}
		hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
	}
	getHotelFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelPolicies, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
	for _, hotelFacilities := range getHotelFacilities {
		hotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
			HotelID: hotelFacilities.HotelID,
			Name:    hotelFacilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, hotelFacilitiesResponse)
	}
	getHotelRoom, err := u.hotelRoomRepo.GetHotelRoomByID2(hotelOrder.HotelRoomID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelRoomImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(getHotelRoom.ID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelRoomImageResponses []dtos.HotelRoomImageResponse
	for _, hotelRoomImage := range getHotelRoomImage {
		hotelRoomImageResponse := dtos.HotelRoomImageResponse{
			HotelID:     hotelRoomImage.HotelID,
			HotelRoomID: hotelRoomImage.ID,
			ImageUrl:    hotelRoomImage.ImageUrl,
		}
		hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
	}
	getHotelRoomFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByHotelRoomID(getHotelRoom.ID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
	for _, hotelRoomFacilities := range getHotelRoomFacilities {
		hotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
			HotelID:     hotelRoomFacilities.HotelID,
			HotelRoomID: hotelRoomFacilities.ID,
			Name:        hotelRoomFacilities.Name,
		}
		hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, hotelRoomFacilitiesResponse)
	}
	getPayment, err := u.paymentRepo.GetPaymentByID(uint(hotelOrder.PaymentID))
	if err != nil {
		return hotelOrderResponses, err
	}
	getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByHotelOrderID(hotelOrder.ID)
	if err != nil {
		return hotelOrderResponses, err
	}

	var travelerDetailResponses []dtos.TravelerDetailResponse

	for _, travelerDetail := range getTravelerDetail {
		travelerDetailResponse := dtos.TravelerDetailResponse{
			ID:           int(travelerDetail.ID),
			Title:        travelerDetail.Title,
			FullName:     travelerDetail.FullName,
			IDCardNumber: *travelerDetail.IDCardNumber,
		}
		travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
	}

	hotelOrderResponses = dtos.HotelOrderResponse{
		HotelOrderID:     int(hotelOrder.ID),
		QuantityAdult:    hotelOrder.QuantityAdult,
		QuantityInfant:   hotelOrder.QuantityInfant,
		NumberOfNight:    hotelOrder.NumberOfNight,
		DateStart:        helpers.FormatDateToYMD(&hotelOrder.DateStart),
		DateEnd:          helpers.FormatDateToYMD(&hotelOrder.DateEnd),
		Price:            hotelOrder.Price,
		TotalAmount:      hotelOrder.TotalAmount,
		NameOrder:        hotelOrder.NameOrder,
		EmailOrder:       hotelOrder.EmailOrder,
		PhoneNumberOrder: hotelOrder.PhoneNumberOrder,
		SpecialRequest:   hotelOrder.SpecialRequest,
		HotelOrderCode:   hotelOrder.HotelOrderCode,
		IsCheckIn:        hotelOrder.IsCheckIn,
		IsCheckOut:       hotelOrder.IsCheckOut,
		Status:           hotelOrder.Status,
		Hotel: dtos.HotelByIDResponses{
			HotelID:         getHotel.ID,
			Name:            getHotel.Name,
			Class:           getHotel.Class,
			Description:     getHotel.Description,
			PhoneNumber:     getHotel.PhoneNumber,
			Email:           getHotel.Email,
			Address:         getHotel.Address,
			HotelImage:      hotelImageResponses,
			HotelFacilities: hotelFacilitiesResponses,
			HotelPolicy: dtos.HotelPoliciesResponse{
				HotelID:            getHotelPolicies.HotelID,
				IsCheckInCheckOut:  getHotelPolicies.IsCheckInCheckOut,
				TimeCheckIn:        getHotelPolicies.TimeCheckIn,
				TimeCheckOut:       getHotelPolicies.TimeCheckOut,
				IsPolicyCanceled:   getHotelPolicies.IsPolicyCanceled,
				PolicyMinimumAge:   getHotelPolicies.PolicyMinimumAge,
				IsPolicyMinimumAge: getHotelPolicies.IsPolicyMinimumAge,
				IsCheckInEarly:     getHotelPolicies.IsCheckInEarly,
				IsCheckOutOverdue:  getHotelPolicies.IsCheckOutOverdue,
				IsBreakfast:        getHotelPolicies.IsBreakfast,
				TimeBreakfastStart: getHotelPolicies.TimeBreakfastStart,
				TimeBreakfastEnd:   getHotelPolicies.TimeBreakfastEnd,
				IsSmoking:          getHotelPolicies.IsSmoking,
				IsPet:              getHotelPolicies.IsPet,
			},
			HotelRoom: dtos.HotelRoomHotelIDResponse{
				HotelRoomID:       getHotelRoom.ID,
				HotelID:           getHotelRoom.HotelID,
				Name:              getHotelRoom.Name,
				SizeOfRoom:        getHotelRoom.SizeOfRoom,
				QuantityOfRoom:    getHotelRoom.QuantityOfRoom,
				Description:       getHotelRoom.Description,
				NormalPrice:       getHotelRoom.NormalPrice,
				Discount:          getHotelRoom.Discount,
				DiscountPrice:     getHotelRoom.DiscountPrice,
				NumberOfGuest:     getHotelRoom.NumberOfGuest,
				MattressSize:      getHotelRoom.MattressSize,
				NumberOfMattress:  getHotelRoom.NumberOfMattress,
				HotelRoomImage:    hotelRoomImageResponses,
				HotelRoomFacility: hotelRoomFacilitiesResponses,
			},
		},
		Payment: &dtos.PaymentResponses{
			ID:            int(getPayment.ID),
			Type:          getPayment.Type,
			ImageUrl:      getPayment.ImageUrl,
			Name:          getPayment.Name,
			AccountName:   getPayment.AccountName,
			AccountNumber: getPayment.AccountNumber,
		},
		TravelerDetail: travelerDetailResponses,
		CreatedAt:      hotelOrder.CreatedAt,
		UpdatedAt:      hotelOrder.UpdatedAt,
	}

	return hotelOrderResponses, nil
}

// GetHotelOrderByID godoc
// @Summary      Get Hotel Order User by ID
// @Description  Get Hotel Order User by ID
// @Tags         User - Order
// @Accept       json
// @Produce      json
// @Param hotel_order_id query int true "Hotel Order ID"
// @Param update_check_in query bool false "Use this params if update status order check in"
// @Param update_check_out query bool false "Use this params if update status order check out"
// @Success      200 {object} dtos.HotelOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/order/hotel/detail/midtrans [get]
// @Security BearerAuth
func (u *hotelOrderUsecase) GetHotelOrderByID2(userID, hotelOrderId uint, isCheckIn, isCheckOut bool) (dtos.HotelOrderResponse2, error) {
	var hotelOrderResponses dtos.HotelOrderResponse2

	hotelOrder, err := u.hotelOrderRepo.GetHotelOrderByID2(hotelOrderId, userID)
	if err != nil {
		return hotelOrderResponses, err
	}

	InitiateCoreApiClient()

	res, err := c.CheckTransaction(hotelOrder.HotelOrderCode)
	if res.TransactionStatus == "settlement" {
		hotelOrder.Status = "paid"
		_, _ = u.hotelOrderRepo.UpdateHotelOrder2(hotelOrder)
	}
	if res.TransactionStatus == "expire" {
		hotelOrder.Status = "canceled"
		_, _ = u.hotelOrderRepo.UpdateHotelOrder2(hotelOrder)
	}
	if res.TransactionStatus == "" {
		hotelOrder.Status = "canceled"
		_, _ = u.hotelOrderRepo.UpdateHotelOrder2(hotelOrder)
	}

	if hotelOrder.IsCheckIn == false && hotelOrder.IsCheckOut == false && isCheckIn == true {
		hotelOrder.IsCheckIn = true
		_, _ = u.hotelOrderRepo.UpdateHotelOrder2(hotelOrder)
	}

	if hotelOrder.IsCheckIn == true && hotelOrder.IsCheckOut == false && isCheckOut == true {
		hotelOrder.IsCheckOut = true
		hotelOrder.Status = "done"
		_, _ = u.hotelOrderRepo.UpdateHotelOrder2(hotelOrder)
	}

	getHotel, err := u.hotelRepo.GetHotelByID2(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelImageResponses []dtos.HotelImageResponse
	for _, hotelImage := range getHotelImage {
		hotelImageResponse := dtos.HotelImageResponse{
			HotelID:  hotelImage.HotelID,
			ImageUrl: hotelImage.ImageUrl,
		}
		hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
	}
	getHotelFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelPolicies, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotelOrder.HotelID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
	for _, hotelFacilities := range getHotelFacilities {
		hotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
			HotelID: hotelFacilities.HotelID,
			Name:    hotelFacilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, hotelFacilitiesResponse)
	}
	getHotelRoom, err := u.hotelRoomRepo.GetHotelRoomByID(hotelOrder.HotelRoomID)
	if err != nil {
		return hotelOrderResponses, err
	}
	getHotelRoomImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(getHotelRoom.ID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelRoomImageResponses []dtos.HotelRoomImageResponse
	for _, hotelRoomImage := range getHotelRoomImage {
		hotelRoomImageResponse := dtos.HotelRoomImageResponse{
			HotelID:     hotelRoomImage.HotelID,
			HotelRoomID: hotelRoomImage.ID,
			ImageUrl:    hotelRoomImage.ImageUrl,
		}
		hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
	}
	getHotelRoomFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByHotelRoomID(getHotelRoom.ID)
	if err != nil {
		return hotelOrderResponses, err
	}
	var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
	for _, hotelRoomFacilities := range getHotelRoomFacilities {
		hotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
			HotelID:     hotelRoomFacilities.HotelID,
			HotelRoomID: hotelRoomFacilities.ID,
			Name:        hotelRoomFacilities.Name,
		}
		hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, hotelRoomFacilitiesResponse)
	}
	getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByHotelOrderID(hotelOrder.ID)
	if err != nil {
		return hotelOrderResponses, err
	}

	var travelerDetailResponses []dtos.TravelerDetailResponse

	for _, travelerDetail := range getTravelerDetail {
		travelerDetailResponse := dtos.TravelerDetailResponse{
			ID:           int(travelerDetail.ID),
			Title:        travelerDetail.Title,
			FullName:     travelerDetail.FullName,
			IDCardNumber: *travelerDetail.IDCardNumber,
		}
		travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
	}

	hotelOrderResponses = dtos.HotelOrderResponse2{
		PaymentURL:       hotelOrder.PaymentURL,
		HotelOrderID:     int(hotelOrder.ID),
		QuantityAdult:    hotelOrder.QuantityAdult,
		QuantityInfant:   hotelOrder.QuantityInfant,
		NumberOfNight:    hotelOrder.NumberOfNight,
		DateStart:        helpers.FormatDateToYMD(&hotelOrder.DateStart),
		DateEnd:          helpers.FormatDateToYMD(&hotelOrder.DateEnd),
		Price:            hotelOrder.Price,
		TotalAmount:      hotelOrder.TotalAmount,
		NameOrder:        hotelOrder.NameOrder,
		EmailOrder:       hotelOrder.EmailOrder,
		PhoneNumberOrder: hotelOrder.PhoneNumberOrder,
		SpecialRequest:   hotelOrder.SpecialRequest,
		HotelOrderCode:   hotelOrder.HotelOrderCode,
		IsCheckIn:        hotelOrder.IsCheckIn,
		IsCheckOut:       hotelOrder.IsCheckOut,
		Status:           hotelOrder.Status,
		Hotel: dtos.HotelByIDResponses{
			HotelID:         getHotel.ID,
			Name:            getHotel.Name,
			Class:           getHotel.Class,
			Description:     getHotel.Description,
			PhoneNumber:     getHotel.PhoneNumber,
			Email:           getHotel.Email,
			Address:         getHotel.Address,
			HotelImage:      hotelImageResponses,
			HotelFacilities: hotelFacilitiesResponses,
			HotelPolicy: dtos.HotelPoliciesResponse{
				HotelID:            getHotelPolicies.HotelID,
				IsCheckInCheckOut:  getHotelPolicies.IsCheckInCheckOut,
				TimeCheckIn:        getHotelPolicies.TimeCheckIn,
				TimeCheckOut:       getHotelPolicies.TimeCheckOut,
				IsPolicyCanceled:   getHotelPolicies.IsPolicyCanceled,
				PolicyMinimumAge:   getHotelPolicies.PolicyMinimumAge,
				IsPolicyMinimumAge: getHotelPolicies.IsPolicyMinimumAge,
				IsCheckInEarly:     getHotelPolicies.IsCheckInEarly,
				IsCheckOutOverdue:  getHotelPolicies.IsCheckOutOverdue,
				IsBreakfast:        getHotelPolicies.IsBreakfast,
				TimeBreakfastStart: getHotelPolicies.TimeBreakfastStart,
				TimeBreakfastEnd:   getHotelPolicies.TimeBreakfastEnd,
				IsSmoking:          getHotelPolicies.IsSmoking,
				IsPet:              getHotelPolicies.IsPet,
			},
			HotelRoom: dtos.HotelRoomHotelIDResponse{
				HotelRoomID:       getHotelRoom.ID,
				HotelID:           getHotelRoom.HotelID,
				Name:              getHotelRoom.Name,
				SizeOfRoom:        getHotelRoom.SizeOfRoom,
				QuantityOfRoom:    getHotelRoom.QuantityOfRoom,
				Description:       getHotelRoom.Description,
				NormalPrice:       getHotelRoom.NormalPrice,
				Discount:          getHotelRoom.Discount,
				DiscountPrice:     getHotelRoom.DiscountPrice,
				NumberOfGuest:     getHotelRoom.NumberOfGuest,
				MattressSize:      getHotelRoom.MattressSize,
				NumberOfMattress:  getHotelRoom.NumberOfMattress,
				HotelRoomImage:    hotelRoomImageResponses,
				HotelRoomFacility: hotelRoomFacilitiesResponses,
			},
		},
		TravelerDetail: travelerDetailResponses,
		CreatedAt:      hotelOrder.CreatedAt,
		UpdatedAt:      hotelOrder.UpdatedAt,
	}

	return hotelOrderResponses, nil
}

// CreateHotelOrder godoc
// @Summary      Order Hotel
// @Description  Order Hotel
// @Tags         User - Hotel
// @Accept       json
// @Produce      json
// @Param        request body dtos.HotelOrderInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.HotelOrderCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/hotel/order [post]
// @Security BearerAuth
func (u *peminjamanUsecase) CreatePeminjaman(userID uint, peminjamanInput dtos.PeminjamanInput) (dtos.PeminjamanResponse, error) {
	var peminjamanResponse dtos.PeminjamanResponse

	// Validasi input
	if peminjamanInput.WaktuPeminjaman == "" {
		return peminjamanResponse, errors.New("Waktu peminjaman harus diisi")
	}

	tanggalPeminjaman, err := time.Parse("2006-01-02", peminjamanInput.TanggalPeminjaman)
	if err != nil {
		// Handle error saat parsing tanggal
		return peminjamanResponse, err
	}

	// Membuat objek Peminjaman dengan TanggalPeminjaman yang sudah di-parse
	createPeminjaman := models.Peminjaman{
		UserID:            userID,
		TanggalPeminjaman: tanggalPeminjaman,
		WaktuPeminjaman:   peminjamanInput.WaktuPeminjaman,
		Description:       peminjamanInput.Description,
		Status:            "request", // Default status
	}

	// Membuat peminjaman di repository
	createdPeminjaman, err := u.peminjamanRepo.CreatePeminjaman(createPeminjaman)
	if err != nil {
		return peminjamanResponse, err
	}

	// Proses penambahan surat rekomendasi jika ada
	for _, suratRekomendasiImageInput := range peminjamanInput.SuratRekomendasiImage {
		suratRekomendasiImage := models.SuratRekomendasiImage{
			PeminjamanID: createdPeminjaman.ID,
			ImageUrl:     suratRekomendasiImageInput.ImageUrl,
		}
		_, err := u.suratRekomendasiImageRepo.CreateSuratRekomendasiImage(suratRekomendasiImage)
		if err != nil {
			return peminjamanResponse, err
		}
	}

	// Mendapatkan detail lab terkait
	getLab, err := u.labRepo.GetLabByID(createdPeminjaman.LabID)
	if err != nil {
		return peminjamanResponse, err
	}

	// Mendapatkan user terkait
	getUser, err := u.userRepo.GetUserByID(createdPeminjaman.ID)
	if err != nil {
		return peminjamanResponse, err
	}


	// Mengisi response dengan data yang diperlukan
	peminjamanResponse = dtos.PeminjamanResponse{
		PeminjamanID:      int(createdPeminjaman.ID),
		TanggalPeminjaman: helpers.FormatDateToYMD(&createdPeminjaman.TanggalPeminjaman),
		WaktuPeminjaman:   createdPeminjaman.WaktuPeminjaman,
		Description:       createdPeminjaman.Description,
		Status:            createdPeminjaman.Status,
		Lab: dtos.LabByIDResponses{
			LabID:       getLab.ID,
			Name:        getLab.Name,
			Description: getLab.Description,
		},
		User: &dtos.UserInformationResponses{
			ID:             getUser.ID,
			FullName:       getUser.FullName,
			Email:          getUser.Email,
			NIMNIP:         getUser.NIMNIP,
			ProfilePicture: getUser.ProfilePicture,
		},
		CreatedAt: createdPeminjaman.CreatedAt,
		UpdatedAt: createdPeminjaman.UpdatedAt,
	}

	return peminjamanResponse, nil
}

func (u *peminjamanUsecase) UpdatePeminjaman(userID, peminjamanID uint, status string) (dtos.PeminjamanResponse, error) {
	var peminjamanResponses dtos.PeminjamanResponse

	peminjaman, err := u.peminjamanRepo.GetPeminjamanByID(peminjamanID, userID)
	if err != nil {
		return peminjamanResponses, err
	}
	if peminjaman.Status == status || status == "request" {
		return peminjamanResponses, errors.New("Failed to update peminjaman status")
	}
	peminjaman.Status = status
	peminjaman, err = u.peminjamanRepo.UpdatePeminjaman(peminjaman)
	if err != nil {
		return peminjamanResponses, err
	}

	if peminjaman.ID > 0 && peminjaman.Status == "accept" {
		createNotification := models.Notification{
			UserID:       userID,
			TemplateID:   3,
			PeminjamanID: peminjaman.ID,
		}

		_, err = u.notificationRepo.CreateNotification(createNotification)
		if err != nil {
			return peminjamanResponses, err
		}
	}

	if peminjaman.ID > 0 && peminjaman.Status == "reject" {
		createNotification := models.Notification{
			UserID:       userID,
			TemplateID:   8,
			PeminjamanID: peminjaman.ID,
		}

		_, err = u.notificationRepo.CreateNotification(createNotification)
		if err != nil {
			return peminjamanResponses, err
		}
	}

	getLab, err := u.labRepo.GetLabByID2(peminjaman.LabID)
	if err != nil {
		return peminjamanResponses, err
	}

	getLabImage, err := u.labImageRepo.GetAllLabImageByID(peminjaman.LabID)
	if err != nil {
		return peminjamanResponses, err
	}

	var labImageResponses []dtos.LabImageResponse
	for _, labImage := range getLabImage {
		labImageResponse := dtos.LabImageResponse{
			LabID:    labImage.LabID,
			ImageUrl: labImage.ImageUrl,
		}
		labImageResponses = append(labImageResponses, labImageResponse)
	}

	getSuratRekomendasiImage, err := u.suratRekomendasiImageRepo.GetAllSuratRekomendasiImageByID(peminjaman.ID)
	if err != nil {
		return peminjamanResponses, err
	}

	var suratRekomendasiImageResponses []dtos.SuratRekomendasiImageResponse
	for _, suratRekomendasiImage := range getSuratRekomendasiImage {
		suratRekomendasiImageResponse := dtos.SuratRekomendasiImageResponse{
			PeminjamanID: suratRekomendasiImage.PeminjamanID,
			ImageUrl:     suratRekomendasiImage.ImageUrl,
		}
		suratRekomendasiImageResponses = append(suratRekomendasiImageResponses, suratRekomendasiImageResponse)
	}

	peminjamanResponses = dtos.PeminjamanResponse{
		PeminjamanID:          int(peminjaman.ID),
		TanggalPeminjaman:     helpers.FormatDateToYMD(&peminjaman.TanggalPeminjaman),
		WaktuPeminjaman:       peminjaman.WaktuPeminjaman,
		SuratRekomendasiImage: suratRekomendasiImageResponses,
		Status:                peminjaman.Status,
		Lab: dtos.LabByIDResponses{
			LabID:       getLab.ID,
			Name:        getLab.Name,
			Description: getLab.Description,
		},
		CreatedAt: peminjaman.CreatedAt,
		UpdatedAt: peminjaman.UpdatedAt,
	}

	return peminjamanResponses, nil
}
