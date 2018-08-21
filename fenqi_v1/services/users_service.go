package services

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	Endpoint        string = "oss-cn-hangzhou.aliyuncs.com" //http://zcmlcimg.
	AccessKeyId     string = "wU7F2aVGFB683FzB"
	AccessKeySecret string = "cPQJuzJno8RvKSwSbvKA0LlKFOLbZx"
	Bucketname      string = "xwbxjd"
	Imghost         string = "https://xjstatic.zcmlc.com/"
	Upload_dir      string = "api/identify/"
)

//图片上传到阿里云
func UploadAliyun(filename, filepath string) (error, string) {
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		return err, "1"
	}
	bucket, err := client.Bucket(Bucketname)
	if err != nil {
		return err, "2"
	}
	path := Upload_dir + time.Now().Format("200601/20060102/")
	path += filename
	err = bucket.PutObjectFromFile(path, filepath)
	if err != nil {
		return err, "3"
	}
	path = Imghost + path
	return err, path
}

/*//图片上传到阿里云
func UploadActivityAliyun(filename, filepath string) (error, string) {
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		return err, "1"
	}
	bucket, err := client.Bucket(Bucketname)
	if err != nil {
		return err, "2"
	}
	path := "api/activity/image/"
	path += filename
	err = bucket.PutObjectFromFile(path, filepath)
	if err != nil {
		return err, "3"
	}
	path = Imghost + path
	return err, path
}

func CheckIsIdentifyFailId(uid int) bool {
	result := models.GetIdentifyFaildUids(uid)
	return result
}

func GetAnthenCount(uid int) bool {
	result := models.GetAnthenCount(uid)
	return result
}

//给用户通讯录去重
func DistinctList(list []models.UsersTelephoneDirectory) (distinctList []models.UsersTelephoneDirectory) {
	var phones []string
	for _, v := range list {
		if !strings.Contains(strings.Join(phones, ","), v.ContactPhoneNumber) {
			phones = append(phones, v.ContactPhoneNumber)
			distinctList = append(distinctList, v)
		}
	}
	return
}

type HttpResult struct {
	Ret   int     //状态吗
	Msg   string  //错误信息
	Score float64 //返回分
}

func HttpRequestXgScore(uid int) (xgScore float64, errmsg string) {
	status := "failure"
	errType := 0
	defer func() {
		models.AddXgRequestLog(uid, errType, xgScore, status, errmsg, "后台通过") //添加请求日志
	}()
	postData := `{"Uid":` + strconv.Itoa(uid) + `}`
	dataByte, err := http.HttpPost(utils.XJD_API_URL+utils.XgScoreRequestUrl, postData)
	if err != nil {
		errType = 1
		errmsg = "err==" + err.Error() + ",params==" + string(postData)
		return
	}
	var result HttpResult
	err = json.Unmarshal(dataByte, &result)
	if err != nil {
		errType = 2
		errmsg = "err==" + err.Error() + ",params==" + string(postData)
		return
	}
	if result.Ret != 200 {
		errType = 3
		errmsg = "ret=" + strconv.Itoa(result.Ret) + ",msg=" + result.Msg + ",params==" + string(postData)
		return
	}
	status = "success"
	xgScore = result.Score
	return
}

func XgScoreRiskScore(authCount, uid int, xgScore float64) {
	//西瓜信用分
	var riskScore float64
	var returnMsg string
	{
		returnMsg = "【个人信息】西瓜信用分" + strconv.FormatFloat(xgScore, 'f', 2, 64)
		switch {
		case xgScore < 550:
			riskScore = utils.XgRiskScore550
		case xgScore >= 550 && xgScore < 590:
			riskScore = utils.XgRiskScore595
		case xgScore >= 590 && xgScore < 615:
			riskScore = utils.XgRiskScore615
		case xgScore >= 615 && xgScore < 635:
			riskScore = utils.XgRiskScore635
		case xgScore >= 635 && xgScore < 655:
			riskScore = utils.XgRiskScore655
		case xgScore >= 655:
			riskScore = utils.XgRiskScoreGt655
		}
		models.AddFkRiskScore(uid, authCount, riskScore, returnMsg)
	}
}*/
