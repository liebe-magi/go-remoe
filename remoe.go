// Package remoe : Nature Remo EおよびNature Remo E Liteを用いて、データを取得するパッケージ
package remoe

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// RemoClient : Nature RemoのAPIにアクセスするクライアントとなる構造体
type RemoClient struct {
	// Token : Nature Remoのアクセストークン
	Token string
}

// remoInfo : Nature Remoから取得したデータを格納する構造体
type remoInfo struct {
	ID     string `json:"id"`
	Device struct {
		Name              string    `json:"name"`
		ID                string    `json:"id"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		MacAddress        string    `json:"mac_address"`
		BtMacAddress      string    `json:"bt_mac_address"`
		SerialNumber      string    `json:"serial_number"`
		FirmwareVersion   string    `json:"firmware_version"`
		TemperatureOffset int       `json:"temperature_offset"`
		HumidityOffset    int       `json:"humidity_offset"`
	} `json:"device,omitempty"`
	Model struct {
		ID           string `json:"id"`
		Manufacturer string `json:"manufacturer"`
		Name         string `json:"name"`
		Image        string `json:"image"`
	} `json:"model,omitempty"`
	Type       string        `json:"type"`
	Nickname   string        `json:"nickname"`
	Image      string        `json:"image"`
	Settings   interface{}   `json:"settings"`
	Aircon     interface{}   `json:"aircon"`
	Signals    []interface{} `json:"signals"`
	SmartMeter struct {
		EchonetliteProperties []struct {
			Name      string    `json:"name"`
			Epc       int       `json:"epc"`
			Val       string    `json:"val"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"echonetlite_properties"`
	} `json:"smart_meter,omitempty"`
}

// RawData : 取得した生データを格納する構造体
type RawData struct {
	ModelID                                  string
	Coefficient                              int
	CumulativeElectricEnergyEffectiveDigits  int
	CumulativeElectricEnergyUnit             int
	NormalDirectionCumulativeElectricEnergy  int
	ReverseDirectionCumulativeElectricEnergy int
	MeasuredInstantaneous                    int
}

// GetRawData : データを取得する関数
func (r *RemoClient) GetRawData() ([]RawData, error) {
	var info []remoInfo

	client := resty.New()
	resp, err := client.R().SetAuthToken(r.Token).Get("https://api.nature.global/1/appliances")
	if err != nil {
		return []RawData{}, err
	}

	err = json.Unmarshal(resp.Body(), &info)
	if err != nil {
		return []RawData{}, err
	}

	var rawDataList []RawData
	for _, d := range info {
		if len(d.SmartMeter.EchonetliteProperties) != 0 {
			var rd RawData
			rd.ModelID = d.Model.ID
			for _, p := range d.SmartMeter.EchonetliteProperties {
				switch p.Epc {
				case 211:
					rd.Coefficient, err = strconv.Atoi(p.Val)
					if err != nil {
						return []RawData{}, err
					}
				case 215:
					rd.CumulativeElectricEnergyEffectiveDigits, err = strconv.Atoi(p.Val)
					if err != nil {
						return []RawData{}, err
					}
				case 224:
					rd.NormalDirectionCumulativeElectricEnergy, err = strconv.Atoi(p.Val)
					if err != nil {
						return []RawData{}, err
					}
				case 225:
					rd.CumulativeElectricEnergyUnit, err = strconv.Atoi(p.Val)
					if err != nil {
						return []RawData{}, err
					}
				case 227:
					rd.ReverseDirectionCumulativeElectricEnergy, err = strconv.Atoi(p.Val)
					if err != nil {
						return []RawData{}, err
					}
				case 231:
					rd.MeasuredInstantaneous, err = strconv.Atoi(p.Val)
					if err != nil {
						return []RawData{}, err
					}
				default:
					return []RawData{}, fmt.Errorf("未定義のEPC値が検出されました : %d", p.Epc)
				}
			}
			rawDataList = append(rawDataList, rd)
		}
	}

	return rawDataList, nil
}

// NewClient : クライアントを作成する関数
func NewClient(token string) RemoClient {
	return RemoClient{Token: token}
}

// GetPowerCunsumption : 積算電力消費量を計算する関数
func GetPowerCunsumption(r RawData) float64 {
	return (float64(r.NormalDirectionCumulativeElectricEnergy * r.Coefficient)) / (10 * float64(r.CumulativeElectricEnergyUnit))
}

// GetPowerCunsumptionDiff : 特定の地点かjらの積算電力消費量の差を計算する関数
func GetPowerCunsumptionDiff(r RawData, p float64) float64 {
	now := GetPowerCunsumption(r)
	if now >= p {
		return now - p
	}
	max := math.Pow(10, float64(r.CumulativeElectricEnergyEffectiveDigits+1))
	return max - p + now
}
