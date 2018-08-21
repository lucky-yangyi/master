var isLogged = 0;
var isMakeCall = 0;
var calling = 0;
var continueTime = 0;
var phoneTimer;
var ringing = 0;
var callingTimer;
var callingSeconds = 0;
function LoginORLogout(aid, apwd, adn, eid) {
    var result;
    if (isLogged == 0) {

        result = objocx.TxnProcess("Agent/Login", "AgentId=" + aid + "&AgentPwd=" + apwd + "&AgentDN=" + adn + "&EntId=" + eid);
        if (result == "0") {
            $(".clickToLogin").html("一键登出");
        }

    } else {
        result = objocx.TxnProcess("Agent/Logout", "");
        if (result == "0") {
            $(".clickToLogin").html("一键登录");
        }

    }
}
function SetReady() {
    var result = objocx.TxnProcess("State/Ready", "");
}
function StatusReset() {
    var result = objocx.TxnProcess("State/Reset", "");
}
function SetNotReady() {
    var result = objocx.TxnProcess("State/NotReady", "");
}

function SetAway(substate) {
    var result = objocx.TxnProcess("State/Away", "SubState=" + substate);
}

function DoAnswer() {
    var result = objocx.TxnProcess("SIP/Answer", "");
}

function DoHangup() {
    var result = objocx.TxnProcess("CallControl/Hangup", "");
    $(".clickToCall").html("拨号");
}
function doMonitor(AgentId) {//监听
    var result = objocx.TxnProcess("CallControl/SilentMonitor", "AgentID=" + AgentId);
}
function getTelephone(type) {
    var result = objocx.TxnProcess("getAgentCallInfo", "");
    return GetQueryString(result, type)
}
function NoANICall(outTarget) {
    var result = objocx.TxnProcess("CallControl/MakeCall", "TargetDN=" + outTarget + "&MakeCallType=2");
    if (result == 0) {
        // DisabledReadyAndNotReady();
        isMakeCall = 1;
    }
}

function InterCall(interTarget) {
    var result = objocx.TxnProcess("CallControl/MakeCall", "TargetDN=" + interTarget + "&MakeCallType=3");
    if (result == 0) {
        // DisabledReadyAndNotReady();
        isMakeCall = 1;
    }
}

function ANICall(target, ani) {
    var result = objocx.TxnProcess("CallControl/MakeCall", "TargetDN=" + target + "&MakeCallType=5" + "&ShowANI=" + ani);
    if (result == 0) {
        // DisabledReadyAndNotReady();
        $(".clickToCall").html("挂断");
        callingTimer = setInterval(function(){
            var tips = "通话时间"+formatSeconds(callingSeconds)
            callingSeconds++;
            $('.modal-tips').html(tips)
        },1000)
        dialog.init({
            title : "通话中",
            tips : "通话时间"+formatSeconds(callingSeconds),
            buttons : [{
                text :"挂断",
                style : "green",
                func : function(){
                    DoHangup();
                    $(".J_modal").remove();
                    clearInterval(callingTimer);
                    callingSeconds = 0
                }
            }]
        })
        isMakeCall = 1;
    }
}

function recordTime(callType){//记录通话时间
    var OCXaid = OCXaid,OCXapwd=OCXapwd;
    var result = objocx.TxnProcess("getAgentCallInfo","");
    var sessionId = GetQueryString1(result,"SessionID");
    //var sessionId = GetQueryString1(result,"A");
    var phone = GetQueryString(result,callType);
    var type = "",callClass = "";
    switch(callType){
        case "ANI=TEL" : type = "呼入",callClass = "呼入"; break;
        case "DNIS=TEL" : type = "呼出",callClass = "呼出"; break;
        case "CA_CALL_EVENT_OP_DISCONNECT" : type = "挂断"; break;
        case "CA_CALL_EVENT_TP_DISCONNECT" : type = "挂断"; break;
        case "CA_CALL_EVENT_TRANSFER" : type = "挂断" ;break;//转移 视为挂断
        default : type="挂断"
    }
    var url = "";
    if(type == "挂断"){
        url = "../../taperecord/hanguptaperecord";
    }else{
        url = "../../taperecord/addtaperecord";
    }
    var data = {
        callType : type,
        sessionId : sessionId,
        phone : phone,
        user : OCXaid,
        password : OCXapwd,
        callClass : callClass,
        systemId : systemId
    };
    $.ajax({//记录开始时间
        type : "get",
        url : url,
        data : data,
        success : function(rslt){

        }
    })
}

function ReadyAgents() {
    var result = objocx.TxnProcess("Data/GetOnlineAgent", "ClientReqType=READY");
    var ressplit = result.split("|");
    // document.all.txtResult.value = ressplit.length;
    //var agent_array = new Array();
    // var agent_list = document.getElementById("intargetdn");
    // agent_list.length = 0;
    // for (var i = 0; i < ressplit.length - 1; i++) {
    //     var temp = ressplit[i].split(':');
    //     //alert(temp[0]);
    //     agent_list.add(new Option(temp[0]));
    // }
    //document.all.txtResult.value = agent_list.length + 10;
}

function TxnProcess(func, prama) {
    var result = objocx.TxnProcess(document.all.func.value, document.all.param.value);
}

function FillSubAway() {
    var sublist = objocx.TxnProcess("Data/SubAway", "");
    // if (sublist != "null") {
    //     var subsplit = sublist.split(';');
    //     // document.all.txtResult.value = subsplit.length;
    //     //var agent_array = new Array();
    //     var sub_list = document.getElementById("subaway");
    //     sub_list.length = 0;
    //     for (var i = 0; i < subsplit.length - 1; i++) {
    //         var subtemp = subsplit[i].split('=');
    //         //alert(temp[0]);
    //         sub_list.add(new Option(subtemp[0], subtemp[1]));
    //     }
    // }
}


function GetQueryString(backString, name) {  //获取电话号码，显示在页面中  xx=xxx&xxx=xx:111
    var reg = new RegExp("(^|&)" + name + ":([^&]*)(&|$)", "i");
    var r = backString.match(reg);  //获取url中"?"符后的字符串并正则匹配
    var context = "";
    if (r != null)
        context = r[2];
    reg = null;
    r = null;
    return context == null || context == "" || context == "undefined" ? "" : context;
}
function GetQueryString1(backString, name) {  //获取录音sessionID，发送给后端  规则为  xx=xxx&xx=xxx
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = backString.match(reg);  //获取url中"?"符后的字符串并正则匹配
    var context = "";
    if (r != null)
        context = r[2];
    reg = null;
    r = null;
    return context == null || context == "" || context == "undefined" ? "" : context;
}