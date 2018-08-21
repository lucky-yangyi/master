package models

import (
	"time"
)

//学信网数据
type MdbMxXXWData struct {
	Uid        int       `bson:"uid"`
	MXXXWData  MXXXWData `bson:"data"`
	CreateTime time.Time `bson:"createtime"`
}

type MXXXWData struct {
	StudentInfoList []struct {
		SchoolName        string `bson:"school_name"`         //学校名称
		EduForm           string `bson:"edu_form"`            //学习形式
		EnrollmentTime    string `bson:"enrollment_time"`     //入学时间
		Level             string `bson:"level"`               //层次
		YardName          string `json:"yard_name"`           //分院
		LeaveSchoolTime   string `bson:"leave_school_time"`   //毕业时间
		Specialty         string `bson:"specialty"`           //专业
		Department        string `bson:"department"`          //系
		Status            string `bson:"status"`              //状态
		LengthOfSchooling string `bson:"length_of_schooling"` //学制
		ClassName         string `bson:"class_name"`          //班级
		EduType           string `bson:"edu_type"`            //学历类别
		StudentId         string `bson:"student_id"`          //学号
	} `bson:"studentInfo_list"`
	EducationList []struct {
		EnrollmentTime string `bson:"enrollment_time"` //入学时间
		EduLevel       string `bson:"edu_level"`       //层次
		EduForm        string `bson:"edu_form"`        //学习形式
		GraduateTime   string `bson:"graduate_time"`   //毕业时间
		EduType        string `bson:"edu_type"`        //学历类别
		Graduate       string `bson:"graduate"`        //状态
		GraduateSchool string `bson:"graduate_school"` //学校名称
		Specialty      string `bson:"specialty"`       //专业
	} `bson:"education_list"`
}

//鹏元MongoDB信息
type MdbPyInfo struct {
	Uid  int `bson:"uid"`
	Data struct {
		ReturnValue struct {
			CisReport []struct {
				PersonApplyScoreInfo struct {
					Score string `bson:"score"`
				} `bson:"PersonApplyScoreInfo"`
			} `bson:"cisReport"`
		} `bson:"returnValue"`
	} `bson:"data"`
	Createtime string `bson:"createtime"`
}
