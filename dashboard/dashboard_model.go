package dashboard

// Ref: Saving enumerated values to a database
// https://stackoverflow.com/a/25374979/12976234
type DocType struct {
	String string
}

var (
	Unknown  = DocType{""}
	Dir      = DocType{"dir"}
	Filetype = DocType{"file"}
)

type (
	CashflowOut struct {
		Id           int64  `json:"id"`
		Date         string `json:"date"`
		Note         string `json:"note"`
		Type         string `json:"type"`
		IdrAmout     string `json:"idr_amount"`
		ProveFileUrl string `json:"prove_file_url"`
	}
	CashflowRes struct {
		TotalCash   string `json:"total_cash"`
		IncomeCash  string `json:"income_cash"`
		OutcomeCash string `json:"outcome_cash"`
	}
	DocumentOut struct {
		IsPrivate bool   `json:"is_private"`
		Id        int64  `json:"id"`
		DirId     int64  `json:"dir_id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Url       string `json:"url"`
	}
	MemberOut struct {
		Id                string `json:"id"`
		Username          string `json:"username"`
		Name              string `json:"name"`
		WaPhone           string `json:"wa_phone"`
		OtherPhone        string `json:"other_phone"`
		HomestayName      string `json:"homestay_name"`
		HomestayAddress   string `json:"homestay_address"`
		HomestayLatitude  string `json:"homestay_latitude"`
		HomestayLongitude string `json:"homestay_longitude"`
		ProfilePicUrl     string `json:"profile_pic_url"`
		IsAdmin           bool   `json:"is_admin"`
		IsApproved        bool   `json:"is_approved"`
	}
	DuesOut struct {
		Id        int64  `json:"id"`
		Date      string `json:"date"`
		IdrAmount string `json:"idr_amount"`
	}
	BlogOut struct {
		Id           int64  `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	PositionOut struct {
		Level int16  `json:"level"`
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
	}
	LatestHistoryRes struct {
		Id          int64  `json:"id"`
		Content     string `json:"content"`
		ContentText string `json:"content_text"`
	}
	MembersDuesOut struct {
		Id            int64  `json:"id"`
		MemberId      string `json:"member_id"`
		Status        string `json:"status"`
		Name          string `json:"name"`
		ProfilePicUrl string `json:"profile_pic_url"`
		PayDate       string `json:"pay_date"`
	}
	FindOrgPeriodGoalRes struct {
		Id          int64  `json:"id"`
		Vision      string `json:"vision"`
		VisionText  string `json:"vision_text"`
		Mission     string `json:"mission"`
		MissionText string `json:"mission_text"`
	}
	StructureMemberOut struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		ProfilePicUrl string `json:"profile_pic_url"`
	}
	StructurePositionOut struct {
		Id      uint64               `json:"id"`
		Name    string               `json:"name"`
		Level   int16                `json:"level"`
		Members []StructureMemberOut `json:"members"`
	}
	StructureRes struct {
		Id        uint64                 `json:"id"`
		StartDate string                 `json:"start_date"`
		EndDate   string                 `json:"end_date"`
		Positions []StructurePositionOut `json:"positions"`
		Vision    string                 `json:"vision"`
		Mission   string                 `json:"mission"`
	}
	PeriodRes struct {
		IsActive  bool   `json:"is_active"`
		Id        uint64 `json:"id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	ImageOut struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
	}
)
