;
(function ($) {
    $.fn.extend({
        "my_page": function (form, source,callback) {
            var $this = this;
            //定义分页结构体
            var pageinfo = {
                url: $(this).data("url"),
                currentPage: $(this).attr("currentPage") * 1, // 当前页码
                pageCount: $(this).attr("pageCount") * 1,// 总页码
                pageSize: $(this).attr("pageSize") * 1
            };
            // var init = getUrlString(window.location.href,"start_time");
            // if(!isBlank(init)){
            //     $("#startDate").val(decodeURIComponent(init));//填充开始时间 
            // }
            var url = $(this).data("url");
            var submitFlag = true;

            var this_select = $(this).parent().siblings().find(".pageSize");
            this_select.val(pageinfo.pageSize).bind("change",function () {
                var pageSize = $(this).val();
                redirectTo(1,pageSize);
            });
            $(form).off("submit").on("submit", function () {
                // debugger;
                if (submitFlag) {
                    submitFlag = false;
                    var data = $(form).serialize();
                    if (data) {
                        url = url.split('?')[0] + '?' + data + "&pageSize=" + this_select.val();
                    }
                    
	                var r = Math.random().toFixed(9)
                    $.zget(url, {r:r}, function (result) {
                        if (!!source) {
                            $(source).html(result);
                        } else {
                            history.pushState({ html: result }, "现金分期", url);
                            $('.wrapper').html(result);
                            SearchCondition(url);
                        }
                        execjs(result);
                        submitFlag = true;
                    });
                }
                return false; // 阻止表单自动提交事件
            });
            if (pageinfo.pageCount < 2) {
                return false;
            }
            //初始起始页数、结束页数
            var start = 0,
                end = 5;
            // pageinfo.currentPage 当前页码
            // pageinfo.pageCount 总页码
            if (pageinfo.currentPage >= 5) {
                start = pageinfo.currentPage - 3;
            }
            if (pageinfo.currentPage >= 5) {
                end = pageinfo.currentPage + 2;
            } else if (pageinfo.currentPage < 5) {
                end = pageinfo.currentPage + 4;
            } else {
                end = pageinfo.pageCount;
            }

            var html = [];
            html.push("<ul class='page_content'>");
            if (pageinfo.currentPage < 5) {
                end = 5;
                start = 0;
            }
            if (pageinfo.currentPage > pageinfo.pageCount - 5) {
                end = pageinfo.pageCount;
                start = pageinfo.pageCount - 5;
            }
            if (pageinfo.pageCount > 5) {
                html.push("<li class='home_page'><a>首页</a></li>");
            }
            if (pageinfo.pageCount > 1) {
                html.push("<li class='page_prev'><a>上一页</a></li>");
            }

            if (pageinfo.currentPage == pageinfo.pageCount - 4 && pageinfo.currentPage != 2) {
                start = pageinfo.currentPage - 3;
                end = pageinfo.currentPage + 2;
            }

            if (pageinfo.pageCount <= 5) {
                start = 0;
                end = pageinfo.pageCount;
            }

            for (var i = start; i < end; i++) {
                if ((i + 1) == pageinfo.currentPage)
                    html.push("<li class='active'><a >" + (i + 1) + "</a></li>");
                else
                    html.push("<li class='nomal'><a >" + (i + 1) + "</a></li>");
            }
            if (pageinfo.pageCount > 1) {
                html.push("<li class='page_next'><a >下一页</a></li>");
            }
            if (pageinfo.pageCount > 5) {
                html.push("<li class='last_page'><a>末页</a></li>");
            }
            html.push("</ul>");
            html.push("<div style='margin-left: 10px;position: relative;top:-10px;display: inline-block'><span>共 " + pageinfo.pageCount + " 页</span>&nbsp;&nbsp;跳转到&nbsp;<input class='specifiedPage' style='width: 50px;height: 30px;text-align: center;border: 1px solid #ccc;border-radius: 2px' type='number' min='1' max='" + pageinfo.pageCount + "'>页&nbsp;&nbsp;<span class='pageSkip' style='display:inline-block;width: 54px;height: 30px;border: 1px solid #ccc;text-align: center;line-height: 30px;border-radius:3px;color: #65CEA7;cursor: pointer '>确定</span></div>");
            $this.html(html.join(""));

            //绑定数据处理函数
            $(".specifiedPage").keydown(function (event) {
                if (event.keyCode == 13) {
                    return false;
                }
            });
            $this.find(".pageSkip").bind('click', function () {
                if ($(".specifiedPage").val() < 1 || $(".specifiedPage").val() > pageinfo.pageCount) {
                    return;
                }
                redirectTo($(".specifiedPage").val());
            });
            $this.find(".last_page a").bind("click", function () {
                redirectTo(pageinfo.pageCount);
            });
            $this.find(".home_page a").bind("click", function () {
                redirectTo(1);
            });
            $this.find(".nomal a").bind("click", function () {
                redirectTo($(this).html());
            });
            $this.find(".page_prev a").bind("click", function () {
                if (pageinfo.currentPage == 1) {
                    $(".page_prev a").attr("disabled", true);
                } else {
                    $(".page_prev a").attr("disabled", false);
                    redirectTo(pageinfo.currentPage - 1);
                }
            });
            $this.find(".page_next a").bind("click", function () {
                if (pageinfo.currentPage == pageinfo.pageCount) {
                    $(".page_next a").attr("disabled", true);
                } else {
                    $(".page_next a").attr("disabled", false);
                    redirectTo(pageinfo.currentPage + 1);
                }
            });


            function redirectTo(page,pageSize) {
                var data = $(form).serialize();
                var url = pageinfo.url;
                var pageSize = pageSize || pageinfo.pageSize;
                // if(url.indexOf("?")==-1)
                //     url+="?";
                // else
                //     url+="&";
                if (data) {
                    url = url.split('?')[0] + '?' + data + "&page=" + page + "&pageSize=" + pageSize;
                } else {
                    if (url.indexOf("?") == -1)
                        url += "?" + "page=" + page + "&pageSize=" + pageSize;
                    else
                        url += "&" + "page=" + page + "&pageSize=" + pageSize;
                }
                var r = Math.random().toFixed(9)
                $.zget(url, {r:r}, function (result) {
                    if (!!source) {
                        $(source).html(result);
                        execjs(result);
                    } else if(callback){
                        callback(result);
                    }else {
                        history.pushState({ html: result }, "现金分期", url);
                        $('.wrapper').html(result);
                        execjs(result);
                        SearchCondition()
                    }
                })
            }
            return $this;
        }
    });
})(jQuery);

function SearchCondition(url) {
    //已选择的条件还原
    //空格变+号的处理
    var windowhref = '';

    if (!!url) {
        windowhref = url.replace(/\+/g, "%20");
    } else {
        windowhref = (window.location.href).replace(/\+/g, "%20");
    }
    var href = decodeURIComponent(windowhref).split("?");
    if (href[1]) {
        var cons = href[1].split("&");
        for (var i = 0; i < cons.length; i++) {
            var key = cons[i].split("=")[0];
            var value = cons[i].split("=")[1];
            if (value) {
                var par = $("[name='" + key + "']");
                if (par.length < 1) {
                    continue;
                }
                //var tagtype = par[0].tagName;//INPUT
                var tagtype = par[0].type; //type
                if (tagtype == "text" || tagtype == "tel") {
                    $(par).val(value);
                } else if (tagtype == "checkbox") {
                    par.attr("checked", true);
                } else if (tagtype == "select-one") {
                    par.children("option").each(function () {
                        if ($(this).val() == value || $(this).text() == value) {
                            $(this).attr("selected", true);
                        }
                    });
                } else if (tagtype == "select") {
                    par.children("option").each(function () {
                        if ($(this).val() == value || $(this).text() == value) {
                            $(this).attr("selected", true);
                        }
                    });
                } else if (tagtype == "radio") {
                    par.each(function () {
                        if ($(this).val() == value) {
                            $(this).attr("checked", true);
                            return false;
                        }
                    });
                }
            }
        }
    }
}

$(function () {
    SearchCondition()
    //var clickflag=true;
    //$("body").on("click",".search_btn",function(){
    //    var $_this=$(this);
    //    if(clickflag){
    //        clickflag=false;
    //        var url=$_this.data("url");
    //        var form=$_this.parents("form").serialize();
    //        $.zget(url,{form:form},function(result){
    //            console.log(result);
    //            $('.wrapper').html(result);
    //            execjs(result);
    //            clickflag=true;
    //        });
    //    }
    //});
});
