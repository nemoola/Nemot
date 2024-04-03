package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type EEW struct {
	urlBuilder url.URL
	time       struct {
		now       time.Time
		today     string
		timestamp string
	}
	Img  bytes.Buffer
	data EEWData
}

func NewEEW() *EEW {
	return &EEW{
		urlBuilder: url.URL{
			Scheme: "http",
			Host:   "www.kmoni.bosai.go.jp",
		},
	}
}

func (eew *EEW) GenerateEEWIMG() bytes.Buffer {
	urlBuilder := eew.urlBuilder
	urlBuilder.Path = fmt.Sprintf("/data/map_img/PSWaveImg/eew/%s/%s.eew.gif", eew.time.today, eew.time.timestamp)
	res, _ := http.Get(urlBuilder.String())
	waveImg, _ := gif.Decode(res.Body)

	urlBuilder.Path = fmt.Sprintf("/data/map_img/EstShindoImg/eew/%s/%s.eew.gif", eew.time.today, eew.time.timestamp)
	res, _ = http.Get(urlBuilder.String())
	estShindoImg, _ := gif.Decode(res.Body)

	urlBuilder.Path = fmt.Sprintf("/data/map_img/RealTimeImg/jma_s/%s/%s.jma_s.gif", eew.time.today, eew.time.timestamp)
	res, _ = http.Get(urlBuilder.String())
	shindoImg, _ := gif.Decode(res.Body)
	if res.StatusCode != 200 {
		time.Sleep(100 * time.Millisecond)
		return eew.GenerateEEWIMG()
	}

	baseImgFile, _ := os.Open("assets/base_map_w.gif")
	defer baseImgFile.Close()
	scaleImgFile, _ := os.Open("assets/nied_jma_s_w_scale.gif")
	defer scaleImgFile.Close()
	baseImg, _ := gif.Decode(baseImgFile)
	scaleImg, _ := gif.Decode(scaleImgFile)

	img := image.NewRGBA(baseImg.Bounds())
	draw.Draw(img, img.Bounds(), baseImg, baseImg.Bounds().Min, draw.Src)
	draw.Draw(img, img.Bounds(), scaleImg, scaleImg.Bounds().Max.Sub(baseImg.Bounds().Max), draw.Over)
	draw.Draw(img, img.Bounds(), waveImg, baseImg.Bounds().Min, draw.Over)
	draw.Draw(img, img.Bounds(), estShindoImg, baseImg.Bounds().Min, draw.Over)
	draw.Draw(img, img.Bounds(), shindoImg, baseImg.Bounds().Min, draw.Over)

	buf := bytes.Buffer{}
	_ = png.Encode(&buf, img)
	return buf
}

type EEWData struct {
	Result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		IsAuth  bool   `json:"is_auth"`
	} `json:"result"`
	ReportTime    string `json:"report_time"`
	RegionCode    string `json:"region_code"`
	RequestTime   string `json:"request_time"`
	RegionName    string `json:"region_name"`
	Longitude     string `json:"longitude"`
	IsCancel      bool   `json:"is_cancel"`
	Depth         string `json:"depth"`
	Calcintensity string `json:"calcintensity"`
	IsFinal       bool   `json:"is_final"`
	IsTraining    bool   `json:"is_training"`
	Latitude      string `json:"latitude"`
	OriginTime    string `json:"origin_time"`
	Security      struct {
		Realm string `json:"realm"`
		Hash  string `json:"hash"`
	} `json:"security"`
	Magunitude      string `json:"magunitude"`
	ReportNum       string `json:"report_num"`
	RequestHypoType string `json:"request_hypo_type"`
	ReportID        string `json:"report_id"`
}

func (eew *EEW) GetEEWData() EEWData {
	urlBuilder := eew.urlBuilder
	urlBuilder.Path = fmt.Sprintf("/webservice/hypo/eew/%s.json", eew.time.timestamp)
	res, _ := http.Get(urlBuilder.String())
	body, _ := io.ReadAll(res.Body)

	var eewData EEWData
	_ = json.Unmarshal(body, &eewData)
	return eewData
}

func (eew *EEW) GetEEW() {
	eew.time.now = time.Now().Add(-550 * time.Millisecond)
	eew.time.today = eew.time.now.Format("20060102")
	eew.time.timestamp = eew.time.now.Format("20060102150405")

	eew.data = eew.GetEEWData()
	eew.Img = eew.GenerateEEWIMG()
}
