package dtos

type DashboardResponse struct {
	CountUser  interface{}              `json:"count_user"`
	CountLab interface{}                `json:"count_lab"`
	CountPeminjaman interface{}         `json:"count_peminjaman"`
	CountJadwal interface{}             `json:"count_jadwal"`
	NewPeminjaman   []map[string]interface{} `json:"new_peminjaman"`
	NewJadwal   []map[string]interface{} `json:"new_jadwal"`
	NewUser    []map[string]interface{} `json:"new_user"`
}