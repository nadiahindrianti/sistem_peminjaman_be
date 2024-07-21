package dtos

type DashboardResponse struct {
	CountUser            map[string]interface{}   `json:"count_user"`
	CountLab             map[string]interface{}   `json:"count_lab"`
	CountPeminjaman      map[string]interface{}   `json:"count_peminjaman"`
	CountJadwal          map[string]interface{}   `json:"count_jadwal"`
	NewPeminjaman        []map[string]interface{} `json:"new_peminjaman"`
	NewJadwal            []map[string]interface{} `json:"new_jadwal"`
	NewUser              []map[string]interface{} `json:"new_user"`
	UserTeraktifMeminjam []UserTeraktifMeminjam   `json:"user_teraktif_meminjam"`
}

type UserTeraktifMeminjam struct {
	FullName          string `json:"full_name"`
	JumlahPeminjaman  int    `json:"jumlah_peminjaman"`
}

type DashboardfilterResponse struct {
	CountUser  interface{}              `json:"count_user"`
	CountLab interface{}                `json:"count_lab"`
	CountPeminjaman interface{}         `json:"count_peminjaman"`
	CountJadwal interface{}             `json:"count_jadwal"`
}


type FilterDashboardByMonthRequest struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}