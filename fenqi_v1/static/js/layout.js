Date.prototype.format = function (format) {
	var o = {
		"M+": this.getMonth() + 1, //month
		"d+": this.getDate(), //day
		"h+": this.getHours(), //hour
		"m+": this.getMinutes(), //minute
		"s+": this.getSeconds(), //second
		"q+": Math.floor((this.getMonth() + 3) / 3), //quarter
		"S": this.getMilliseconds() //millisecond
	}
	if (/(y+)/.test(format)) format = format.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
	for (var k in o) if (new RegExp("(" + k + ")").test(format))
		format = format.replace(RegExp.$1,RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length));
	return format;
}
var orderCountDownTimer_Arr = [],orderCountDownTime_Arr = [],orderTimeOut_Arr = [],creditTimeOut_Arr = [];
$(function () {
	jQuery('.main-content').css({ 'min-height': $(window).height() });
	// countdownanimation(3)
	$.zpost('/', {}, function (result) {
		if (result && result.ret == 200 && !!result.data) {
			var menu = '';
			$.each(result.data, function (i, m) {
				if (m.ChildMenu && m.ChildMenu.length > 0) {
					menu += '<li class="menu-list"><a href="#" ><span>' + m.DisplayName + '</span></a><ul class="sub-menu-list">';
					var subMenu = '';
					$.each(m.ChildMenu, function (j, n) {
						subMenu += '<li><a class="norefresh" url="' + n.ControlUrl + '" href="' + n.ControlUrl + (n.HomeUrl ? eval(n.HomeUrl) : "") + '"> ' + n.DisplayName + '</a></li>';
					});
					menu += subMenu + '</ul></li>';
				} else {
					menu += '<li><a class="norefresh"  url="' + m.ControlUrl + '" href="' + m.ControlUrl + (m.HomeUrl ? eval(m.HomeUrl) : "") + '" ><span>' + m.DisplayName + '</span></a></li>';
				}
			});
			$('.js-left-nav').append(menu);
			leftSelect();
		}
	});
	function GetDateStr(AddDayCount) {
		var dd = new Date();
		dd.setDate(dd.getDate() + AddDayCount);//获取AddDayCount天后的日期  
		var y = dd.getFullYear();
		var m = (dd.getMonth() + 1) < 10 ? "0" + (dd.getMonth() + 1) : (dd.getMonth() + 1);//获取当前月份的日期，不足10补0  
		var d = dd.getDate() < 10 ? "0" + dd.getDate() : dd.getDate();//获取当前几号，不足10补0  
		return y + "-" + m + "-" + d + " 00:00:00";
	}


	//左边菜单加选中状态
	function leftSelect() {
		var pathname = location.pathname;
		if (pathname != "/") {
			$('.js-left-nav .norefresh').filter(function () {
				return $(this).attr('url') == pathname;
			}).closest("li").addClass('active').parents('.menu-list').addClass('nav-active');
		}
	}

	$("body").delegate(".menu-list > a", "click", function () {
		var parent = jQuery(this).parent();
		var sub = parent.find('> ul');
		//if(!jQuery('body').hasClass('left-side-collapsed')) {
		if (sub.is(':visible')) {
			sub.slideUp(200, function () {//增加展开动画
				sub.css("display", "");
				parent.removeClass("nav-active")
				//mainContentHeightAdjust();
			});
			// sub.slideUp(200, function(){
			//    parent.removeClass('nav-active');
			//    jQuery('.main-content').css({height: ''});
			//    mainContentHeightAdjust();
			// });
		} else {
			visibleSubMenuClose();
			sub.slideDown(200, function () {
				sub.css("display", "");//解决某些情况下无法收起bug
				parent.addClass("nav-active")
				//mainContentHeightAdjust();
			});
		}
		//}
		return false;
	});
	$("body").delegate(".js-left-nav .norefresh", "click", function () {
		// if ($(this).closest("li").hasClass('active')){
		// 	return false;
		// }
		$('.js-left-nav .active').removeClass('active');
		$(this).closest("li").addClass('active');
		if (!jQuery('body').hasClass('left-side-collapsed') && $('.js-left-nav .nav-active').length && $(this).closest(".nav-active").length == 0) {
			visibleSubMenuClose();
			// mainContentHeightAdjust();
		}
		var url = $(this).attr('href');
		for(var i = 0;i<orderCountDownTimer_Arr.length;i++){
			clearInterval(orderCountDownTimer_Arr[i])
		}
		
		for(var j = 0;j<orderCountDownTime_Arr.length;j++){
			clearInterval(orderCountDownTime_Arr[j])
		}
        for(var x = 0;x<creditTimeOut_Arr.length;x++){
            clearTimeout(creditTimeOut_Arr[x])
        }
        for(var y = 0;y<orderTimeOut_Arr.length;y++){
            clearTimeout(orderTimeOut_Arr[y])
        }

		getpage(url, 1);
		//$('.wrapper').zload(url)
		return false;
	});
	window.onpopstate = function (e) {
		if (e.state) {
			$('.js-left-nav .active').removeClass('active');
			visibleSubMenuClose();
			leftSelect();
			// mainContentHeightAdjust();
			$('.wrapper').empty().append(e.state.html);
			execjs(e.state.html);
		}
	}



	function visibleSubMenuClose() {
		jQuery('.menu-list').each(function () {
			var t = jQuery(this);
			if (t.hasClass('nav-active')) {
				t.find('> ul').slideUp(200, function () {
					t.removeClass('nav-active');
				});
			}
		});
	}

	//function mainContentHeightAdjust() {
	// Adjust main content height
	// var docHeight = jQuery(document).height();
	// if(docHeight > jQuery('.main-content').height()){
	//   	jQuery('.main-content').height(docHeight);
	// }
	//}
	$('html').on("keyup",".onlyNumber",function(){
		this.value=this.value.replace(/\D/g,'');
	}).on("afterpaste",".onlyNumber",function(){
		this.value=this.value.replace(/\D/g,'');
	})
})

function getpage(url, is_condition) {
	var r = Math.random().toFixed(9)
	$.zget(url, {r : r}, function (result) {
		history.pushState({ html: result }, "现金分期", url);
		$('.wrapper')[0].innerHTML = result;// 不用Jq操作 是因为某些情况下，jq会执行请求到的页面的js语句， 所以采用原生防止执行
		execjs(result);
		if (is_condition) {
			SearchCondition(url)
		}
	});
	return false;
}

$('html').on('click', '.skip', function () {
	var href = $(this).attr('href');
	getpage(href);
	return false;
})

function getUrlString(str, name) {
	var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
	var r = str.substr(str.indexOf("?")).substr(1).match(reg);
	var context = "";
	if (r != null)
		context = r[2];
	reg = null;
	r = null;
	return context == null || context == "" || context == "undefined" ? "" : context;
}
function isBlank(obj) {
	return (!obj || $.trim(obj) === "");
}
function countdownanimation(seconds,msg,callback) {
	var modal = '<div class="countdown-modal"><div class="countdown-time">' + seconds + '</div></div>';
	$("body").append(modal);
	var countdown = setInterval(function () {
		var second = parseInt($(".countdown-time").html());
		if (second > 1) {
			$(".countdown-time").html(--second)
		} else if(second == 1){
			$(".countdown-time").html(msg);
		}else {
			$(".countdown-modal").remove()
			clearInterval(countdown);
			if(callback){callback()};
		}
	}, 1000)
}



function execjs(html) {
	// 第一步：匹配加载的页面中是否含有js
	var regDetectJs = /<script(.|\n)*?>(.|\n|\r\n)*?<\/script>/ig;
	var jsContained = html.match(regDetectJs);
	// 第二步：如果包含js，则一段一段的取出js再加载执行
	if (jsContained) {
		// 分段取出js正则
		var regGetJS = /<script(.|\n)*?>((.|\n|\r\n)*)?<\/script>/im;

		// 按顺序分段执行js
		var jsNums = jsContained.length;
		for (var i = 0; i < jsNums; i++) {
			var jsSection = jsContained[i].match(regGetJS);

			if (jsSection[2]) {
				if (window.execScript) {
					// 给IE的特殊待遇
					window.execScript(jsSection[2]);
				} else {
					// 给其他大部分浏览器用的
					window.eval(jsSection[2]);
				}
			}
		}
	}
}

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
