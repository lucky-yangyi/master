package services

import (
	"fenqi_v1/utils"
	"strconv"
	"strings"
	"time"
)

// 基础参数
type PdfBaseParameter struct {
	BorrowerName    string    // 借款人姓名
	BorrowerIdCard  string    // 借款人身份证号
	LoanDate        time.Time //放款日期
	TotalLoanAmount float64   // 待还总金额
	HandleName      string    // 所属催收员姓名
	HandlePhone     string    //所属催收员联系电话
	ContractCode    string    //合同编号
}

//缴款通知书
func GeneratePDFDemandNote(param *PdfBaseParameter) *utils.XJDPdf {
	loanDate := strings.Split(param.LoanDate.Format("2006-01-02"), "-")
	loanDateYear := loanDate[0]
	loanDateMonth := loanDate[1]
	loanDateDay := loanDate[2]
	currDate := strings.Split(time.Now().Format("2006-01-02"), "-")
	currYear := currDate[0]
	currMonth := currDate[1]
	currDay := currDate[2]
	var name = ""
	if len(param.HandleName) > 2 {
		name = param.HandleName[0:3] + "经理"
	}
	xjdpdf := utils.NewXJDPdf(nil)
	xjdpdf.ZAddPage()
	xjdpdf.SetFontFamily(utils.FONT_NAME_BOLD).WriteInCenter("缴款通知书", 21).DoBR(utils.GetBR(18))

	xjdpdf.SetDefaultFontSize().SetFontFamily(utils.FONT_NAME).DoBR(utils.GetBR(18))
	xjdpdf.WriteInCenterWithLineFixedlength(param.BorrowerName, 8, 13).Write("先生/女士（身份证号码：").
		WriteInCenterWithLineFixedlength(param.BorrowerIdCard, 18, 13).Write("）：").DoBR(utils.GetBR(18))
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 18, `  我方受杭州有个金融服务外包有限公司（以下简称有个金服）委托，就您办理的“花无忧”贷款逾期一事致函给您，希望您能尽快予以处理，避免因进一步的法律行动给您造成更大的损失和不利法律后果。`)

	xjdpdf.Write("    您于").WriteInCenterWithLineFixedlength(loanDateYear, 4, 13).Write("年").
		WriteInCenterWithLineFixedlength(loanDateMonth, 4, 13).Write("月").
		WriteInCenterWithLineFixedlength(loanDateDay, 4, 13).Write("日起向有个金服申请“花无忧”贷款，截至").
		WriteInCenterWithLineFixedlength(currYear, 4, 13).Write("年").
		WriteInCenterWithLineFixedlength(currMonth, 4, 13).Write("月").
		WriteInCenterWithLineFixedlength(currDay, 4, 13).DoBR(utils.GetBR(18))
	xjdpdf.Write(`日您已累计欠款共计人民币`).WriteInCenterWithLineFixedlength(strconv.FormatFloat(param.TotalLoanAmount, 'f', -1, 64), 10, 13).Write("元（具体金额以有个金服系统查询为准）。").DoBR(utils.GetBR(18))

	xjdpdf.SetDefaultFontSize().Write(`    我们郑重通知您，请您于收到此函之日起`).WriteInCenterWithLineFixedlength("2", 4, 13).WriteMassage(`日内，通过有个金服指定的还款方式缴`).DoBR(utils.GetBR(18))
	xjdpdf.Write(`清全部款项。若未在本函规定时间内偿还，有个金服将有权采取以下措施：`).DoBR(utils.GetBR(18))

	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 18, `  1、依法就您涉嫌“贷款诈骗”向公安机关申请立案，一旦罪名成立，不但要清偿全部欠款本息，而且面临判处有期徒刑等刑罚。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 18, `  2、依法向人民法院提起诉讼，请求法院对您名下财产采取查封、扣押、冻结、划拨等法律措施。届时，您将承担欠款本息、诉讼费、保全费、执行费、拍卖费等全部费用。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 18, `  3、将您的逾期状况如实上报中国人民银行，这将严重影响您在央行个人征信系统中的信用记录，也会对您日后与银行等金融机构发生借贷业务产生重大影响。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 18, `  4、在合法范围内寄送《律师函》至您户籍地、居住地、所在单位或上级单位，由此造成的不利影响由您本人承担。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 18, `  如果您收到此函时已还清欠款，则不必理会此函，由此给您带来的不便，敬请谅解。`).DoBR(utils.GetBR(18)).DoBR(utils.GetBR(18)).DoBR(utils.GetBR(18))

	xjdpdf.SetDefaultFontSize().Write("    联系人：").WriteInCenterWithLineFixedlength(name, 8, 13).DoBR(utils.GetBR(18))
	xjdpdf.SetDefaultFontSize().Write("    联系电话：").WriteInCenterWithLineFixedlength(param.HandlePhone, 11, 13).DoBR(utils.GetBR(18))
	xjdpdf.SetDefaultFontSize().Write("    客服热线："+utils.Service_Moblie).WriteAnyPlace(0.64, currYear+" 年 "+currMonth+" 月 "+currDay+" 日", 13).DoBR(utils.GetBR(18))

	return xjdpdf
}

//债务律师催告函
func GeneratePDFAttorneyLetter(param *PdfBaseParameter) *utils.XJDPdf {
	loanDate := strings.Split(param.LoanDate.Format("2006-01-02"), "-")
	loanDateYear := loanDate[0]
	loanDateMonth := loanDate[1]
	loanDateDay := loanDate[2]
	//合同编号
	currDate := strings.Split(time.Now().Format("2006-01-02"), "-")
	currYear := currDate[0]
	currMonth := currDate[1]
	currDay := currDate[2]
	var name = ""
	if len(param.HandleName) > 2 {
		name = param.HandleName[0:3] + "经理"
	}
	xjdpdf := utils.NewXJDPdf(nil)
	xjdpdf.ZAddPage()
	xjdpdf.SetFontFamily(utils.FONT_NAME_BOLD).WriteInCenter("债务律师催告函", 21).DoBR(utils.GetBR(15))

	xjdpdf.SetDefaultFontSize().SetFontFamily(utils.FONT_NAME)
	xjdpdf.WriteInCenterWithLineFixedlength(param.BorrowerName, 8, 13).Write("先生/女士：　　　　　　　　　身份证号：").WriteInCenterWithLineFixedlength(param.BorrowerIdCard, 18, 13).DoBR(utils.GetBR(15)).DoBR(utils.GetBR(15))

	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  受花无无忧借款平台委托，我们向您郑重函告如下：`)
	xjdpdf.Write("    据我们了解，您").WriteInCenterWithLineFixedlength(loanDateYear, 4, 13).Write("年").
		WriteInCenterWithLineFixedlength(loanDateMonth, 4, 13).Write("月").
		WriteInCenterWithLineFixedlength(loanDateDay, 4, 13).Write("日通过花无无忧借款平台申请了金额为人民币元的").DoBR(utils.GetBR(15)).
		Write("一笔个人消费贷款／现金贷款（合同编号：").
		WriteInCenterWithLineFixedlength(param.ContractCode, 24, 13).
		Write("），但您并未如约").DoBR(utils.GetBR(15)).
		Write("还款。截止").WriteInCenterWithLineFixedlength(currYear, 4, 13).Write("年").
		WriteInCenterWithLineFixedlength(currMonth, 4, 13).Write("月").
		WriteInCenterWithLineFixedlength(currDay, 4, 13).
		Write(`日，您拖欠还款的总金额已达到人民币`).WriteInCenterWithLineFixedlength(strconv.FormatFloat(param.TotalLoanAmount, 'f', -1, 64), 10, 13).Write("元，并且").DoBR(utils.GetBR(15)).
		Write("超过还款期限天。但经花无无忧借款平台屡次催收，您仍拒不偿还。").DoBR(utils.GetBR(15))

	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  根据《中华人民共和国刑法》第二百二十四条、第一百九十三条，全国人民代表大会常务委员会《关于惩治破坏金融秩序犯罪的决定》第十条的规定，您的行为有可能因涉嫌合同诈骗或贷款诈骗而被追究刑事责任。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  基于以上事实和法律，我们催请您知悉和处理以下事项：`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  一、我们依法有权视情况决定是否就您涉嫌犯罪一事向公安机关报案，并发协查函到您户籍地派出所。一旦国家机关判令相关罪名成立，您不但要清偿全部欠款，且有可能面临相应刑事处罚。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  二、我们依法有权视情况决定是否通过专业律师事务所向人民法院提起诉讼，请求法院判令您清偿全部欠款并支付全部由此产生的律师费、诉讼费、财产保全费等款项。同时，我们依法有权向人民法院申请财产保全，请求法院依法查封、扣押、冻结您相应金额的个人财产（包括但不限于房产、车辆、存款、股票、债券等）。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  三、我们依法有权视情况决定是否委派专员至您的住所地、户籍所在地、您所在的居委会或村委会、工作单位等联络拜访并发送相关法律函件，促使相关国家机关配合我们对您的调查。`)
	xjdpdf.Write("  ").WritePassageAnyWidth(39, 0, 15, `  四、请您在收到本函后立即来电联系办理还款事宜。`)
	xjdpdf.Write("    ").WritePassageAnyWidth(39, 0, 15, `  请您联系以下人员`)
	xjdpdf.Write("      联系人：").WriteInCenterWithLineFixedlength(name, 11, 13).DoBR(utils.GetBR(15))
	xjdpdf.Write("      电  话：").WriteInCenterWithLineFixedlength(param.HandlePhone, 11, 13).DoBR(utils.GetBR(15))
	xjdpdf.Write("    ").WritePassageAnyWidth(39, 0, 15, `  客服热线电话：`+utils.Service_Moblie)

	xjdpdf.WriteAnyPlace(0.64, currYear+" 年 "+currMonth+" 月 "+currDay+" 日", 13).DoBR(utils.GetBR(15)).DoBR(utils.GetBR(15))

	xjdpdf.SetDefaultFontSize().SetFontFamily(utils.FONT_NAME_BOLD).Write("备注：如您已还清上述欠款，则不必理会此函。")

	return xjdpdf
}
