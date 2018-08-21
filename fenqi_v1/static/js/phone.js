var valueCM = "";
function formatSeconds(s) {
	var t;
	if(s > -1){
		var hour = Math.floor(s/3600);
		var min = Math.floor(s/60) % 60;
		var sec = s % 60;
		if(hour < 10) {
			t = '0'+ hour + ":";
		} else {
			t = hour + ":";
		}

		if(min < 10){t += "0";}
		t += min + ":";
		if(sec < 10){t += "0";}
		t += sec;
	}
	return t;
}
switch (value) {
	case "OCX_LOADED": valueCM = "OCX控件加载完成"; break;
	case "CA_CALL_EVENT_ALERTING": valueCM = "坐席振铃"; break;
	case "CA_CALL_EVENT_CONNECTED": valueCM = "双方通话"; break;
	case "CA_CALL_EVENT_OP_DISCONNECT": valueCM = "客户方挂断"; break;
	case "CA_CALL_EVENT_TP_DISCONNECT": valueCM = "坐席方挂断"; break;
	case "CA_CALL_EVENT_INTERNAL_ALERTING_TP": valueCM = "内呼本方振铃"; break;
	case "CA_CALL_EVENT_INTERNAL_CONNECTED_TP": valueCM = "内呼本方接通"; break;
	case "CA_CALL_EVENT_INTERNAL_ALERTING_OP": valueCM = "内呼对方振铃"; break;
	case "CA_CALL_EVENT_OUTBOUND_ALERTING_TP": valueCM = "外呼本方振 铃"; break;
	case "CA_CALL_EVENT_OUTBOUND_CONNECTED_TP": valueCM = "外呼本方接通"; break;
	case "CA_CALL_EVENT_OUTBOUND_ALERTING_OP": valueCM = "外呼对方振铃"; break;
	case "CA_CALL_EVENT_OUTBOUND_CONNECTED_OP": valueCM = "外呼对方接通"; break;
	case "CA_CALL_EVENT_HOLD": valueCM = "坐席保持"; break;
	case "CA_CALL_EVENT_HOLD_RETRIEVE": valueCM = "保持接回"; break;
	case "CA_CALL_EVENT_AGENT_BEHOLD": valueCM = "坐席被保持"; break;
	case "CA_CALL_EVENT_AGENT_BEUNHOLD": valueCM = "坐席被保持接回"; break;
	case "CA_CALL_EVENT_CONSULT_ALERTING_TP": valueCM = "咨询本方振铃"; break;
	case "CA_CALL_EVENT_CONSULT_CONNECTED_TP": valueCM = "咨询本方通话"; break;
	case "CA_CALL_EVENT_CONSULT_ALTERTING_OP": valueCM = "咨询对方振铃"; break;
	case "CA_CALL_EVENT_CONSULT_CONNECTED_OP": valueCM = "咨询对方通话"; break;
	case "CA_CALL_EVENT_CONSULT_RETRIEVE": valueCM = "咨询接回"; break;
	case "CA_CALL_EVENT_TRANSFER": valueCM = "转移"; break;
	case "CA_CALL_EVENT_CONFERENCE": valueCM = "会议"; break;
	case "CA_CALL_EVENT_MONITOR_ALERTING": valueCM = "监听振铃"; break;
	case "CA_CALL_EVENT_MONITOR": valueCM = "监听"; break;
	case "CA_CALL_EVENT_INTRUDE": valueCM = "强插"; break;
	case "CA_CALL_EVENT_TEARDOWN": valueCM = "强拆"; break;
	case "CA_CALL_EVENT_CLEAR_CALL": valueCM = "全拆"; break;
	case "CA_CALL_EVENT_THIRD_PARTY_DISCONNECT": valueCM = "第三方挂断"; break;
	case "CA_CALL_EVENT_FOURTH_PARTY_DISCONNECT": valueCM = "第四方挂断"; break;
	case "CA_CALL_EVENT_SINGLESTEPTRANSFER_CONNECTED_TP": valueCM = "单转通话(本方)"; break;
	case "CA_CALL_EVENT_SINGLESTEPTRANSFER_CONNECTED_OP": valueCM = "单转通话(对方)"; break;
	case "CA_CALL_EVENT_SINGLESTEPTRANSFER": valueCM = "单转振铃(本方)"; break;
	case "SFER_ALERTING_TP": valueCM = ""; break;
	case "CA_CALL_EVENT_SINGLESTEPTRANSFER_ALERTING_OP": valueCM = "单转振铃(对方)"; break;
	case "EVENT_SERVER_CLOSED": valueCM = "服务器主动关闭"; break;
	case "EVENT_LOGIN_SUCCESS": valueCM = "登陆成功"; break;
	case "EVENT_LOGIN_FAILED": valueCM = "登陆失败"; break;
	case "EVENT_LOGOUT_SUCCESS": valueCM = "登出成功"; break;
	case "EVENT_LOGOUT_FAILED": valueCM = "登出失败"; break;
	case "EVENT_SET_STATE_SUCCESS": valueCM = "设置状态成功"; break;
	case "EVENT_SET_STATE_FAILED": valueCM = "设置状态失败"; break;
	case "EVENT_CC_SUCCESS": valueCM = "呼叫控制调用成功"; break;
	case "EVENT_CC_FAILED": valueCM = "呼叫控制调用失败"; break;
	case "EVENT_MODIFY_SUCCESS": valueCM = "修改密码成功"; break;
	case "EVENT_MODIFY_FAILED": valueCM = "修改密码失败"; break;
	case "AGENTSTATUS_READY": valueCM = "空闲状态(可以接电话)"; break;
	case "AGENTSTATUS_AWAY": valueCM = "离开状态（只可以重置）"; break;
	case "AGENTSTATUS_NOTREADY": valueCM = "客座席未就绪状态(置忙, 刚登录后的默认状态)"; break;
	case "AGENTSTATUS_LOCKED": valueCM = "预锁定状态(已被服务端选中，稍后有电话呼入)"; break;
	case "AGENTSTATUS_WORKING": valueCM = "工作状态(通话中)"; break;
	case "AGENTSTATUS_ACW": valueCM = "事后整理状态"; break;
	case "AGENTSTATUS_OTHERWORK": valueCM = "座席未就绪状态(置忙, 刚登录后的默认状态)"; break;
	case "AGENTSTATUS_LOGOUT": valueCM = "登录前的状态"; break;
	case "AGENTSTATUS_UNKOWN": valueCM = "未知状态"; break;
	case "CONN_MESSAGE": valueCM = "连接消息"; break;
	case "RECONN_MESSAGE": valueCM = "重连消息"; break;
	case "BUSYREQUEST_MESSAGE": valueCM = "服务端忙消息"; break;
	case "RESULT_GET_IDLE_AGENT_LIST_RESP_SUCCESS": valueCM = "返回获取在线空闲坐席列表成功"; break;
	case "RESULT_GET_ALL_AGENT_LIST_RESP_SUCCESS": valueCM = "返回获取所有坐席列表成功"; break;
	case "RESULT_GET_WORKING_AGENT_LIST_RESP_SUCCESS": valueCM = "返回在线工作坐席成功"; break;
	case "RESULT_GET_ONLINE_AGENT_LIST_RESP_SUCCESS": valueCM = "返回获取在线坐席列表成功"; break;
	case "RESULT_GET_ALL_SKILL_NAME_RESP_SUCCESS": valueCM = "返回获取所有技能组列表成功"; break;
	case "RESULT_RESP_FAILED": valueCM = "返回操作失败"; break;
	case "RESULT_GET_ASSOCIATE_DATA_RESP_SUCCESS": valueCM = "返回获取随路数据成功"; break;
	case "RESULT_SET_ASSOCIATE_DATA_RESP_SUCCESS": valueCM = "返回设置随路数据成功"; break;
	case "WM_SIP_RETURN_SUCCESS": valueCM = "软电话登录成功"; break;
	case "WM_SIP_RETURN_GENERALERROR": valueCM = "普通错误"; break;
	case "WM_SIP_RETURN_FATALERROR": valueCM = "严重错误"; break;
	case "WM_SIP_RETURN_BADARGUMENT": valueCM = "参数错误"; break;
	case "WM_SIP_RETURN_TIMEOUT": valueCM = "时间超时"; break;
	case "WM_SIP_RETURN_OTHER": valueCM = "其他错误"; break;
	case "AGENT_DETECTION_FAILED_DCMSWEBSERVICE_ERROR": valueCM = "dcmsWebservice服务不可用，请查看IP及端口"; break;
	case "AGENT_DETECTION_FAILED_UCDS_ERROR": valueCM = "UCDS地址或端口不可用；或服务未启动"; break;
	case "AGENT_DETECTION_FAILED_UCDS_IP_UNAVAILABLE": valueCM = "登录UCDS 地址不可达，IP地址ping不通"; break;
	case "AGENT_DETECTION_FAILED_GLS_IP_UNAVAILABLE": valueCM = "由GLS获取UCDS信息，GLS地址不可达"; break;
	case "AGENT_DETECTION_FAILED_GLS_ERROR": valueCM = "由GLS获取UCDS信息，GLS地址不准确（IP可ping同，但不是GLS服务地址）或端口错误"; break;
	case "AGENT_DETECTION_FAILED_GLS_NOUCDSINFO": valueCM = "由GLS获取UCDS信息，获取不到UCDS信息"; break;
	case "8": valueCM = "非4.1平台有门户获取UCDS信息，门户webservice服务不可用，（扩展用，当前不可用这种登录方式）"; break;
	case "9": valueCM = "非4.1平台有门户获取UCDS信息，获取不到ucds信息（扩展用，当前不可用这种登录方式）"; break;
	case "10": valueCM = "非4.1平台有门户获取UCDS信息，UCDS信息错误（扩展用，当前不可用这种登录方式）"; break;
	case "AGENT_DETECTION_FAILED_SIPSERVER_IP_UNAVAILABLE": valueCM = "使用OCX内置SIP软电话时，SIP服务器地址不可达"; break;
	case "AGENT_DETECTION_FAILED_SIPSERVER_ERROR": valueCM = "使用OCX内置SIP软电话时，SIP服务器地址不准确或端口错误"; break;
	case "AGENT_DETECTION_FALLED_DCMSURL_LOSEWEBSEVICE": valueCM = "DCMS门户webservice地址缺少dcmsWebservice字符串"; break;
	case "AGENT_DETECTION_FAILED_PARAMS_ERROR": valueCM = "企业ID、坐席ID或密码错误"; break;
	case "AGENT_DETECTION_DCMSURL_PASSED": valueCM = "DCMS的webservice接口URL可用"; break;
	case "AGENT_DETECTION_UCDS_PASSED": valueCM = "UCDS服务可用"; break;
	case "AGENT_DETECTION_GLS_PASSED": valueCM = "GLS服务可用"; break;
	case "AGENT_DETECTION_SIPSERVER_PASSED": valueCM = "SIP服务器可以"; break;
	case "SSED": valueCM = ""; break;
	case "AGENT_DETECTION_PASSED": valueCM = "坐席检测通过"; break;
}
if (value == "OCX_LOADED"){		//OCX控件加载事件
	//  var OCXaid = {{.QnAccount}},OCXapwd={{.QnPassword}},OCXadn = {{.QnAccount}},OCXeid="0101290166";
	// if(!(OCXaid&&OCXapwd&&OCXaid!=""&&OCXapwd!="")){
	// 	$("#UsbossViewer").hide();
	// 	$(".clickToLogin").hide();
	// 	return
	// }
}
if (value == "EVENT_LOGIN_SUCCESS") {
	// document.getElementById("callingPhone").innerHTML = " ";
	clearInterval(phoneTimer);
	continueTime = 0;
	isLogged = 1;
    $(".call_outSide").show();
	FillSubAway();
	// SetReady();//登录自动置闲
}
if (value == "EVENT_LOGOUT_SUCCESS") {
	isLogged = 0;
	// document.getElementById("callingPhone").innerHTML = " ";
	clearInterval(phoneTimer);
	continueTime = 0;
	if(calling == 1){//如果在通话状态，登出视为挂断
		recordTime("EVENT_LOGOUT_SUCCESS");
	}
}
if (value == "AGENTSTATUS_NOTREADY"){		// 置忙状态
	
}
if (value == "AGENTSTATUS_READY"){			// 置闲状态

}
if (value == "AGENTSTATUS_ACW") {			// 事后整理状态
	isMakeCall = 0;
}
if (value == "AGENTSTATUS_AWAY"){			// 离开状态

}
if (value == "EVENT_CC_FAILED") {       	// 操作失败

	if (isMakeCall == 1) {			//外呼操作失败，返回置忙状态

	}

}
if (value == "CA_CALL_EVENT_TP_DISCONNECT"){	// 坐席挂断事件

}
if (value == "CA_CALL_EVENT_ALERTING"){	// 入呼叫振铃

}
if(value == "CA_CALL_EVENT_OUTBOUND_CONNECTED_OP"){ //外呼对方接通
	
	var ssk = getTelephone("DNIS=TEL");
	recordTime("DNIS=TEL");
	clearInterval(phoneTimer);
	continueTime = 0;
	phoneTimer = setInterval(function(){
		continueTime++;
		// document.getElementById("callingPhone").innerHTML = "通话时间"+""+formatSeconds(continueTime);
	},1000);
}else if(value == "CA_CALL_EVENT_ALERTING"){ //坐席振铃事件 
	ringing = 1;
}else if (value == "CA_CALL_EVENT_CONNECTED"){	// 入呼叫接通事件
	calling = 1;
	
}
if(value == 'CA_CALL_EVENT_OUTBOUND_CONNECTED_TP'){//外呼本方接通  代表已拨打
	$.ajax({
		url : "/telesale/teleconnect",
		data : {
			telephone : $(".clickToCall.calling").data("telephone"),
			id : $(".clickToCall.calling").data("id")
		},
		success : function(rslt){
			
		}
	})
}
if (value == "AGENTSTATUS_WORKING") {
	calling = 1
}
if (value == "CA_CALL_EVENT_INTERNAL_ALERTING_TP") {
	
}
if (value == "CA_CALL_EVENT_INTERNAL_CONNECTED_TP") {

}
if (value == "CA_CALL_EVENT_OUTBOUND_ALERTING_TP") {
	isMakeCall = 0;
}
if (value == "CA_CALL_EVENT_OUTBOUND_CONNECTED_TP") {
	isMakeCall = 0;
}
if (value == "WM_SIP_RETURN_SUCCESS") { //软电话登录成功 视为 登录成功
	// SetReady();//登录自动置闲
}
if(value == "AGENTSTATUS_NOTREADY"){//重置
	// SetReady();//点击重置自动置闲
}
if(value == "CA_CALL_EVENT_TRANSFER"){//转移 视为挂断
	calling = 0;
	clearInterval(phoneTimer);
	continueTime = 0;
	recordTime("CA_CALL_EVENT_TRANSFER");
}

if(value == "CA_CALL_EVENT_CONNECTED"){
	
	if(ringing == 1){//先振铃后双方通话  代表呼入接通
		var ssk = getTelephone("ANI=TEL");
		recordTime("ANI=TEL");
		clearInterval(phoneTimer);
		continueTime = 0;
		phoneTimer = setInterval(function(){
			continueTime++;
			// document.getElementById("callingPhone").innerHTML = "通话时间"+"  "+formatSeconds(continueTime);
		},1000);
	}
}


if (value == "CA_CALL_EVENT_OP_DISCONNECT") { //客户方挂断
	$("#mp3_main").attr("src","/static/MyPlayer/guaduan.mp3");
	// document.getElementById("callingPhone").innerHTML = " ";
	ringing = 0
	calling = 0;
	clearInterval(phoneTimer);
	continueTime = 0;
	SetNotReady()
	$(".clickToCall").html("拨号")
	$(".J_modal").remove();
	recordTime("CA_CALL_EVENT_OP_DISCONNECT");
}
if(value == "CA_CALL_EVENT_TP_DISCONNECT"){//坐席方挂断
	// document.getElementById("callingPhone").innerHTML = " ";
	calling = 0;
	ringing = 0
	$("#mp3_main").attr("src","/static/MyPlayer/guaduan.mp3");
	clearInterval(phoneTimer);
	continueTime = 0;
	SetNotReady() 
	$(".clickToCall").html("拨号");
	$(".J_modal").remove();
	recordTime("CA_CALL_EVENT_TP_DISCONNECT");
}

document.getElementById("Rt").value += valueCM + "\r\n";
document.getElementById("Rt").scrollTop = document.getElementById("Rt").scrollHeight;
document.getElementById("Rt").style.display = "block";
clearTimeout(showTimer);
var showTimer = setTimeout(function () {
	document.getElementById("Rt").style.display = "none";
}, 3000)