package requests

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/dto/security_dto"
)

type UnixTime struct {
	time.Time
}

func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var ts int64
	if err := json.Unmarshal(b, &ts); err != nil {
		return err
	}

	u.Time = time.Unix(ts, 0)

	return nil
}

type TgLogin struct {
	TgID        int64    `json:"id"`
	Name        string   `json:"first_name"`
	Surname     string   `json:"last_name"`
	Username    string   `json:"username"`
	Avatar      string   `json:"photo_url"`
	AuthDate    UnixTime `json:"auth_date"`
	Hash        string   `json:"hash"`
	CheckString string
}

// ComputeCheckString — вычисляет строку по правилам: key=value, сортировка, join "\n".
func (req *TgLogin) ComputeCheckString() string {
	m := map[string]string{
		"id":         strconv.FormatInt(req.TgID, 10),
		"first_name": req.Name,
		"auth_date":  strconv.FormatInt(req.AuthDate.Unix(), 10),
	}

	if req.Surname != "" {
		m["last_name"] = req.Surname
	}

	if req.Username != "" {
		m["username"] = req.Username
	}

	if req.Avatar != "" {
		m["photo_url"] = req.Avatar
	}

	var parts []string
	for k, v := range m {
		parts = append(parts, k+"="+v)
	}

	sort.Strings(parts)

	return strings.Join(parts, "\n")
}

// UnmarshalJSON fill standard fields, then set CheckString via ComputeCheckString.
func (req *TgLogin) UnmarshalJSON(data []byte) error {
	type Alias TgLogin

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(req),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	req.CheckString = req.ComputeCheckString()

	return nil
}

func (req *TgLogin) ToModel() *security_dto.TgLoginDto {
	return &security_dto.TgLoginDto{
		TgID:        req.TgID,
		Name:        req.Name,
		Surname:     req.Surname,
		Username:    req.Username,
		Avatar:      req.Avatar,
		AuthDate:    req.AuthDate.Time,
		Hash:        req.Hash,
		CheckString: req.CheckString,
	}
}
